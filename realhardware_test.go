package gologix

import (
	"flag"
	"fmt"
	"log"
	"testing"
)

func TestRealHardware(t *testing.T) {
	flag.Parse()
	plc := &PLC{IPAddress: "192.168.2.241"}
	plc.Connect()
	defer plc.Disconnect()
	//plc.ReadAll(1)
	//plc.read_single("program:Shed.Temp1", CIPTypeREAL, 1)
	//ReadAndPrint[float32](plc, "program:Shed.Temp1")
	ReadAndPrint[int32](plc, "TestDint") // 36

	// these two tests don't work yet
	ReadAndPrint[bool](plc, "TestDint.0") // should be false
	ReadAndPrint[bool](plc, "TestDint.2") // should be true

	ReadAndPrint[int32](plc, "TestDintArr[0]")
	ReadAndPrint[int32](plc, "TestDintArr[2]")
	ReadAndPrint[int32](plc, "TestUDT.Field1")
	ReadAndPrint[int16](plc, "TestInt")
	//v, err := plc.read_single("TestDintArr[1]", CIPTypeDINT, 2)
	//if err != nil {
	//fmt.Printf("Problem with reading two elements of array. %v\n", err)
	//} else {
	//fmt.Printf("two element value: %v\n", v)
	//}

	v2, err := ReadArray[int32](plc, "TestDintArr", 9)
	if err != nil {
		t.Errorf("Problem with reading two elements of array. %v\n", err)
	} else {
		fmt.Printf("two element value new method: %v\n", v2)
	}
	v3, err := ReadArray[TestUDT](plc, "TestUDTArr[2]", 2)
	if err != nil {
		t.Errorf("Problem with reading two elements of array. %v\n", err)
	} else {
		fmt.Printf("two elements of UDT : %+v\n", v3)
	}
	test_strarr := false
	// string array read is untested:
	if test_strarr {
		v4, err := ReadArray[string](plc, "TestStrArr", 3)
		if err != nil {
			t.Errorf("Problem with reading two elements of array. %v\n", err)
		} else {
			fmt.Printf("two element value new method: %v\n", v4)
		}
	}

	//ReadAndPrint[bool](plc, "TestBool")
	//ReadAndPrint[float32](plc, "TestReal")
	ReadAndPrint[string](plc, "TestString")

	value, err := Read[TestUDT](plc, "TestUDT")
	if err != nil {
		t.Errorf("Problem reading udt. %v\n", err)
	}
	fmt.Printf("UDT: %+v\n", value)

	//tags := []string{"TestInt", "TestReal"}
	tags := MultiReadStr{}

	err = plc.read_multi(&tags, CIPTypeDWORD, 1)
	if err != nil {
		t.Errorf("Error reading multi. %v\n", err)
	}

	fmt.Printf("Values: %+v", tags)

}

func ReadAndPrint[T GoLogixTypes](plc *PLC, path string) {
	value, err := Read[T](plc, path)
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
