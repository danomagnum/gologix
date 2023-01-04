package gologix

import (
	"fmt"
	"log"
)

func (h *serverTCPHandler) cipConnectedWrite(items []cipItem) error {
	var l byte // length in words
	item := items[1]
	err := item.Unmarshal(l)
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
	var elements uint16
	err = item.Unmarshal(&elements)
	if err != nil {
		return fmt.Errorf("problem unmarshaling element count %w", err)
	}

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
