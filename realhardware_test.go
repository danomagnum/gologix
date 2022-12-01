package gologix

import (
	"flag"
	"fmt"
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
	ReadAndPrint[int32](client, "TestDint") // 36

	// these two tests don't work yet
	ReadAndPrint[bool](client, "TestDint.0") // should be false
	ReadAndPrint[bool](client, "TestDint.2") // should be true

	ReadAndPrint[int32](client, "TestDintArr[0]")
	ReadAndPrint[int32](client, "TestDintArr[2]")
	ReadAndPrint[int32](client, "TestUDT.Field1")
	ReadAndPrint[int16](client, "TestInt")
	//v, err := client.read_single("TestDintArr[1]", CIPTypeDINT, 2)
	//if err != nil {
	//fmt.Printf("Problem with reading two elements of array. %v\n", err)
	//} else {
	//fmt.Printf("two element value: %v\n", v)
	//}

	v2, err := ReadArray[int32](client, "TestDintArr", 9)
	if err != nil {
		t.Errorf("Problem with reading two elements of array. %v\n", err)
	} else {
		fmt.Printf("two element value new method: %v\n", v2)
	}
	v3, err := ReadArray[TestUDT](client, "TestUDTArr[2]", 2)
	if err != nil {
		t.Errorf("Problem with reading two elements of array. %v\n", err)
	} else {
		fmt.Printf("two elements of UDT : %+v\n", v3)
	}
	test_strarr := false
	// string array read is untested:
	if test_strarr {
		v4, err := ReadArray[string](client, "TestStrArr", 3)
		if err != nil {
			t.Errorf("Problem with reading two elements of array. %v\n", err)
		} else {
			fmt.Printf("two element value new method: %v\n", v4)
		}
	}

	//ReadAndPrint[bool](client, "TestBool")
	//ReadAndPrint[float32](client, "TestReal")
	ReadAndPrint[string](client, "TestString")

	value, err := Read[TestUDT](client, "TestUDT")
	if err != nil {
		t.Errorf("Problem reading udt. %v\n", err)
	}
	fmt.Printf("UDT: %+v\n", value)

	//tags := []string{"TestInt", "TestReal"}
	tags := MultiReadStr{}

	err = client.read_multi(&tags, CIPTypeDWORD, 1)
	if err != nil {
		t.Errorf("Error reading multi. %v\n", err)
	}

	fmt.Printf("Values: %+v", tags)

}

func ReadAndPrint[T GoLogixTypes](client *Client, path string) {
	value, err := Read[T](client, path)
	if err != nil {
		log.Printf("Problem reading %s. %v", path, err)
		return
	}
	log.Printf("%s: %v", path, value)
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
