package gologix_tests

import (
	"encoding/xml"
	"io"
	"os"
	"testing"

	"github.com/danomagnum/gologix/l5x"
)

func TestDecodeL5X(t *testing.T) {
	var l5xData l5x.RSLogix5000Content

	f, err := os.Open("gologix_tests_Program.L5X")
	if err != nil {
		t.Error(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Error(err)
	}

	err = xml.Unmarshal(b, &l5xData)
	if err != nil {
		t.Error(err)
	}

	result, err := l5x.LoadTags(l5xData)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)

}
