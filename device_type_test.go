package gologix

import (
	"testing"
)

func TestKnownDeviceTypeName(t *testing.T) {
	deviceType := DeviceType(0x02)
	name := deviceType.Name()
	if name != "AC Drive" {
		t.Errorf("expected 'AC Drive', got '%s'", name)
	}
}

func TestUnknownDeviceTypeName(t *testing.T) {
	deviceType := DeviceType(0x99)
	name := deviceType.Name()
	if name != "Unknown" {
		t.Errorf("expected 'Unknown', got '%s'", name)
	}
}
