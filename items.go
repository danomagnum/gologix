package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type CIPItemID uint16

const (
	CIPItem_Null            CIPItemID = 0
	CIPItem_UnconnectedData CIPItemID = 0xB2
)

func ReadItems(r io.Reader) ([]CIPItem, error) {

	var count uint16

	err := binary.Read(r, binary.LittleEndian, &count)
	if err != nil {
		return nil, fmt.Errorf("couldn't read item count. %w", err)
	}

	items := make([]CIPItem, 0, 2) // usually have 2 items.

	for i := 0; i < int(count); i++ {
		var hdr CIPItemHeader
		err := binary.Read(r, binary.LittleEndian, &hdr)
		if err != nil {
			return nil, fmt.Errorf("couldn't read item %d header. %w", i, err)
		}
		var item CIPItem
		item.Data = make([]byte, hdr.Length)
		err = binary.Read(r, binary.LittleEndian, &item.Data)
		if err != nil {
			return nil, fmt.Errorf("couldn't read item %d (hdr: %+v) data. %w", i, hdr, err)
		}
		items = append(items, item)
	}
	return items, nil
}

type CIPItem struct {
	Header CIPItemHeader
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
		Header: CIPItemHeader{
			ID: id,
		},
	}
	c.Pack(str)
	return c
}

// serialize a sturcture into the item's data.  Updates the header accordingly
func (item *CIPItem) Pack(str any) {
	binary.Write(item, binary.LittleEndian, str)
}

func (item *CIPItem) Bytes() []byte {
	b := bytes.Buffer{}
	binary.Write(&b, binary.LittleEndian, item.Header)
	b.Write(item.Data)
	return b.Bytes()
}

func (item *CIPItem) Reset() {
	item.Pos = 0
}

type CIPItemHeader struct {
	ID     CIPItemID
	Length uint16
}

type CIPItemsHeader struct {
	InterfaceHandle uint32
	SequenceCounter uint16
	Count           uint16
}

func BuildItemsBytes(items []CIPItem) *[]byte {

	b := new(bytes.Buffer)

	item_hdr := CIPItemsHeader{
		InterfaceHandle: 0,
		SequenceCounter: 0,
		Count:           uint16(len(items)),
	}
	binary.Write(b, binary.LittleEndian, item_hdr)

	for _, item := range items {
		b.Write(item.Bytes())
	}

	out := b.Bytes()

	return &out
}
