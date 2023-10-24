package gologix

import (
	"encoding/binary"
	"fmt"
	"log"
)

type cipAttributeResponseHdr struct {
	SequenceCount   uint16
	ResponseService CIPService
	_               byte
	Status          uint16
}

func (client *Client) GetAttrSingle(class CIPClass, instance CIPInstance, attr CIPAttribute) (*CIPItem, error) {

	err := client.checkConnection()
	if err != nil {
		return nil, fmt.Errorf("could not start single read: %w", err)
	}

	reqitems := make([]CIPItem, 2)
	reqitems[0] = NewItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: uint16(sequencer()),
		Service:       CIPService_GetAttributeSingle,
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
		return nil, err
	}

	var resphdr cipAttributeResponseHdr
	items[1].DeSerialize(&resphdr)
	//dat := make([]byte, hdr.Length-26)
	//items[1].DeSerialize(&dat)
	return &items[1], nil
}

// this is specifically the response for a GetAttrList service on a
// controller info object with requested attributes of 1,2,3,4,10
type msgGetControllerPropList struct {
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
// don't exactly know what they are, just going off of what 1756-pm020_-en-p.pdf says on page 51
func (client *Client) GetControllerPropList() (msgGetControllerPropList, error) {

	item, err := client.GetAttrList(CipObject_ControllerInfo, 1, 1, 2, 3, 4, 10)
	if err != nil {
		return msgGetControllerPropList{}, err
	}

	result := msgGetControllerPropList{}
	err = item.DeSerialize(&result)
	if err != nil {
		return msgGetControllerPropList{}, fmt.Errorf("couldn't read data. %w", err)
	}
	if verbose {
		log.Printf("Result: %+v", result)
	}

	return result, nil
}

// read multiple attributes.
// This function returns a cip item containing all the attributes.  They should be in the following format:
//
//		Attr1_ID     uint16
//		Attr1_Status uint16
//		Attr1_Value  []byte
//	    Attr2_ID     uint16
//		Attr2_Status uint16
//		Attr2_Value  []byte
//	    ...
//	    AttrN_ID     uint16
//		AttrN_Status uint16
//		AttrN_Value  []byte
//
// CIP expects you to know the data type for each of the values so you'll have to parse it manually one at a time to figure out
// where the subsequent fields are in the binary data.
// If there are no variable-length fields in the attributes you are getting, the best way is to create the equivalent
// struct as above with the proper types for the AttrX_Value instead of []byte and do an item.Serialize(&InstanceOfMyType)
func (client *Client) GetAttrList(class CIPClass, instance CIPInstance, attrs ...CIPAttribute) (*CIPItem, error) {

	reqitems := make([]CIPItem, 2)
	//reqitems[0] = cipItem{Header: cipItemHeader{ID: cipItem_Null}}
	reqitems[0] = NewItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	p, err := Serialize(
		class, instance,
	)
	if err != nil {
		return nil, fmt.Errorf("couldn't build path. %w", err)
	}

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: uint16(sequencer()),
		Service:       CIPService_GetAttributeList,
		PathLength:    byte(p.Len() / 2),
	}

	reqitems[1] = NewItem(cipItem_ConnectedData, readmsg)
	reqitems[1].Serialize(p.Bytes())
	reqitems[1].Serialize(uint16(len(attrs)))
	for i := range attrs {
		reqitems[1].Serialize(uint16(attrs[i]))
	}
	reqitems[1].Serialize(byte(1))
	reqitems[1].Serialize(byte(0))
	reqitems[1].Serialize(uint16(1))

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
		return nil, err
	}

	var resphdr cipAttributeResponseHdr
	items[1].DeSerialize(&resphdr)
	//dat := make([]byte, hdr.Length-26)
	//items[1].DeSerialize(&dat)

	// There is a count before there is any result data - we'll remove that here.
	// If needed, the item can always be Reset() to get all the data from the start.
	_, err = items[1].Int16() // Count
	if err != nil {
		return nil, err
	}

	return &items[1], nil

}

// Generic CIP Message
//
// CIP expects you to know the data type for each of the values so you'll have to parse the resulting CIPItem yourself.
//
// If there are no variable-length fields in the response data you are expecting, the best way may be to create the equivalent
// struct with the proper types and do an item.Serialize(&InstanceOfMyType)
func (client *Client) GenericCIPMessage(service CIPService, class CIPClass, instance CIPInstance, msg_data []byte) (*CIPItem, error) {

	reqitems := make([]CIPItem, 2)
	//reqitems[0] = cipItem{Header: cipItemHeader{ID: cipItem_Null}}
	reqitems[0] = NewItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	p, err := Serialize(
		class, instance,
	)
	if err != nil {
		return nil, fmt.Errorf("couldn't build path. %w", err)
	}

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: uint16(sequencer()),
		Service:       service,
		PathLength:    byte(p.Len() / 2),
	}

	reqitems[1] = NewItem(cipItem_ConnectedData, readmsg)
	reqitems[1].Serialize(p.Bytes())
	reqitems[1].Serialize(msg_data)

	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, SerializeItems(reqitems))

	if err != nil {
		return nil, err
	}
	_ = hdr

	read_result_header := msgCIPResultHeader{}
	err = binary.Read(data, binary.LittleEndian, &read_result_header)
	if err != nil {
		return nil, fmt.Errorf("problpm reading read result header: %w", err)
	}
	items, err := ReadItems(data)
	if err != nil {
		return nil, fmt.Errorf("problem reading items: %w", err)
	}

	// There is a count before there is any result data - we'll remove that here.
	// If needed, the item can always be Reset() to get all the data from the start.
	_, err = items[1].Int16() // Count
	if err != nil {
		return nil, fmt.Errorf("problem reading sequence counter: %w", err)
	}
	s_response, err := items[1].Int16()
	if err != nil {
		return nil, fmt.Errorf("problem reading service response: %w", err)
	}

	if CIPService(s_response).UnResponse() != service {
		return nil, fmt.Errorf("expected service response 0x%X but got 0x%X", service, CIPService(s_response).UnResponse())
	}
	status, err := items[1].Int16()
	if err != nil {
		return nil, fmt.Errorf("problem getting resposne status: %w", err)
	}
	if status != 0 {
		return &items[1], fmt.Errorf("got status of 0x%X instead of 0", status)
	}

	return &items[1], nil

}
