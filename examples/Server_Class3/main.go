// This example shows how to create a server that can handle incoming "class 3" cip messages.
// These are messages generated in a program from a msg instruction.
// From the PLCs perspective we are just another PLC.  When it does a write request to a specific tag we accept that data
// and when it does a read we send it the data associated with that tag on our end.
//
// Because all cip messages come in on the same port, we have to support various independent messages coming in from multiple controllers.
// otherwise you'd never be able to communicate with more than one controller per computer.  As you'll see below this is actually
// fairly simple and takes advantage of the CIP routing path in a msg instruction.  All messages come to our IP address and then we split them
// up based on the cip path after that.  The example uses 1, 0 and 1,1 which equates to backplane slot 0 and backplane slot 1 respectively.
//
// If you look at the screenshots of MSG insructions in this folder you'll see how the read and write are setup.  Note that on the
// connection tab of the msg setup the path is "gologix, 1, 0" this is because there is a generic ethernet module in the IO config
// with the same address as the computer used for the screenshots.  You can just type the IP address in here instead of "gologix".
//
// Note that you won't be able to have multiple servers on a single computer.  They bind to the EIP ports on TCP and UDP so you'll need
// to multiplex multiple connections through one program.
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
	r.Handle(path1.Bytes(), &p1)

	// set up some default tags.
	// using TagWrite() and TagRead() are treadsafe if needed.
	// otherwise you can lock p1.Mutex and manipulate p1.Data yourself
	p1.TagWrite("testtag1", int32(12345))
	p1.TagWrite("testtag2", float32(543.21))

	// a different memory based tag provider at slot 1 on the virtual "backplane" this would be "2,xxx.xxx.xxx.xxx,1,1" in the msg connection path
	p2 := gologix.MapTagProvider{}
	path2, err := gologix.ParsePath("1,1")
	if err != nil {
		fmt.Printf("problem parsing path. %v", err)
		os.Exit(1)
	}
	r.Handle(path2.Bytes(), &p2)

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
