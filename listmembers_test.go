package gologix

import (
	"testing"
)

func TestMembersList(t *testing.T) {

	client := &Client{IPAddress: "192.168.2.241"}
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
