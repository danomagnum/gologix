package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func (h *serverTCPHandler) unconnectedData(item CIPItem) error {
	var service CIPService
	var err error
	err = item.DeSerialize(&service)
	if err != nil {
		return fmt.Errorf("problem unmarshling service %w", err)
	}
	switch service {
	case CIPService_ForwardOpen:
		item.Reset()
		err = h.forwardOpen(item)
		if err != nil {
			return fmt.Errorf("problem handling forward open. %w", err)
		}

	case CIPService_LargeForwardOpen:
		item.Reset()
		err = h.largeForwardOpen(item)
		if err != nil {
			return fmt.Errorf("problem handling large forward open. %w", err)
		}
	case CIPService_ForwardClose:
		item.Reset()
		err = h.forwardClose(item)
		if err != nil {
			return fmt.Errorf("problem handling forward close. %w", err)
		}

	// CIPService_GetAttributeAll arrives as a direct UCMM request (not wrapped
	// in an UnconnectedSend) from clients such as FactoryTalk Linx and
	// pycomm3.LogixDriver, which probe Identity Class 0x01 Instance 1 before
	// they accept further communication with the device.
	case CIPService_GetAttributeAll:
		err = h.unconnectedGetAttributeAll(item)
		if err != nil {
			return fmt.Errorf("problem handling get attributes all. %w", err)
		}

	case 0x52:
		// unconnected send?
		var pathsize byte
		err = item.DeSerialize(&pathsize)
		if err != nil {
			return fmt.Errorf("error getting path size 305. %w", err)
		}
		path := make([]byte, pathsize*2)
		err = item.DeSerialize(&path)
		if err != nil {
			return fmt.Errorf("error getting path. %w", err)
		}
		var timeout uint16
		err = item.DeSerialize(&timeout)
		if err != nil {
			return fmt.Errorf("error getting timeout. %w", err)
		}
		var embedded_size uint16
		err = item.DeSerialize(&embedded_size)
		if err != nil {
			return fmt.Errorf("error getting embedded size. %w", err)
		}
		var emService CIPService
		err = item.DeSerialize(&emService)
		if err != nil {
			return fmt.Errorf("error getting embedded service. %w", err)
		}
		switch emService {
		case CIPService_Write:
			return h.unconnectedServiceWrite(item)
		case CIPService_Read:
			return h.unconnectedServiceRead(item)
		case CIPService_GetAttributeSingle:
			return h.unconnectedServiceGetAttrSingle(item)
		case CIPService_GetAttributeAll:
			return h.unconnectedGetAttributeAll(item)
		default:
			return fmt.Errorf("don't know how to handle service '%v'", emService)

		}
	}
	return nil
}

// unconnectedGetAttributeAll handles CIP Get_Attributes_All (service 0x01).
// External CIP clients (FactoryTalk Linx shortcut bind, pycomm3
// LogixDriver.open(), RSLogix MSG path browse) probe Identity Object
// (Class 0x01 Instance 1) immediately after RegisterSession and refuse to
// proceed if no response arrives. The caller has already consumed the
// service byte; the remaining wire layout is `pathSize:byte` followed by
// `pathSize*2` bytes of EPATH segments.
func (h *serverTCPHandler) unconnectedGetAttributeAll(item CIPItem) error {
	var pathSize byte
	if err := item.DeSerialize(&pathSize); err != nil {
		return fmt.Errorf("get_attributes_all: read path size: %w", err)
	}
	pathBytes := make([]byte, int(pathSize)*2)
	if err := item.DeSerialize(&pathBytes); err != nil {
		return fmt.Errorf("get_attributes_all: read path: %w", err)
	}

	class, instance, ok := parseClassInstancePath(pathBytes)
	if !ok {
		return fmt.Errorf("get_attributes_all: unsupported path segments % x", pathBytes)
	}
	if class != uint16(CipObject_Identity) || instance != 1 {
		return fmt.Errorf("get_attributes_all: only Class 0x01 Instance 1 supported, got class %d instance %d", class, instance)
	}

	payload, err := buildIdentityGetAttributesAllResponse(h.server.Attributes)
	if err != nil {
		return fmt.Errorf("get_attributes_all: build response: %w", err)
	}
	return h.sendUnconnectedRRDataReply(CIPService_GetAttributeAll, payload)
}

