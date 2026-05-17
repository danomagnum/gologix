package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"net"
	"testing"
	"time"
)

// fakeCIPServer is a net.Pipe-based test fake that lets tests script byte-level
// CIP responses, including non-OK general statuses like 0x06 PartialTransfer
// that the real server-side code cannot currently emit on demand.
type fakeCIPServer struct {
	t         *testing.T
	serverEnd net.Conn
	requests  chan fakeRequest
	done      chan struct{}
}

// fakeRequest is a parsed view of a single CIP request the client sent over the wire.
type fakeRequest struct {
	rawPayload []byte // everything after the 24-byte EIP header
	command    CIPCommand
	seq        uint16
	service    CIPService
	pathBytes  []byte
	elements   uint16
	// offset is populated only for CIPService_FragRead requests; zero otherwise.
	offset uint32
}

// newFakeCIPClient returns a Client wired to a fakeCIPServer over net.Pipe.
// The client is pre-configured as already connected — no real handshake runs.
func newFakeCIPClient(t *testing.T) (*Client, *fakeCIPServer) {
	t.Helper()
	clientEnd, serverEnd := net.Pipe()

	fs := &fakeCIPServer{
		t:         t,
		serverEnd: serverEnd,
		requests:  make(chan fakeRequest, 8),
		done:      make(chan struct{}),
	}
	go fs.readLoop()

	client := NewClient("fake")
	client.conn = clientEnd
	client.connStatus = connectionStatusConnected
	client.AutoConnect = false
	client.SessionHandle = 0xDEADBEEF
	client.OTNetworkConnectionID = 0x00112233
	client.SocketTimeout = 500 * time.Millisecond
	client.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	client.ioi_cache = make(map[string]*tagIOI)
	// Skip the firmware probe (GetAttributeSingle on Identity object) that newIOI
	// triggers on first call. Tests provide synthetic responses for reads only.
	client.knownFirmware = 99

	t.Cleanup(func() {
		_ = clientEnd.Close()
		fs.close()
	})

	return client, fs
}

// readLoop pulls EIP-framed requests off the pipe, parses what it can, and pushes
// onto fs.requests for the test to consume via awaitRequest.
func (fs *fakeCIPServer) readLoop() {
	defer close(fs.done)
	for {
		var hdr eipHeader
		if err := binary.Read(fs.serverEnd, binary.LittleEndian, &hdr); err != nil {
			return
		}
		payload := make([]byte, hdr.Length)
		if _, err := io.ReadFull(fs.serverEnd, payload); err != nil {
			return
		}
		req, err := parseFakeRequest(hdr.Command, payload)
		if err != nil {
			// Push a minimal request so the test can still see what arrived.
			fs.requests <- fakeRequest{rawPayload: payload, command: hdr.Command}
			continue
		}
		fs.requests <- req
	}
}

// close shuts down the fake's server-side socket; the readLoop will exit.
func (fs *fakeCIPServer) close() {
	_ = fs.serverEnd.Close()
}

// awaitRequest returns the next request the client sent, or fails the test if
// nothing arrives within the timeout.
func (fs *fakeCIPServer) awaitRequest(timeout time.Duration) fakeRequest {
	fs.t.Helper()
	select {
	case r := <-fs.requests:
		return r
	case <-time.After(timeout):
		fs.t.Fatalf("timed out waiting for client request")
		return fakeRequest{}
	}
}

// replyConnectedRead sends a SendUnitData response wrapping a connected-read reply
// with the given general status (0x00 OK, 0x06 PartialTransfer, etc.), CIP type,
// and payload bytes (the data portion only — type and unknown bytes are added here).
func (fs *fakeCIPServer) replyConnectedRead(svc CIPService, seq uint16, status byte, cipType CIPType, data []byte) {
	fs.t.Helper()
	// Payload after status: Type + Unknown(0) + data bytes
	var inner bytes.Buffer
	inner.WriteByte(byte(cipType))
	inner.WriteByte(0)
	inner.Write(data)
	body := buildConnectedReply(seq, 0x00112233, svc, status, inner.Bytes())
	fs.sendFrame(cipCommandSendUnitData, body)
}

