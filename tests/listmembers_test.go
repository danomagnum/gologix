package gologix_tests

import (
	"log"
	"testing"

	"github.com/danomagnum/gologix"
)

func TestMembersList(t *testing.T) {
	t.Skip("controller specific test")

	tcs := getTestConfig()
	for _, tc := range tcs.PlcList {
		t.Run(tc.PlcAddress, func(t *testing.T) {
			client := gologix.NewClient(tc.PlcAddress)
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

			err = client.ListAllTags(0)

			if err != nil {
				t.Errorf("problem getting tag list: %v", err)
				return
			}

			kt, ok := client.KnownTags["program:gologix_tests.readudt"]
			if !ok {
				t.Errorf("TestUDT not found")
				return
			}

			members, err := client.ListMembers(uint32(kt.Info.Template_ID()))
			if err != nil {
				t.Error(err)
			}

			log.Printf("Members: %+v", members)

			if members.Name != "TestUDT" {
				t.Errorf("expected name to be 'TestUDT', got '%v'", members.Name)
			}

			if len(members.Members) != 2 {
				t.Errorf("expected 2 members, got %v", len(members.Members))
			}

			if members.Members[0].Name != "Field1" {
				t.Errorf("expected first field to be 'Field1', got '%v'", members.Members[0].Name)
			}

			if members.Members[1].Name != "Field2" {
				t.Errorf("expected second field to be 'Field2', got '%v'", members.Members[1].Name)
			}

			if members.Members[0].Info.Type != 0xC4 {
				t.Errorf("expected first field type to be 0xC4, got %v", members.Members[0].Info.Type)
			}

			if members.Members[1].Info.Type != 0xCA {
				t.Errorf("expected second field type to be 0xCA, got %v", members.Members[1].Info.Type)
			}

			if members.Members[0].Info.Offset != 0 {
				t.Errorf("expected first field offset to be 0, got %v", members.Members[0].Info.Offset)
			}
			if members.Members[1].Info.Offset != 4 {
				t.Errorf("expected second field offset to be 4, got %v", members.Members[1].Info.Offset)
			}
		})
	}

}
