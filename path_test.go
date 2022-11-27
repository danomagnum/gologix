package main

import (
	"testing"
)

func TestPath_gen(t *testing.T) {
	tar := []byte{0x01, 0x00, 0x12, 0x0C, 0x31, 0x37, 0x32, 0x2E, 0x32, 0x35, 0x2E, 0x35, 0x38, 0x2E, 0x31, 0x31, 0x01, 0x01}
	path := "1,0,2,172.25.58.11,1,1"
	res, err := pathgen(path)

	if err != nil {
		t.Errorf("Error in pathgen for %s. %v", path, err)
	}

	if !check_bytes(tar, res) {
		t.Errorf("Wrong Value for result.  \nWanted %v. \nGot    %v", tar, res)
	}
	t.Logf("Wrong Value for result.  \nWanted %v. \nGot    %v", tar, res)

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
