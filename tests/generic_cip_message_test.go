package gologix_tests

import (
	"log"
	"testing"
	"time"

	"github.com/danomagnum/gologix"
)

// This test uses the GenericCIPMessage function to read attributes from a controller.  In this case it is reading
// the time object's usec since the unix epoch.
func TestGenericCIPMessage1(t *testing.T) {
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

			path, err := gologix.Serialize(gologix.CipObject_TIME, gologix.CIPInstance(1))
			if err != nil {
				t.Errorf("could not serialize path: %v", err)
				return
			}
			r, err := client.GenericCIPMessage(gologix.CIPService_GetAttributeList, path.Bytes(), []byte{0x01, 0x00, 0x0B, 0x00})
			if err != nil {
				t.Errorf("bad result: %v", err)
				return
			}
			type response_str struct {
				Attr_Count int16
				Attr_ID    uint16
				Status     uint16
				Usecs      int64 // the microseconds since the unix epoch.
			}

			rs := response_str{}
			err = r.DeSerialize(&rs)
			if err != nil {
				t.Errorf("could not deserialize response structure: %v", err)
				return
			}

			log.Printf("result: us: %v / %v", rs.Usecs, time.UnixMicro(int64(rs.Usecs)))
		})
	}
}

/*
func TestGenericCIPMessage2(t *testing.T) {
	tc := getTestConfig()
	client := gologix.NewClient(tc.PLC_Address)
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

	path, err := gologix.Serialize(gologix.CIPClass(0x04), gologix.CIPInstance(0x303), gologix.CIPAttribute(0x03))
	if err != nil {
		t.Errorf("could not serialize path: %v", err)
		return
	}
	r, err := client.GenericCIPMessage(0x14, path.Bytes(), []byte{})
	if err != nil {
		t.Errorf("bad result: %v", err)
		return
	}
	type response_str struct {
		Attr_Count int16
		Attr_ID    uint16
		Status     uint16
		Usecs      int64 // the microseconds since the unix epoch.
	}

	rs := response_str{}
	err = r.DeSerialize(&rs)
	if err != nil {
		t.Errorf("could not deserialize response structure: %v", err)
		return
	}

	log.Printf("result: us: %v / %v", rs.Usecs, time.UnixMicro(int64(rs.Usecs)))

}
*/
