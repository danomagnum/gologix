package main

import (
	"fmt"
	"log"
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

func (plc *PLC) read_single(tag string, datatype CIPType, elements uint16) error {
	ioi := BuildIOI(tag, datatype)

	ioi_header := CIPIOIHeader{
		Service: CIPService_FragRead,
		Size:    byte(len(ioi.Buffer) / 2),
	}
	ioi_footer := CIPIOIFooter{
		Elements: 1,
		Offset:   0,
	}
	// I think the read message consists of two items because item 1 says
	// "the next item has the details" and item 2 says "the details are
	// an ioi of this size". This is just speculation though.
	cip_header := CIPCommonPacketConnected{}
	cip_header.InterfaceHandle = 0
	cip_header.Timeout = 0
	cip_header.ItemCount = 2
	cip_header.Item1ID = 0xA1
	cip_header.Item1Length = 0x04
	cip_header.Item1 = plc.conn.OTNetworkConnectionID
	cip_header.Item2ID = 0xB1
	cip_header.Item2Length = uint16(SizeOf(ioi_header, ioi.Buffer, ioi_footer)) + 2
	log.Printf("item 2 length %v", cip_header.Item2Length)
	cip_header.Sequence = plc.conn.SequenceCounter

	plc.conn.Send(CIPCommandSendUnitData, cip_header, ioi_header, ioi.Buffer, ioi_footer)
	hdr, data, err := plc.conn.recv_data()
	if err != nil {
		return err
	}
	fmt.Printf("single read complete.\n Got header\n %v\ndata\n %v\n", hdr, data)
	return nil
}
