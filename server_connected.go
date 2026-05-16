package gologix

import (
	"encoding/binary"
	"errors"
	"fmt"
)

func (h *serverTCPHandler) cipConnectedWrite(items []CIPItem) error {
	var l byte // length in words
	item := items[1]
	err := item.DeSerialize(&l)
	if err != nil {
		return fmt.Errorf("problem deserializing item length %w", err)
	}

	tag, err := getTagFromPath(&item)
	if err != nil {
		return fmt.Errorf("problem deserializing tag bytes %w", err)
	}

	results, err := parseWriteValues(&item)
	if err != nil {
		return fmt.Errorf("problem parsing write values for tag %s: %w", tag, err)
	}
	qty := uint16(len(results))
	h.server.Logger.Debug("write", "tag", tag, "value", results)

	if items[0].Header.ID != cipItem_ConnectionAddress {
		return fmt.Errorf("expected a connection address item in item 0. got %v", items[0].Header.ID)
	}
	items[0].Reset()
	var connID uint32
	err = items[0].DeSerialize(&connID)
	if err != nil {
		return fmt.Errorf("problem deserializing connection ID (%+v) %w", items[0], err)
	}

	conn, err := h.server.ConnMgr.GetByOT(connID)
	if err != nil {
		return fmt.Errorf("no server handler for %v. %w", connID, err)
	}

	p, err := h.server.Router.Resolve(conn.Path)
	if err != nil {
		return fmt.Errorf("problem finding tag provider for %v(%v). %w", connID, conn.Path, err)
	}
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

	// path is part of the forward open we've previously received.

	return h.sendUnitDataReply(CIPService_Write)
}

func (h *serverTCPHandler) connectedMulti(items []CIPItem) error {
	items[0].Reset()
	var connID uint32
	err := items[0].DeSerialize(&connID)
	if err != nil {
		return fmt.Errorf("couldn't get connection ID from item 0: %w", err)
	}
	connection, err := h.server.ConnMgr.GetByOT(connID)
	if err != nil {
		return fmt.Errorf("couldn't get connection with ID %v: %w", connID, err)
	}

	path := connection.Path

	provider, err := h.server.Router.Resolve(path)
	if err != nil {
		return fmt.Errorf("problem finding tag provider for %v. %w", path, err)
	}
	p := provider

	items[1].Reset()
	item := items[1]

	var seq uint16
	err = item.DeSerialize(&seq)
	if err != nil {
		return fmt.Errorf("error getting sequence ID: %w", err)
	}
	var service CIPService
	err = item.DeSerialize(&service)
	if err != nil {
		return fmt.Errorf("error getting service: %w", err)
	}
	var pathlen byte
	err = item.DeSerialize(&pathlen)
	if err != nil {
		return fmt.Errorf("error getting path len: %w", err)
	}

	// need to handle this differently!
	pb := make([]byte, pathlen*2)
	err = item.DeSerialize(&pb)
	if err != nil {
		return fmt.Errorf("problem dumping path bytes from multi-svc header")
	}
	qty, err := item.Uint16()
	if err != nil {
		return fmt.Errorf("problem getting multisvc item count: %w", err)
	}
	offsets := make([]int16, qty)
	results := make([][]byte, qty)
	for i := range offsets {
		offsets[i], err = item.Int16()
		if err != nil {
			return fmt.Errorf("problem getting multisvc offset for item %d: %w", i, err)
		}
	}

	for i := range results {
		item.Pos = int(offsets[i]) + 8 // seek to the start of the multi-service payload data.
		var svc CIPService
		err = item.DeSerialize(&svc)
		if err != nil {
			return fmt.Errorf("problem getting multiread service for item %d: %w", i, err)
		}

		svcReqSize, err := item.Byte()
		if err != nil {
			return fmt.Errorf("problem getting service req payload size for item %d: %w", i, err)
		}
		_ = svcReqSize // don't need this yet.  It is in 16 bit words.

		switch svc {
		case CIPService_FragRead,
			CIPService_Read:
			tagname, err := getTagFromPath(&item)
			if err != nil {
				return fmt.Errorf("problem getting multiread tagname for item %d: %w", i, err)
			}
			qty, err := item.Uint16()
			if err != nil {
				return fmt.Errorf("error getting write qty: %w", err)
			}

			result, err := p.TagRead(tagname, int16(qty))
			if err != nil {
				return fmt.Errorf("problem getting data for %s from provider. %w", tagname, err)
			}

			// build this portion of the response msg

			typ, _ := GoVarToCIPType(result)
			rhdr := msgMultiReadResult{Service: svc.AsResponse(), Status: 0, Type: typ}
			b, err := Serialize(rhdr)
			if err != nil {
				return fmt.Errorf("problem serializing header for %s: %w", tagname, err)
			}
			if typ == CIPTypeSTRING {
				res_str, ok := result.(string)
				if !ok {
					return errors.New("expected a string but didn't get one")
				}
				b_dat, err := Serialize(cipStringPacker(res_str))
				if err != nil {
					return fmt.Errorf("problem serializing data for %s: %w", tagname, err)
				}
				b.Truncate(4)
				_, err = b_dat.WriteTo(b)
				if err != nil {
					return fmt.Errorf("problem combining header and data for %s: %w", tagname, err)
				}
			} else {
				b_dat, err := Serialize(result)
				if err != nil {
					return fmt.Errorf("problem serializing data for %s: %w", tagname, err)
				}
				_, err = b_dat.WriteTo(b)
				if err != nil {
					return fmt.Errorf("problem combining header and data for %s: %w", tagname, err)
				}
			}

			results[i] = b.Bytes()
		case CIPService_Write,
			CIPService_FragWrite:

			tag, err := getTagFromPath(&item)
			if err != nil {
				return fmt.Errorf("problem deserializing tag bytes %w", err)
			}

			var typ CIPType
			err = item.DeSerialize(&typ)
			if err != nil {
				return fmt.Errorf("problem deserializing cip type %w", err)
			}
			var reserved byte
			err = item.DeSerialize(&reserved)
			if err != nil {
				return fmt.Errorf("problem deserializing reserved byte %w", err)
			}
			var qty uint16
			err = item.DeSerialize(&qty)
			if err != nil {
				return fmt.Errorf("problem deserializing element count %w", err)
			}

			writeTags := make([]any, qty)
			for i := 0; i < int(qty); i++ {
				writeTags[i], err = typ.readValue(&item)
				if err != nil {
					return fmt.Errorf("problem getting write element %d: %w", i, err)
				}
			}
			h.server.Logger.Debug("Write", "tag", tag, "value", writeTags)

			if qty > 1 {
				err = p.TagWrite(tag, writeTags)
				if err != nil {
					return fmt.Errorf("problem writing tag %v. %w", tag, err)
				}
			} else {
				err = p.TagWrite(tag, writeTags[0])
				if err != nil {
					return fmt.Errorf("problem writing tag %v. %w", tag, err)
				}
			}

			results[i] = []byte{byte(svc.AsResponse()), 0x00, 0x00, 0x00}

		}

	}

	response := make([]any, qty+2)
	response[0] = qty

	pos := (int(qty) + 1) * 2
	jump_table := make([]uint16, qty)
	for i := range jump_table {
		jump_table[i] = uint16(pos)
		pos += len(results[i])
	}
	response[1] = jump_table
	for i := range results {
		response[2+i] = results[i]
	}

	return h.sendConnectedReply(CIPService_MultipleService, seq, connection.OT, response...)
}