// sendFrame writes a full EIP frame (header + body) back to the client.
func (fs *fakeCIPServer) sendFrame(cmd CIPCommand, body []byte) {
	fs.t.Helper()
	hdr := eipHeader{
		Command:       cmd,
		Length:        uint16(len(body)),
		SessionHandle: 0xDEADBEEF,
	}
	if err := binary.Write(fs.serverEnd, binary.LittleEndian, hdr); err != nil {
		fs.t.Fatalf("fake server failed to write EIP header: %v", err)
	}
	if _, err := fs.serverEnd.Write(body); err != nil {
		fs.t.Fatalf("fake server failed to write body: %v", err)
	}
}

// parseFakeRequest decodes a SendUnitData request payload enough to expose the
// service, sequence counter, path, element count, and (for FragRead) offset.
func parseFakeRequest(cmd CIPCommand, payload []byte) (fakeRequest, error) {
	req := fakeRequest{rawPayload: payload, command: cmd}
	if cmd != cipCommandSendUnitData {
		return req, nil
	}
	r := bytes.NewReader(payload)

	var pre msgPreItemData
	if err := binary.Read(r, binary.LittleEndian, &pre); err != nil {
		return req, fmt.Errorf("pre-item header: %w", err)
	}
	items, err := readItems(r)
	if err != nil {
		return req, fmt.Errorf("read items: %w", err)
	}
	if len(items) < 2 {
		return req, fmt.Errorf("expected >=2 items, got %d", len(items))
	}

	connData := items[1]
	connData.Reset()
	if err := connData.DeSerialize(&req.seq); err != nil {
		return req, fmt.Errorf("seq: %w", err)
	}
	if err := connData.DeSerialize(&req.service); err != nil {
		return req, fmt.Errorf("service: %w", err)
	}
	var pathWords byte
	if err := connData.DeSerialize(&pathWords); err != nil {
		return req, fmt.Errorf("path length: %w", err)
	}
	pathLen := int(pathWords) * 2
	req.pathBytes = make([]byte, pathLen)
	if pathLen > 0 {
		if err := connData.DeSerialize(&req.pathBytes); err != nil {
			return req, fmt.Errorf("path bytes: %w", err)
		}
	}
	if err := connData.DeSerialize(&req.elements); err != nil {
		return req, fmt.Errorf("elements: %w", err)
	}
	if req.service == CIPService_FragRead {
		if err := connData.DeSerialize(&req.offset); err != nil {
			return req, fmt.Errorf("offset: %w", err)
		}
	}
	return req, nil
}

// buildConnectedReply assembles the SendUnitData body for a connected-mode reply.
// Layout: msgPreItemData + item count + item0 ConnectionAddress + item1 ConnectedData.
// Item 1 carries: seq u16 + service u8 + reserved=0 + status u8 + statusExtended=0 + inner.
func buildConnectedReply(seq uint16, otConnID uint32, svc CIPService, status byte, inner []byte) []byte {
	var item1 bytes.Buffer
	_ = binary.Write(&item1, binary.LittleEndian, seq)
	item1.WriteByte(byte(svc.AsResponse()))
	item1.WriteByte(0)      // reserved
	item1.WriteByte(status) // general status
	item1.WriteByte(0)      // status extended word count
	item1.Write(inner)

	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.LittleEndian, uint32(0)) // InterfaceHandle
	_ = binary.Write(&buf, binary.LittleEndian, uint16(0)) // Timeout
	_ = binary.Write(&buf, binary.LittleEndian, uint16(2)) // item count

	// item 0: ConnectionAddress (id, length=4, payload=otConnID)
	_ = binary.Write(&buf, binary.LittleEndian, cipItem_ConnectionAddress)
	_ = binary.Write(&buf, binary.LittleEndian, uint16(4))
	_ = binary.Write(&buf, binary.LittleEndian, otConnID)

	// item 1: ConnectedData
	_ = binary.Write(&buf, binary.LittleEndian, cipItem_ConnectedData)
	_ = binary.Write(&buf, binary.LittleEndian, uint16(item1.Len()))
	buf.Write(item1.Bytes())

	return buf.Bytes()
}

