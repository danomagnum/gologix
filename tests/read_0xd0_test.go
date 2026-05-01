package gologix_tests

import (
	"log"
	"testing"

	"github.com/danomagnum/gologix"
)

type str struct {
	Dat [2]byte
}

func TestRead0xd0(t *testing.T) {
	// This test is for a specific case where the tag definition has a type of 0xd0 (unknown string type)
	// and an element count of 4. We want to ensure that the code can handle this case without error.

	client := gologix.NewClient("localhost")
	err := client.Connect()
	if err != nil {
		log.Printf("Error opening client: %v", err)
		return
	}
	defer client.Disconnect()

	var text str

	err = client.Read("TEXT", &text)
	if err != nil {
		log.Fatalf("Read failed: %v", err)
	}

	log.Printf("Tag'TEXT' : '%v'", text)

}
