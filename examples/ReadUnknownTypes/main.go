package main

import (
	"log"

	"github.com/danomagnum/gologix"
)

// Demo program for readng an INT tag named "TestInt" in the controller when you don't know the data type
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

	// define a variable of type any if you don't know the type of the tag.
	var tag1 any
	// call the read function.
	// note that tag names are case insensitive.
	// also note that for you need to use a pointer.
	err = client.Read("testint", &tag1)
	if err != nil {
		log.Printf("error reading testint. %v", err)
	}
	// do whatever you want with the value.  You'll probably have to type assert it to figure out what it is.
	switch x := tag1.(type) {
	case int16:
		log.Printf("tag1 has value %d", x)
	default:
		log.Printf("tag1 has type %T", x)
	}

	// If you try to read a UDT as an unknown type, you'll get a packed byte slice of the data back.
	// To read a UDT and get usable data you need to define a struct that matches the UDT definition in the controller.
	// See the simple read example for how to do that.
	var dat any
	err = client.Read("Program:Gologix_Tests.ReadUDT", &dat)
	if err != nil {
		log.Printf("error reading udt. %v", err)
	}

	log.Printf("dat has value %v", dat)

	// to read multiple unknown types at once you can use the ReadMap function

	// define a map of string to any.  The keys are the tag names and the values are nil since we don't know the types.
	// if you assigned a value to the map, it will be used as the type for the read.
	mr := make(map[string]any)
	mr["TestInt"] = nil
	mr["Program:Gologix_Tests.ReadUDT"] = nil

	err = client.ReadMap(mr)
	if err != nil {
		log.Printf("error reading map. %v", err)
	}
	// now you can use the values in the map.  The types are determined by the tag type on the PLC.
	log.Printf("TestInt has type %T and value %v", mr["TestInt"], mr["TestInt"])
	log.Printf("ReadUDT has type %T and value %v", mr["Program:Gologix_Tests.ReadUDT"], mr["Program:Gologix_Tests.ReadUDT"])

}
