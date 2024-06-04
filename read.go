package gologix

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

// Read a single tag into data.  Data should be a pointer to the variable where the data will be deposited.
//
// If the data type is not known at read time, use the Read_single function with a CIPType_Unknown type
//
// # To efficiently read multiple tags at once, use the ReadMulti, ReadMap, or ReadList functions
//
// If the data type does not match what is returned by the controller you will get an error.
func (client *Client) Read(tag string, data any) error {
	err := client.checkConnection()
	if err != nil {
		return fmt.Errorf("could not start read: %w", err)
	}
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
		val, err := client.Read_single(tag, CIPTypeStruct, 1)
		if err != nil {
			return err
		}
		cast, ok := val.([]byte)
		if !ok {
			return errors.New("couldn't convert to byte slice")
		}
		b := bytes.NewBuffer(cast)
		//err = binary.Read(b, binary.LittleEndian, data)
		_, err = Unpack(b, data)
		if err != nil {
			return fmt.Errorf("couldn't parse str data. %w", err)
		}
	}
	// Fallback to reflect-based decoding.
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Pointer:
		// a pointer to a struct.
		val, err := client.Read_single(tag, CIPTypeStruct, 1)
		if err != nil {
			return err
		}
		cast, ok := val.([]byte)
		if !ok {
			return errors.New("couldn't convert to byte slice")
		}
		b := bytes.NewBuffer(cast)

		//err = binary.Read(b, binary.LittleEndian, data)
		_, err = Unpack(b, data)
		if err != nil {
			return fmt.Errorf("couldn't parse str data. %w", err)
		}

	case reflect.Slice:
		// slice of structs.
		elements := uint16(v.Len())
		val, err := client.Read_single(tag, CIPTypeStruct, elements)
		if err != nil {
			return err
		}
		dat, ok := val.([]byte)
		if !ok {
			return fmt.Errorf("couldn't convert to bytes")
		}

		b := bytes.NewBuffer(dat)
		//TODO: unpack here instead of just a read.
		err = binary.Read(b, binary.LittleEndian, data)
		if err != nil {
			return fmt.Errorf("couldn't parse str data element %w", err)
		}

	}

	return nil
}

// Read a single tag with the datatype given by a parameter instead of inferred from a pointer.
//
// To read data of an unknown type, use CIPTypeUnknown for the data type.
//
// The data is returned as an interface{} so you'll probably have to type assert it.
func (client *Client) Read_single(tag string, datatype CIPType, elements uint16) (any, error) {

	err := client.checkConnection()
	if err != nil {
		return nil, fmt.Errorf("could not start single read: %w", err)
	}

	ioi, err := client.newIOI(tag, datatype)

	if err != nil {
		return nil, err
	}

	reqItems := make([]CIPItem, 2)
	reqItems[0] = newItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	readMsg := msgCIPConnectedServiceReq{
		SequenceCount: uint16(sequencer()),
		Service:       CIPService_Read,
		PathLength:    byte(len(ioi.Bytes()) / 2),
	}
	// setup item
	reqItems[1] = newItem(cipItem_ConnectedData, readMsg)
	// add path
	reqItems[1].Serialize(ioi.Bytes())
	// add service specific data
	reqItems[1].Serialize(elements)

	itemData, err := serializeItems(reqItems)
	if err != nil {
		return nil, err
	}
	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, itemData)
	if err != nil {
		return nil, err
	}
	_ = hdr

	read_result_header := msgCIPResultHeader{}
	err = binary.Read(data, binary.LittleEndian, &read_result_header)
	if err != nil {
		client.Logger.Printf("Problem reading read result header. %v", err)
	}
	items, err := readItems(data)
	if err != nil {
		client.Logger.Printf("Problem reading items. %v", err)
		return 0, err
	}
	if len(items) != 2 {
		return 0, fmt.Errorf("wrong Number of Items. Expected 2 but got %v", len(items))
	}
	var hdr2 msgCIPReadResultData
	err = items[1].DeSerialize(&hdr2)
	if err != nil {
		return 0, fmt.Errorf("problem reading item 2's header. %w", err)
	}

	if hdr2.Type == CIPTypeStruct {
		if datatype == CIPTypeSTRING {
			str_hdr := cipStringHeader{}
			err = items[1].DeSerialize(&str_hdr)
			if err != nil {
				return nil, fmt.Errorf("couldn't unpack struct header. %w", err)
			}
			str := make([]byte, str_hdr.Length)
			err = items[1].DeSerialize(&str)
			if err != nil {
				return nil, fmt.Errorf("couldn't unpack struct data. %w", err)
			}
			return str, nil
		}
		str_hdr := cipStructHeader{}
		err = items[1].DeSerialize(&str_hdr)
		if err != nil {
			return nil, fmt.Errorf("couldn't unpack struct header. %w", err)
		}
		str := items[1].Data[items[1].Pos:]
		return str, nil
	}
	if elements == 1 {
		if datatype == CIPTypeBOOL && hdr2.Type != CIPTypeBOOL && ioi.BitAccess {
			// we have requested a bool from some other type.  Maybe a bit access?
			value, err := readValue(hdr2.Type, &items[1])
			if err != nil {
				return nil, fmt.Errorf("problem reading bool tag %s: %w", tag, err)
			}
			return getBit(hdr2.Type, value, ioi.BitPosition)
		}
		// not a struct so we can read the value directly
		value, err := readValue(hdr2.Type, &items[1])
		if err != nil {
			return nil, fmt.Errorf("problem reading tag %s: %w", tag, err)
		}
		return value, nil
	} else {
		value := make([]any, elements)
		for i := 0; i < int(elements); i++ {
			value[i], err = readValue(hdr2.Type, &items[1])
			if err != nil {
				return nil, fmt.Errorf("problem reading element %d of %s: %w", i, tag, err)
			}

		}
		return value, nil

	}
}

