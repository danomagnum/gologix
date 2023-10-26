package gologix_tests

import (
	"log"
	"testing"
	"time"

	"github.com/danomagnum/gologix"
	"github.com/danomagnum/gologix/cipclass"
	"github.com/danomagnum/gologix/cipservice"
)

// This test uses the GenericCIPMessage function to read attributes from a controller.  In this case it is reading
// the time object's usec since the unix epoch.
func TestGenericCIPMessage1(t *testing.T) {
	client := gologix.NewClient("192.168.2.241")
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

	r, err := client.GenericCIPMessage(cipservice.GetAttributeList, cipclass.CipObject_TIME, 1, []byte{0x01, 0x00, 0x0B, 0x00})
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
