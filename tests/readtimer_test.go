package gologix_tests

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/danomagnum/gologix"
	"github.com/danomagnum/gologix/lgxtypes"
)

func TestTimerRead(t *testing.T) {
	var tmr lgxtypes.TIMER

	tc := getTestConfig()
	client := gologix.NewClient(tc.PLC_Address)
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

	err = client.Write("Program:gologix_tests.trigger_Timer", false)
	if err != nil {
		t.Errorf("failed to turn off the timer before starting: %v", err)
		return
	}
	//have, err := gologix.ReadPacked[udt2](client, "Program:gologix_tests.ReadUDT2")
	err = client.Read("Program:gologix_tests.TestTimer", &tmr)
	log.Printf("timer 1: %+v", tmr)

	if err != nil {
		t.Errorf("problem reading timer data: %v", err)
		return
	}

	if tmr.PRE != 2345 {
		t.Errorf("Expected preset of 2,345 but got %d ", tmr.PRE)
	}

	if tmr.ACC != 0 {
		t.Errorf("Expected ACC of 0 but got %d", tmr.ACC)
	}

	if tmr.DN {
		t.Error("Expected timer !DN")
	}

	if tmr.EN {
		t.Error("Expected timer !EN")
	}

	if tmr.TT {
		t.Error("Expected timer !TT")
	}

	err = client.Write("Program:gologix_tests.trigger_Timer", true)
	if err != nil {
		t.Errorf("problem starting timer: %v", err)
	}

	// the task this timer lives in is set to a 50 ms rate so we should expect
	// that after 500 ms we should be between 449 and 551 ms elapsed on the timer
	// (accounting for minimal networking latency)
	time.Sleep(time.Millisecond * 500)

	err = client.Read("Program:gologix_tests.TestTimer", &tmr)
	log.Printf("timer 2: %+v", tmr)
	if err != nil {
		t.Errorf("problem reading timer data: %v", err)
		return
	}

	if tmr.PRE != 2345 {
		t.Errorf("Expected preset of 2,345 but got %d ", tmr.PRE)
	}

	if tmr.ACC < 449 || tmr.ACC > 551 {
		t.Errorf("Expected ACC between 449 and 551 but got %d", tmr.ACC)
	}

	if tmr.DN {
		t.Error("Expected timer !DN")
	}
	err = client.Write("Program:gologix_tests.trigger_Timer", true)
	if err != nil {
		t.Errorf("problem starting timer: %v", err)
	}

	time.Sleep(time.Millisecond * 2000)

	err = client.Read("Program:gologix_tests.TestTimer", &tmr)
	log.Printf("timer 3: %+v", tmr)
	if err != nil {
		t.Errorf("problem reading timer data: %v", err)
		return
	}

	if tmr.PRE != 2345 {
		t.Errorf("Expected preset of 2,345 but got %d ", tmr.PRE)
	}

	if tmr.ACC <= 2345 {
		t.Errorf("Expected ACC at least 2,345 but got %d", tmr.ACC)
	}

	if !tmr.DN {
		t.Error("Expected timer DN")
	}

	if !tmr.EN {
		t.Error("Expected timer EN")
	}

	if tmr.TT {
		t.Error("Expected timer !TT")
	}

	err = client.Write("Program:gologix_tests.trigger_Timer", false)
	if err != nil {
		t.Errorf("problem resetting timer: %v", err)
	}

	// make sure we can go the other way and recover it.
	b := bytes.Buffer{}
	_, err = gologix.Pack(&b, tmr)
	if err != nil {
		t.Errorf("problem packing data: %v", err)
	}
	var tmr2 lgxtypes.TIMER
	_, err = gologix.Unpack(&b, &tmr2)
	if err != nil {
		t.Errorf("problem unpacking timer: %v", err)
	}

	if tmr.ACC != tmr2.ACC {
		t.Errorf("ACC didn't recover properly.  %d != %d", tmr.ACC, tmr2.ACC)
	}

	if tmr.PRE != tmr2.PRE {
		t.Errorf("PRE didn't recover properly.  %d != %d", tmr.PRE, tmr2.PRE)
	}

	if tmr.DN != tmr2.DN {
		t.Errorf("DN didn't recover properly.  %v != %v", tmr.DN, tmr2.DN)
	}

	if tmr.TT != tmr2.TT {
		t.Errorf("TT didn't recover properly.  %v != %v", tmr.TT, tmr2.TT)
	}

	if tmr.EN != tmr2.EN {
		t.Errorf("EN didn't recover properly.  %v != %v", tmr.EN, tmr2.EN)
	}

}

func TestTimerStructRead(t *testing.T) {

	tc := getTestConfig()
	client := gologix.NewClient(tc.PLC_Address)
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

	x := struct {
		Field0 int32
		Flag1  bool
		Flag2  bool
		Timer  lgxtypes.TIMER
		Field1 int32
	}{}

	//have, err := gologix.ReadPacked[udt2](client, "Program:gologix_tests.ReadUDT2")
	err = client.Read("Program:gologix_tests.TestTimerStruct", &x)
	if err != nil {
		t.Errorf("problem reading timer data: %v", err)
		return
	}

	if x.Timer.PRE != 8765 {
		t.Errorf("Expected preset of 8765 but got %d ", x.Timer.PRE)
	}

	if x.Field0 != 44444 {
		t.Errorf("Expected field0 of 44444 but got %d", x.Field0)
	}

	if x.Field1 != 55555 {
		t.Errorf("Expected field1 of 55555 but got %d", x.Field1)
	}

}
