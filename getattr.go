package gologix

import (
	"encoding/binary"
	"fmt"
	"log"
)

func (client *Client) GetAttrSingle(class CIPClass, instance CIPInstance, attr CIPAttribute) ([]CIPItem, error) {

	err := client.checkConnection()
	if err != nil {
		return nil, fmt.Errorf("could not start single read: %w", err)
	}

	reqitems := make([]CIPItem, 2)
	reqitems[0] = NewItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: uint16(sequencer()),
		Service:       cipService_GetAttributeSingle,
		PathLength:    3,
	}
	// setup item
	reqitems[1] = NewItem(cipItem_ConnectedData, readmsg)
	// add path
	reqitems[1].Serialize(class)
	reqitems[1].Serialize(instance)
	reqitems[1].Serialize(attr)

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
		return []CIPItem{}, err
	}

	return items, nil
}
