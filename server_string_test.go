package gologix

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"
	"time"
)

// TestParseWriteValuesString crafts the exact byte sequence a Logix STRING
// write request carries on the wire (matching what pylogix sends — see
// .pi/test-evidence/2026-05-15/37-wire-stage3b.log) and verifies the
// parser extracts a Go string with the trimmed payload. Locks in the
// type_info_length-in-bytes interpretation that bit us during external
// validation.
func TestParseWriteValuesString(t *testing.T) {
	const want = "pylogix-round-trip-check"

	var buf bytes.Buffer
	buf.WriteByte(0xA0)                            // typ = CIPTypeStruct
	buf.WriteByte(0x02)                            // type_info_length = 2 bytes
	_ = binary.Write(&buf, binary.LittleEndian, cipStringStructCRC)
	_ = binary.Write(&buf, binary.LittleEndian, uint16(1)) // qty
	_ = binary.Write(&buf, binary.LittleEndian, uint32(len(want)))
	data := make([]byte, cipStringDataLen)
	copy(data, want)
	buf.Write(data)

	item := CIPItem{Data: buf.Bytes()}
	results, err := parseWriteValues(&item)
	if err != nil {
		t.Fatalf("parseWriteValues: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d elements; want 1", len(results))
	}
	got, ok := results[0].(string)
	if !ok {
		t.Fatalf("element 0 type = %T; want string", results[0])
	}
	if got != want {
		t.Fatalf("element 0 = %q; want %q", got, want)
	}
}

// TestParseWriteValuesAtomic confirms the historic atomic path keeps
// working — type_info_length is the 0x00 high byte of the uint16 DataType
// the gologix client serializes; parser must treat that as "no type info"
// and fall back to readValue.
func TestParseWriteValuesAtomic(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteByte(byte(CIPTypeDINT))                          // typ = DINT (0xC4)
	buf.WriteByte(0x00)                                       // high byte of uint16 DataType
	_ = binary.Write(&buf, binary.LittleEndian, uint16(2))    // qty
	_ = binary.Write(&buf, binary.LittleEndian, int32(42))    // element 0
	_ = binary.Write(&buf, binary.LittleEndian, int32(-1337)) // element 1

	item := CIPItem{Data: buf.Bytes()}
	results, err := parseWriteValues(&item)
	if err != nil {
		t.Fatalf("parseWriteValues: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d elements; want 2", len(results))
	}
	if got, want := results[0].(int32), int32(42); got != want {
		t.Fatalf("element 0 = %d; want %d", got, want)
	}
	if got, want := results[1].(int32), int32(-1337); got != want {
		t.Fatalf("element 1 = %d; want %d", got, want)
	}
}

// TestParseWriteValuesUnknownStruct rejects struct writes whose CRC does
// not match the Logix STRING UDT. Anything else (custom UDTs) is out of
// scope for this fix and should surface as an explicit error rather than
// the historic silent-success.
func TestParseWriteValuesUnknownStruct(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteByte(0xA0)                                       // typ = CIPTypeStruct
	buf.WriteByte(0x02)                                       // type_info_length = 2 bytes
	_ = binary.Write(&buf, binary.LittleEndian, uint16(0xBEEF))
	_ = binary.Write(&buf, binary.LittleEndian, uint16(1))

	item := CIPItem{Data: buf.Bytes()}
	if _, err := parseWriteValues(&item); err == nil {
		t.Fatal("expected error for non-STRING struct CRC; got nil")
	}
}

// TestCipStringPackerShape locks in the per-element STRING wire layout the
// Logix STRING UDT puts on the wire: 4-byte type segment (0xA0 0x02 +
// StructTypeCRC LE) + 4-byte LEN + 82-byte fixed DATA = 90 bytes total.
// External CIP clients (pylogix, MSG instructions, FactoryTalk) refuse to
// extract values when DATA is shorter than 82 bytes, so the fixed slot is
// non-negotiable.
func TestCipStringPackerShape(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want uint32
	}{
		{name: "short", in: "Hello", want: 5},
		{name: "empty", in: "", want: 0},
		{name: "exact-82", in: string(make([]byte, 82)), want: 82},
		{name: "over-82", in: string(make([]byte, 200)), want: 82},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := cipStringPacker(tc.in).Bytes()
			if len(b) != cipStringSlotLen {
				t.Fatalf("packer wrote %d bytes; want %d", len(b), cipStringSlotLen)
			}
			if b[0] != 0xA0 || b[1] != 0x02 {
				t.Fatalf("type segment header = % x; want a0 02", b[0:2])
			}
			if crc := binary.LittleEndian.Uint16(b[2:4]); crc != cipStringStructCRC {
				t.Fatalf("CRC = 0x%04X; want 0x%04X", crc, cipStringStructCRC)
			}
			if got := binary.LittleEndian.Uint32(b[4:8]); got != tc.want {
				t.Fatalf("LEN = %d; want %d", got, tc.want)
			}
			// DATA region (after LEN) is exactly 82 bytes by construction.
			if len(b[8:]) != cipStringDataLen {
				t.Fatalf("DATA region = %d bytes; want %d", len(b[8:]), cipStringDataLen)
			}
		})
	}
}

// TestServerStringReadRoundTrip verifies a CIP client can read a STRING tag
// from a gologix server in the same process. Uses the hardcoded EIP port
// 44818 so the test skips when another process is bound there.
func TestServerStringReadRoundTrip(t *testing.T) {
	if probe, err := net.Listen("tcp", "0.0.0.0:44818"); err != nil {
		t.Skipf("port 44818 unavailable: %v", err)
	} else {
		probe.Close()
	}

	router := PathRouter{}
	provider := MapTagProvider{}
	path, err := ParsePath("1,0")
	if err != nil {
		t.Fatalf("parse path: %v", err)
	}
	router.Handle(path.Bytes(), &provider)

	const tag = "teststring"
	const want = "Hello World"
	if err := provider.TagWrite(tag, want); err != nil {
		t.Fatalf("seed tag: %v", err)
	}

	srv := NewServer(&router)
	go func() { _ = srv.Serve() }()
	defer func() {
		if srv.TCPListener != nil {
			srv.TCPListener.Close()
		}
		if srv.UDPListener != nil {
			srv.UDPListener.Close()
		}
	}()

	// Wait briefly for the listener to come up.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:44818", 100*time.Millisecond)
		if err == nil {
			conn.Close()
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	client := NewClient("127.0.0.1")
	if err := client.Connect(); err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer client.Disconnect()

	var got string
	if err := client.Read(tag, &got); err != nil {
		t.Fatalf("client read: %v", err)
	}
	if got != want {
		t.Fatalf("read-back mismatch: want %q, got %q", want, got)
	}
}
