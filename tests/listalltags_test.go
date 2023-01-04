package gologix_tests

import (
	"log"
	"testing"

	"github.com/danomagnum/gologix"
)

func TestList(t *testing.T) {

	client := gologix.NewClient("192.168.2.241")
	err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Disconnect()

	err = client.ListAllTags(0)
	if err != nil {
		t.Error(err)
		return
	}

	log.Printf("Tags: %+v\n", client.KnownTags["testdintarr"])

	// check that we picked up all the test tags properly
	tests := make(map[string]gologix.KnownTag)
	tests["testdintarr"] = gologix.KnownTag{
		Name:        "TestDintArr",
		Type:        gologix.CIPTypeDINT,
		Array_Order: []int{10},
	}
	tests["testdint"] = gologix.KnownTag{
		Name:        "TestDint",
		Type:        gologix.CIPTypeDINT,
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
