package gologix_tests

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/danomagnum/gologix"
)

func TestServer(t *testing.T) {
	path := "1,0"

	r := gologix.PathRouter{}

	mapTagProvider := gologix.MapTagProvider{
		Data: make(map[string]any),
	}

	path1, err := gologix.ParsePath(path)
	if err != nil {
		log.Printf("problem parsing path. %v", err)
		os.Exit(1)
	}

	r.Handle(path1.Bytes(), &mapTagProvider)

	s := gologix.NewServer(&r)

	err = createTags(&mapTagProvider)
	if err != nil {
		log.Fatalf("Failed to create tags: %v", err)
	}

	go func(s *gologix.Server) {
		err := s.Serve()
		if err != nil {
			log.Fatalf("Failed to run server: %v", err)
			panic(err)
		}
	}(s)
	time.Sleep(1 * time.Second)

	client := gologix.NewClient("127.0.0.1")
	err = client.Connect()
	if err != nil {
		t.Errorf("problem connecting to server. %v", err)
		return
	}
	defer client.Disconnect()

	var x bool
	err = client.Read("Bool1", &x)
	if err != nil {
		t.Errorf("problem reading Bool1. %v", err)
	}
	if x != true {
		t.Errorf("Expected Bool1 to be %v but got %v", true, x)
	}
	err = client.Read("Bool0", &x)
	if err != nil {
		t.Errorf("problem reading Bool0. %v", err)
	}
	if x != false {
		t.Errorf("Expected Bool0 to be %v but got %v", false, x)
	}

	var testByte byte
	err = client.Read("testbyte", &testByte)
	if err != nil {
		t.Errorf("problem reading testbyte. %v", err)
	}
	if testByte != 0x01 {
		t.Errorf("Expected testbyte to be %v but got %v", 0x01, testByte)
	}

	var testInt32 int32
	err = client.Read("testint32", &testInt32)
	if err != nil {
		t.Errorf("problem reading testint32. %v", err)
	}
	if testInt32 != 12345 {
		t.Errorf("Expected testint32 to be %v but got %v", 12345, testInt32)
	}

	var testString string
	err = client.Read("teststring", &testString)
	if err != nil {
		t.Errorf("problem reading teststring. %v", err)
	}
	if testString != "Hello World" {
		t.Errorf("Expected teststring to be %q but got %q", "Hello World", testString)
	}

	var testFloat64 float64
	err = client.Read("testfloat64", &testFloat64)
	if err != nil {
		t.Errorf("problem reading testfloat64. %v", err)
	}
	if testFloat64 != 10238.21 {
		t.Errorf("Expected testfloat64 to be %v but got %v", 10238.21, testFloat64)
	}

	testBoolArray := make([]bool, 10)
	err = client.Read("testboolarray", &testBoolArray)
	if err != nil {
		t.Errorf("problem reading testboolarray. %v", err)
	}
	/*
		expectedBoolArray := []bool{true, false, true, false, true, false, true, false, true, false}
		if len(testBoolArray) != len(expectedBoolArray) || !equalBoolSlices(testBoolArray, expectedBoolArray) {
			t.Errorf("Expected testboolarray to be %v but got %v", expectedBoolArray, testBoolArray)
		}
	*/

}

func createTags(mtp *gologix.MapTagProvider) error {

	err := mtp.TagWrite("Bool1", true)
	if err != nil {
		return err
	}
	err = mtp.TagWrite("Bool0", false)
	if err != nil {
		return err
	}

	err = mtp.TagWrite("testbyte", byte(0x01))
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testint32", int32(12345))
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testint64", int64(12345678))
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testdint", int32(12))
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testint8", int8(-16))
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testint", int16(3))
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testuint64", uint64(1234567))
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testuint32", uint32(1234))
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testuint", uint16(123))
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testfloat32", float32(543.21))
	if err != nil {
		return err
	}
	err = mtp.TagWrite("teststring", "Hello World")
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testfloat64", float64(10238.21))
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testboolarray", []bool{true, false, true, false, true, false, true, false, true, false})
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testbytearray", []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x01, 0x02, 0x03})
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testint8array", []int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testint16array", []int16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testint32array", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testint64array", []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testuint8array", []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testuint16array", []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testuint32array", []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testuint64array", []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testfloat32array", []float32{1.234, 2.567, 3.172, 4.234, 5.342, 6.31, 7.521, 8.42, 9.23, 10.123})
	if err != nil {
		return err
	}
	err = mtp.TagWrite("testfloat64array", []float64{10.123, 2.321, 3.313, 4.543, 5.123, 6.12, 7.423, 8.32, 9.12, 10.64})
	if err != nil {
		return err
	}
	err = mtp.TagWrite("teststringarray", []string{"Hello1", "World1", "Hello2", "World2", "Hello3", "World3", "Hello4", "World4", "Hello5", "World5"})
	if err != nil {
		return err
	}

	return nil
}
