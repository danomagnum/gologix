package gologix

import (
	"encoding/binary"
	"net"
	"testing"
	"time"
)

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
