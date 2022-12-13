package gologix

import (
	"strings"
	"testing"
)

func TestMultiWriteToDict(t *testing.T) {

	type TestUDT struct {
		Field1 int32
		Field2 float32
	}

	type test_str struct {
		TestSint          byte    `gologix:"TestSint"`
		TestInt           int16   `gologix:"TestInt"`
		TestDint          int32   `gologix:"TestDint"`
		TestReal          float32 `gologix:"TestReal"`
		TestDintArr0      int32   `gologix:"testdintarr[0]"`
		TestDintArr0_0    bool    `gologix:"testdintarr[0].0"`
		TestDintArr0_9    bool    `gologix:"testdintarr[0].9"`
		TestDintArr2      int32   `gologix:"testdintarr[2]"`
		TestUDTField1     int32   `gologix:"testudt.field1"`
		TestUDTField2     float32 `gologix:"testudt.field2"`
		TestUDTArr2Field1 int32   `gologix:"testudtarr[2].field1"`
		TestUDTArr2Field2 float32 `gologix:"testudtarr[2].field2"`
	}

	type test_str2 struct {
		TestSint          byte    `gologix:"TestSint"`
		TestInt           int16   `gologix:"TestInt"`
		TestDint          int32   `gologix:"TestDint"`
		TestReal          float32 `gologix:"TestReal"`
		TestDintArr0      int32   `gologix:"testdintarr[0]"`
		TestDintArr0_0    bool    `gologix:"testdintarr[0].0"`
		TestDintArr0_9    bool    `gologix:"testdintarr[0].9"`
		TestDintArr2      int32   `gologix:"testdintarr[2]"`
		TestUDT           TestUDT `gologix:"testudt"`
		TestUDTArr2Field1 int32   `gologix:"testudtarr[2].field1"`
		TestUDTArr2Field2 float32 `gologix:"testudtarr[2].field2"`
	}

	read := test_str{
		TestSint:          117,
		TestInt:           999,
		TestDint:          36,
		TestReal:          93.45,
		TestDintArr0:      4351,
		TestDintArr0_0:    true,
		TestDintArr0_9:    false,
		TestDintArr2:      4353,
		TestUDTField1:     85456,
		TestUDTField2:     123.456,
		TestUDTArr2Field1: 16,
		TestUDTArr2Field2: 15.0,
	}
	read2 := test_str2{
		TestSint:          117,
		TestInt:           999,
		TestDint:          36,
		TestReal:          93.45,
		TestDintArr0:      4351,
		TestDintArr0_0:    true,
		TestDintArr0_9:    false,
		TestDintArr2:      4353,
		TestUDT:           TestUDT{85456, 123.456},
		TestUDTArr2Field1: 16,
		TestUDTArr2Field2: 15.0,
	}
	wants := map[string]interface{}{
		"testsint":             byte(117),
		"testint":              int16(999),
		"testdint":             int32(36),
		"testreal":             float32(93.45),
		"testdintarr[0]":       int32(4351),
		"testdintarr[0].0":     true,
		"testdintarr[0].9":     false,
		"testdintarr[2]":       int32(4353),
		"testudt.field1":       int32(85456),
		"testudt.field2":       float32(123.456),
		"testudtarr[2].field1": int32(16),
		"testudtarr[2].field2": float32(15.0),
	}

	have, err := multi_to_dict(read)
	if err != nil {
		t.Errorf("problem creating dict. %v", err)
	}
	for k := range have {
		if wants[strings.ToLower(k)] != have[k] {
			t.Errorf("key %s is not a match. Have %v want %v", k, have[k], wants[k])
		}

	}

	have, err = multi_to_dict(read2)
	if err != nil {
		t.Errorf("problem creating dict 2. %v", err)
	}
	for k := range have {
		if wants[strings.ToLower(k)] != have[k] {
			t.Errorf("key %s is not a match. Have %v want %v", k, have[k], wants[k])
		}

	}

}

func TestStructToDict(t *testing.T) {

	type TestUDT struct {
		Field1 int32
		Field2 float32
	}

	udt := TestUDT{Field1: 15, Field2: 5.1}
	d, err := udt_to_dict("prefix", udt)
	if err != nil {
		t.Errorf("problem creating dict. %v", err)
	}
	want := map[string]interface{}{
		"prefix.Field1": int32(15),
		"prefix.Field2": float32(5.1),
	}

	for k := range d {
		if d[k] != want[k] {
			t.Errorf("want %v got %v", want[k], d[k])
		}
	}

}
