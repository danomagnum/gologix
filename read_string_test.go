package gologix

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
	"time"
)

// stringTagData builds the tag-data portion of a single STRING read reply as a
// Logix controller sends it: struct CRC, LEN, DATA, then the 2-byte alignment
// pad. length is written as given so tests can claim more than is really
// there (issue #72).
func stringTagData(length uint32, data []byte) []byte {
	var b bytes.Buffer
	_ = binary.Write(&b, binary.LittleEndian, cipStringStructCRC)
	_ = binary.Write(&b, binary.LittleEndian, length)
	b.Write(data)
	b.Write(make([]byte, cipStringStructPad))
	return b.Bytes()
}

// stringArrayTagData builds the tag-data portion of a STRING array read reply:
// one struct CRC, then LEN + DATA + pad per element. lengths lets a test claim
// a bogus LEN for an element; each element's DATA is the full 82-byte slot.
func stringArrayTagData(lengths []uint32, values []string) []byte {
	var b bytes.Buffer
	_ = binary.Write(&b, binary.LittleEndian, cipStringStructCRC)
	for i, v := range values {
		_ = binary.Write(&b, binary.LittleEndian, lengths[i])
		b.Write(paddedStringData(v))
		b.Write(make([]byte, cipStringStructPad))
	}
	return b.Bytes()
}

// paddedStringData returns the fixed-width DATA slot of a Logix STRING holding s.
func paddedStringData(s string) []byte {
	d := make([]byte, cipStringDataLen)
	copy(d, s)
	return d
}

func assertStringLengthError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("want an error for a STRING header whose length exceeds the data, got nil")
	}
	if !strings.Contains(err.Error(), "string length") {
		t.Fatalf("want a string-length error, got: %v", err)
	}
}

// bogusStringLengths are STRING headers whose Length field does not match the
// bytes that follow. The near-max value is the case from issue #72: taken at
// face value it asks for a 4 GiB allocation.
var bogusStringLengths = []struct {
	name   string
	length uint32
	data   []byte
}{
	{"length exceeds data", 100, []byte("0123456789")},
	{"length near max uint32", 0xFFFFFFFF, []byte{1, 2, 3, 4}},
}

func TestReadStringRoundTrip(t *testing.T) {
	client, fs := newFakeCIPClient(t)

	done := make(chan error, 1)
	var got string
	go func() {
		done <- client.Read("MyString", &got)
	}()

	req := fs.awaitRequest(time.Second)
	fs.replyConnectedRead(CIPService_Read, req.seq, 0x00, CIPTypeStruct, stringTagData(5, paddedStringData("hello")))

	if err := <-done; err != nil {
		t.Fatalf("Read returned error: %v", err)
	}
	if got != "hello" {
		t.Fatalf("got %q, want %q", got, "hello")
	}
}

func TestReadStringRejectsBogusLength(t *testing.T) {
	for _, c := range bogusStringLengths {
		t.Run(c.name, func(t *testing.T) {
			client, fs := newFakeCIPClient(t)

			done := make(chan error, 1)
			var got string
			go func() {
				done <- client.Read("MyString", &got)
			}()

			req := fs.awaitRequest(time.Second)
			fs.replyConnectedRead(CIPService_Read, req.seq, 0x00, CIPTypeStruct, stringTagData(c.length, c.data))

			assertStringLengthError(t, <-done)
		})
	}
}

func TestReadStringArrayRoundTrip(t *testing.T) {
	client, fs := newFakeCIPClient(t)

	done := make(chan error, 1)
	got := make([]string, 2)
	go func() {
		done <- client.Read("MyStrings", got)
	}()

	req := fs.awaitRequest(time.Second)
	data := stringArrayTagData([]uint32{2, 2}, []string{"hi", "yo"})
	fs.replyConnectedRead(CIPService_Read, req.seq, 0x00, CIPTypeStruct, data)

	if err := <-done; err != nil {
		t.Fatalf("Read returned error: %v", err)
	}
	if got[0] != "hi" || got[1] != "yo" {
		t.Fatalf("got %q, want [hi yo]", got)
	}
}

