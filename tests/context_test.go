package gologix_tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/danomagnum/gologix"
)

// A context that's already done must make the *WithContext methods bail out
// immediately with the context's error instead of sending a request to the PLC
// and waiting for a response (or the socket timeout).
func TestReadWriteWithContextBailsImmediately(t *testing.T) {
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

			tag := "Program:gologix_tests.ReadDint"

			// Warm up the tag/firmware caches with a normal request so the
			// canceled-context calls below don't do an incidental round trip
			// (e.g. the firmware probe newIOI triggers on first use).
			var warm int32
			if err := client.Read(tag, &warm); err != nil {
				t.Fatalf("warmup read failed: %v", err)
			}

			t.Run("canceled", func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				start := time.Now()
				var got int32
				err := client.ReadWithContext(ctx, tag, &got)
				elapsed := time.Since(start)
				if !errors.Is(err, context.Canceled) {
					t.Errorf("Read: want context.Canceled, got %v", err)
				}
				if elapsed > 200*time.Millisecond {
					t.Errorf("Read took %v, want it to bail out immediately", elapsed)
				}

				start = time.Now()
				err = client.WriteWithContext(ctx, tag, warm)
				elapsed = time.Since(start)
				if !errors.Is(err, context.Canceled) {
					t.Errorf("Write: want context.Canceled, got %v", err)
				}
				if elapsed > 200*time.Millisecond {
					t.Errorf("Write took %v, want it to bail out immediately", elapsed)
				}
			})

			t.Run("deadline exceeded", func(t *testing.T) {
				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
				defer cancel()

				start := time.Now()
				var got int32
				err := client.ReadWithContext(ctx, tag, &got)
				elapsed := time.Since(start)
				if !errors.Is(err, context.DeadlineExceeded) {
					t.Errorf("Read: want context.DeadlineExceeded, got %v", err)
				}
				if elapsed > 200*time.Millisecond {
					t.Errorf("Read took %v, want it to bail out immediately", elapsed)
				}

				start = time.Now()
				err = client.WriteWithContext(ctx, tag, warm)
				elapsed = time.Since(start)
				if !errors.Is(err, context.DeadlineExceeded) {
					t.Errorf("Write: want context.DeadlineExceeded, got %v", err)
				}
				if elapsed > 200*time.Millisecond {
					t.Errorf("Write took %v, want it to bail out immediately", elapsed)
				}
			})
		})
	}
}
