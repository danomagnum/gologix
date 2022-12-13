package gologix_tests

import (
	"gologix"
	"testing"
)

func TestWriteMulti(t *testing.T) {
	client := gologix.NewClient("192.168.2.241")
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()

	udt := TestUDT{Field1: 215, Field2: 25.1}

	err = client.Write("WriteUDTs[2]", udt)
	if err != nil {
		t.Errorf("error writing. %v", err)
	}

}
