package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// The write equivalent to ReadMulti.  value should be a struct where each field has a tag of the form `gologix:"tagname"` that maps
// what tag in the controller it corresponds to.
//
// To write multiple tags without creating a tagged struct, look at WriteMap()
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

// write a single value to a single tag.
// the type of value must correspond to the type of tag in the controller
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
	reqitems[1].Serialize(ioi_header)
	reqitems[1].Serialize(ioi.Buffer)
	reqitems[1].Serialize(datatype)
	reqitems[1].Serialize(byte(2))
	reqitems[1].Serialize(typecrc)
	reqitems[1].Serialize(elements)
	reqitems[1].Serialize(UDTdata)

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
	if hdr2.Status != 0 {
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
	reqitems[1].Serialize(ioi_header)
	reqitems[1].Serialize(ioi.Buffer)
	reqitems[1].Serialize(ioi_footer)
	reqitems[1].Serialize(value)

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
	if hdr2.Status != 0 {
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
