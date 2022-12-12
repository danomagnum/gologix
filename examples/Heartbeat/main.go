package main

import (
	"fmt"
	"gologix"
	"time"
)

// Demo program for readng an INT tag named "TestInt" in the controller.
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
		fmt.Printf("Error opening client. %v", err)
		return
	}
	// setup a disconnect.  If you don't disconnect you might have trouble reconnecting
	defer client.Disconnect()

	hbChan := Heartbeat[bool](client, "TestHeartBeat", time.Second, time.Second*10)

	for {
		hbStatusChange := <-hbChan

		if hbStatusChange {
			fmt.Printf("heartbeat now OK")
		} else {
			fmt.Printf("heartbeat now bad")
		}

	}

}

// monitor tag for changes.
// poll at an interval of pollrate
// sends a single true on the output channel if the value starts to change
// sends a single false on the output channel if the value stops changing for timeout
func Heartbeat[T gologix.GoLogixTypes](client *gologix.Client, tag string, pollrate, timeout time.Duration) <-chan bool {

	hbstatus := make(chan bool)

	go func() {
		// initialize heart beat variable
		var lastHB T
		var newHB T
		lastHB_Time := time.Now()
		ok := false

		// poll rate for the heartbeat tag
		ticker := time.NewTicker(pollrate)
		for range ticker.C {
			client.Read(tag, &newHB)
			if newHB != lastHB {
				lastHB_Time = time.Now()
				if !ok {
					hbstatus <- true
				}
				ok = true
			}
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
