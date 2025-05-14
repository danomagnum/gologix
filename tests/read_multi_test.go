package gologix_tests

import (
	"log"
	"testing"

	"github.com/danomagnum/gologix"
)

func TestReadMultiNew(t *testing.T) {
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

			type multiread struct {
				TestInt  int16   `gologix:"TestInt"`
				TestDint int32   `gologix:"TestDint"`
				TestArr  []int32 `gologix:"TestDintArr[2]"`
			}
			var mr multiread
			mr.TestArr = make([]int32, 5)

			// call the read multi function with the structure passed in as a pointer.
			err = client.ReadMulti(&mr)
			if err != nil {
				log.Printf("error reading testint. %v", err)
			}

			if mr.TestInt != 999 {
				t.Errorf("TestInt expected 999 but got %d", mr.TestInt)
			}
			if mr.TestDint != 36 {
				t.Errorf("TestDint expected 36 but got %d", mr.TestDint)
			}
			if mr.TestArr[0] != 4353 {
				t.Errorf("TestArr[0] expected 4353 but got %d", mr.TestArr[0])
			}
			if mr.TestArr[1] != 4354 {
				t.Errorf("TestArr[1] expected 4354 but got %d", mr.TestArr[1])
			}
			if mr.TestArr[2] != 4355 {
				t.Errorf("TestArr[2] expected 4355 but got %d", mr.TestArr[2])
			}
			if mr.TestArr[3] != 4356 {
				t.Errorf("TestArr[3] expected 4356 but got %d", mr.TestArr[3])
			}
		})
	}
}

func TestReadMap(t *testing.T) {
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

			mr := make(map[string]any)
			mr["TestInt"] = int16(0)
			mr["TestDint"] = int32(0)
			mr["TestDintArr[2]"] = make([]int32, 5)

			// call the read multi function with the structure passed in as a pointer.
			err = client.ReadMap(mr)
			if err != nil {
				log.Printf("error reading testint. %v", err)
			}

			if mr["TestInt"].(int16) != 999 {
				t.Errorf("TestInt expected 999 but got %d", mr["TestInt"])
			}
			if mr["TestDint"].(int32) != 36 {
				t.Errorf("TestDint expected 36 but got %d", mr["TestDint"])
			}
			arr, ok := mr["TestDintArr[2]"].([]any)
			if !ok {
				t.Error("didn't get an int32 slice for TestDintArr[2]")
			}
			if arr[0].(int32) != 4353 {
				t.Errorf("TestArr[0] expected 4353 but got %d", arr[0])
			}
			if arr[1].(int32) != 4354 {
				t.Errorf("TestArr[1] expected 4354 but got %d", arr[1])
			}
			if arr[2].(int32) != 4355 {
				t.Errorf("TestArr[2] expected 4355 but got %d", arr[2])
			}
			if arr[3].(int32) != 4356 {
				t.Errorf("TestArr[3] expected 4356 but got %d", arr[3])
			}
		})
	}
}

// if you don't know the types of the tags, you can use the ReadMap function with nil values.
// this will read the tags and fill in the values.  The types are determined by the tag type on the PLC.
// Note that by specifying the type as in TestReadMap, you can read arrays and such.  In the case of
// unknown types, we always read a length of 1.
func TestReadMapAny(t *testing.T) {
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

			mr := make(map[string]any)
			mr["TestInt"] = nil
			mr["TestDint"] = nil
			mr["TestDintArr[2]"] = nil

			// call the read multi function with the structure passed in as a pointer.
			err = client.ReadMap(mr)
			if err != nil {
				log.Printf("error reading testint. %v", err)
			}

			if mr["TestInt"].(int16) != 999 {
				t.Errorf("TestInt expected 999 but got %d", mr["TestInt"])
			}
			if mr["TestDint"].(int32) != 36 {
				t.Errorf("TestDint expected 36 but got %d", mr["TestDint"])
			}
			if mr["TestDintArr[2]"].(int32) != 4353 {
				t.Errorf("TestArr[0] expected 4353 but got %d", mr["TestDintArr[2]"])
			}
		})
	}
}

func TestReadMultiWithGaps(t *testing.T) {
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

			type multiread struct {
				TestInt  int16 `gologix:"TestInt"`
				NotATag  int32
				TestDint int32   `gologix:"TestDint"`
				TestArr  []int32 `gologix:"TestDintArr[2]"`
			}
			var mr multiread
			mr.TestArr = make([]int32, 5)

			// call the read multi function with the structure passed in as a pointer.
			err = client.ReadMulti(&mr)
			if err != nil {
				log.Printf("error reading testint. %v", err)
			}

			if mr.TestInt != 999 {
				t.Errorf("TestInt expected 999 but got %d", mr.TestInt)
			}
			if mr.NotATag != 0 {
				t.Errorf("NotATag expected 0 but got %d", mr.NotATag)
			}
			if mr.TestDint != 36 {
				t.Errorf("TestDint expected 36 but got %d", mr.TestDint)
			}
			if mr.TestArr[0] != 4353 {
				t.Errorf("TestArr[0] expected 4353 but got %d", mr.TestArr[0])
			}
			if mr.TestArr[1] != 4354 {
				t.Errorf("TestArr[1] expected 4354 but got %d", mr.TestArr[1])
			}
			if mr.TestArr[2] != 4355 {
				t.Errorf("TestArr[2] expected 4355 but got %d", mr.TestArr[2])
			}
			if mr.TestArr[3] != 4356 {
				t.Errorf("TestArr[3] expected 4356 but got %d", mr.TestArr[3])
			}
		})
	}
}
