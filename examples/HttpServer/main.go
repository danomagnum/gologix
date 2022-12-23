package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/danomagnum/gologix"
)

func main() {

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

	// create the ethernet/ip class 3 message server
	s := gologix.NewServer(&r)
	go s.Serve()

	// this is the function that will handle the web requests
	// we'll get a lock on the tag provider, marshal the data into json, and then return that
	send_json := func(w http.ResponseWriter, req *http.Request) {
		p1.Mutex.Lock()
		defer p1.Mutex.Unlock()
		enc := json.NewEncoder(w)
		err := enc.Encode(p1.Data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
	}

	// set up a web request handler and start the server
	mux := http.NewServeMux()
	mux.HandleFunc("/", send_json)
	http.ListenAndServe("localhost:8080", mux)

}
