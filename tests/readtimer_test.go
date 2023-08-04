package gologix_tests

import (
	"log"
	"testing"
	"time"

	"github.com/danomagnum/gologix"
)

func TestTimerRead(t *testing.T) {
	var tmr gologix.LogixTIMER

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
	// that after 500 ms we should be between 450 and 550 ms elapsed on the timer
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

	if tmr.ACC < 450 || tmr.ACC > 550 {
		t.Errorf("Expected ACC between 450 and 550 but got %d", tmr.ACC)
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

}
