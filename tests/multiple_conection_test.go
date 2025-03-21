package gologix_tests

import (
	"testing"

	"github.com/danomagnum/gologix"
)

func TestMultipleConns(t *testing.T) {
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

			client2 := gologix.NewClient(tc.PlcAddress)
			err = client2.Connect()
			if err != nil {
				t.Error(err)
				return
			}
			defer func() {
				err := client2.Disconnect()
				if err != nil {
					t.Errorf("problem disconnecting. %v", err)
				}
			}()

			tag := "Program:gologix_tests.ReadDints[0]"
			have := make([]int32, 5)
			want := []int32{4351, 4352, 4353, 4354, 4355}

			err = client.Read(tag, have)
			if err != nil {
				t.Errorf("Problem reading %s. %v", tag, err)
				return
			}
			for i := range want {
				if have[i] != want[i] {
					t.Errorf("index %d wanted %v got %v", i, want[i], have[i])
				}
			}

			err = client2.Read(tag, have)
			if err != nil {
				t.Errorf("Problem reading %s. %v", tag, err)
				return
			}
			for i := range want {
				if have[i] != want[i] {
					t.Errorf("index %d wanted %v got %v", i, want[i], have[i])
				}
			}
		})
	}
}
