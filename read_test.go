package gologix

import (
	"errors"
	"fmt"
	"log"
	"net"
	"testing"
	"time"
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

	err := client.ReadMulti(&read, CIPTypeStruct, 1)
	if err != nil {
		t.Errorf("Problem reading. %v", err)
		return
	}
	if read != wants {
		t.Errorf("wanted %v got %v", wants, read)
	}
}

func TestReadTimeout(t *testing.T) {
	//t.Skip("requires timeout that is too long")

	client := &Client{IPAddress: "192.168.2.241"}
	client.SocketTimeout = time.Minute
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()
	fmt.Println("sleeping for 2 minutes. to let the comms timeout")
	for t := 0; t < 24; t++ {
		log.Printf("%d\n", t*5)
		time.Sleep(time.Second * 5)
	}
	//client.Conn.Close()
	//fmt.Println("\nsleep complete")
	value, err := client.ReadSingle("testint", CIPTypeINT, 1)
	if err != nil {
		t.Errorf("problem reading. %v", err)
	}
	log.Printf("value: %v\n", value)
}

func TestErrs(t *testing.T) {
	e := fmt.Errorf("this is an error %v", 3)
	err := net.OpError{Err: e}

	if errors.Is(&err, &net.OpError{}) {
		log.Printf("ok")
	}

}

// these are the four example structs defined in 1756-PM020H-EN-P page 61
type test_STRUCT_A struct {
	Limits uint32 // two bits packed here.
	Travel uint32
	Errors uint32 // actually a byte in the controller. But comes as a dint.
	Wear   float32
}

type test_STRUCT_B struct {
	PilotOn     uint32 // actually a bool
	HourlyCount [12]uint16
	Rate        float32
}

type test_STRUCT_C struct {
	HoursFull  uint32 // actually a bool
	Today      test_STRUCT_B
	SampleTime uint32
	Shipped    uint32
}

type test_STRUCT_D struct {
	MyInt   uint32 // actually a uint16
	MyFloat float32
	MyArray [8]test_STRUCT_C
	MyPID   float32
}
