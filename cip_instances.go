package gologix

import "encoding/binary"

// Used to indicate how many bytes are used for the data. If they are more than 8 bits,
// they actually actually take n+1 bytes.  First byte after specifier is a 0
type CIPInstanceSize byte

const (
	CIPInstance_8bit  CIPInstanceSize = 0x24
	CIPInstance_16bit CIPInstanceSize = 0x25
)

type CIPInstance uint16

func (p CIPInstance) Bytes() []byte {
	if p < 256 {
		b := make([]byte, 2)
		b[0] = byte(CIPInstance_8bit)
		b[1] = byte(p)
		return b
	} else {

		b := make([]byte, 4)
		b[0] = byte(CIPInstance_16bit)
		binary.LittleEndian.PutUint16(b[2:], uint16(p))
		return b
	}
}
