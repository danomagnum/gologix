package gologix

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

// send takes the command followed by all the structures that need
// concatenated together.
//
// It builds the appropriate header for all the data, puts the packet together, and then sends it.
func (client *Client) send(cmd CIPCommand, msgs ...any) error {
	// calculate size of all message parts
	size := 0
	for _, msg := range msgs {
		size += binary.Size(msg)
	}
	// build header based on size
	hdr := client.newEIPHeader(cmd, size)

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
		n, err := client.conn.Write(b[written:])
		if err != nil {
			return err
		}
		written += n
	}
	//log.Printf("Sent: %v", b)
	return nil

}

// sends one message and gets one response in a mutex-protected way.
func (client *Client) send_recv_data(cmd CIPCommand, msgs ...any) (EIPHeader, *bytes.Reader, error) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	err := client.send(cmd, msgs...)
	if err != nil {
		return EIPHeader{}, nil, err
	}
	return client.recv_data()

}

// recv_data reads the header and then the number of words it specifies.
func (client *Client) recv_data() (EIPHeader, *bytes.Reader, error) {

	hdr := EIPHeader{}
	var err error
	err = binary.Read(client.conn, binary.LittleEndian, &hdr)
	if err != nil {
		return hdr, nil, err
	}
	//log.Printf("Header: %v", hdr)
	//data_size := hdr.Length * 2
	data_size := hdr.Length
	data := make([]byte, data_size)
	if data_size > 0 {
		err = binary.Read(client.conn, binary.LittleEndian, &data)
	}
	//log.Printf("Buffer: %v", data)
	buf := bytes.NewReader(data)
	return hdr, buf, err

}

func (client *Client) newEIPHeader(cmd CIPCommand, size int) (hdr EIPHeader) {

	client.HeaderSequenceCounter++

	hdr.Command = uint16(cmd)
	//hdr.Command = 0x0070
	hdr.Length = uint16(size)
	hdr.SessionHandle = client.SessionHandle
	hdr.Status = 0
	hdr.Context = client.Context
	hdr.Options = 0

	return

}
