package gologix_tests

import (
	"testing"

	"github.com/danomagnum/gologix"
)

func TestListServices(t *testing.T) {
	// This test is a placeholder for the actual implementation
	tcs := getTestConfig()
	for _, tc := range tcs.GenericCIPTests {
		t.Run(tc.Address, func(t *testing.T) {
			client := gologix.NewClient(tc.Address)
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

			services, err := client.ListServices()
			if err != nil {
				t.Error(err)
				return
			}

			if services == nil {
				t.Error("services is nil")
				return
			}

			for i, service := range services {
				if service.Name != tc.Services[i].Name {
					t.Errorf("Name mismatch on entry %d. Have %s. Want %s.", i, service.Name, tc.Services[i].Name)
				}

				if service.Capabilities != gologix.ServiceCapabilityFlags(tc.Services[i].Capabilities) {
					t.Errorf("Capabilities mismatch on entry %d. Have %d. Want %d.", i, service.Capabilities, tc.Services[i].Capabilities)
				}
			}
		})
	}
}
