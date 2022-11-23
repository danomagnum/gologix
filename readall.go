package main

type EmbeddedMessage struct {
	Size          uint16
	Service       CIPService
	PathLength    byte
	Path          [4]byte
	Data          [4]uint16
	RoutePathSize byte
	Reserved      byte
	PathSegment   uint16
}

type ReaddAllData struct {
	//Sequence    uint16
	Service     CIPService
	PathLength  byte
	RequestPath [4]byte
	Timeout     uint16
	Message     EmbeddedMessage
}

func (plc *PLC) ReadAll(start_instance byte) error {
	plc.readSequencer += 1

	// have to start at 1.
	if start_instance == 0 {
		start_instance = 1
	}

	reqitems := make([]CIPItem, 2)
	reqitems[0] = CIPItem{Header: CIPItemHeader{ID: CIPItem_Null}}

	readmsg := ReaddAllData{
		//Sequence:    plc.readSequencer,
		Service:     CIPService_FragRead,
		PathLength:  2,
		RequestPath: [4]byte{0x20, 0x06, 0x24, 0x01},
		Timeout:     0,
		Message: EmbeddedMessage{
			Size:          14,
			Service:       CIPService_GetInstanceAttributeList,
			PathLength:    2,
			Path:          [4]byte{0x20, 0x6b, 0x24, start_instance},
			Data:          [4]uint16{3, 1, 2, 8},
			RoutePathSize: 1,
			Reserved:      0,
			PathSegment:   1,
		},
	}

	reqitems[1] = NewItem(CIPItem_UnconnectedData, readmsg)

	plc.conn.Send(CIPCommandSendRRData, BuildItemsBytes(reqitems))
	hdr, data, err := plc.conn.recv_data()
	if err != nil {
		return err
	}
	_ = hdr
	_ = data

	return nil
}
