package main

import (
	"fmt"
	"log"
	"os"

	"github.com/danomagnum/gologix"
)

// Demo program showing how to use ListIdentity and ListServices
// to discover information about EtherNet/IP devices on the network.
func main() {
	var err error

	// Setup the client with the PLC IP address
	// Replace with your actual device IP
	client := gologix.NewClient("192.168.2.244")

	// Connect using parameters in the client struct
	err = client.Connect()
	if err != nil {
		log.Printf("Error opening client. %v", err)
		os.Exit(1)
	}
	// Setup a deferred disconnect for proper cleanup
	defer client.Disconnect()

	// List Identity - Get information about the device
	fmt.Println("Listing Identity Information:")
	identity, err := client.ListIdentity()
	if err != nil {
		log.Printf("Error listing identity: %v", err)
		os.Exit(1)
	}

	// Display the identity information
	fmt.Printf("Identity Response: %+v\n", identity)
	fmt.Println()

	// List Services - Get information about available services
	fmt.Println("Listing Available Services:")
	services, err := client.ListServices()
	if err != nil {
		log.Printf("Error listing services: %v", err)
		os.Exit(1)
	}

	// Display the services information
	fmt.Printf("Services Response: %+v\n", services)
	fmt.Println()

	// You can further process and display specific fields from the
	// identity and services responses based on their actual structure
}
