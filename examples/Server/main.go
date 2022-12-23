package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/danomagnum/gologix"
)

func main() {

	////////////////////////////////////////////////////
	// First we set up the tag providers.
	//
	// Each one will have a path and an object that fulfills the gologix.TagProvider interface
	// We set those up and then pass them to the Router object.
	// here we're using the build in map tag provider which just maps cip reads and writes to/from a go map
	//
	// In theory, it should be easy to make a tag provider that works with a sql database or pumps messages onto a channel
	// or whatever else you might need.
	////////////////////////////////////////////////////

	r := gologix.PathRouter{}

	// one memory based tag provider at slot 0 on the virtual "backplane"
	// a message to 2,xxx.xxx.xxx.xxx,1,0 (compact logix) will get this provider (where xxx.xxx.xxx.xxx is the ip address)
	// the message path before the IP address in the msg instruction will be different based on the actual controller
	// you're using, but the part after the IP address is what this matches
	p1 := gologix.MapTagProvider{}
	path1, err := gologix.ParsePath("1,0")
	if err != nil {
		fmt.Printf("problem parsing path. %v", err)
		os.Exit(1)
	}
	r.AddHandler(path1.Bytes(), &p1)

	// a different memory based tag provider at slot 1 on the virtual "backplane"
	p2 := gologix.MapTagProvider{}
	path2, err := gologix.ParsePath("1,1")
	if err != nil {
		fmt.Printf("problem parsing path. %v", err)
		os.Exit(1)
	}
	r.AddHandler(path2.Bytes(), &p2)

	s := gologix.NewServer(&r)
	go s.Serve()

	t := time.NewTicker(time.Second * 5)
	for {
		<-t.C
		p1.Mutex.Lock()
		log.Printf("Data 1: %v", p1.Data)
		p1.Mutex.Unlock()

		p2.Mutex.Lock()
		log.Printf("Data 2: %v", p2.Data)
		p2.Mutex.Unlock()

	}

}
