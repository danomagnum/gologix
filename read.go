package gologix

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"reflect"
)

func (client *Client) read_single(tag string, datatype CIPType, elements uint16) (any, error) {
	ioi := NewIOI(tag, datatype)
	// you have to change this read sequencer every time you make a new tag request.  If you don't, you
	// won't get an error but it will return the last value you requested again.
	// You don't even have to keep incrementing it.  just going back and forth between 1 and 0 works OK.
	plc.readSequencer += 1

	reqitems := make([]CIPItem, 2)
	reqitems[0] = NewItem(CIPItem_ConnectionAddress, &plc.OTNetworkConnectionID)

	// right now I'm putting the IOI data into the cip Item, but I suspect it might actually be that the readsequencer is
	// the item's data and the service code actually starts the next portion of the message.  But the item's header length reflects
	// the total data so maybe not.
	reqitems[1] = CIPItem{Header: CIPItemHeader{ID: CIPItem_ConnectedData}}
	//reqitems[1].Marshal(ioi_header)
	//reqitems[1].Marshal(ioi.Buffer)
	//reqitems[1].Marshal(ioi_footer)
	reqitems[1].Marshal(plc.readSequencer)
	reqitems[1].Marshal(ioi.Service(CIPService_Read).Bytes())
	reqitems[1].Marshal(elements)

	plc.Send(CIPCommandSendUnitData, MarshalItems(reqitems))
	hdr, data, err := plc.recv_data()
	if err != nil {
		return nil, err
	}
	_ = hdr

	read_result_header := CIPReadResultHeader{}
	err = binary.Read(data, binary.LittleEndian, &read_result_header)
	if err != nil {
		log.Printf("Problem reading read result header. %v", err)
	}
	items, err := ReadItems(data)
	if err != nil {
		log.Printf("Problem reading items. %v", err)
		return 0, nil
	}
	if len(items) != 2 {
		return 0, fmt.Errorf("wrong Number of Items. Expected 2 but got %v", len(items))
	}
	var hdr2 CIPReadResultData
	err = items[1].Unmarshal(&hdr2)
	if err != nil {
		return 0, fmt.Errorf("problem reading item 2's header. %w", err)
	}

	if hdr2.Type == CIPTypeStruct {
		if datatype == CIPTypeSTRING {
			str_hdr := CIPStringHeader{}
			err = items[1].Unmarshal(&str_hdr)
			if err != nil {
				return nil, fmt.Errorf("couldn't unpack struct header. %w", err)
			}
			str := make([]byte, str_hdr.Length)
			err = items[1].Unmarshal(&str)
			if err != nil {
				return nil, fmt.Errorf("couldn't unpack struct data. %w", err)
			}
			return str, nil
		}
		str_hdr := CIPStructHeader{}
		err = items[1].Unmarshal(&str_hdr)
		if err != nil {
			return nil, fmt.Errorf("couldn't unpack struct header. %w", err)
		}
		str := items[1].Data[items[1].Pos:]
		return str, nil
	}
	if elements == 1 {
		// not a struct so we can read the value directly
		value := readValue(hdr2.Type, &items[1])
		return value, nil
	} else {
		value := make([]any, elements)
		for i := 0; i < int(elements); i++ {
			value[i] = readValue(hdr2.Type, &items[1])

		}
		return value, nil

	}
}

func ReadArray[T GoLogixTypes](client *Client, tag string, elements uint16) ([]T, error) {
	t := make([]T, elements)
	ct := GoVarToCIPType(t[0])
	val, err := plc.read_single(tag, ct, elements)
	if err != nil {
		return t, err
	}

	if ct == CIPTypeStruct {
		b, ok := val.([]byte)
		if !ok {
			return t, fmt.Errorf("couldn't cast to bytes. %w", err)
		}
		return parseArrayStuct[T](b, elements)
	}

	//cast, ok := val.([]T)
	cast, ok := val.([]any)
	if !ok {
		return t, fmt.Errorf("couldn't cast array. %w", err)
	}
	for i, v := range cast {
		if ct == CIPTypeSTRING {
			// v should be a byte slice
			cast2, ok := any(v).([]byte)
			if !ok {
				return t, errors.New("couldn't convert to byte slice")
			}
			s := string(cast2)
			t[i], ok = any(s).(T)
			if !ok {
				return t, errors.New("couldn't convert to string")
			}
		}
		cast2, ok := v.(T)
		if !ok {
			return t, errors.New("couldn't convert to correct type")
		}
		t[i] = cast2
	}
	return t, nil

}

