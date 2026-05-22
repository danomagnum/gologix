package gologix

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"
	"time"
)

// TestBuildListIdentityItemBody locks the wire bytes of the per-item body in
// a List Identity reply against the layout a real ControlLogix / CompactLogix
// emits. The big-endian fields (sin_family / sin_port / sin_addr) and the
// little-endian Identity attributes share one buffer, so the order and the
// endianness are both load-bearing.
func TestBuildListIdentityItemBody(t *testing.T) {
	srv := NewServer(nil)
	srv.Attributes[6] = uint32(0xC01E663C)
	srv.Attributes[7] = "gologix Server_Class3"

	ip := net.IPv4(192, 168, 1, 50)
	body, err := buildListIdentityItemBody(srv.Attributes, ip, 44818)
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	// Fixed prefix (EncapProtoVersion LE + sockaddr_in BE + sin_zero).
	wantPrefix := []byte{
		0x01, 0x00, // EncapProtoVersion = 1
		0x00, 0x02, // sin_family = AF_INET, big-endian
		0xAF, 0x12, // sin_port = 44818, big-endian
		0xC0, 0xA8, 0x01, 0x32, // sin_addr = 192.168.1.50
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // sin_zero
	}
	if !bytes.HasPrefix(body, wantPrefix) {
		t.Fatalf("prefix mismatch:\nwant % x\ngot  % x", wantPrefix, body[:len(wantPrefix)])
	}

	// Identity attributes follow the prefix in little-endian.
	r := bytes.NewReader(body[len(wantPrefix):])
	var vendor, deviceType, productCode, revision, status int16
	var serial uint32
	for i, target := range []any{&vendor, &deviceType, &productCode, &revision, &status, &serial} {
		if err := binary.Read(r, binary.LittleEndian, target); err != nil {
			t.Fatalf("read attr %d: %v", i, err)
		}
	}
	if vendor != srv.Attributes[1].(int16) {
		t.Errorf("VendorID = %d; want %d", vendor, srv.Attributes[1].(int16))
	}
	if deviceType != srv.Attributes[2].(int16) {
		t.Errorf("DeviceType = %d; want %d", deviceType, srv.Attributes[2].(int16))
	}
	if productCode != srv.Attributes[3].(int16) {
		t.Errorf("ProductCode = %d; want %d", productCode, srv.Attributes[3].(int16))
	}
	if revision != srv.Attributes[4].(int16) {
		t.Errorf("Revision = 0x%04X; want 0x%04X", revision, srv.Attributes[4].(int16))
	}
	if status != srv.Attributes[5].(int16) {
		t.Errorf("Status = 0x%04X; want 0x%04X", status, srv.Attributes[5].(int16))
	}
	if serial != uint32(0xC01E663C) {
		t.Errorf("SerialNumber = 0x%08X; want 0xC01E663C", serial)
	}

	// SHORT_STRING ProductName + 1 byte State follow.
	nameLen, err := r.ReadByte()
	if err != nil {
		t.Fatalf("read name length: %v", err)
	}
	want := srv.Attributes[7].(string)
	if int(nameLen) != len(want) {
		t.Fatalf("name length byte = %d; want %d", nameLen, len(want))
	}
	name := make([]byte, nameLen)
	if _, err := r.Read(name); err != nil {
		t.Fatalf("read name: %v", err)
	}
	if string(name) != want {
		t.Errorf("ProductName = %q; want %q", string(name), want)
	}
	state, err := r.ReadByte()
	if err != nil {
		t.Fatalf("read state: %v", err)
	}
	if state != 0xFF {
		t.Errorf("State = 0x%02X; want 0xFF", state)
	}
	if rem := r.Len(); rem != 0 {
		t.Errorf("trailing bytes left: %d", rem)
	}
}