// connectedRead handles both CIPService_Read (0x4C) and CIPService_FragRead
// (0x52). Callers MUST pass the request service so reply and error frames
// echo the matching response code (0xCC for Read, 0xD2 for FragRead) — some
// CIP clients (pylogix on older firmware paths) tolerate a mismatched code,
// but Logix MSG instructions and stricter SCADA stacks reject it.
func (h *serverTCPHandler) connectedRead(reqSvc CIPService, items []CIPItem) error {
	items[0].Reset()
	var connID uint32
	err := items[0].DeSerialize(&connID)
	if err != nil {
		err2 := h.sendConnectedError(reqSvc, 0, 0, CIPStatus_InvalidParameter, 0)
		return fmt.Errorf("couldn't get connection ID from item 0: %w / %v", err, err2)
	}
	connection, err := h.server.ConnMgr.GetByOT(connID)
	if err != nil {
		err2 := h.sendConnectedError(reqSvc, 0, 0, CIPStatus_ConnectionLost, 0)
		return fmt.Errorf("couldn't get connection with ID %v: %w /%v", connID, err, err2)
	}

	items[1].Reset()
	item := items[1]

	var seq uint16
	err = item.DeSerialize(&seq)
	if err != nil {
		err2 := h.sendConnectedError(reqSvc, seq, connection.OT, CIPStatus_InvalidParameter, 0)
		return fmt.Errorf("error getting sequence ID: %w / %v", err, err2)
	}
	var pathlen uint16
	err = item.DeSerialize(&pathlen)
	if err != nil {
		err2 := h.sendConnectedError(reqSvc, seq, connection.OT, CIPStatus_InvalidParameter, 0)
		return fmt.Errorf("error getting path len: %w / %v", err, err2)
	}
	tag, err := getTagFromPath(&item)
	if err != nil {
		err2 := h.sendConnectedError(reqSvc, seq, connection.OT, CIPStatus_InvalidParameter, 0)
		return fmt.Errorf("couldn't parse path: %w / %v", err, err2)
	}
	var qty uint16
	err = item.DeSerialize(&qty)
	if err != nil {
		err2 := h.sendConnectedError(reqSvc, seq, connection.OT, CIPStatus_InvalidParameter, 0)
		return fmt.Errorf("error getting write qty: %w / %v", err, err2)
	}

	path := connection.Path

	provider, err := h.server.Router.Resolve(path)
	if err != nil {
		err2 := h.sendConnectedError(reqSvc, seq, connection.OT, CIPStatus_PathDestinationUnknown, 0)
		return fmt.Errorf("problem finding tag provider for %v. %w: %v", path, err, err2)
	}
	p := provider
	result, err := p.TagRead(tag, int16(qty))
	if err != nil {
		err2 := h.sendConnectedError(reqSvc, seq, connection.OT, CIPStatus_InvalidMemberID, 0)
		return fmt.Errorf("problem getting data from provider. %w / %v", err, err2)
	}
	typ, _ := GoVarToCIPType(result)

	if typ == CIPTypeSTRING {
		res_str, ok := result.(string)
		if !ok {
			err2 := h.sendConnectedError(reqSvc, seq, connection.OT, CIPStatus_InvalidAttributeValue, 0)
			return fmt.Errorf("was expecting a string but didn't get one: %w", err2)
		}
		return h.sendConnectedReply(reqSvc, seq, connection.OT, cipStringPacker(res_str))
	} else {
		return h.sendConnectedReply(reqSvc, seq, connection.OT, typ, byte(0), result)
	}
}

