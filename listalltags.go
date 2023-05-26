package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"strings"
)

type msgtagResultDataHeader struct {
	InstanceID uint32
	NameLength uint16
}

type TagInfo struct {
	Type       CIPType
	TypeInfo   byte
	Dimension1 uint32
	Dimension2 uint32
	Dimension3 uint32
}

// see page 42 of 1756-PM020H-EN-P
func (f TagInfo) PreDefined() bool {
	val := binary.LittleEndian.Uint16([]byte{byte(f.Type), f.TypeInfo})
	var mask uint16 = 0b0000_1111_1111_1111
	var mask2 uint16 = 0x0FFF
	_ = mask2
	val2 := val & mask
	log.Print(val, mask, mask2, val2)
	//return !((val > 0x0100) && (val < 0x0EFF))
	return val2 <= 0x0100
}

// see page 42 of 1756-PM020H-EN-P
func (f TagInfo) Atomic() bool {
	val := binary.LittleEndian.Uint16([]byte{byte(f.Type), f.TypeInfo})
	return val&0b1001_0000_0000_0000 == 0
}

// see page 42 of 1756-PM020H-EN-P
func (f TagInfo) Template_ID() uint16 {
	val := binary.LittleEndian.Uint16([]byte{byte(f.Type), f.TypeInfo})
	template_mask := uint16(0b0000_0111_1111_1111)
	bit12 := uint16(1 << 12)
	bit15 := uint16(1 << 15)
	b12_set := val&bit12 != 0
	b15_set := val&bit15 != 0
	if !b15_set || b12_set {
		// not a template
		return 0
	}

	return val & template_mask

}

type msgListInstanceHeader struct {
	Service       CIPService
	Reserved      byte
	SequenceCount uint16
	Status        uint16
}

