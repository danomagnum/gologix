package gologix

import "testing"

// TestNewClientHostPortParsing locks in how NewClient splits a target string
// into IpAddress/Port. NewClient has no error return, so "bad" port input
// doesn't fail the call -- it's rejected by falling back to the default CIP
// port (44818) instead of carrying through garbage.
func TestNewClientHostPortParsing(t *testing.T) {
	cases := []struct {
		name     string
		target   string
		wantHost string
		wantPort uint
	}{
		{name: "host only", target: "192.168.1.100", wantHost: "192.168.1.100", wantPort: portDefault},
		{name: "host with port", target: "192.168.1.100:44818", wantHost: "192.168.1.100", wantPort: 44818},
		{name: "host with non-default port", target: "192.168.1.100:12345", wantHost: "192.168.1.100", wantPort: 12345},
		{name: "hostname with port", target: "plc.local:44818", wantHost: "plc.local", wantPort: 44818},
		{name: "ipv6 with port", target: "[::1]:44818", wantHost: "::1", wantPort: 44818},
		{name: "ipv6 without port", target: "::1", wantHost: "::1", wantPort: portDefault},
		{name: "non-numeric port falls back to default", target: "192.168.1.100:abc", wantHost: "192.168.1.100", wantPort: portDefault},
		{name: "port out of uint16 range falls back to default", target: "192.168.1.100:99999", wantHost: "192.168.1.100", wantPort: portDefault},
		{name: "negative port falls back to default", target: "192.168.1.100:-1", wantHost: "192.168.1.100", wantPort: portDefault},
		{name: "trailing colon with no port falls back to default", target: "192.168.1.100:", wantHost: "192.168.1.100", wantPort: portDefault},
		{name: "malformed too many colons falls back to whole string as host", target: "192.168.1.100:44818:extra", wantHost: "192.168.1.100:44818:extra", wantPort: portDefault},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewClient(tc.target)
			if c.Controller.IpAddress != tc.wantHost {
				t.Errorf("IpAddress = %q; want %q", c.Controller.IpAddress, tc.wantHost)
			}
			if c.Controller.Port != tc.wantPort {
				t.Errorf("Port = %d; want %d", c.Controller.Port, tc.wantPort)
			}
		})
	}
}
