package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
)

func (plc *PLC) Connect() error {
	if plc.Size == 0 {
		plc.Size = 508
	}

	if ioi_cache == nil {
		ioi_cache = make(map[string]*IOI)
	}
	return plc.connect(plc.IPAddress)
}

func (plc *PLC) register_session() error {
	reg_msg := CIPMessage_Register{}
	reg_msg.ProtocolVersion = 1
	reg_msg.OptionFlag = 0

	err := plc.Send(CIPCommandRegisterSession, reg_msg) // 0x65 is register session
	if err != nil {
		return fmt.Errorf("couldn't send connect req %w", err)
	}
	//binary.Write(conn.Conn, binary.LittleEndian, register_msg)
	resp_hdr, resp_data, err := plc.recv_data()
	if err != nil {
		return fmt.Errorf("couldn't get connect response %w", err)
	}
	plc.SessionHandle = resp_hdr.SessionHandle
	log.Printf("Session Handle %v", plc.SessionHandle)
	_ = resp_data
	return nil

}

// To connect we first send a register session command.
// based on the reply we get from that we send a forward open command.
func (plc *PLC) connect(ip string) error {
	if plc.Connected {
		return nil
	}
	var err error
	plc.Conn, err = net.Dial("tcp", ip+CIP_Port)
	if err != nil {
		return err
	}

	err = plc.register_session()
	if err != nil {
		return err
	}

	plc.ConnectionSize = 4002
	// we have to do something different for small connection sizes.
	fwd_open := plc.build_forward_open_large()
	s := binary.Size(fwd_open)
	_ = s
	items0 := make([]CIPItem, 2)
	items0[0] = CIPItem{Header: CIPItemHeader{ID: CIPItem_Null}}
	items0[1] = fwd_open
	err = plc.Send(CIPCommandSendRRData, BuildItemsBytes(items0))
	if err != nil {
		return err
	}
	hdr, dat, err := plc.recv_data()
	if err != nil {
		return err
	}
	_ = hdr
	if hdr.Status != 0x00 {
		return fmt.Errorf("large Forward Open Failed. code %v", hdr.Status)
	}
	// header before items
	preitem := PreItemData{}
	err = binary.Read(dat, binary.LittleEndian, &preitem)
	if err != nil {
		return fmt.Errorf("problem reading items header from forward open req. %w", err)
	}

	items, err := ReadItems(dat)
	if err != nil {
		return fmt.Errorf("problem reading items from forward open req. %w", err)
	}

	forwardopenresp := EIPForwardOpen_Reply{}
	err = items[1].Unmarshal(&forwardopenresp)
	if err != nil {
		return fmt.Errorf("error unmarshaling forward open response. %w", err)
	}
	log.Printf("ForwardOpen: %+v", forwardopenresp)
	plc.OTNetworkConnectionID = forwardopenresp.OTConnectionID
	log.Printf("Connection ID: OT=%d, TO=%d", forwardopenresp.OTConnectionID, forwardopenresp.TOConnectionID)

	plc.Connected = true
	return nil

}

type PreItemData struct {
	Handle  uint32
	Timeout uint16
}

type EIPForwardOpen_Reply struct {
	Service        CIPService
	Unknown2       [3]byte
	OTConnectionID uint32
	TOConnectionID uint32
	Unknown3       uint16
}

// in this message T is for target and O is for originator so
// TO is target -> originator and OT is originator -> target
type EIPForwardOpen_Standard struct {
	Service                byte
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
}

type EIPForwardOpen_Large struct {
	Service   CIPService
	PathSize  byte
	ClassType byte
	Class     byte

	InstanceType byte
	Instance     byte
	Priority     byte
	TimeoutTicks byte

	OTConnectionID         uint32
	TOConnectionID         uint32
	ConnectionSerialNumber uint16
	VendorID               uint16

	OriginatorSerialNumber uint32
	Multiplier             uint32
	OTRPI                  uint32
	OTNetworkConnParams    uint32

	TORPI               uint32
	TONetworkConnParams uint32
	TransportTrigger    byte
	PathLen             byte
}

func (plc *PLC) build_forward_open_large() CIPItem {
	item := CIPItem{Header: CIPItemHeader{ID: CIPItem_UnconnectedData}}
	var msg EIPForwardOpen_Large

	p := Paths(
		PathPortBuild([]byte{0x00}, 1, true),
		PathLogicalBuild(LogicalTypeClassID, uint32(CIPObject_MessageRouter), true),
		PathLogicalBuild(LogicalTypeInstanceID, 0x01, true),
	)

	plc.ConnectionSerialNumber = uint16(rand.Uint32())
	ConnectionParams := uint32(0x4200)
	ConnectionParams = ConnectionParams << 16 // for long packet
	ConnectionParams += uint32(plc.ConnectionSize)

	msg.Service = CIPService_LargeForwardOpen
	msg.PathSize = 0x02
	msg.ClassType = 0x20
	msg.Class = 0x06
	msg.InstanceType = 0x24
	msg.Instance = 0x01
	msg.Priority = 0x0A
	msg.TimeoutTicks = 0x0E
	//msg.OTConnectionID = 0x05318008
	msg.OTConnectionID = rand.Uint32() //0x20000002
	msg.TOConnectionID = rand.Uint32()
	msg.ConnectionSerialNumber = plc.ConnectionSerialNumber
	msg.VendorID = CIP_VendorID
	msg.OriginatorSerialNumber = CIP_SerialNumber
	msg.Multiplier = 0x03
	msg.OTRPI = 0x00201234
	msg.OTNetworkConnParams = ConnectionParams
	msg.TORPI = 0x00204001
	msg.TONetworkConnParams = ConnectionParams
	msg.TransportTrigger = 0xA3
	// The path is formatted like this.
	// byte 0: number of 16 bit words
	// byte 1: 000. .... path segment type (port segment = 0)
	// byte 1: ...0 .... extended link address (0 = false)
	// byte 1: .... 0001 port (backplane = 1)
	// byte 2: n/a
	// byte 3: 001. .... path segment type (logical segment = 1)
	// byte 3: ...0 00.. logical segment type class ID (0)
	// byte 3: .... ..00 logical segment format: 8-bit (0)
	// byte 4: path segment 0x20
	// byte 5: 001. .... path segment type (logical segment = 1)
	// byte 5: ...0 01.. logical segment type: Instance ID = 1
	// byte 5: .... ..00 logical segment format: 8-bit (0)
	// byte 6: path segment instance 0x01
	//msg.Path = [6]byte{0x01, 0x00, 0x20, 0x02, 0x24, 0x01} // TODO: build path automatically
	// the 0x00 here is the slot number of the controller?
	msg.PathLen = byte(len(p) / 2)
	//msg.Path = [6]byte{0x01, 0x00, 0x20, 0x02, 0x24, 0x01} // TODO: build path automatically
	item.Marshal(msg)
	item.Marshal(p)

	return item
}

type CIPMessage_Register struct {
	ProtocolVersion uint16
	OptionFlag      uint16
}