// the gist here is that we want to do a fragmented read (since there will undoubtedly be more than one packet's worth)
// of the instance attribute list of the symbol objects.
//
// see 1756-PM020H-EN-P March 2022 page 39
// also see https://forums.mrclient.com/index.php?/topic/40626-reading-and-writing-io-tags-in-plc/
func (client *Client) ListAllTags(start_instance uint32) error {

	// have to start at 1.
	if start_instance == 0 {
		start_instance = 1
	}

	reqitems := make([]cipItem, 2)
	//reqitems[0] = cipItem{Header: cipItemHeader{ID: cipItem_Null}}
	reqitems[0] = NewItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	p, err := Serialize(
		cipObject_Symbol, CIPInstance(start_instance),
	)
	if err != nil {
		return fmt.Errorf("couldn't build path. %w", err)
	}

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: uint16(sequencer()),
		Service:       cipService_GetInstanceAttributeList,
		PathLength:    byte(p.Len() / 2),
	}

	// setup item
	reqitems[1] = NewItem(cipItem_ConnectedData, readmsg)
	// add path
	reqitems[1].Serialize(p.Bytes())
	// add service specific data
	number_of_attr_to_receive := 3
	attr1_symbol_name := 1
	attr2_symbol_type := 2
	attr8_arraydims := 8
	reqitems[1].Serialize([4]uint16{uint16(number_of_attr_to_receive), uint16(attr1_symbol_name), uint16(attr2_symbol_type), uint16(attr8_arraydims)})

	// if we have an empty path, we don't include this (for micro800)
	// TODO: figure out what this is
	if client.Path.Len() != 0 {
		reqitems[1].Serialize(byte(1))
		reqitems[1].Serialize(byte(0))
		reqitems[1].Serialize(uint16(1))
	}

	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, SerializeItems(reqitems))
	if err != nil {
		return err
	}
	_ = hdr
	_ = data
	//data_hdr := ListInstanceHeader{}
	//binary.Read(data, binary.LittleEndian, &data_hdr)

	// first six bytes are zero.
	padding := make([]byte, 6)
	_, err = data.Read(padding)
	if err != nil {
		return fmt.Errorf("problem getting padding bytes. %w", err)
	}

	resp_items, err := ReadItems(data)
	if err != nil {
		return fmt.Errorf("couldn't parse items. %w", err)
	}

	// get ready to read tag info from item 1 data
	data2 := bytes.NewBuffer(resp_items[1].Data)
	//data2.Next(4)
	data_hdr := msgListInstanceHeader{}
	err = binary.Read(data2, binary.LittleEndian, &data_hdr)
	if err != nil {
		return fmt.Errorf("problem reading list instance header. %w", err)
	}

	tag_hdr := new(msgtagResultDataHeader)
	tag_ftr := new(TagInfo)
	for data2.Len() > 0 {

		err = binary.Read(data2, binary.LittleEndian, tag_hdr)
		if err != nil {
			return fmt.Errorf("problem reading tag header. %w", err)
		}
		tag_name := make([]byte, tag_hdr.NameLength)
		err = binary.Read(data2, binary.LittleEndian, &tag_name)
		if err != nil {
			return fmt.Errorf("problem reading tag name. %w", err)
		}
		tag_string := string(tag_name)

		// the end of the tagname has to be aligned on a 16 bit word
		//tagname_alignment := tag_hdr.NameLength % 2
		//if tagname_alignment != 0 {
		//data2.Next(int(tagname_alignment))
		//}
		err = binary.Read(data2, binary.LittleEndian, tag_ftr)
		if err != nil {
			return fmt.Errorf("problem reading tag footer. %w", err)
		}

		kt := KnownTag{
			Name:     string(tag_name),
			Info:     *tag_ftr,
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

		if !isValidTag(tag_string, *tag_ftr) {
			continue
		}

		if len(tag_string) > 8 {
			if strings.HasPrefix(tag_string, "Program:") {
				_, err = client.ListSubTags(tag_string, 1)
				if err != nil {
					return err
				}
				continue
			}
		}

		if tag_ftr.Template_ID() != 0 {
			log.Printf("found UDT of some sort %s", kt.Name)
		}

		if tag_ftr.Template_ID() != 0 && !tag_ftr.PreDefined() {
			if verbose {
				log.Printf("Looking up template for Tag: '%s' ", tag_string)
			}
			u, err := client.ListMembers(uint32(tag_ftr.Template_ID()))
			if err != nil {
				log.Printf("problem reading member list for %v, %+v, %+v, %v, %v", tag_string, tag_hdr, tag_ftr, tag_ftr.Template_ID(), tag_ftr.PreDefined())
				//return err
			} else {
				kt.UDT = &u
				log.Printf("Successful member read for %s", kt.Name)
			}
		}

		client.KnownTags[strings.ToLower(string(tag_name))] = kt

		if verbose {
			log.Printf("Tag: '%s' Instance: %d Type: %s/%d[%d,%d,%d].  Template %d",
				tag_name,
				tag_hdr.InstanceID,
				tag_ftr.Type,
				tag_ftr.TypeInfo,
				tag_ftr.Dimension1,
				tag_ftr.Dimension2,
				tag_ftr.Dimension3,
				tag_ftr.Template_ID(),
			)
		}
		start_instance = tag_hdr.InstanceID

	}
	if verbose {
		log.Printf("Status: %v", hdr.Status)
	}

	if data_hdr.Status == 6 && start_instance < 200 {
		err = client.ListAllTags(start_instance)
		if err != nil {
			return err
		}
	}

	return nil
}

// per 1756-PM020H-EN-P page 43 there are some conditions in which we should discard the tags
// because they aren't valid for reading/writing.
func isValidTag(tag_string string, tag_ftr TagInfo) bool {
	if !tag_ftr.Atomic() && tag_ftr.PreDefined() && verbose {
		log.Printf("Skipping Tag: '%s' Type: %s/%d[%d,%d,%d].  Template %d",
			tag_string,
			tag_ftr.Type,
			tag_ftr.TypeInfo,
			tag_ftr.Dimension1,
			tag_ftr.Dimension2,
			tag_ftr.Dimension3,
			tag_ftr.Template_ID(),
		)
	}
	if tag_string[:2] == "__" {
		if verbose {
			log.Printf("Skipping Tag: '%s' because it starts with '__'", tag_string)
		}
		return false
	}
	if strings.Contains(tag_string, ":") {
		if !strings.HasPrefix(tag_string, "Program") {
			log.Printf("Skipping Tag: '%s' because it has a : but doesn't start with 'program' ", tag_string)
			return false
		}

	}
	return true
}
