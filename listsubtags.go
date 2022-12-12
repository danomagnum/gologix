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
func (client *Client) ListSubTags(roottag string, start_instance uint32) error {
	if verbose {
		log.Printf("readall for %v", start_instance)
	}

	// have to start at 1.
	if start_instance == 0 {
		start_instance = 1
	}

	ioi, err := client.newIOI(roottag, 16)
	if err != nil {
		return fmt.Errorf("bad IOI gen. %w", err)
	}

	reqitems := make([]CIPItem, 2)
	//reqitems[0] = CIPItem{Header: CIPItemHeader{ID: CIPItem_Null}}
	reqitems[0] = NewItem(CIPItem_ConnectionAddress, &client.OTNetworkConnectionID)

	p, err := Serialize(
		ioi.Buffer,
		CIPObject_Symbol, CIPInstance(start_instance),
	)
	if err != nil {
		return fmt.Errorf("couldn't build path. %w", err)
	}

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: client.Sequencer(),
		Service:       CIPService_GetInstanceAttributeList,
		PathLength:    byte(p.Len() / 2),
	}

	reqitems[1] = NewItem(CIPItem_ConnectedData, readmsg)
	reqitems[1].Marshal(p.Bytes())
	number_of_attr_to_receive := 3
	attr1_symbol_name := 1
	attr2_symbol_type := 2
	attr8_arraydims := 8
	//reqitems[1].Marshal([4]uint16{3, 1, 2, 8})
	reqitems[1].Marshal([4]uint16{uint16(number_of_attr_to_receive), uint16(attr1_symbol_name), uint16(attr2_symbol_type), uint16(attr8_arraydims)})
	reqitems[1].Marshal(byte(1))
	reqitems[1].Marshal(byte(0))
	reqitems[1].Marshal(uint16(1))

	hdr, data, err := client.send_recv_data(CIPCommandSendUnitData, MarshalItems(reqitems))
	if err != nil {
		return err
	}
	_ = hdr
	_ = data
	//data_hdr := ListInstanceHeader{}
	//binary.Read(data, binary.LittleEndian, &data_hdr)

	// first six bytes are zero.
	padding := make([]byte, 6)
	data.Read(padding)

	resp_items, err := ReadItems(data)
	if err != nil {
		return fmt.Errorf("couldn't parse items. %w", err)
	}

	// get ready to read tag info from item 1 data
	data2 := bytes.NewBuffer(resp_items[1].Data)
	//data2.Next(4)
	data_hdr := msgListInstanceHeader{}
	binary.Read(data2, binary.LittleEndian, &data_hdr)

	tag_hdr := new(msgtagResultDataHeader)
	tag_ftr := new(msgtagResultDataFooter)
	for data2.Len() > 0 {

		binary.Read(data2, binary.LittleEndian, tag_hdr)
		newtag_bytes := make([]byte, tag_hdr.NameLength)
		binary.Read(data2, binary.LittleEndian, &newtag_bytes)
		newtag_name := fmt.Sprintf("%s.%s", roottag, string(newtag_bytes))

		// the end of the tagname has to be aligned on a 16 bit word
		//tagname_alignment := tag_hdr.NameLength % 2
		//if tagname_alignment != 0 {
		//data2.Next(int(tagname_alignment))
		//}
		binary.Read(data2, binary.LittleEndian, tag_ftr)

		kt := KnownTag{
			Name:     newtag_name,
			Type:     tag_ftr.Type,
			Class:    CIPClass(tag_ftr.TypeInfo),
			Instance: CIPInstance(tag_hdr.InstanceID),
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
		client.KnownTags[strings.ToLower(newtag_name)] = kt

		if verbose {
			log.Printf("Tag: '%s' Instance: %d Type: %s/%d[%d,%d,%d]",
				newtag_name,
				tag_hdr.InstanceID,
				tag_ftr.Type,
				tag_ftr.TypeInfo,
				tag_ftr.Dimension1,
				tag_ftr.Dimension2,
				tag_ftr.Dimension3,
			)
		}
		start_instance = tag_hdr.InstanceID

	}
	if verbose {
		log.Printf("Status: %v", hdr.Status)
	}

	if data_hdr.Status == 6 && start_instance < 200 {
		client.ListSubTags(roottag, start_instance)
	}

	return nil
}
