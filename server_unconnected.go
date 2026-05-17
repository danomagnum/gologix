package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"
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

	// Direct UCMM Get_Attribute_Single — FactoryTalk Linx fires this against
	// Identity (and a handful of other classes) right after Get_Attributes_All
	// to validate the peer. Routing it here avoids the silent-drop trap.
	case CIPService_GetAttributeSingle:
		err = h.unconnectedServiceGetAttrSingle(item)
		if err != nil {
			return fmt.Errorf("problem handling get attribute single. %w", err)
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
			// Reply with an explicit ServiceNotSupported CIP error so the
			// originator knows the device is alive but the service isn't
			// implemented. Returning a Go error here used to bubble up to
			// the connection loop and silently drop the reply, which made
			// FactoryTalk Linx (and any other strict client) interpret the
			// device as dead and tear down the shortcut bind. See upstream
			// issue danomagnum/gologix#13.
			h.server.Logger.Warn("unsupported embedded service in UnconnectedSend", "service", emService)
			return h.sendUnconnectedErrorReply(emService, CIPStatus_ServiceNotSupported)
		}

	default:
		// Top-level service code that we don't recognize — same rule as the
		// embedded-service default above: reply with ServiceNotSupported
		// instead of silently dropping. Without this, FactoryTalk Linx
		// retries forever waiting for a response that never comes.
		h.server.Logger.Warn("unsupported unconnected service", "service", service)
		return h.sendUnconnectedErrorReply(service, CIPStatus_ServiceNotSupported)
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
		// The path is malformed (or uses segment types we don't decode).
		// Treat it as a path error rather than a silent drop.
		h.server.Logger.Warn("get_attributes_all: malformed path", "bytes", pathBytes)
		return h.sendUnconnectedErrorReply(CIPService_GetAttributeAll, CIPStatus_PathSegmentError)
	}

	switch class {
	case uint16(CipObject_Identity):
		if instance != 1 {
			return h.sendUnconnectedErrorReply(CIPService_GetAttributeAll, CIPStatus_PathDestinationUnknown)
		}
		payload, err := buildIdentityGetAttributesAllResponse(h.server.Attributes)
		if err != nil {
			return fmt.Errorf("get_attributes_all: build identity response: %w", err)
		}
		return h.sendUnconnectedRRDataReply(CIPService_GetAttributeAll, payload)

	case 0x64:
		// Class 0x64 = Logix Program Object. pycomm3.LogixDriver.get_plc_name()
		// probes Instance 1 of this class immediately after Identity to learn
		// the controller name. The Logix wire layout starts with a uint32
		// ProgramNumber followed by a SHORT_STRING ProgramName. We expose the
		// product name from the Identity attributes as the program name so
		// the same Server.Attributes map remains the single source of truth.
		if instance != 1 {
			return h.sendUnconnectedErrorReply(CIPService_GetAttributeAll, CIPStatus_PathDestinationUnknown)
		}
		payload, err := buildProgramObjectGetAttributesAllResponse(h.server.Attributes)
		if err != nil {
			return fmt.Errorf("get_attributes_all: build program response: %w", err)
		}
		return h.sendUnconnectedRRDataReply(CIPService_GetAttributeAll, payload)

	default:
		// Class isn't implemented at all — answer formally so the originator
		// keeps the session alive instead of timing out on silence.
		h.server.Logger.Warn("get_attributes_all: class not implemented", "class", class, "instance", instance)
		return h.sendUnconnectedErrorReply(CIPService_GetAttributeAll, CIPStatus_PathDestinationUnknown)
	}
}

// buildSymbolObjectInstanceListResponse renders the payload for the
// Symbol Object (Class 0x6B) Get_Instance_Attribute_List service (0x55).
// Logix exposes every tag as one instance of Class 0x6B. The response
// is a stream of {InstanceID uint32, SymbolName SHORT_STRING, SymbolType
// uint16} triplets in instance-id order. We synthesize stable
// instance IDs from the alphabetical position of the tag name so the
// same provider produces the same IDs across restarts — clients cache
// these IDs for subsequent reads.
//
// startInstance lets callers iterate when the response would exceed the
// transport limit (~500 bytes). For now Server_Class3 exposes only a
// handful of tags so we emit everything in one reply, but the loop is
// written so a future paginator can drop in without restructuring.
func buildSymbolObjectInstanceListResponse(tags []ServerTagInfo, startInstance uint32) []byte {
	// Stable ordering keeps instance IDs deterministic across reads.
	sort.Slice(tags, func(i, j int) bool { return tags[i].Name < tags[j].Name })

	var buf bytes.Buffer
	for i, tag := range tags {
		instanceID := uint32(i + 1) // CIP instances are 1-based
		if instanceID < startInstance {
			continue
		}
		_ = binary.Write(&buf, binary.LittleEndian, instanceID)
		name := tag.Name
		if len(name) > 255 {
			name = name[:255]
		}
		_ = binary.Write(&buf, binary.LittleEndian, uint16(len(name)))
		buf.WriteString(name)
		_ = binary.Write(&buf, binary.LittleEndian, uint16(tag.Type))
	}
	return buf.Bytes()
}

// buildProgramObjectGetAttributesAllResponse renders the payload for
// Get_Attributes_All on Class 0x64 (Logix Program Object) Instance 1. The
// wire layout that pycomm3.get_plc_name() expects is a SHORT_STRING holding
// the controller name; the gologix server doesn't have programs in the
// Logix sense, so we surface the ProductName from the Identity attributes.
func buildProgramObjectGetAttributesAllResponse(attrs map[CIPAttribute]any) ([]byte, error) {
	name, ok := attrs[7].(string)
	if !ok {
		return nil, fmt.Errorf("program object: identity attribute 7 (ProductName) missing or wrong type: %T", attrs[7])
	}
	if len(name) > 255 {
		name = name[:255]
	}
	var buf bytes.Buffer
	buf.WriteByte(byte(len(name)))
	buf.WriteString(name)
	return buf.Bytes(), nil
}

