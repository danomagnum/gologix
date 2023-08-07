package gologix_tests

import (
	"bytes"
	"testing"

	"github.com/danomagnum/gologix"
	"github.com/danomagnum/gologix/lgxtypes"
)

func TestCounterRead(t *testing.T) {
	var cnt lgxtypes.COUNTER

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

	//have, err := gologix.ReadPacked[udt2](client, "Program:gologix_tests.ReadUDT2")
	err = client.Read("Program:gologix_tests.TestCounter", &cnt)

	if err != nil {
		t.Errorf("problem reading counter data: %v", err)
		return
	}

	if cnt.PRE != 562855 {
		t.Errorf("Expected preset of 2,345 but got %d ", cnt.PRE)
	}

	if cnt.ACC != 632 {
		t.Errorf("Expected ACC of 0 but got %d", cnt.ACC)
	}

	if cnt.DN {
		t.Error("Expected counter !DN")
	}

	if cnt.CU {
		t.Error("Expected counter !CU")
	}

	if cnt.CD {
		t.Error("Expected counter !CD")
	}

	// make sure we can go the other way and recover it.
	b := bytes.Buffer{}
	_ = gologix.Pack(&b, gologix.CIPPack{}, cnt)
	var cnt2 lgxtypes.COUNTER
	_, err = gologix.Unpack(&b, gologix.CIPPack{}, &cnt2)
	if err != nil {
		t.Errorf("problem unpacking timer: %v", err)
	}

	if cnt.ACC != cnt2.ACC {
		t.Errorf("ACC didn't recover properly.  %d != %d", cnt.ACC, cnt2.ACC)
	}

	if cnt.PRE != cnt2.PRE {
		t.Errorf("PRE didn't recover properly.  %d != %d", cnt.PRE, cnt2.PRE)
	}

	if cnt.DN != cnt2.DN {
		t.Errorf("DN didn't recover properly.  %v != %v", cnt.DN, cnt2.DN)
	}

	if cnt.CU != cnt2.CU {
		t.Errorf("CU didn't recover properly.  %v != %v", cnt.CU, cnt2.CU)
	}

	if cnt.CD != cnt2.CD {
		t.Errorf("CD didn't recover properly.  %v != %v", cnt.CD, cnt2.CD)
	}

}
