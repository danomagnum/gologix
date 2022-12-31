package gologix

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type IOProvider struct {
	Mutex sync.Mutex
	Data  map[string]any
}

var io_read_test_counter byte = 0

// this gets called with the IO setup forward open as the items
func (p *IOProvider) IORead() ([]byte, error) {
	io_read_test_counter++
	return []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, io_read_test_counter}, nil
}

// TODO: handle connection closing.
// it looks like a forward close will come in which has a connection serial number that should have been in the foward open to start with.
// we'll have to know that is tied to the correct io provider
// also if we don't get an input IO message in a certain time we should abandon the output IO messages
func (p *IOProvider) ioRead(fwd_open msgEIPForwardOpen_Standard, rpi time.Duration) {
	dat := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	t := time.NewTicker(rpi)
	seq := uint32(0)
	for {
		seq++
		<-t.C
		// every RPI send the message.
		items := make([]cipItem, 2)
		items[0] = NewItem(cipItem_SequenceAddress, nil)
		items[0].Marshal(fwd_open.TOConnectionID)
		items[0].Marshal(seq)
		items[1] = NewItem(cipItem_ConnectedData, nil)
		items[1].Marshal(uint16(seq))
		//items[1].Marshal(uint32(1)) // connection properties. 1 = running. (not used on response)
		items[1].Marshal(dat)

		// TODO: stop hardcoding the IP address.  Get it from the connection handler that starts the comms.
		conn, err := net.Dial("udp", "192.168.2.241:2222")
		if err != nil {
			log.Printf("problem connecting UDP. %v", err)
			continue
		}
		payload := *MarshalItems(items)
		payload = payload[6:]
		log.Printf("writing udp io payload %v", payload)
		_, err = conn.Write(payload)
		if err != nil {
			log.Printf("problem writing %v", err)
		}
		conn.Close()

	}

}

func (p *IOProvider) IOWrite(items []cipItem) error {
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
	err := items[1].Unmarshal(&seq_counter)
	if err != nil {
		return fmt.Errorf("problem getting sequence counter. %w", err)
	}
	err = items[1].Unmarshal(&header)
	if err != nil {
		return fmt.Errorf("problem getting header. %w", err)
	}
	payload := make([]byte, items[1].Header.Length-6)
	err = items[1].Unmarshal(&payload)
	if err != nil {
		return fmt.Errorf("problem getting raw data. %w", err)
	}
	log.Printf("got IO input %v", payload)

	return nil
}

func (p *IOProvider) TagRead(tag string, qty int16) (any, error) {
	return 0, nil
}

func (p *IOProvider) TagWrite(tag string, value any) error {
	return nil
}
