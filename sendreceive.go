package main

import (
	"bytes"
	"encoding/binary"
)

type EIPHeader struct {
	Command       uint16
	Length        uint16
	SessionHandle uint32
	Status        uint32
	Context       uint64 // 8 bytes you can do whatever you want with. They'll be echoed back.
	Options       uint32
}

// Send takes the command followed by all the structures that need
// concatenated together.
//
// It builds the appropriate header for all the data, puts the packet together, and then sends it.
func (plc *PLC) Send(cmd CIPCommand, msgs ...any) error {
	// calculate size of all message parts
	size := 0
	for _, msg := range msgs {
		size += binary.Size(msg)
	}
	// build header based on size
	hdr := plc.NewEIPHeader(cmd, size)

	// initialize a buffer and add the header to it.
	// the 24 is from the header size
	b := make([]byte, 0, size+24)
	buf := bytes.NewBuffer(b)
	binary.Write(buf, binary.LittleEndian, hdr)

	// add all message components to the buffer.
	for _, msg := range msgs {
		binary.Write(buf, binary.LittleEndian, msg)
	}

	b = buf.Bytes()
	// write the packet buffer to the tcp connection
	written := 0
	for written < len(b) {
		n, err := plc.Conn.Write(b[written:])
		if err != nil {
			return err
		}
		written += n
	}
	//log.Printf("Sent: %v", b)
	return nil

}

// recv_data reads the header and then the number of words it specifies.
func (plc *PLC) recv_data() (EIPHeader, *bytes.Reader, error) {

	hdr := EIPHeader{}
	var err error
	err = binary.Read(plc.Conn, binary.LittleEndian, &hdr)
	if err != nil {
		return hdr, nil, err
	}
	//log.Printf("Header: %v", hdr)
	//data_size := hdr.Length * 2
	data_size := hdr.Length
	data := make([]byte, data_size)
	if data_size > 0 {
		err = binary.Read(plc.Conn, binary.LittleEndian, &data)
	}
	//log.Printf("Buffer: %v", data)
	buf := bytes.NewReader(data)
	return hdr, buf, err

}

func (plc *PLC) NewEIPHeader(cmd CIPCommand, size int) (hdr EIPHeader) {

	plc.SequenceCounter++

	hdr.Command = uint16(cmd)
	//hdr.Command = 0x0070
	hdr.Length = uint16(size)
	hdr.SessionHandle = plc.SessionHandle
	hdr.Status = 0
	hdr.Context = plc.Context
	hdr.Options = 0

	return

}
