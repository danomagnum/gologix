package main

import (
	"log"

	"github.com/danomagnum/gologix"
)

// Demo program for doing a generic CIP message.
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

	// for generic messages we need to create the cip path ourselves.  The serialize function can be used to do this.
	path, err := gologix.Serialize(gologix.CipObject_RunMode, gologix.CIPInstance(1))
	if err != nil {
		log.Printf("could not serialize path: %v", err)
		return
	}

	// This generic message would probably stop the controller, but you'd have to figure out how to elevate
	// the privileges associated with your connection first.  As it stands, you will probably get an 0x0F status code
	// and it won't do anything.
	resp, err := client.GenericCIPMessage(gologix.CIPService_Stop, path.Bytes(), []byte{})
	if err != nil {
		log.Printf("problem stopping PLC: %v", err)
		return
	}
	_ = resp

}
