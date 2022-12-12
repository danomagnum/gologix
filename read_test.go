package gologix

import (
	"errors"
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

func TestReadArrNew(t *testing.T) {
	client := &Client{IPAddress: "192.168.2.241"}
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()
	tag := "TestDintArr[0]"
	have := make([]int32, 5)
	want := []int32{4351, 4352, 4353, 4354, 4355}

	err = client.Read(tag, have)
	if err != nil {
		t.Errorf("Problem reading %s. %v", tag, err)
		return
	}
	for i := range want {
		if have[i] != want[i] {
			t.Errorf("index %d wanted %v got %v", i, want[i], have[i])
		}
	}
}

func TestReadNewUDT(t *testing.T) {
	client := &Client{IPAddress: "192.168.2.241"}
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()
	tag := "TestUDTArr[0]"
	have := TestUDT{}
	want := TestUDT{Field1: 20, Field2: 19.0}
	err = client.Read(tag, &have)
	if err != nil {
		t.Errorf("failed to read. %v", err)
	}
	if have != want {
		fmt.Printf("have: %+v, want: %+v", have, want)
	}
}
func TestReadNewUDTArr(t *testing.T) {
	client := &Client{IPAddress: "192.168.2.241"}
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()
	tag := "TestUDTArr[0]"
	have := make([]TestUDT, 5)
	want := []TestUDT{
		{Field1: 20, Field2: 19.0},
		{Field1: 18, Field2: 17.0},
		{Field1: 16, Field2: 15.0},
		{Field1: 14, Field2: 13.0},
		{Field1: 12, Field2: 11.0},
	}
	err = client.Read(tag, have)
	if err != nil {
		t.Errorf("failed to read. %v", err)
	}
	if len(have) != len(want) {
		t.Errorf("didn't get the right number of elements. got %v wanted %v", len(have), len(want))
	}
	for i := range want {
		if have[i] != want[i] {
			t.Errorf("have: %+v, want: %+v\n", have[i], want[i])
		}
	}
}

func TestReadNew(t *testing.T) {
	client := &Client{IPAddress: "192.168.2.241"}
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()

	testReadNew(t, client, "TestSint", byte(117))
	testReadNew(t, client, "TestDint", int32(36))
	testReadNew(t, client, "TestDint.0", false)
	testReadNew(t, client, "TestDint.2", true)
	testReadNew(t, client, "TestReal", float32(93.45))
	testReadNew(t, client, "TestDintArr[0]", int32(4351))
	testReadNew(t, client, "TestDintArr[0].0", true)
	testReadNew(t, client, "TestDintArr[0].1", true)
	testReadNew(t, client, "TestDintArr[0].2", true)
	testReadNew(t, client, "TestDintArr[0].3", true)
	testReadNew(t, client, "TestDintArr[0].4", true)
	testReadNew(t, client, "TestDintArr[0].5", true)
	testReadNew(t, client, "TestDintArr[0].6", true)
	testReadNew(t, client, "TestDintArr[0].7", true)
	testReadNew(t, client, "TestDintArr[0].8", false)
	testReadNew(t, client, "TestDintArr[0].9", false)
	testReadNew(t, client, "TestDintArr[0].10", false)
	testReadNew(t, client, "TestDintArr[0].11", false)
	testReadNew(t, client, "TestDintArr[0].12", true)
	testReadNew(t, client, "TestDintArr[0].13", false)
	testReadNew(t, client, "TestDintArr[0].14", false)
	testReadNew(t, client, "TestDintArr[0].15", false)
	testReadNew(t, client, "TestDintArr[2]", int32(4353))
	testReadNew(t, client, "TestUDT.Field1", int32(85456))
	testReadNew(t, client, "TestUDT.Field2", float32(123.456))
	testReadNew(t, client, "TestUDTArr[2].Field1", int32(16))
	testReadNew(t, client, "TestUDTArr[2].Field2", float32(15.0))

}

func testReadNew[T ComparableGoLogixTypes](t *testing.T, client *Client, tag string, want T) {

	t.Run(tag, func(t *testing.T) {
		//tag, want := "TestInt", int16(999)
		var have T

		err := client.Read(tag, &have)
		if err != nil {
			t.Errorf("Problem reading %s. %v", tag, err)
			return
		}
		if have != want {
			t.Errorf("wanted %v got %v", want, have)
		}
	},
	)
}

func TestReadSingle(t *testing.T) {
	var tests = []struct {
		tag   string
		wants any
	}{
		{"TestInt", int16(999)},
		{"TestSint", byte(117)},
		{"TestDint", int32(36)},
		{"TestDint.0", false},
		{"TestDint.2", true},
		{"TestReal", float32(93.45)},
		{"TestDintArr[0]", int32(4351)},
		{"TestDintArr[0].0", true},
		{"TestDintArr[0].1", true},
		{"TestDintArr[0].2", true},
		{"TestDintArr[0].3", true},
		{"TestDintArr[0].4", true},
		{"TestDintArr[0].5", true},
		{"TestDintArr[0].6", true},
		{"TestDintArr[0].7", true},
		{"TestDintArr[0].8", false},
		{"TestDintArr[0].9", false},
		{"TestDintArr[0].10", false},
		{"TestDintArr[0].11", false},
		{"TestDintArr[0].12", true},
		{"TestDintArr[0].13", false},
		{"TestDintArr[0].14", false},
		{"TestDintArr[0].15", false},
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
		TestDintArr0_0    bool    `gologix:"testdintarr[0].0"`
		TestDintArr0_9    bool    `gologix:"testdintarr[0].9"`
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
		TestDintArr0_0:    true,
		TestDintArr0_9:    false,
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
	t.Skip("requires timeout that is too long")

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
	var value int16
	err = client.Read("testint", &value)
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
