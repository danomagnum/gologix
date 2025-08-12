package gologix

import (
	"encoding/binary"
	"fmt"
)

type cipAttributeResponseHdr struct {
	SequenceCount   uint16
	ResponseService CIPService
	_               byte
	Status          uint16
}

// GetAttrSingle retrieves a single attribute value from a specified CIP object.
//
// This function executes a CIP GetAttributeSingle service to read a specific attribute
// from a target object/instance combination. It provides direct access to device
// attributes that may not be available through standard tag operations.
//
// Parameters:
//   - class: The CIP object class (e.g., CipObject_Identity, CipObject_ControllerInfo)
//   - instance: The instance number of the object (typically 1 for singleton objects)
//   - attr: The attribute number to read from the specified object instance
//
// Returns a CIPItem containing the raw attribute data that must be deserialized
// based on the expected data type for the specific attribute being read.
//
// The caller must know the expected data type and structure of the attribute
// response in advance to properly parse the returned data.
//
// Common Use Cases:
//
//  1. Reading device vendor information:
//     // Get Vendor ID (attribute 1) from Identity Object
//     result, err := client.GetAttrSingle(
//     gologix.CipObject_Identity,
//     gologix.CIPInstance(1),
//     gologix.CIPAttribute(1),
//     )
//     if err != nil {
//     log.Fatal(err)
//     }
//
//     var vendorID uint16
//     err = result.DeSerialize(&vendorID)
//     fmt.Printf("Vendor ID: %d\n", vendorID)
//
//  2. Reading device product information:
//     // Get Product Code (attribute 3) from Identity Object
//     result, err := client.GetAttrSingle(
//     gologix.CipObject_Identity,
//     gologix.CIPInstance(1),
//     gologix.CIPAttribute(3),
//     )
//     if err != nil {
//     log.Fatal(err)
//     }
//
//     var productCode uint16
//     err = result.DeSerialize(&productCode)
//     fmt.Printf("Product Code: %d\n", productCode)
//
//  3. Reading controller-specific attributes:
//     // Get controller information
//     result, err := client.GetAttrSingle(
//     gologix.CipObject_ControllerInfo,
//     gologix.CIPInstance(1),
//     gologix.CIPAttribute(1), // Controller property
//     )
//     if err != nil {
//     log.Fatal(err)
//     }
//
//     var controllerProp uint16
//     err = result.DeSerialize(&controllerProp)
//
//  4. Reading current time from controller:
//     // Get microseconds since Unix epoch from Time Object
//     result, err := client.GetAttrSingle(
//     gologix.CipObject_TIME,
//     gologix.CIPInstance(1),
//     gologix.CIPAttribute(11), // Time attribute
//     )
//     if err != nil {
//     log.Fatal(err)
//     }
//
//     var timeUs int64
//     err = result.DeSerialize(&timeUs)
//     currentTime := time.UnixMicro(timeUs)
//
// Error Handling:
// The function returns errors for communication failures, connection issues,
// or invalid object/attribute combinations. Some attributes may require
// specific access privileges or may not be supported on all device types.
//
// Data Type Considerations:
// Different attributes return different data types (uint16, uint32, strings, etc.).
// Consult the device documentation or CIP specification for attribute data types.
func (client *Client) GetAttrSingle(class CIPClass, instance CIPInstance, attr CIPAttribute) (*CIPItem, error) {

	err := client.checkConnection()
	if err != nil {
		return nil, fmt.Errorf("could not start single read: %w", err)
	}

	reqitems := make([]CIPItem, 2)
	reqitems[0] = newItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)
	pl := class.Len() + instance.Len() + attr.Len()
	pl = pl / 2

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: uint16(sequencer()),
		Service:       CIPService_GetAttributeSingle,
		PathLength:    byte(pl),
	}
	// setup item
	reqitems[1] = newItem(cipItem_ConnectedData, readmsg)
	// add path
	err = reqitems[1].Serialize(class)
	if err != nil {
		return nil, fmt.Errorf("could not serialize class: %w", err)
	}
	err = reqitems[1].Serialize(instance)
	if err != nil {
		return nil, fmt.Errorf("could not serialize instance: %w", err)
	}
	err = reqitems[1].Serialize(attr)
	if err != nil {
		return nil, fmt.Errorf("could not serialize attribute: %w", err)
	}

	itemdata, err := serializeItems(reqitems)
	if err != nil {
		return nil, err
	}
	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, itemdata)
	if err != nil {
		return nil, err
	}
	_ = hdr

	read_result_header := msgCIPResultHeader{}
	err = binary.Read(data, binary.LittleEndian, &read_result_header)
	if err != nil {
		client.Logger.Warn("Problem reading read result header.", "error", err)
	}
	items, err := readItems(data)
	if err != nil {
		client.Logger.Warn("Problem reading items.", "error", err)
		return nil, err
	}

	var resphdr cipAttributeResponseHdr
	err = items[1].DeSerialize(&resphdr)
	if err != nil {
		return nil, fmt.Errorf("couldn't read response header. %w", err)
	}
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
	client.Logger.Debug("Controller Prop", "result", result)

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
	reqItems := make([]CIPItem, 2)
	reqItems[0] = newItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	path, err := Serialize(class, instance)
	if err != nil {
		return nil, fmt.Errorf("could not build path. %w", err)
	}

	readMsg := msgCIPConnectedServiceReq{
		SequenceCount: uint16(client.sequenceNumber.Add(1)),
		Service:       CIPService_GetAttributeList,
		PathLength:    byte(path.Len() / 2),
	}

	reqItems[1] = newItem(cipItem_ConnectedData, readMsg)
	err = reqItems[1].Serialize(path.Bytes())
	if err != nil {
		return nil, fmt.Errorf("could not serialize path. %w", err)
	}
	err = reqItems[1].Serialize(uint16(len(attrs)))
	if err != nil {
		return nil, fmt.Errorf("could not serialize attribute count. %w", err)
	}
	for i := range attrs {
		err = reqItems[1].Serialize(uint16(attrs[i]))
		if err != nil {
			return nil, fmt.Errorf("could not serialize attribute %d. %w", i, err)
		}
	}

	itemData, err := serializeItems(reqItems)
	if err != nil {
		return nil, err
	}
	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, itemData)

	if err != nil {
		return nil, err
	}
	_ = hdr

	read_result_header := msgCIPResultHeader{}
	err = binary.Read(data, binary.LittleEndian, &read_result_header)
	if err != nil {
		client.Logger.Warn("Problem reading read result header.", "error", err)
		return nil, err
	}
	items, err := readItems(data)
	if err != nil {
		client.Logger.Warn("Problem reading items.", "error", err)
		return nil, err
	}

	var respHdr cipAttributeResponseHdr
	err = items[1].DeSerialize(&respHdr)
	if err != nil {
		return nil, fmt.Errorf("couldn't read response header. %w", err)
	}
	if respHdr.Status != uint16(CIPStatus_OK) {
		return &items[1], fmt.Errorf("response header has status 0x%X (%v)", respHdr.Status, CIPStatus(respHdr.Status))
	}

	// There is a count before there is any result data - we'll remove that here.
	// If needed, the item can always be Reset() to get all the data from the start.
	_, err = items[1].Int16() // Count
	if err != nil {
		return nil, err
	}

	return &items[1], nil

}