func Read[T GoLogixTypes](client *Client, tag string) (T, error) {
	var t T
	//fmt.Printf("reading type %T", t)
	ct := GoVarToCIPType(t)
	val, err := plc.read_single(tag, ct, 1)
	if err != nil {
		return t, err
	}
	if ct == CIPTypeStruct {
		// val should be a byte slice
		cast, ok := val.([]byte)
		if !ok {
			return t, errors.New("couldn't convert to byte slice")
		}
		b := bytes.NewBuffer(cast)
		err := binary.Read(b, binary.LittleEndian, &t)
		if err != nil {
			return t, fmt.Errorf("couldn't parse str data. %w", err)
		}
		return t, nil
	}
	if ct == CIPTypeSTRING {
		// val should be a byte slice
		cast, ok := val.([]byte)
		if !ok {
			return t, errors.New("couldn't convert to byte slice")
		}
		s := string(cast)
		t, ok = any(s).(T)
		if !ok {
			return t, errors.New("couldn't convert to string")
		}
		return t, nil
	}
	cast, ok := val.(T)
	if !ok {
		return t, errors.New("couldn't convert to correct type")
	}
	return cast, nil

}

type CIPStringHeader struct {
	Unknown uint16
	Length  uint32
}
type CIPStructHeader struct {
	Unknown uint16
}

