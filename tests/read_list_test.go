package gologix_tests

import (
	"testing"

	"github.com/danomagnum/gologix"
)

// bug report (issue 8): read list fails if one of the tags is a string.
func TestReadListWithString(t *testing.T) {
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

			tags := make([]string, 5)
			types := make([]any, 5)
			elements := make([]int, 5)

			tags[0] = "program:gologix_tests.ReadBool"     // false
			tags[1] = "program:gologix_tests.ReadDint"     // 36
			tags[2] = "program:gologix_tests.ReadString"   // "Somethig"
			tags[3] = "program:gologix_tests.ReadReal"     // 93.45
			tags[4] = "program:gologix_tests.ReadDints[0]" // 4351

			types[0] = true
			types[1] = int32(0)
			types[2] = ""
			types[3] = float32(0)
			types[4] = int32(0)

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
		})
	}

}

// bug report (issue 8): read list fails if one of the tags is a string.
func TestReadMultiWithString(t *testing.T) {
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
		})
	}

}

// bug report (issue 59): ReadList batch fails with []string types when reading
// more than one element. The single-element STRING path was already handled
// (issue 8); the multi-element branch fell through to the generic readValue,
// which returns "don't know what to do with a struct" for CIPTypeStruct.
func TestReadListWithStringArray(t *testing.T) {
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

			count := 10
			if client.ConnectionSize < 600 {
				count = 5
			}

			tags := []string{"program:gologix_tests.ReadStrings"}
			types := []any{make([]string, count)}
			elements := []int{count}

			vals, err := client.ReadList(tags, types, elements)
			if err != nil {
				t.Errorf("ReadList shouldn't have failed but did: %v", err)
				return
			}

			if len(vals) != 1 {
				t.Fatalf("expected 1 result, got %d", len(vals))
			}

			got, ok := vals[0].([]any)
			if !ok {
				t.Fatalf("expected []any, got %T", vals[0])
			}

			want := []string{"a", "b", "cd", "efg", "hijk", "lmnop", "qrstuvw", "xyz123", "0123456789", "9876543210"}
			if len(got) != count {
				t.Fatalf("expected %d strings, got %d", count, len(got))
			}

			for i, w := range want[:count] {
				g, ok := got[i].(string)
				if !ok {
					t.Errorf("element %d: expected string, got %T", i, got[i])
					continue
				}
				if g != w {
					t.Errorf("element %d: expected %q, got %q", i, w, g)
				}
			}
		})
	}
}

// TODO: When partial transfer continuations are implemented, add the ability to read all 10 strings.  Right now,
// when the PLC gets a small forward open, you can't read them all (looking at you v20 enbt)
func TestReadWithStringArray(t *testing.T) {
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

			count := 10
			if client.ConnectionSize < 600 {
				count = 5
			}
			tag := "program:gologix_tests.ReadStrings"
			got := make([]string, count)

			err = client.Read(tag, got)
			if err != nil {
				t.Errorf("ReadList shouldn't have failed but did: %v", err)
				return
			}

			want := []string{"a", "b", "cd", "efg", "hijk", "lmnop", "qrstuvw", "xyz123", "0123456789", "9876543210"}
			if len(got) != len(want[:count]) {
				t.Fatalf("expected %d strings, got %d", count, len(got))
			}

			for i, w := range want[:count] {
				g := got[i]
				if g != w {
					t.Errorf("element %d: expected %q, got %q", i, w, g)
				}
			}
		})
	}
}
