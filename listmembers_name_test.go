package gologix

import "testing"

func TestCleanTemplateName(t *testing.T) {
	tests := []struct {
		raw  string
		want string
	}{
		{"MyUDT\x00", "MyUDT"},
		{"myAOI;n0_1234\x00", "myAOI"},
		{"myAOI_Interface;n0_1234\x00", "myAOI_Interface"},
		{"STRING\x00", "STRING"},
		{"NoTerminator", "NoTerminator"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := cleanTemplateName(tt.raw); got != tt.want {
			t.Errorf("cleanTemplateName(%q) = %q, want %q", tt.raw, got, tt.want)
		}
	}
}
