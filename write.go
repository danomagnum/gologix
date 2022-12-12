package gologix

import (
	"encoding/binary"
	"fmt"
	"log"
	"reflect"
)

// write a single value to a single tag.
func (client *Client) Write(tag string, value any) error {
	//service = 0x4D // CIPService_Write
	datatype := GoVarToCIPType(value)
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
		Sequence: client.Sequencer(),
		Service:  CIPService_Write,
		Size:     byte(len(ioi.Buffer) / 2),
	}
	ioi_footer := msgCIPWriteIOIFooter{
		DataType: uint16(datatype),
		Elements: elements,
	}

	reqitems := make([]CIPItem, 2)
	reqitems[0] = NewItem(CIPItem_ConnectionAddress, &client.OTNetworkConnectionID)
	reqitems[1] = CIPItem{Header: CIPItemHeader{ID: CIPItem_ConnectedData}}
	reqitems[1].Marshal(ioi_header)
	reqitems[1].Marshal(ioi.Buffer)
	reqitems[1].Marshal(ioi_footer)
	reqitems[1].Marshal(value)

	hdr, data, err := client.send_recv_data(CIPCommandSendUnitData, MarshalItems(reqitems))
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
