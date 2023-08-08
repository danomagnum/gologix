package gologix

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"sync"
)

// this type satisfies the TagChannelProvider interface to provide class 1 IO support
// It has to be defined with an input and output struct that consist of only GoLogixTypes
// It will then serialize the input data and send it to the PLC at the requested rate.
// When the PLC sends an IO output message, that gets deserialized and sent to all destination channels.
//
// get a new destintaion channel by calling GetOutputData.  You can then receive from this to get data as it comes in.
//
// update the input data with the SetInputData(Tin) function
//
// it does not handle class 3 tag reads or writes.
type IOChannelProvider[Tin, Tout any] struct {
	inMutex     sync.Mutex
	in          Tin
	outChannels []chan Tout
}

// this gets called with the IO setup forward open as the items
func (p *IOChannelProvider[Tin, Tout]) IORead() ([]byte, error) {
	p.inMutex.Lock()
	defer p.inMutex.Unlock()
	b := bytes.Buffer{}
	_ = Pack(&b, CIPPack{}, p.in)
	dat := b.Bytes()
	return dat, nil
}

func (p *IOChannelProvider[Tin, Tout]) IOWrite(items []CIPItem) error {
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

	var out Tout
	_, err = Unpack(b, CIPPack{}, &out)
	if err != nil {
		return fmt.Errorf("problem unpacking data into output struct %w", err)
	}
	for i := range p.outChannels {
		select {
		case p.outChannels[i] <- out:
		default:
			log.Printf("problem sending. channel full?")
		}
	}

	return nil
}

func (p *IOChannelProvider[Tin, Tout]) TagRead(tag string, qty int16) (any, error) {
	return 0, errors.New("not implemented")
}

func (p *IOChannelProvider[Tin, Tout]) TagWrite(tag string, value any) error {
	return errors.New("not implemented")
}

// returns the most udpated copy of the output data
// this output data is what the PLC is writing to us
func (p *IOChannelProvider[Tin, Tout]) GetOutputDataChannel() <-chan Tout {
	newout := make(chan Tout)
	p.outChannels = append(p.outChannels, newout)
	return newout
}

// update the input data thread safely
// this input data is what the PLC receives
func (p *IOChannelProvider[Tin, Tout]) SetInputData(newin Tin) {
	p.inMutex.Lock()
	defer p.inMutex.Unlock()
	p.in = newin
}
