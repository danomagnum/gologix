package gologix

import "fmt"

func (h *handler) unconnectedData(item cipItem) error {
	var service CIPService
	var err error
	item.Unmarshal(&service)
	switch service {
	case cipService_ForwardOpen:
		item.Reset()
		err = h.forwardOpen(item)
		if err != nil {
			return fmt.Errorf("problem handling forward open. %w", err)
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
		fmt.Printf("Path 312: %v", path)
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

		}
	}
	return nil
}

func (h *handler) unconnectedServiceWrite(item cipItem) error {
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
		fmt.Printf("read %s as %s * %v = %v", tag, typ, qty, item.Data[item.Pos:])
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
	fmt.Printf("Path 211: %v\n", path)
	fmt.Printf("read %s to %v as %s * %v = %v", tag, path, typ, qty, results)

	provider, err := h.server.Router.Resolve(path)
	if err != nil {
		return fmt.Errorf("problem finding tag provider for %v. %w", path, err)
	}
	p := provider
	if qty > 1 {
		p.Write(tag, results)
	} else {
		p.Write(tag, results[0])
	}

	return h.sendUnconnectedRRDataReply(cipService_Write)

}
