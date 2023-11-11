package main

import (
	"log"

	"github.com/danomagnum/gologix"
)

// Demo program for readng multiple tags at once where the tags are in a go map.
// the keys in the map are the tag names, and the values need to be the correct type
// for the tag.  The ReadMap function will update the values in the map to the current values
// in the controller.
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

	// define a struct where fields have the tag to read from the controller specified
	// note that tag names are case insensitive.
	m := make(map[string]any)
	m["TestInt"] = int16(0)
	m["TestDint"] = int32(0)
	m["TestDintArr[2]"] = make([]int32, 5)

	// call the read multi function with the structure passed in as a pointer.
	err = client.ReadMulti(m)
	if err != nil {
		log.Printf("error reading testint. %v", err)
	}
	// do whatever you want with the values
	log.Printf("multiread map has values %+v", m)

}
