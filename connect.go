package gologix

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

func (client *Client) Connect() error {
	if client.ConnectionSize == 0 {
		client.ConnectionSize = 4000
		//client.ConnectionSize = 508
	}

	if client.Port == "" {
		client.Port = ":44818"
	}
	if client.VendorID == 0 {

		client.VendorID = 0x1776
	}
	if client.SerialNumber == 0 {
		client.SerialNumber = 42
	}

	// default path is backplane -> slot 0
	var err error
	if client.Path == nil {
		client.Path, err = Serialize(CIPPort{PortNo: 1}, CIPAddress(0))
		if err != nil {
			return fmt.Errorf("can't setup default path. %w", err)

		}
	}

	if ioi_cache == nil {
		ioi_cache = make(map[string]*tagIOI)
	}
	return client.connect()
}

func (client *Client) register_session() error {
	reg_msg := msgCIPRegister{}
	reg_msg.ProtocolVersion = 1
	reg_msg.OptionFlag = 0

	//binary.Write(conn.Conn, binary.LittleEndian, register_msg)
	//client.Send(cipCommandRegisterSession, reg_msg)
	resp_hdr, resp_data, err := client.send_recv_data(cipCommandRegisterSession, reg_msg)
	//resp_hdr, resp_data, err := client.recv_data(cipCommandRegisterSession, reg_msg)
	if err != nil {
		return fmt.Errorf("couldn't get connect response %w", err)
	}
	client.SessionHandle = resp_hdr.SessionHandle
	log.Printf("Session Handle %v", client.SessionHandle)
	_ = resp_data
	return nil

}

func (client *Client) keepalive() {
	if client.SocketTimeout == 0 {
		return
	}
	og_props, err := client.GetControllerPropList()
	if err != nil {
		log.Printf("keepalive prop list failed. %v", err)
		return
	}
	err = client.ListAllTags(0)
	if err != nil {
		log.Printf("keepalive list tags failed. %v", err)
		return
	}
	t := time.NewTicker(client.SocketTimeout / 4)
	for {
		select {
		case <-t.C:
			if client.Connected {
				new_props, err := client.GetControllerPropList()
				if err != nil {
					log.Printf("keepalive failed. %v", err)
					return
				}
				if !new_props.Match(og_props) {
					log.Printf("controller change detected. re-analyzing types.\n Was: %+v\n  Is: %+v", og_props, new_props)
					err := client.ListAllTags(0)
					if err != nil {
						log.Printf("keepalive list tags failed. %v", err)
						return
					}
					og_props = new_props

				}

			}
		case <-client.cancel_keepalive:
			t.Stop()
			return

		}
	}

}

