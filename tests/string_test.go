package gologix_tests

import (
	"testing"

	"github.com/danomagnum/gologix"
)

type LilString struct {
	Length uint32
	Data   [25]byte
}

func (ls *LilString) String() string {
	return string(ls.Data[:ls.Length])
}

type BigString struct {
	Length uint32
	Data   [100]byte
}

func (bs *BigString) String() string {
	return string(bs.Data[:bs.Length])
}

func TestReadCustomString(t *testing.T) {
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

			tag := "Program:gologix_tests.LilString"
			have := LilString{}
			want := LilString{Length: 15}
			copy(want.Data[:], "Im a lil string")

			err = client.Read(tag, &have)
			if err != nil {
				t.Errorf("Problem reading %s. %v", tag, err)
				return
			}
			if have.Length != want.Length {
				t.Errorf("Length mismatch. have: %d want: %d", have.Length, want.Length)
			}
			if string(have.Data[:have.Length]) != string(want.Data[:want.Length]) {
				t.Errorf("Data mismatch. have: %s want: %s", string(have.Data[:have.Length]), string(want.Data[:want.Length]))
			}

			tag = "Program:gologix_tests.BigString"
			have2 := BigString{}
			want2 := BigString{Length: 16}
			copy(want2.Data[:], "Im a big string.")
			err = client.Read(tag, &have2)
			if err != nil {
				t.Errorf("Problem reading %s. %v", tag, err)
				return
			}
			if have2.Length != want2.Length {
				t.Errorf("Length mismatch. have: %d want: %d", have2.Length, want2.Length)
			}
			if string(have2.Data[:have2.Length]) != string(want2.Data[:want2.Length]) {
				t.Errorf("Data mismatch. have: %s want: %s", string(have2.Data[:have2.Length]), string(want2.Data[:want2.Length]))
			}

		})
	}
}

func TestReadCustomStringArray(t *testing.T) {
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

			tag := "Program:gologix_tests.LilStrings"
			have := make([]LilString, 10)
			want := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

			err = client.Read(tag, have)
			if err != nil {
				t.Errorf("Problem reading %s. %v", tag, err)
				return
			}

			for i := 0; i < len(want); i++ {
				if have[i].String() != want[i] {
					t.Errorf("Data mismatch at index %d. have: %s want: %s", i, have[i].String(), want[i])
				}
			}

			tag = "Program:gologix_tests.BigStrings"
			have2 := make([]BigString, 10)
			want2 := []string{"ZERO", "ONE", "TWO", "THREE", "FOUR", "FIVE", "SIX", "SEVEN", "EIGHT", "NINE"}
			if client.ConnectionSize <= 600 {
				// Reduce the array size for smaller connections
				have2 = make([]BigString, 4)
			}

			err = client.Read(tag, have2)
			if err != nil {
				t.Errorf("Problem reading %s. %v", tag, err)
				return
			}
			for i := 0; i < len(have2); i++ {
				if have2[i].String() != want2[i] {
					t.Errorf("Data mismatch at index %d. have: %s want: %s", i, have2[i].String(), want2[i])
				}
			}

		})
	}
}
