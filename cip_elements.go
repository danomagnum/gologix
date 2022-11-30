package gologix

import "encoding/binary"

// Here are the objects

type CIPElement uint32

// Used to indicate how many bytes are used for the data. If they are more than 8 bits,
// they actually actually take n+1 bytes.  First byte after specifier is a 0
type CIPElementType byte

const (
	CIPElement_8bit  CIPElementType = 0x28
	CIPElement_16bit CIPElementType = 0x29
	CIPElement_32bit CIPElementType = 0x2A
)

func (p CIPElement) Bytes() []byte {
	if p < 256 {
		b := make([]byte, 2)
		b[0] = byte(CIPElement_8bit)
		b[1] = byte(p)
		return b
	} else if p < 65536 {

		b := make([]byte, 4)
		b[0] = byte(CIPElement_16bit)
		binary.LittleEndian.PutUint16(b[2:], uint16(p))
		return b
	} else {

		b := make([]byte, 6)
		b[0] = byte(CIPElement_16bit)
		binary.LittleEndian.PutUint32(b[2:], uint32(p))
		return b
	}
}
