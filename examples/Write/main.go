package main

import (
	"fmt"
	"gologix"
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
		fmt.Printf("Error opening client. %v", err)
		return
	}
	// setup a disconnect.  If you don't disconnect you might have trouble reconnecting
	defer client.Disconnect()

	// the variable you want to write needs to be the proper type for the controller tag.  See GoVarToLogixType() for info

	mydint := int32(12345)

	err = client.Write("WriteUDTs[5].Field1", mydint)
	if err != nil {
		fmt.Printf("error writing to tag 'WriteUDTs[5].Field1'. %v", err)
	}

	// no error = write OK.

}
