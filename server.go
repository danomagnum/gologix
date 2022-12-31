package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

type serverConnection struct {
	ID   uint16
	OT   uint32
	TO   uint32
	RPI  time.Duration
	Path []byte
	Open bool
}

type serverConnectionManager struct {
	Connections []*serverConnection
	Lock        sync.RWMutex
}

func (cm *serverConnectionManager) Init() {
	cm.Connections = make([]*serverConnection, 0, 32)
}

func (cm *serverConnectionManager) Add(conn *serverConnection) {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	cm.Connections = append(cm.Connections, conn)
}

func (cm *serverConnectionManager) GetByID(ID uint16) (*serverConnection, error) {
	cm.Lock.RLock()
	defer cm.Lock.RUnlock()
	for _, conn := range cm.Connections {
		if conn.ID == ID {
			return conn, nil
		}
	}
	return nil, fmt.Errorf("couldn't find connection %v by ID", ID)
}

func (cm *serverConnectionManager) GetByOT(OT uint32) (*serverConnection, error) {
	cm.Lock.RLock()
	defer cm.Lock.RUnlock()
	for _, conn := range cm.Connections {
		if conn.OT == OT {
			return conn, nil
		}
	}
	return nil, fmt.Errorf("couldn't find connection %v by OT", OT)
}
func (cm *serverConnectionManager) GetByTO(TO uint32) (*serverConnection, error) {
	cm.Lock.RLock()
	defer cm.Lock.RUnlock()
	for _, conn := range cm.Connections {
		if conn.TO == TO {
			return conn, nil
		}
	}
	return nil, fmt.Errorf("couldn't find connection %v by TO", TO)
}

func (cm *serverConnectionManager) CloseByID(ID uint16) error {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	for i, conn := range cm.Connections {
		if conn.ID == ID {
			conn.Open = false
			if len(cm.Connections) == 1 {
				cm.Connections = make([]*serverConnection, 0, 32)
				return nil
			}
			cm.Connections[i] = cm.Connections[len(cm.Connections)-1]
			cm.Connections = cm.Connections[:len(cm.Connections)-1]
			return nil
		}
	}
	return fmt.Errorf("couldn't find connection %v by ID", ID)
}

func (cm *serverConnectionManager) CloseByOT(OT uint32) error {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	for i, conn := range cm.Connections {
		if conn.OT == OT {
			conn.Open = false
			if len(cm.Connections) == 1 {
				cm.Connections = make([]*serverConnection, 0, 32)
				return nil
			}
			cm.Connections[i] = cm.Connections[len(cm.Connections)-1]
			cm.Connections = cm.Connections[:len(cm.Connections)-1]
			return nil
		}
	}
	return fmt.Errorf("couldn't find connection %v by OT", OT)
}
func (cm *serverConnectionManager) CloseByTO(TO uint32) error {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	for i, conn := range cm.Connections {
		if conn.TO == TO {
			conn.Open = false
			if len(cm.Connections) == 1 {
				cm.Connections = make([]*serverConnection, 0, 32)
				return nil
			}
			cm.Connections[i] = cm.Connections[len(cm.Connections)-1]
			cm.Connections = cm.Connections[:len(cm.Connections)-1]
			return nil
		}
	}
	return fmt.Errorf("couldn't find connection %v by TO", TO)
}

type Server struct {
	TCPListener   net.Listener
	UDPListener   net.PacketConn
	OTConnections map[uint32]*handler
	ConnMgr       serverConnectionManager
	ConnMutex     sync.Mutex
	Router        *PathRouter
}

type handler struct {
	conn           net.Conn
	server         *Server
	handle         uint32
	options        uint32
	context        uint64
	OTConnectionID uint32
	TOConnectionID uint32
	Path           []byte

	UnitDataSequencer uint16
}

func NewServer(r *PathRouter) *Server {
	s := Server{}
	s.OTConnections = make(map[uint32]*handler)
	s.Router = r
	return &s
}

func (srv *Server) Serve() error {
	srv.ConnMgr.Init()

	var err error
	srv.TCPListener, err = net.Listen("tcp", "0.0.0.0:44818")
	log.Printf("Listening on TCP port 44818")
	if err != nil {
		return fmt.Errorf("couldn't open tcp listener. %v", err)
	}

	srv.UDPListener, err = net.ListenPacket("udp", "0.0.0.0:2222")
	log.Printf("Listening on UDP port 2222")
	if err != nil {
		return fmt.Errorf("couldn't open udp listener. %v", err)
	}
	go srv.serveUDP()
	return srv.serveTCP()

}

