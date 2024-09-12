package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"strings"
)

// the gist here is that we want to do a fragmented read (since there will undoubtedly be more than one packet's worth)
// of the instance attribute list of the symbol objects.
//
// see 1756-PM020H-EN-P March 2022 page 39
// also see https://forums.mrclient.com/index.php?/topic/40626-reading-and-writing-io-tags-in-plc/
func (client *Client) ListSubTags(Program *KnownProgram, start_instance uint32) ([]KnownTag, error) {

	new_kts := make([]KnownTag, 0, 100)
	client.Logger.Debug("readall", "start id", start_instance)

	// have to start at 1.
	if start_instance == 0 {
		start_instance = 1
	}

	reqitems := make([]CIPItem, 2)
	//reqitems[0] = cipItem{Header: cipItemHeader{ID: cipItem_Null}}
	reqitems[0] = newItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	p, err := Serialize(
		Program.Bytes(),
		CipObject_Symbol, CIPInstance(start_instance),
	)
	if err != nil {
		return new_kts, fmt.Errorf("couldn't build path. %w", err)
	}

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: uint16(sequencer()),
		Service:       CIPService_GetInstanceAttributeList,
		PathLength:    byte(p.Len() / 2),
	}

	reqitems[1] = newItem(cipItem_ConnectedData, readmsg)
	reqitems[1].Serialize(p.Bytes())
	number_of_attr_to_receive := 3
	attr1_symbol_name := 1
	attr2_symbol_type := 2
	attr8_arraydims := 8
	reqitems[1].Serialize([4]uint16{uint16(number_of_attr_to_receive), uint16(attr1_symbol_name), uint16(attr2_symbol_type), uint16(attr8_arraydims)})

	itemdata, err := serializeItems(reqitems)
	if err != nil {
		return nil, fmt.Errorf("problem serializing items: %w", err)
	}
	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, itemdata)
	if err != nil {
		return new_kts, err
	}
	_ = hdr
	_ = data

	// first six bytes are zero.
	padding := make([]byte, 6)
	_, err = data.Read(padding)
	if err != nil {
		return nil, fmt.Errorf("problem reading padding. %w", err)
	}

	resp_items, err := readItems(data)
	if err != nil {
		return new_kts, fmt.Errorf("couldn't parse items. %w", err)
	}

	// get ready to read tag info from item 1 data
	data2 := bytes.NewBuffer(resp_items[1].Data)
	//data2.Next(4)
	data_hdr := msgListInstanceHeader{}
	err = binary.Read(data2, binary.LittleEndian, &data_hdr)
	if err != nil {
		return nil, fmt.Errorf("problem reading tag header. %w", err)
	}

	tag_hdr := new(msgtagResultDataHeader)
	tag_ftr := new(TagInfo)
	for data2.Len() > 0 {

		err = binary.Read(data2, binary.LittleEndian, tag_hdr)
		if err != nil {
			return nil, fmt.Errorf("problem reading tag header. %w", err)
		}
		newtag_bytes := make([]byte, tag_hdr.NameLength)
		err = binary.Read(data2, binary.LittleEndian, &newtag_bytes)
		if err != nil {
			return nil, fmt.Errorf("problem reading tag header. %w", err)
		}
		newtag_name := fmt.Sprintf("program:%s.%s", Program.Name, string(newtag_bytes))

		// the end of the tagname has to be aligned on a 16 bit word
		//tagname_alignment := tag_hdr.NameLength % 2
		//if tagname_alignment != 0 {
		//data2.Next(int(tagname_alignment))
		//}
		err = binary.Read(data2, binary.LittleEndian, tag_ftr)
		if err != nil {
			return nil, fmt.Errorf("problem reading tag footer. %w", err)
		}

		if strings.ToLower(newtag_name) == "program:gologix_tests.testtimer" {
			log.Printf("found it.")
		}

		kt := KnownTag{
			Name:     newtag_name,
			Info:     *tag_ftr,
			Instance: CIPInstance(tag_hdr.InstanceID),
			Parent:   Program,
		}
		if tag_ftr.Dimension3 != 0 {
			kt.Array_Order = make([]int, 3)
			kt.Array_Order[2] = int(tag_ftr.Dimension3)
			kt.Array_Order[1] = int(tag_ftr.Dimension2)
			kt.Array_Order[0] = int(tag_ftr.Dimension1)
		} else if tag_ftr.Dimension2 != 0 {
			kt.Array_Order = make([]int, 2)
			kt.Array_Order[1] = int(tag_ftr.Dimension2)
			kt.Array_Order[0] = int(tag_ftr.Dimension1)
		} else if tag_ftr.Dimension1 != 0 {
			kt.Array_Order = make([]int, 1)
			kt.Array_Order[0] = int(tag_ftr.Dimension1)
		} else {
			kt.Array_Order = make([]int, 0)
		}
		if !isValidTag(string(newtag_bytes), *tag_ftr) {
			continue
		}
		client.KnownTags[strings.ToLower(newtag_name)] = kt
		new_kts = append(new_kts, kt)

		start_instance = tag_hdr.InstanceID

	}

	if data_hdr.Status == 6 {
		_, err = client.ListSubTags(Program, start_instance)
		if err != nil {
			return new_kts, fmt.Errorf("problem listing subtags. %w", err)
		}
	}

	return new_kts, nil
}
