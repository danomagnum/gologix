package main

import (
	"log"

	"github.com/danomagnum/gologix"
)

// Demo program for reading elements from a DINT array named "TestDintArr" in the controller.
func main() {
	var err error

	// setup the client.  If you need a different path you'll have to set that.
	client := gologix.NewClient("192.168.2.241")

	// for example, to have a controller on slot 1 instead of 0 you could do this
	// client.Path, err = gologix.Serialize(gologix.CIPPort{PortNo: 1}, gologix.CIPAddress(1))
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

	// define a variable with a type that matches the tag you want to read.
	// In this case we're creating an int32 slice of length 5 to hold 5 elements from the DINT array.
	tag1 := make([]int32, 5)

	// call the read function.
	// This reads 5 elements from TestDintArr starting at index 2.
	// note that tag names are case insensitive.
	// also note that for atomic types and structs you need to use a pointer.
	// for slices you don't use a pointer.
	err = client.Read("TestDintArr[2]", tag1)
	if err != nil {
		log.Printf("error reading TestDintArr. %v", err)
	}
	// do whatever you want with the value
	log.Printf("tag1 has value %d", tag1)

}
