package gologix

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

// Read reads a single tag value from the PLC into the provided data variable.
//
// The data parameter must be a pointer to a variable with a type that matches the PLC tag's data type.
// Supported Go types are mapped to CIP types as documented in types.go.
//
// For arrays and slices, data should be a pre-allocated slice of the correct length:
//   - Use a pointer (&variable) for scalar values and structs
//   - Use a slice directly (no pointer) for arrays
//
// Examples:
//
//	var intTag int16
//	err := client.Read("TestInt", &intTag)  // Read INT tag
//
//	var realTag float32
//	err := client.Read("TestReal", &realTag)  // Read REAL tag
//
//	intArray := make([]int32, 5)
//	err := client.Read("TestDintArr[2]", intArray)  // Read 5 elements starting at index 2
//
//	var udtTag MyStruct
//	err := client.Read("MyUDTTag", &udtTag)  // Read UDT into struct
//
// For reading multiple tags efficiently, use ReadMulti, ReadMap, or ReadList instead.
// If the tag data type is unknown at compile time, use Read_single with CIPType_Unknown.
//
// Returns an error if the connection fails, the tag doesn't exist, or there's a type mismatch.
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
	case *int8:
		v, err := read[int8](client, tag)
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
		count := elements / 32
		if count*32 != elements {
			return fmt.Errorf("slice length must be a multiple of 32 for []bool, got %d", elements)
		}
		if count == 1 { // special case for 1 element slice - have to read it as an atomic uint32
			v, err := read[uint32](client, tag)
			if err != nil {
				return err
			}
			for i := 0; i < 32; i++ {
				data[i] = (v & (1 << i)) != 0
			}
			return nil
		}
		v, err := readArray[uint32](client, tag, uint16(count))
		if err != nil {
			return err
		}
		if len(v) != count {
			return fmt.Errorf("got %d instead of %d elements", len(v), elements)
		}
		for word := range v {
			for i := 0; i < 32; i++ {
				bit := word*32 + i
				data[bit] = (v[word] & (1 << i)) != 0
			}
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
	case []int8:
		elements := len(data)
		v, err := readArray[int8](client, tag, uint16(elements))
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
	case []string:
		elements := len(data)
		v, err := readArray[string](client, tag, uint16(elements))
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
	case *any:
		// could be anything?
		val, err := client.Read_single(tag, CIPTypeStruct, 1)
		if err != nil {
			return err
		}
		reflect.ValueOf(data).Elem().Set(reflect.ValueOf(val))

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
		//err = binary.Read(b, binary.LittleEndian, data)
		_, err = Unpack(b, data)
		if err != nil {
			return fmt.Errorf("couldn't parse str data element %w", err)
		}

	}

	return nil
}