// cipStringStructCRC is the StructTypeCRC the Logix STRING UDT serializes
// with on the wire (see lgxtypes.STRING.TypeAbbr). External CIP clients
// (pylogix, Kepware, Ignition, MSG instructions) tag both reads and
// writes of native STRINGs with this value, prefixed by the 0xA0
// CIPTypeStruct byte.
const cipStringStructCRC uint16 = 0x0FCE

// cipStringDataLen is the fixed DATA buffer width of a Logix STRING UDT
// (SINT[82]). The wire payload always carries 82 bytes regardless of the
// real string length, padded with zeros past LEN. External clients treat
// anything shorter as malformed and silently fail to extract the value.
const cipStringDataLen = 82

// cipStringSlotLen is the full per-element wire footprint of a STRING:
// 4-byte type segment (0xA0 0x02 + StructTypeCRC LE) + 4-byte LEN +
// 82-byte DATA = 90 bytes.
const cipStringSlotLen = 4 + 4 + cipStringDataLen

type cipStringPacker string

func (c cipStringPacker) Len() int {
	return cipStringSlotLen
}
func (c cipStringPacker) Bytes() []byte {
	b := make([]byte, cipStringSlotLen)
	b[0] = 0xA0
	b[1] = 0x02
	binary.LittleEndian.PutUint16(b[2:], cipStringStructCRC)
	// LEN reflects the real string length, clamped to the 82-byte payload
	// so the header never claims more bytes than the slot can carry.
	l := len(c)
	if l > cipStringDataLen {
		l = cipStringDataLen
	}
	binary.LittleEndian.PutUint32(b[4:], uint32(l))
	copy(b[8:], c)
	return b
}

