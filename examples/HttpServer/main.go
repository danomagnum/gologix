// this example is for a web server that allows GET requests to read tags and POST requests to write them.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/danomagnum/gologix"
)

var Connections = make(map[string]*gologix.Client)

func main() {

	configfile, err := os.Open("config.json")
	if err != nil {
		log.Panicf("couldn't open config.json. %v", err)
	}
	jd := json.NewDecoder(configfile)
	err = jd.Decode(&Config)
	if err != nil {
		log.Panicf("Problem reading config.json: %v", err)
	}
	log.Printf("Config: %+v", Config)

	for _, plcconf := range Config.PLCs {
		path, err := gologix.ParsePath(plcconf.Path)
		if err != nil {
			log.Printf("problem with plc connection %s. Can't parse path. %v", plcconf.Name, err)
			continue
		}
		conn := gologix.Client{
			IPAddress: plcconf.Address,
			Path:      path,
		}

		err = conn.Connect()
		if err != nil {
			log.Printf("problem with plc connection %s. Can't connect. %v", plcconf.Name, err)
			continue
		}
		Connections[plcconf.Name] = &conn
	}

	// set up a web request handler and start the server
	mux := http.NewServeMux()
	mux.HandleFunc("/", httpreq)
	connection_addr := fmt.Sprintf("%s:%d", Config.Server.Address, Config.Server.Port)
	http.ListenAndServe(connection_addr, mux)

}

func httpreq(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request Path: %v", r.URL.Path)
	switch r.Method {
	case "GET":
		// reads
		httpread(w, r)
	case "POST":
		// writes
		httpwrite(w, r)
	default:
		// not supported
	}
}

func httpread(w http.ResponseWriter, r *http.Request) {
	conn, path, err := parsePLC(r.URL.Path)
	if err != nil {
		// problem getting connection.
		w.Write([]byte(fmt.Sprintf("Problem connecting. %v", err)))
		return
	}
	value, err := conn.Read_single(path, gologix.CIPType(0), 1)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Problem reading. %v", err)))
		return
	}
	w.Write([]byte(fmt.Sprintf("Value: %v", value)))

}

func httpwrite(w http.ResponseWriter, r *http.Request) {
	// not implemented yet.
}

func parsePLC(reqPath string) (*gologix.Client, string, error) {
	if reqPath[0] == '/' {
		reqPath = reqPath[1:]
	}
	parts := strings.Split(reqPath, "/")
	if len(parts) == 0 {
		return nil, "", fmt.Errorf("could not get PLC from %v", reqPath)
	}

	client, ok := Connections[parts[0]]
	if !ok {
		return nil, "", fmt.Errorf("unknown PLC '%v'", parts[0])
	}
	switch len(parts) {
	case 2:
		return client, parts[1], nil
	case 1:
		return client, "", nil
	default:
		return nil, "", errors.New("bad path")

	}

}
