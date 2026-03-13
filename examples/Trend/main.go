package main

import (
	"fmt"
	"log"
	"time"

	"github.com/danomagnum/gologix"
)

// Demo program for setting up a PLC-side trend using a DataTable buffer.
//
// A DataTable buffer (CIP class 0xB2) is created on the PLC, tags are associated
// with it, configured with a sample rate, and then all tag values are read in a
// single request. This is the same mechanism RSLinx uses for grouped/trend reads.
//
// Three approaches are demonstrated:
//  1. AddTag      — add tags one at a time (separate 0x4E calls)
func main() {

	// Setup the client.
	client := gologix.NewClient("192.168.2.241")

	err := client.Connect()
	if err != nil {
		log.Printf("Error opening client. %v", err)
		return
	}
	defer client.Disconnect()

	// =========================================================================
	// 1. AddTag — add tags one at a time
	// =========================================================================
	// Each call sends a separate CIP 0x4E message to the PLC.

	fmt.Println("=== AddTag (one at a time) ===")

	type TrendData struct {
		CycleCount       int32   `gologix:"CycleCount"`
		CycleCountEighth float32 `gologix:"CycleFloatEighth"`
	}

	t, err := gologix.NewStructTrend[TrendData](client, time.Millisecond*10, 100)
	if err != nil {
		log.Fatalf("Error setting trend attributes: %v", err)
	}

	err = t.StartTrend()
	if err != nil {
		log.Fatalf("Error starting trend: %v", err)
	}

	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		dat, err := t.ReadAll()
		if err != nil {
			log.Fatalf("Error updating trend: %v", err)
		}
		fmt.Printf("Trend data: %+v\n", dat)
	}

	err = t.StopTrend()
	if err != nil {
		log.Fatalf("Error stopping trend: %v", err)
	}

	err = t.StopTrend()
	if err != nil {
		log.Fatalf("Error stopping trend: %v", err)
	}
}
