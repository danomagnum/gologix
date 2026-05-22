package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"strings"
	"time"
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
	Attributes  map[CIPAttribute]any
	Logger      *slog.Logger
}

// an instance of serverTCPHandler will be created for every incoming connection to the EIP tcp port.
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
	srv.Logger = slog.Default()
	srv.ConnMgr.Init(srv.Logger)
	srv.Router = r
	srv.Attributes = make(map[CIPAttribute]any)
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
	srv.ConnMgr.Init(srv.Logger)

	var err error
	srv.TCPListener, err = net.Listen("tcp", "0.0.0.0:44818")
	srv.Logger.Info("Listening on TCP port 44818")
	if err != nil {
		return fmt.Errorf("couldn't open tcp listener. %w", err)
	}

	srv.UDPListener, err = net.ListenPacket("udp", "0.0.0.0:2222")
	srv.Logger.Info("Listening on UDP port 2222")
	if err != nil {
		return fmt.Errorf("couldn't open udp listener. %v", err)
	}

	// we'll start two server goroutines and then wait for either of them to error out on the error channel.

	errCh := make(chan error)

	go func() {
		err := srv.serveUDP()
		if err != nil {
			errCh <- fmt.Errorf("problem serving UDP. %w", err)
		}
	}()

	go func() {
		err := srv.serveTCP()
		if err != nil {
			errCh <- fmt.Errorf("problem serving TCP. %w", err)
		}
	}()

	// we will wait forever for one of the serve goroutines to let us know they crashed.
	// then we'll close them both and check for errors, combining them all together and returning it.
	err = <-errCh
	final_err := newMultiError(err)

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
			srv.Logger.Error("problem with tcp accept", "error", err)
			continue
		}
		// create a new handler and kick off its serve method to handle the connection
		h := serverTCPHandler{conn: conn, server: srv}
		go func() {
			err := h.serve(srv)
			if err != nil {
				srv.Logger.Warn("Error on connection %v. %v", h.conn.RemoteAddr().String(), err)
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
// dispatches it accordingly to the proper router endpoint.
func (srv *Server) serveUDP() error {
	bufSize := 4096
	for {
		b := make([]byte, 4096)
		buf := bytes.NewBuffer(b)
		n, addr, err := srv.UDPListener.ReadFrom(b)
		if n == 0 {
			srv.Logger.Debug("Read 0 bytes on udp listener.")
			continue
		}
		if n == bufSize {
			srv.Logger.Debug("udp buffer size not big enough!")
			continue
		}
		if err != nil {
			srv.Logger.Debug("problem with udp accept", "error", err)
			continue
		}
		_ = addr // don't need this yet.
		// we've read a packet on udp so we need to parse the eip data

		items, err := readItems(buf)
		if err != nil {
			srv.Logger.Debug("problem reading udp items", "error", err)
			continue
		}
		if len(items) != 2 {
			srv.Logger.Debug("expected 2 items", "got", len(items))
			continue
		}
		if items[0].Header.ID == cipItem_SequenceAddress {
			// this is an IO message (output data from the controller to us as an "io adapter" in the hardware tree)
			io_info := cipIOSeqAccessData{}
			err := items[0].DeSerialize(&io_info)
			if err != nil {
				srv.Logger.Debug("problem reading sequence address info.")
				continue
			}
			h, err := srv.ConnMgr.GetByOT(io_info.ConnectionID)
			if err != nil {
				srv.Logger.Debug("couldn't find handler to handle IO message", "connection", io_info.ConnectionID)
				continue
			}

			tp, err := srv.Router.Resolve(h.Path)
			if err != nil {
				srv.Logger.Debug("couldn't find tag provider to handle IO message", "connection", io_info.ConnectionID, "path", h.Path)
				continue
			}
			err = tp.IOWrite(items)
			if err != nil {
				srv.Logger.Debug("Problem writing IO to tag provider", "error", err)
				continue
			}

		}

	}
}

func (h *serverTCPHandler) serve(srv *Server) error {
	srv.Logger.Info("new connection", "remote addr", h.conn.RemoteAddr().String())
	// if this function ends, close and remove ourselves from the server's open connection list.
	defer func() {
		h.conn.Close()
	}()

	for {
		var eipHdr eipHeader
		err := binary.Read(h.conn, binary.LittleEndian, &eipHdr)
		if err != nil {
			return fmt.Errorf("problem reading eip header. %w", err)
		}
		h.context = eipHdr.Context
		h.server.Logger.Debug("New Context", "context", h.context)
		switch eipHdr.Command {
		case cipCommandRegisterSession:
			err = h.registerSession(eipHdr)
			if err != nil {
				return fmt.Errorf("problem with register session %w", err)
			}
		case cipCommandSendRRData:
			// this is things like forward opens
			err = h.sendRRData(eipHdr)
			if err != nil {
				return fmt.Errorf("problem with sendRRData %w", err)
			}
		case cipCommandSendUnitData:
			// this is things like writes and reads
			err = h.sendUnitData(eipHdr)
			if err != nil {
				return fmt.Errorf("problem with sendUnitData %w", err)
			}
		case cipCommandListServices:
			err = h.sendListServicesData(eipHdr)
			if err != nil {
				return fmt.Errorf("problem with sendListServices %w", err)
			}
		case cipCommandListIdentity:
			err = h.sendListIdentityReply(eipHdr)
			if err != nil {
				return fmt.Errorf("problem with sendListIdentity %w", err)
			}

		}
	}

}

func (h *serverTCPHandler) sendUnitData(hdr eipHeader) error {
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
	items, err := readItems(h.conn)
	if err != nil {
		return fmt.Errorf("problem reading items for rrData %w", err)
	}
	if len(items) != 2 {
		return fmt.Errorf("expected 2 items. got %v", len(items))
	}
	// item 0 is the connected data item
	if items[0].Header.ID != cipItem_ConnectionAddress {
		return fmt.Errorf("should have had a connected data item in position 0. got %v", items[0].Header.ID)
	}
	var connId uint32
	err = items[0].DeSerialize(&connId)
	if err != nil {
		return fmt.Errorf("problem deserializing connection ID %w", err)
	}

	err = items[1].DeSerialize(&h.UnitDataSequencer)
	if err != nil {
		return fmt.Errorf("problem deserializing unit data seq %w", err)
	}
	var service CIPService
	err = items[1].DeSerialize(&service)
	if err != nil {
		return fmt.Errorf("problem deserializing service %w", err)
	}
	switch service {
	case CIPService_Write:
		err = h.cipConnectedWrite(items)
		if err != nil {
			return fmt.Errorf("problem handling write. %w", err)
		}
	case CIPService_FragRead:
		err = h.connectedRead(CIPService_FragRead, items)
		if err != nil {
			return fmt.Errorf("problem handling frag read. %w", err)
		}
	case CIPService_Read:
		err = h.connectedRead(CIPService_Read, items)
		if err != nil {
			return fmt.Errorf("problem handling non-frag read. %w", err)
		}
	case CIPService_MultipleService:
		err = h.connectedMulti(items)
		if err != nil {
			return fmt.Errorf("problem handling multi service. %w", err)
		}
	case CIPService_GetAttributeSingle:
		err = h.connectedGetAttr(items)
		if err != nil {
			return fmt.Errorf("problem handling getAttrSingle %w", err)
		}
	case CIPService_GetAttributeAll:
		// pycomm3.LogixDriver.get_plc_name() and similar clients probe
		// Class 0x64 (Program Object) Instance 1 over the established
		// connection after Forward Open. Previously this fell through to
		// the silent default which left the client hanging until timeout.
		err = h.connectedGetAttributeAll(items)
		if err != nil {
			return fmt.Errorf("problem handling connected getAttributesAll %w", err)
		}
	default:
		// Any other service that lands here would have silently dropped
		// before, leaving strict CIP clients (FactoryTalk Linx, Studio
		// 5000, pycomm3) waiting forever for a reply. Send a formal CIP
		// error response so the client sees "alive but service not
		// supported" and can recover gracefully.
		h.server.Logger.Warn("connected: unsupported service", "service", service)
		return h.sendUnitDataErrorReply(service, CIPStatus_ServiceNotSupported)
	}
	h.server.Logger.Debug("send unit data service requested", "service", service)
	return nil
}

// sendUnitDataErrorReply emits a CIP error response on an established
// connection (SendUnitData transport) so the client sees a formal status
// code instead of silence. Mirrors sendUnitDataReply but populates the
// non-zero CIPStatus required for error semantics.
func (h *serverTCPHandler) sendUnitDataErrorReply(s CIPService, status CIPStatus) error {
	items := make([]CIPItem, 2)
	items[0] = newItem(cipItem_ConnectionAddress, h.TOConnectionID)
	items[1] = newItem(cipItem_ConnectedData, nil)
	resp := msgWriteResultHeader{
		SequenceCount: h.UnitDataSequencer,
		Service:       s.AsResponse(),
		Status:        status,
	}
	if err := items[1].Serialize(resp); err != nil {
		return fmt.Errorf("serialize connected error header: %w", err)
	}
	itemData, err := serializeItems(items)
	if err != nil {
		return fmt.Errorf("serialize connected error items: %w", err)
	}
	return h.send(cipCommandSendUnitData, itemData)
}

// connectedGetAttributeAll handles CIP Get_Attributes_All (service 0x01)
// arriving on an already established Class 3 connection. The Logix
// Program Object (Class 0x64) Instance 1 path is the one pycomm3 probes
// after Forward Open to learn the controller name; routing that through
// the same Identity attribute map we already expose keeps a single
// source of truth for the device's product name string.
func (h *serverTCPHandler) connectedGetAttributeAll(items []CIPItem) error {
	items[1].Reset()
	item := items[1]

	if _, err := item.Uint16(); err != nil {
		return fmt.Errorf("connected GAA: read seq: %w", err)
	}
	if _, err := item.Byte(); err != nil {
		return fmt.Errorf("connected GAA: read service byte: %w", err)
	}
	pathSize, err := item.Byte()
	if err != nil {
		return fmt.Errorf("connected GAA: read path size: %w", err)
	}
	pathBytes := make([]byte, int(pathSize)*2)
	if err := item.DeSerialize(&pathBytes); err != nil {
		return fmt.Errorf("connected GAA: read path: %w", err)
	}

	class, instance, err := parseClassInstancePath(bytes.NewBuffer(pathBytes))
	if err != nil {
		h.server.Logger.Warn("connected GAA: malformed path", "bytes", pathBytes, "error", err)
		return h.sendUnitDataErrorReply(CIPService_GetAttributeAll, CIPStatus_PathSegmentError)
	}

	switch class {
	case CipObject_Identity:
		if instance != 1 {
			return h.sendUnitDataErrorReply(CIPService_GetAttributeAll, CIPStatus_PathDestinationUnknown)
		}
		payload, err := buildIdentityGetAttributesAllResponse(h.server.Attributes)
		if err != nil {
			return fmt.Errorf("connected GAA identity: %w", err)
		}
		return h.sendUnitDataReplyWithPayload(CIPService_GetAttributeAll, payload)

	case 0x64:
		if instance != 1 {
			return h.sendUnitDataErrorReply(CIPService_GetAttributeAll, CIPStatus_PathDestinationUnknown)
		}
		payload, err := buildProgramObjectGetAttributesAllResponse(h.server.Attributes)
		if err != nil {
			return fmt.Errorf("connected GAA program: %w", err)
		}
		return h.sendUnitDataReplyWithPayload(CIPService_GetAttributeAll, payload)

	default:
		h.server.Logger.Warn("connected GAA: class not implemented", "class", class, "instance", instance)
		return h.sendUnitDataErrorReply(CIPService_GetAttributeAll, CIPStatus_PathDestinationUnknown)
	}
}

// sendUnitDataReplyWithPayload is sendUnitDataReply with an additional
// raw payload appended after the standard write-result header. Used by
// the connected Get_Attributes_All handler to ship the Identity /
// Program Object attribute bytes back to the client.
func (h *serverTCPHandler) sendUnitDataReplyWithPayload(s CIPService, payload []byte) error {
	items := make([]CIPItem, 2)
	items[0] = newItem(cipItem_ConnectionAddress, h.TOConnectionID)
	items[1] = newItem(cipItem_ConnectedData, nil)
	resp := msgWriteResultHeader{
		SequenceCount: h.UnitDataSequencer,
		Service:       s.AsResponse(),
	}
	if err := items[1].Serialize(resp); err != nil {
		return fmt.Errorf("serialize connected reply header: %w", err)
	}
	if err := items[1].Serialize(payload); err != nil {
		return fmt.Errorf("serialize connected reply payload: %w", err)
	}
	itemData, err := serializeItems(items)
	if err != nil {
		return fmt.Errorf("serialize connected reply items: %w", err)
	}
	return h.send(cipCommandSendUnitData, itemData)
}

func (h *serverTCPHandler) sendUnitDataReply(s CIPService) error {
	items := make([]CIPItem, 2)
	items[0] = newItem(cipItem_ConnectionAddress, h.TOConnectionID)
	items[1] = newItem(cipItem_ConnectedData, nil)
	resp := msgWriteResultHeader{
		SequenceCount: h.UnitDataSequencer,
		Service:       s.AsResponse(),
	}
	err := items[1].Serialize(resp)
	if err != nil {
		return fmt.Errorf("problem serializing write result header %w", err)
	}
	itemData, err := serializeItems(items)
	if err != nil {
		return fmt.Errorf("could not serialize items: %w", err)
	}
	return h.send(cipCommandSendUnitData, itemData)
}

func (h *serverTCPHandler) sendRRData(hdr eipHeader) error {
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
	items, err := readItems(h.conn)
	if err != nil {
		return fmt.Errorf("problem reading items for rrData %w", err)
	}
	if len(items) != 2 {
		return fmt.Errorf("expected 2 items. got %v", len(items))
	}
	switch items[1].Header.ID {
	case cipItem_ConnectedData:
		var service CIPService
		err = items[1].DeSerialize(&service)
		if err != nil {
			return fmt.Errorf("failed to get service. %w", err)
		}
		return h.connectedRead(service, items)
	case cipItem_UnconnectedData:
		return h.unconnectedData(items[1])
	}
	return nil
}

func (h *serverTCPHandler) forwardClose(i CIPItem) error {
	h.server.Logger.Info("got forward close", "remote", h.conn.RemoteAddr())
	var fwd_close msgEIPForwardClose // TODO: Transition this to the new ForwardOpen
	err := i.DeSerialize(&fwd_close)
	if err != nil {
		return fmt.Errorf("problem parsing forward close %w", err)
	}
	h.server.Logger.Info("Closing connection", "serial", fwd_close.ConnectionSerialNumber)
	err = h.server.ConnMgr.CloseByID(fwd_close.ConnectionSerialNumber)
	if err != nil {
		return fmt.Errorf("couldn't close open connection with ID %v. %v", fwd_close.ConnectionSerialNumber, err)
	}

	path := make([]byte, fwd_close.PathSize*2)
	err = i.DeSerialize(&path)
	if err != nil {
		return fmt.Errorf("problem parsing forward close path %w", err)
	}
	return nil
}

// in this message T is for target and O is for originator so
// TO is target -> originator and OT is originator -> target
type msgEIPForwardOpen_Standard struct {
	Service                CIPService
	PathSize               byte
	ClassType              byte
	Class                  byte
	InstanceType           byte
	Instance               byte
	Priority               byte
	TimeoutTicks           byte
	OTConnectionID         uint32
	TOConnectionID         uint32
	ConnectionSerialNumber uint16
	VendorID               uint16
	OriginatorSerialNumber uint32
	Multiplier             uint32
	OTRPI                  uint32
	OTNetworkConnParams    uint16
	TORPI                  uint32
	TONetworkConnParams    uint16
	TransportTrigger       byte
	ConnPathSize           byte
}

type msgEIPForwardOpen_Large struct {
	// service
	Service CIPService
	// path
	PathSize     byte
	ClassType    cipClassSize
	Class        byte
	InstanceType cipInstanceSize
	Instance     byte

	// service specific data
	Priority               byte
	TimeoutTicks           byte
	OTConnectionID         uint32
	TOConnectionID         uint32
	ConnectionSerialNumber uint16
	VendorID               uint16
	OriginatorSerialNumber uint32
	Multiplier             uint32
	OTRPI                  uint32
	OTNetworkConnParams    uint32
	TORPI                  uint32
	TONetworkConnParams    uint32
	TransportTrigger       byte
	ConnPathSize           byte
}

func (h *serverTCPHandler) largeForwardOpen(i CIPItem) error {
	h.server.Logger.Info("got large forward open", "remote", h.conn.RemoteAddr())
	var fwd_open msgEIPForwardOpen_Large // TODO: Transition this to the new ForwardOpen
	err := i.DeSerialize(&fwd_open)
	if err != nil {
		return fmt.Errorf("problem with fwd open parsing %w", err)
	}
	fwd_path := make([]byte, fwd_open.ConnPathSize*2)
	err = i.DeSerialize(&fwd_path)
	if err != nil {
		return fmt.Errorf("problem with fwd open path parsing %w", err)
	}
	h.server.Logger.Debug("forward open msg", "open", fwd_open, "path", fwd_path)
	path := fwd_path[:2]

	//preItem := msgPreItemData{Handle: 0, Timeout: 0}
	items := make([]CIPItem, 2)
	items[0] = CIPItem{Header: cipItemHeader{ID: cipItem_Null}}
	items[1] = newItem(cipItem_UnconnectedData, nil)

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
	// Echo the requested RPI back as the Actual Packet Interval (API). CIP
	// spec ODVA Vol 1 §3-5.6 allows the target to commit to a different API
	// than what was requested, but it MUST report a non-zero value the
	// originator can honor. Real Logix controllers always echo the requested
	// RPI verbatim for Class 3 explicit messaging connections.
	//
	// Previously this was hardcoded to 0, which caused real ControlLogix
	// MSG instructions to silently discard responses on connected reads (the
	// firmware-level scheduler treats API=0 as invalid timing and tears
	// down the response path even though the connection itself stays open).
	fwOpenRespHdr := msgEIPForwardOpen_Standard_Reply{
		Service:                fwd_open.Service.AsResponse(),
		OTConnectionID:         h.OTConnectionID,
		TOConnectionID:         h.TOConnectionID,
		ConnectionSerialNumber: fwd_open.ConnectionSerialNumber,
		VendorID:               fwd_open.VendorID,
		OriginatorSerialNumber: fwd_open.OriginatorSerialNumber,
		OTApi:                  fwd_open.OTRPI,
		TOApi:                  fwd_open.TORPI,
	}

	err = items[1].Serialize(fwOpenRespHdr)
	if err != nil {
		return fmt.Errorf("problem serializing fwd open response %w", err)
	}

	itemData, err := serializeItems(items)
	if err != nil {
		return fmt.Errorf("could not serialize items: %w", err)
	}
	err = h.send(cipCommandSendRRData, itemData)
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
		h.server.Logger.Warn("Large Forward Open IO Connections Not Yet Supported")
	}

	return nil
}

func (h *serverTCPHandler) forwardOpen(i CIPItem) error {
	h.server.Logger.Info("got small forward open", "remote", h.conn.RemoteAddr())
	var fwd_open msgEIPForwardOpen_Standard // TODO: Transition this to the new forward open
	err := i.DeSerialize(&fwd_open)
	if err != nil {
		return fmt.Errorf("problem with fwd open parsing %w", err)
	}
	fwd_path := make([]byte, fwd_open.ConnPathSize*2)
	err = i.DeSerialize(&fwd_path)
	if err != nil {
		return fmt.Errorf("problem with fwd open path parsing %w", err)
	}

	// if this is a produce/consume, the fwd_path here will have a few segments in it.
	// First is a "keying" segment.  It's 4 words of 0's in the wireshark samples I've got.
	// Then there is a 8-bit class segment 0x69 which I've named "IO Class" until i know better.
	// Then there is intance 1
	// Then there is an ansii class 0x91 of the consumed tag name.
	// then a "connection point" segment of 1
	// then a "simple data segment" of some mystery data.

	h.server.Logger.Debug("forward open msg", "open", fwd_open, "path", fwd_path)
	path := fwd_path[:2]

	//preItem := msgPreItemData{Handle: 0, Timeout: 0}
	items := make([]CIPItem, 2)
	items[0] = CIPItem{Header: cipItemHeader{ID: cipItem_Null}}
	items[1] = newItem(cipItem_UnconnectedData, nil)

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
	// See largeForwardOpen above for why API must echo the requested RPI
	// instead of zero. Same fix applies to the standard (small) Forward
	// Open path because real Logix MSG instructions hit this code path for
	// any read/write that fits in the standard connection size envelope.
	fwOpenRespHdr := msgEIPForwardOpen_Standard_Reply{
		Service:                fwd_open.Service.AsResponse(),
		OTConnectionID:         h.OTConnectionID,
		TOConnectionID:         h.TOConnectionID,
		ConnectionSerialNumber: fwd_open.ConnectionSerialNumber,
		VendorID:               fwd_open.VendorID,
		OriginatorSerialNumber: fwd_open.OriginatorSerialNumber,
		OTApi:                  fwd_open.OTRPI,
		TOApi:                  fwd_open.TORPI,
	}

	err = items[1].Serialize(fwOpenRespHdr)
	if err != nil {
		return fmt.Errorf("problem serializing forward open response %w", err)
	}

	itemData, err := serializeItems(items)
	if err != nil {
		return fmt.Errorf("could not serialize items: %w", err)
	}
	err = h.send(cipCommandSendRRData, itemData)
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
			return fmt.Errorf("no tag provider for path %v", path)
		}
		go h.ioConnection(fwd_open, tp, cipConnection)
	}

	return nil
}

