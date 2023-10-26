// this example is for a web server that accepts "push" message from a controller.
// the controller writes using a MSG instruction to this server and that value and tag will show up when doing a http GET
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/danomagnum/gologix"
	"github.com/danomagnum/gologix/cippath"
)

func main() {

	r := gologix.PathRouter{}

	// one memory based tag provider at slot 0 on the virtual "backplane"
	// a message to 2,xxx.xxx.xxx.xxx,1,0 (compact logix) will get this provider (where xxx.xxx.xxx.xxx is the ip address
	// of the computer running this server) The message path before the IP address in the msg instruction will be different
	// based on the actual controller you're using, but the part after the IP address is what this matches
	p1 := gologix.MapTagProvider{}
	path1, err := cippath.ParsePath("1,0")
	if err != nil {
		log.Printf("problem parsing path. %v", err)
		os.Exit(1)
	}
	r.Handle(path1.Bytes(), &p1)

	// create the ethernet/ip class 3 message server
	s := gologix.NewServer(&r)
	go s.Serve()

	// this is the function that will handle the web requests
	// we'll get a lock on the tag provider, marshal the data into json, and then return that
	// you could get fancier and allow retrieving specific tag keys from the data map.
	// you could also create a function that supports posting to specific keys that could be read back
	// with "CIP Data Table Read" msgs in the controller.  You'd need to be careful with types though - the json
	// library likes to use float64 for values.
	send_json := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		p1.Mutex.Lock()
		defer p1.Mutex.Unlock()
		enc := json.NewEncoder(w)
		err := enc.Encode(p1.Data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// set up a web request handler and start the server
	mux := http.NewServeMux()
	mux.HandleFunc("/", send_json)
	http.ListenAndServe("localhost:8080", mux)

}
