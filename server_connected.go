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

	results := make([]any, qty)
	for i := 0; i < int(qty); i++ {
		results[i], err = typ.readValue(&item)
		if err != nil {
			return fmt.Errorf("problem reading element %d: %w", i, err)
		}
	}
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

func (h *serverTCPHandler) connectedRead(items []CIPItem) error {
	items[0].Reset()
	var connID uint32
	err := items[0].DeSerialize(&connID)
	if err != nil {
		err2 := h.sendConnectedError(CIPService_FragRead, 0, 0, CIPStatus_InvalidParameter, 0)
		return fmt.Errorf("couldn't get connection ID from item 0: %w / %v", err, err2)
	}
	connection, err := h.server.ConnMgr.GetByOT(connID)
	if err != nil {
		err2 := h.sendConnectedError(CIPService_FragRead, 0, 0, CIPStatus_ConnectionLost, 0)
		return fmt.Errorf("couldn't get connection with ID %v: %w /%v", connID, err, err2)
	}

	items[1].Reset()
	item := items[1]

	var seq uint16
	err = item.DeSerialize(&seq)
	if err != nil {
		err2 := h.sendConnectedError(CIPService_FragRead, seq, connection.OT, CIPStatus_InvalidParameter, 0)
		return fmt.Errorf("error getting sequence ID: %w / %v", err, err2)
	}
	var pathlen uint16
	err = item.DeSerialize(&pathlen)
	if err != nil {
		err2 := h.sendConnectedError(CIPService_FragRead, seq, connection.OT, CIPStatus_InvalidParameter, 0)
		return fmt.Errorf("error getting path len: %w / %v", err, err2)
	}
	tag, err := getTagFromPath(&item)
	if err != nil {
		err2 := h.sendConnectedError(CIPService_FragRead, seq, connection.OT, CIPStatus_InvalidParameter, 0)
		return fmt.Errorf("couldn't parse path: %w / %v", err, err2)
	}
	var qty uint16
	err = item.DeSerialize(&qty)
	if err != nil {
		err2 := h.sendConnectedError(CIPService_FragRead, seq, connection.OT, CIPStatus_InvalidParameter, 0)
		return fmt.Errorf("error getting write qty: %w / %v", err, err2)
	}

	path := connection.Path

	provider, err := h.server.Router.Resolve(path)
	if err != nil {
		err2 := h.sendConnectedError(CIPService_FragRead, seq, connection.OT, CIPStatus_PathDestinationUnknown, 0)
		return fmt.Errorf("problem finding tag provider for %v. %w: %v", path, err, err2)
	}
	p := provider
	result, err := p.TagRead(tag, int16(qty))
	if err != nil {
		err2 := h.sendConnectedError(CIPService_FragRead, seq, connection.OT, CIPStatus_InvalidMemberID, 0)
		return fmt.Errorf("problem getting data from provider. %w / %v", err, err2)
	}
	typ, _ := GoVarToCIPType(result)

	if typ == CIPTypeSTRING {
		res_str, ok := result.(string)
		if !ok {
			err2 := h.sendConnectedError(CIPService_FragRead, seq, connection.OT, CIPStatus_InvalidAttributeValue, 0)
			return fmt.Errorf("was expecting a string but didn't get one: %w", err2)
		}
		return h.sendConnectedReply(CIPService_FragRead, seq, connection.OT, cipStringPacker(res_str))
	} else {
		return h.sendConnectedReply(CIPService_FragRead, seq, connection.OT, typ, byte(0), result)
	}
}

type cipStringPacker string

func (c cipStringPacker) Len() int {
	return 8 + len(c)
}
func (c cipStringPacker) Bytes() []byte {
	l := len(c)
	b := make([]byte, 8+l)
	b[0] = 0xA0
	b[1] = 0x02
	b[2] = 0xCE
	b[3] = 0x0F
	binary.LittleEndian.PutUint32(b[4:], uint32(l))
	copy(b[8:], c)
	return b
}

func (h *serverTCPHandler) sendConnectedReply(s CIPService, seq uint16, connID uint32, payload ...any) error {
	items := make([]CIPItem, 2)
	items[0] = newItem(cipItem_ConnectionAddress, connID)
	items[1] = newItem(cipItem_ConnectedData, seq)
	resp := msgUnconnWriteResultHeader{
		Service: s.AsResponse(),
	}
	items[1].Serialize(resp)
	for i := range payload {
		items[1].Serialize(payload[i])
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
	items[1].Serialize(resp)
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