func (h *serverTCPHandler) ioConnection(fwd_open msgEIPForwardOpen_Standard, tp CIPEndpoint, conn *serverConnection) {
	rpi := time.Duration(fwd_open.TORPI) * time.Microsecond
	h.server.Logger.Info("IO Connection", "rpi", rpi)
	t := time.NewTicker(rpi)
	seq := uint32(0)

	// get the address to send the response back to, trim off the port number, and add the eip udp port number (2222) back on.
	// I suspect this will break with IPV6 since it uses colons in the IP address itself.
	remote := h.conn.RemoteAddr().String()
	remote = strings.Split(remote, ":")[0]
	addr := fmt.Sprintf("%s:2222", remote)

	udpConn, err := net.Dial("udp", addr)
	if err != nil {
		h.server.Logger.Error("Problem connecting UDP", "error", err)
		return
	}
	defer udpConn.Close()

	for {
		seq++
		<-t.C
		if !conn.Open {
			h.server.Logger.Warn("connection closed. no longer sending IO messages", "conn", *conn)
			return

		}

		dat, err := tp.IORead()
		if err != nil {
			h.server.Logger.Warn("problem getting IO data from provider", "error", err)
			continue
		}

		// every RPI send the message.
		items := make([]CIPItem, 2)
		items[0] = newItem(cipItem_SequenceAddress, nil)
		err = items[0].Serialize(fwd_open.TOConnectionID, seq)
		if err != nil {
			h.server.Logger.Warn("problem serializing connection ID", "error", err)
			return
		}

		items[1] = newItem(cipItem_ConnectedData, nil)
		err = items[1].Serialize(uint16(seq), dat)
		if err != nil {
			h.server.Logger.Warn("problem serializing seq number", "error", err)
			return
		}

		p, err := serializeItems(items)
		if err != nil {
			h.server.Logger.Warn("could not serialize items", "error", err)
			return
		}
		payload := *p
		payload = payload[6:]
		_, err = udpConn.Write(payload)
		if err != nil {
			h.server.Logger.Warn("problem writing", "error", err)
		}

	}
}

