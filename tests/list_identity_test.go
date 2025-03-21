package gologix_tests

import (
	"testing"

	"github.com/danomagnum/gologix"
)

func TestListIdentity(t *testing.T) {
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

			identity, err := client.ListIdentity()
			if err != nil {
				t.Error(err)
				return
			}

			if identity == nil {
				t.Error("identity is nil")
				return
			}

			if identity.ProductName != tc.ProductName {
				t.Errorf("ProductName mismatch. Have %s. Want %s.", identity.ProductName, tc.ProductName)
			}

			version16 := uint16(tc.SoftwareVersionMinor)<<8 | uint16(tc.SoftwareVersionMajor)
			if identity.Revision != version16 {
				t.Errorf("Revision mismatch. Have %d. Want %d.", identity.Revision, version16)
			}

			if identity.SerialNumber != tc.SerialNumber {
				t.Errorf("SerialNumber mismatch. Have %d. Want %d.", identity.SerialNumber, tc.SerialNumber)
			}

			if identity.ProductCode != tc.ProductCode {
				t.Errorf("ProductCode mismatch. Have %d. Want %d.", identity.ProductCode, tc.ProductCode)
			}

		})
	}
}
