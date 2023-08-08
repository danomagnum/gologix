package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type CIPItemID uint16

const (
	cipItem_Null                CIPItemID = 0x0000
	cipItem_ListIdentityReponse CIPItemID = 0x000C
	cipItem_ConnectionAddress   CIPItemID = 0x00A1
	cipItem_ConnectedData       CIPItemID = 0x00B1
	cipItem_UnconnectedData     CIPItemID = 0x00B2
	cipItem_ListServiceResponse CIPItemID = 0x0100
	cipItem_SockAddrInfo_OT     CIPItemID = 0x8000 // socket address info
	cipItem_SockAddrInfo_TO     CIPItemID = 0x8001 // socket address info
	cipItem_SequenceAddress     CIPItemID = 0x8002
)

// ReadItems takes an io.Reader positioned at the count of items in the data stream.
// It then reads each item from the data stream into an Item structure and returns a slice of all items.
func ReadItems(r io.Reader) ([]CIPItem, error) {

	var count uint16

	err := binary.Read(r, binary.LittleEndian, &count)
	if err != nil {
		return nil, fmt.Errorf("couldn't read item count. %w", err)
	}

	items := make([]CIPItem, count) // usually have 2 items.

	for i := 0; i < int(count); i++ {
		//var hdr cipItemHeader
		err := binary.Read(r, binary.LittleEndian, &(items[i].Header))
		if err != nil {
			return nil, fmt.Errorf("couldn't read item %d header. %w", i, err)
		}
		items[i].Data = make([]byte, items[i].Header.Length)
		err = binary.Read(r, binary.LittleEndian, &(items[i].Data))
		if err != nil {
			return nil, fmt.Errorf("couldn't read item %d (hdr: %+v) data. %w", i, items[i].Header, err)
		}
	}
	return items, nil
}

type CIPItem struct {
	Header cipItemHeader
	Data   []byte
	Pos    int
}

func (item *CIPItem) Read(p []byte) (n int, err error) {
	if item.Pos >= len(item.Data) {
		return 0, io.EOF
	}
	n = copy(p, item.Data[item.Pos:])
	item.Pos += n
	return
}
func (item *CIPItem) Write(p []byte) (n int, err error) {
	item.Data = append(item.Data, p...)
	n = len(p)
	item.Header.Length = uint16(len(item.Data))
	return
}

// create an item given an item id and data structure
func NewItem(id CIPItemID, str any) CIPItem {
	c := CIPItem{
		Header: cipItemHeader{
			ID: id,
		},
	}
	if str != nil {
		c.Serialize(str)
	}
	return c
}

// Serialize a sturcture into the item's data.
//
// If called more than once the []byte data for the additional structures is appended to the
// end of the item's data buffer.
//
// The data length in the item's header is updated to match.
func (item *CIPItem) Serialize(str any) {
	switch x := str.(type) {
	case Serializable:
		err := binary.Write(item, binary.LittleEndian, x.Bytes())
		if err != nil {
			log.Printf("Problem writing serializable item: %v", err)
		}
	default:
		err := binary.Write(item, binary.LittleEndian, str)
		if err != nil {
			log.Printf("Problem writing default item: %v", err)
		}
	}
}

// DeSerialize an item's data into the given sturcture.
//
// If called more than once the []byte data for the additional structures is continuously
// read from the current position of the item's data buffer.
func (item *CIPItem) DeSerialize(str any) error {
	return binary.Read(item, binary.LittleEndian, str)
}

func (item *CIPItem) Bytes() []byte {
	b := bytes.Buffer{}
	err := binary.Write(&b, binary.LittleEndian, item.Header)
	if err != nil {
		log.Printf("problem writing data. %v", err)
		return b.Bytes()
	}
	_, err = b.Write(item.Data)
	if err != nil {
		log.Printf("problem writing item data. %v", err)
		return b.Bytes()
	}
	return b.Bytes()
}

// Sets the items data position back to zero.  Can be used to overrite the item's internal data or to re-read the item's data
func (item *CIPItem) Reset() {
	item.Pos = 0
}

// This is the header for a single item in an item list.
type cipItemHeader struct {
	ID     CIPItemID
	Length uint16 // bytes of data to follow
}

// This is the header for multiple items
type cipItemsHeader struct {
	InterfaceHandle uint32
	SequenceCounter uint16
	Count           uint16
}

// SerializeItems takes a slice of items and generates the appropriate byte pattern for the packet
//
// A lot of the time, item0 ends up being the "null" item with no Data section.
//
// A typical item structure will look like this:
// byte		info          	Field
// 0		Items Header	InterfaceHandle
// 1
// 2
// 3
// 4	             		SequenceCounter
// 5
// 6	               		ItemCount
// 7		Item0 Header	Item ID
// 8
// 9 	            		Length (bytes) = N0
// 10
// 11 		Item0 Data   	Byte 0
// ...
// 11+N0	Item1 Header	Item ID
// 12+N0
// 13+N0	           		Length (bytes) = N1
// 14+N0
// 15+N0	Item1 Data   	Byte 0
// ...  repeat for all items...
func SerializeItems(items []CIPItem) *[]byte {

	b := new(bytes.Buffer)

	item_hdr := cipItemsHeader{
		InterfaceHandle: 0,
		SequenceCounter: 0,
		Count:           uint16(len(items)),
	}
	err := binary.Write(b, binary.LittleEndian, item_hdr)
	if err != nil {
		log.Printf("problem serializing item header into b. %v", err)
	}

	for _, item := range items {
		b.Write(item.Bytes())
	}

	out := b.Bytes()

	return &out
}
