package gologix_tests

import (
	"encoding/json"
	"fmt"
	"os"
)

type TestConfig struct {
	PlcAddress           string `json:"PLC_Address"`
	ProductCode          uint16 `json:"ProductCode"`
	SoftwareVersionMajor byte   `json:"SoftwareVersionMajor"`
	SoftwareVersionMinor byte   `json:"SoftwareVersionMinor"`
	SerialNumber         uint32 `json:"SerialNumber"`
	ProductName          string `json:"ProductName"`
}

func getTestConfig() TestConfig {

	f, err := os.Open("test_config.json")
	if err != nil {
		panic(fmt.Sprintf("couldn't open file: %v", err))
	}
	defer f.Close()
	dec := json.NewDecoder(f)

	var tc TestConfig

	dec.Decode(&tc)

	return tc
}
