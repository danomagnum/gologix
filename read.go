package gologix

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"reflect"
)

func (client *Client) Read(tag string, data any) error {
	switch data := data.(type) {
	case *bool:
		v, err := read[bool](client, tag)
		if err != nil {
			return err
		}
		*data = v
		return nil
	case *byte:
		v, err := read[byte](client, tag)
		if err != nil {
			return err
		}
		*data = v
		return nil
	case *int16:
		v, err := read[int16](client, tag)
		if err != nil {
			return err
		}
		*data = v
		return nil
	case *uint16:
		v, err := read[uint16](client, tag)
		if err != nil {
			return err
		}
		*data = v
		return nil
	case *int32:
		v, err := read[int32](client, tag)
		if err != nil {
			return err
		}
		*data = v
		return nil
	case *uint32:
		v, err := read[uint32](client, tag)
		if err != nil {
			return err
		}
		*data = v
		return nil
	case *int64:
		v, err := read[int64](client, tag)
		if err != nil {
			return err
		}
		*data = v
		return nil
	case *uint64:
		v, err := read[uint64](client, tag)
		if err != nil {
			return err
		}
		*data = v
		return nil
	case *float32:
		v, err := read[float32](client, tag)
		if err != nil {
			return err
		}
		*data = v
		return nil
	case *float64:
		v, err := read[float64](client, tag)
		if err != nil {
			return err
		}
		*data = v
		return nil
	case *string:
		v, err := read[string](client, tag)
		if err != nil {
			return err
		}
		*data = v
		return nil

	case []bool:
		elements := len(data)
		v, err := readArray[bool](client, tag, uint16(elements))
		if err != nil {
			return err
		}
		if len(v) != elements {
			return fmt.Errorf("got %d instead of %d elements", len(v), elements)
		}
		for i := range data {
			data[i] = v[i]
		}
		return nil
	case []byte:
		elements := len(data)
		v, err := readArray[byte](client, tag, uint16(elements))
		if err != nil {
			return err
		}
		if len(v) != elements {
			return fmt.Errorf("got %d instead of %d elements", len(v), elements)
		}
		for i := range data {
			data[i] = v[i]
		}
		return nil
	case []int16:
		elements := len(data)
		v, err := readArray[int16](client, tag, uint16(elements))
		if err != nil {
			return err
		}
		if len(v) != elements {
			return fmt.Errorf("got %d instead of %d elements", len(v), elements)
		}
		for i := range data {
			data[i] = v[i]
		}
		return nil
	case []uint16:
		elements := len(data)
		v, err := readArray[uint16](client, tag, uint16(elements))
		if err != nil {
			return err
		}
		if len(v) != elements {
			return fmt.Errorf("got %d instead of %d elements", len(v), elements)
		}
		for i := range data {
			data[i] = v[i]
		}
		return nil
	case []int32:
		elements := len(data)
		v, err := readArray[int32](client, tag, uint16(elements))
		if err != nil {
			return err
		}
		if len(v) != elements {
			return fmt.Errorf("got %d instead of %d elements", len(v), elements)
		}
		for i := range data {
			data[i] = v[i]
		}
		return nil
	case []uint32:
		elements := len(data)
		v, err := readArray[uint32](client, tag, uint16(elements))
		if err != nil {
			return err
		}
		if len(v) != elements {
			return fmt.Errorf("got %d instead of %d elements", len(v), elements)
		}
		for i := range data {
			data[i] = v[i]
		}
		return nil
	case []int64:
		elements := len(data)
		v, err := readArray[int64](client, tag, uint16(elements))
		if err != nil {
			return err
		}
		if len(v) != elements {
			return fmt.Errorf("got %d instead of %d elements", len(v), elements)
		}
		for i := range data {
			data[i] = v[i]
		}
		return nil
	case []uint64:
		elements := len(data)
		v, err := readArray[uint64](client, tag, uint16(elements))
		if err != nil {
			return err
		}
		if len(v) != elements {
			return fmt.Errorf("got %d instead of %d elements", len(v), elements)
		}
		for i := range data {
			data[i] = v[i]
		}
		return nil
	case []float32:
		elements := len(data)
		v, err := readArray[float32](client, tag, uint16(elements))
		if err != nil {
			return err
		}
		if len(v) != elements {
			return fmt.Errorf("got %d instead of %d elements", len(v), elements)
		}
		for i := range data {
			data[i] = v[i]
		}
		return nil
	case []float64:
		elements := len(data)
		v, err := readArray[float64](client, tag, uint16(elements))
		if err != nil {
			return err
		}
		if len(v) != elements {
			return fmt.Errorf("got %d instead of %d elements", len(v), elements)
		}
		for i := range data {
			data[i] = v[i]
		}
		return nil
	case []interface{}:
		// a pointer to a struct.
		fmt.Printf("slice of interfaces detected")
		val, err := client.read_single(tag, CIPTypeStruct, 1)
		if err != nil {
			return err
		}
		cast, ok := val.([]byte)
		if !ok {
			return errors.New("couldn't convert to byte slice")
		}
		b := bytes.NewBuffer(cast)
		err = binary.Read(b, binary.LittleEndian, data)
		if err != nil {
			return fmt.Errorf("couldn't parse str data. %w", err)
		}
	}
	// Fallback to reflect-based decoding.
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Pointer:
		// a pointer to a struct.
		fmt.Printf("struct detected")
		val, err := client.read_single(tag, CIPTypeStruct, 1)
		if err != nil {
			return err
		}
		cast, ok := val.([]byte)
		if !ok {
			return errors.New("couldn't convert to byte slice")
		}
		b := bytes.NewBuffer(cast)
		err = binary.Read(b, binary.LittleEndian, data)
		if err != nil {
			return fmt.Errorf("couldn't parse str data. %w", err)
		}

	case reflect.Slice:
		// slice of structs.
		elements := uint16(v.Len())
		fmt.Printf("slice of structs length %d detected", elements)
		val, err := client.read_single(tag, CIPTypeStruct, elements)
		if err != nil {
			return err
		}
		dat, ok := val.([]byte)
		if !ok {
			return fmt.Errorf("couldn't convert to bytes")
		}

		b := bytes.NewBuffer(dat)
		err = binary.Read(b, binary.LittleEndian, data)
		if err != nil {
			return fmt.Errorf("couldn't parse str data element %w", err)
		}
	}

	return nil
}

