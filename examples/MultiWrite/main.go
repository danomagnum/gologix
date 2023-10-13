package main

import (
	"log"

	"github.com/danomagnum/gologix"
)

// Demo program for writing multiple tags at once.
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
	// setup a disconnect.  If you don't disconnect you might have trouble reconnecting
	defer client.Disconnect()

	// Set up a map[string]any of tag:value pairs.

	write_map := make(map[string]any)
	write_map["program:gologix_tests.MultiWriteInt"] = int16(123)
	write_map["program:gologix_tests.MultiWriteReal"] = float32(456.7)
	write_map["program:gologix_tests.MultiWriteDint"] = int32(891011)
	write_map["program:gologix_tests.MultiWriteString"] = "Test String!"
	write_map["program:gologix_tests.MultiWriteBool"] = true

	err = client.WriteMap(write_map)
	if err != nil {
		log.Printf("error writing to multiple tags at once: %v", err)
	}

	log.Printf("No Errors!")
	// no error = write OK.

}
