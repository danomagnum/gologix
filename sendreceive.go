package gologix

import (
	"bytes"
	"context"
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

// prepares a bytes buffer for sending a message
func (client *Client) sendMsgBuild(cmd CIPCommand, msgs ...any) ([]byte, error) {
	// calculate messageLen of all message parts
	messageLen := 0
	for _, msg := range msgs {
		messageLen += binary.Size(msg)
	}
	// build header based on size
	hdr := client.newEIPHeader(cmd, messageLen)

	// initialize a buffer and add the header to it.
	// the 24 is from the header size
	b := make([]byte, 0, messageLen+24)
	buf := bytes.NewBuffer(b)
	err := binary.Write(buf, binary.LittleEndian, hdr)
	if err != nil {
		return nil, fmt.Errorf("problem writing header to buffer. %w", err)
	}

	// add all message components to the buffer.
	for _, msg := range msgs {
		err = binary.Write(buf, binary.LittleEndian, msg)
		if err != nil {
			return nil, fmt.Errorf("problem writing msg to buffer. %w", err)
		}
	}

	return buf.Bytes(), nil
}

// deadline returns the earlier of now+timeout and ctx's deadline (if any).
func deadline(ctx context.Context, timeout time.Duration) time.Time {
	d := time.Now().Add(timeout)
	if ctxDeadline, ok := ctx.Deadline(); ok && ctxDeadline.Before(d) {
		return ctxDeadline
	}
	return d
}

// send takes the command followed by all the structures that need
// concatenated together.
//
// It builds the appropriate header for all the data, puts the packet together, and then sends it.
func (client *Client) sendData(ctx context.Context, b []byte) error {
	// write the packet buffer to the tcp connection
	written := 0
	for written < len(b) {
		if err := ctx.Err(); err != nil {
			return err
		}
		timeout := client.SocketTimeout
		if timeout == 0 {
			timeout = time.Second
		}
		err := client.conn.SetWriteDeadline(deadline(ctx, timeout))
		if err != nil {
			return fmt.Errorf("problem setting write deadline: %w", err)
		}
		n, err := client.conn.Write(b[written:])
		if err != nil {
			return fmt.Errorf("problem writing to socket: %w", err)
		}
		written += n
	}
	return nil

}

// sends one message and gets one response in a mutex-protected way.
func (client *Client) send_recv_data(ctx context.Context, cmd CIPCommand, msgs ...any) (eipHeader, *bytes.Buffer, error) {
	if err := ctx.Err(); err != nil {
		return eipHeader{}, nil, err
	}
	buffer, err := client.sendMsgBuild(cmd, msgs...)
	if err != nil {
		return eipHeader{}, nil, fmt.Errorf("error preparing to send message: %w", err)
	}
	client.mutex.Lock()

	err = client.sendData(ctx, buffer)
	if err != nil {
		client.mutex.Unlock()
		err2 := client.Disconnect()
		if err2 != nil {
			return eipHeader{}, nil, fmt.Errorf("error disconnecting after send error %w: %w", err, err2)
		}
		return eipHeader{}, nil, fmt.Errorf("error sending data resulting in forced disconnect: %w", err)
	}

	hdr, buf, err := client.recvData(ctx)
	client.mutex.Unlock()
	if err != nil {
		err2 := client.Disconnect()
		if err2 != nil {
			return hdr, buf, fmt.Errorf("error disconnecting after recvError %w: %w", err, err2)
		}
		return hdr, buf, fmt.Errorf("error receiving data resulting in forced disconnect: %w", err)
	}
	return hdr, buf, nil
}

// recv_data reads the header and then the number of words it specifies.
func (client *Client) recvData(ctx context.Context) (eipHeader, *bytes.Buffer, error) {

	hdr := eipHeader{}
	if err := ctx.Err(); err != nil {
		return hdr, nil, err
	}
	timeout := client.SocketTimeout
	if timeout == 0 {
		timeout = time.Second
	}
	err := client.conn.SetReadDeadline(deadline(ctx, timeout))
	if err != nil {
		return hdr, nil, fmt.Errorf("problem setting read deadline: %w", err)
	}
	err = binary.Read(client.conn, binary.LittleEndian, &hdr)
	if err != nil {
		return hdr, nil, fmt.Errorf("problem reading header from socket: %w", err)
	}
	err = client.conn.SetReadDeadline(deadline(ctx, timeout))
	if err != nil {
		return hdr, nil, fmt.Errorf("problem setting read deadline: %w", err)
	}
	data_size := hdr.Length
	data := make([]byte, data_size)
	if data_size > 0 {

		err = binary.Read(client.conn, binary.LittleEndian, &data)
		if err != nil {
			return hdr, nil, fmt.Errorf("problem reading socket payload: %w", err)
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
	hdr.Length = uint16(size)
	hdr.SessionHandle = client.SessionHandle
	hdr.Status = 0
	hdr.Context = client.Context
	hdr.Options = 0

	return

}
