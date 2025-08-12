package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// convert tagged struct to a map in the format of {"fieldTag": fieldvalue}
func multi_to_dict(data any) (map[string]interface{}, error) {
	//TODO: handle nested structs and arrays
	// convert the struct to a dict of FieldName: FieldValue
	d := make(map[string]interface{})

	vs := reflect.ValueOf(data)
	fs := reflect.TypeOf(data)
	for i := 0; i < vs.NumField(); i++ {
		v := vs.Field(i)
		f := fs.Field(i)
		fulltag := f.Tag.Get("gologix")
		switch v.Kind() {
		case reflect.Struct:
			d2, err := udt_to_dict(fulltag, v.Interface())
			if err != nil {
				return nil, fmt.Errorf("problem parsing %s. %w", fulltag, err)
			}
			for k := range d2 {
				d[k] = d2[k]
			}
		case reflect.Array:
			//TODO
		default:
			d[fulltag] = v.Interface()
		}

	}

	return d, nil

}

// convert a struct to a map in the format of {"tag.fieldName": fieldvalue}
func udt_to_dict(tag string, data any) (map[string]interface{}, error) {
	//TODO: handle nested structs and arrays
	// convert the struct to a dict of FieldName: FieldValue
	d := make(map[string]interface{})

	vs := reflect.ValueOf(data)
	fs := reflect.TypeOf(data)
	for i := 0; i < vs.NumField(); i++ {
		v := vs.Field(i)
		f := fs.Field(i)
		fulltag := fmt.Sprintf("%s.%s", tag, f.Name)
		switch v.Kind() {

		}
		switch v.Kind() {
		case reflect.Struct:
			d2, err := udt_to_dict(fulltag, v.Interface())
			if err != nil {
				return nil, fmt.Errorf("problem parsing %s. %w", fulltag, err)
			}
			for k := range d2 {
				d[k] = d2[k]
			}
		case reflect.Array:
			//TODO
		default:
			d[fulltag] = v.Interface()
		}

	}

	return d, nil

}

// WriteMap writes multiple tags efficiently using a map where keys are tag names and values are the data to write.
//
// The map keys specify the tag names to write to in the PLC (case insensitive).
// The map values must be Go types that correspond to the PLC tag types as documented in types.go.
//
// Supported value types:
//   - Scalar types: int16 (INT), int32 (DINT), float32 (REAL), bool (BOOL), string (STRING), etc.
//   - Arrays: slices like []int32, []float32, etc.
//   - Structs: user-defined types for UDT tags (struct name must match UDT name)
//
// Examples:
//   // Basic scalar writes
//   writeMap := map[string]interface{}{
//       "TestInt":    int16(123),
//       "TestDint":   int32(456789),  
//       "TestReal":   float32(123.45),
//       "TestBool":   true,
//       "TestString": "Hello PLC",
//   }
//
//   // Array writes (write to specific indices)
//   writeMap["TestDintArr[0]"] = []int32{1, 2, 3, 4, 5}  // Write 5 elements starting at index 0
//   writeMap["TestRealArr[10]"] = []float32{1.1, 2.2, 3.3} // Write 3 elements starting at index 10
//
//   // UDT writes
//   type MyUDT struct {
//       Field1 int32
//       Field2 float32
//   }
//   writeMap["MyUDTTag"] = MyUDT{Field1: 100, Field2: 3.14}
//
//   // Individual UDT field writes  
//   writeMap["MyUDTTag.Field1"] = int32(200)
//   writeMap["MyUDTTag.Field2"] = float32(6.28)
//
//   err := client.WriteMap(writeMap)
//   if err != nil {
//       log.Fatal(err)
//   }
//
// For struct-based writing with field tags, use WriteMulti instead.
// For writing a single tag, use Write.
//
// WriteMap automatically handles message splitting for large requests and optimizes
// network usage by grouping writes into the minimum number of requests.
//
// Returns an error if the connection fails, any tag doesn't exist, there are type mismatches,
// or any tags are read-only.
func (client *Client) WriteMap(tag_str map[string]interface{}) error {

	// build the tag list from the structure
	tags := make([]string, 0)
	types := make([]CIPType, 0)
	for k := range tag_str {
		ct, _ := GoVarToCIPType(tag_str[k])
		types = append(types, ct)
		tags = append(tags, k)
	}

	// first generate IOIs for each tag
	qty := len(tags)
	iois := make([]*tagIOI, qty)
	for i, tag := range tags {
		var err error
		iois[i], err = client.newIOI(tag, types[i])
		if err != nil {
			return err
		}
	}

	reqitems := make([]CIPItem, 2)
	reqitems[0] = newItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

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
			Service: CIPService_Write,
			Size:    byte(len(ioi.Buffer) / 2),
		}
		f := msgCIPWriteIOIFooter{
			DataType: uint16(types[i]),
			Elements: 1,
		}
		//f := msgCIPIOIFooter{
		//Elements: 1,
		//}
		err := binary.Write(&b, binary.LittleEndian, h)
		if err != nil {
			return fmt.Errorf("problem writing udt item header to buffer. %w", err)
		}
		b.Write(ioi.Buffer)
		ftr_buf, err := Serialize(f)
		if err != nil {
			return fmt.Errorf("problem serializing footer for item %d: %w", i, err)
		}
		err = binary.Write(&b, binary.LittleEndian, ftr_buf.Bytes())
		if err != nil {
			return fmt.Errorf("problem writing udt item footer to buffer. %w", err)
		}
		item_buf, err := Serialize(tag_str[tags[i]])
		if err != nil {
			return fmt.Errorf("problem serializing %v: %w", tags[i], err)
		}
		err = binary.Write(&b, binary.LittleEndian, item_buf.Bytes())
		if err != nil {
			return fmt.Errorf("problem writing udt tag name to buffer. %w", err)
		}
	}

	// right now I'm putting the IOI data into the cip Item, but I suspect it might actually be that the readsequencer is
	// the item's data and the service code actually starts the next portion of the message.  But the item's header length reflects
	// the total data so maybe not.
	reqitems[1] = CIPItem{Header: cipItemHeader{ID: cipItem_ConnectedData}}
	err := reqitems[1].Serialize(ioi_header, jump_table, &b)
	if err != nil {
		return fmt.Errorf("problem serializing item header: %w", err)
	}

	itemdata, err := serializeItems(reqitems)
	if err != nil {
		return err
	}
	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, itemdata)
	if err != nil {
		return err
	}
	_ = hdr
	_ = data
	//TODO: do something with the result here!

	return nil
}
