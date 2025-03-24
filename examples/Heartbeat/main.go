package main

import (
	"log"
	"time"

	"github.com/danomagnum/gologix"
)

// Demo program for monitoring a heartbeat tag named "TestHeartBeat" in the controller.
// This example demonstrates how to detect when a value stops changing, which can be used
// to determine if a remote process is still running.
func main() {
	var err error

	// setup the client.  If you need a different path you'll have to set that.
	client := gologix.NewClient("192.168.2.241")

	// for example, to have a controller on slot 1 instead of 0 you could do this
	//client.Path, err = gologix.Serialize(gologix.CIPPort{PortNo: 1}, gologix.CIPAddress(1))
	// or this
	// client.Path, err = gologix.ParsePath("1,1")

	// connect using parameters in the client struct
	err = client.Connect()
	if err != nil {
		log.Printf("Error opening client. %v", err)
		return
	}
	// setup a deffered disconnect.  If you don't disconnect you might have trouble reconnecting because
	// you won't have sent the close forward open.  You'll have to wait for the CIP connection to time out
	// if that happens (about a minute)
	defer client.Disconnect()

	// the client should be threadsafe so kicking off goroutines doing reads/writes in the background like this should be no issue.
	hbChan := Heartbeat[bool](client, "TestHeartBeat", time.Second, time.Second*10)

	for {
		hbStatusChange := <-hbChan

		if hbStatusChange {
			log.Printf("heartbeat now OK")
		} else {
			log.Printf("heartbeat now bad")
		}

	}

}

// monitor tag for changes in a background goroutine.
// poll at an interval of pollrate
// sends a single true on the output channel if the value starts to change
// sends a single false on the output channel if the value stops changing for timeout
func Heartbeat[T gologix.GoLogixTypes](client *gologix.Client, tag string, pollrate, timeout time.Duration) <-chan bool {

	hbstatus := make(chan bool)

	go func() {
		// heart beat value variables
		var lastHB T
		var newHB T

		// hb last change reference
		lastHB_Time := time.Now()

		// hb status
		ok := false

		ticker := time.NewTicker(pollrate)
		for range ticker.C {
			err := client.Read(tag, &newHB)
			if err != nil {
				// heartbeat read failed. send an edge triggered message about that.
				if ok {
					hbstatus <- false
					ok = false
				}
				continue
			}
			// if the value changes, update the timestamp and set to ok flag to true.
			if newHB != lastHB {
				lastHB_Time = time.Now()
				if !ok {
					hbstatus <- true
				}
				ok = true
			}

			// see if we've timed out.  Send an edge triggered message out.
			if time.Since(lastHB_Time) > timeout {
				if ok {
					hbstatus <- false
					ok = false
				}
			}
			lastHB = newHB
		}
	}()

	return hbstatus

}
