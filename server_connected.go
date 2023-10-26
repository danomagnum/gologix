package gologix

import (
	"fmt"
	"log"

	"github.com/danomagnum/gologix/cipclass"
	"github.com/danomagnum/gologix/cipservice"
	"github.com/danomagnum/gologix/ciptype"
	"github.com/danomagnum/gologix/eipcommand"
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

	var typ ciptype.CIPType
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
		results[i], err = typ.ReadValue(&item)
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

	return h.sendUnitDataReply(cipservice.Write)
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
	typ := ciptype.GoVarToCIPType(result)
	log.Printf("read %s to %v elements: %v %v. Value = %v\n", tag, path, qty, typ, result)

	return h.sendConnectedReadReply(cipservice.FragRead, seq, connection.OT, typ, byte(0), result)

}

func (h *serverTCPHandler) sendConnectedReadReply(s cipservice.CIPService, seq uint16, connID uint32, payload ...any) error {
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
	return h.send(eipcommand.SendUnitData, SerializeItems(items))
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
	log.Printf("got connection id %v = %+v", connID, connection)

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

	var cls cipclass.CIPClass
	err = cls.Read(&item)
	if err != nil {
		return fmt.Errorf("could not read class: %w", err)
	}

	var inst cipclass.CIPInstance
	err = inst.Read(&item)
	if err != nil {
		return fmt.Errorf("could not read instance: %w", err)
	}

	if cls != 1 || inst != 1 {
		return fmt.Errorf("only support class 1 instance 1 so far. got %d:%d", cls, inst)
	}

	var attr cipclass.CIPAttribute
	err = attr.Read(&item)
	if err != nil {
		return fmt.Errorf("could not read attribute ID: %w", err)
	}

	result, ok := h.server.Attributes[attr]
	if !ok {
		return fmt.Errorf("bad attribute %d", attr)
	}

	return h.sendConnectedReadReply(cipservice.FragRead, seq, connection.OT, result)

}