func (srv *Server) serveTCP() error {
	for {
		conn, err := srv.TCPListener.Accept()
		if err != nil {
			log.Printf("problem with tcp accept. %v", err)
			continue
		}
		h := handler{conn: conn, server: srv}
		go func() {
			err := h.serve(srv)
			if err != nil {
				log.Printf("Error on connnection %v. %v", h.conn.RemoteAddr().String(), err)
			}
		}()
	}
}

type cipIOSeqAccessData struct {
	ConnectionID  uint32
	SequenceCount uint32
}

func (srv *Server) serveUDP() error {
	bufsize := 4096
	for {
		b := make([]byte, 4096)
		buf := bytes.NewBuffer(b)
		n, addr, err := srv.UDPListener.ReadFrom(b)
		if n == 0 {
			log.Print("Read 0 bytes on udp listener.")
			continue
		}
		if n == bufsize {
			log.Print("udp buffer size not big enough!")
			continue
		}
		if err != nil {
			log.Printf("problem with udp accept. %v", err)
			continue
		}
		_ = addr // don't need this yet.
		// we've read a packet on udp so we need to parse the eip data

		items, err := ReadItems(buf)
		if err != nil {
			log.Printf("problem reading udp items. %v", err)
			continue
		}
		if len(items) != 2 {
			log.Printf("expected 2 items but got %v", len(items))
			continue
		}
		log.Printf("got %v on udp items.", items)
		if items[0].Header.ID == cipItem_SequenceAddress {
			// this is an IO message (output data from the controller to us as an "io adapter")
			io_info := cipIOSeqAccessData{}
			err := items[0].Unmarshal(&io_info)
			if err != nil {
				log.Printf("problem reading sequence address info.")
				continue
			}
			h, ok := srv.OTConnections[io_info.ConnectionID]
			if !ok {
				log.Printf("couldn't find handler for connection %v to handle IO message", io_info.ConnectionID)
				continue
			}

			tp, err := srv.Router.Resolve(h.Path)
			if err != nil {
				log.Printf("couldn't find tag provider for connection %v at path %v to handle IO message", io_info.ConnectionID, h.Path)
				continue
			}
			err = tp.IOWrite(items)
			if err != nil {
				log.Printf("Problem writing IO to tag provider. %v", err)
				continue
			}

		}

	}
}

func (h *handler) serve(srv *Server) error {
	log.Printf("new connection from %v", h.conn.RemoteAddr().String())
	// if this function ends, close and remove ourselves from the server's open connection list.
	defer func() {
		h.conn.Close()
		srv.ConnMutex.Lock()
		delete(srv.OTConnections, h.OTConnectionID)
		srv.ConnMutex.Unlock()
	}()

	for {
		var eiphdr EIPHeader
		err := binary.Read(h.conn, binary.LittleEndian, &eiphdr)
		if err != nil {
			return fmt.Errorf("problem reading eip header. %w", err)
		}
		h.context = eiphdr.Context
		fmt.Printf("context: %v\n", h.context)
		switch eiphdr.Command {
		case cipCommandRegisterSession:
			err = h.registerSession(eiphdr)
			if err != nil {
				return fmt.Errorf("problem with register session %w", err)
			}
		case cipCommandSendRRData:
			// this is things like forward opens
			err = h.sendRRData(eiphdr)
			if err != nil {
				return fmt.Errorf("problem with sendrrdata %w", err)
			}
		case cipCommandSendUnitData:
			// this is things like writes and reads
			err = h.sendUnitData(eiphdr)
			if err != nil {
				return fmt.Errorf("problem with sendunitdata %w", err)
			}

		}

	}

}

func (h *handler) sendUnitData(hdr EIPHeader) error {
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
	// item 0 is the connected data item
	if items[0].Header.ID != cipItem_ConnectionAddress {
		return fmt.Errorf("should have had a connected data item in position 0. got %v", items[0].Header.ID)
	}
	var connid uint32
	items[0].Unmarshal(&connid)

	items[1].Unmarshal(&h.UnitDataSequencer)
	var service CIPService
	items[1].Unmarshal(&service)
	switch service {
	case cipService_Write:
		err = h.cipConnectedWrite(items)
		if err != nil {
			return fmt.Errorf("problem handling write. %w", err)
		}
	default:
		log.Printf("Got unknown service %d", service)
	}
	log.Printf("send unit data service requested: %v", service)
	return nil
}

func (h *handler) cipFragRead(item *cipItem) error {
	if item.Header.ID != cipItem_UnconnectedData {
		return fmt.Errorf("expected unconnected frag read. got %v", item.Header.ID)
	}
	fmt.Printf("frag read data: %v", item.Data)
	//return h.sendUnitDataReply(cipService_FragRead)
	return h.sendUnconnectedUnitDataReply(cipService_FragRead)

}

