package gologix

import "encoding/binary"

// Here are the objects

type CIPElement uint32

// Used to indicate how many bytes are used for the data. If they are more than 8 bits,
// they actually actually take n+1 bytes.  First byte after specifier is a 0
type cipElementType byte

const (
	cipElement_8bit  cipElementType = 0x28
	cipElement_16bit cipElementType = 0x29
	cipElement_32bit cipElementType = 0x2A
)

func (p CIPElement) Bytes() []byte {
	if p < 256 {
		b := make([]byte, 2)
		b[0] = byte(cipElement_8bit)
		b[1] = byte(p)
		return b
	} else if p < 65536 {

		b := make([]byte, 4)
		b[0] = byte(cipElement_16bit)
		binary.LittleEndian.PutUint16(b[2:], uint16(p))
		return b
	} else {

		b := make([]byte, 6)
		b[0] = byte(cipElement_16bit)
		binary.LittleEndian.PutUint32(b[2:], uint32(p))
		return b
	}
}
