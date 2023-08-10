package gologix

import (
	"fmt"
	"testing"
)

func TestPath(t *testing.T) {
	var tests = []struct {
		path string
		want []byte
	}{
		{
			"1,0,2,172.25.58.11,1,1",
			[]byte{0x01, 0x00, 0x12, 0x0C, 0x31, 0x37, 0x32, 0x2E, 0x32, 0x35, 0x2E, 0x35, 0x38, 0x2E, 0x31, 0x31, 0x01, 0x01},
		},
		{
			"1,0,32,2,36,1",
			[]byte{0x01, 0x00, 0x20, 0x02, 0x24, 0x01},
		},
		{
			"1,0",
			[]byte{0x01, 0x00},
		},
	}

	for _, tt := range tests {

		testname := fmt.Sprintf("path: %s", tt.path)
		t.Run(testname, func(t *testing.T) {
			res, err := ParsePath(tt.path)
			if err != nil {
				t.Errorf("Error in pathgen for %s. %v", tt.path, err)
			}
			if !check_bytes(res.Bytes(), tt.want) {
				t.Errorf("Wrong Value for result.  \nWanted %v. \nGot    %v", tt.want, res)
			}
		})
	}

}

func check_bytes(s0, s1 []byte) bool {
	if len(s1) != len(s0) {
		return false
	}
	for i := range s0 {
		if s0[i] != s1[i] {
			return false
		}

	}
	return true
}

func TestPathBuild(t *testing.T) {
	client := Client{}
	client.SocketTimeout = 0

	pmp_ioi, err := client.NewIOI("Program:MainProgram", 16)
	if err != nil {
		t.Errorf("problem creating pmp ioi. %v", err)
	}

	tests := []struct {
		name string
		path []any
		want []byte
	}{
		{
			name: "connection manager only",
			path: []any{CipObject_ConnectionManager},
			want: []byte{0x20, 0x06},
		},
		{
			name: "backplane to slot 0",
			path: []any{CIPPort{PortNo: 1}, CIPAddress(0)},
			want: []byte{0x01, 0x00},
		},
		{
			name: "connection manager instance 1",
			path: []any{CipObject_ConnectionManager, CIPInstance(1)},
			want: []byte{0x20, 0x06, 0x24, 0x01},
		},
		{
			name: "Symbol Object Instance 0",
			path: []any{CipObject_Symbol, CIPInstance(0)},
			want: []byte{0x20, 0x6B, 0x24, 0x00},
		},
		{
			name: "Symbol Object Instance 0 of tag 'Program:MainProgram'",
			path: []any{pmp_ioi, CipObject_Symbol, CIPInstance(0)},
			want: []byte{0x91, 0x13, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x3a, 0x6d, 0x61, 0x69,
				0x6e, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x00, 0x20, 0x6B, 0x24, 0x00},
		},
		{
			name: "Template Attributes Instance 0x02E9",
			path: []any{CipObject_Template, CIPInstance(0x02E9)},
			want: []byte{0x20, 0x6C, 0x25, 0x00, 0xE9, 0x02},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			have, err := Serialize(tt.path...)
			if err != nil {
				t.Errorf("Problem building path. %v", err)
			}
			if !check_bytes(have.Bytes(), tt.want) {
				t.Errorf("ResultMismatch.\n Have %v\n Want %v\n", have.Bytes(), tt.want)
			}
		})
	}

}
