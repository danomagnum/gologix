package gologix

import (
	"bytes"
	"encoding/binary"
)

// todo: move sequence to a different struct and combine msgCIPIOIHeader and CIPMultiIOIHeader
type msgCIPIOIHeader struct {
	Sequence uint16
	Service  CIPService
	Size     byte
}

type msgCIPMultiIOIHeader struct {
	Service CIPService
	Size    byte
}

// this is the generic connected message.
// it goes into an item (always item[1]?) and is followed up with
// a valid path.  The item specifies the cipService that goes with the message
type msgCIPConnectedServiceReq struct {
	SequenceCount uint16
	Service       CIPService
	PathLength    byte
}

type msgCIPConnectedMultiServiceReq struct {
	Sequence     uint16
	Service      CIPService
	PathSize     byte
	Path         [4]byte
	ServiceCount uint16
}

type msgCIPWriteIOIFooter struct {
	DataType uint16
	Elements uint16
}

func (ftr msgCIPWriteIOIFooter) Bytes() []byte {
	if ftr.DataType == uint16(CIPTypeSTRING) {
		b := []byte{0xA0, 0x02, 0xCE, 0x0F, 0x00, 0x00}
		binary.LittleEndian.PutUint16(b[4:], ftr.Elements)
		return b
	}

	b := bytes.Buffer{}
	binary.Write(&b, binary.LittleEndian, ftr)
	return b.Bytes()

}
func (ftr msgCIPWriteIOIFooter) Len() int {
	if ftr.DataType == uint16(CIPTypeSTRING) {
		return 6
	}

	return 4
}

type msgCIPIOIFooter struct {
	Elements uint16
}

type msgCIPResultHeader struct {
	InterfaceHandle uint32
	Timeout         uint16
}

// This should be everything before the actual result value data
// so you can read this off the buffer and be in the correct position to
// read the actual value as the type indicated by Type
type msgCIPReadResultData struct {
	SequenceCounter uint16
	Service         CIPService
	Status          [3]byte
	Type            CIPType
	Unknown         byte
}
