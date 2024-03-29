package main

import (
	"bytes"
	"log"

	"github.com/danomagnum/gologix"
)

// Demo program for readng an INT tag named "TestInt" in the controller.
func main() {
	var err error

	// setup the client.  If you need a different path you'll have to set that.
	client := gologix.NewClient("192.168.2.244")
	// micro8xx use no path.  So an empty buffer will give us that.
	client.Path = &bytes.Buffer{}

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

	// define a variable with a type that matches the tag you want to read.  In this case it is an INT so
	// int16 or uint16 will work.
	input_dat := make([]int32, 8)

	// call the read function.
	// note that tag names are case insensitive.
	// also note that for atomic types and structs you need to use a pointer.
	// for slices you don't use a pointer.
	//
	// As far as I can tell you can't read program scope tags
	err = client.Read("inputs", input_dat)
	if err != nil {
		log.Printf("error reading 'input' tag. %v\n", err)
	}
	// do whatever you want with the value
	log.Printf("input_dat has value %d\n", input_dat)

}
