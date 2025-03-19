package gologix_tests

import (
	"encoding/json"
	"fmt"
	"os"
)

type TestConfig struct {
	PlcList []struct {
		PlcAddress           string `json:"PLC_Address"`
		ProductCode          uint16 `json:"ProductCode"`
		SoftwareVersionMajor byte   `json:"SoftwareVersionMajor"`
		SoftwareVersionMinor byte   `json:"SoftwareVersionMinor"`
		SerialNumber         uint32 `json:"SerialNumber"`
		ProductName          string `json:"ProductName"`
	} `json:"PLC_List"`

	ListIdentify []struct {
		Address              string `json:"Device_Address"`
		Vendor               uint16 `json:"Vendor"`
		DeviceType           uint16 `json:"DeviceType"`
		ProductCode          uint16 `json:"ProductCode"`
		SoftwareVersionMajor uint16 `json:"SoftwareVersionMajor"`
		SoftwareVersionMinor uint16 `json:"SoftwareVersionMinor"`
		Status               uint16 `json:"Status"`
		SerialNumber         uint32 `json:"SerialNumber"`
		ProductName          string `json:"ProductName"`
		State                uint8  `json:"State"`
	} `json:"ListIdentify"`

	ListServices []struct {
		Address  string `json:"Device_Address"`
		Services []struct {
			Name         string `json:"Name"`
			Capabilities uint16 `json:"Capabilities"`
		} `json:"Services"`
	} `json:"ListServices"`
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
