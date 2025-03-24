package main

import (
	"log"

	"github.com/danomagnum/gologix"
)

// Demo program for listing all tags in a controller and reading their values.
// This example shows how to:
// - Connect to an Allen Bradley PLC
// - List all available tags in the controller
// - Read specific tags by name
// - Read all tags in a loop
func main() {
	var err error

	// setup the client with the PLC's IP address
	client := gologix.NewClient("192.168.2.241")

	// for example, to have a controller on slot 1 instead of 0 you could do this
	//client.Path, err = gologix.Serialize(gologix.CIPPort{PortNo: 1}, gologix.CIPAddress(1))
	// or this
	// client.Path, err = gologix.ParsePath("1,1")
	client.KeepAliveAutoStart = false

	// connect using parameters in the client struct
	err = client.Connect()
	if err != nil {
		log.Printf("Error opening client. %v", err)
		return
	}
	// setup a deferred disconnect.
	defer client.Disconnect()

	// update the client's list of tags.
	err = client.ListAllTags(0)
	if err != nil {
		log.Printf("Error getting tag list. %v", err)
		return
	}

	// Reading specific tags as examples
	var y int32
	err = client.Read("Program:gologix_tests.ReadDint", &y)
	if err != nil {
		log.Printf("Error reading tag. %v", err)
		return
	}
	var x float32
	err = client.Read("Program:gologix_tests.ReadReal", &x)
	if err != nil {
		log.Printf("Error reading tag. %v", err)
		return
	}

	log.Printf("Found %d tags.", len(client.KnownTags))
	// loop through the tag list and read all tags
	for tagname := range client.KnownTags {
		tag := client.KnownTags[tagname]
		log.Printf("%s: %v", tag.Name, tag.Info.Type)

		// Handle array tags by accessing first element
		qty := uint16(1)
		if tag.Info.Dimension1 != 0 {
			tagname = tagname + "[0]"
			x := tag.Info.Atomic()
			qty = uint16(tag.Info.Dimension1)
			_ = x
		}
		if tag.UDT == nil && !tag.Info.Atomic() {
			//log.Print("Not Atomic or UDT")
			continue
		}
		if tag.UDT != nil {
			log.Printf("%s size = %d", tag.Name, tag.UDT.Size())
		}

		// Read and display the value of each tag
		val, err := client.Read_single(tagname, tag.Info.Type, qty)
		if err != nil {
			log.Printf("Error!  Problem reading tag %s. %v", tagname, err)
			continue
		}
		log.Printf("     = %v", val)
	}

	log.Printf("Found %d tags.", len(client.KnownTags))
}