// To connect we first send a register session command.
// based on the reply we get from that we send a forward open command.
func (client *Client) connect() error {
	if client.Connected {
		return nil
	}
	client.KnownTags = make(map[string]KnownTag)
	var err error
	client.conn, err = net.Dial("tcp", client.IPAddress+client.Port)
	if err != nil {
		return err
	}

	err = client.register_session()
	if err != nil {
		return err
	}

	if client.ConnectionSize == 0 {
		client.ConnectionSize = 4002
	}
	// we have to do something different for small connection sizes.
	fwd_open, err := client.NewForwardOpenLarge()
	if err != nil {
		return fmt.Errorf("couldn't create forward open. %w", err)
	}
	s := binary.Size(fwd_open)
	_ = s
	items0 := make([]cipItem, 2)
	items0[0] = cipItem{Header: cipItemHeader{ID: cipItem_Null}}
	items0[1] = fwd_open
	hdr, dat, err := client.send_recv_data(cipCommandSendRRData, MarshalItems(items0))
	if err != nil {
		return err
	}
	_ = hdr
	if hdr.Status != 0x00 {
		return fmt.Errorf("large Forward Open Failed. code %v", hdr.Status)
	}
	// header before items
	preitem := msgPreItemData{}
	err = binary.Read(dat, binary.LittleEndian, &preitem)
	if err != nil {
		return fmt.Errorf("problem reading items header from forward open req. %w", err)
	}

	items, err := ReadItems(dat)
	if err != nil {
		return fmt.Errorf("problem reading items from forward open req. %w", err)
	}

	fwopenresphdr := msgCIPMessageRouterResponse{}
	err = items[1].Unmarshal(&fwopenresphdr)
	if err != nil {
		return fmt.Errorf("error unmarshaling forward open response header. %w", err)
	}
	extended_status := make([]byte, fwopenresphdr.Status_Len*2)
	if fwopenresphdr.Status_Len != 0 {
		err = items[1].Unmarshal(&extended_status)
		if err != nil {
			return fmt.Errorf("error unmarshaling forward open response header extended status. %w", err)
		}
	}

	if fwopenresphdr.Status != 0x00 {
		log.Printf("bad status on large forward open header. got %x. Falling back to small forard open", fwopenresphdr.Status)
		client.ConnectionSize = 502

		// we have to do something different for small connection sizes.
		fwd_open, err := client.NewForwardOpenStandard()
		if err != nil {
			return fmt.Errorf("couldn't create forward open. %w", err)
		}
		s := binary.Size(fwd_open)
		_ = s
		items0 := make([]cipItem, 2)
		items0[0] = cipItem{Header: cipItemHeader{ID: cipItem_Null}}
		items0[1] = fwd_open
		hdr, dat, err := client.send_recv_data(cipCommandSendRRData, MarshalItems(items0))
		if err != nil {
			return err
		}
		_ = hdr
		if hdr.Status != 0x00 {
			return fmt.Errorf("small Forward Open Failed. code %v", hdr.Status)
		}
		// header before items
		preitem = msgPreItemData{}
		err = binary.Read(dat, binary.LittleEndian, &preitem)
		if err != nil {
			return fmt.Errorf("problem reading items header from forward open req. %w", err)
		}

		items, err = ReadItems(dat)
		if err != nil {
			return fmt.Errorf("problem reading items from forward open req. %w", err)
		}

		fwopenresphdr = msgCIPMessageRouterResponse{}
		err = items[1].Unmarshal(&fwopenresphdr)
		if err != nil {
			return fmt.Errorf("error unmarshaling forward open response header. %w", err)
		}
		extended_status = make([]byte, fwopenresphdr.Status_Len*2)
		if fwopenresphdr.Status_Len != 0 {
			err = items[1].Unmarshal(&extended_status)
			if err != nil {
				return fmt.Errorf("error unmarshaling forward open response header extended status. %w", err)
			}
		}

		if fwopenresphdr.Status != 0x00 {
			return fmt.Errorf("bad status on both forward opens. got %x", fwopenresphdr.Status)
		}

	}

	forwardopenresp := msgEIPForwardOpen_Reply{}
	err = items[1].Unmarshal(&forwardopenresp)
	if err != nil {
		return fmt.Errorf("error unmarshaling forward open response. %w", err)
	}
	log.Printf("ForwardOpen: %+v", forwardopenresp)
	client.OTNetworkConnectionID = forwardopenresp.OTConnectionID
	log.Printf("Connection ID: OT=%d, TO=%d", forwardopenresp.OTConnectionID, forwardopenresp.TOConnectionID)

	client.Connected = true

	client.cancel_keepalive = make(chan struct{})
	go client.keepalive()

	return nil

}

type msgPreItemData struct {
	Handle  uint32
	Timeout uint16
}

type msgCIPMessageRouterResponse struct {
	Service    CIPService
	Reserved   byte // always 0
	Status     byte // result status
	Status_Len byte // additional result word count - can be zero
}

type msgEIPForwardOpen_Reply struct {
	//Service        cipService
	//Reserved       byte // always 0
	//Status         byte // result status
	//Status_Len     byte // additional result word count - can be zero
	OTConnectionID uint32
	TOConnectionID uint32
	Unknown3       uint16
}

type msgEIPForwardClose struct {
	Service                CIPService
	PathSize               byte
	ClassType              byte
	Class                  byte
	InstanceType           byte
	Instance               byte
	Priority               byte
	TimeoutTicks           byte
	ConnectionSerialNumber uint16
	VendorID               uint16
	OriginatorSerialNumber uint32
	ConnPathSize           byte
}

// in this message T is for target and O is for originator so
// TO is target -> originator and OT is originator -> target
type msgEIPForwardOpen_Standard struct {
	Service                CIPService
	PathSize               byte
	ClassType              byte
	Class                  byte
	InstanceType           byte
	Instance               byte
	Priority               byte
	TimeoutTicks           byte
	OTConnectionID         uint32
	TOConnectionID         uint32
	ConnectionSerialNumber uint16
	VendorID               uint16
	OriginatorSerialNumber uint32
	Multiplier             uint32
	OTRPI                  uint32
	OTNetworkConnParams    uint16
	TORPI                  uint32
	TONetworkConnParams    uint16
	TransportTrigger       byte
	ConnPathSize           byte
}

type msgEIPForwardOpen_Large struct {
	// service
	Service CIPService
	// path
	PathSize     byte
	ClassType    CIPClassSize
	Class        byte
	InstanceType cipInstanceSize
	Instance     byte

	// service specific data
	Priority               byte
	TimeoutTicks           byte
	OTConnectionID         uint32
	TOConnectionID         uint32
	ConnectionSerialNumber uint16
	VendorID               uint16
	OriginatorSerialNumber uint32
	Multiplier             uint32
	OTRPI                  uint32
	OTNetworkConnParams    uint32
	TORPI                  uint32
	TONetworkConnParams    uint32
	TransportTrigger       byte
	PathLen                byte
}

