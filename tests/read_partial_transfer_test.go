package gologix_tests

import (
	"testing"

	"github.com/danomagnum/gologix"
)

// TestReadPartialTransferHardwareLongDints exercises the FragRead loop against
// a real ControlLogix / CompactLogix controller. With a 256-byte CIP connection,
// the 400-byte response payload for Program:gologix_tests.LongDints (DINT[100])
// cannot fit in a single packet and the controller must split the reply across
// multiple CIPService_FragRead requests with a cumulative byte offset. Without
// the partial-transfer handling in Read_single this returns a parse error or
// truncated data; with the handling in place the full 100-element slice arrives
// intact and matches the seeded values.
//
// Requires the canonical tests/gologix_tests_Program.L5X imported on the
// controller and an entry in tests/test_config.json pointing at it.
func TestReadPartialTransferHardwareLongDints(t *testing.T) {
	tcs := getTestConfig()
	for _, tc := range tcs.TagReadWriteTests {
		if tc.Skip {
			continue
		}
		t.Run(tc.PlcAddress, func(t *testing.T) {
			client := gologix.NewClient(tc.PlcAddress)
			// 256-byte connection so the 400-byte LongDints[100] response
			// cannot fit in a single CIP packet, forcing the controller to
			// split the reply across multiple FragRead requests.
			client.ConnectionSize = 256

			if err := client.Connect(); err != nil {
				t.Fatalf("connect: %v", err)
			}
			defer func() {
				if err := client.Disconnect(); err != nil {
					t.Errorf("disconnect: %v", err)
				}
			}()

			if err := client.ListAllTags(1); err != nil {
				t.Fatalf("ListAllTags: %v", err)
			}

			tag := "Program:gologix_tests.LongDints"
			have := make([]int32, 100)
			if err := client.Read(tag, have); err != nil {
				t.Fatalf("read %s: %v", tag, err)
			}

			// Spot-check the same indices the existing ReadList test verifies.
			checks := []struct {
				index int
				want  int32
			}{
				{0, 5556},
				{89, 2329},
				{90, 888884},
				{99, 232058},
			}
			for _, c := range checks {
				if have[c.index] != c.want {
					t.Errorf("LongDints[%d] = %d, want %d", c.index, have[c.index], c.want)
				}
			}
		})
	}
}