func (h *handler) cipConnectedWrite(items []cipItem) error {
	var l byte // length in words
	item := items[1]
	item.Unmarshal(l)
	var path_type SegmentType
	item.Unmarshal(&path_type)
	if path_type != SegmentTypeExtendedSymbolic {
		return fmt.Errorf("only support symbolic writes. got segment type %v", path_type)
	}
	var tag_length byte
	item.Unmarshal(&tag_length)
	tag_bytes := make([]byte, tag_length)
	item.Unmarshal(&tag_bytes)
	tag := string(tag_bytes)

	// string will be padded with a null if odd length
	if (tag_length % 2) == 1 {
		var b byte
		item.Unmarshal(&b)
	}

	var typ CIPType
	item.Unmarshal(&typ)
	var reserved byte
	item.Unmarshal(&reserved)
	var elements uint16
	item.Unmarshal(&elements)

	fmt.Printf("tag: %s", tag)
	for i := 0; i < int(elements); i++ {
		v := typ.readValue(&item)
		fmt.Printf("value: %v", v)
	}

	// path is part of the forward open we've previously received.

	return h.sendUnitDataReply(cipService_Write)
}

func (h *handler) sendUnitDataReply(s CIPService) error {
	items := make([]cipItem, 2)
	items[0] = NewItem(cipItem_ConnectionAddress, h.TOConnectionID)
	items[1] = NewItem(cipItem_ConnectedData, nil)
	resp := msgWriteResultHeader{
		SequenceCount: h.UnitDataSequencer,
		Service:       s.AsResponse(),
	}
	items[1].Marshal(resp)
	return h.send(cipCommandSendUnitData, MarshalItems(items))
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
	switch items[1].Header.ID {
	case cipItem_ConnectedData:
		return h.connectedData(items[1])
	case cipItem_UnconnectedData:
		return h.unconnectedData(items[1])
	}
	return nil
}
func (h *handler) connectedData(item cipItem) error {
	var service CIPService
	var err error
	item.Unmarshal(&service)
	item.Reset()
	switch service {
	case cipService_FragRead:
		err = h.cipFragRead(&item)
		if err != nil {
			return fmt.Errorf("problem handling frag read. %w", err)
		}
	default:
		log.Printf("Got unknown service %d", service)
	}
	log.Printf("sendrrdata service requested: %v", service)
	return nil
}

func getTagFromPath(item *cipItem) (string, error) {
	var prefix byte
	err := item.Unmarshal(&prefix)
	if err != nil {
		return "", fmt.Errorf("problem getting path prefix. %w", err)
	}
	if prefix != 0x91 {
		return "", fmt.Errorf("only support reading by tag name. TODO: support other things?. %w", err)
	}
	var tag_len byte
	err = item.Unmarshal(&tag_len)
	if err != nil {
		return "", fmt.Errorf("problem getting tag len. %w", err)
	}
	b := make([]byte, tag_len)
	err = item.Unmarshal(b)
	if err != nil {
		return "", fmt.Errorf("problem reading tag path. %w", err)
	}
	if tag_len%2 == 1 {
		var pad byte
		err = item.Unmarshal(&pad)
		if err != nil {
			return "", fmt.Errorf("problem reading pad byte. %w", err)
		}
	}
	return string(b), nil

}

func (h *handler) forwardClose(i cipItem) error {
	log.Printf("got forward close from %v", h.conn.RemoteAddr())
	var fwd_close msgEIPForwardClose
	err := i.Unmarshal(&fwd_close)
	if err != nil {
		log.Printf("problem parsing forward close. %v", err)
		return fmt.Errorf("problem parsing forward close %w", err)
	}
	log.Printf("Closing connection %v", fwd_close.ConnectionSerialNumber)
	err = h.server.ConnMgr.CloseByID(fwd_close.ConnectionSerialNumber)
	if err != nil {
		log.Printf("couldn't close open connection with ID %v. %v", fwd_close.ConnectionSerialNumber, err)
		return fmt.Errorf("couldn't close open connection with ID %v. %v", fwd_close.ConnectionSerialNumber, err)
	}

	path := make([]byte, fwd_close.PathSize*2)
	err = i.Unmarshal(&path)
	if err != nil {
		log.Printf("problem parsing forward close path. %v", err)
		return fmt.Errorf("problem parsing forward close path %w", err)
	}
	return nil
}

