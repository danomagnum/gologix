package gologix

import (
	"encoding/binary"
	"fmt"
	"log"
	"reflect"
)

func (client *Client) WriteMulti(value any) error {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Struct {
		d, err := multi_to_dict(value)
		if err != nil {
			return fmt.Errorf("problem creating keyvalue dict %w", err)
		}
		return client.writeDict(d)
	}
	return fmt.Errorf("value must be a struct with gologix tags")
}

// write a single value to a single tag.
func (client *Client) Write(tag string, value any) error {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Struct {
		d, err := udt_to_dict(tag, value)
		if err != nil {
			return fmt.Errorf("problem creating keyvalue dict %w", err)
		}
		return client.writeDict(d)
	}
	return client.write_single(tag, value)
}

// write a single value to a single tag.
func (client *Client) write_single(tag string, value any) error {
	//service = 0x4D // cipService_Write
	datatype := GoVarToCIPType(value)
	ioi, err := client.NewIOI(tag, datatype)
	if err != nil {
		return fmt.Errorf("problem generating IOI. %w", err)
	}
	elements := uint16(1)

	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Slice {
		elements = uint16(v.Len())
	}

	ioi_header := msgCIPIOIHeader{
		Sequence: client.Sequencer(),
		Service:  cipService_Write,
		Size:     byte(len(ioi.Buffer) / 2),
	}
	ioi_footer := msgCIPWriteIOIFooter{
		DataType: uint16(datatype),
		Elements: elements,
	}

	reqitems := make([]cipItem, 2)
	reqitems[0] = NewItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)
	reqitems[1] = cipItem{Header: cipItemHeader{ID: cipItem_ConnectedData}}
	reqitems[1].Marshal(ioi_header)
	reqitems[1].Marshal(ioi.Buffer)
	reqitems[1].Marshal(ioi_footer)
	reqitems[1].Marshal(value)

	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, MarshalItems(reqitems))
	if err != nil {
		return err
	}
	if hdr.Status != 0 {
		return fmt.Errorf("got non-success status %d when writing", hdr.Status)
	}
	read_result_header := msgCIPResultHeader{}
	err = binary.Read(data, binary.LittleEndian, &read_result_header)
	if err != nil {
		log.Printf("Problem reading read result header. %v", err)
	}
	items, err := ReadItems(data)
	if err != nil {
		return fmt.Errorf("problem reading items from write, %w", err)
	}
	var hdr2 msgWriteResultHeader
	err = items[1].Unmarshal(&hdr2)
	if err != nil {
		return fmt.Errorf("problem unmarshaling write response header, %w", err)
	}
	if hdr2.Status != 0 {
		extended := uint16(0)
		if hdr2.StatusExtended == 1 {
			items[1].Unmarshal(&extended)
		}
		return fmt.Errorf("got status %d:%d:%d instead of 0 in write response", hdr2.Status, hdr2.StatusExtended, extended)
	}
	return err
}

type msgWriteResultHeader struct {
	SequenceCount  uint16
	Service        CIPService
	Reserved       byte
	Status         byte
	StatusExtended byte
}

type msgUnconnWriteResultHeader struct {
	Service        CIPService
	Reserved       byte
	Status         byte
	StatusExtended byte
}
