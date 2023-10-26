package gologix_tests

import (
	"flag"
	"log"
	"testing"

	"github.com/danomagnum/gologix"
	"github.com/danomagnum/gologix/ciptype"
)

func TestRealHardware(t *testing.T) {
	flag.Parse()
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
	//client.ReadAll(1)
	//client.read_single("program:Shed.Temp1", CIPTypeREAL, 1)
	//ReadAndPrint[float32](client, "program:Shed.Temp1")
	read[int32](t, client, "Program:gologix_tests.ReadDint") // 36

	// these two tests don't work yet
	read[bool](t, client, "Program:gologix_tests.ReadDint.0") // should be false
	read[bool](t, client, "Program:gologix_tests.ReadDint.2") // should be true

	read[int32](t, client, "Program:gologix_tests.ReadDints[0]")
	read[int32](t, client, "Program:gologix_tests.ReadDints[2]")
	read[int32](t, client, "Program:gologix_tests.ReadUDT.Field1")
	read[int16](t, client, "Program:gologix_tests.ReadInt")
	//v, err := client.read_single("Program:gologix_tests.ReadDintArr[1]", CIPTypeDINT, 2)
	//if err != nil {
	//log.Printf("Problem with reading two elements of array. %v\n", err)
	//} else {
	//log.Printf("two element value: %v\n", v)
	//}

	/*
		v2, err := readArray[int32](client, "Program:gologix_tests.ReadDintArr", 9)
		if err != nil {
			t.Errorf("Problem with reading two elements of array. %v\n", err)
		}
		for i := range v2 {
			want := int32(4351 + i)
			if v2[i] != want {
				t.Errorf("Problem with reading nine elements of array. got %v want %v\n", v2[i], want)
			}
		}
	*/
	/*
		v3, err := readArray[TestUDT](client, "Program:gologix_tests.ReadUDTArr[2]", 2)
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
			v4, err := readArray[string](client, "Program:gologix_tests.ReadStrArr", 3)
			if err != nil {
				t.Errorf("Problem with reading two elements of array. %v\n", err)
			} else {
				t.Logf("two element value new method: %v\n", v4)
			}
		}
	*/

	//ReadAndPrint[bool](client, "Program:gologix_tests.ReadBool")
	//ReadAndPrint[float32](client, "Program:gologix_tests.ReadReal")
	read[string](t, client, "Program:gologix_tests.ReadString")

	var ut TestUDT
	err = client.Read("Program:gologix_tests.ReadUDT", &ut)
	if err != nil {
		t.Errorf("Problem reading udt. %v\n", err)
	}

	//tags := []string{"Program:gologix_tests.ReadInt", "TestReal"}
	tags := MultiReadStr{}

	err = client.ReadMulti(&tags)
	if err != nil {
		t.Errorf("Error reading multi. %v\n", err)
	}

}

func read[T ciptype.GoLogixTypes](t *testing.T, client *gologix.Client, path string) {
	var have T
	err := client.Read(path, &have)
	if err != nil {
		t.Errorf("Problem reading %s. %v", path, err)
	}
	//t.Logf("%s: %v", path, value)
}

type MultiReadStr struct {
	TI int16   `gologix:"Program:gologix_tests.ReadInt"`
	TD int32   `gologix:"Program:gologix_tests.ReadDint"`
	TR float32 `gologix:"Program:gologix_tests.ReadReal"`
}

type TestUDT struct {
	Field1 int32
	Field2 float32
}

func TestReadKnown(t *testing.T) {

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

	err = client.ListAllTags(0)
	if err != nil {
		t.Error(err)
		return
	}

	log.Printf("Tags: %+v\n", client.KnownTags["program:gologix_tests.readint"])

	v := int16(0)
	err = client.Read("Program:gologix_tests.ReadInt", &v)

	if err != nil {
		t.Errorf("problem reading. %v", err)
	}

	if v != 999 {
		t.Errorf("problem reading 'TestInt'. got %v want %v", v, 999)

	}

}
