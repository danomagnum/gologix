package gologix_tests

import (
	"testing"

	"github.com/danomagnum/gologix"
)

func TestWrite(t *testing.T) {
	client := gologix.NewClient("192.168.2.241")
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()

	write_and_check[int16](t, client, "program:gologix_tests.writeint", 0, 10, 20)
	write_and_check[int32](t, client, "program:gologix_tests.writedint", 0, 10, 20)
	//write_and_check[byte](t, client, "writesint", 0, 10, 20)
	write_and_check(t, client, "program:gologix_tests.writebool", false, true, false, true, false)
	write_and_check[float32](t, client, "program:gologix_tests.writereal", 0, 12.4, 5353.0281, 4)

	write_and_check[int32](t, client, "program:gologix_tests.writeudt.field1", 0, 5, 281, 46)
	write_and_check[float32](t, client, "program:gologix_tests.writeudt.field2", 0, 12.4, 5353.0281, 4)

}

func write_and_check[T gologix.GoLogixTypes](t *testing.T, client *gologix.Client, tag string, values ...T) {
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
