package gologix_tests

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/danomagnum/gologix"
	"github.com/danomagnum/gologix/lgxtypes"
)

func TestControl(t *testing.T) {
	var ctrl lgxtypes.CONTROL

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

	wants := []lgxtypes.CONTROL{
		{LEN: 8563, POS: 3324, EN: true, EU: false, DN: false, EM: false, ER: false, UL: false, IN: false, FD: false},
		{LEN: 0, POS: 0, EN: false, EU: true, DN: false, EM: false, ER: false, UL: false, IN: false, FD: false},
		{LEN: 0, POS: 0, EN: false, EU: false, DN: true, EM: false, ER: false, UL: false, IN: false, FD: false},
		{LEN: 0, POS: 0, EN: false, EU: false, DN: false, EM: true, ER: false, UL: false, IN: false, FD: false},
		{LEN: 0, POS: 0, EN: false, EU: false, DN: false, EM: false, ER: true, UL: false, IN: false, FD: false},
		{LEN: 0, POS: 0, EN: false, EU: false, DN: false, EM: false, ER: false, UL: true, IN: false, FD: false},
		{LEN: 0, POS: 0, EN: false, EU: false, DN: false, EM: false, ER: false, UL: false, IN: true, FD: false},
		{LEN: 0, POS: 0, EN: false, EU: false, DN: false, EM: false, ER: false, UL: false, IN: false, FD: true},
	}

	for i := range wants {
		//have, err := gologix.ReadPacked[udt2](client, "Program:gologix_tests.ReadUDT2")
		err = client.Read(fmt.Sprintf("Program:gologix_tests.TestControl[%d]", i), &ctrl)
		if err != nil {
			t.Errorf("problem reading %d. %v", i, err)
			return
		}

		compareControl(fmt.Sprintf("test %d", i), wants[i], ctrl, t)

		b := bytes.Buffer{}
		_ = gologix.Pack(&b, gologix.CIPPack{}, ctrl)
		var ctrl2 lgxtypes.CONTROL
		_, err = gologix.Unpack(&b, gologix.CIPPack{}, &ctrl2)
		if err != nil {
			t.Errorf("problem unpacking %d: %v", i, err)
		}

		compareControl(fmt.Sprintf("rebuild test %d", i), ctrl, ctrl2, t)
	}

}

func compareControl(name string, want, have lgxtypes.CONTROL, t *testing.T) {

	if have.LEN != want.LEN {
		t.Errorf("%s LEN mismatch. Have %d want %d", name, have.LEN, want.LEN)
	}
	if have.POS != want.POS {
		t.Errorf("%s POS mismatch. Have %d want %d", name, have.POS, want.POS)
	}
	if have.EN != want.EN {
		t.Errorf("%s EN mismatch. Have %v want %v", name, have.EN, want.EN)
	}
	if have.EU != want.EU {
		t.Errorf("%s EU mismatch. Have %v want %v", name, have.EU, want.EU)
	}
	if have.DN != want.DN {
		t.Errorf("%s DN mismatch. Have %v want %v", name, have.DN, want.DN)
	}
	if have.EM != want.EM {
		t.Errorf("%s EM mismatch. Have %v want %v", name, have.EM, want.EM)
	}
	if have.ER != want.ER {
		t.Errorf("%s ER mismatch. Have %v want %v", name, have.ER, want.ER)
	}
	if have.UL != want.UL {
		t.Errorf("%s UL mismatch. Have %v want %v", name, have.UL, want.UL)
	}
	if have.IN != want.IN {
		t.Errorf("%s IN mismatch. Have %v want %v", name, have.IN, want.IN)
	}
	if have.FD != want.FD {
		t.Errorf("%s FD mismatch. Have %v want %v", name, have.FD, want.FD)
	}

}
