package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// WriteMulti writes multiple tags efficiently in a single request using a struct with gologix field tags.
//
// The value parameter must be a struct where each field has a `gologix:"tagname"` tag 
// that specifies which PLC tag to write to. Field types must correspond to the correct 
// CIP types as documented in types.go.
//
// Example:
//   type MyWriteTags struct {
//       IntTag     int16     `gologix:"TestInt"`
//       RealTag    float32   `gologix:"TestReal"`  
//       DintTag    int32     `gologix:"TestDint"`
//       BoolTag    bool      `gologix:"TestBool"`
//       StringTag  string    `gologix:"TestString"`
//   }
//
//   writeValues := MyWriteTags{
//       IntTag:    123,
//       RealTag:   456.78,
//       DintTag:   999888,
//       BoolTag:   true,
//       StringTag: "Hello PLC",
//   }
//   
//   err := client.WriteMulti(writeValues)
//
// For UDT tags, nest the struct data:
//   type MyUDT struct {
//       Field1 int32
//       Field2 float32  
//   }
//   type WriteTags struct {
//       UDTTag MyUDT `gologix:"MyUDTTag"`
//   }
//
// For writing using a map instead of a struct, use WriteMap.
// For writing a single tag, use Write.
//
// WriteMulti automatically handles message splitting for large requests to stay
// within connection size limits.
func (client *Client) WriteMulti(value any) error {
	err := client.checkConnection()
	if err != nil {
		return fmt.Errorf("could not start multi write: %w", err)
	}
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Struct {
		d, err := multi_to_dict(value)
		if err != nil {
			return fmt.Errorf("problem creating keyvalue dict %w", err)
		}
		return client.WriteMap(d)
	}
	return fmt.Errorf("value must be a struct with gologix tags")
}

// Write writes a single value to a single tag in the PLC.
//
// The value parameter must be a Go type that corresponds to the PLC tag's data type.
// Supported type mappings are documented in types.go.
//
// For scalar values and arrays, pass the value directly (not a pointer).
// For UDT (User Defined Type) tags, pass a struct where the struct type name
// matches the UDT name in the PLC, and field types match the UDT field types.
//
// Examples:
//   // Write scalar values
//   err := client.Write("TestInt", int16(123))        // Write to INT tag
//   err := client.Write("TestDint", int32(456789))    // Write to DINT tag  
//   err := client.Write("TestReal", float32(123.45))  // Write to REAL tag
//   err := client.Write("TestBool", true)             // Write to BOOL tag
//   err := client.Write("TestString", "Hello World")  // Write to STRING tag
//
//   // Write arrays
//   intArray := []int32{1, 2, 3, 4, 5}
//   err := client.Write("TestDintArr[0]", intArray)   // Write 5 elements starting at index 0
//
//   // Write UDT struct (struct name must match UDT name)
//   type MyUDT struct {
//       Field1 int32
//       Field2 float32
//   }
//   udtValue := MyUDT{Field1: 100, Field2: 3.14}
//   err := client.Write("MyUDTTag", udtValue)
//
//   // Write to nested UDT field
//   err := client.Write("MyUDTTag.Field1", int32(200))
//
// For writing multiple tags efficiently, use WriteMulti or WriteMap instead.
//
// Returns an error if the connection fails, the tag doesn't exist, there's a type mismatch,
// or the tag is read-only.
func (client *Client) Write(tag string, value any) error {
	err := client.checkConnection()
	if err != nil {
		return fmt.Errorf("could not start write: %w", err)
	}
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Struct {
		return client.write_udt(tag, value)
	}
	return client.write_single(tag, value)
}

