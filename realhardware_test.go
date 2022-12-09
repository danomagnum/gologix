package gologix

import (
	"flag"
	"log"
	"testing"
)

func TestRealHardware(t *testing.T) {
	flag.Parse()
	client := &Client{IPAddress: "192.168.2.241"}
	client.Connect()
	defer client.Disconnect()
	//client.ReadAll(1)
	//client.read_single("program:Shed.Temp1", CIPTypeREAL, 1)
	//ReadAndPrint[float32](client, "program:Shed.Temp1")
	ReadAndPrint[int32](t, client, "TestDint") // 36

	// these two tests don't work yet
	ReadAndPrint[bool](t, client, "TestDint.0") // should be false
	ReadAndPrint[bool](t, client, "TestDint.2") // should be true

	ReadAndPrint[int32](t, client, "TestDintArr[0]")
	ReadAndPrint[int32](t, client, "TestDintArr[2]")
	ReadAndPrint[int32](t, client, "TestUDT.Field1")
	ReadAndPrint[int16](t, client, "TestInt")
	//v, err := client.read_single("TestDintArr[1]", CIPTypeDINT, 2)
	//if err != nil {
	//log.Printf("Problem with reading two elements of array. %v\n", err)
	//} else {
	//log.Printf("two element value: %v\n", v)
	//}

	v2, err := ReadArray[int32](client, "TestDintArr", 9)
	if err != nil {
		t.Errorf("Problem with reading two elements of array. %v\n", err)
	}
	for i := range v2 {
		want := int32(4351 + i)
		if v2[i] != want {
			t.Errorf("Problem with reading nine elements of array. got %v want %v\n", v2[i], want)
		}
	}
	v3, err := ReadArray[TestUDT](client, "TestUDTArr[2]", 2)
	if err != nil {
		t.Errorf("Problem with reading two elements of array. %v\n", err)
	}
	want := TestUDT{16, 15.0}
	if v3[0] != want {
		t.Errorf("problem reading two elements of array. got %v want %v", v3[0], want)
	}
	want = TestUDT{14, 13.0}
	if v3[1] != want {
		t.Errorf("problem reading two elements of array. got %v want %v", v3[1], want)
	}
	test_strarr := false
	// string array read is untested:
	if test_strarr {
		v4, err := ReadArray[string](client, "TestStrArr", 3)
		if err != nil {
			t.Errorf("Problem with reading two elements of array. %v\n", err)
		} else {
			t.Logf("two element value new method: %v\n", v4)
		}
	}

	//ReadAndPrint[bool](client, "TestBool")
	//ReadAndPrint[float32](client, "TestReal")
	ReadAndPrint[string](t, client, "TestString")

	_, err = Read[TestUDT](client, "TestUDT")
	if err != nil {
		t.Errorf("Problem reading udt. %v\n", err)
	}

	//tags := []string{"TestInt", "TestReal"}
	tags := MultiReadStr{}

	err = client.ReadMulti(&tags, CIPTypeDWORD, 1)
	if err != nil {
		t.Errorf("Error reading multi. %v\n", err)
	}

}

func ReadAndPrint[T GoLogixTypes](t *testing.T, client *Client, path string) {
	_, err := Read[T](client, path)
	if err != nil {
		t.Errorf("Problem reading %s. %v", path, err)
	}
	//t.Logf("%s: %v", path, value)
}

type MultiReadStr struct {
	TI int16   `gologix:"TestInt"`
	TD int32   `gologix:"TestDint"`
	TR float32 `gologix:"TestReal"`
}

type TestUDT struct {
	Field1 int32
	Field2 float32
}

func TestReadKnown(t *testing.T) {

	client := &Client{IPAddress: "192.168.2.241"}
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()

	client.ListAllTags(0)

	log.Printf("Tags: %+v\n", client.KnownTags["testdintarr"])

	v, err := client.read_single("TestInt", CIPTypeINT, 1)

	if err != nil {
		t.Error(err)
	}

	if v.(int16) != 999 {
		t.Errorf("problem reading 'TestInt'. got %v want %v", v, 999)

	}

}
