package main

import (
	"log"

	"github.com/danomagnum/gologix"
	"github.com/danomagnum/gologix/canned"
)

func main() {
	// This is an example of using functions in the canned package to do common tasks.
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

	forceStatus, err := canned.GetForces(client)
	if err != nil {
		log.Printf("error reading forces. %v", err)
		return
	}
	if forceStatus.Exist() {
		log.Printf("Forces exist in the controller")
	} else {
		log.Printf("Forces do not exist in the controller")
	}

	if forceStatus.Enabled() {
		log.Printf("Forces are enabled in the controller")
	} else {
		log.Printf("Forces are not enabled in the controller")
	}

	faultCode, err := canned.GetFaults(client)
	if err != nil {
		log.Printf("error reading fault code. %v", err)
		return
	}
	if faultCode.MajorType != 0 || faultCode.MinorType != 0 {
		log.Printf("Fault code: %d:%d", faultCode.MajorType, faultCode.MinorType)
		log.Printf("Fault Description: %s", faultCode.Events[0])
	} else {
		log.Printf("No fault code")
	}

}
