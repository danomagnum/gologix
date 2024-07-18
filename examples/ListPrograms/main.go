package main

import (
	"fmt"
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
	client.KeepAliveAutoStart = false

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

	err = client.ListAllPrograms()
	if err != nil {
		log.Printf("Error getting program list. %v", err)
		return
	}

	for _, p := range client.KnownPrograms {
		if p.Name == "gologix_tests" {
			client.ListSubTags(p, 1)
		}
	}

	fmt.Printf("Found %d tags.\n", len(client.KnownTags))

}
