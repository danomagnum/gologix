package gologix_tests

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/danomagnum/gologix"
)

func TestReadArrNew(t *testing.T) {
	client := gologix.NewClient("192.168.2.241")
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
	client := gologix.NewClient("192.168.2.241")
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
	client := gologix.NewClient("192.168.2.241")
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

func TestReadBoolPack(t *testing.T) {
	client := gologix.NewClient("192.168.2.241")
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()

	type udt2 struct {
		Field1 int32
		Flag1  bool
		Flag2  bool
		Field2 int32
	}

	s, err := gologix.ReadPacked[udt2](client, "TestUDT2")
	if err != nil {
		t.Errorf("couldn't read. %v", err)
	}
	/*
		b := make([]byte, 12)
		err = client.Read("TestUDT2", &b)
		if err != nil {
			t.Errorf("couldn't read as bytes. %v", err)
		}
		fmt.Printf("bytes: %v\n", b)
		gologix.Unpack(bytes.NewBuffer(b), gologix.CIPPack{}, &s)
	*/
	fmt.Printf("got %+v\n", s)

}

func TestReadNew(t *testing.T) {
	client := gologix.NewClient("192.168.2.241")
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
	testReadNew(t, client, "TestString", "Something")

}

func testReadNew[T gologix.GoLogixTypes](t *testing.T, client *gologix.Client, tag string, want T) {

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

	client := gologix.NewClient("192.168.2.241")
	client.Connect()
	defer client.Disconnect()

	err := client.ReadMulti(&read)
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

	client := gologix.NewClient("192.168.2.241")
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
