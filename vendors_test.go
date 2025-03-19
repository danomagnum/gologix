package gologix

import (
	"testing"
)

func TestKnownVendorName(t *testing.T) {
	vendor := VendorId(1)
	name := vendor.Name()
	if name != "Rockwell Automation/Allen-Bradley" {
		t.Errorf("expected 'Rockwell Automation/Allen-Bradley', got '%s'", name)
	}
}

func TestUnknownVendorName(t *testing.T) {
	vendor := VendorId(9999)
	name := vendor.Name()
	if name != "Reserved" {
		t.Errorf("expected 'Reserved', got '%s'", name)
	}
}
