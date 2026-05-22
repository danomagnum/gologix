package gologix

import (
	"bytes"
	"encoding/binary"
	"testing"
)

// TestParseClassInstancePath covers the small EPATH decoder used by the
// Identity Object probe handler. The shape FT Linx and pycomm3 send is
// `0x20 <class> 0x24 <instance>` (8-bit Class + 8-bit Instance). Anything
// else is rejected.
func TestParseClassInstancePath(t *testing.T) {
	cases := []struct {
		name     string
		in       []byte
		class    CIPClass
		instance CIPInstance
		ok       bool
	}{
		{name: "Identity Class 1 Instance 1", in: []byte{0x20, 0x01, 0x24, 0x01}, class: 1, instance: 1, ok: true},
		{name: "MessageRouter Class 2 Instance 1", in: []byte{0x20, 0x02, 0x24, 0x01}, class: 2, instance: 1, ok: true},
		{name: "wrong length", in: []byte{0x20, 0x01}, ok: false},
		{name: "16-bit class segment", in: []byte{0x21, 0x00, 0x01, 0x00, 0x24, 0x01}, class: 1, instance: 1, ok: true},
		{name: "16-bit class segment and instance", in: []byte{0x21, 0x00, 0x01, 0x00, 0x25, 0x00, 0x01, 0x00}, class: 1, instance: 1, ok: true},
		{name: "missing class header", in: []byte{0x00, 0x01, 0x24, 0x01}, ok: false},
		{name: "missing instance header", in: []byte{0x20, 0x01, 0x00, 0x01}, ok: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			class, instance, err := parseClassInstancePath(bytes.NewBuffer(tc.in))
			wasok := err == nil
			if wasok != tc.ok {
				t.Fatalf("ok=%v; want %v got %v", !wasok, tc.ok, err)
			}
			if !tc.ok {
				return
			}
			if class != tc.class || instance != tc.instance {
				t.Fatalf("class=%d instance=%d; want %d %d", class, instance, tc.class, tc.instance)
			}
		})
	}
}

// TestBuildIdentityGetAttributesAllResponse renders an Identity response
// against the same Attributes map NewServer seeds and locks the wire bytes
// in: 7 attributes back-to-back, ProductName as a SHORT_STRING (length byte
// then ASCII). Regression guard against accidental drift in field order
// or sizes.
func TestBuildIdentityGetAttributesAllResponse(t *testing.T) {
	srv := NewServer(nil)
	srv.Attributes[6] = uint32(0xDEADBEEF) // make Serial deterministic

	got, err := buildIdentityGetAttributesAllResponse(srv.Attributes)
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	productName, _ := srv.Attributes[7].(string)
	wantLen := 2 + 2 + 2 + 2 + 2 + 4 + 1 + len(productName) // 15 + name
	if len(got) != wantLen {
		t.Fatalf("payload length = %d; want %d", len(got), wantLen)
	}

	r := bytes.NewReader(got)
	var vendor, deviceType, productCode, revision, status int16
	var serial uint32
	for i, target := range []any{&vendor, &deviceType, &productCode, &revision, &status, &serial} {
		if err := binary.Read(r, binary.LittleEndian, target); err != nil {
			t.Fatalf("read field %d: %v", i, err)
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
	if serial != uint32(0xDEADBEEF) {
		t.Errorf("SerialNumber = 0x%08X; want 0xDEADBEEF", serial)
	}

	nameLen, err := r.ReadByte()
	if err != nil {
		t.Fatalf("read name length: %v", err)
	}
	if int(nameLen) != len(productName) {
		t.Fatalf("SHORT_STRING length byte = %d; want %d", nameLen, len(productName))
	}
	nameBytes := make([]byte, nameLen)
	if _, err := r.Read(nameBytes); err != nil {
		t.Fatalf("read name bytes: %v", err)
	}
	if string(nameBytes) != productName {
		t.Errorf("ProductName = %q; want %q", string(nameBytes), productName)
	}
}

// TestBuildIdentityGetAttributesAllResponseTypeChecks verifies the
// attribute-map type contract — missing or wrong-typed entries surface as
// errors instead of crashing the server goroutine.
func TestBuildIdentityGetAttributesAllResponseTypeChecks(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(map[CIPAttribute]any)
	}{
		{name: "missing VendorID", mutate: func(a map[CIPAttribute]any) { delete(a, 1) }},
		{name: "wrong VendorID type", mutate: func(a map[CIPAttribute]any) { a[1] = "not-an-int16" }},
		{name: "wrong Serial type", mutate: func(a map[CIPAttribute]any) { a[6] = int32(0) }},
		{name: "wrong ProductName type", mutate: func(a map[CIPAttribute]any) { a[7] = 42 }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := NewServer(nil)
			tc.mutate(srv.Attributes)
			if _, err := buildIdentityGetAttributesAllResponse(srv.Attributes); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

// TestBuildIdentityGetAttributesAllResponseTruncatesLongName guards the
// SHORT_STRING length byte's 255-byte limit. A ProductName longer than
// 255 chars must be truncated, not overflow into the next packet.
func TestBuildIdentityGetAttributesAllResponseTruncatesLongName(t *testing.T) {
	srv := NewServer(nil)
	srv.Attributes[6] = uint32(0)
	long := make([]byte, 300)
	for i := range long {
		long[i] = 'A'
	}
	srv.Attributes[7] = string(long)

	got, err := buildIdentityGetAttributesAllResponse(srv.Attributes)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	// Skip the 14-byte fixed prefix; next byte is the SHORT_STRING length.
	if got[14] != 255 {
		t.Fatalf("SHORT_STRING length byte = %d; want 255 (truncation)", got[14])
	}
	if len(got) != 14+1+255 {
		t.Fatalf("payload length = %d; want %d", len(got), 14+1+255)
	}
}
