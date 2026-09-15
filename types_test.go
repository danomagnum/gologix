package gologix

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
	"time"
)

// TestGetBitSINT locks in that getBit accepts the int8 that readValue produces
// for a SINT. Before the fix getBit asserted byte, so every SINT bit read
// failed with a misleading "must be 0-7" error (issue #71).
func TestGetBitSINT(t *testing.T) {
	cases := []struct {
		raw    byte
		bitpos int
		want   bool
	}{
		{0x08, 3, true},
		{0x08, 2, false},
		{0x80, 7, true}, // sign bit of a negative int8
		{0x7F, 7, false},
	}
	for _, c := range cases {
		v, err := readValue(CIPTypeSINT, bytes.NewReader([]byte{c.raw}))
		if err != nil {
			t.Fatalf("readValue(SINT 0x%02X): %v", c.raw, err)
		}
		got, err := getBit(CIPTypeSINT, v, c.bitpos)
		if err != nil {
			t.Errorf("getBit(SINT 0x%02X, bit %d): %v", c.raw, c.bitpos, err)
			continue
		}
		if got != c.want {
			t.Errorf("getBit(SINT 0x%02X, bit %d) = %v; want %v", c.raw, c.bitpos, got, c.want)
		}
	}
}

// TestGetBitOutOfRange: a bit position outside the type's width must be an
// error. Returning (false, nil) is indistinguishable from a bit that is
// genuinely off, so a typo in a tag address became good data (issue #71).
func TestGetBitOutOfRange(t *testing.T) {
	cases := []struct {
		name   string
		typ    CIPType
		v      any
		bitpos int
	}{
		{"BOOL bit 1", CIPTypeBOOL, true, 1},
		{"SINT bit 8", CIPTypeSINT, int8(-1), 8},
		{"BYTE bit -1", CIPTypeBYTE, byte(0xFF), -1},
		{"INT bit 16", CIPTypeINT, int16(-1), 16},
		{"DINT bit 40", CIPTypeDINT, int32(-1), 40},
		{"UDINT bit 32", CIPTypeUDINT, uint32(0xFFFFFFFF), 32},
		{"LINT bit 64", CIPTypeLINT, int64(-1), 64},
		{"LWORD bit -1", CIPTypeLWORD, uint64(1), -1},
	}
	for _, c := range cases {
		got, err := getBit(c.typ, c.v, c.bitpos)
		if err == nil {
			t.Errorf("%s: got (%v, nil); want error", c.name, got)
		}
	}
}

// TestGetBitWrongGoType keeps the existing contract: a value whose Go type
// does not match the CIP type is still an error, not a silent false.
func TestGetBitWrongGoType(t *testing.T) {
	if _, err := getBit(CIPTypeDINT, int16(1), 0); err == nil {
		t.Fatal("DINT with an int16 value: want error, got nil")
	}
}

// TestReadSINTBitAccess drives the user-visible path for issue #71:
// Read("Tag.3", *bool) against a SINT tag goes through readValue then getBit.
func TestReadSINTBitAccess(t *testing.T) {
	client, fs := newFakeCIPClient(t)

	done := make(chan error, 1)
	var got bool
	go func() {
		done <- client.Read("MySint.3", &got)
	}()

	req := fs.awaitRequest(time.Second)
	if req.service != CIPService_Read {
		t.Fatalf("expected service Read, got %v", req.service)
	}
	fs.replyConnectedRead(CIPService_Read, req.seq, 0x00, CIPTypeSINT, []byte{0x08})

	if err := <-done; err != nil {
		t.Fatalf("Read returned error: %v", err)
	}
	if !got {
		t.Fatal("bit 3 of SINT 0x08 read as false; want true")
	}
}

// TestReadMultiOutOfRangeBit covers the batch path. readList used to log and
// skip a getBit error, leaving a nil result that ReadMulti then handed to
// reflect.Value.Set, which panics. An out-of-range bit position must come back
// from ReadMulti as an error.
func TestReadMultiOutOfRangeBit(t *testing.T) {
	client, fs := newFakeCIPClient(t)

	type tags struct {
		Bit bool `gologix:"MyInt.20"`
	}
	done := make(chan error, 1)
	go func() {
		var got tags
		done <- client.ReadMulti(&got)
	}()

	req := fs.awaitRequest(time.Second)
	var data bytes.Buffer
	_ = binary.Write(&data, binary.LittleEndian, int16(-1))
	fs.replyMultiRead(req.seq, [][]byte{multiReadReplyBody(CIPTypeINT, 0x00, data.Bytes())})

	err := <-done
	if err == nil {
		t.Fatal("ReadMulti of bit 20 on an INT: want error, got nil")
	}
	if !strings.Contains(err.Error(), "out of range") {
		t.Fatalf("want the out-of-range error from getBit, got: %v", err)
	}
}
