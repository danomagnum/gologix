package main

import (
	"fmt"
	"log"

	"github.com/danomagnum/gologix"
)

// Demo program for listing programs and tags from a PLC.
func main() {
	var err error

	// Setup the client to connect to a PLC with the specified IP address
	client := gologix.NewClient("192.168.2.241")

	// For example, to have a controller on slot 1 instead of 0 you could do this:
	//client.Path, err = gologix.Serialize(gologix.CIPPort{PortNo: 1}, gologix.CIPAddress(1))
	// or this:
	// client.Path, err = gologix.ParsePath("1,1")
	client.KeepAliveAutoStart = false

	// Connect to the PLC using parameters in the client struct
	err = client.Connect()
	if err != nil {
		log.Printf("Error opening client. %v", err)
		return
	}
	// Setup a deferred disconnect. If you don't disconnect properly, you might have trouble reconnecting
	// as the CIP connection will need to time out (about a minute)
	defer client.Disconnect()

	// List all programs on the PLC
	err = client.ListAllPrograms()
	if err != nil {
		log.Printf("Error getting program list. %v", err)
		return
	}

	// For a specific program, list its subtags with a depth of 1
	for _, p := range client.KnownPrograms {
		if p.Name == "gologix_tests" {
			client.ListSubTags(p, 1)
		}
	}

	// Display the number of tags found
	fmt.Printf("Found %d tags.\n", len(client.KnownTags))
}