// GenericCIPMessage sends a custom CIP (Common Industrial Protocol) message to the device.
//
// This is an advanced function that provides direct access to the CIP messaging layer,
// allowing for custom commands beyond the standard tag read/write operations. It enables
// interaction with any CIP object, service, and attributes on EIP devices.
//
// Parameters:
//   - service: The CIP service code to execute (e.g., CIPService_GetAttributeSingle, CIPService_Stop)
//   - path: Serialized CIP path specifying the target object/instance/attribute
//   - msg_data: Additional data payload for the service (empty []byte{} if not needed)
//
// Returns a CIPItem containing the response data that must be manually parsed based on
// the expected response format for the specific service and object being accessed.
//
// Path Construction:
// The path parameter must be constructed using the Serialize() function:
//
//	// Target a specific object and instance
//	path, err := gologix.Serialize(gologix.CipObject_TIME, gologix.CIPInstance(1))
//
//	// Target a specific attribute of an object
//	path, err := gologix.Serialize(
//	    gologix.CIPClass(1),        // Identity Object
//	    gologix.CIPInstance(1),     // Instance 1
//	    gologix.CIPAttribute(1),    // Vendor ID attribute
//	)
//
// Response Parsing:
// The returned CIPItem must be deserialized based on the expected response structure:
//
//	type TimeResponse struct {
//	    AttrCount int16
//	    AttrID    uint16
//	    Status    uint16
//	    Usecs     int64  // Microseconds since Unix epoch
//	}
//
//	response := TimeResponse{}
//	err = result.DeSerialize(&response)
//
// Common Use Cases:
//
//  1. Reading controller time:
//     path, _ := gologix.Serialize(gologix.CipObject_TIME, gologix.CIPInstance(1))
//     result, err := client.GenericCIPMessage(
//     gologix.CIPService_GetAttributeList,
//     path.Bytes(),
//     []byte{0x01, 0x00, 0x0B, 0x00}, // Request attribute 11 (time)
//     )
//
//  2. Accessing controller run mode:
//     path, _ := gologix.Serialize(gologix.CipObject_RunMode, gologix.CIPInstance(1))
//     result, err := client.GenericCIPMessage(
//     gologix.CIPService_GetAttributeSingle,
//     path.Bytes(),
//     []byte{},
//     )
//
//  3. Controller control operations (requires elevated privileges):
//     path, _ := gologix.Serialize(gologix.CipObject_RunMode, gologix.CIPInstance(1))
//     result, err := client.GenericCIPMessage(
//     gologix.CIPService_Stop,  // Will likely fail without proper privileges
//     path.Bytes(),
//     []byte{},
//     )
//
//  4. Reading device identity information:
//     path, _ := gologix.Serialize(
//     gologix.CipObject_Identity,
//     gologix.CIPInstance(1),
//     gologix.CIPAttribute(1), // Vendor ID
//     )
//     result, err := client.GenericCIPMessage(
//     gologix.CIPService_GetAttributeSingle,
//     path.Bytes(),
//     []byte{},
//     )
//
// Error Handling:
// The function returns errors for communication issues, but CIP-level errors
// (like access denied, object not found) may be embedded in the response data.
// Always check both the error return and parse the response status codes.
//
// Security Considerations:
// Many advanced operations require elevated privileges. Operations like starting/stopping
// the controller, modifying safety parameters, or accessing security objects will
// typically return privilege violation errors (0x0F) unless proper authentication
// and authorization have been established.
// Authorization is a secret that only rockwell knows.
//
// Note: This function is intended for advanced users who need direct CIP protocol
// access. For standard tag operations, use Read(), Write(), and related functions.
func (client *Client) GenericCIPMessage(service CIPService, path, msg_data []byte) (*CIPItem, error) {

	reqitems := make([]CIPItem, 2)
	//reqitems[0] = cipItem{Header: cipItemHeader{ID: cipItem_Null}}
	reqitems[0] = newItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: uint16(sequencer()),
		Service:       service,
		PathLength:    byte(len(path) / 2),
	}

	reqitems[1] = newItem(cipItem_ConnectedData, readmsg)
	err := reqitems[1].Serialize(path)
	if err != nil {
		return nil, fmt.Errorf("could not serialize path: %w", err)
	}
	err = reqitems[1].Serialize(msg_data)
	if err != nil {
		return nil, fmt.Errorf("could not serialize msg data: %w", err)
	}

	itemdata, err := serializeItems(reqitems)
	if err != nil {
		return nil, err
	}
	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, itemdata)

	if err != nil {
		return nil, err
	}
	_ = hdr

	read_result_header := msgCIPResultHeader{}
	err = binary.Read(data, binary.LittleEndian, &read_result_header)
	if err != nil {
		return nil, fmt.Errorf("problpm reading read result header: %w", err)
	}
	items, err := readItems(data)
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
