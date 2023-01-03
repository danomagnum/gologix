package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

// this is specifically the response for a GetAttrList service on a
// controller info object with requested attributes of 1,2,3,4,10
type msgGetControllerPropList struct {
	SequenceCount   uint16
	Service         CIPService
	Reserved        byte
	Status          byte
	Status_extended byte
	Count           uint16

	Attr1_ID     uint16
	Attr1_Status uint16
	Attr1        uint16

	Attr2_ID     uint16
	Attr2_Status uint16
	Attr2        uint16

	Attr3_ID     uint16
	Attr3_Status uint16
	Attr3        uint32

	Attr4_ID     uint16
	Attr4_Status uint16
	Attr4        uint32

	Attr5_ID     uint16
	Attr5_Status uint16
	Attr5        uint32
}

func (old msgGetControllerPropList) Match(new msgGetControllerPropList) bool {
	if new.Attr1 != old.Attr1 || new.Attr1_Status != old.Attr1_Status {
		return false
	}
	if new.Attr2 != old.Attr2 || new.Attr2_Status != old.Attr2_Status {
		return false
	}
	if new.Attr3 != old.Attr3 || new.Attr3_Status != old.Attr3_Status {
		return false
	}
	if new.Attr4 != old.Attr4 || new.Attr4_Status != old.Attr4_Status {
		return false
	}
	if new.Attr5 != old.Attr5 || new.Attr5_Status != old.Attr5_Status {
		return false
	}
	return true
}

// read the general controller information.
// these properties indicate if the controller has been modified.  Could indicate a logic change or a tag was added or removed.
func (client *Client) GetControllerPropList() (msgGetControllerPropList, error) {

	reqitems := make([]cipItem, 2)
	//reqitems[0] = cipItem{Header: cipItemHeader{ID: cipItem_Null}}
	reqitems[0] = NewItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	p, err := Serialize(
		cipObject_ControllerInfo, CIPInstance(1),
		//cipObject_Symbol, cipInstance(start_instance),
	)
	if err != nil {
		return msgGetControllerPropList{}, fmt.Errorf("couldn't build path. %w", err)
	}

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: client.Sequencer(),
		Service:       cipService_GetAttributeList,
		PathLength:    byte(p.Len() / 2),
	}

	reqitems[1] = NewItem(cipItem_ConnectedData, readmsg)
	reqitems[1].Marshal(p.Bytes())
	number_of_attr_to_receive := 5
	reqitems[1].Marshal([]uint16{
		uint16(number_of_attr_to_receive),
		uint16(1),
		uint16(2),
		uint16(3),
		uint16(4),
		uint16(10),
	})
	reqitems[1].Marshal(byte(1))
	reqitems[1].Marshal(byte(0))
	reqitems[1].Marshal(uint16(1))

	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, MarshalItems(reqitems))
	if err != nil {
		return msgGetControllerPropList{}, err
	}
	_ = hdr
	_ = data
	//data_hdr := ListInstanceHeader{}
	//binary.Read(data, binary.LittleEndian, &data_hdr)

	// first six bytes are zero.
	padding := make([]byte, 6)
	_, err = data.Read(padding)
	if err != nil {
		return msgGetControllerPropList{}, fmt.Errorf("couldn't read data. %w", err)
	}

	resp_items, err := ReadItems(data)
	if err != nil {
		return msgGetControllerPropList{}, fmt.Errorf("couldn't parse items %w", err)
	}

	// get ready to read tag info from item 1 data
	data2 := bytes.NewBuffer(resp_items[1].Data)

	result := msgGetControllerPropList{}
	err = binary.Read(data2, binary.LittleEndian, &result)
	if err != nil {
		return msgGetControllerPropList{}, fmt.Errorf("couldn't read data. %w", err)
	}
	if verbose {
		log.Printf("Result: %+v\n\n", result)
	}

	return result, nil
}
