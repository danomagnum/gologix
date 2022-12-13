package gologix_tests

import (
	"testing"

	"github.com/danomagnum/gologix"
)

func TestMembersList(t *testing.T) {

	client := gologix.NewClient("192.168.2.241")
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()

	_, err = client.ListMembers(1658)
	if err != nil {
		t.Error(err)
	}

}
