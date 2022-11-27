package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
)

// Send takes the command followed by all the structures that need
// concatenated together.  It builds the header, puts the packet together,
// and then sends it.
func (conn *PLC) Send(cmd CIPCommand, msgs ...any) error {
	// calculate size of all message parts
	size := 0
	for _, msg := range msgs {
		size += binary.Size(msg)
	}
	// build header based on size
	hdr := conn.BuildHeader(cmd, size)

	// initialize a buffer and add the header to it.
	// the 24 is from the header size
	b := make([]byte, 0, size+24)
	buf := bytes.NewBuffer(b)
	binary.Write(buf, binary.LittleEndian, hdr)

	// add all message components to the buffer.
	for _, msg := range msgs {
		binary.Write(buf, binary.LittleEndian, msg)
	}

	b = buf.Bytes()
	// write the packet buffer to the tcp connection
	written := 0
	for written < len(b) {
		n, err := conn.Conn.Write(b[written:])
		if err != nil {
			return err
		}
		written += n
	}
	//log.Printf("Sent: %v", b)
	return nil

}

// recv_data reads the header and then the number of words it specifies.
func (conn *PLC) recv_data() (EIPHeader, *bytes.Reader, error) {

	hdr := EIPHeader{}
	var err error
	err = binary.Read(conn.Conn, binary.LittleEndian, &hdr)
	if err != nil {
		return hdr, nil, err
	}
	//log.Printf("Header: %v", hdr)
	//data_size := hdr.Length * 2
	data_size := hdr.Length
	data := make([]byte, data_size)
	if data_size > 0 {
		err = binary.Read(conn.Conn, binary.LittleEndian, &data)
	}
	//log.Printf("Buffer: %v", data)
	buf := bytes.NewReader(data)
	return hdr, buf, err

}

func (conn *PLC) BuildHeader(cmd CIPCommand, size int) (hdr EIPHeader) {

	conn.SequenceCounter++

	hdr.Command = uint16(cmd)
	//hdr.Command = 0x0070
	hdr.Length = uint16(size)
	hdr.SessionHandle = conn.SessionHandle
	hdr.Status = 0
	hdr.Context = conn.Context
	hdr.Options = 0

	return

}

const CIP_Port = ":44818"
const CIP_VendorID = 0x1776

// To connect we first send a register session command.
// based on the reply we get from that we send a forward open command.
func (conn *PLC) connect(ip string) error {
	if conn.Connected {
		return nil
	}
	var err error
	conn.Conn, err = net.Dial("tcp", ip+CIP_Port)
	if err != nil {
		return err
	}

	reg_msg := CIPMessage_Register{}
	reg_msg.ProtocolVersion = 1
	reg_msg.OptionFlag = 0

	err = conn.Send(CIPCommandRegisterSession, reg_msg) // 0x65 is register session
	if err != nil {
		log.Panicf("Couldn't send connect req %v", err)
	}
	//binary.Write(conn.Conn, binary.LittleEndian, register_msg)
	resp_hdr, resp_data, err := conn.recv_data()
	if err != nil {
		return err
	}
	conn.SessionHandle = resp_hdr.SessionHandle
	log.Printf("Session Handle %v", conn.SessionHandle)
	_ = resp_data

	conn.ConnectionSize = 4002
	// we have to do something different for small connection sizes.
	fwd_open := conn.build_forward_open_large()
	s := binary.Size(fwd_open)
	_ = s
	items0 := make([]CIPItem, 2)
	items0[0] = CIPItem{Header: CIPItemHeader{ID: CIPItem_Null}}
	items0[1] = NewItem(CIPItem_UnconnectedData, fwd_open)
	err = conn.Send(CIPCommandSendRRData, BuildItemsBytes(items0))
	if err != nil {
		return err
	}
	hdr, dat, err := conn.recv_data()
	if err != nil {
		return err
	}
	_ = hdr
	if hdr.Status == 0x01 {
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
	err = binary.Read(&items[1], binary.LittleEndian, &forwardopenresp)
	//err = binary.Read(dat, binary.LittleEndian, &forwardopenresp)
	if err != nil {
		log.Printf("Error Reading. %v", err)
	}
	log.Printf("ForwardOpen: %+v", forwardopenresp)
	conn.OTNetworkConnectionID = forwardopenresp.OTConnectionID

	conn.Connected = true
	return nil

}

// to disconect we send two items - a null item and an unconnected data item for the unregister service
func (conn *PLC) Disconnect() error {
	if !conn.Connected {
		return nil
	}
	conn.Connected = false
	var err error

	items := make([]CIPItem, 2)
	items[0] = CIPItem{} // null item

	reg_msg := CIPMessage_UnRegister{
		Service:                CIPService_ForwardClose,
		CipPathSize:            0x02,
		ClassType:              CIPClass_8bit,
		Class:                  0x06,
		InstanceType:           CIPInstance_8bit,
		Instance:               0x01,
		Priority:               0x0A,
		TimeoutTicks:           0x0E,
		ConnectionSerialNumber: conn.ConnectionSerialNumber,
		VendorID:               CIP_VendorID,
		OriginatorSerialNumber: CIP_SerialNumber,
		PathSize:               3,                                           // 16 bit words
		Path:                   [6]byte{0x01, 0x00, 0x20, 0x02, 0x24, 0x01}, // TODO: generate paths automatically
	}

	items[1] = NewItem(CIPItem_UnconnectedData, reg_msg)

	err = conn.Send(CIPCommandSendRRData, BuildItemsBytes(items)) // 0x65 is register session
	if err != nil {
		log.Panicf("Couldn't send unconnect req %v", err)
	}
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

const CIP_SerialNumber = 42

func (conn *PLC) build_forward_open_large() CIPItem {
	item := CIPItem{Header: CIPItemHeader{ID: CIPItem_UnconnectedData}}
	var msg EIPForwardOpen_Large

	p := Paths(
		PathPortBuild([]byte{0x00}, 1, true),
		PathLogicalBuild(LogicalTypeClassID, uint32(CIPObject_MessageRouter), true),
		PathLogicalBuild(LogicalTypeInstanceID, 0x01, true),
	)

	conn.ConnectionSerialNumber = uint16(rand.Uint32())
	ConnectionParams := uint32(0x4200)
	ConnectionParams = ConnectionParams << 16 // for long packet
	ConnectionParams += uint32(conn.ConnectionSize)

	msg.Service = CIPService_LargeForwardOpen
	msg.PathSize = 0x02
	msg.ClassType = 0x20
	msg.Class = 0x06
	msg.InstanceType = 0x24
	msg.Instance = 0x01
	msg.Priority = 0x0A
	msg.TimeoutTicks = 0x0E
	//msg.OTConnectionID = 0x05318008
	msg.OTConnectionID = 0x20000002
	msg.TOConnectionID = rand.Uint32()
	msg.ConnectionSerialNumber = conn.ConnectionSerialNumber
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
