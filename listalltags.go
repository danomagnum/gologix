package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

type EmbeddedMessage struct {
	Size       uint16
	Service    CIPService
	PathLength byte
	/*
		Path          []byte
		Data          [4]uint16
		RoutePathSize byte
		Reserved      byte
		PathSegment   uint16
	*/
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
	TypeInfo   byte
	Dimension1 uint32
	Dimension2 uint32
	Dimension3 uint32
}

type ListInstanceHeader struct {
	Unknown uint16
	Status  uint16
}

func (plc *PLC) ListAllTags(start_instance uint32) error {
	plc.readSequencer += 1
	fmt.Printf("readall for %v", start_instance)

	// have to start at 1.
	if start_instance == 0 {
		start_instance = 1
	}

	reqitems := make([]CIPItem, 2)
	reqitems[0] = CIPItem{Header: CIPItemHeader{ID: CIPItem_Null}}

	p := Paths(
		MarshalPathLogical(LogicalTypeClassID, 0x6B, true),
		MarshalPathLogical(LogicalTypeInstanceID, start_instance, true),
	)

	readmsg := ReaddAllData{
		//Sequence:    plc.readSequencer,
		Service:     CIPService_FragRead,
		PathLength:  2,
		RequestPath: [4]byte{0x20, 0x06, 0x24, 0x01},
		Timeout:     0,
		Message: EmbeddedMessage{
			Size:       14,
			Service:    CIPService_GetInstanceAttributeList,
			PathLength: byte(len(p) / 2),
			/*Path:          p,
			Data:          [4]uint16{3, 1, 2, 8},
			RoutePathSize: 1,
			Reserved:      0,
			PathSegment:   1,
			*/
		},
	}

	reqitems[1] = NewItem(CIPItem_UnconnectedData, readmsg)
	reqitems[1].Marshal(p)
	reqitems[1].Marshal([4]uint16{3, 1, 2, 8})
	reqitems[1].Marshal(byte(1))
	reqitems[1].Marshal(byte(0))
	reqitems[1].Marshal(uint16(1))

	plc.Send(CIPCommandSendRRData, MarshalItems(reqitems))
	hdr, data, err := plc.recv_data()
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

		// the end of the tagname has to be aligned on a 16 bit word
		//tagname_alignment := tag_hdr.NameLength % 2
		//if tagname_alignment != 0 {
		//data2.Next(int(tagname_alignment))
		//}
		binary.Read(data2, binary.LittleEndian, tag_ftr)

		log.Printf("Tag: '%s' Instance: %d Type: %s/%d[%d,%d,%d]",
			tag_name,
			tag_hdr.InstanceID,
			tag_ftr.Type,
			tag_ftr.TypeInfo,
			tag_ftr.Dimension1,
			tag_ftr.Dimension2,
			tag_ftr.Dimension3,
		)
		start_instance = tag_hdr.InstanceID

	}
	log.Printf("Status: %v", hdr.Status)
	log.Printf("item1 Status: %v", data_hdr.Status)
	// eventually keep going past 200
	//if tries < 3 {
	//tries++
	//return plc.ReadAll(start_instance)
	//}
	if data_hdr.Status == 6 && start_instance < 200 {
		plc.ListAllTags(start_instance)

	}

	return nil
}
