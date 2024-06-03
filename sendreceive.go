package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

type eipHeader struct {
	Command       CIPCommand
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
	err := binary.Write(buf, binary.LittleEndian, hdr)
	if err != nil {
		return fmt.Errorf("problem writing header to buffer. %w", err)
	}

	// add all message components to the buffer.
	for _, msg := range msgs {
		err = binary.Write(buf, binary.LittleEndian, msg)
		if err != nil {
			return fmt.Errorf("problem writing msg to buffer. %w", err)
		}
	}

	b = buf.Bytes()
	// write the packet buffer to the tcp connection
	written := 0
	for written < len(b) {
		if client.SocketTimeout != 0 {
			client.conn.SetWriteDeadline(time.Now().Add(client.SocketTimeout))
		} else {
			client.conn.SetWriteDeadline(time.Now().Add(time.Second))
		}
		n, err := client.conn.Write(b[written:])
		if err != nil {
			err2 := client.Disconnect()
			return fmt.Errorf("%w: %v", err, err2)
		}
		written += n
	}
	return nil

}

// sends one message and gets one response in a mutex-protected way.
func (client *Client) send_recv_data(cmd CIPCommand, msgs ...any) (eipHeader, *bytes.Buffer, error) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	err := client.send(cmd, msgs...)
	if err != nil {
		return eipHeader{}, nil, err
	}
	return client.recv_data()

}

// recv_data reads the header and then the number of words it specifies.
func (client *Client) recv_data() (eipHeader, *bytes.Buffer, error) {

	hdr := eipHeader{}
	var err error
	if client.SocketTimeout != 0 {
		client.conn.SetReadDeadline(time.Now().Add(client.SocketTimeout))
	} else {
		client.conn.SetReadDeadline(time.Now().Add(time.Second))
	}
	err = binary.Read(client.conn, binary.LittleEndian, &hdr)
	if err != nil {
		err2 := client.Disconnect()
		return hdr, nil, fmt.Errorf("problem reading header from socket: %w: %v", err, err2)
	}
	if client.SocketTimeout != 0 {
		client.conn.SetReadDeadline(time.Now().Add(client.SocketTimeout))
	} else {
		client.conn.SetReadDeadline(time.Now().Add(time.Second))
	}
	data_size := hdr.Length
	data := make([]byte, data_size)
	if data_size > 0 {

		err = binary.Read(client.conn, binary.LittleEndian, &data)
		if err != nil {
			err2 := client.Disconnect()
			return hdr, nil, fmt.Errorf("problem reading socket payload: %w: %v", err, err2)
		}
	}
	buf := bytes.NewBuffer(data)
	return hdr, buf, err

}

func (client *Client) DebugCloseConn() {
	client.conn.Close()
}

func (client *Client) newEIPHeader(cmd CIPCommand, size int) (hdr eipHeader) {

	client.HeaderSequenceCounter++

	hdr.Command = cmd
	//hdr.Command = 0x0070
	hdr.Length = uint16(size)
	hdr.SessionHandle = client.SessionHandle
	hdr.Status = 0
	hdr.Context = client.Context
	hdr.Options = 0

	return

}
