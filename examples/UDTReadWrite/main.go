package main

import (
	"log"

	"github.com/danomagnum/gologix"
)

// Demo program for readng a UDT tag named "TestUDT" in the controller.
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

	// define a struct that matches the tag you want to read.
	// the struct must have the same fields and types as the UDT in the controller.
	// the name of the struct MUST match the name of the UDT in the controller.
	// the name of any sub-structures must also match the name of the UDT in the controller.
	// field names don't matter and don't have to match.  They do have to be exported though so the
	// reflect package can access them to set the values.
	type myUDT struct {
		Field1 int32
		Field2 float32
	}
	var tag1 myUDT
	// call the read function.
	// note that tag names are case insensitive.
	// also note that you need to use a pointer.
	err = client.Read("TestUDT", &tag1)
	if err != nil {
		log.Printf("error reading testudt. %v", err)
	}
	// do whatever you want with the value
	log.Printf("tag1 has value %+v", tag1)

	//
	// Writing example.
	//
	tag1.Field1 = 5
	tag1.Field2 = 12.4

	err = client.Write("TestUDT", tag1)
	if err != nil {
		log.Printf("error writing testudt. %v", err)
	}

}