// write a single UDT struct to a tag.  The UDT *must* be named the same as the struct type and have the same field types.
// field names don't matter but type names do.  go types will be converted to CIP types as appropriate, but any nested structs
// must be named the same as the UDT on the plc.
func (client *Client) write_udt(tag string, value any) error {
	//service = 0x4D // cipService_Write
	datatype := CIPTypeStruct
	ioi, err := client.newIOI(tag, datatype)
	if err != nil {
		return fmt.Errorf("problem generating IOI. %w", err)
	}
	elements := uint16(1)

	ioi_header := msgCIPIOIHeader{
		Sequence: uint16(sequencer()),
		Service:  CIPService_Write,
		Size:     byte(len(ioi.Buffer) / 2),
	}

	_, typecrc, err := TypeEncode(value)
	if err != nil {
		return fmt.Errorf("problem encoding type. %w", err)
	}

	UDTdata := bytes.NewBuffer([]byte{})

	_, err = Pack(UDTdata, value)
	if err != nil {
		return fmt.Errorf("problem packing data. %w", err)
	}

	reqitems := make([]CIPItem, 2)
	reqitems[0] = newItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)
	reqitems[1] = CIPItem{Header: cipItemHeader{ID: cipItem_ConnectedData}}
	err = reqitems[1].Serialize(ioi_header)
	if err != nil {
		return fmt.Errorf("problem serializing ioi header. %w", err)
	}
	err = reqitems[1].Serialize(ioi.Buffer)
	if err != nil {
		return fmt.Errorf("problem serializing ioi buffer. %w", err)
	}
	err = reqitems[1].Serialize(datatype)
	if err != nil {
		return fmt.Errorf("problem serializing datatype. %w", err)
	}

	err = reqitems[1].Serialize(byte(2))
	if err != nil {
		return fmt.Errorf("problem serializing number of elements. %w", err)
	}

	err = reqitems[1].Serialize(typecrc)
	if err != nil {
		return fmt.Errorf("problem serializing type crc. %w", err)
	}

	err = reqitems[1].Serialize(elements)
	if err != nil {
		return fmt.Errorf("problem serializing elements. %w", err)
	}

	err = reqitems[1].Serialize(UDTdata)
	if err != nil {
		return fmt.Errorf("problem serializing UDT data. %w", err)
	}

	itemdata, err := serializeItems(reqitems)
	if err != nil {
		return err
	}
	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, itemdata)
	if err != nil {
		return err
	}
	if hdr.Status != 0 {
		return fmt.Errorf("got non-success status %d when writing", hdr.Status)
	}
	read_result_header := msgCIPResultHeader{}
	err = binary.Read(data, binary.LittleEndian, &read_result_header)
	if err != nil {
		client.Logger.Warn("Problem reading read result header", "error", err)
	}
	items, err := readItems(data)
	if err != nil {
		return fmt.Errorf("problem reading items from write, %w", err)
	}
	var hdr2 msgWriteResultHeader
	err = items[1].DeSerialize(&hdr2)
	if err != nil {
		return fmt.Errorf("problem deserializing write response header, %w", err)
	}
	if hdr2.Status != CIPStatus_OK {
		extended := uint16(0)
		if hdr2.StatusExtended == 1 {
			err = items[1].DeSerialize(&extended)
			return fmt.Errorf("got status %d:%d:%d instead of 0 in write response.  Problem getting extended status %w",
				hdr2.Status,
				hdr2.StatusExtended,
				extended,
				err)

		}
		return fmt.Errorf("got status %d:%d:%d instead of 0 in write response", hdr2.Status, hdr2.StatusExtended, extended)
	}
	return err
}

// write a single value to a single tag.
func (client *Client) write_single(tag string, value any) error {
	//service = 0x4D // cipService_Write
	datatype, _ := GoVarToCIPType(value)
	ioi, err := client.newIOI(tag, datatype)
	if err != nil {
		return fmt.Errorf("problem generating IOI. %w", err)
	}
	elements := uint16(1)

	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Slice {
		elements = uint16(v.Len())
	}

	ioi_header := msgCIPIOIHeader{
		Sequence: uint16(sequencer()),
		Service:  CIPService_Write,
		Size:     byte(len(ioi.Buffer) / 2),
	}
	ioi_footer := msgCIPWriteIOIFooter{
		DataType: uint16(datatype),
		Elements: elements,
	}

	reqitems := make([]CIPItem, 2)
	reqitems[0] = newItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)
	reqitems[1] = CIPItem{Header: cipItemHeader{ID: cipItem_ConnectedData}}
	err = reqitems[1].Serialize(ioi_header)
	if err != nil {
		return fmt.Errorf("problem serializing ioi header. %w", err)
	}

	err = reqitems[1].Serialize(ioi.Buffer)
	if err != nil {
		return fmt.Errorf("problem serializing ioi buffer. %w", err)
	}

	err = reqitems[1].Serialize(ioi_footer)
	if err != nil {
		return fmt.Errorf("problem serializing ioi footer. %w", err)
	}

	err = reqitems[1].Serialize(value)
	if err != nil {
		return fmt.Errorf("problem serializing value. %w", err)
	}

	itemdata, err := serializeItems(reqitems)
	if err != nil {
		return err
	}
	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, itemdata)
	if err != nil {
		return err
	}
	if hdr.Status != 0 {
		return fmt.Errorf("got non-success status %d when writing", hdr.Status)
	}
	read_result_header := msgCIPResultHeader{}
	err = binary.Read(data, binary.LittleEndian, &read_result_header)
	if err != nil {
		client.Logger.Warn("Problem reading read result header", "error", err)
	}
	items, err := readItems(data)
	if err != nil {
		return fmt.Errorf("problem reading items from write, %w", err)
	}
	var hdr2 msgWriteResultHeader
	err = items[1].DeSerialize(&hdr2)
	if err != nil {
		return fmt.Errorf("problem deserializing write response header, %w", err)
	}
	if hdr2.Status != CIPStatus_OK {
		extended := uint16(0)
		if hdr2.StatusExtended == 1 {
			err = items[1].DeSerialize(&extended)
			return fmt.Errorf("got status %d:%d:%d instead of 0 in write response.  Problem getting extended status %w",
				hdr2.Status,
				hdr2.StatusExtended,
				extended,
				err)

		}
		return fmt.Errorf("got status %d:%d:%d instead of 0 in write response", hdr2.Status, hdr2.StatusExtended, extended)
	}
	return err
}

type msgWriteResultHeader struct {
	SequenceCount  uint16
	Service        CIPService
	Reserved       byte
	Status         CIPStatus
	StatusExtended byte
}

type msgUnconnWriteResultHeader struct {
	Service        CIPService
	Reserved       byte
	Status         CIPStatus
	StatusExtended byte
}
