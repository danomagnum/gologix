package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

// returns whether the type of the tag is a pre-defined type like a DINT, SINT, INT, REAL, etc...
// see page 42 of 1756-PM020H-EN-P
func (f TagInfo) PreDefined() bool {
	val := binary.LittleEndian.Uint16([]byte{byte(f.Type), f.TypeInfo})
	var mask uint16 = 0b0000_1111_1111_1111
	var mask2 uint16 = 0x0FFF
	_ = mask2
	val2 := val & mask
	return val2 <= 0x0100
}

// returns whether the type of the tag is an atomic type like a DINT, SINT, INT, REAL, etc...
// see page 42 of 1756-PM020H-EN-P
func (f TagInfo) Atomic() bool {
	val := binary.LittleEndian.Uint16([]byte{byte(f.Type), f.TypeInfo})
	return val&0b1001_0000_0000_0000 == 0
}

// The template ID is basically the type of the tag.  Probably a udt.
// see page 42 of 1756-PM020H-EN-P
func (f TagInfo) Template_ID() uint16 {
	val := binary.LittleEndian.Uint16([]byte{byte(f.Type), f.TypeInfo})
	template_mask := uint16(0b0000_1111_1111_1111) // spec says first 11 bits, but built-in types use 12th.
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

// ListAllTags discovers and catalogs all tags available in the PLC, starting from the specified instance ID.
//
// This function performs a comprehensive scan of the PLC's symbol table to build a complete inventory
// of available tags. The discovered tags are stored in the client's KnownTags map for later reference.
//
// The function automatically handles:
//   - Controller-scoped tags (global tags)
//   - Program-scoped tags (tags within program instances)
//   - Array tags with dimension information
//   - UDT (User Defined Type) tags with structure definitions
//   - Fragmented responses for large tag lists
//
// Parameters:
//   - start_instance: The symbol instance ID to start scanning from (use 0 to scan all tags)
//
// After successful completion, access the discovered tags via:
//   - client.KnownTags: map[string]KnownTag containing all discovered tags (keyed by lowercase name)
//   - client.KnownPrograms: slice of discovered program instances
//   - client.KnownTypes: map of UDT definitions
//
// Examples:
//
//	// Discover all tags in the PLC
//	err := client.ListAllTags(0)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Access discovered tags
//	for tagName, tag := range client.KnownTags {
//	    fmt.Printf("Tag: %s, Type: %s, Array: %v\n",
//	        tag.Name, tag.Info.Type, tag.Array_Order)
//	}
//
//	// Find a specific tag
//	if tag, exists := client.KnownTags["testint"]; exists {
//	    fmt.Printf("Found tag: %s with type %s\n", tag.Name, tag.Info.Type)
//	}
//
//	// Check if it's an array
//	if len(tag.Array_Order) > 0 {
//	    fmt.Printf("Array dimensions: %v\n", tag.Array_Order)
//	}
//
// Note: This operation can take significant time on PLCs with large numbers of tags.
//
// The function automatically calls ListAllPrograms() first when start_instance is 0
// to ensure program-scoped tags are properly discovered and categorized.
func (client *Client) ListAllTags(start_instance uint32) error {
	const minimumTagValue = 1
	if start_instance < minimumTagValue {
		start_instance = minimumTagValue
	}

	// if we are starting from scratch, we should list all the programs first so we have
	// their instance IDs when we come across program scoped tags.
	if start_instance == 1 {
		err := client.ListAllPrograms()
		if err != nil {
			return fmt.Errorf("problem listing programs: %w", err)
		}
		for _, p := range client.KnownPrograms {
			_, err := client.ListSubTags(p, 1)
			if err != nil {
				return fmt.Errorf("problem listing sub tags for %s: %w", p.Name, err)
			}

		}
	}

	reqItems := make([]CIPItem, 2)
	reqItems[0] = newItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	path, err := Serialize(
		CipObject_Symbol, CIPInstance(start_instance),
	)
	if err != nil {
		return fmt.Errorf("couldn't build path. %w", err)
	}

	readMsg := msgCIPConnectedServiceReq{
		SequenceCount: uint16(sequencer()),
		Service:       CIPService_GetInstanceAttributeList,
		PathLength:    byte(path.Len() / 2),
	}

	// setup item
	reqItems[1] = newItem(cipItem_ConnectedData, readMsg)
	// add path
	err = reqItems[1].Serialize(path.Bytes())
	if err != nil {
		return fmt.Errorf("problem serializing path: %w", err)
	}
	// add service specific data
	number_of_attr_to_receive := 3
	attr1_symbol_name := 1
	attr2_symbol_type := 2
	attr8_arrayDims := 8
	err = reqItems[1].Serialize([4]uint16{
		uint16(number_of_attr_to_receive),
		uint16(attr1_symbol_name),
		uint16(attr2_symbol_type),
		uint16(attr8_arrayDims),
	})
	if err != nil {
		return fmt.Errorf("problem serializing item data: %w", err)
	}

	itemData, err := serializeItems(reqItems)
	if err != nil {
		return fmt.Errorf("problem serializing items: %w", err)
	}
	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, itemData)
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

	resp_items, err := readItems(data)
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
		start_instance = tag_hdr.InstanceID

		tag_name := make([]byte, tag_hdr.NameLength)
		err = binary.Read(data2, binary.LittleEndian, &tag_name)
		if err != nil {
			return fmt.Errorf("problem reading tag name. %w", err)
		}
		tag_string := string(tag_name)

		// the end of the tagName has to be aligned on a 16 bit word
		//tagName_alignment := tag_hdr.NameLength % 2
		//if tagName_alignment != 0 {
		//data2.Next(int(tagName_alignment))
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
				// we have a program scoped tag.  These are read separately.
				continue
			}
		}

		if tag_ftr.Template_ID() != 0 {
			client.Logger.Debug("found UDT of some sort", "name", kt.Name)
		}

		if tag_ftr.Template_ID() != 0 { //&& !tag_ftr.PreDefined() {
			client.Logger.Debug("Looking up template", "tag name", tag_string)
			u, err := client.ListMembers(uint32(tag_ftr.Template_ID()))
			if err != nil {
				client.Logger.Error("problem reading member list",
					"string", tag_string,
					"header", tag_hdr,
					"footer", tag_ftr,
					"template ID", tag_ftr.Template_ID(),
					"predefined", tag_ftr.PreDefined())
				//return err
			} else {
				kt.UDT = &u
				client.Logger.Debug("Successful member read for %s", "name", kt.Name)
				client.KnownTypes[u.Name] = u
			}
		}

		client.KnownTags[strings.ToLower(string(tag_name))] = kt

	}

	if data_hdr.Status == uint16(CIPStatus_PartialTransfer) { //} && start_instance < 200 {
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
	// Check for empty string
    if tag_string == "" {
        return false
    }
    
    // Check for system tags (start with __)
    if strings.HasPrefix(tag_string, "__") {
        return false
    }
	if strings.Contains(tag_string, ":") {
		if !strings.HasPrefix(tag_string, "Program") {
			return false
		}

	}
	_ = tag_ftr
	return true
}

