package main

import (
	"log"

	"github.com/danomagnum/gologix"
)

// Demo program for reading attributes from a logix PLC
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

	// the attributes of the Identity object are as follows:
	//
	// 1 - vendor ID
	// 2 - Device Type
	// 3 - Product Code
	// 4 - Revision
	// 5 - Status
	// 6 - Serial Number
	// 7 - Product Name
	item, err := client.GetAttrSingle(gologix.CipObject_Identity, 1, 1)
	if err != nil {
		log.Fatalf("problem reading attribute 1: %v", err)
	}

	// The GetAttrSingle returns an entire CIP Object (some attributes may have complex data responses)
	// so you have to parse them yourself.  For simple data types you can just call the appropriate
	// function on the item to get the value directly.

	// Instance 1 of the identity object is a 16 bit vendor ID
	vendor, err := item.Int16()
	if err != nil {
		// We could have an error here if there wasn't enough data returned for example.
		log.Fatalf("problem getting vendor ID from the response item: %v", err)
	}
	log.Printf("Vendor ID: %X", vendor)

	item, err = client.GetAttrList(gologix.CipObject_Identity, 1, 1, 2, 3, 4)
	if err != nil {
		log.Fatalf("problem reading list: %v", err)
	}

	// The GetAttrList returns an entire CIP Object as described in the function docs.
	// Since we have to know before hand what the types are for the attributes we are
	// reading, we can create a struct to parse the response all at once.

	type AttrResults struct {
		Attr1_ID     uint16
		Attr1_Status uint16
		Attr1_Value  uint16

		Attr2_ID     uint16
		Attr2_Status uint16
		Attr2_Value  uint16

		Attr3_ID     uint16
		Attr3_Status uint16
		Attr3_Value  uint32

		Attr4_ID     uint16
		Attr4_Status uint16
		Attr4_Value  uint32
	}

	var ar AttrResults

	err = item.DeSerialize(&ar)
	if err != nil {
		log.Fatalf("Problem parsing attribute results: %v", err)
	}
	log.Printf("Attr results: %+v", ar)

}
