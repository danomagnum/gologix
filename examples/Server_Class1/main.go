// This example shows how to set up a class 1 IO connection.  To support multiple connections you should use the "Ethernet Bridge" module
// in the IO tree.  Then you should add one "CIP Mdoule" per virtual IO rack position.
//
// I think multiple readers should work.  Multiple writers would also appear to work but they would step on each other.
//
// You should create your own class that fulfills the TagProvider interface with the IORead and IOWrite methods completed where you handle the
// marshaling and unmarshaling of data properly.
//
// I think you should be able to have class 3 tag providers AND class 1 tag providers at the same time for the same path, BUT you'll have to
// combine their logic into a single class since the router will resolve all messages to the same place.  For this reason it might be easiest
// to keep class 3 tag providers and class 1 tag providers segregated to different "slots" on the "backplane"
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

type InStr struct {
	Data  [9]byte
	Count byte
}
type OutStr struct {
	Data [10]byte
}

func main() {

	////////////////////////////////////////////////////
	// First we set up the tag providers.
	//
	// Each one will have a path and an object that fulfills the gologix.TagProvider interface
	// We set those up and then pass them to the Router object.
	// here we're using the build in io tag provider which just has 10 bytes of inputs and 10 bytes of outputs
	//
	////////////////////////////////////////////////////

	r := gologix.PathRouter{}

	// define the Input and Output instances.  (Input and output here is from the plc's perspective)
	inInstance := InStr{}
	outInstance := OutStr{}

	// an IO handler in slot 2
	//p3 := gologix.IOProvider[InStr, OutStr]{}
	p3 := gologix.IOProvider[InStr, OutStr]{
		In:  &inInstance,
		Out: &outInstance,
	}
	path3, err := gologix.ParsePath("1,2")
	if err != nil {
		fmt.Printf("problem parsing path. %v", err)
		os.Exit(1)
	}
	r.AddHandler(path3.Bytes(), &p3)

	s := gologix.NewServer(&r)
	go s.Serve()

	t := time.NewTicker(time.Second)

	for {
		<-t.C
		inInstance.Count++
		log.Printf("PLC Input: %v", inInstance)
		log.Printf("PLC Output: %v", outInstance)
	}

}
