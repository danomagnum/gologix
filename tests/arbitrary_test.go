package gologix_tests

import (
	"testing"

	"github.com/danomagnum/gologix"
)

func TestArb(t *testing.T) {
	path, err := gologix.ParsePath("1,0")
	if err != nil {
		t.Errorf("problem parsing path. %v", err)
	}

	client := gologix.NewClient("192.168.2.241")
	err = client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()

	ioi, err := client.NewIOI("program:gologix_tests.ReadInt", gologix.CIPTypeINT)
	if err != nil {
		t.Errorf("problem parsing ioi. %v", err)
	}

	client.ArbitraryMessage(0x4D, path, ioi)

}
