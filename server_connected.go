package gologix

import (
	"fmt"
	"log"
)

func (h *serverTCPHandler) cipConnectedWrite(items []CIPItem) error {
	var l byte // length in words
	item := items[1]
	err := item.DeSerialize(&l)
	if err != nil {
		return fmt.Errorf("problem deserializing item length %w", err)
	}
	var path_type SegmentType
	err = item.DeSerialize(&path_type)
	if err != nil {
		return fmt.Errorf("problem deserializing path type %w", err)
	}
	if path_type != SegmentTypeExtendedSymbolic {
		return fmt.Errorf("only support symbolic writes. got segment type %v", path_type)
	}
	var tag_length byte
	err = item.DeSerialize(&tag_length)
	if err != nil {
		return fmt.Errorf("problem deserializing tag length %w", err)
	}
	tag_bytes := make([]byte, tag_length)
	err = item.DeSerialize(&tag_bytes)
	if err != nil {
		return fmt.Errorf("problem deserializing tag bytes %w", err)
	}
	tag := string(tag_bytes)

	// string will be padded with a null if odd length
	if (tag_length % 2) == 1 {
		var b byte
		err = item.DeSerialize(&b)
		if err != nil {
			return fmt.Errorf("problem deserializing odd length pad byte %w", err)
		}
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
	log.Printf("tag: %s", tag)
	for i := 0; i < int(qty); i++ {
		results[i], err = typ.readValue(&item)
		if err != nil {
			return fmt.Errorf("problem reading element %d: %w", i, err)
		}
	}
	log.Printf("value: %v", results)

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

	return h.sendUnitDataReply(cipService_Write)
}

func (h *serverTCPHandler) connectedData(items []CIPItem) error {
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
	log.Printf("got connection id %v = %+v", connID, connection)

	items[1].Reset()
	return h.connectedFragRead(connection, items[1])

}

func (h *serverTCPHandler) connectedFragRead(connection *serverConnection, item CIPItem) error {

	var seq uint16
	err := item.DeSerialize(&seq)
	if err != nil {
		return fmt.Errorf("error getting sequence ID: %w", err)
	}
	var pathlen uint16
	err = item.DeSerialize(&pathlen)
	if err != nil {
		return fmt.Errorf("error getting path len: %w", err)
	}
	tag, err := getTagFromPath(&item)
	if err != nil {
		return fmt.Errorf("couldn't parse path: %w", err)
	}
	var qty uint16
	err = item.DeSerialize(&qty)
	if err != nil {
		return fmt.Errorf("error getting write qty: %w", err)
	}

	path := connection.Path

	provider, err := h.server.Router.Resolve(path)
	if err != nil {
		return fmt.Errorf("problem finding tag provider for %v. %w", path, err)
	}
	p := provider
	result, err := p.TagRead(tag, int16(qty))
	if err != nil {
		return fmt.Errorf("problem getting data from provider. %w", err)
	}
	typ := GoVarToCIPType(result)
	log.Printf("read %s to %v elements: %v %v. Value = %v\n", tag, path, qty, typ, result)

	return h.sendConnectedReadReply(cipService_FragRead, seq, connection.OT, typ, byte(0), result)

}

func (h *serverTCPHandler) sendConnectedReadReply(s CIPService, seq uint16, connID uint32, payload ...any) error {
	items := make([]CIPItem, 2)
	items[0] = NewItem(cipItem_ConnectionAddress, connID)
	items[1] = NewItem(cipItem_ConnectedData, seq)
	resp := msgUnconnWriteResultHeader{
		Service: s.AsResponse(),
	}
	items[1].Serialize(resp)
	for i := range payload {
		items[1].Serialize(payload[i])
	}
	return h.send(cipCommandSendUnitData, SerializeItems(items))
}
