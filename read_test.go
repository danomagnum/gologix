package gologix

import (
	"testing"
)

func TestReadSingle(t *testing.T) {
	var tests = []struct {
		tag   string
		wants any
	}{
		{"TestInt", int16(999)},
		{"TestSint", byte(117)},
		{"TestDint", int32(36)},
		{"TestReal", float32(93.45)},
		{"TestDintArr[0]", int32(4351)},
		{"TestDintArr[2]", int32(4353)},
		{"TestUDT.Field1", int32(85456)},
		{"TestUDT.Field2", float32(123.456)},
		{"TestUDTArr[2].Field1", int32(16)},
		{"TestUDTArr[2].Field2", float32(15.0)},
	}

	client := &Client{IPAddress: "192.168.2.241"}
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()

	for _, tt := range tests {

		t.Run(tt.tag, func(t2 *testing.T) {
			ct := GoVarToCIPType(tt.wants)
			value, err := client.read_single(tt.tag, ct, 1)
			if err != nil {
				t2.Errorf("Problem reading %s. %v", tt.tag, err)
				return
			}
			if value != tt.wants {
				t2.Errorf("wanted %v got %v", tt.wants, value)
			}
		})
	}

}

func TestReadMulti(t *testing.T) {
	type test_str struct {
		TestSint          byte    `gologix:"TestSint"`
		TestInt           int16   `gologix:"TestInt"`
		TestDint          int32   `gologix:"TestDint"`
		TestReal          float32 `gologix:"TestReal"`
		TestDintArr0      int32   `gologix:"testdintarr[0]"`
		TestDintArr2      int32   `gologix:"testdintarr[2]"`
		TestUDTField1     int32   `gologix:"testudt.field1"`
		TestUDTField2     float32 `gologix:"testudt.field2"`
		TestUDTArr2Field1 int32   `gologix:"testudtarr[2].field1"`
		TestUDTArr2Field2 float32 `gologix:"testudtarr[2].field2"`
	}
	read := test_str{}
	wants := test_str{
		TestSint:          117,
		TestInt:           999,
		TestDint:          36,
		TestReal:          93.45,
		TestDintArr0:      4351,
		TestDintArr2:      4353,
		TestUDTField1:     85456,
		TestUDTField2:     123.456,
		TestUDTArr2Field1: 16,
		TestUDTArr2Field2: 15.0,
	}

	client := &Client{IPAddress: "192.168.2.241"}
	client.Connect()
	defer client.Disconnect()

	err := client.read_multi(&read, CIPTypeStruct, 1)
	if err != nil {
		t.Errorf("Problem reading. %v", err)
		return
	}
	if read != wants {
		t.Errorf("wanted %v got %v", wants, read)
	}
}
