package main

import (
	"bytes"
	"log"
	"os"
	"time"

	"github.com/danomagnum/gologix"
)

// Demo program for testing PLC connection reliability by reading
// a tag at increasingly longer intervals until the connection drops.
func main() {
	var err error

	// Setup the client with the PLC IP address
	client := gologix.NewClient("192.168.2.244")
	client.Controller.Path = &bytes.Buffer{}

	// For example, to have a controller on slot 1 instead of 0 you could do this
	// client.Path, err = gologix.Serialize(gologix.CIPPort{PortNo: 1}, gologix.CIPAddress(1))
	// or this
	// client.Path, err = gologix.ParsePath("1,1")

	// Connect using parameters in the client struct
	err = client.Connect()
	if err != nil {
		log.Printf("Error opening client. %v", err)
		os.Exit(1)
	}
	// Setup a deferred disconnect. If you don't disconnect properly, you might have
	// trouble reconnecting because you won't have sent the close forward open.
	// You'll have to wait for the CIP connection to time out if that happens (about a minute)
	defer client.Disconnect()

	// Define a variable with a type that matches the tag you want to read.
	// In this case, we're using int32 for the 'inputs[0]' tag
	var tag1 int32

	// Initial read of the tag to verify connection
	err = client.Read("inputs[0]", &tag1)
	if err != nil {
		log.Printf("error reading inputs[0]. %v", err)
		os.Exit(1)
	}

	// Start with a delay of 287 seconds and increase by 1 second each iteration
	t := time.Second * 287
	for {
		time.Sleep(t)
		err = client.Read("inputs[0]", &tag1)
		if err != nil {
			log.Printf("error reading inputs[0] after %v seconds: %v", t, err)
			break
		}
		// Log the value and time interval
		log.Printf("after %v seconds tag1 has value %d", t, tag1)
		t += time.Second * 1
	}
}
