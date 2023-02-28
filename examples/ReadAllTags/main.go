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
		//log.Printf("%s: %v", tag.Name, tag.Info.Type)
		if tag.Name == "Program:MotionTest.Target3" ||
			tag.Name == "Program:Garage.Pressure_History" ||
			tag.Name == "Program:MainProgram.pBaseINTArray" ||
			tag.Name == "Program:gologix_tests.WriteSints" ||
			tag.Name == "Program:gologix_tests.ReadDints" {
			log.Print("this one")
		}
		if tagname == "local:1:i" {
			log.Print("should have one!")
		}
		//if !tag.Info.Atomic() {
		//log.Print("Not Atomic")
		//continue
		//}
		if tag.Info.Dimension1 != 0 {
			tagname = tagname + "[0]"
		}
		if tag.UDT == nil && !tag.Info.Atomic() {
			//log.Print("Not Atomic or UDT")
			continue
		}
		if tag.UDT != nil {
			log.Printf("%s size = %d", tag.Name, tag.UDT.Size())
		}

		if false {
			val, err := client.Read_single(tagname, tag.Info.Type, 1)
			if err != nil {
				log.Printf("Error!  Problem reading tag %s. %v", tagname, err)
				continue
			}
			log.Printf("     = %v", val)
		}
	}

}
