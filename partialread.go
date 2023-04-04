package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

func (client *Client) ReadListPartial(tags []string, types []CIPType) ([]any, error) {
	n := 0
	n_new := 0
	var err error
	var req []cipItem
	total := len(tags)
	results := make([]any, total)
	msgs := 0

	for n < total {
		msgs += 1
		n_new, req, err = client.GetIOIsThatFit(tags[n:], types[n:])
		n += n_new
		if err != nil {
			return nil, err
		}
		_ = req     // TODO: do read and get results
		_ = results // TODO: append results
	}
	log.Printf("Took %d messages to read %d tags", msgs, n)
	return results, nil
}

func (client *Client) GetIOIsThatFit(tags []string, types []CIPType) (int, []cipItem, error) {
	// first generate IOIs for each tag
	qty := len(tags)

	ioi_header := msgCIPConnectedMultiServiceReq{
		Sequence:     client.Sequencer(),
		Service:      cipService_MultipleService,
		PathSize:     2,
		Path:         [4]byte{0x20, 0x02, 0x24, 0x01},
		ServiceCount: uint16(qty),
	}

	reqitems := make([]cipItem, 2)
	reqitems[0] = NewItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	mainhdr_size := binary.Size(ioi_header)
	ioihdr_size := binary.Size(msgCIPMultiIOIHeader{})
	ioiftr_size := binary.Size(msgCIPIOIFooter{})
	item0_size := len(reqitems[0].Bytes())

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
			return 0, nil, err
		}

		jump_table[i] = uint16(b.Len())

		h := msgCIPMultiIOIHeader{
			Service: cipService_Read,
			Size:    byte(len(ioi.Buffer) / 2),
		}
		f := msgCIPIOIFooter{
			Elements: 1,
		}

		// Calculate the size of the data once we add this ioi to the list.
		newSize := item0_size                                  // length of the 0 item
		newSize += mainhdr_size                                // lenght of the multi-read header
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
			return 0, nil, fmt.Errorf("problem writing cip IO header to buffer. %w", err)
		}
		b.Write(ioi.Buffer)
		err = binary.Write(&b, binary.LittleEndian, f)
		if err != nil {
			return 0, nil, fmt.Errorf("problem writing ioi buffer to msg buffer. %w", err)
		}

		n = i + 1
	}

	// truncate the slices to the actual size that fit
	jump_table = jump_table[:n]

	// now that we know how long the jump table actually is, we need to go through
	// and offset the jmp table values by the length of the jump table.
	// because it currently holds the distance in to &b, but we need the distance from the start
	// of the jump table
	jumpTableSize := n * 2 // 2 bytes + 2 bytes per jump entry
	for i := 0; i < n; i++ {
		jump_table[i] += uint16(jumpTableSize)
	}

	reqitems[1] = cipItem{Header: cipItemHeader{ID: cipItem_ConnectedData}}
	reqitems[1].Serialize(ioi_header)
	reqitems[1].Serialize(jump_table)
	reqitems[1].Serialize(&b)

	log.Printf("Fit %d tags into %d bytes.  Total bytes: %d", n, client.ConnectionSize, binary.Size(SerializeItems(reqitems)))

	return n, reqitems, nil

}
