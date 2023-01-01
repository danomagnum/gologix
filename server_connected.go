package gologix

import (
	"fmt"
	"log"
)

func (h *serverTCPHandler) cipConnectedWrite(items []cipItem) error {
	var l byte // length in words
	item := items[1]
	item.Unmarshal(l)
	var path_type SegmentType
	item.Unmarshal(&path_type)
	if path_type != SegmentTypeExtendedSymbolic {
		return fmt.Errorf("only support symbolic writes. got segment type %v", path_type)
	}
	var tag_length byte
	item.Unmarshal(&tag_length)
	tag_bytes := make([]byte, tag_length)
	item.Unmarshal(&tag_bytes)
	tag := string(tag_bytes)

	// string will be padded with a null if odd length
	if (tag_length % 2) == 1 {
		var b byte
		item.Unmarshal(&b)
	}

	var typ CIPType
	item.Unmarshal(&typ)
	var reserved byte
	item.Unmarshal(&reserved)
	var elements uint16
	item.Unmarshal(&elements)

	fmt.Printf("tag: %s", tag)
	for i := 0; i < int(elements); i++ {
		v := typ.readValue(&item)
		fmt.Printf("value: %v", v)
	}

	// path is part of the forward open we've previously received.

	return h.sendUnitDataReply(cipService_Write)
}

func (h *serverTCPHandler) connectedData(item cipItem) error {
	var service CIPService
	var err error
	item.Unmarshal(&service)
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
