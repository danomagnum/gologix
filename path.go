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
type segmentType byte

const (
	segmentTypeExtendedSymbolic segmentType = 0x91
)

// represents the port number on a CIP device
//
// If you're going to serialize this class to bytes for transimssion be sure to use one of the gologix
// serialization functions or call Bytes() to get the properly formatted data.
type CIPPort struct {
	PortNo       byte
	ExtensionLen byte
}

func (p CIPPort) Len() int {
	return 2
}

func (p CIPPort) Bytes() []byte {
	if p.ExtensionLen != 0 {
		return []byte{p.PortNo, p.ExtensionLen}

	}
	return []byte{p.PortNo}
}

// This function takes a CIP path in the format of 0,1,192.168.2.1,0,1 and converts it into the proper equivalent byte slice.
//
// The most common use is probably setting up the communication path on a new client.
func ParsePath(path string) (*bytes.Buffer, error) {
	if path == "" {
		return new(bytes.Buffer), nil
	}
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

// PathBuilder constructs a CIP path incrementally using a fluent builder pattern.
// Each method appends one logical segment to the path and returns the same *PathBuilder,
// so calls can be chained:
//
//	path := new(gologix.PathBuilder).
//	    Class(gologix.CIPClass_MessageRouter).
//	    Instance(1).
//	    Bytes()
//
// If any method encounters an error the builder stores it in Err and all subsequent
// method calls become no-ops, so you only need to check Err once at the end.
type PathBuilder struct {
	b bytes.Buffer
	// Err holds the first error encountered while building the path.
	// All methods check this field before writing; once set, further calls are no-ops.
	Err error
}

// NewPathBuilder creates and returns a new PathBuilder with an empty path and no error.
// You could also just do a PathBuilder{} but if trying to use it directly you would need to make it a pointer like
// &(PathBuilder{}) to be able to chain methods, so this is a little more convenient.
func NewPathBuilder() *PathBuilder {
	return &PathBuilder{}
}

// Custom appends raw bytes directly to the path without any interpretation.
// This is useful when you need to include a pre-encoded segment that has no
// dedicated builder method.
func (pb *PathBuilder) Custom(data []byte) *PathBuilder {
	if pb.Err != nil {
		return pb
	}
	_, err := pb.b.Write(data)
	if err != nil {
		pb.Err = fmt.Errorf("error writing segment %v: %w", data, err)
	}
	return pb

}

// Bytes returns the accumulated path as a raw byte slice.
// Check Err before using the returned slice to confirm no error occurred during construction.
func (pb *PathBuilder) Bytes() []byte {
	return pb.b.Bytes()
}

// Class appends a logical class segment for the given CIPClass to the path.
// Class segments identify the object class being addressed (e.g. the Message Router class).
// There are a number of predefined classes in the gologix package.
// If using a constant, there is no need to cast it to CIPClass, e.g. Class(4) is fine.
// If using a variable, you will have to cast it to CIPClass, e.g. Class(gologix.CipClass(myVar)).
func (pb *PathBuilder) Class(classID CIPClass) *PathBuilder {
	if pb.Err != nil {
		return pb
	}
	_, err := pb.b.Write(classID.Bytes())
	if err != nil {
		pb.Err = fmt.Errorf("error writing class ID %v: %w", classID, err)
	}
	return pb
}

// Instance appends a logical instance segment for the given CIPInstance to the path.
// Instance segments select a specific instance of the class identified by the preceding Class segment.
// If using a constant, there is no need to cast it to CIPInstance, e.g. Instance(1) is fine.
// If using a variable, you will have to cast it to CIPInstance, e.g. Instance(gologix.CipInstance(myVar)).
func (pb *PathBuilder) Instance(instanceID CIPInstance) *PathBuilder {
	if pb.Err != nil {
		return pb
	}
	_, err := pb.b.Write(instanceID.Bytes())
	if err != nil {
		pb.Err = fmt.Errorf("error writing instance ID %v: %w", instanceID, err)
	}
	return pb
}

// Attribute appends a logical attribute segment for the given CIPAttribute to the path.
// Attribute segments target a specific attribute within an object instance.
// If using a constant, there is no need to cast it to CIPAttribute, e.g. Attribute(1) is fine.
// If using a variable, you will have to cast it to CIPAttribute, e.g. Attribute(gologix.CipAttribute(myVar)).
func (pb *PathBuilder) Attribute(attributeID CIPAttribute) *PathBuilder {
	if pb.Err != nil {
		return pb
	}
	_, err := pb.b.Write(attributeID.Bytes())
	if err != nil {
		pb.Err = fmt.Errorf("error writing attribute ID %v: %w", attributeID, err)
	}
	return pb
}

// Symbolic appends an ANSI extended symbolic segment for the given tag or symbol name.
// This is the standard way to address a named controller tag, e.g. "MyTag" or "Program:MyProg.MyTag".
// If a tag has nested paths (such as a member of a struct) you need to add a separate symbolic segment
// for each level of the path, e.g. "MyProg.MyTag.MyField" becomes Symbolic("MyProg").Symbolic("MyTag").Symbolic("MyField").
func (pb *PathBuilder) Symbolic(symbol string) *PathBuilder {
	if pb.Err != nil {
		return pb
	}
	str, err := marshalIOIPart(symbol)
	if err != nil {
		pb.Err = fmt.Errorf("error writing symbolic ID %v: %w", symbol, err)
	}
	pb.b.Write(str)
	return pb
}

// Element appends a logical element segment for the given CIPElement index.
// Element segments address a specific element within an array tag.
// Used for accessing array elements, e.g. MyArray[3] would be Symbolic("MyArray").Element(3).
func (pb *PathBuilder) Element(elementID CIPElement) *PathBuilder {
	if pb.Err != nil {
		return pb
	}
	_, err := pb.b.Write(elementID.Bytes())
	if err != nil {
		pb.Err = fmt.Errorf("error writing element ID %v: %w", elementID, err)
	}
	return pb
}

// Port appends a port segment with the given port number.
// Port segments route the message through a specific communication port on the device.
// Use together with Address to hop across network boundaries (e.g. through a ControlLogix backplane).
func (pb *PathBuilder) Port(portNo byte) *PathBuilder {
	if pb.Err != nil {
		return pb
	}
	err := binary.Write(&(pb.b), binary.LittleEndian, portNo)
	if err != nil {
		pb.Err = fmt.Errorf("error writing port number %v: %w", portNo, err)
	}
	return pb
}

// Address appends a link-address segment for the given CIPAddress (node address).
// This is typically used after a Port segment to specify the destination node on that port.
func (pb *PathBuilder) Address(address CIPAddress) *PathBuilder {
	if pb.Err != nil {
		return pb
	}
	_, err := pb.b.Write(address.Bytes())
	if err != nil {
		pb.Err = fmt.Errorf("error writing address %v: %w", address, err)
	}
	return pb
}

// Slot appends a backplane slot number to the path.
// This is a shorthand for the common case of routing to a module in a specific chassis slot
// (e.g. slot 0 for the first module in a ControlLogix rack).
func (pb *PathBuilder) Slot(slotNo uint8) *PathBuilder {
	if pb.Err != nil {
		return pb
	}
	err := binary.Write(&(pb.b), binary.LittleEndian, slotNo)
	if err != nil {
		pb.Err = fmt.Errorf("error writing address %v: %w", slotNo, err)
	}
	return pb
}

// Parse parses a human-readable path string in the same format a msg instruction uses in logix (e.g. "0,1,192.168.2.1,0,1")
// The resulting bytes are appended to the path being built.
func (pb *PathBuilder) Parse(path string) *PathBuilder {
	if pb.Err != nil {
		return pb
	}
	buf, err := ParsePath(path)
	if err != nil {
		pb.Err = fmt.Errorf("error parsing path %v: %w", path, err)
		return pb
	}
	_, err = pb.b.Write(buf.Bytes())
	if err != nil {
		pb.Err = fmt.Errorf("error writing parsed path %v: %w", path, err)
	}
	return pb
}
