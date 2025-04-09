package main

import (
	"encoding/xml"
	"io"
	"log"
	"os"

	"github.com/danomagnum/gologix/l5x"
)

// This is an example of how to read tags from an L5X file.
//
// the tags will be loaded into a map[string]any where the key is the tag name and the value is the tag value.
// any program tags will be in a nested map[string]any on a key of "program:<programname>".
// structures will be in nested map[string]any's.
// arrays will be in slices.
func main() {
	var l5xData l5x.RSLogix5000Content

	f, err := os.Open("gologix_tests_Program.L5X")
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	err = xml.Unmarshal(b, &l5xData)
	if err != nil {
		log.Fatal(err)
	}

	tags, err := l5x.LoadTags(l5xData)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range tags {
		log.Printf("%s: %v\n", k, v)
	}

	tagComments, err := l5x.LoadTagComments(l5xData)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range tagComments {
		log.Printf("%s: %v\n", k, v)
	}

	rungComments, err := l5x.LoadRungComments(l5xData)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range rungComments {
		log.Printf("%s: %v\n", k, v)
	}

}
