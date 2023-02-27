package main

import (
	"log"

	"github.com/danomagnum/gologix"
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
		log.Printf("Error opening client. %v", err)
		return
	}
	// setup a deffered disconnect.
	defer client.Disconnect()

	// update the client's list of tags.
	err = client.ListAllTags(0)

	if err != nil {
		log.Printf("Error getting tag list. %v", err)
		return
	}

	log.Printf("Found %d tags.", len(client.KnownTags))
	// list through the tag list.
	for tagname := range client.KnownTags {
		tag := client.KnownTags[tagname]
		log.Printf("%s: %v", tag.Name, tag.Info.Type)
		if !tag.Info.Atomic() {
			log.Print("Not Atomic")
			continue
		}
		val, err := client.Read_single(tagname, tag.Info.Type, 1)
		if err != nil {
			log.Printf("Error!  Problem reading tag %s. %v", tagname, err)
			continue
		}
		log.Printf("     = %v", val)
	}

}
