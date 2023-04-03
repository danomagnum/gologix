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
	defer func() {
		err := client.Disconnect()
		if err != nil {
			t.Errorf("problem disconnecting. %v", err)
		}
	}()
	tag := "Program:gologix_tests.ReadDints[0]"
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
	defer func() {
		err := client.Disconnect()
		if err != nil {
			t.Errorf("problem disconnecting. %v", err)
		}
	}()
	tag := "Program:gologix_tests.ReadUDTs[0]"
	have := TestUDT{}
	want := TestUDT{Field1: 20, Field2: 19.0}
	err = client.Read(tag, &have)
	if err != nil {
		t.Errorf("failed to read. %v", err)
	}
	if have != want {
		log.Printf("have: %+v, want: %+v", have, want)
	}
}
func TestReadNewUDTArr(t *testing.T) {
	client := gologix.NewClient("192.168.2.241")
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err := client.Disconnect()
		if err != nil {
			t.Errorf("problem disconnecting. %v", err)
		}
	}()
	tag := "Program:gologix_tests.ReadUDTs[0]"
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
	defer func() {
		err := client.Disconnect()
		if err != nil {
			t.Errorf("problem disconnecting. %v", err)
		}
	}()

	type udt2 struct {
		Field1 int32
		Flag1  bool
		Flag2  bool
		Field2 int32
	}

	//have, err := gologix.ReadPacked[udt2](client, "Program:gologix_tests.ReadUDT2")
	var have udt2
	err = client.Read("Program:gologix_tests.ReadUDT2", &have)
	if err != nil {
		t.Errorf("couldn't read. %v", err)
	}
	want := udt2{
		Field1: 654321,
		Flag1:  true,
		Flag2:  false,
		Field2: 123456,
	}
	if have != want {
		t.Errorf("have %v want %v", have, want)
	}

}

func TestReadNew(t *testing.T) {
	client := gologix.NewClient("192.168.2.241")
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err := client.Disconnect()
		if err != nil {
			t.Errorf("problem disconnecting. %v", err)
		}
	}()

	testReadNew(t, client, "Program:gologix_tests.ReadSint", byte(117))
	testReadNew(t, client, "Program:gologix_tests.ReadDint", int32(36))
	testReadNew(t, client, "Program:gologix_tests.ReadDint.0", false)
	testReadNew(t, client, "Program:gologix_tests.ReadDint.2", true)
	testReadNew(t, client, "Program:gologix_tests.ReadReal", float32(93.45))
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0]", int32(4351))
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].0", true)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].1", true)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].2", true)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].3", true)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].4", true)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].5", true)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].6", true)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].7", true)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].8", false)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].9", false)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].10", false)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].11", false)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].12", true)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].13", false)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].14", false)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[0].15", false)
	testReadNew(t, client, "Program:gologix_tests.ReadDints[2]", int32(4353))
	testReadNew(t, client, "Program:gologix_tests.ReadUDT.Field1", int32(85456))
	testReadNew(t, client, "Program:gologix_tests.ReadUDT.Field2", float32(123.456))
	testReadNew(t, client, "Program:gologix_tests.ReadUDTs[2].Field1", int32(16))
	testReadNew(t, client, "Program:gologix_tests.ReadUDTs[2].Field2", float32(15.0))
	testReadNew(t, client, "Program:gologix_tests.ReadString", "Something")

}

func testReadNew[T gologix.GoLogixTypes](t *testing.T, client *gologix.Client, tag string, want T) {

	t.Run(tag, func(t *testing.T) {
		//tag, want := "Program:gologix_tests:ReadInt", int16(999)
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
		TestSint          byte    `gologix:"Program:gologix_tests.ReadSint"`
		TestInt           int16   `gologix:"Program:gologix_tests.ReadInt"`
		TestDint          int32   `gologix:"Program:gologix_tests.ReadDint"`
		TestReal          float32 `gologix:"Program:gologix_tests.ReadReal"`
		TestDintArr0      int32   `gologix:"Program:gologix_tests.Readdints[0]"`
		TestDintArr0_0    bool    `gologix:"Program:gologix_tests.Readdints[0].0"`
		TestDintArr0_9    bool    `gologix:"Program:gologix_tests.Readdints[0].9"`
		TestDintArr2      int32   `gologix:"Program:gologix_tests.Readdints[2]"`
		TestUDTField1     int32   `gologix:"Program:gologix_tests.Readudt.field1"`
		TestUDTField2     float32 `gologix:"Program:gologix_tests.Readudt.field2"`
		TestUDTArr2Field1 int32   `gologix:"Program:gologix_tests.Readudts[2].field1"`
		TestUDTArr2Field2 float32 `gologix:"Program:gologix_tests.Readudts[2].field2"`
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
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err := client.Disconnect()
		if err != nil {
			t.Errorf("problem disconnecting. %v", err)
		}
	}()

	err = client.ReadMulti(&read)
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
	defer func() {
		err := client.Disconnect()
		if err != nil {
			t.Errorf("problem disconnecting. %v", err)
		}
	}()
	log.Println("sleeping for 2 minutes. to let the comms timeout")
	for t := 0; t < 24; t++ {
		log.Printf("%d\n", t*5)
		time.Sleep(time.Second * 5)
	}
	//client.Conn.Close()
	//log.Println("\nsleep complete")
	var value int16
	err = client.Read("Program:gologix_tests.Readint", &value)
	if err != nil {
		t.Errorf("problem reading. %v", err)
	}
	log.Printf("value: %v\n", value)
}

// Right now this test will fail, but eventually the read function will be updated to do multiple requests and merge the
// results back together at which point it will pass.
func TestReadTooManyTags(t *testing.T) {
	client := gologix.NewClient("192.168.2.241")
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err := client.Disconnect()
		if err != nil {
			t.Errorf("problem disconnecting. %v", err)
		}
	}()
	tag := "Program:gologix_tests.ReadDints"

	tags := make([]string, 0)
	types := make([]gologix.CIPType, 0)

	tagcount := 4000

	for i := 0; i < tagcount; i++ {
		tags = append(tags, fmt.Sprintf("%s[%d]", tag, i))
		types = append(types, gologix.CIPTypeDINT)
	}

	_, err = client.ReadList(tags, types)
	if err == nil {
		t.Error("Should have failed but didn't.")
		return
	}
}
