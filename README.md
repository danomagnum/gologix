# gologix

gologix is a communication driver written in native Go that lets you easily read/write values from tags in Rockwell Automation ControlLogix and CompactLogix PLCs over Ethernet/IP using Go. PLCs that use CIP over Ethernet/IP are supported (ControlLogix, CompactLogix, Micro820). Models like PLC5, SLC, and MicroLogix that use PCCC instead of CIP are *not* supported.

It is modeled after pylogix with changes to make it usable in Go.

### Your First Client Program:

There are a few examples in the `examples` folder. Here is an abridged version of the `/examples/SimpleRead` example. See the actual example for a more thorough description of what is going on.

```go
package main

import (
	"fmt"
	"log"
	"github.com/danomagnum/gologix"
)

func main() {
	client := gologix.NewClient("192.168.2.241")
	err := client.Connect()
	if err != nil {
		log.Printf("Error opening client: %v", err)
		return
	}
	defer client.Disconnect()

	var tag1 int16
	err = client.Read("testint", &tag1)
	if err != nil {
		log.Printf("Error reading testint: %v", err)
		return
	}
	log.Printf("tag1 has value %d", tag1)
}
```

### Your First Server Program:

There are a few examples in the `examples` folder. Here is an abridged version of the `/examples/Server_Class3` example. See the actual example(s) for a more thorough description of what is going on. Basically, it listens to incoming MSG instructions doing CIP Data Table Writes and CIP Data Table Reads and maps the data to/from an internal Go map. You can then access the data through that map as long as you acquire the lock on it.

```go
package main

import (
	"log"
	"os"
	"time"
	"github.com/danomagnum/gologix"
)

func main() {
	r := gologix.PathRouter{}

	p1 := gologix.MapTagProvider{}
	path1, err := gologix.ParsePath("1,0")
	if err != nil {
		log.Printf("Problem parsing path: %v", err)
		os.Exit(1)
	}
	r.AddHandler(path1.Bytes(), &p1)

	s := gologix.NewServer(&r)
	go s.Serve()

	t := time.NewTicker(time.Second * 5)
	for {
		<-t.C
		p1.Mutex.Lock()
		log.Printf("Data 1: %v", p1.Data)
		p1.Mutex.Unlock()
	}
}
```

### Canned Functions

There is a `canned` package that can be used for common features such as reading controller fault codes or the status of forces. Look at the `/examples/Canned` directory to see how to use these. You can also use the code in `/canned/` as good examples of how to do particular things. Pull requests for extensions to the `canned` package are welcome (as are all pull requests).

### Other Features

- Can behave as a class 1 or class 3 server, allowing push messages from a PLC (class 3 via MSG instruction) or implicit messaging (class 1). See the *server examples*.
- You can read/write multiple tags at once by defining a struct with each field tagged with `gologix:"tagname"`. See `MultiRead` in the examples directory.
- To read multiple items from an array, pass a slice to the `Read` method.
- To read more than one arbitrary tag at once, use the `ReadList` method. The first parameter is a slice of tag names, and the second parameter is a slice of each tag's type.
- You can read UDTs if you define an equivalent struct to blit the data into. Arrays of UDTs also work (see limitation below about UDTs with packed bools).

There is also a `Server` type that lets you receive MSG instructions from the controller. See "Server" in the examples folder. It currently handles reads and writes of atomic data types (SINT, INT, DINT, REAL). You could use this to create a "push" mechanism instead of having to poll the controller for data changes.

### Limitations

- You cannot write multiple items from an array at once yet, but you can do them piecewise if needed.
- You can write to BOOL tags but NOT to bits of integers yet (e.g., "MyBool" is OK, but "MyDint.3" is NOT). You can read from either just fine. A "write with mask" implementation may be needed for this.
- No UDTs or arrays in the server yet.

## License

This project is licensed under the MIT license.

## Acknowledgements

- pylogix
- go-ethernet-ip
