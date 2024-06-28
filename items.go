package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type CIPItemID uint16

const (
	cipItem_Null                 CIPItemID = 0x0000
	cipItem_ListIdentityResponse CIPItemID = 0x000C
	cipItem_ConnectionAddress    CIPItemID = 0x00A1
	cipItem_ConnectedData        CIPItemID = 0x00B1
	cipItem_UnconnectedData      CIPItemID = 0x00B2
	cipItem_ListServiceResponse  CIPItemID = 0x0100
	cipItem_SockAddrInfo_OT      CIPItemID = 0x8000 // socket address info
	cipItem_SockAddrInfo_TO      CIPItemID = 0x8001 // socket address info
	cipItem_SequenceAddress      CIPItemID = 0x8002
)

// readItems takes an io.Reader positioned at the count of items in the data stream.
// It then reads each item from the data stream into an Item structure and returns a slice of all items.
func readItems(r io.Reader) ([]CIPItem, error) {

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

// The CIPItem is one of the core abstractions this library uses.
//
// When a response comes back from the controller it is structured in a CIPItem which can then
// be used to deserialize it with its various methods.
//
// There are methods such as Int16(), Uint16(), Int32(), Float32(), etc... that read the type specified off the item and advance
// the item buffer position.  This can be used to read one-at-a-time the data from the item.  Alternatively, you can use the
// Serialize() method to dump data into a pre-defined struct if that is more convenient for your application.
//
// This class also satisfies the io.reader and io.writer interfaces so you can use it for those kind of operations if needed.
type CIPItem struct {
	Header cipItemHeader
	Data   []byte
	Pos    int
}

// Allows the CIPItem to behave as an io.Reader
//
// Reads bytes from the current position in the item's buffer to the end or
// until we reach the requested read size.
func (item *CIPItem) Read(p []byte) (n int, err error) {
	if item.Pos >= len(item.Data) {
		return 0, io.EOF
	}
	n = copy(p, item.Data[item.Pos:])
	item.Pos += n
	return
}

// Allows the CIPItem to behave as an io.Writer.
//
// Appends bytes to the end of the item's buffer.
func (item *CIPItem) Write(p []byte) (n int, err error) {
	item.Data = append(item.Data, p...)
	n = len(p)
	item.Header.Length = uint16(len(item.Data))
	return
}

// create an item given an item id and data structure
func newItem(id CIPItemID, str any) CIPItem {
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

// returns all unprocessed bytes remaining in the item's buffer.
func (item *CIPItem) Rest() []byte {
	return item.Data[item.Pos:]
}

// Retrieve the next byte in the item's buffer and increment the buffer position.
func (item *CIPItem) Byte() (byte, error) {
	if len(item.Data) <= item.Pos {
		return 0, fmt.Errorf("item out of data")
	}
	b := item.Data[item.Pos]
	item.Pos++
	return b, nil
}

// Retrieve 2 bytes from the buffer and increment the buffer position by 2.  The
// data is interpreted as a 16 bit unsigned integer.
func (item *CIPItem) Uint16() (uint16, error) {
	if len(item.Data) <= item.Pos+1 {
		return 0, fmt.Errorf("item out of data")
	}
	val := binary.LittleEndian.Uint16(item.Data[item.Pos:])
	item.Pos += 2
	return val, nil
}

// Retrieve 2 bytes from the buffer and increment the buffer position by 2.  The
// data is interpreted as a 16 bit signed integer.
func (item *CIPItem) Int16() (int16, error) {
	if len(item.Data) <= item.Pos+1 {
		return 0, fmt.Errorf("item out of data")
	}
	val := binary.LittleEndian.Uint16(item.Data[item.Pos:])
	item.Pos += 2
	return int16(val), nil
}

// Retrieve 4 bytes from the buffer and increment the buffer position by 4.  The
// data is interpreted as a 32 bit unsigned integer.
func (item *CIPItem) Uint32() (uint32, error) {
	if len(item.Data) <= item.Pos+3 {
		return 0, fmt.Errorf("item out of data")
	}
	val := binary.LittleEndian.Uint32(item.Data[item.Pos:])
	item.Pos += 4
	return val, nil
}

// Retrieve 4 bytes from the buffer and increment the buffer position by 4.  The
// data is interpreted as a 32 bit signed integer.
func (item *CIPItem) Int32() (int32, error) {
	if len(item.Data) <= item.Pos+3 {
		return 0, fmt.Errorf("item out of data")
	}
	val := binary.LittleEndian.Uint32(item.Data[item.Pos:])
	item.Pos += 4
	return int32(val), nil
}

// Retrieve 8 bytes from the buffer and increment the buffer position by 8.  The
// data is interpreted as a 64 bit unsigned integer.
func (item *CIPItem) Uint64() (uint64, error) {
	if len(item.Data) <= item.Pos+7 {
		return 0, fmt.Errorf("item out of data")
	}
	val := binary.LittleEndian.Uint64(item.Data[item.Pos:])
	item.Pos += 8
	return val, nil
}

// Retrieve 8 bytes from the buffer and increment the buffer position by 8.  The
// data is interpreted as a 64 bit signed integer.
func (item *CIPItem) Int64() (int64, error) {
	if len(item.Data) <= item.Pos+7 {
		return 0, fmt.Errorf("item out of data")
	}
	val := binary.LittleEndian.Uint64(item.Data[item.Pos:])
	item.Pos += 8
	return int64(val), nil
}

// Retrieve 4 bytes from the buffer and increment the buffer position by 4.  The
// data is interpreted as a 32 bit floating point number.
func (item *CIPItem) Float32() (float32, error) {
	if len(item.Data) <= item.Pos+3 {
		return 0, fmt.Errorf("item out of data")
	}
	var val float32
	err := binary.Read(item, binary.LittleEndian, &val)
	return val, err
}

// Retrieve 8 bytes from the buffer and increment the buffer position by 8.  The
// data is interpreted as a 64 bit floating point number.
func (item *CIPItem) Float64() (float64, error) {
	if len(item.Data) <= item.Pos+7 {
		return 0, fmt.Errorf("item out of data")
	}
	var val float64
	err := binary.Read(item, binary.LittleEndian, &val)
	return val, err
}

// Serialize a structure into the item's data.
//
// If called more than once the []byte data for the additional structures is appended to the
// end of the item's data buffer.
//
// The data length in the item's header is updated to match.
func (item *CIPItem) Serialize(str any) error {
	switch x := str.(type) {
	case string:
		strLen := uint32(len(x))
		err := binary.Write(item, binary.LittleEndian, strLen)
		if err != nil {
			return fmt.Errorf("problem writing string header: %v", err)
		}
		if strLen%2 == 1 {
			strLen++
		}
		//b := make([]byte, strLen)
		b := make([]byte, 84)
		copy(b, x)
		err = binary.Write(item, binary.LittleEndian, b)
		if err != nil {
			return fmt.Errorf("problem writing string payload: %v", err)
		}

	case Serializable:
		err := binary.Write(item, binary.LittleEndian, x.Bytes())
		if err != nil {
			return fmt.Errorf("problem writing serializable item: %v", err)
		}
	default:
		err := binary.Write(item, binary.LittleEndian, str)
		if err != nil {
			return fmt.Errorf("problem writing default item: %v", err)
		}
	}
	return nil
}

// DeSerialize an item's data into the given structure.
//
// The position in the item's buffer is updated to account for the number of bytes
// required.
func (item *CIPItem) DeSerialize(str any) error {
	return binary.Read(item, binary.LittleEndian, str)
}

func (item *CIPItem) Bytes() ([]byte, error) {
	b := bytes.Buffer{}
	err := binary.Write(&b, binary.LittleEndian, item.Header)
	if err != nil {
		return b.Bytes(), fmt.Errorf("problem writing data. %v", err)
	}
	_, err = b.Write(item.Data)
	if err != nil {
		return b.Bytes(), fmt.Errorf("problem writing item data. %v", err)
	}
	return b.Bytes(), nil
}

// Sets the items data position back to zero without removing the data.
// Can be used to overwrite the item's internal data or to re-read the item's data
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

// serializeItems takes a slice of items and generates the appropriate byte pattern for the packet
//
// A lot of the time, item0 ends up being the "null" item with no Data section.
//
// A typical item structure will look like this:
//
//	byte	info        	Field
//	0   	Items Header	InterfaceHandle
//	1
//	2
//	3
//	4                   	SequenceCounter
//	5
//	6                   	ItemCount
//	7   	Item0 Header	Item ID
//	8
//	9                   	Length (bytes) = N0
//	10
//	11  	Item0 Data   	Byte 0
//	...
//	11+N0	Item1 Header	Item ID
//	12+N0
//	13+N0	            	Length (bytes) = N1
//	14+N0
//	15+N0	Item1 Data   	Byte 0
//	...  repeat for all items...
func serializeItems(items []CIPItem) (*[]byte, error) {

	b := new(bytes.Buffer)

	item_hdr := cipItemsHeader{
		InterfaceHandle: 0,
		SequenceCounter: 0,
		Count:           uint16(len(items)),
	}
	err := binary.Write(b, binary.LittleEndian, item_hdr)
	if err != nil {
		return nil, fmt.Errorf("problem serializing item header into b. %v", err)
	}

	for i, item := range items {
		b2, err := item.Bytes()
		if err != nil {
			return nil, fmt.Errorf("problem byte-ing item %d: %w", i, err)
		}
		b.Write(b2)
	}

	out := b.Bytes()

	return &out, nil
}
