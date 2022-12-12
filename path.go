package gologix

// based on code from https://github.com/loki-os/go-ethernet-ip

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

// The path is formatted like this.
// byte 0: number of 16 bit words
// byte 1: 000. .... path segment type (port segment = 0)
// byte 1: ...0 .... extended link address (0 = false)
// byte 1: .... 0001 port (backplane = 1)
// byte 2: n/a
// byte 3: 001. .... path segment type (logical segment = 1)
// byte 3: ...0 00.. logical segment type class ID (0)
// byte 3: .... ..00 logical segment format: 8-bit (0)
// byte 4: path segment 0x20
// byte 5: 001. .... path segment type (logical segment = 1)
// byte 5: ...0 01.. logical segment type: Instance ID = 1
// byte 5: .... ..00 logical segment format: 8-bit (0)
// byte 6: path segment instance 0x01
// so on...
//msg.Path = [6]byte{0x01, 0x00, 0x20, 0x02, 0x24, 0x01}

// bits 5,6,7 (counting from 0) are the segment type
type SegmentType byte

const (
	SegmentTypeExtendedSymbolic SegmentType = 0x91
)

type CIPPort struct {
	PortNo       byte
	ExtensionLen byte
}

func (p CIPPort) Bytes() []byte {
	if p.ExtensionLen != 0 {
		return []byte{p.PortNo, p.ExtensionLen}

	}
	return []byte{p.PortNo}
}

// this function takes a CIP path in the format of 0,1,192.168.2.1,0,1 and converts it into the proper equivalent byte slice.
func ParsePath(path string) (*bytes.Buffer, error) {
	// get rid of any spaces and square brackets
	path = strings.ReplaceAll(path, " ", "")
	path = strings.ReplaceAll(path, "[", "")
	path = strings.ReplaceAll(path, "]", "")
	// split on commas
	parts := strings.Split(path, ",")

	byte_path := make([]byte, 0, len(parts))

	for _, part := range parts {
		// first see if this looks like an IP address.
		is_ip := strings.Contains(part, ".")
		if is_ip {
			// for some god forsaken reason the path doesn't use the ip address as actual bytes but as an ascii string.
			// we first have to set bit 5 in the previous byte to say we're using an extended address for this part.
			last_pos := len(byte_path) - 1
			last_byte := byte_path[last_pos]
			byte_path[last_pos] = last_byte | 1<<4
			l := len(part)
			byte_path = append(byte_path, byte(l))
			string_bytes := []byte(part)
			byte_path = append(byte_path, string_bytes...)
			continue
		}
		// not an IP address
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("problem converting %v to number. %w", part, err)
		}
		if val < 0 || val > 255 {
			return nil, fmt.Errorf("number out of range. %v", part)
		}
		byte_path = append(byte_path, byte(val))
	}

	return bytes.NewBuffer(byte_path), nil
}

type Serializable interface {
	Bytes() []byte
}

// given a list of structures, serialize them with the Bytes() method if available,
// otherwise serialize with binary.write()
func Serialize(strs ...any) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	for _, str := range strs {
		switch serializable_str := str.(type) {
		case Serializable:
			// if the struct is serializable, we should use its Bytes() function to get its
			// representation instead of binary.Write
			_, err := b.Write(serializable_str.Bytes())
			if err != nil {
				return nil, err
			}
		case any:
			binary.Write(b, binary.LittleEndian, str)
		}
	}
	return b, nil
}
