package gologix

import "encoding/binary"

// Used to indicate how many bytes are used for the data. If they are more than 8 bits,
// they actually actually take n+1 bytes.  First byte after specifier is a 0
type cipInstanceSize byte

const (
	cipInstance_8bit  cipInstanceSize = 0x24
	cipInstance_16bit cipInstanceSize = 0x25
)

type CIPInstance uint16

func (p CIPInstance) Bytes() []byte {
	if p < 256 {
		b := make([]byte, 2)
		b[0] = byte(cipInstance_8bit)
		b[1] = byte(p)
		return b
	} else {

		b := make([]byte, 4)
		b[0] = byte(cipInstance_16bit)
		binary.LittleEndian.PutUint16(b[2:], uint16(p))
		return b
	}
}
func (p CIPInstance) Len() int {
	if p < 256 {
		return 2
	}
	return 4
}

type JustBytes []byte

func (p JustBytes) Bytes() []byte {
	if len(p) == 1 {
		b := make([]byte, len(p)+1)
		b[0] = byte(cipInstance_8bit)
		copy(b[1:], p)
		return b
	} else {
		b := make([]byte, len(p)+2)
		b[0] = byte(cipInstance_16bit)
		copy(b[2:], p)
		return b
	}

}