func readArray[T GoLogixTypes](client *Client, tag string, elements uint16) ([]T, error) {
	t := make([]T, elements)
	ct, _ := GoVarToCIPType(t[0])
	val, err := client.Read_single(tag, ct, elements)
	if err != nil {
		return t, err
	}

	if ct == CIPTypeStruct {
		b, ok := val.([]byte)
		if !ok {
			return t, fmt.Errorf("couldn't cast to bytes. %w", err)
		}
		return parseArrayStruct[T](b, elements)
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
	ct, _ := GoVarToCIPType(t)
	val, err := client.Read_single(tag, ct, 1)
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

type cipStringHeader struct {
	Unknown uint16
	Length  uint32
}
type cipStructHeader struct {
	Unknown uint16
}

// Tag_str is a pointer to a struct with each field tagged with a `gologix:"TAGNAME"` tag that specifies the tag on the client.
// The types of each field need to correspond to the correct CIP type as mapped in types.go
//
// To read multiple tags without creating a tagged struct, use the ReadList() or ReadMap() functions instead.
func (client *Client) ReadMulti(tag_str any) error {
	switch x := tag_str.(type) {
	case map[string]any:
		return client.ReadMap(x)
	}

	err := client.checkConnection()
	if err != nil {
		return fmt.Errorf("could not start multi read: %w", err)
	}

	// build the tag list from the structure by reflecting through the tags on the fields of the struct.
	T := reflect.TypeOf(tag_str).Elem()
	vf := reflect.VisibleFields(T)
	tags := make([]string, 0)
	types := make([]CIPType, 0)
	elements := make([]int, 0)
	tag_map := make(map[string]int)
	val := reflect.ValueOf(tag_str).Elem()
	for i := range vf {
		field := vf[i]
		tagPath, ok := field.Tag.Lookup("gologix")
		if !ok {
			continue
		}
		v := val.Field(i).Interface()
		ct, elem := GoVarToCIPType(v)
		types = append(types, ct)
		elements = append(elements, elem)
		tags = append(tags, tagPath)
		tag_map[tagPath] = i
	}

	result_values, err := client.ReadList(tags, types, elements)
	if err != nil {
		return fmt.Errorf("problem in read list: %w", err)
	}

	// now unpack the result values back into the given structure
	for i, tag := range tags {
		fieldNo := tag_map[tag]
		val := result_values[i]

		v := reflect.ValueOf(&tag_str).Elem().Elem().Elem()

		fieldVal := v.Field(fieldNo)

		if fieldVal.Type().Kind() != reflect.Slice {
			fieldVal.Set(reflect.ValueOf(val))
		} else {
			l := fieldVal.Len()
			for j := 0; j < l; j++ {
				fieldVal.Index(j).Set(reflect.ValueOf(val).Index(j).Elem())
			}
		}
	}

	return nil
}

type tagDesc struct {
	TagName  string
	TagType  CIPType
	Elements int
}

func (client *Client) readList(tags []tagDesc) ([]any, error) {

	// first generate IOIs for each tag
	qty := len(tags)
	iois := make([]*tagIOI, qty)
	for i, tag := range tags {
		var err error
		iois[i], err = client.newIOI(tag.TagName, tag.TagType)
		if err != nil {
			return nil, err
		}
	}

	reqItems := make([]CIPItem, 2)
	reqItems[0] = newItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	ioi_header := msgCIPConnectedMultiServiceReq{
		Sequence:     uint16(sequencer()),
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
			Elements: uint16(tags[i].Elements),
		}
		err := binary.Write(&b, binary.LittleEndian, h)
		if err != nil {
			return nil, fmt.Errorf("problem writing cip IO header to buffer. %w", err)
		}
		b.Write(ioi.Buffer)
		err = binary.Write(&b, binary.LittleEndian, f)
		if err != nil {
			return nil, fmt.Errorf("problem writing ioi buffer to msg buffer. %w", err)
		}
		// TODO: calculate the actual message size, not just the IOI data size.
		// TODO: We also need to calculate the response size we expect from the PLC and split
		//       into multiple messages on that also.
		if b.Len() > int(client.ConnectionSize) {
			// TODO: split this read up into multiple messages.
			return nil, fmt.Errorf("maximum read message size is %d", client.ConnectionSize)
		}
	}

	// right now I'm putting the IOI data into the cip Item, but I suspect it might actually be that the read sequencer is
	// the item's data and the service code actually starts the next portion of the message.  But the item's header length reflects
	// the total data so maybe not.
	reqItems[1] = CIPItem{Header: cipItemHeader{ID: cipItem_ConnectedData}}
	reqItems[1].Serialize(ioi_header)
	reqItems[1].Serialize(jump_table)
	reqItems[1].Serialize(&b)

	itemData, err := serializeItems(reqItems)
	if err != nil {
		return nil, err
	}
	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, itemData)
	if err != nil {
		return nil, err
	}
	_ = hdr

	if hdr.Status != 0 {
		return nil, fmt.Errorf("problem reading tags. Status %v", CIPStatus(hdr.Status))
	}

	read_result_header := msgCIPResultHeader{}
	err = binary.Read(data, binary.LittleEndian, &read_result_header)
	if err != nil {
		client.Logger.Printf("Problem reading read result header. %v", err)
	}
	items, err := readItems(data)
	if err != nil {
		return nil, fmt.Errorf("problem reading items. %w", err)
	}
	if len(items) != 2 {
		return nil, fmt.Errorf("wrong Number of Items. Expected 2 but got %v", len(items))
	}
	rItem := items[1]
	var reply_hdr msgMultiReadResultHeader
	//err = binary.Read(&rItem, binary.LittleEndian, &reply_hdr)
	reply_hdr.SequenceCount, err = rItem.Uint16()
	if err != nil {
		return nil, fmt.Errorf("problem reading reply header sequence count. %w", err)
	}
	byt, err := rItem.Byte()
	reply_hdr.Service = CIPService(byt)
	if err != nil {
		return nil, fmt.Errorf("problem reading reply header service code. %w", err)
	}
	_, err = rItem.Byte()
	if err != nil {
		return nil, fmt.Errorf("problem reading reply header padding byte. %w", err)
	}
	reply_hdr.Status, err = rItem.Uint16()
	if err != nil {
		return nil, fmt.Errorf("problem reading reply header status. %w", err)
	}
	if reply_hdr.Status != uint16(CIPStatus_OK) {
		return nil, fmt.Errorf("service returned status %v", CIPStatus(reply_hdr.Status))
	}
	reply_hdr.Reply_Count, err = rItem.Uint16()
	if err != nil {
		return nil, fmt.Errorf("problem reading reply header item count. %w", err)
	}

	offset_table := make([]uint16, reply_hdr.Reply_Count)
	err = binary.Read(&rItem, binary.LittleEndian, &offset_table)
	if err != nil {
		return nil, fmt.Errorf("problem reading offset table. %w", err)
	}
	rb, err := rItem.Bytes()
	if err != nil {
		return nil, err
	}
	result_values := make([]interface{}, reply_hdr.Reply_Count)
	for i := 0; i < int(reply_hdr.Reply_Count); i++ {
		offset := offset_table[i] + 10 // offset doesn't start at 0 in the item
		myBytes := bytes.NewBuffer(rb[offset:])
		rHdr := msgMultiReadResult{}
		err = binary.Read(myBytes, binary.LittleEndian, &rHdr)
		if err != nil {
			return nil, fmt.Errorf("problem reading multi result header. %w", err)
		}

		// bit 8 of the service indicates whether it is a response service
		if !rHdr.Service.IsResponse() {
			return nil, fmt.Errorf("wasn't a response service. Got %v", rHdr.Service)
		}
		rHdr.Service = rHdr.Service.UnResponse()
		if rHdr.Status != 0 {
			return nil, fmt.Errorf("problem reading %v. Status %v", tags[i], rHdr.Status)
		}
		if tags[i].Elements == 1 {
			if tags[i].TagType == CIPTypeBOOL && rHdr.Type != CIPTypeBOOL && iois[i].BitAccess {
				// we have requested a bool from some other type.  Maybe a bit access?
				value, err := readValue(rHdr.Type, myBytes)
				if err != nil {
					return nil, fmt.Errorf("problem reading tag %v: %w", tags[i], err)
				}
				val, err := getBit(rHdr.Type, value, iois[i].BitPosition)
				if err != nil {
					client.Logger.Printf("problem reading value for this guy")
					continue
				}
				result_values[i] = val
			} else if tags[i].TagType == CIPTypeSTRING {
				str_hdr := cipStringHeader{}
				err = binary.Read(myBytes, binary.LittleEndian, &str_hdr)
				if err != nil {
					return nil, fmt.Errorf("couldn't unpack string struct header. %w", err)
				}
				str := make([]byte, str_hdr.Length)
				err = binary.Read(myBytes, binary.LittleEndian, str)
				if err != nil {
					return nil, fmt.Errorf("couldn't unpack struct data. %w", err)
				}
				result_values[i] = string(str)
			} else {
				result_values[i], err = rHdr.Type.readValue(myBytes)
				if err != nil {
					return nil, fmt.Errorf("problem reading tag %v: %w", tags[i], err)
				}
			}

			if verbose {
				client.Logger.Printf("Result %d @ %d. %+v. value: %v.\n", i, offset, rHdr, result_values[i])
			}
		} else {
			// multi-element type.
			val := make([]any, tags[i].Elements)
			for respIndex := 0; respIndex < tags[i].Elements; respIndex++ {
				value, err := readValue(rHdr.Type, myBytes)
				if err != nil {
					return nil, fmt.Errorf("problem reading tag %v: %w", tags[i], err)
				}
				val[respIndex] = value
			}
			result_values[i] = val
		}
	}

	return result_values, nil

}