func (client *Client) read_single(tag string, datatype CIPType, elements uint16) (any, error) {
	ioi, err := client.newIOI(tag, datatype)

	if err != nil {
		return nil, err
	}

	reqitems := make([]CIPItem, 2)
	reqitems[0] = NewItem(CIPItem_ConnectionAddress, &client.OTNetworkConnectionID)

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: client.Sequencer(),
		Service:       CIPService_Read,
		PathLength:    byte(len(ioi.Bytes()) / 2),
	}
	// setup item
	reqitems[1] = NewItem(CIPItem_ConnectedData, readmsg)
	// add path
	reqitems[1].Marshal(ioi.Bytes())
	// add service specific data
	reqitems[1].Marshal(elements)

	//client.Send(CIPCommandSendUnitData, MarshalItems(reqitems))
	hdr, data, err := client.send_recv_data(CIPCommandSendUnitData, MarshalItems(reqitems))
	if err != nil {
		return nil, err
	}
	_ = hdr

	read_result_header := msgCIPReadResultHeader{}
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
	var hdr2 msgCIPReadResultData
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
		if datatype == CIPTypeBOOL && hdr2.Type != CIPTypeBOOL && ioi.BitAccess {
			// we have requested a bool from some other type.  Maybe a bit access?
			value := readValue(hdr2.Type, &items[1])
			return getBit(hdr2.Type, value, ioi.BitPosition)
		}
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