func (h *serverTCPHandler) registerSession(hdr eipHeader) error {

	reg_msg := msgCIPRegister{}
	err := binary.Read(h.conn, binary.LittleEndian, &reg_msg)
	if err != nil {
		return fmt.Errorf("problem reading register session message. %w", err)
	}
	h.server.Logger.Debug("register message", "msg", reg_msg)
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
func (h *serverTCPHandler) send(cmd CIPCommand, msgs ...any) error {
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

func (h *serverTCPHandler) newEIPHeader(cmd CIPCommand, size int) (hdr eipHeader) {

	hdr.Command = cmd
	//hdr.Command = 0x0070
	hdr.Length = uint16(size)
	hdr.SessionHandle = h.handle
	hdr.Status = 0
	hdr.Context = h.context
	hdr.Options = h.options

	return

}

type ServiceCapabilityFlags uint16

const (
	ServiceCapabilityFlag_CipEncapsulation  ServiceCapabilityFlags = 1 << 5
	ServiceCapabilityFlag_SupportsClass1UDP ServiceCapabilityFlags = 1 << 8
)

type listServicesReply struct {
	Count    uint16
	TypeCode uint16
	Length   uint16
	Version  uint16
	CapFlags ServiceCapabilityFlags
	Name     [16]byte
}

func (h *serverTCPHandler) sendListServicesData(hdr eipHeader) error {
	// on a list services request there is no more data to read.
	response := listServicesReply{
		Count:    1,
		TypeCode: 0x0100,
		Length:   10,
		Version:  1,
		CapFlags: ServiceCapabilityFlag_CipEncapsulation | ServiceCapabilityFlag_SupportsClass1UDP,
		Name:     [16]byte{'C', 'o', 'm', 'm', 'u', 'n', 'i', 'c', 'a', 't', 'i', 'o', 'n', 's', ' ', ' '},
	}
	return h.send(cipCommandSendRRData, response)
}
