package gologix

import (
	"fmt"
	"log"
)

func (h *serverTCPHandler) unconnectedData(item cipItem) error {
	var service CIPService
	var err error
	err = item.Unmarshal(&service)
	if err != nil {
		return fmt.Errorf("problem unmarshling service %w", err)
	}
	switch service {
	case cipService_ForwardOpen:
		item.Reset()
		err = h.forwardOpen(item)
		if err != nil {
			return fmt.Errorf("problem handling forward open. %w", err)
		}
	case cipService_ForwardClose:
		item.Reset()
		err = h.forwardClose(item)
		if err != nil {
			return fmt.Errorf("problem handling forward close. %w", err)
		}

	case 0x52:
		// unconnected send?
		var pathsize byte
		err = item.Unmarshal(&pathsize)
		if err != nil {
			return fmt.Errorf("error getting path size 305. %w", err)
		}
		path := make([]byte, pathsize*2)
		err = item.Unmarshal(&path)
		if err != nil {
			return fmt.Errorf("error getting path. %w", err)
		}
		var timeout uint16
		err = item.Unmarshal(&timeout)
		if err != nil {
			return fmt.Errorf("error getting timeout. %w", err)
		}
		var embedded_size uint16
		err = item.Unmarshal(&embedded_size)
		if err != nil {
			return fmt.Errorf("error getting embedded size. %w", err)
		}
		var emService CIPService
		err = item.Unmarshal(&emService)
		if err != nil {
			return fmt.Errorf("error getting embedded service. %w", err)
		}
		switch emService {
		case cipService_Write:
			return h.unconnectedServiceWrite(item)
		case cipService_Read:
			return h.unconnectedServiceRead(item)

		default:
			return fmt.Errorf("don't know how to handle service '%v'", emService)

		}
	}
	return nil
}

func (h *serverTCPHandler) unconnectedServiceWrite(item cipItem) error {
	var reserved byte
	err := item.Unmarshal(&reserved)
	if err != nil {
		return fmt.Errorf("error getting reserved byte. %w", err)
	}
	tag, err := getTagFromPath(&item)
	if err != nil {
		return fmt.Errorf("couldn't parse path. %w", err)
	}
	var typ CIPType
	err = item.Unmarshal(&typ)
	if err != nil {
		return fmt.Errorf("error getting write type. %w", err)
	}
	var pad byte
	err = item.Unmarshal(&pad)
	if err != nil {
		return fmt.Errorf("error getting pad. %w", err)
	}
	var qty uint16
	err = item.Unmarshal(&qty)
	if err != nil {
		return fmt.Errorf("error getting write qty. %w", err)
	}
	// TODO: read structs gracefully.
	if typ == CIPTypeStruct {
		log.Printf("read %s as %s * %v = %v", tag, typ, qty, item.Data[item.Pos:])
		return h.sendUnconnectedRRDataReply(cipService_Write)
	}
	results := make([]any, qty)
	for i := 0; i < int(qty); i++ {
		results[i] = typ.readValue(&item)
	}
	var path_size uint16
	err = item.Unmarshal(&path_size)
	if err != nil {
		return fmt.Errorf("couldn't get path size 374. %w", err)
	}
	path := make([]byte, 2*path_size)
	err = item.Unmarshal(&path)
	if err != nil {
		return fmt.Errorf("couldn't get path. %w", err)
	}
	log.Printf("Got Unconn Write from %v. Tag: %s Path:%v Type:%s Qty:%v Value:%v", h.conn.RemoteAddr(), tag, path, typ, qty, results)

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

	return h.sendUnconnectedRRDataReply(cipService_Write)

}

func (h *serverTCPHandler) unconnectedServiceRead(item cipItem) error {
	var reserved byte
	err := item.Unmarshal(&reserved)
	if err != nil {
		return fmt.Errorf("error getting reserved byte. %w", err)
	}
	tag, err := getTagFromPath(&item)
	if err != nil {
		return fmt.Errorf("couldn't parse path. %w", err)
	}
	var qty uint16
	err = item.Unmarshal(&qty)
	if err != nil {
		return fmt.Errorf("error getting write qty. %w", err)
	}
	var path_size uint16
	err = item.Unmarshal(&path_size)
	if err != nil {
		return fmt.Errorf("couldn't get path size 374. %w", err)
	}
	path := make([]byte, 2*path_size)
	err = item.Unmarshal(&path)
	if err != nil {
		return fmt.Errorf("couldn't get path. %w", err)
	}
	log.Printf("Path read: %v\n", path)
	log.Printf("read %s to %v elements: %v", tag, path, qty)

	provider, err := h.server.Router.Resolve(path)
	if err != nil {
		return fmt.Errorf("problem finding tag provider for %v. %w", path, err)
	}
	p := provider
	result, err := p.TagRead(tag, int16(qty))
	if err != nil {
		return fmt.Errorf("problem getting data from provider. %w", err)
	}
	log.Printf("read %s to %v elements: %v. Value = %v\n", tag, path, qty, result)
	typ := GoVarToCIPType(result)

	return h.sendUnconnectedRRDataReply(cipService_Read, typ, byte(0), result)

}

func (h *serverTCPHandler) sendUnconnectedRRDataReply(s CIPService, payload ...any) error {
	items := make([]cipItem, 2)
	items[0] = NewItem(cipItem_Null, nil)
	items[1] = NewItem(cipItem_UnconnectedData, nil)
	resp := msgUnconnWriteResultHeader{
		Service: s.AsResponse(),
	}
	items[1].Marshal(resp)
	for i := range payload {
		items[1].Marshal(payload[i])
	}
	return h.send(cipCommandSendRRData, MarshalItems(items))
}

func (h *serverTCPHandler) sendUnconnectedUnitDataReply(s CIPService) error {
	items := make([]cipItem, 2)
	items[0] = NewItem(cipItem_Null, nil)
	items[1] = NewItem(cipItem_UnconnectedData, nil)
	resp := msgWriteResultHeader{
		SequenceCount: h.UnitDataSequencer,
		Service:       s.AsResponse(),
	}
	items[1].Marshal(resp)
	return h.send(cipCommandSendUnitData, MarshalItems(items))
}