// parseClassInstancePath decodes a 4-byte EPATH carrying an 8-bit Class
// followed by an 8-bit Instance segment (`0x20 <class> 0x24 <instance>` —
// the canonical Identity Object probe shape). Anything else falls outside
// the scope of this minimal implementation and is rejected by the caller.
func parseClassInstancePath(p []byte) (class, instance uint16, ok bool) {
	if len(p) != 4 {
		return 0, 0, false
	}
	if p[0] != 0x20 || p[2] != 0x24 {
		return 0, 0, false
	}
	return uint16(p[1]), uint16(p[3]), true
}

// buildIdentityGetAttributesAllResponse renders the payload portion of a
// Get_Attributes_All reply on the Identity Object (attributes 1..7 in the
// order the spec mandates). The Attributes map is the same one the existing
// Get_Attribute_Single handler reads from — the field types match the
// historic int16/uint32/string layout NewServer seeds.
//
// Wire layout:
//
//	VendorID    UINT  (2)
//	DeviceType  UINT  (2)
//	ProductCode UINT  (2)
//	Revision    USINT,USINT (2)
//	Status      WORD  (2)
//	Serial      UDINT (4)
//	ProductName SHORT_STRING (1 + N)
func buildIdentityGetAttributesAllResponse(attrs map[CIPAttribute]any) ([]byte, error) {
	vendor, ok := attrs[1].(int16)
	if !ok {
		return nil, fmt.Errorf("attribute 1 (VendorID) missing or wrong type: %T", attrs[1])
	}
	deviceType, ok := attrs[2].(int16)
	if !ok {
		return nil, fmt.Errorf("attribute 2 (DeviceType) missing or wrong type: %T", attrs[2])
	}
	productCode, ok := attrs[3].(int16)
	if !ok {
		return nil, fmt.Errorf("attribute 3 (ProductCode) missing or wrong type: %T", attrs[3])
	}
	revision, ok := attrs[4].(int16)
	if !ok {
		return nil, fmt.Errorf("attribute 4 (Revision) missing or wrong type: %T", attrs[4])
	}
	status, ok := attrs[5].(int16)
	if !ok {
		return nil, fmt.Errorf("attribute 5 (Status) missing or wrong type: %T", attrs[5])
	}
	serial, ok := attrs[6].(uint32)
	if !ok {
		return nil, fmt.Errorf("attribute 6 (SerialNumber) missing or wrong type: %T", attrs[6])
	}
	productName, ok := attrs[7].(string)
	if !ok {
		return nil, fmt.Errorf("attribute 7 (ProductName) missing or wrong type: %T", attrs[7])
	}
	if len(productName) > 255 {
		productName = productName[:255]
	}

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, vendor); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, deviceType); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, productCode); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, revision); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, status); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, serial); err != nil {
		return nil, err
	}
	buf.WriteByte(byte(len(productName)))
	buf.WriteString(productName)
	return buf.Bytes(), nil
}

func (h *serverTCPHandler) unconnectedServiceWrite(item CIPItem) error {
	var reserved byte
	err := item.DeSerialize(&reserved)
	if err != nil {
		return fmt.Errorf("error getting reserved byte. %w", err)
	}
	tag, err := getTagFromPath(&item)
	if err != nil {
		return fmt.Errorf("couldn't parse path. %w", err)
	}
	var typ CIPType
	err = item.DeSerialize(&typ)
	if err != nil {
		return fmt.Errorf("error getting write type. %w", err)
	}
	var pad byte
	err = item.DeSerialize(&pad)
	if err != nil {
		return fmt.Errorf("error getting pad. %w", err)
	}
	var qty uint16
	err = item.DeSerialize(&qty)
	if err != nil {
		return fmt.Errorf("error getting write qty. %w", err)
	}
	// TODO: read structs gracefully.
	if typ == CIPTypeStruct {
		return h.sendUnconnectedRRDataReply(CIPService_Write)
	}
	results := make([]any, qty)
	for i := 0; i < int(qty); i++ {
		results[i], err = typ.readValue(&item)
		if err != nil {
			return fmt.Errorf("problem reading element %d: %w", i, err)
		}
	}
	var path_size uint16
	err = item.DeSerialize(&path_size)
	if err != nil {
		return fmt.Errorf("couldn't get path size 374. %w", err)
	}
	path := make([]byte, 2*path_size)
	err = item.DeSerialize(&path)
	if err != nil {
		return fmt.Errorf("couldn't get path. %w", err)
	}

	provider, err := h.server.Router.Resolve(path)
	if err != nil {
		return fmt.Errorf("problem finding tag provider for %v. %w", path, err)
	}
	p := provider
	if qty > 1 {
		err = p.TagWrite(tag, results)
		if err != nil {
			return fmt.Errorf("problem writing tag %v. %w", tag, err)
		}
	} else {
		err = p.TagWrite(tag, results[0])
		if err != nil {
			return fmt.Errorf("problem writing tag %v. %w", tag, err)
		}
	}

	return h.sendUnconnectedRRDataReply(CIPService_Write)

}

