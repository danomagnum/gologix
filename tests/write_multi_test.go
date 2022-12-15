package gologix_tests

import (
	"testing"

	"github.com/danomagnum/gologix"
)

func TestWriteMulti(t *testing.T) {
	client := gologix.NewClient("192.168.2.241")
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()

	want := TestUDT{Field1: 215, Field2: 25.1}

	err = client.Write("Program:gologix_tests.WriteUDTs[2]", want)
	if err != nil {
		t.Errorf("error writing. %v", err)
	}

	have := TestUDT{}
	err = client.Read("Program:gologix_tests.WriteUDTs[2]", &have)
	if err != nil {
		t.Errorf("error reading. %v", err)
	}

	if have != want {
		t.Errorf("have %v. Want %v", have, want)

	}

}
