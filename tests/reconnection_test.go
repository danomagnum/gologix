package gologix_tests

import (
	"testing"

	"github.com/danomagnum/gologix"
)

func TestReconnection(t *testing.T) {
	client := gologix.NewClient("192.168.2.241")
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err := client.Disconnect()
		if err != nil {
			t.Errorf("problem disconnecting. %v", err)
		}
	}()
	tag := "Program:gologix_tests.ReadDints[0]"
	have := make([]int32, 5)
	want := []int32{4351, 4352, 4353, 4354, 4355}

	// first we'll do a read which should succeed
	err = client.Read(tag, have)
	if err != nil {
		t.Errorf("Problem reading %s. %v", tag, err)
		return
	}
	for i := range want {
		if have[i] != want[i] {
			t.Errorf("index %d wanted %v got %v", i, want[i], have[i])
		}
	}

	// then we'll close the socket.
	client.DebugCloseConn()

	// then read again.  This should fail, but cause a "proper" disconnect.
	err = client.Read(tag, have)
	if err == nil {
		t.Errorf("read should have failed but didn't.")
		return
	}

	// now read should work again because AutoConnect = true on the client.
	err = client.Read(tag, have)
	if err != nil {
		t.Errorf("Problem reading after reconnect %s. %v", tag, err)
		return
	}
	for i := range want {
		if have[i] != want[i] {
			t.Errorf("index %d wanted %v got %v", i, want[i], have[i])
		}
	}
}

func TestNoReconnection(t *testing.T) {
	client := gologix.NewClient("192.168.2.241")
	client.AutoConnect = false
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err := client.Disconnect()
		if err != nil {
			t.Errorf("problem disconnecting. %v", err)
		}
	}()
	tag := "Program:gologix_tests.ReadDints[0]"
	have := make([]int32, 5)
	want := []int32{4351, 4352, 4353, 4354, 4355}

	// first we'll do a read which should succeed
	err = client.Read(tag, have)
	if err != nil {
		t.Errorf("Problem reading %s. %v", tag, err)
		return
	}
	for i := range want {
		if have[i] != want[i] {
			t.Errorf("index %d wanted %v got %v", i, want[i], have[i])
		}
	}

	// then we'll close the socket.
	client.DebugCloseConn()

	// then read again.  This should fail, but cause a "proper" disconnect.
	err = client.Read(tag, have)
	if err == nil {
		t.Errorf("read should have failed but didn't.")
		return
	}

	// now read should work again because AutoConnect = true on the client.
	err = client.Read(tag, have)
	if err == nil {
		t.Errorf("read should have failed again but didn't.")
		return
	}

}