func (h *serverTCPHandler) unconnectedServiceGetAttrSingle(item CIPItem) error {

	path_size, err := item.Uint16()
	if err != nil {
		return fmt.Errorf("couldn't read data len: %w", err)
	}

	if path_size != 3 {
		return fmt.Errorf("currently only support getattrsingle path size of 3. got %d", path_size)
	}

	var cls CIPClass
	err = cls.Read(&item)
	if err != nil {
		return fmt.Errorf("could not read class: %w", err)
	}

	var inst CIPInstance
	err = inst.Read(&item)
	if err != nil {
		return fmt.Errorf("could not read instance: %w", err)
	}

	if cls != 1 || inst != 1 {
		return fmt.Errorf("only support class 1 instance 1 so far. got %d:%d", cls, inst)
	}

	var attr CIPAttribute
	err = attr.Read(&item)
	if err != nil {
		return fmt.Errorf("could not read attribute ID: %w", err)
	}

	val, ok := h.server.Attributes[attr]
	if !ok {
		return fmt.Errorf("bad attribute %d", attr)
	}

	typ, _ := GoVarToCIPType(val)

	return h.sendUnconnectedRRDataReply(CIPService_GetAttributeSingle, typ, byte(0), val)
}

func (h *serverTCPHandler) unconnectedServiceRead(item CIPItem) error {
	var reserved byte
	err := item.DeSerialize(&reserved)
	if err != nil {
		return fmt.Errorf("error getting reserved byte. %w", err)
	}
	tag, err := getTagFromPath(&item)
	if err != nil {
		return fmt.Errorf("couldn't parse path. %w", err)
	}
	var qty uint16
	err = item.DeSerialize(&qty)
	if err != nil {
		return fmt.Errorf("error getting write qty. %w", err)
	}
	var path_size uint16
	err = item.DeSerialize(&path_size)
	if err != nil {
		return fmt.Errorf("couldn't get path size 374. %w", err)
	}
	path := make([]byte, 2*path_size)
	err = item.DeSerialize(&path)
	if err != nil {
		return fmt.Errorf("couldn't get path. %w", err)
	}
	h.server.Logger.Debug("read", "tag", tag, "path", path, "qty", qty)

	provider, err := h.server.Router.Resolve(path)
	if err != nil {
		return fmt.Errorf("problem finding tag provider for %v. %w", path, err)
	}
	p := provider
	result, err := p.TagRead(tag, int16(qty))
	if err != nil {
		return fmt.Errorf("problem getting data from provider. %w", err)
	}
	h.server.Logger.Debug("Read", "tag", tag, "path", path, "qty", qty, "results", result)
	typ, _ := GoVarToCIPType(result)

	return h.sendUnconnectedRRDataReply(CIPService_Read, typ, byte(0), result)

}

func (h *serverTCPHandler) sendUnconnectedRRDataReply(s CIPService, payload ...any) error {
	items := make([]CIPItem, 2)
	items[0] = newItem(cipItem_Null, nil)
	items[1] = newItem(cipItem_UnconnectedData, nil)
	resp := msgUnconnWriteResultHeader{
		Service: s.AsResponse(),
	}
	err := items[1].Serialize(resp)
	if err != nil {
		return fmt.Errorf("problem serializing unconnected data header. %w", err)
	}
	for i := range payload {
		err = items[1].Serialize(payload[i])
		if err != nil {
			return fmt.Errorf("problem serializing unconnected data payload. %w", err)
		}
	}
	itemdata, err := serializeItems(items)
	if err != nil {
		return err
	}
	return h.send(cipCommandSendRRData, itemdata)
}

func (h *serverTCPHandler) sendUnconnectedUnitDataReply(s CIPService) error {
	items := make([]CIPItem, 2)
	items[0] = newItem(cipItem_Null, nil)
	items[1] = newItem(cipItem_UnconnectedData, nil)
	resp := msgWriteResultHeader{
		SequenceCount: h.UnitDataSequencer,
		Service:       s.AsResponse(),
	}
	err := items[1].Serialize(resp)
	if err != nil {
		return fmt.Errorf("problem serializing unconnected data header. %w", err)
	}
	itemdata, err := serializeItems(items)
	if err != nil {
		return err
	}
	return h.send(cipCommandSendUnitData, itemdata)
}