// parseClassInstancePath decodes the CIP EPATH segments commonly used to
// address a Class + Instance: an 8-bit Logical Class segment followed by
// either an 8-bit Logical Instance segment (the canonical Identity probe
// `0x20 <class> 0x24 <instance>` = 4 bytes) or a 16-bit Logical Instance
// segment (`0x20 <class> 0x25 0x00 <instance_lo> <instance_hi>` = 6
// bytes) that pycomm3.LogixDriver and the Symbol Object iterator both
// emit. Returning ok=false on anything else lets the caller respond
// with a formal PathSegmentError instead of silently dropping the
// request.
//
// Reference: ODVA Vol 1 §C-1.4.2 (logical segments).
func parseClassInstancePath(p []byte) (class, instance uint16, ok bool) {
	switch len(p) {
	case 4:
		if p[0] != 0x20 || p[2] != 0x24 {
			return 0, 0, false
		}
		return uint16(p[1]), uint16(p[3]), true
	case 6:
		// 8-bit class + 16-bit instance: `20 <class> 25 00 <inst_lo> <inst_hi>`
		if p[0] != 0x20 || p[2] != 0x25 || p[3] != 0x00 {
			return 0, 0, false
		}
		instance = uint16(p[4]) | uint16(p[5])<<8
		return uint16(p[1]), instance, true
	}
	return 0, 0, false
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
	results, err := parseWriteValues(&item)
	if err != nil {
		return fmt.Errorf("problem parsing write values for tag %s: %w", tag, err)
	}
	qty := uint16(len(results))
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

	// CIP wire format for Get_Attribute_Single requests: 1-byte path size in
	// words followed by `pathSize*2` bytes of EPATH segments. This is the
	// same shape in both the direct UCMM case (where the dispatcher just
	// consumed the service byte) and the embedded-in-UnconnectedSend case
	// (where the wrapper's outer fields were already drained before reaching
	// here). The previous `item.Uint16()` read silently consumed an extra
	// byte and only happened to match `!= 3` because nothing else was read
	// from the item after `attr.Read`; the moment FactoryTalk Linx started
	// probing real Identity attributes (0x15, 0x16, ...) it tripped the
	// "Path segment error" branch instead of looking up the attribute.
	pathSize, err := item.Byte()
	if err != nil {
		return fmt.Errorf("couldn't read path size: %w", err)
	}
	if pathSize != 3 {
		h.server.Logger.Debug("getattrsingle: unsupported path size", "size", pathSize)
		return h.sendUnconnectedErrorReply(CIPService_GetAttributeSingle, CIPStatus_PathSegmentError)
	}

	var cls CIPClass
	if err := cls.Read(&item); err != nil {
		return fmt.Errorf("could not read class: %w", err)
	}
	var inst CIPInstance
	if err := inst.Read(&item); err != nil {
		return fmt.Errorf("could not read instance: %w", err)
	}
	var attr CIPAttribute
	if err := attr.Read(&item); err != nil {
		return fmt.Errorf("could not read attribute ID: %w", err)
	}

	// Class 0x01 Identity is the only class with attribute data in this
	// minimal server. Other classes FactoryTalk Linx probes (0x47 DLR,
	// 0xF4 Port, ...) should produce a formal "PathDestinationUnknown"
	// error so the originator can mark them as not-implemented and proceed.
	if cls != CIPClass(CipObject_Identity) || inst != 1 {
		h.server.Logger.Debug("getattrsingle: class/instance not implemented", "class", cls, "instance", inst, "attr", attr)
		return h.sendUnconnectedErrorReply(CIPService_GetAttributeSingle, CIPStatus_PathDestinationUnknown)
	}

	val, ok := h.server.Attributes[attr]
	if !ok {
		// Attribute not in our Identity attribute map. Real Logix
		// controllers respond with status 0x14 AttributeNotSupported, which
		// FactoryTalk Linx interprets as "device alive but doesn't expose
		// this attribute" — exactly what we want here.
		h.server.Logger.Debug("getattrsingle: identity attribute not supported", "attr", attr)
		return h.sendUnconnectedErrorReply(CIPService_GetAttributeSingle, CIPStatus_AttributeNotSupported)
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

// sendUnconnectedErrorReply emits a CIP error response carrying a non-zero
// General Status. Used whenever the dispatcher receives a service / class /
// attribute it doesn't implement: instead of returning a Go error (which
// leaves the wire silent and makes the originator think the device died),
// we respond with the appropriate CIP status code so strict clients like
// FactoryTalk Linx, Studio 5000 path browser, and pycomm3 can mark the
// peer as alive-but-unsupported and continue gracefully.
//
// Reference: ODVA Vol 1 §2-4.5 General Status Codes; upstream issue
// danomagnum/gologix#13 documents the user-visible impact of silent drops.
func (h *serverTCPHandler) sendUnconnectedErrorReply(s CIPService, status CIPStatus) error {
	items := make([]CIPItem, 2)
	items[0] = newItem(cipItem_Null, nil)
	items[1] = newItem(cipItem_UnconnectedData, nil)
	resp := msgUnconnWriteResultHeader{
		Service: s.AsResponse(),
		Status:  status,
	}
	if err := items[1].Serialize(resp); err != nil {
		return fmt.Errorf("serialize error reply header: %w", err)
	}
	itemdata, err := serializeItems(items)
	if err != nil {
		return fmt.Errorf("serialize error reply items: %w", err)
	}
	return h.send(cipCommandSendRRData, itemdata)
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
