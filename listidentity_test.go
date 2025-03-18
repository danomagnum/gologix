package gologix

import (
	"testing"
)

func TestParseListIdentityResponse(t *testing.T) {
	itemBytes := []byte{1, 0, 0, 2, 175, 18, 192, 168, 1, 65, 0, 0, 0, 0, 0, 0, 0, 0, 21, 5, 12, 0, 1, 0, 1, 2, 52, 0, 3, 232, 34, 23, 20, 73, 81, 32, 83, 101, 110, 115, 111, 114, 32, 78, 101, 116, 32, 83, 121, 115, 116, 101, 109, 3}
	response := listIdentityResponeBody{}
	err := response.ParseFromBytes(itemBytes)
	if err != nil {
		t.Fatalf("error parsing list identity response: %v", err)
	}

	expected := listIdentityResponeBody{
		EncapProtocolVersion: 1,
		SocketAddress: listIdentitySocketAddress{
			Family:  2,
			Port:    44818,
			Address: 3232235841,
			Zero0:   0,
			Zero1:   0,
		},
		Vendor:       0x0515,
		DeviceType:   12,
		ProductCode:  1,
		Revision:     0x0201,
		Status:       0x0034,
		SerialNumber: 0x1722E803,
		ProductName:  "IQ Sensor Net System",
		State:        3,
	}

	if response != expected {
		t.Fatalf("expected %v, got %v", expected, response)
	}
}