func parseArrayStruct[T GoLogixTypes](dat []byte, elements uint16) ([]T, error) {
	t := make([]T, elements)
	// val should be a byte slice
	b := bytes.NewBuffer(dat)
	for i := 0; i < int(elements); i++ {
		//err := binary.Read(b, binary.LittleEndian, &t[i])
		_, err := Unpack(b, &t[i])
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

// Function for reading multiple tags at once where the tags are in a go map.
// the keys in the map are the tag names, and the values need to be the correct type
// for the tag.  The ReadMap function will update the values in the map to the current values
// in the controller.
//
// Example:
//
//		m := make(map[string]any) // define the map
//		m["TestInt"] = int16(0) // the controller has a tag "TestInt" that is an INT
//		m["TestDint"] = int32(0) // the controller has a tag "TestDint" that is a DINT
//	    err = client.ReadMulti(&mr) // do the read.
func (client *Client) ReadMap(m map[string]any) error {

	err := client.checkConnection()
	if err != nil {
		return fmt.Errorf("could not start multi read: %w", err)
	}

	total := len(m)
	tags := make([]tagDesc, total)
	indexes := make([]string, total)
	i := 0
	for k := range m {
		v := m[k]
		ct, elem := GoVarToCIPType(v)
		tags[i] = tagDesc{TagName: k, TagType: ct, Elements: elem}
		indexes[i] = k
		i++
	}
	result_values := make([]any, 0, total)

	n := 0
	msgs := 0
	n_new := 0
	for n < total {
		msgs += 1
		n_new, err = client.countIOIsThatFit(tags[n:])
		if err != nil {
			return err
		}
		subResults, err := client.readList(tags[n : n+n_new])
		n += n_new
		if err != nil {
			return err
		}
		result_values = append(result_values, subResults...)

	}

	// now unpack the result values back into the given structure
	for i := range result_values {
		m[indexes[i]] = result_values[i]
	}

	return nil
}
