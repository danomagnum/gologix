package main

import (
	"fmt"
	"log"
	"time"

	"github.com/danomagnum/gologix"
	"github.com/danomagnum/gologix/cipclass"
	"github.com/danomagnum/gologix/cipservice"
)

// this program will read the current controller time out of the PLC using a custom generic CIP message.
// this could also be done with the GetAttrList() function (see tests/getattr_test.go) and I would recommend that
// method for this specific purpose, but this example should apply to other devices like drives and other services
// on PLCs as desired.
func main() {
	client := gologix.NewClient("192.168.2.241")
	err := client.Connect()
	if err != nil {
		fmt.Print(err)
		return
	}
	defer func() {
		err := client.Disconnect()
		if err != nil {
			fmt.Printf("problem disconnecting. %v", err)
		}
	}()

	r, err := client.GenericCIPMessage(cipservice.GetAttributeList, cipclass.CipObject_TIME, 1, []byte{0x01, 0x00, 0x0B, 0x00})
	if err != nil {
		fmt.Printf("bad result: %v", err)
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
		fmt.Printf("could not deserialize response structure: %v", err)
		return
	}

	log.Printf("result: us: %v / %v", rs.Usecs, time.UnixMicro(int64(rs.Usecs)))

}
