package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

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

type tagResultDataHeader struct {
	InstanceID uint32
	NameLength uint16
}

type tagResultDataFooter struct {
	Type       CIPType
	Dimension1 uint32
	Dimension2 uint32
	Dimension3 uint32
}

type ListInstanceHeader struct {
	Unknown uint16
	Status  uint16
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
	//data_hdr := ListInstanceHeader{}
	//binary.Read(data, binary.LittleEndian, &data_hdr)
	padding := make([]byte, 6)
	data.Read(padding)

	resp_items, err := ReadItems(data)
	if err != nil {
		log.Panic("Couldn't parse items")
	}
	data2 := bytes.NewBuffer(resp_items[1].Data)
	//data2.Next(4)
	data_hdr := ListInstanceHeader{}
	binary.Read(data2, binary.LittleEndian, &data_hdr)

	tag_hdr := new(tagResultDataHeader)
	tag_ftr := new(tagResultDataFooter)
	for data2.Len() > 0 {
		binary.Read(data2, binary.LittleEndian, tag_hdr)
		tag_name := make([]byte, tag_hdr.NameLength)
		binary.Read(data2, binary.LittleEndian, &tag_name)
		if tag_hdr.NameLength%2 == 1 {
			data2.Next(1)
		}
		binary.Read(data2, binary.LittleEndian, tag_ftr)

		log.Printf("Tag: '%s' Instance: %d Type: %s[%d,%d,%d]",
			tag_name,
			tag_hdr.InstanceID,
			tag_ftr.Type,
			tag_ftr.Dimension1,
			tag_ftr.Dimension2,
			tag_ftr.Dimension3,
		)
		start_instance = byte(tag_hdr.InstanceID)

	}
	log.Printf("Status: %v", hdr.Status)
	log.Printf("item1 Status: %v", data_hdr.Status)
	// eventually keep going past 200
	if data_hdr.Status == 6 && start_instance < 200 {
		// continue
		//plc.ReadAll(start_instance)

	}

	return nil
}
