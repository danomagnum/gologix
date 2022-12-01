package gologix

import "testing"

func TestList2(t *testing.T) {

	plc := &PLC{IPAddress: "192.168.2.241"}
	err := plc.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer plc.Disconnect()

	plc.ListAllTags2(0)

}
