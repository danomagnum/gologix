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

func (client *Client) GetControllerPropList() (msgGetControllerPropList, error) {

	reqitems := make([]CIPItem, 2)
	//reqitems[0] = CIPItem{Header: CIPItemHeader{ID: CIPItem_Null}}
	reqitems[0] = NewItem(CIPItem_ConnectionAddress, &client.OTNetworkConnectionID)

	p, err := Serialize(
		CIPObject_ControllerInfo, CIPInstance(1),
		//CIPObject_Symbol, CIPInstance(start_instance),
	)
	if err != nil {
		return msgGetControllerPropList{}, fmt.Errorf("couldn't build path. %w", err)
	}

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: client.Sequencer(),
		Service:       CIPService_GetAttributeList,
		PathLength:    byte(p.Len() / 2),
	}

	reqitems[1] = NewItem(CIPItem_ConnectedData, readmsg)
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

	hdr, data, err := client.send_recv_data(CIPCommandSendUnitData, MarshalItems(reqitems))
	if err != nil {
		return msgGetControllerPropList{}, err
	}
	_ = hdr
	_ = data
	//data_hdr := ListInstanceHeader{}
	//binary.Read(data, binary.LittleEndian, &data_hdr)

	// first six bytes are zero.
	padding := make([]byte, 6)
	data.Read(padding)

	resp_items, err := ReadItems(data)
	if err != nil {
		return msgGetControllerPropList{}, fmt.Errorf("couldn't parse items %w", err)
	}

	// get ready to read tag info from item 1 data
	data2 := bytes.NewBuffer(resp_items[1].Data)

	result := msgGetControllerPropList{}
	binary.Read(data2, binary.LittleEndian, &result)
	log.Printf("Result: %+v\n\n", result)

	return result, nil
}
