package gologix

import (
	"fmt"
	"testing"
)

func TestList(t *testing.T) {

	client := &Client{IPAddress: "192.168.2.241"}
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()

	client.ListAllTags(0)

	fmt.Printf("Tags: %+v\n", client.KnownTags["testdintarr"])

	// check that we picked up all the test tags properly
	tests := make(map[string]KnownTag)
	tests["testdintarr"] = KnownTag{
		Name:        "TestDintArr",
		Type:        CIPTypeDINT,
		Array_Order: []int{10},
	}
	tests["testdint"] = KnownTag{
		Name:        "TestDint",
		Type:        CIPTypeDINT,
		Array_Order: []int{},
	}

	for k := range tests {
		t.Run(k, func(t *testing.T) {
			have := client.KnownTags[k]
			want := tests[k]

			if have.Name != want.Name {
				t.Errorf("Name Mismatch. Have %s. Want %s.", have.Name, want.Name)
			}
			if have.Type != want.Type {
				t.Errorf("Type Mismatch. Have %s. Want %s.", have.Type, want.Type)
			}
			if !compare_array_order(have.Array_Order, want.Array_Order) {
				t.Errorf("Array Order Mismatch. Have %v. Want %v.", have.Array_Order, want.Array_Order)

			}
		})
	}

}

func compare_array_order(have, want []int) bool {
	if len(have) != len(want) {
		return false
	}

	for i := range have {
		if have[i] != want[i] {
			return false
		}
	}
	return true
}
