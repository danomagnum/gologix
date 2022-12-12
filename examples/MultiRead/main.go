package main

import (
	"fmt"
	"gologix"
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

	// define a struct where fields have the tag to read from the controller specified
	// note that tag names are case insensitive.
	type multiread struct {
		TestInt  int16 `gologix:"TestInt"`
		TestDint int32 `gologix:"TestDint"`
	}
	var mr multiread

	// call the read multi function with the structure passed in as a pointer.
	err = client.ReadMulti(&mr)
	if err != nil {
		fmt.Printf("error reading testint. %v", err)
	}
	// do whatever you want with the values
	fmt.Printf("multiread struct has values %+v", mr)

}
