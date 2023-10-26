// This example shows how to set up a class 1 IO connection.  To support multiple connections you should use the "Ethernet Bridge" module
// in the IO tree.  Then you should add one "CIP Mdoule" per virtual IO rack position.
//
// I think multiple readers should work.  Multiple writers would also appear to work but they would step on each other.
//
// You should create your own class that fulfills the TagProvider interface with the IORead and IOWrite methods completed where you handle the
// serializing and deserializing of data properly.
//
// I think you should be able to have class 3 tag providers AND class 1 tag providers at the same time for the same path, BUT you'll have to
// combine their logic into a single class since the router will resolve all messages to the same place.  For this reason it might be easiest
// to keep class 3 tag providers and class 1 tag providers segregated to different "slots" on the "backplane"
//
// Note that you won't be able to have multiple servers on a single computer.  They bind to the EIP ports on TCP and UDP so you'll need
// to multiplex multiple connections through one program.
//
// In the current inarnation of this server it doesn't matter what assembly instance IDs you select in the controller, although you could create your own
// TagProvider that changed behavior based on that.
package main

import (
	"log"
	"os"
	"time"

	"github.com/danomagnum/gologix"
	"github.com/danomagnum/gologix/cippath"
)

// these types will be the input and output data section for the io connection.
// the input/output nomenclature is from the PLC's point of view - Input goes to the PLC and output
// comes to us.
//
// the size (in bytes) of these structures has to match the size you set up in the IO tree for the IO connection.
// presumably you can also use other formats than bytes for the data type, but the sizes still have to match.
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

	// an IO handler in slot 2
	//p3 := gologix.IOProvider[InStr, OutStr]{}
	p3 := gologix.IOChannelProvider[InStr, OutStr]{}
	path3, err := cippath.ParsePath("1,2")
	if err != nil {
		log.Printf("problem parsing path. %v", err)
		os.Exit(1)
	}
	r.Handle(path3.Bytes(), &p3)

	s := gologix.NewServer(&r)
	go s.Serve()

	t := time.NewTicker(time.Second)
	data_chan := p3.GetOutputDataChannel()

	for {
		select {
		case <-t.C:
			// time to update the input data
			inInstance.Count++
			p3.SetInputData(inInstance)
		case outdat := <-data_chan:
			log.Printf("PLC Output: %v", outdat)
		}
	}

}
