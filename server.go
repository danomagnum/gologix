package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
)

type Server struct {
	Listener    net.Listener
	Connections map[net.Addr]net.Conn
	ConnMutex   sync.Mutex
}

type handler struct {
	conn    net.Conn
	handle  uint32
	options uint32
	context uint64
}

func NewServer() *Server {
	s := Server{}
	s.Connections = make(map[net.Addr]net.Conn)
	return &s
}

func (srv *Server) Serve() error {

	var err error
	log.Printf("Listening on port 44818")
	srv.Listener, err = net.Listen("tcp", "0.0.0.0:44818")
	if err != nil {
		return fmt.Errorf("couldn't open listener. %v", err)
	}

	for {
		conn, err := srv.Listener.Accept()
		if err != nil {
			log.Printf("problem with accept. %v", err)
			continue
		}
		srv.ConnMutex.Lock()
		srv.Connections[conn.RemoteAddr()] = conn
		srv.ConnMutex.Unlock()
		h := handler{conn: conn}
		go func() {
			err := h.serve(srv)
			if err != nil {
				log.Printf("Error on connnection %v. %v", h.conn.RemoteAddr().String(), err)
			}
		}()
	}

}

func (h *handler) serve(srv *Server) error {
	log.Printf("new connection from %v", h.conn.RemoteAddr().String())
	// if this function ends, close and remove ourselves from the server's open connection list.
	defer func() {
		h.conn.Close()
		srv.ConnMutex.Lock()
		delete(srv.Connections, h.conn.RemoteAddr())
		srv.ConnMutex.Unlock()
	}()

	for {
		var eiphdr EIPHeader
		err := binary.Read(h.conn, binary.LittleEndian, &eiphdr)
		if err != nil {
			return fmt.Errorf("problem reading eip header. %w", err)
		}
		switch eiphdr.Command {
		case cipCommandRegisterSession:
			err = h.registerSession(eiphdr)
			if err != nil {
				return fmt.Errorf("problem with register session %w", err)
			}
		case cipCommandSendRRData:
			err = h.sendRRData(eiphdr)
			if err != nil {
				return fmt.Errorf("problem with sendrrdata %w", err)
			}

		}

	}

}
func (h *handler) sendRRData(hdr EIPHeader) error {
	var interface_handle uint32
	var timeout uint16
	binary.Read(h.conn, binary.LittleEndian, interface_handle)
	binary.Read(h.conn, binary.LittleEndian, timeout)
	log.Printf("ih: %x. timeout: %x", interface_handle, timeout)
	items, err := ReadItems(h.conn)
	if err != nil {
		return fmt.Errorf("problem reading items for rrdata %w", err)
	}
	log.Printf("items: %+v", items)
	if len(items) != 2 {
		return fmt.Errorf("expected 2 items. got %v", len(items))
	}
	var service CIPService
	items[1].Unmarshal(&service)
	items[1].Reset()
	switch service {
	case cipService_ForwardOpen:
		err = h.forwardOpen(items[1])
		if err != nil {
			return fmt.Errorf("problem handling forward open. %w", err)
		}
	default:
		log.Printf("Got unkown service %d", service)
	}
	log.Printf("service requested: %v", service)
	return nil
}

func (h *handler) forwardOpen(i cipItem) error {
	log.Printf("got small forward open from %v", h.conn.RemoteAddr())
	var fwd_open msgEIPForwardOpen_Standard
	err := i.Unmarshal(&fwd_open)
	if err != nil {
		return fmt.Errorf("problem with fwd open parsing %w", err)
	}
	log.Printf("forward open msg: %v", fwd_open)

	preitem := msgPreItemData{}
	items := make([]cipItem, 2)
	items[0] = cipItem{Header: cipItemHeader{ID: cipItem_Null}}
	reply := msgEIPForwardOpen_Reply{
		OTConnectionID: rand.Uint32(),
		TOConnectionID: rand.Uint32(),
	}
	h.send(cipCommandSendRRData, preitem, reply)

	return nil
}

func (h *handler) registerSession(hdr EIPHeader) error {

	reg_msg := msgCIPRegister{}
	err := binary.Read(h.conn, binary.LittleEndian, &reg_msg)
	if err != nil {
		return fmt.Errorf("problem reading register session message. %w", err)
	}
	log.Printf("register message: %+v", reg_msg)
	h.handle = hdr.SessionHandle
	if h.handle == 0 {
		h.handle = rand.Uint32()
	}
	h.options = hdr.Options

	err = h.send(cipCommandRegisterSession, reg_msg)
	if err != nil {
		return fmt.Errorf("problem sending register response. %w", err)
	}

	return nil
}

// send takes the command followed by all the structures that need
// concatenated together.
//
// It builds the appropriate header for all the data, puts the packet together, and then sends it.
func (h *handler) send(cmd CIPCommand, msgs ...any) error {
	// calculate size of all message parts
	size := 0
	for _, msg := range msgs {
		size += binary.Size(msg)
	}
	// build header based on size
	hdr := h.newEIPHeader(cmd, size)

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
		n, err := h.conn.Write(b[written:])
		if err != nil {
			return err
		}
		written += n
	}
	//log.Printf("Sent: %v", b)
	return nil

}

func (h *handler) newEIPHeader(cmd CIPCommand, size int) (hdr EIPHeader) {

	hdr.Command = cmd
	//hdr.Command = 0x0070
	hdr.Length = uint16(size)
	hdr.SessionHandle = h.handle
	hdr.Status = 0
	hdr.Context = h.context
	hdr.Options = h.options

	return

}
