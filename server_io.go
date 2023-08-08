package gologix

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
)

// this type satisfies the TagProvider interface to provide class 1 IO support
// It has to be defined with an input and output struct that consist of only GoLogixTypes
// It will then serialize the input data and send it to the PLC at the requested rate.
// When the PLC sends an IO output message, that gets deserialized into the output structure.
//
// If you are going to access In or Out, be sure to lock the appropriate Mutex first to prevent data race.
// Remember that they are pointers here so the locks also need to apply to the original data that was pointed at.
//
// it does not handle class 3 tag reads or writes.
type IOProvider[Tin, Tout any] struct {
	InMutex  sync.Mutex
	OutMutex sync.Mutex
	In       *Tin
	Out      *Tout
}

var io_read_test_counter byte = 0

// this gets called with the IO setup forward open as the items
func (p *IOProvider[Tin, Tout]) IORead() ([]byte, error) {
	p.InMutex.Lock()
	defer p.InMutex.Unlock()
	io_read_test_counter++
	b := bytes.Buffer{}
	_ = Pack(&b, CIPPack{}, *(p.In))
	dat := b.Bytes()
	return dat, nil
}

func (p *IOProvider[Tin, Tout]) IOWrite(items []CIPItem) error {
	if len(items) != 2 {
		return fmt.Errorf("expeted 2 items but got %v", len(items))
	}
	if items[1].Header.ID != cipItem_ConnectedData {
		return fmt.Errorf("expected item 2 to be a connected data item but got %v", items[1].Header.ID)
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

	p.OutMutex.Lock()
	defer p.OutMutex.Unlock()

	_, err = Unpack(b, CIPPack{}, p.Out)
	if err != nil {
		return fmt.Errorf("problem unpacking data into output struct %w", err)
	}

	return nil
}

func (p *IOProvider[Tin, Tout]) TagRead(tag string, qty int16) (any, error) {
	return 0, errors.New("not implemented")
}

func (p *IOProvider[Tin, Tout]) TagWrite(tag string, value any) error {
	return errors.New("not implemented")
}

// returns the most udpated copy of the output data
// this output data is what the PLC is writing to us
func (p *IOProvider[Tin, Tout]) GetOutputData() Tout {
	p.OutMutex.Lock()
	defer p.OutMutex.Unlock()
	t := *p.Out
	return t
}

// update the input data thread safely
// this input data is what the PLC receives
func (p *IOProvider[Tin, Tout]) SetInputData(newin Tin) {
	p.InMutex.Lock()
	defer p.InMutex.Unlock()
	p.In = &newin
}
