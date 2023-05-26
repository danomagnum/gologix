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
		_, err = Unpack(b, CIPPack{}, data)
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
		_, err = Unpack(b, CIPPack{}, data)
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

func (client *Client) Read_single(tag string, datatype CIPType, elements uint16) (any, error) {

	err := client.checkConnection()
	if err != nil {
		return nil, fmt.Errorf("could not start single read: %w", err)
	}

	ioi, err := client.NewIOI(tag, datatype)

	if err != nil {
		return nil, err
	}

	reqitems := make([]cipItem, 2)
	reqitems[0] = NewItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: uint16(sequencer()),
		Service:       cipService_Read,
		PathLength:    byte(len(ioi.Bytes()) / 2),
	}
	// setup item
	reqitems[1] = NewItem(cipItem_ConnectedData, readmsg)
	// add path
	reqitems[1].Serialize(ioi.Bytes())
	// add service specific data
	reqitems[1].Serialize(elements)

	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, SerializeItems(reqitems))
	if err != nil {
		return nil, err
	}
	_ = hdr

	read_result_header := msgCIPResultHeader{}
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
	err = items[1].DeSerialize(&hdr2)
	if err != nil {
		return 0, fmt.Errorf("problem reading item 2's header. %w", err)
	}

	if hdr2.Type == CIPTypeStruct {
		if datatype == CIPTypeSTRING {
			str_hdr := CIPStringHeader{}
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
		str_hdr := CIPStructHeader{}
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
	ct := GoVarToCIPType(t[0])
	val, err := client.Read_single(tag, ct, elements)
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

type CIPStringHeader struct {
	Unknown uint16
	Length  uint32
}
type CIPStructHeader struct {
	Unknown uint16
}

// tag_str is a pointer to a struct with each field tagged with a `gologix:"TAGNAME"` tag that specifies the tag on the client.
// The types of each field need to correspond to the correct CIP type as mapped in types.go
func (client *Client) ReadMulti(tag_str any) error {

	err := client.checkConnection()
	if err != nil {
		return fmt.Errorf("could not start multi read: %w", err)
	}

	// build the tag list from the structure by reflecing through the tags on the fields of the struct.
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
	result_values, err := client.ReadList(tags, types)
	if err != nil {
		return fmt.Errorf("problem in read list: %w", err)
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

func (client *Client) ReadList(tags []string, types []CIPType) ([]any, error) {

	err := client.checkConnection()
	if err != nil {
		return nil, fmt.Errorf("could not start list read: %w", err)
	}

	// first generate IOIs for each tag
	qty := len(tags)
	iois := make([]*tagIOI, qty)
	for i, tag := range tags {
		var err error
		iois[i], err = client.NewIOI(tag, types[i])
		if err != nil {
			return nil, err
		}
	}

	reqitems := make([]cipItem, 2)
	reqitems[0] = NewItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	ioi_header := msgCIPConnectedMultiServiceReq{
		Sequence:     uint16(sequencer()),
		Service:      cipService_MultipleService,
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
			Service: cipService_Read,
			Size:    byte(len(ioi.Buffer) / 2),
		}
		f := msgCIPIOIFooter{
			Elements: 1,
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
		// TODO: We also need to caculate the response size we expect from the PLC and split
		//       into multiple messages on that also.
		if b.Len() > client.ConnectionSize {
			// TODO: split this read up into mulitple messages.
			return nil, fmt.Errorf("maximum read message size is %d", client.ConnectionSize)
		}
	}

	// right now I'm putting the IOI data into the cip Item, but I suspect it might actually be that the readsequencer is
	// the item's data and the service code actually starts the next portion of the message.  But the item's header length reflects
	// the total data so maybe not.
	reqitems[1] = cipItem{Header: cipItemHeader{ID: cipItem_ConnectedData}}
	reqitems[1].Serialize(ioi_header)
	reqitems[1].Serialize(jump_table)
	reqitems[1].Serialize(&b)

	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, SerializeItems(reqitems))
	if err != nil {
		return nil, err
	}
	_ = hdr

	read_result_header := msgCIPResultHeader{}
	err = binary.Read(data, binary.LittleEndian, &read_result_header)
	if err != nil {
		log.Printf("Problem reading read result header. %v", err)
	}
	items, err := ReadItems(data)
	if err != nil {
		return nil, fmt.Errorf("problem reading items. %w", err)
	}
	if len(items) != 2 {
		return nil, fmt.Errorf("wrong Number of Items. Expected 2 but got %v", len(items))
	}
	ritem := items[1]
	var reply_hdr msgMultiReadResultHeader
	err = binary.Read(&ritem, binary.LittleEndian, &reply_hdr)
	if err != nil {
		return nil, fmt.Errorf("problem reading reply header. %w", err)
	}
	offset_table := make([]uint16, reply_hdr.Reply_Count)
	err = binary.Read(&ritem, binary.LittleEndian, &offset_table)
	if err != nil {
		return nil, fmt.Errorf("problem reading offset table. %w", err)
	}
	rb := ritem.Bytes()
	result_values := make([]interface{}, reply_hdr.Reply_Count)
	for i := 0; i < int(reply_hdr.Reply_Count); i++ {
		offset := offset_table[i] + 10 // offset doesn't start at 0 in the item
		mybytes := bytes.NewBuffer(rb[offset:])
		rhdr := msgMultiReadResult{}
		err = binary.Read(mybytes, binary.LittleEndian, &rhdr)
		if err != nil {
			return nil, fmt.Errorf("problem reading multi result header. %w", err)
		}

		// bit 8 of the service indicates whether it is a response service
		if !rhdr.Service.IsResponse() {
			return nil, fmt.Errorf("wasn't a response service. Got %v", rhdr.Service)
		}
		rhdr.Service = rhdr.Service.UnResponse()
		if rhdr.Status != 0 {
			return nil, fmt.Errorf("problem reading %v. Status %v", tags[i], rhdr.Status)
		}
		if types[i] == CIPTypeBOOL && rhdr.Type != CIPTypeBOOL && iois[i].BitAccess {
			// we have requested a bool from some other type.  Maybe a bit access?
			value, err := readValue(rhdr.Type, mybytes)
			if err != nil {
				return nil, fmt.Errorf("problem reading tag %s: %w", tags[i], err)
			}
			val, err := getBit(rhdr.Type, value, iois[i].BitPosition)
			if err != nil {
				log.Printf("problem reading value for this guy")
				continue
			}
			result_values[i] = val
		} else {
			result_values[i], err = rhdr.Type.readValue(mybytes)
			if err != nil {
				return nil, fmt.Errorf("problem reading tag %s: %w", tags[i], err)
			}
		}

		if verbose {
			log.Printf("Result %d @ %d. %+v. value: %v.\n", i, offset, rhdr, result_values[i])
		}
	}

	return result_values, nil

}

func parseArrayStuct[T GoLogixTypes](dat []byte, elements uint16) ([]T, error) {
	t := make([]T, elements)
	// val should be a byte slice
	b := bytes.NewBuffer(dat)
	for i := 0; i < int(elements); i++ {
		//err := binary.Read(b, binary.LittleEndian, &t[i])
		_, err := Unpack(b, CIPPack{}, &t[i])
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