// TestFrameListIdentityResponse asserts the CPF wrapper around the identity
// item body matches what a List Identity reply carries on the wire: 1 item
// count, item type 0x000C, item length, then the body.
func TestFrameListIdentityResponse(t *testing.T) {
	body := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	framed := frameListIdentityResponse(body)
	want := append(
		[]byte{
			0x01, 0x00, // item count = 1
			0x0C, 0x00, // item type = 0x000C
			0x04, 0x00, // body length = 4
		},
		body...,
	)
	if !bytes.Equal(framed, want) {
		t.Fatalf("framed mismatch:\nwant % x\ngot  % x", want, framed)
	}
}

// TestSendListIdentityReplyEndToEnd boots a real gologix server in-process
// and drives a List Identity request through the TCP listener exactly like
// pycomm3 does. The test fails when port 44818 is busy so it stays out of
// the way of other hardware tests sharing the same bind.
func TestSendListIdentityReplyEndToEnd(t *testing.T) {
	if probe, err := net.Listen("tcp", "0.0.0.0:44818"); err != nil {
		t.Skipf("port 44818 unavailable: %v", err)
	} else {
		probe.Close()
	}

	srv := NewServer(&PathRouter{})
	srv.Attributes[6] = uint32(0xDEADBEEF)
	srv.Attributes[7] = "gologix-listidentity-test"
	go func() { _ = srv.Serve() }()
	defer func() {
		if srv.TCPListener != nil {
			srv.TCPListener.Close()
		}
		if srv.UDPListener != nil {
			srv.UDPListener.Close()
		}
	}()

	deadline := time.Now().Add(2 * time.Second)
	var conn net.Conn
	var err error
	for time.Now().Before(deadline) {
		conn, err = net.DialTimeout("tcp", "127.0.0.1:44818", 100*time.Millisecond)
		if err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatalf("set deadline: %v", err)
	}

	reqHdr := eipHeader{
		Command: cipCommandListIdentity,
		Length:  0,
	}
	if err := binary.Write(conn, binary.LittleEndian, reqHdr); err != nil {
		t.Fatalf("write request header: %v", err)
	}

	var respHdr eipHeader
	if err := binary.Read(conn, binary.LittleEndian, &respHdr); err != nil {
		t.Fatalf("read response header: %v", err)
	}
	if respHdr.Command != cipCommandListIdentity {
		t.Fatalf("response command = 0x%X; want 0x%X", respHdr.Command, cipCommandListIdentity)
	}
	if respHdr.Length == 0 {
		t.Fatalf("response carried no body")
	}

	payload := make([]byte, respHdr.Length)
	if _, err := conn.Read(payload); err != nil {
		t.Fatalf("read payload: %v", err)
	}

	// Item count + item header (4 bytes) + body.
	if len(payload) < 6 {
		t.Fatalf("payload too small: %d bytes", len(payload))
	}
	count := binary.LittleEndian.Uint16(payload[0:2])
	itemType := binary.LittleEndian.Uint16(payload[2:4])
	itemLen := binary.LittleEndian.Uint16(payload[4:6])
	if count != 1 {
		t.Errorf("item count = %d; want 1", count)
	}
	if CIPItemID(itemType) != cipIdentityItemType {
		t.Errorf("item type = 0x%04X; want 0x%04X", itemType, cipIdentityItemType)
	}
	body := payload[6:]
	if len(body) != int(itemLen) {
		t.Errorf("body length = %d; item header said %d", len(body), itemLen)
	}

	// EncapProtoVersion (LE) + sockaddr_in (BE) prefix.
	if got := binary.LittleEndian.Uint16(body[0:2]); got != 1 {
		t.Errorf("EncapProtoVersion = %d; want 1", got)
	}
	if got := binary.BigEndian.Uint16(body[2:4]); got != 2 {
		t.Errorf("sin_family = %d; want 2 (AF_INET, big-endian)", got)
	}
	if got := binary.BigEndian.Uint16(body[4:6]); got != 44818 {
		t.Errorf("sin_port = %d; want 44818 (big-endian)", got)
	}
}