// Read_single reads a single tag with an explicitly specified data type instead of inferring the type from a pointer.
//
// This function allows you to read tags when the data type is not known at compile time. Use CIPType_Unknown
// to read a tag of unknown type - the PLC will return the actual data type and value.
//
// Parameters:
//   - tag: The name of the tag to read (case insensitive)
//   - datatype: The expected CIP data type (use CIPType_Unknown for auto-detection)
//   - elements: Number of elements to read (1 for scalar values)
//
// Returns the data as interface{}, which you'll need to type assert to the appropriate Go type.
//
// Examples:
//
//	// Read a tag of known type
//	value, err := client.Read_single("TestInt", CIPType_INT, 1)
//	intValue := value.(int16)
//
//	// Read a tag of unknown type
//	value, err := client.Read_single("UnknownTag", CIPType_Unknown, 1)
//
//	// Read multiple elements
//	values, err := client.Read_single("IntArray", CIPType_INT, 5)
//	intSlice := values.([]interface{})
//
// For strongly-typed reading, use the Read function instead. For multiple tags, use ReadMulti or ReadMap.
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
	err = reqItems[1].Serialize(ioi, elements)
	if err != nil {
		return nil, fmt.Errorf("problem serializing ioi: %w", err)
	}

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
		client.Logger.Warn("Problem reading read result header", "error", err)
	}
	items, err := readItems(data)
	if err != nil {
		client.Logger.Warn("Problem reading items", "error", err)
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
			if elements == 1 {
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
			response := make([]any, elements)

			for i := 0; i < int(elements); i++ {
				str_hdr := cipStringHeader{}
				err = items[1].DeSerialize(&str_hdr)
				if err != nil {
					return nil, fmt.Errorf("couldn't unpack struct header. %w", err)
				}
				str := make([]byte, 82)
				err = items[1].DeSerialize(&str)
				if err != nil {
					return nil, fmt.Errorf("couldn't unpack struct data. %w", err)
				}
				response[i] = str[:str_hdr.Length]
			}
			return response, nil
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
			continue
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
	StructTypeCRC uint16
}

// ReadMulti reads multiple tags efficiently in a single request using struct field tags or a map.
//
// This function supports two input types:
//
//  1. Struct with gologix field tags:
//     Each field should have a `gologix:"tagname"` tag specifying the PLC tag to read.
//     Field types must match the corresponding CIP types as documented in types.go.
//
//     Example:
//     type MyTags struct {
//     IntTag    int16     `gologix:"TestInt"`
//     RealTag   float32   `gologix:"TestReal"`
//     ArrayTag  []int32   `gologix:"TestDintArr[2]"`  // Read 5 elements starting at index 2
//     }
//     var tags MyTags
//     tags.ArrayTag = make([]int32, 5)  // Pre-allocate slice
//     err := client.ReadMulti(&tags)
//
//  2. Map[string]any:
//     Keys are tag names, values are variables with correct types.
//     The function updates the map values with data from the PLC.
//
//     Example:
//     m := map[string]any{
//     "TestInt":         int16(0),
//     "TestReal":        float32(0),
//     "TestDintArr[2]":  make([]int32, 5),
//     }
//     err := client.ReadMulti(m)
//
// For struct-based reading without field tags, or when working with maps exclusively,
// use ReadMap instead. For reading tags with different types, use ReadList.
//
// ReadMulti automatically splits large requests across multiple messages if needed
// to stay within connection size limits.
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
	taglist := make([]tagDesc, 0, len(vf))
	tags := make([]string, 0)
	tag_map := make(map[string]int)
	val := reflect.ValueOf(tag_str).Elem()
	for i := range vf {
		field := vf[i]
		tagPath, ok := field.Tag.Lookup("gologix")
		if !ok || tagPath == "" {
			continue
		}
		v := val.Field(i).Interface()
		t, elem := GoVarToCIPType(v)
		tags = append(tags, tagPath)
		tag_map[tagPath] = i
		taglist = append(taglist, tagDesc{
			TagName:  tagPath,
			TagType:  t,
			Elements: elem,
			Struct:   v,
		})
	}

	result_values, err := client.readList(taglist)
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
	Struct   any
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

		totalMsgSize := b.Len() + SizeOf(h) + SizeOf(f) + 2 // 2 for the jump table entry
		// TODO: calculate the actual message size, not just the IOI data size.
		// TODO: We also need to calculate the response size we expect from the PLC and split
		//       into multiple messages on that also.
		if totalMsgSize > int(client.ConnectionSize) {
			first_part := tags[:i]
			rest := tags[i:]

			results0, err := client.readList(first_part)
			if err != nil {
				return nil, fmt.Errorf("problem reading first part of tags: %w", err)
			}
			results1, err := client.readList(rest)
			if err != nil {
				return nil, fmt.Errorf("problem reading second part of tags: %w", err)
			}
			return append(results0, results1...), nil

		}
	}

	// right now I'm putting the IOI data into the cip Item, but I suspect it might actually be that the read sequencer is
	// the item's data and the service code actually starts the next portion of the message.  But the item's header length reflects
	// the total data so maybe not.
	reqItems[1] = CIPItem{Header: cipItemHeader{ID: cipItem_ConnectedData}}
	err := reqItems[1].Serialize(ioi_header, jump_table, &b)
	if err != nil {
		return nil, fmt.Errorf("problem serializing item header: %w", err)
	}

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
		client.Logger.Warn("Problem reading read result header", "error", err)
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
	_ = reply_hdr.SequenceCount
	byt, err := rItem.Byte()
	reply_hdr.Service = CIPService(byt)
	if err != nil {
		return nil, fmt.Errorf("problem reading reply header service code. %w", err)
	}
	_ = reply_hdr.Service
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
		if rHdr.Status != uint16(CIPStatus_OK) {
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
					client.Logger.Warn("problem reading value for this guy")
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
			} else if rHdr.Type == CIPTypeStruct {
				typehash := cipStructHeader{}
				err := binary.Read(myBytes, binary.LittleEndian, &typehash)
				if err != nil {
					return nil, fmt.Errorf("couldn't unpack struct header. %w", err)
				}
				//dat := make([]byte, binary.Size(tags[i].Struct))
				if tags[i].Struct == nil {
					// just return the byte array.
					result_values[i] = myBytes.Bytes()
					continue
				}
				x := reflect.New(reflect.TypeOf(tags[i].Struct)).Interface()
				err = binary.Read(myBytes, binary.LittleEndian, x)
				if err != nil {
					return nil, fmt.Errorf("couldn't unpack struct data. %w", err)
				}
				// depointer the result to the correct type.
				y := reflect.ValueOf(x).Elem().Interface()
				result_values[i] = y
			} else {
				result_values[i], err = rHdr.Type.readValue(myBytes)
				if err != nil {
					return nil, fmt.Errorf("problem reading tag %v: %w", tags[i], err)
				}
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

// ReadMap reads multiple tags efficiently using a map where keys are tag names and values define the expected types.
//
// The map keys specify the tag names to read from the PLC (case insensitive).
// The map values must be initialized with the correct Go types that correspond to the PLC tag types.
// After the function completes successfully, the map values are updated with the current PLC values.
//
// Supported value types must match CIP types as documented in types.go:
//   - Scalar types: int16 (INT), int32 (DINT), float32 (REAL), bool (BOOL), etc.
//   - Arrays: pre-allocated slices like []int32, []float32, etc.
//   - Strings: string type for STRING tags
//   - Structs: user-defined types for UDT tags
//
// Examples:
//
//	// Basic scalar tags
//	m := map[string]any{
//	    "TestInt":    int16(0),      // Read INT tag
//	    "TestDint":   int32(0),      // Read DINT tag
//	    "TestReal":   float32(0),    // Read REAL tag
//	    "TestBool":   false,         // Read BOOL tag
//	    "TestString": "",            // Read STRING tag
//	}
//
//	// Array tags (specify starting index and pre-allocate slice)
//	m["TestDintArr[2]"] = make([]int32, 5)  // Read 5 DINTs starting at index 2
//	m["TestRealArr[0]"] = make([]float32, 10) // Read 10 REALs starting at index 0
//
//	// UDT tags
//	type MyUDT struct {
//	    Field1 int32
//	    Field2 float32
//	}
//	m["MyUDTTag"] = MyUDT{}
//
//	err := client.ReadMap(m)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Access the read values
//	intValue := m["TestInt"].(int16)
//	arrayValue := m["TestDintArr[2]"].([]int32)
//
// ReadMap automatically handles message splitting for large requests and optimizes
// network usage by grouping tags into the minimum number of requests.
//
// For struct-based reading with field tags, use ReadMulti instead.
// For reading tags with unknown types, set map values to nil.
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
		tags[i] = tagDesc{
			TagName:  k,
			TagType:  ct,
			Elements: elem,
			Struct:   v,
		}
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