func readArray[T GoLogixTypes](client *Client, tag string, elements uint16) ([]T, error) {
	t := make([]T, elements)
	ct := GoVarToCIPType(t[0])
	val, err := client.read_single(tag, ct, elements)
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

func read[T GoLogixTypes](client *Client, tag string) (T, error) {
	var t T
	ct := GoVarToCIPType(t)
	val, err := client.read_single(tag, ct, 1)
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

// tag_str is a struct with each field tagged with a `gologix:"TAGNAME"` tag that specifies the tag on the client.
// The types of each field need to correspond to the correct CIP type as mapped in types.go
func (client *Client) ReadMulti(tag_str any, datatype CIPType, elements uint16) error {

	// build the tag list from the structure
	T := reflect.TypeOf(tag_str).Elem()
	vf := reflect.VisibleFields(T)
	tags := make([]string, 0)
	types := make([]CIPType, 0)
	tag_map := make(map[string]int)
	val := reflect.ValueOf(tag_str).Elem()
	for i := range vf {
		field := vf[i]
		tagpath, ok := field.Tag.Lookup("gologix")
		v := val.Field(i).Interface()
		ct := GoVarToCIPType(v)
		types = append(types, ct)
		if !ok {
			continue
		}
		tags = append(tags, tagpath)
		tag_map[tagpath] = i
	}

	// first generate IOIs for each tag
	qty := len(tags)
	iois := make([]*tagIOI, qty)
	for i, tag := range tags {
		var err error
		iois[i], err = client.newIOI(tag, datatype)
		if err != nil {
			return err
		}
	}
	// you have to change this read sequencer every time you make a new tag request.  If you don't, you
	// won't get an error but it will return the last value you requested again.
	// You don't have to keep incrementing it.  just going back and forth between 1 and 0 works OK.

	reqitems := make([]CIPItem, 2)
	reqitems[0] = NewItem(CIPItem_ConnectionAddress, &client.OTNetworkConnectionID)

	ioi_header := msgCIPConnectedMultiServiceReq{
		Sequence:     client.Sequencer(),
		Service:      CIPService_MultipleService,
		PathSize:     2,
		Path:         [4]byte{0x20, 0x02, 0x24, 0x01},
		ServiceCount: uint16(qty),
	}

	b := bytes.Buffer{}
	// we now have to build up the jump table for each IOI.
	// and pack all the IOIs together into b
	jump_table := make([]uint16, qty)
	jump_start := 2 + qty*2 // 2 bytes + 2 bytes per jump entry
	for i := 0; i < qty; i++ {
		jump_table[i] = uint16(jump_start + b.Len())
		ioi := iois[i]
		h := msgCIPMultiIOIHeader{
			Service: CIPService_Read,
			Size:    byte(len(ioi.Buffer) / 2),
		}
		f := msgCIPIOIFooter{
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

	hdr, data, err := client.send_recv_data(CIPCommandSendUnitData, MarshalItems(reqitems))
	if err != nil {
		return err
	}
	_ = hdr

	read_result_header := msgCIPReadResultHeader{}
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
	var reply_hdr msgMultiReadResultHeader
	binary.Read(&ritem, binary.LittleEndian, &reply_hdr)
	offset_table := make([]uint16, reply_hdr.Reply_Count)
	binary.Read(&ritem, binary.LittleEndian, &offset_table)
	rb := ritem.Bytes()
	result_values := make([]interface{}, reply_hdr.Reply_Count)
	for i := 0; i < int(reply_hdr.Reply_Count); i++ {
		offset := offset_table[i] + 10 // offset doesn't start at 0 in the item
		mybytes := bytes.NewBuffer(rb[offset:])
		rhdr := msgMultiReadResult{}
		binary.Read(mybytes, binary.LittleEndian, &rhdr)

		// bit 8 of the service indicates whether it is a response service
		if !rhdr.Service.IsResponse() {
			return fmt.Errorf("wasn't a response service. Got %v", rhdr.Service)
		}
		rhdr.Service = rhdr.Service.UnResponse()
		if rhdr.Status != 0 {
			return fmt.Errorf("problem reading %v. Status %v", tags[i], rhdr.Status)
		}
		if types[i] == CIPTypeBOOL && rhdr.Type != CIPTypeBOOL && iois[i].BitAccess {
			// we have requested a bool from some other type.  Maybe a bit access?
			value := readValue(rhdr.Type, &items[1])
			val, err := getBit(rhdr.Type, value, iois[i].BitPosition)
			if err != nil {
				log.Printf("problem reading value for this guy")
				continue
			}
			result_values[i] = val
		} else {
			result_values[i] = rhdr.Type.readValue(mybytes)
		}

		if verbose {
			log.Printf("Result %d @ %d. %+v. value: %v.\n", i, offset, rhdr, result_values[i])
		}
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

type msgMultiReadResultHeader struct {
	SequenceCount uint16
	Service       CIPService
	Reserved      byte
	Status        uint16
	Reply_Count   uint16
}

type msgMultiReadResult struct {
	Service   CIPService
	Reserved  byte
	Status    uint16
	Type      CIPType
	Reserved2 byte
}