// parseWriteValues decodes the payload portion of a CIP write request after
// the path has been consumed. The wire is laid out as:
//
//	[typ: byte] [type_info_length_bytes: byte]
//	[type_info: type_info_length_bytes bytes]      -- only for struct types
//	[qty: uint16]
//	for each element:
//	    atomic: typ.readValue
//	    STRING struct (CRC 0x0FCE): [LEN: uint32] [DATA: 82 bytes]
//
// For atomic writes type_info_length_bytes is 0 — the field doubles as the
// high byte of the uint16 DataType the gologix client serializes — so the
// parser falls back to the historic readValue loop. For structured writes
// (typ == CIPTypeStruct) the segment carries 2 bytes of StructTypeCRC; only
// the Logix STRING UDT (0x0FCE) is supported here. Other UDTs return an
// error rather than the previous silent success or readValue panic.
//
// Both cipConnectedWrite and unconnectedServiceWrite go through this helper
// so the wire format stays in one place.
func parseWriteValues(item *CIPItem) ([]any, error) {
	var typ CIPType
	if err := item.DeSerialize(&typ); err != nil {
		return nil, fmt.Errorf("error reading write type: %w", err)
	}
	var typeInfoLen byte
	if err := item.DeSerialize(&typeInfoLen); err != nil {
		return nil, fmt.Errorf("error reading type-info length: %w", err)
	}

	// Read whatever type info bytes the segment declares before qty.
	var structCRC uint16
	if typeInfoLen > 0 {
		typeInfo := make([]byte, int(typeInfoLen))
		if err := item.DeSerialize(&typeInfo); err != nil {
			return nil, fmt.Errorf("error reading %d bytes of type info: %w", len(typeInfo), err)
		}
		if typ == CIPTypeStruct && len(typeInfo) >= 2 {
			structCRC = binary.LittleEndian.Uint16(typeInfo[0:2])
		}
	}

	var qty uint16
	if err := item.DeSerialize(&qty); err != nil {
		return nil, fmt.Errorf("error reading element count: %w", err)
	}

	results := make([]any, qty)
	if typ == CIPTypeStruct {
		if structCRC != cipStringStructCRC {
			return nil, fmt.Errorf("server only supports STRING struct writes (CRC 0x%04X); got CRC 0x%04X", cipStringStructCRC, structCRC)
		}
		for i := 0; i < int(qty); i++ {
			var slen uint32
			if err := item.DeSerialize(&slen); err != nil {
				return nil, fmt.Errorf("error reading STRING LEN for element %d: %w", i, err)
			}
			data := make([]byte, cipStringDataLen)
			if err := item.DeSerialize(&data); err != nil {
				return nil, fmt.Errorf("error reading STRING DATA for element %d: %w", i, err)
			}
			if slen > cipStringDataLen {
				slen = cipStringDataLen
			}
			results[i] = string(data[:slen])
		}
		return results, nil
	}

	for i := 0; i < int(qty); i++ {
		val, err := typ.readValue(item)
		if err != nil {
			return nil, fmt.Errorf("error reading element %d: %w", i, err)
		}
		results[i] = val
	}
	return results, nil
}

func (h *serverTCPHandler) sendConnectedReply(s CIPService, seq uint16, connID uint32, payload ...any) error {
	items := make([]CIPItem, 2)
	items[0] = newItem(cipItem_ConnectionAddress, connID)
	items[1] = newItem(cipItem_ConnectedData, seq)
	resp := msgUnconnWriteResultHeader{
		Service: s.AsResponse(),
	}
	err := items[1].Serialize(resp)
	if err != nil {
		return fmt.Errorf("couldn't serialize error response: %w", err)
	}
	for i := range payload {
		err = items[1].Serialize(payload[i])
		if err != nil {
			return fmt.Errorf("problem serializing unconnected data payload. %w", err)
		}
	}
	itemdata, err := serializeItems(items)
	if err != nil {
		return fmt.Errorf("could not serialize items: %w", err)
	}
	return h.send(cipCommandSendUnitData, itemdata)
}

func (h *serverTCPHandler) sendConnectedError(s CIPService, seq uint16, connID uint32, status CIPStatus, statusExtended byte) error {
	items := make([]CIPItem, 2)
	items[0] = newItem(cipItem_ConnectionAddress, connID)
	items[1] = newItem(cipItem_ConnectedData, seq)
	resp := msgUnconnWriteResultHeader{
		Service:        s.AsResponse(),
		Status:         status,
		StatusExtended: statusExtended,
	}
	err := items[1].Serialize(resp)
	if err != nil {
		return fmt.Errorf("couldn't serialize error response: %w", err)
	}
	itemdata, err := serializeItems(items)
	if err != nil {
		return fmt.Errorf("could not serialize items: %w", err)
	}
	return h.send(cipCommandSendUnitData, itemdata)
}

func (h *serverTCPHandler) connectedGetAttr(items []CIPItem) error {
	items[0].Reset()
	var connID uint32
	err := items[0].DeSerialize(&connID)
	if err != nil {
		return fmt.Errorf("couldn't get connection ID from item 0: %w", err)
	}
	connection, err := h.server.ConnMgr.GetByOT(connID)
	if err != nil {
		return fmt.Errorf("couldn't get connection with ID %v: %w", connID, err)
	}

	items[1].Reset()
	return h.getAttrSingle(connection, items[1])

}

func (h *serverTCPHandler) getAttrSingle(connection *serverConnection, item CIPItem) error {

	var seq uint16
	err := item.DeSerialize(&seq)
	if err != nil {
		return fmt.Errorf("error getting sequence ID: %w", err)
	}

	_, err = item.Byte()
	if err != nil {
		return fmt.Errorf("couldn't read cmd type: %w", err)
	}

	path_size, err := item.Byte()
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

	result, ok := h.server.Attributes[attr]
	if !ok {
		return fmt.Errorf("bad attribute %d", attr)
	}

	return h.sendConnectedReply(CIPService_FragRead, seq, connection.OT, result)

}