func (h *handler) forwardOpen(i cipItem) error {
	log.Printf("got small forward open from %v", h.conn.RemoteAddr())
	var fwd_open msgEIPForwardOpen_Standard
	err := i.Unmarshal(&fwd_open)
	if err != nil {
		return fmt.Errorf("problem with fwd open parsing %w", err)
	}
	fwd_path := make([]byte, fwd_open.ConnPathSize*2)
	err = i.Unmarshal(&fwd_path)
	if err != nil {
		return fmt.Errorf("problem with fwd open path parsing %w", err)
	}
	log.Printf("forward open msg: %v @ %v", fwd_open, fwd_path)
	h.Path = fwd_path[:2]

	//preitem := msgPreItemData{Handle: 0, Timeout: 0}
	items := make([]cipItem, 2)
	items[0] = cipItem{Header: cipItemHeader{ID: cipItem_Null}}
	items[1] = NewItem(cipItem_UnconnectedData, nil)

	if fwd_open.TOConnectionID == 0 {
		h.TOConnectionID = rand.Uint32()
	} else {
		h.TOConnectionID = fwd_open.TOConnectionID
	}
	fwd_open.TOConnectionID = h.TOConnectionID
	if fwd_open.OTConnectionID == 0 {
		h.OTConnectionID = rand.Uint32()
	} else {
		h.OTConnectionID = fwd_open.OTConnectionID
	}
	fwd_open.OTConnectionID = h.OTConnectionID
	fwopenresphdr := msgEIPForwardOpen_Standard_Reply{
		Service:                fwd_open.Service.AsResponse(),
		OTConnectionID:         h.OTConnectionID,
		TOConnectionID:         h.TOConnectionID,
		ConnectionSerialNumber: fwd_open.ConnectionSerialNumber,
		VendorID:               fwd_open.VendorID,
		OriginatorSerialNumber: fwd_open.OriginatorSerialNumber,
		OTAPI:                  0,
		TOAPI:                  0,
	}

	items[1].Marshal(fwopenresphdr)

	h.send(cipCommandSendRRData, MarshalItems(items))

	cipConnection := &serverConnection{
		TO:   fwd_open.TOConnectionID,
		OT:   fwd_open.OTConnectionID,
		ID:   fwd_open.ConnectionSerialNumber,
		RPI:  time.Duration(fwd_open.TORPI) * time.Microsecond,
		Open: true,
	}

	h.server.ConnMgr.Add(cipConnection)

	// set us up as the connection handler for this connection ID
	h.server.ConnMutex.Lock()
	h.server.OTConnections[h.OTConnectionID] = h
	h.server.ConnMutex.Unlock()

	if fwd_open.TransportTrigger == 1 {
		// this is a cyclic IO connection
		tp, err := h.server.Router.Resolve(h.Path[:2])
		if err != nil {
			log.Printf("No tag provider for path %v", h.Path[:2])
			return fmt.Errorf("no tag provider for path %v", h.Path[:2])
		}
		go h.ioConnection(fwd_open, tp, cipConnection)
	}

	return nil
}

func (h *handler) ioConnection(fwd_open msgEIPForwardOpen_Standard, tp TagProvider, conn *serverConnection) {
	rpi := time.Duration(fwd_open.TORPI) * time.Microsecond
	log.Printf("IO RPI of %v", rpi)
	t := time.NewTicker(rpi)
	seq := uint32(0)
	for {
		seq++
		<-t.C
		if !conn.Open {
			log.Printf("connection %+v closed. no longer sending IO messages", *conn)
			return

		}

		dat, err := tp.IORead()
		if err != nil {
			log.Printf("problem getting IO data from provider %v", err)
			continue
		}

		// every RPI send the message.
		items := make([]cipItem, 2)
		items[0] = NewItem(cipItem_SequenceAddress, nil)
		items[0].Marshal(fwd_open.TOConnectionID)
		items[0].Marshal(seq)
		items[1] = NewItem(cipItem_ConnectedData, nil)
		items[1].Marshal(uint16(seq))
		items[1].Marshal(dat)

		// get the address to send the response back to, trim off the port number, and add the eip udp port number (2222) back on.
		remote := h.conn.RemoteAddr().String()
		remote = strings.Split(remote, ":")[0]
		addr := fmt.Sprintf("%s:2222", remote)

		conn, err := net.Dial("udp", addr)
		if err != nil {
			log.Printf("problem connecting UDP. %v", err)
			continue
		}

		payload := *MarshalItems(items)
		payload = payload[6:]
		log.Printf("writing udp io payload %v", payload)
		_, err = conn.Write(payload)
		if err != nil {
			log.Printf("problem writing %v", err)
		}
		conn.Close()

	}
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
