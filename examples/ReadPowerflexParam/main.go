package main

import (
	"fmt"
	"log"

	"github.com/danomagnum/gologix"
)

func main() {
	// PowerFlex IP Address - replace with your drive's IP address
	ipAddress := "192.168.2.246"

	// Create a new client with the PowerFlex IP address
	// The PowerFlex drive typically uses port 44818 (default EIP port)
	c := gologix.NewClient(ipAddress)
	err := c.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	defer c.Disconnect()

	ParamNo := 44 // param 44 = maximum hertz * 100

	path, err := gologix.Serialize(
		gologix.CipObject_Parameter,
		gologix.CIPInstance(ParamNo),
		gologix.CIPAttribute(1),
	)
	if err != nil {
		log.Printf("could not serialize path: %v", err)
		return
	}

	result, err := c.GenericCIPMessage(gologix.CIPService_GetAttributeSingle, path.Bytes(), []byte{})

	if err != nil {
		log.Fatalf("Failed to read parameter: %v", err)
	}

	// Parse the response
	val, err := result.Int16()
	if err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}
	maxhz := float64(val) / 100
	fmt.Printf("Parameter %d: %3.2f\n", ParamNo, maxhz)

	// you can also read multiple parameters at once using the GetAttributeList service
	path, err = gologix.Serialize(gologix.CipObject_DPIParams, gologix.CIPInstance(0))
	if err != nil {
		log.Printf("could not serialize path: %v", err)
		return
	}

	// a scattered read is used to read multiple parameters at once and takes a list of parameter IDs as 32 bit ints
	// Actually, the param IDs are 16 bit ints followed by a 16 bit WRITE value.  But since we're reading we can ignore the WRITE value and
	// use int32 to pad that area with zeroes.  If you swap to a ScatteredWrite you'd have to change this accordingly.
	params, err := gologix.Serialize([]int32{
		1,  // output freq
		2,  // command freq
		3,  // output current
		4,  // output voltage
		7,  // fault code
		27, // drive temp
	})

	if err != nil {
		log.Printf("could not serialize params: %v", err)
		return
	}

	r, err := c.GenericCIPMessage(gologix.CIPService_ScatteredRead, path.Bytes(), params.Bytes())

	if err != nil {
		log.Fatalf("Failed to read parameters: %v", err)
	}

	// the powerflex 525 uses 16 bit integers for all of its parameters
	// so we can use a struct to parse the response.
	// the response will contain the parameter ID followed by the value
	// for each parameter requested.  We don't need the parameter ID so we use _ to skip it.
	var resultData = struct {
		_             int16
		OutputFreq    int16
		_             int16
		CommandFreq   int16
		_             int16
		OutputCurrent int16
		_             int16
		OutputVoltage int16
		_             int16
		FaultCode     int16
		_             int16
		DriveTemp     int16
	}{}

	err = r.DeSerialize(&resultData)
	if err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	log.Printf("Drive Status: %+v", resultData)

}