// TestFakeServerSmokeRead validates the test fixtures themselves: a normal status
// read response built by the helpers must be parsed correctly by Read_single.
// Guards against regressions in the fixture code before it is used to drive RED
// tests for partial transfer.
func TestFakeServerSmokeRead(t *testing.T) {
	client, fs := newFakeCIPClient(t)

	done := make(chan error, 1)
	var got int32
	go func() {
		done <- client.Read("MyDint", &got)
	}()

	req := fs.awaitRequest(time.Second)
	if req.service != CIPService_Read {
		t.Fatalf("expected service Read, got %v", req.service)
	}
	if req.elements != 1 {
		t.Fatalf("expected elements=1, got %d", req.elements)
	}

	var data bytes.Buffer
	_ = binary.Write(&data, binary.LittleEndian, int32(42))
	fs.replyConnectedRead(CIPService_Read, req.seq, 0x00, CIPTypeDINT, data.Bytes())

	if err := <-done; err != nil {
		t.Fatalf("Read returned error: %v", err)
	}
	if got != 42 {
		t.Fatalf("got %d, want 42", got)
	}
}

// dintFragment describes one server-side response in a multi-fragment exchange.
// status: general status byte (0x00 OK terminates, 0x06 PartialTransfer continues,
// anything else is an error the client must surface).
// values: int32 elements carried in this fragment's data section.
type dintFragment struct {
	status byte
	values []int32
}

// serveDintFragments drives a multi-fragment DINT read end-to-end from the fake
// server side. It pulls one request per fragment off the wire, validates that
// each fragment after the first is a FragRead with the cumulative byte offset,
// and writes the canned reply.
func serveDintFragments(t *testing.T, fs *fakeCIPServer, fragments []dintFragment) {
	t.Helper()
	var receivedOffset uint32
	var emittedBytes uint32
	for i, frag := range fragments {
		req := fs.awaitRequest(time.Second)
		if i == 0 {
			if req.service != CIPService_Read {
				t.Errorf("fragment 0: expected service Read, got %v", req.service)
				return
			}
		} else {
			if req.service != CIPService_FragRead {
				t.Errorf("fragment %d: expected service FragRead, got %v", i, req.service)
				return
			}
			if req.offset != receivedOffset {
				t.Errorf("fragment %d: expected offset %d, got %d", i, receivedOffset, req.offset)
				return
			}
		}
		var data bytes.Buffer
		for _, v := range frag.values {
			_ = binary.Write(&data, binary.LittleEndian, v)
		}
		fs.replyConnectedRead(req.service, req.seq, frag.status, CIPTypeDINT, data.Bytes())
		emittedBytes += uint32(data.Len())
		receivedOffset = emittedBytes
	}
}

// makeDintRange returns a slice of count int32 values starting at start, used to
// build expectation arrays for partial-transfer tests.
func makeDintRange(start, count int) []int32 {
	out := make([]int32, count)
	for i := range out {
		out[i] = int32(start + i)
	}
	return out
}

