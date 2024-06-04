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
	client.Controller.Path = &bytes.Buffer{}

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

	err = client.ListAllTags(0)
	if err != nil {
		log.Printf("Error reading tags: %v", err)
		return
	}
	log.Printf("All Tags:")
	for i := range client.KnownTags {
		tag := client.KnownTags[i]
		log.Printf("%v: / %v[%v]", tag.Name, tag.Info.Type, tag.Info.Dimension1)
	}

}
