package gologix_tests

import (
	"gologix"
	"testing"
)

func TestWrite(t *testing.T) {
	client := gologix.NewClient("192.168.2.241")
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()

	write_and_check[int16](t, client, "writeint", 0, 10, 20)
	write_and_check[int32](t, client, "writedint", 0, 10, 20)
	//write_and_check[byte](t, client, "writesint", 0, 10, 20)
	write_and_check(t, client, "writebool", false, true, false, true, false)
	write_and_check[float32](t, client, "writereal", 0, 12.4, 5353.0281, 4)

	write_and_check[int32](t, client, "writeudt.field1", 0, 5, 281, 46)
	write_and_check[float32](t, client, "writeudt.field2", 0, 12.4, 5353.0281, 4)

	/*
		have := []int32{12, 34, 56, 78}
		err = client.Write("writedints", have)
		if err != nil {
			t.Errorf("failed to write dint array %v", err)
		}
	*/

	/*
		have := TestUDT{Field1: 5, Field2: 7.2}
		err = client.Write("TestUDT", have)
		if err != nil {
			t.Errorf("failed to write udt %v", err)
		}
	*/

}

func write_and_check[T gologix.ComparableGoLogixTypes](t *testing.T, client *gologix.Client, tag string, values ...T) {
	var err error
	var have T
	t.Run(tag, func(t *testing.T) {
		for _, want := range values {
			err = client.Write(tag, want)
			if err != nil {
				t.Errorf("problem writing. %v", err)
			}
			err = client.Read(tag, &have)
			if err != nil {
				t.Errorf("problem reading back. %v", err)
			}
			if have != want {
				t.Errorf("want %v. have %v", want, have)
			}

		}
	})

}
