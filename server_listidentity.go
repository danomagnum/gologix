package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

// cipIdentityItemType is the CPF item ID a List Identity reply carries its
// payload under. The encoding is documented in CIP Vol 2, Section 2-4.1 and
// matches what a ControlLogix/CompactLogix emits when probed.
const cipIdentityItemType CIPItemID = 0x000C

// buildListIdentityItemBody renders the per-item payload an EtherNet/IP
// List Identity reply carries: a fixed-layout slice combining one
// little-endian Encapsulation Protocol Version field, a big-endian
// sockaddr_in (the controller's bound interface), the 7 attributes of the
// Identity Object, and a one-byte state field. The big-endian portion is
// non-negotiable — every native Logix client (FactoryTalk Linx, pycomm3,
// RSLogix browse) reads sin_family / sin_port / sin_addr that way.
//
// Inputs come from the same Server.Attributes map that
// unconnectedGetAttributeAll reads, so both Identity paths stay in sync.
func buildListIdentityItemBody(attrs map[CIPAttribute]any, localIP net.IP, localPort uint16) ([]byte, error) {
	vendor, ok := attrs[1].(int16)
	if !ok {
		return nil, fmt.Errorf("attribute 1 (VendorID) missing or wrong type: %T", attrs[1])
	}
	deviceType, ok := attrs[2].(int16)
	if !ok {
		return nil, fmt.Errorf("attribute 2 (DeviceType) missing or wrong type: %T", attrs[2])
	}
	productCode, ok := attrs[3].(int16)
	if !ok {
		return nil, fmt.Errorf("attribute 3 (ProductCode) missing or wrong type: %T", attrs[3])
	}
	revision, ok := attrs[4].(int16)
	if !ok {
		return nil, fmt.Errorf("attribute 4 (Revision) missing or wrong type: %T", attrs[4])
	}
	status, ok := attrs[5].(int16)
	if !ok {
		return nil, fmt.Errorf("attribute 5 (Status) missing or wrong type: %T", attrs[5])
	}
	serial, ok := attrs[6].(uint32)
	if !ok {
		return nil, fmt.Errorf("attribute 6 (SerialNumber) missing or wrong type: %T", attrs[6])
	}
	productName, ok := attrs[7].(string)
	if !ok {
		return nil, fmt.Errorf("attribute 7 (ProductName) missing or wrong type: %T", attrs[7])
	}
	if len(productName) > 255 {
		productName = productName[:255]
	}

	var addr [4]byte
	if ip4 := localIP.To4(); ip4 != nil {
		copy(addr[:], ip4)
	}

	var buf bytes.Buffer
	// EncapProtocolVersion — little-endian uint16.
	_ = binary.Write(&buf, binary.LittleEndian, uint16(1))
	// SocketAddress — sockaddr_in in network (big-endian) byte order.
	_ = binary.Write(&buf, binary.BigEndian, uint16(2)) // AF_INET
	_ = binary.Write(&buf, binary.BigEndian, localPort)
	buf.Write(addr[:])
	// sin_zero — 8 bytes of padding required by the sockaddr_in layout.
	buf.Write(make([]byte, 8))
	// Identity attributes — little-endian.
	_ = binary.Write(&buf, binary.LittleEndian, vendor)
	_ = binary.Write(&buf, binary.LittleEndian, deviceType)
	_ = binary.Write(&buf, binary.LittleEndian, productCode)
	_ = binary.Write(&buf, binary.LittleEndian, revision)
	_ = binary.Write(&buf, binary.LittleEndian, status)
	_ = binary.Write(&buf, binary.LittleEndian, serial)
	buf.WriteByte(byte(len(productName)))
	buf.WriteString(productName)
	// State — one byte, 0xFF when the device offers no state attribute.
	buf.WriteByte(0xFF)
	return buf.Bytes(), nil
}

// frameListIdentityResponse wraps the per-item identity body in the CPF
// structure an EtherNet/IP List Identity reply mandates: 1 item count
// followed by an item header (type 0x000C plus the body length) and the
// body itself.
func frameListIdentityResponse(itemBody []byte) []byte {
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.LittleEndian, uint16(1))                  // Item count.
	_ = binary.Write(&buf, binary.LittleEndian, cipIdentityItemType)        // Item type.
	_ = binary.Write(&buf, binary.LittleEndian, uint16(len(itemBody)))      // Item length.
	buf.Write(itemBody)
	return buf.Bytes()
}

// sendListIdentityReply handles an EtherNet/IP List Identity request that
// arrives on the TCP listener. The handler reads the local TCP address to
// fill the sockaddr portion of the reply and emits the framed identity
// payload via the standard send helper.
func (h *serverTCPHandler) sendListIdentityReply(hdr eipHeader) error {
	ip, port := extractLocalSocket(h.conn.LocalAddr())
	body, err := buildListIdentityItemBody(h.server.Attributes, ip, port)
	if err != nil {
		return fmt.Errorf("build list identity body: %w", err)
	}
	return h.send(cipCommandListIdentity, frameListIdentityResponse(body))
}

// extractLocalSocket pulls the IPv4 address and port from a net.Addr coming
// out of a TCP or UDP listener. It falls back to 0.0.0.0:44818 when the
// address does not surface a port (most often during in-process tests using
// net.Pipe-style endpoints).
func extractLocalSocket(addr net.Addr) (net.IP, uint16) {
	switch a := addr.(type) {
	case *net.TCPAddr:
		return a.IP, uint16(a.Port)
	case *net.UDPAddr:
		return a.IP, uint16(a.Port)
	}
	return net.IPv4zero, 44818
}