func (client *Client) NewForwardOpenLarge() (cipItem, error) {
	item := cipItem{Header: cipItemHeader{ID: cipItem_UnconnectedData}}
	var msg msgEIPForwardOpen_Large
	rand.Seed(time.Now().Unix())

	p, err := Serialize(
		client.Path,
		//CIPPort{PortNo: 1}, CIPAddress(0),
		cipObject_MessageRouter, CIPInstance(1))
	if err != nil {
		return item, fmt.Errorf("couldn't build path. %w", err)
	}

	client.ConnectionSerialNumber = uint16(rand.Uint32())
	ConnectionParams := uint32(0x4200)
	ConnectionParams = ConnectionParams << 16 // for long packet
	ConnectionParams += uint32(client.ConnectionSize)

	msg.Service = cipService_LargeForwardOpen
	// this next section is the path
	msg.PathSize = 0x02 // length in words
	msg.ClassType = cipClass_8bit
	msg.Class = byte(cipObject_ConnectionManager)
	msg.InstanceType = cipInstance_8bit
	msg.Instance = 0x01
	// end of path
	msg.Priority = 0x0A     // 0x0A means normal multiplier (about 1 second?)
	msg.TimeoutTicks = 0x0E // number of "priority" ticks (0x0E = 14 * Priority = ~1 sec => ~ 14 seconds.)
	//msg.OTConnectionID = 0x05318008
	msg.OTConnectionID = rand.Uint32() //0x20000002
	msg.TOConnectionID = rand.Uint32()
	msg.ConnectionSerialNumber = client.ConnectionSerialNumber
	msg.VendorID = client.VendorID
	msg.OriginatorSerialNumber = client.SerialNumber
	msg.Multiplier = 0x03
	msg.OTRPI = 0x00201234
	msg.OTNetworkConnParams = ConnectionParams
	msg.TORPI = 0x00204001
	msg.TONetworkConnParams = ConnectionParams
	msg.TransportTrigger = 0xA3
	msg.PathLen = byte(p.Len() / 2)
	item.Marshal(msg)
	item.Marshal(p.Bytes())

	return item, nil
}

type msgCIPRegister struct {
	ProtocolVersion uint16
	OptionFlag      uint16
}

func (client *Client) NewForwardOpenStandard() (cipItem, error) {
	item := cipItem{Header: cipItemHeader{ID: cipItem_UnconnectedData}}
	var msg msgEIPForwardOpen_Standard

	rand.Seed(time.Now().Unix())

	p, err := Serialize(
		client.Path,
		cipObject_MessageRouter, CIPInstance(1))
	if err != nil {
		return item, fmt.Errorf("couldn't build path. %w", err)
	}

	client.ConnectionSerialNumber = uint16(rand.Uint32())
	ConnectionParams := uint16(0x43F6)

	msg.Service = cipService_ForwardOpen
	// this next section is the path
	msg.PathSize = 0x02 // length in words
	msg.ClassType = byte(cipClass_8bit)
	msg.Class = byte(cipObject_ConnectionManager)
	msg.InstanceType = byte(cipInstance_8bit)
	msg.Instance = 0x01
	// end of path
	msg.Priority = 0x07     // 0x0A means normal multiplier (about 1 second?)
	msg.TimeoutTicks = 0xE9 // number of "priority" ticks (0x0E = 14 * Priority = ~1 sec => ~ 14 seconds.)
	//msg.OTConnectionID = 0x05318008
	msg.OTConnectionID = 0 //0x20000002
	msg.TOConnectionID = rand.Uint32()
	msg.ConnectionSerialNumber = client.ConnectionSerialNumber
	msg.VendorID = client.VendorID
	msg.OriginatorSerialNumber = client.SerialNumber
	msg.Multiplier = 0x00
	msg.OTRPI = 0x007270E0
	msg.OTNetworkConnParams = uint16(ConnectionParams)
	msg.TORPI = 0x007270E0
	msg.TONetworkConnParams = uint16(ConnectionParams)
	msg.TransportTrigger = 0xA3
	msg.ConnPathSize = byte(p.Len() / 2)
	item.Marshal(msg)
	item.Marshal(p.Bytes())

	return item, nil
}

type msgEIPForwardOpen_Standard_Reply struct {
	Service                CIPService
	Reserved               byte
	Status                 byte
	StatusLen              byte
	OTConnectionID         uint32
	TOConnectionID         uint32
	ConnectionSerialNumber uint16
	VendorID               uint16
	OriginatorSerialNumber uint32
	OTAPI                  uint32
	TOAPI                  uint32
	ReplySize              byte
	Reserved2              byte
}