// tag_str is a struct with each field tagged with a `gologix:"TAGNAME"` tag that specifies the tag on the PLC.
// The types of each field need to correspond to the correct CIP type as mapped in types.go
func (client *Client) read_multi(tag_str any, datatype CIPType, elements uint16) error {

	// build the tag list from the structure
	T := reflect.TypeOf(tag_str).Elem()
	vf := reflect.VisibleFields(T)
	tags := make([]string, 0)
	tag_map := make(map[string]int)
	for i := range vf {
		field := vf[i]
		tagpath, ok := field.Tag.Lookup("gologix")
		if !ok {
			continue
		}
		tags = append(tags, tagpath)
		tag_map[tagpath] = i
	}

	// first generate IOIs for each tag
	qty := len(tags)
	iois := make([]*IOI, qty)
	for i, tag := range tags {
		iois[i] = NewIOI(tag, datatype)
	}
	// you have to change this read sequencer every time you make a new tag request.  If you don't, you
	// won't get an error but it will return the last value you requested again.
	// You don't have to keep incrementing it.  just going back and forth between 1 and 0 works OK.
	plc.readSequencer += 1

	reqitems := make([]CIPItem, 2)
	reqitems[0] = NewItem(CIPItem_ConnectionAddress, &plc.OTNetworkConnectionID)

	ioi_header := CIPMultiServiceHeader{
		Sequence:     plc.readSequencer,
		Service:      CIPService_MultipleService,
		PathSize:     2,
		Path:         [4]byte{0x20, 0x02, 0x24, 0x01},
		ServiceCount: uint16(qty),
	}

	b := bytes.Buffer{}
	// we now have to build up the jump table for each IOI.
	jump_table := make([]uint16, qty)
	jump_start := 2 + qty*2 // 2 bytes + 2 bytes per jump entry
	for i := 0; i < qty; i++ {
		jump_table[i] = uint16(jump_start + b.Len())
		ioi := iois[i]
		h := CIPMultiIOIHeader{
			Service: CIPService_Read,
			Size:    byte(len(ioi.Buffer) / 2),
		}
		f := CIPIOIFooter{
			Elements: elements,
		}
		binary.Write(&b, binary.LittleEndian, h)
		b.Write(ioi.Buffer)
		binary.Write(&b, binary.LittleEndian, f)
	}

	// right now I'm putting the IOI data into the cip Item, but I suspect it might actually be that the readsequencer is
	// the item's data and the service code actually starts the next portion of the message.  But the item's header length reflects
	// the total data so maybe not.
	reqitems[1] = CIPItem{Header: CIPItemHeader{ID: CIPItem_ConnectedData}}
	reqitems[1].Marshal(ioi_header)
	reqitems[1].Marshal(jump_table)
	reqitems[1].Marshal(b.Bytes())

	plc.Send(CIPCommandSendUnitData, MarshalItems(reqitems))
	hdr, data, err := plc.recv_data()
	if err != nil {
		return err
	}
	_ = hdr

	read_result_header := CIPReadResultHeader{}
	err = binary.Read(data, binary.LittleEndian, &read_result_header)
	if err != nil {
		log.Printf("Problem reading read result header. %v", err)
	}
	items, err := ReadItems(data)
	if err != nil {
		return fmt.Errorf("problem reading items. %w", err)
	}
	if len(items) != 2 {
		return fmt.Errorf("wrong Number of Items. Expected 2 but got %v", len(items))
	}
	ritem := items[1]
	var reply_hdr MultiReadResultHeader
	binary.Read(&ritem, binary.LittleEndian, &reply_hdr)
	offset_table := make([]uint16, reply_hdr.Reply_Count)
	binary.Read(&ritem, binary.LittleEndian, &offset_table)
	rb := ritem.Bytes()
	result_values := make([]interface{}, reply_hdr.Reply_Count)
	for i := 0; i < int(reply_hdr.Reply_Count); i++ {
		offset := offset_table[i] + 10 // offset doesn't start at 0 in the item
		mybytes := bytes.NewBuffer(rb[offset:])
		rhdr := MultiReadResult{}
		binary.Read(mybytes, binary.LittleEndian, &rhdr)

		// bit 8 of the service indicates whether it is a response service
		if !rhdr.Service.IsResponse() {
			return fmt.Errorf("wasn't a response service. Got %v", rhdr.Service)
		}
		rhdr.Service = rhdr.Service.UnResponse()
		if rhdr.Status != 0 {
			return fmt.Errorf("problem reading %v. Status %v", tags[i], rhdr.Status)
		}

		result_values[i] = rhdr.Type.readValue(mybytes)

		fmt.Printf("Result %d @ %d. %+v. value: %v.\n", i, offset, rhdr, result_values[i])
	}

	// now unpack the result values back into the given structure
	for i, tag := range tags {
		fieldno := tag_map[tag]
		val := result_values[i]

		v := reflect.ValueOf(&tag_str).Elem().Elem().Elem()

		fieldVal := v.Field(fieldno)
		fieldVal.Set(reflect.ValueOf(val))

		if err != nil {
			return fmt.Errorf("problem populating field %v with tag %v of value %v", fieldno, tag, val)
		}

	}

	return nil
}

func parseArrayStuct[T GoLogixTypes](dat []byte, elements uint16) ([]T, error) {
	t := make([]T, elements)
	// val should be a byte slice
	b := bytes.NewBuffer(dat)
	for i := 0; i < int(elements); i++ {
		err := binary.Read(b, binary.LittleEndian, &t[i])
		if err != nil {
			return t, fmt.Errorf("couldn't parse str data. %w", err)
		}
	}
	return t, nil
}

type MultiReadResultHeader struct {
	SequenceCount uint16
	Service       CIPService
	Reserved      byte
	Status        uint16
	Reply_Count   uint16
}

type MultiReadResult struct {
	Service   CIPService
	Reserved  byte
	Status    uint16
	Type      CIPType
	Reserved2 byte
}
