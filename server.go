package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/danomagnum/gologix/cipclass"
	"github.com/danomagnum/gologix/cipservice"
	"github.com/danomagnum/gologix/eipcommand"
)

// use NewServer() to get a new server object that is properly initialized.
//
// This is the main server object for handling incoming EIP messages.
// After setting up the server use the Serve() method to listen on the appropriate TCP and UDP ports
type Server struct {
	TCPListener net.Listener
	UDPListener net.PacketConn
	ConnMgr     serverConnectionManager
	Router      *PathRouter
	Attributes  map[cipclass.CIPAttribute]any
}

// an instance of serverTCPHandler will be created for every incomming connection to the EIP tcp port.
// it handles receiving messages, figuring out what kind they are, and using the associated server's ConnMgr and Router to
// dispatch the appropriate code
type serverTCPHandler struct {
	conn           net.Conn
	server         *Server
	handle         uint32
	options        uint32
	context        uint64
	OTConnectionID uint32
	TOConnectionID uint32

	UnitDataSequencer uint16
}

// Sets up a main server object for handling incoming EIP messages according to PathRouter r
// After setting up the server use the Serve() method to listen on the appropriate TCP and UDP ports
func NewServer(r *PathRouter) *Server {
	srv := Server{}
	srv.ConnMgr.Init()
	srv.Router = r
	srv.Attributes = make(map[cipclass.CIPAttribute]any)
	srv.Attributes[1] = int16(0x1776) // vendor ID
	srv.Attributes[2] = int16(0x000E) // device type
	srv.Attributes[3] = int16(0x0001) // Product Code
	srv.Attributes[4] = int16(0x1412) // Revision
	srv.Attributes[5] = int16(0x3060) // Status
	srv.Attributes[6] = rand.Uint32() // Serial
	srv.Attributes[7] = "Gologix"     // Product Name
	return &srv
}

// Start listening on the TCP and UDP ports associated with the Ethernet/IP protocol.
// these are 44818 and 2222 respectively
// as far as I can tell there is never an option to change this on any devices so it is hard coded here.
func (srv *Server) Serve() error {
	srv.ConnMgr.Init()

	var err error
	srv.TCPListener, err = net.Listen("tcp", "0.0.0.0:44818")
	log.Printf("Listening on TCP port 44818")
	if err != nil {
		return fmt.Errorf("couldn't open tcp listener. %w", err)
	}

	srv.UDPListener, err = net.ListenPacket("udp", "0.0.0.0:2222")
	log.Printf("Listening on UDP port 2222")
	if err != nil {
		return fmt.Errorf("couldn't open udp listener. %v", err)
	}

	// we'll start two server goroutines and then wait for either of them to error out on the error channel.

	errch := make(chan error)

	go func() {
		err := srv.serveUDP()
		if err != nil {
			errch <- fmt.Errorf("problem serving UDP. %w", err)
		}
	}()

	go func() {
		err := srv.serveTCP()
		if err != nil {
			errch <- fmt.Errorf("problem serving TCP. %w", err)
		}
	}()

	// we will wait forever for one of the serve goroutines to let us know they crased.
	// then we'll close them both and check for errors, combining them all toghether and returning it.
	err = <-errch
	final_err := NewMultiError(err)

	err = srv.TCPListener.Close()
	if err != nil {
		final_err.Add(fmt.Errorf("err on tcp close: %w", err))
	}
	err = srv.UDPListener.Close()
	if err != nil {
		final_err.Add(fmt.Errorf("err on udp close: %w", err))
	}

	return final_err

}

