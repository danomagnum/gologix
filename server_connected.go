package gologix

import (
	"fmt"
	"log"
)

func (h *serverTCPHandler) cipConnectedWrite(items []cipItem) error {
	var l byte // length in words
	item := items[1]
	err := item.Unmarshal(&l)
	if err != nil {
		return fmt.Errorf("problem unmarshaling item length %w", err)
	}
	var path_type SegmentType
	err = item.Unmarshal(&path_type)
	if err != nil {
		return fmt.Errorf("problem unmarshaling path type %w", err)
	}
	if path_type != SegmentTypeExtendedSymbolic {
		return fmt.Errorf("only support symbolic writes. got segment type %v", path_type)
	}
	var tag_length byte
	err = item.Unmarshal(&tag_length)
	if err != nil {
		return fmt.Errorf("problem unmarshaling tag length %w", err)
	}
	tag_bytes := make([]byte, tag_length)
	err = item.Unmarshal(&tag_bytes)
	if err != nil {
		return fmt.Errorf("problem unmarshaling tag bytes %w", err)
	}
	tag := string(tag_bytes)

	// string will be padded with a null if odd length
	if (tag_length % 2) == 1 {
		var b byte
		err = item.Unmarshal(&b)
		if err != nil {
			return fmt.Errorf("problem unmarshaling odd length pad byte %w", err)
		}
	}

	var typ CIPType
	err = item.Unmarshal(&typ)
	if err != nil {
		return fmt.Errorf("problem unmarshaling cip type %w", err)
	}
	var reserved byte
	err = item.Unmarshal(&reserved)
	if err != nil {
		return fmt.Errorf("problem unmarshaling reserved byte %w", err)
	}
	var qty uint16
	err = item.Unmarshal(&qty)
	if err != nil {
		return fmt.Errorf("problem unmarshaling element count %w", err)
	}

	results := make([]any, qty)
	log.Printf("tag: %s", tag)
	for i := 0; i < int(qty); i++ {
		results[i] = typ.readValue(&item)
	}
	log.Printf("value: %v", results)

	if items[0].Header.ID != cipItem_ConnectionAddress {
		return fmt.Errorf("expected a connection address item in item 0. got %v", items[0].Header.ID)
	}
	items[0].Reset()
	var connID uint32
	err = items[0].Unmarshal(&connID)
	if err != nil {
		return fmt.Errorf("problem unmarshaling connection ID (%+v) %w", items[0], err)
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

func (h *serverTCPHandler) connectedData(item cipItem) error {
	var service CIPService
	var err error
	err = item.Unmarshal(&service)
	if err != nil {
		return fmt.Errorf("problem unmarshaling service %w", err)
	}
	item.Reset()
	switch service {
	case cipService_FragRead:
		err = h.cipFragRead(&item)
		if err != nil {
			return fmt.Errorf("problem handling frag read. %w", err)
		}
	default:
		log.Printf("Got unknown service %d", service)
	}
	log.Printf("sendrrdata service requested: %v", service)
	return nil
}
