package gologix

import (
	"bytes"
	"fmt"
	"sync"
)

type IOProvider[Tin, Tout any] struct {
	Mutex sync.Mutex
	Data  map[string]any
	In    *Tin
	Out   *Tout
}

var io_read_test_counter byte = 0

// this gets called with the IO setup forward open as the items
func (p *IOProvider[Tin, Tout]) IORead() ([]byte, error) {
	io_read_test_counter++
	b := bytes.Buffer{}
	_ = Pack(&b, CIPPack{}, *(p.In))
	dat := b.Bytes()
	return dat, nil
}

func (p *IOProvider[Tin, Tout]) IOWrite(items []cipItem) error {
	if len(items) != 2 {
		return fmt.Errorf("expeted 2 items but got %v", len(items))
	}
	if items[1].Header.ID != cipItem_ConnectedData {
		return fmt.Errorf("expeted item 2 to be a connected data item but got %v", items[1].Header.ID)
	}
	var seq_counter uint32

	// according to wireshark only the least significant 4 bits are used.
	// 00.. ROO?
	// ..0. COO?
	// ...1 // Run/Idle (1 = run)
	var header uint16
	err := items[1].DeSerialize(&seq_counter)
	if err != nil {
		return fmt.Errorf("problem getting sequence counter. %w", err)
	}
	err = items[1].DeSerialize(&header)
	if err != nil {
		return fmt.Errorf("problem getting header. %w", err)
	}
	payload := make([]byte, items[1].Header.Length-6)
	err = items[1].DeSerialize(&payload)
	if err != nil {
		return fmt.Errorf("problem getting raw data. %w", err)
	}
	b := bytes.NewBuffer(payload)
	_, err = Unpack(b, CIPPack{}, p.Out)
	if err != nil {
		return fmt.Errorf("problem unpacking data into output struct %w", err)
	}

	return nil
}

func (p *IOProvider[Tin, Tout]) TagRead(tag string, qty int16) (any, error) {
	return 0, nil
}

func (p *IOProvider[Tin, Tout]) TagWrite(tag string, value any) error {
	return nil
}
