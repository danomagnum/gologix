package gologix

import (
	"testing"
)

func TestParseListServicesResponse(t *testing.T) {
	itemBytes := []byte{1, 0, 32, 1, 67, 111, 109, 109, 117, 110, 105, 99, 97, 116, 105, 111, 110, 115, 0, 0}

	response := CIPListService{}
	err := response.ParseFromBytes(itemBytes)
	if err != nil {
		t.Fatalf("error parsing list services response: %v", err)
	}

	expected := CIPListService{
		EncapProtocolVersion: 1,
		Capabilities:         ServiceCapabilityFlag_CipEncapsulation | ServiceCapabilityFlag_SupportsClass1UDP,
		Name:                 "Communications",
	}

	if response != expected {
		t.Fatalf("expected %v, got %v", expected, response)
	}
}
