package main

import (
	"fmt"
	"time"
)

type PLC struct {
	IPAddress     string
	ProcessorSlot int
	SocketTimeout time.Duration
	// Route
	conn Connection
}

func (plc *PLC) Read_Single(tag string) []byte {
	return nil
}

func (plc *PLC) Connect() error {
	return plc.conn.Connect(plc.IPAddress)
}

type IOI_Header struct {
	InterfaceHandle uint32
	Timeout         uint16
	ItemCount       uint16
	Item1ID         uint16
	Item1Length     uint16
	Item1           uint32
	Item2ID         uint16
	Item2Length     uint16
	Sequence        uint16
}

func (plc *PLC) read_single(tag string, datatype CIPType, elements uint16) error {
	ioi := BuildIOI(tag, datatype)

	// I think the read message consists of two items because item 1 says
	// "the next item has the details" and item 2 says "the details are
	// an ioi of this size". This is just speculation though.
	ioi_header := IOI_Header{}
	ioi_header.InterfaceHandle = 0
	ioi_header.Timeout = 0
	ioi_header.ItemCount = 2
	ioi_header.Item1ID = 0xA1
	ioi_header.Item1Length = 0x04
	ioi_header.Item1 = plc.conn.OTNetworkConnectionID
	ioi_header.Item2ID = 0xB1
	ioi_header.Item2Length = uint16(len(ioi.Buffer) / 2)
	ioi_header.Sequence = plc.conn.SequenceCounter

	plc.conn.Send(CIPService_Read, ioi_header, ioi)
	hdr, data, err := plc.conn.recv_data()
	if err != nil {
		return err
	}
	fmt.Printf("single read complete.\n Got header\n %v\ndata\n %v\n", hdr, data)
	return nil
}
