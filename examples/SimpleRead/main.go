package main

import (
	"fmt"
	"gologix"
)

func main() {

	client := &gologix.Client{IPAddress: "192.168.2.241"}
	err := client.Connect()
	if err != nil {
		fmt.Printf("Error opening client. %v", err)
		return
	}
	defer client.Disconnect()

	var tag1 int16
	err = client.Read("testint", &tag1)
	if err != nil {
		fmt.Printf("error reading testint. %v", err)
	}
	fmt.Printf("tag1 has value %d", tag1)

}
