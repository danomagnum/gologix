package gologix_tests

import (
	"encoding/json"
	"fmt"
	"os"
)

type TestConfig struct {
	PLC_Address string `json:"PLC_Address"`
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