// TestReadPartialTransfer covers the four behaviours the FragRead loop must
// satisfy. The three multi-fragment cases fail today and should pass once
// Read_single detects CIPStatus_PartialTransfer (0x06) and follows up with
// FragRead requests using cumulative byte offsets. The single-fragment case is
// a regression guard that must keep passing through every refactor.
func TestReadPartialTransfer(t *testing.T) {
	cases := []struct {
		name      string
		fragments []dintFragment
		wantLen   int
		wantFirst int32
		wantLast  int32
		wantErr   bool
	}{
		{
			name: "two-fragment DINT[1500] read",
			fragments: []dintFragment{
				{status: byte(CIPStatus_PartialTransfer), values: makeDintRange(0, 900)},
				{status: 0x00, values: makeDintRange(900, 600)},
			},
			wantLen:   1500,
			wantFirst: 0,
			wantLast:  1499,
		},
		{
			name: "three-fragment DINT[2000] read",
			fragments: []dintFragment{
				{status: byte(CIPStatus_PartialTransfer), values: makeDintRange(0, 800)},
				{status: byte(CIPStatus_PartialTransfer), values: makeDintRange(800, 800)},
				{status: 0x00, values: makeDintRange(1600, 400)},
			},
			wantLen:   2000,
			wantFirst: 0,
			wantLast:  1999,
		},
		{
			name: "single fragment DINT[10] (no 0x06) backward-compat",
			fragments: []dintFragment{
				{status: 0x00, values: makeDintRange(0, 10)},
			},
			wantLen:   10,
			wantFirst: 0,
			wantLast:  9,
		},
		{
			name: "error mid-fragment surfaces a clean error",
			fragments: []dintFragment{
				{status: byte(CIPStatus_PartialTransfer), values: makeDintRange(0, 100)},
				{status: byte(CIPStatus_PathDestinationUnknown), values: nil},
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client, fs := newFakeCIPClient(t)

			type result struct {
				val any
				err error
			}
			done := make(chan result, 1)
			go func() {
				v, e := client.Read_single("MyDintArray", CIPTypeDINT, uint16(tc.wantLen))
				done <- result{v, e}
			}()

			serveDintFragments(t, fs, tc.fragments)

			select {
			case r := <-done:
				if tc.wantErr {
					if r.err == nil {
						t.Fatalf("expected error, got value %v", r.val)
					}
					return
				}
				if r.err != nil {
					t.Fatalf("unexpected error: %v", r.err)
				}
				got, ok := r.val.([]any)
				if !ok {
					t.Fatalf("expected []any, got %T", r.val)
				}
				if len(got) != tc.wantLen {
					t.Fatalf("element count: got %d, want %d", len(got), tc.wantLen)
				}
				first, ok := got[0].(int32)
				if !ok {
					t.Fatalf("first element type: got %T, want int32", got[0])
				}
				if first != tc.wantFirst {
					t.Errorf("first element: got %d, want %d", first, tc.wantFirst)
				}
				last, ok := got[len(got)-1].(int32)
				if !ok {
					t.Fatalf("last element type: got %T, want int32", got[len(got)-1])
				}
				if last != tc.wantLast {
					t.Errorf("last element: got %d, want %d", last, tc.wantLast)
				}
			case <-time.After(2 * time.Second):
				t.Fatal("Read_single hung")
			}
		})
	}
}

// TestReadPartialTransferStructIsRejected locks in the explicit scope guard:
// when the first response indicates a structured tag type (Type=0xA0) AND
// status=0x06, Read_single must surface a clear error instead of attempting
// the FragRead loop. Struct partial transfer requires per-fragment
// StructHandle deduplication that this PR does not implement; surfacing an
// error keeps callers from silently consuming corrupt bytes.
func TestReadPartialTransferStructIsRejected(t *testing.T) {
	client, fs := newFakeCIPClient(t)

	type result struct {
		val any
		err error
	}
	done := make(chan result, 1)
	go func() {
		v, e := client.Read_single("MyStructArray", CIPTypeStruct, 4)
		done <- result{v, e}
	}()

	// Reply with status=0x06 and a struct type marker (Type=0xA0). The data
	// portion is arbitrary — the guard fires before any byte parsing.
	req := fs.awaitRequest(time.Second)
	if req.service != CIPService_Read {
		t.Fatalf("expected service Read, got %v", req.service)
	}
	fs.replyConnectedRead(CIPService_Read, req.seq, byte(CIPStatus_PartialTransfer), CIPTypeStruct, []byte{0xCE, 0x0F, 0x01, 0x02, 0x03})

	select {
	case r := <-done:
		if r.err == nil {
			t.Fatalf("expected struct partial-transfer to be rejected, got value %v", r.val)
		}
		if !contains(r.err.Error(), "structured tag types are not yet supported") {
			t.Errorf("error message did not mention scope guard: %v", r.err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Read_single hung instead of returning the scope-guard error")
	}
}

// contains is a tiny substring helper to keep error-message assertions readable
// without importing strings into every test case.
func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
