package gologix

import (
	"testing"
)

func TestSubList(t *testing.T) {

	client := &Client{IPAddress: "192.168.2.241"}
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()

	client.ListSubTags("Program:Shed", 1)
	client.ListSubTags("Program:Shed.fan_run_time", 1)

}

func compare_array_order(have, want []int) bool {
	if len(have) != len(want) {
		return false
	}

	for i := range have {
		if have[i] != want[i] {
			return false
		}
	}
	return true
}