// A STRING array element is a fixed 82-byte slot, so a Length above that used
// to panic with a slice bounds error instead of returning one.
func TestReadStringArrayRejectsBogusLength(t *testing.T) {
	client, fs := newFakeCIPClient(t)

	done := make(chan error, 1)
	got := make([]string, 2)
	go func() {
		done <- client.Read("MyStrings", got)
	}()

	req := fs.awaitRequest(time.Second)
	data := stringArrayTagData([]uint32{200, 2}, []string{"hi", "yo"})
	fs.replyConnectedRead(CIPService_Read, req.seq, 0x00, CIPTypeStruct, data)

	assertStringLengthError(t, <-done)
}

func TestReadListStringRoundTrip(t *testing.T) {
	client, fs := newFakeCIPClient(t)

	type result struct {
		vals []any
		err  error
	}
	done := make(chan result, 1)
	go func() {
		vals, err := client.readList([]tagDesc{{TagName: "MyString", TagType: CIPTypeSTRING, Elements: 1}})
		done <- result{vals, err}
	}()

	req := fs.awaitRequest(time.Second)
	fs.replyMultiRead(req.seq, [][]byte{multiReadReplyBody(CIPTypeStruct, 0x02, stringTagData(5, paddedStringData("hello")))})

	r := <-done
	if r.err != nil {
		t.Fatalf("readList returned error: %v", r.err)
	}
	if len(r.vals) != 1 || r.vals[0] != "hello" {
		t.Fatalf("got %v, want [hello]", r.vals)
	}
}

func TestReadListStringRejectsBogusLength(t *testing.T) {
	for _, c := range bogusStringLengths {
		t.Run(c.name, func(t *testing.T) {
			client, fs := newFakeCIPClient(t)

			done := make(chan error, 1)
			go func() {
				_, err := client.readList([]tagDesc{{TagName: "MyString", TagType: CIPTypeSTRING, Elements: 1}})
				done <- err
			}()

			req := fs.awaitRequest(time.Second)
			fs.replyMultiRead(req.seq, [][]byte{multiReadReplyBody(CIPTypeStruct, 0x02, stringTagData(c.length, c.data))})

			assertStringLengthError(t, <-done)
		})
	}
}

func TestReadListStringArrayRoundTrip(t *testing.T) {
	client, fs := newFakeCIPClient(t)

	type result struct {
		vals []any
		err  error
	}
	done := make(chan result, 1)
	go func() {
		vals, err := client.readList([]tagDesc{{TagName: "MyStrings", TagType: CIPTypeSTRING, Elements: 2}})
		done <- result{vals, err}
	}()

	req := fs.awaitRequest(time.Second)
	data := stringArrayTagData([]uint32{2, 2}, []string{"hi", "yo"})
	fs.replyMultiRead(req.seq, [][]byte{multiReadReplyBody(CIPTypeStruct, 0x02, data)})

	r := <-done
	if r.err != nil {
		t.Fatalf("readList returned error: %v", r.err)
	}
	if len(r.vals) != 1 {
		t.Fatalf("got %d results, want 1", len(r.vals))
	}
	arr, ok := r.vals[0].([]any)
	if !ok || len(arr) != 2 || arr[0] != "hi" || arr[1] != "yo" {
		t.Fatalf("got %v, want [hi yo]", r.vals[0])
	}
}

func TestReadListStringArrayRejectsBogusLength(t *testing.T) {
	client, fs := newFakeCIPClient(t)

	done := make(chan error, 1)
	go func() {
		_, err := client.readList([]tagDesc{{TagName: "MyStrings", TagType: CIPTypeSTRING, Elements: 2}})
		done <- err
	}()

	req := fs.awaitRequest(time.Second)
	data := stringArrayTagData([]uint32{200, 2}, []string{"hi", "yo"})
	fs.replyMultiRead(req.seq, [][]byte{multiReadReplyBody(CIPTypeStruct, 0x02, data)})

	assertStringLengthError(t, <-done)
}
