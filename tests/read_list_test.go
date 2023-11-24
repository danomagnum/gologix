package gologix_tests

import (
	"testing"

	"github.com/danomagnum/gologix"
)

// bug report (issue 8): read list fails if one of the tags is a string.
func TestReadListWithString(t *testing.T) {
	tc := getTestConfig()
	client := gologix.NewClient(tc.PLC_Address)
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err := client.Disconnect()
		if err != nil {
			t.Errorf("problem disconnecting. %v", err)
		}
	}()

	tags := make([]string, 5)
	types := make([]gologix.CIPType, 5)
	elements := make([]int, 5)

	tags[0] = "program:gologix_tests.ReadBool"     // false
	tags[1] = "program:gologix_tests.ReadDint"     // 36
	tags[2] = "program:gologix_tests.ReadString"   // "Somethig"
	tags[3] = "program:gologix_tests.ReadReal"     // 93.45
	tags[4] = "program:gologix_tests.ReadDints[0]" // 4351

	types[0] = gologix.CIPTypeBOOL
	types[1] = gologix.CIPTypeDINT
	types[2] = gologix.CIPTypeSTRING
	types[3] = gologix.CIPTypeREAL
	types[4] = gologix.CIPTypeDINT

	elements[0] = 1
	elements[1] = 1
	elements[2] = 1
	elements[3] = 1
	elements[4] = 1

	vals, err := client.ReadList(tags, types, elements)
	if err != nil {
		t.Errorf("shouldn't have failed but did. %v", err)
		return
	}

	v0, ok0 := vals[0].(bool)
	v1, ok1 := vals[1].(int32)
	v2, ok2 := vals[2].(string)
	v3, ok3 := vals[3].(float32)
	v4, ok4 := vals[4].(int32)

	if !ok0 || !ok1 || !ok2 || !ok3 || !ok4 {
		t.Errorf("A type cast failed.: %v %v %v %v %v ", ok0, ok1, ok2, ok3, ok4)
	}

	if v0 {
		t.Errorf("ReadBool should be false but wasn't.")
	}

	if v1 != 36 {
		t.Errorf("ReadDint should be 36 but was %v", v1)
	}

	if v2 != "Something" {
		t.Errorf("ReadString should be 'Something' but was '%v'", v2)
	}

	if v3 != 93.45 {
		t.Errorf("ReadReal should be 93.45 but was %v", v3)
	}

	if v4 != 4351 {
		t.Errorf("ReadDints[0] should be 4351 but was %v", v4)
	}

}

// bug report (issue 8): read list fails if one of the tags is a string.
func TestReadMultiWithString(t *testing.T) {
	tc := getTestConfig()
	client := gologix.NewClient(tc.PLC_Address)
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err := client.Disconnect()
		if err != nil {
			t.Errorf("problem disconnecting. %v", err)
		}
	}()

	type customRead struct {
		ReadBool   bool    `gologix:"program:gologix_tests.ReadBool"`
		ReadDint   int32   `gologix:"program:gologix_tests.ReadDint"`
		ReadString string  `gologix:"program:gologix_tests.ReadString"`
		ReadReal   float32 `gologix:"program:gologix_tests.ReadReal"`
		ReadDint0  int32   `gologix:"program:gologix_tests.ReadDints[0]"`
	}

	var cr customRead

	err = client.ReadMulti(&cr)
	if err != nil {
		t.Errorf("shouldn't have failed but did. %v", err)
		return
	}

	if cr.ReadBool {
		t.Errorf("ReadBool should be false but wasn't.")
	}

	if cr.ReadDint != 36 {
		t.Errorf("ReadDint should be 36 but was %v", cr.ReadDint)
	}

	if cr.ReadString != "Something" {
		t.Errorf("ReadString should be 'Something' but was '%v'", cr.ReadString)
	}

	if cr.ReadReal != 93.45 {
		t.Errorf("ReadReal should be 93.45 but was %v", cr.ReadReal)
	}

	if cr.ReadDint0 != 4351 {
		t.Errorf("ReadDints[0] should be 4351 but was %v", cr.ReadDint0)
	}

}
