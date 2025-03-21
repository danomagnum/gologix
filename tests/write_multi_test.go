package gologix_tests

import (
	"testing"

	"github.com/danomagnum/gologix"
)

func TestWriteMulti(t *testing.T) {
	tcs := getTestConfig()
	for _, tc := range tcs.TagReadWriteTests {
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

			want := TestUDT{Field1: 215, Field2: 25.1}

			err = client.Write("Program:gologix_tests.WriteUDTs[2]", want)
			if err != nil {
				t.Errorf("error writing. %v", err)
			}

			have := TestUDT{}
			err = client.Read("Program:gologix_tests.WriteUDTs[2]", &have)
			if err != nil {
				t.Errorf("error reading. %v", err)
			}

			if have != want {
				t.Errorf("have %v. Want %v", have, want)

			}
		})
	}
}