// this will listen on the EIP tcp port and kick off a serverTCPHandler for each connection
func (srv *Server) serveTCP() error {
	for {
		conn, err := srv.TCPListener.Accept()
		if err != nil {
			log.Printf("problem with tcp accept. %v", err)
			continue
		}
		// create a new handler and kick off its serve method to handle the connection
		h := serverTCPHandler{conn: conn, server: srv}
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

// this listens on the eip udp port and handles incoming messages.
// for each message that comes it it figures out which connection it belongs to and
// dispatches it accordingly to the proper router engpoint.
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
		if items[0].Header.ID == cipItem_SequenceAddress {
			// this is an IO message (output data from the controller to us as an "io adapter" in the hardware tree)
			io_info := cipIOSeqAccessData{}
			err := items[0].DeSerialize(&io_info)
			if err != nil {
				log.Printf("problem reading sequence address info.")
				continue
			}
			h, err := srv.ConnMgr.GetByOT(io_info.ConnectionID)
			if err != nil {
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

func (h *serverTCPHandler) serve(srv *Server) error {
	log.Printf("new connection from %v", h.conn.RemoteAddr().String())
	// if this function ends, close and remove ourselves from the server's open connection list.
	defer func() {
		h.conn.Close()
	}()

	for {
		var eiphdr EIPHeader
		err := binary.Read(h.conn, binary.LittleEndian, &eiphdr)
		if err != nil {
			return fmt.Errorf("problem reading eip header. %w", err)
		}
		h.context = eiphdr.Context
		log.Printf("context: %v\n", h.context)
		switch eiphdr.Command {
		case eipcommand.RegisterSession:
			err = h.registerSession(eiphdr)
			if err != nil {
				return fmt.Errorf("problem with register session %w", err)
			}
		case eipcommand.SendRRData:
			// this is things like forward opens
			err = h.sendRRData(eiphdr)
			if err != nil {
				return fmt.Errorf("problem with sendrrdata %w", err)
			}
		case eipcommand.SendUnitData:
			// this is things like writes and reads
			err = h.sendUnitData(eiphdr)
			if err != nil {
				return fmt.Errorf("problem with sendunitdata %w", err)
			}

		}

	}

}

func (h *serverTCPHandler) sendUnitData(hdr EIPHeader) error {
	var interface_handle uint32
	var timeout uint16
	err := binary.Read(h.conn, binary.LittleEndian, &interface_handle)
	if err != nil {
		return fmt.Errorf("problem reading interface handle %w", err)
	}
	err = binary.Read(h.conn, binary.LittleEndian, &timeout)
	if err != nil {
		return fmt.Errorf("problem reading timeout %w", err)
	}
	log.Printf("ih: %x. timeout: %x", interface_handle, timeout)
	items, err := ReadItems(h.conn)
	if err != nil {
		return fmt.Errorf("problem reading items for rrdata %w", err)
	}
	//log.Printf("items: %+v", items)
	if len(items) != 2 {
		return fmt.Errorf("expected 2 items. got %v", len(items))
	}
	// item 0 is the connected data item
	if items[0].Header.ID != cipItem_ConnectionAddress {
		return fmt.Errorf("should have had a connected data item in position 0. got %v", items[0].Header.ID)
	}
	var connid uint32
	err = items[0].DeSerialize(&connid)
	if err != nil {
		return fmt.Errorf("problem deserializing connection ID %w", err)
	}

	err = items[1].DeSerialize(&h.UnitDataSequencer)
	if err != nil {
		return fmt.Errorf("problem deserializing unit data seq %w", err)
	}
	var service cipservice.CIPService
	err = items[1].DeSerialize(&service)
	if err != nil {
		return fmt.Errorf("problem deserializing service %w", err)
	}
	switch service {
	case cipservice.Write:
		err = h.cipConnectedWrite(items)
		if err != nil {
			return fmt.Errorf("problem handling write. %w", err)
		}
	case cipservice.FragRead:
		err = h.connectedData(items)
		if err != nil {
			return fmt.Errorf("problem handling frag read. %w", err)
		}
	case cipservice.Read:
		err = h.connectedData(items)
		if err != nil {
			return fmt.Errorf("problem handling frag read. %w", err)
		}
	case cipservice.MultipleService:
		err = h.connectedData(items)
		if err != nil {
			return fmt.Errorf("problem handling multi service. %w", err)
		}
	case cipservice.GetAttributeSingle:
		err = h.connectedGetAttr(items)
		if err != nil {
			return fmt.Errorf("problem handling getAttrSingle %w", err)
		}
	default:
		log.Printf("Got unknown service at send unit data handler %d", service)
	}
	log.Printf("send unit data service requested: %v", service)
	return nil
}

func (h *serverTCPHandler) sendUnitDataReply(s cipservice.CIPService) error {
	items := make([]CIPItem, 2)
	items[0] = NewItem(cipItem_ConnectionAddress, h.TOConnectionID)
	items[1] = NewItem(cipItem_ConnectedData, nil)
	resp := msgWriteResultHeader{
		SequenceCount: h.UnitDataSequencer,
		Service:       s.AsResponse(),
	}
	items[1].Serialize(resp)
	return h.send(eipcommand.SendUnitData, SerializeItems(items))
}

func (h *serverTCPHandler) sendRRData(hdr EIPHeader) error {
	var interface_handle uint32
	var timeout uint16
	err := binary.Read(h.conn, binary.LittleEndian, &interface_handle)
	if err != nil {
		return fmt.Errorf("problem reading interface handle %w", err)
	}
	err = binary.Read(h.conn, binary.LittleEndian, &timeout)
	if err != nil {
		return fmt.Errorf("problem reading timeout %w", err)
	}
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
		var service cipservice.CIPService
		err = items[1].DeSerialize(&service)
		if err != nil {
			return fmt.Errorf("failed to get service. %w", err)
		}
		return h.connectedData(items)
	case cipItem_UnconnectedData:
		return h.unconnectedData(items[1])
	}
	return nil
}

func (h *serverTCPHandler) forwardClose(i CIPItem) error {
	log.Printf("got forward close from %v", h.conn.RemoteAddr())
	var fwd_close msgEIPForwardClose
	err := i.DeSerialize(&fwd_close)
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
	err = i.DeSerialize(&path)
	if err != nil {
		log.Printf("problem parsing forward close path. %v", err)
		return fmt.Errorf("problem parsing forward close path %w", err)
	}
	return nil
}

func (h *serverTCPHandler) largeforwardOpen(i CIPItem) error {
	log.Printf("got large forward open from %v", h.conn.RemoteAddr())
	var fwd_open msgEIPForwardOpen_Large
	err := i.DeSerialize(&fwd_open)
	if err != nil {
		return fmt.Errorf("problem with fwd open parsing %w", err)
	}
	fwd_path := make([]byte, fwd_open.ConnPathSize*2)
	err = i.DeSerialize(&fwd_path)
	if err != nil {
		return fmt.Errorf("problem with fwd open path parsing %w", err)
	}
	log.Printf("forward open msg: %v @ %v", fwd_open, fwd_path)
	path := fwd_path[:2]

	//preitem := msgPreItemData{Handle: 0, Timeout: 0}
	items := make([]CIPItem, 2)
	items[0] = CIPItem{Header: cipItemHeader{ID: cipItem_Null}}
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

	items[1].Serialize(fwopenresphdr)

	err = h.send(eipcommand.SendRRData, SerializeItems(items))
	if err != nil {
		return fmt.Errorf("problem sending response data %w", err)
	}

	cipConnection := &serverConnection{
		TO:   fwd_open.TOConnectionID,
		OT:   fwd_open.OTConnectionID,
		ID:   fwd_open.ConnectionSerialNumber,
		RPI:  time.Duration(fwd_open.TORPI) * time.Microsecond,
		Path: path,
		Open: true,
	}

	h.server.ConnMgr.Add(cipConnection)

	if fwd_open.TransportTrigger == 1 {
		log.Printf("Large Forward Open IO Connections Not Yet Supported")
	}

	return nil
}

func (h *serverTCPHandler) forwardOpen(i CIPItem) error {
	log.Printf("got small forward open from %v", h.conn.RemoteAddr())
	var fwd_open msgEIPForwardOpen_Standard
	err := i.DeSerialize(&fwd_open)
	if err != nil {
		return fmt.Errorf("problem with fwd open parsing %w", err)
	}
	fwd_path := make([]byte, fwd_open.ConnPathSize*2)
	err = i.DeSerialize(&fwd_path)
	if err != nil {
		return fmt.Errorf("problem with fwd open path parsing %w", err)
	}
	log.Printf("forward open msg: %v @ %v", fwd_open, fwd_path)
	path := fwd_path[:2]

	//preitem := msgPreItemData{Handle: 0, Timeout: 0}
	items := make([]CIPItem, 2)
	items[0] = CIPItem{Header: cipItemHeader{ID: cipItem_Null}}
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

	items[1].Serialize(fwopenresphdr)

	err = h.send(eipcommand.SendRRData, SerializeItems(items))
	if err != nil {
		return fmt.Errorf("problem sending response data %w", err)
	}

	cipConnection := &serverConnection{
		TO:   fwd_open.TOConnectionID,
		OT:   fwd_open.OTConnectionID,
		ID:   fwd_open.ConnectionSerialNumber,
		RPI:  time.Duration(fwd_open.TORPI) * time.Microsecond,
		Path: path,
		Open: true,
	}

	h.server.ConnMgr.Add(cipConnection)

	if fwd_open.TransportTrigger == 1 {
		// this is a cyclic IO connection
		tp, err := h.server.Router.Resolve(path)
		if err != nil {
			log.Printf("No tag provider for path %v", path)
			return fmt.Errorf("no tag provider for path %v", path)
		}
		go h.ioConnection(fwd_open, tp, cipConnection)
	}

	return nil
}

func (h *serverTCPHandler) ioConnection(fwd_open msgEIPForwardOpen_Standard, tp TagProvider, conn *serverConnection) {
	rpi := time.Duration(fwd_open.TORPI) * time.Microsecond
	log.Printf("IO RPI of %v", rpi)
	t := time.NewTicker(rpi)
	seq := uint32(0)

	// get the address to send the response back to, trim off the port number, and add the eip udp port number (2222) back on.
	// I suspect this will break with IPV6 since it uses colons in the IP address itself.
	remote := h.conn.RemoteAddr().String()
	remote = strings.Split(remote, ":")[0]
	addr := fmt.Sprintf("%s:2222", remote)

	udpconn, err := net.Dial("udp", addr)
	if err != nil {
		log.Printf("[ERROR] problem connecting UDP. %v", err)
		return
	}
	defer udpconn.Close()

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
		items := make([]CIPItem, 2)
		items[0] = NewItem(cipItem_SequenceAddress, nil)
		items[0].Serialize(fwd_open.TOConnectionID)
		items[0].Serialize(seq)
		items[1] = NewItem(cipItem_ConnectedData, nil)
		items[1].Serialize(uint16(seq))
		items[1].Serialize(dat)

		payload := *SerializeItems(items)
		payload = payload[6:]
		_, err = udpconn.Write(payload)
		if err != nil {
			log.Printf("problem writing %v", err)
		}

	}
}

func (h *serverTCPHandler) registerSession(hdr EIPHeader) error {

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

	err = h.send(eipcommand.RegisterSession, reg_msg)
	if err != nil {
		return fmt.Errorf("problem sending register response. %w", err)
	}

	return nil
}

// send takes the command followed by all the structures that need
// concatenated together.
//
// It builds the appropriate header for all the data, puts the packet together, and then sends it.
func (h *serverTCPHandler) send(cmd eipcommand.CIPCommand, msgs ...any) error {
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
	err := binary.Write(buf, binary.LittleEndian, hdr)
	if err != nil {
		return fmt.Errorf("problem writing header to buffer. %w", err)
	}

	// add all message components to the buffer.
	for _, msg := range msgs {
		err = binary.Write(buf, binary.LittleEndian, msg)
		if err != nil {
			return fmt.Errorf("problem writing message to buffer. %w", err)
		}
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
	return nil

}

func (h *serverTCPHandler) newEIPHeader(cmd eipcommand.CIPCommand, size int) (hdr EIPHeader) {

	hdr.Command = cmd
	//hdr.Command = 0x0070
	hdr.Length = uint16(size)
	hdr.SessionHandle = h.handle
	hdr.Status = 0
	hdr.Context = h.context
	hdr.Options = h.options

	return

}
