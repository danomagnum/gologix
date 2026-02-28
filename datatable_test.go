package gologix

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestDataTableTypeSize(t *testing.T) {
	tests := []struct {
		cipType  CIPType
		expected uint16
	}{
		{CIPTypeBOOL, 1},
		{CIPTypeSINT, 1},
		{CIPTypeBYTE, 1},
		{CIPTypeUSINT, 1},
		{CIPTypeINT, 2},
		{CIPTypeUINT, 2},
		{CIPTypeDINT, 4},
		{CIPTypeUDINT, 4},
		{CIPTypeREAL, 4},
		{CIPTypeLINT, 8},
		{CIPTypeULINT, 8},
		{CIPTypeLREAL, 8},
		{CIPTypeLWORD, 8},
		{CIPTypeSTRING, 88},
		{CIPTypeStruct, 88},
	}

	for _, tt := range tests {
		size, err := dataTableTypeSize(tt.cipType)
		if err != nil {
			t.Errorf("dataTableTypeSize(%v) returned error: %v", tt.cipType, err)
			continue
		}
		if size != tt.expected {
			t.Errorf("dataTableTypeSize(%v) = %d, want %d", tt.cipType, size, tt.expected)
		}
	}
}

func TestDataTableTypeSizeUnsupported(t *testing.T) {
	_, err := dataTableTypeSize(CIPTypeUnknown)
	if err == nil {
		t.Error("dataTableTypeSize(CIPTypeUnknown) should return error")
	}
}

func TestDataTablePath(t *testing.T) {
	tests := []struct {
		instanceID uint16
		expected   []byte
	}{
		{0x0000, []byte{0x20, 0xB2, 0x25, 0x00, 0x00, 0x00}},
		{0x0021, []byte{0x20, 0xB2, 0x25, 0x00, 0x21, 0x00}},
		{0x1234, []byte{0x20, 0xB2, 0x25, 0x00, 0x34, 0x12}},
		{0xFFFF, []byte{0x20, 0xB2, 0x25, 0x00, 0xFF, 0xFF}},
	}

	for _, tt := range tests {
		path := dataTablePath(tt.instanceID)
		if !bytes.Equal(path, tt.expected) {
			t.Errorf("dataTablePath(0x%04X) = %X, want %X", tt.instanceID, path, tt.expected)
		}
	}
}

func TestDataTablePathLength(t *testing.T) {
	// Path must always be 6 bytes = 3 words
	path := dataTablePath(0x0021)
	if len(path) != 6 {
		t.Errorf("dataTablePath should return 6 bytes, got %d", len(path))
	}
	// Path size in words should be 3
	if len(path)/2 != 3 {
		t.Errorf("dataTablePath word count should be 3, got %d", len(path)/2)
	}
}

func TestDecodeDataTableValueDINT(t *testing.T) {
	raw := make([]byte, 4)
	binary.LittleEndian.PutUint32(raw, uint32(0x12345678))

	val, err := decodeDataTableValue(CIPTypeDINT, raw)
	if err != nil {
		t.Fatalf("decodeDataTableValue(DINT) error: %v", err)
	}
	dint, ok := val.(int32)
	if !ok {
		t.Fatalf("expected int32, got %T", val)
	}
	if dint != 0x12345678 {
		t.Errorf("expected 0x12345678, got 0x%X", dint)
	}
}

func TestDecodeDataTableValueREAL(t *testing.T) {
	raw := make([]byte, 4)
	binary.LittleEndian.PutUint32(raw, 0x41200000) // 10.0 in IEEE 754

	val, err := decodeDataTableValue(CIPTypeREAL, raw)
	if err != nil {
		t.Fatalf("decodeDataTableValue(REAL) error: %v", err)
	}
	real, ok := val.(float32)
	if !ok {
		t.Fatalf("expected float32, got %T", val)
	}
	if real != 10.0 {
		t.Errorf("expected 10.0, got %f", real)
	}
}

func TestDecodeDataTableValueINT(t *testing.T) {
	raw := make([]byte, 2)
	binary.LittleEndian.PutUint16(raw, uint16(0x00FF)) // 255 as int16

	val, err := decodeDataTableValue(CIPTypeINT, raw)
	if err != nil {
		t.Fatalf("decodeDataTableValue(INT) error: %v", err)
	}
	intVal, ok := val.(int16)
	if !ok {
		t.Fatalf("expected int16, got %T", val)
	}
	if intVal != 255 {
		t.Errorf("expected 255, got %d", intVal)
	}
}

func TestDecodeDataTableValueBOOL(t *testing.T) {
	val, err := decodeDataTableValue(CIPTypeBOOL, []byte{0x01})
	if err != nil {
		t.Fatalf("decodeDataTableValue(BOOL) error: %v", err)
	}
	boolVal, ok := val.(bool)
	if !ok {
		t.Fatalf("expected bool, got %T", val)
	}
	if !boolVal {
		t.Error("expected true, got false")
	}
}

func TestDecodeDataTableValueSTRING(t *testing.T) {
	// Standard AB STRING: 4 byte length + 82 byte data + 2 pad = 88 bytes
	raw := make([]byte, 88)
	binary.LittleEndian.PutUint32(raw[0:4], 5) // length = 5
	copy(raw[4:9], "Hello")                     // "Hello"

	val, err := decodeDataTableValue(CIPTypeSTRING, raw)
	if err != nil {
		t.Fatalf("decodeDataTableValue(STRING) error: %v", err)
	}
	str, ok := val.(string)
	if !ok {
		t.Fatalf("expected string, got %T", val)
	}
	if str != "Hello" {
		t.Errorf("expected %q, got %q", "Hello", str)
	}
}

func TestDecodeDataTableValueSTRINGEmpty(t *testing.T) {
	raw := make([]byte, 88)
	binary.LittleEndian.PutUint32(raw[0:4], 0) // length = 0

	val, err := decodeDataTableValue(CIPTypeSTRING, raw)
	if err != nil {
		t.Fatalf("decodeDataTableValue(STRING empty) error: %v", err)
	}
	str, ok := val.(string)
	if !ok {
		t.Fatalf("expected string, got %T", val)
	}
	if str != "" {
		t.Errorf("expected empty string, got %q", str)
	}
}

func TestDecodeDataTableValueLREAL(t *testing.T) {
	raw := make([]byte, 8)
	binary.LittleEndian.PutUint64(raw, 0x4024000000000000) // 10.0 in float64

	val, err := decodeDataTableValue(CIPTypeLREAL, raw)
	if err != nil {
		t.Fatalf("decodeDataTableValue(LREAL) error: %v", err)
	}
	lreal, ok := val.(float64)
	if !ok {
		t.Fatalf("expected float64, got %T", val)
	}
	if lreal != 10.0 {
		t.Errorf("expected 10.0, got %f", lreal)
	}
}
