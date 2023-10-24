package gologix

import (
	"encoding/binary"
	"fmt"
	"log"
	"reflect"
)

// The write equivalent to ReadMulti.  value should be a struct where each field has a tag of the form `gologix:"tagname"` that maps
// what tag in the controller it corresponds to
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
		d, err := udt_to_dict(tag, value)
		if err != nil {
			return fmt.Errorf("problem creating keyvalue dict %w", err)
		}
		return client.WriteMap(d)
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
		Sequence: uint16(sequencer()),
		Service:  CIPService_Write,
		Size:     byte(len(ioi.Buffer) / 2),
	}
	ioi_footer := msgCIPWriteIOIFooter{
		DataType: uint16(datatype),
		Elements: elements,
	}

	reqitems := make([]CIPItem, 2)
	reqitems[0] = NewItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)
	reqitems[1] = CIPItem{Header: cipItemHeader{ID: cipItem_ConnectedData}}
	reqitems[1].Serialize(ioi_header)
	reqitems[1].Serialize(ioi.Buffer)
	reqitems[1].Serialize(ioi_footer)
	reqitems[1].Serialize(value)

	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, SerializeItems(reqitems))
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
	Status         byte
	StatusExtended byte
}

type msgUnconnWriteResultHeader struct {
	Service        CIPService
	Reserved       byte
	Status         byte
	StatusExtended byte
}
