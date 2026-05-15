package gologix

import (
	"encoding/binary"
	"testing"
)

// TestSerializeStringSliceEmpty verifies that an empty or nil []string
// passed to CIPItem.Serialize is treated as a no-op: returns nil and
// does not write any payload bytes. write_single declares elements=0
// in the IOI footer for an empty slice, so emitting a partial slot
// would desync the wire format and the controller would reject the
// request with a confusing CIP status.
func TestSerializeStringSliceEmpty(t *testing.T) {
	cases := []struct {
		name string
		in   []string
	}{
		{name: "nil", in: nil},
		{name: "empty", in: []string{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			item := CIPItem{}
			if err := item.Serialize(tc.in); err != nil {
				t.Fatalf("Serialize(%v) returned error: %v", tc.in, err)
			}
			if len(item.Data) != 0 {
				t.Fatalf("Serialize(%v) wrote %d bytes; want 0", tc.in, len(item.Data))
			}
		})
	}
}

// TestSerializeStringSliceShape verifies the wire shape for a non-empty
// []string: each element occupies a fixed 88-byte slot (4-byte LEN +
// 84-byte data buffer). N elements produce exactly N*88 bytes.
// long-over-84 also asserts the LEN header is clamped to 84 so the
// wire header never claims more bytes than the slot actually carries.
func TestSerializeStringSliceShape(t *testing.T) {
	cases := []struct {
		name    string
		in      []string
		wantLen uint32
	}{
		{name: "single", in: []string{"hello"}, wantLen: 5},
		{name: "multi", in: []string{"a", "bb", "ccc"}, wantLen: 1},
		{name: "long-exact-84", in: []string{string(make([]byte, 84))}, wantLen: 84},
		{name: "long-over-84", in: []string{string(make([]byte, 200))}, wantLen: 84},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			item := CIPItem{}
			if err := item.Serialize(tc.in); err != nil {
				t.Fatalf("Serialize(%v) returned error: %v", tc.in, err)
			}
			want := len(tc.in) * 88
			if len(item.Data) != want {
				t.Fatalf("Serialize(%v) wrote %d bytes; want %d", tc.in, len(item.Data), want)
			}
			// LEN header of the first element must match wantLen.
			gotLen := binary.LittleEndian.Uint32(item.Data[0:4])
			if gotLen != tc.wantLen {
				t.Fatalf("first-element LEN header = %d; want %d", gotLen, tc.wantLen)
			}
		})
	}
}

// TestSerializeStringSliceVariadicContinues verifies that an empty
// []string sandwiched between other Serialize arguments does NOT
// short-circuit the variadic loop: arguments before and after the
// empty slice still get serialized.
func TestSerializeStringSliceVariadicContinues(t *testing.T) {
	item := CIPItem{}
	if err := item.Serialize(uint32(0xDEADBEEF), []string{}, uint32(0xCAFEBABE)); err != nil {
		t.Fatalf("Serialize(...) returned error: %v", err)
	}
	// Two uint32 values = 8 bytes. Empty []string contributes 0.
	if len(item.Data) != 8 {
		t.Fatalf("Serialize wrote %d bytes; want 8 (two uint32 with empty []string between)", len(item.Data))
	}
}
