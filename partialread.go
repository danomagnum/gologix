package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/danomagnum/gologix/cipservice"
	"github.com/danomagnum/gologix/ciptype"
)

// Read a list of tags of specified types.
//
// The result slice will be in the same order as the tag list.  Each value in the list will be an
// interface{} so you'll need to type assert to get the values back out.
//
// To read multiple tags at once without type assertion you can use ReadMulti()
func (client *Client) ReadList(tags []string, types []ciptype.CIPType) ([]any, error) {
	err := client.checkConnection()
	if err != nil {
		return nil, fmt.Errorf("could not start list read: %w", err)
	}
	n := 0
	n_new := 0
	total := len(tags)
	results := make([]any, 0, total)
	msgs := 0

	for n < total {
		msgs += 1
		n_new, err = client.countIOIsThatFit(tags[n:], types[n:])
		if err != nil {
			return nil, err
		}
		subresults, err := client.readList(tags[n:n+n_new], types[n:n+n_new])
		n += n_new
		if err != nil {
			return nil, err
		}
		results = append(results, subresults...)

	}

	log.Printf("Took %d messages to read %d tags", msgs, n)
	return results, nil
}

func (client *Client) countIOIsThatFit(tags []string, types []ciptype.CIPType) (int, error) {
	// first generate IOIs for each tag
	qty := len(tags)

	ioi_header := msgCIPConnectedMultiServiceReq{
		Sequence:     uint16(sequencer()),
		Service:      cipservice.MultipleService,
		PathSize:     2,
		Path:         [4]byte{0x20, 0x02, 0x24, 0x01},
		ServiceCount: uint16(qty),
	}

	mainhdr_size := binary.Size(ioi_header)
	ioihdr_size := binary.Size(msgCIPMultiIOIHeader{})
	ioiftr_size := binary.Size(msgCIPIOIFooter{})

	b := bytes.Buffer{}
	// we now have to build up the jump table for each IOI.
	// and pack all the IOIs together into b
	jump_table := make([]uint16, qty)

	// how many ioi's fit in the message
	n := 1

	response_size := 0

	for i, tag := range tags {
		ioi, err := client.NewIOI(tag, types[i])
		if err != nil {
			return 0, err
		}

		jump_table[i] = uint16(b.Len())

		h := msgCIPMultiIOIHeader{
			Service: cipservice.Read,
			Size:    byte(len(ioi.Buffer) / 2),
		}
		f := msgCIPIOIFooter{
			Elements: 1,
		}

		// Calculate the size of the data once we add this ioi to the list.
		newSize := mainhdr_size                                // length of the multi-read header
		newSize += 2 * n                                       // add in the jump table
		newSize += b.Len()                                     // everything we have so far
		newSize += ioihdr_size + len(ioi.Buffer) + ioiftr_size // the new ioi data

		response_size += types[i].Size()
		if newSize > client.ConnectionSize || response_size > client.ConnectionSize {
			// break before adding this ioi to the list since it will push us over.
			// we'll continue with n iois (n only increments after an IOI is added)
			break
		}

		err = binary.Write(&b, binary.LittleEndian, h)
		if err != nil {
			return 0, fmt.Errorf("problem writing cip IO header to buffer. %w", err)
		}
		b.Write(ioi.Buffer)
		err = binary.Write(&b, binary.LittleEndian, f)
		if err != nil {
			return 0, fmt.Errorf("problem writing ioi buffer to msg buffer. %w", err)
		}

		n = i + 1
	}

	log.Printf("Fit %d tags into %d bytes.  ", n, client.ConnectionSize)

	return n, nil

}
