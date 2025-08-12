# gologix

gologix is a communication driver written in native Go that lets you easily read/write values from tags in Rockwell Automation ControlLogix and CompactLogix PLCs over Ethernet/IP using Go. PLCs that use CIP over Ethernet/IP are supported (ControlLogix, CompactLogix, Micro820). Models like PLC5, SLC, and MicroLogix that use PCCC instead of CIP are *not* supported.

It is modeled after pylogix with changes to make it usable in Go with a goal of being similar in usage patterns to modules of the go standard library.

### Client

The `client` lets you read and write tag data to a logix PLC.  

You can read/write single tags (atomic or UDT) with `Read` by providing a pointer to the equivalent go type and it will populate the data similar to the json library. You can also read multiple tags at once with `ReadMulti`, `ReadList`, and `ReadMap`.  You can likewise write with `Write`, `WritMulti`, and `WriteMap`.  If you don't know what tag you want, you can list all the tags in a controller with `ListAllTags` or programs with `ListAllPrograms`.  To do a custom message you can use `GenericCIPMessage`.  There are also a couple other features for advanced use.

There are examples for all these methods in the `examples`, `tests`, and `canned` folders. Here is an abridged version of the `/examples/SimpleRead` example. See the actual example for a more thorough description of what is going on.

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

### Server

The `Server` type lets you receive MSG instructions from the controller. or behave as an IO adapter.

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

### L5X File Support

The `l5x` package provides support for working with RSLogix5000 L5X export files. This allows you to parse project files exported from Studio 5000 to extract tag information, UDT definitions, and program structure without needing a live connection to the PLC.


## License

This project is licensed under the MIT license.

## Acknowledgements

- pylogix
- go-ethernet-ip
