# gologix

gologix is a communication driver written in native go that lets you easily read/write values from tags in Rockwell Automation ControlLogix, and CompactLogix PLC's over Ethernet I/P using GO.  Only PLC's that are programmed with RSLogix5000/Studio5000, models like PLC5, SLC, MicroLogix, or Micro800 are *not* supported.  They use a different protocol, which I have no plans to support at this time.

It is modeled after pylogix with changes to make it usable in go.

### Your First Program:

There are a few examples in the examples folder, here is an abriged version of the SimpleRead example. See the actual example for a more thorough description of what is going on.

```go
package main

import (
	"fmt"
	"gologix"
)

func main() {
	client := gologix.NewClient("192.168.2.241")
	err := client.Connect()
	if err != nil {
		fmt.Printf("Error opening client. %v", err)
		return
	}
	defer client.Disconnect()

	var tag1 int16
	err = client.Read("testint", &tag1)
	if err != nil {
		fmt.Printf("error reading testint. %v", err)
        return
	}
	fmt.Printf("tag1 has value %d", tag1)
}

```



### Other Features

You can read/write multiple tags at once by defining a struct with each field tagged with `gologix:"tagname"`.  see MultiRead in the examples directory.

To read multiple items from an array, pass a slice to the Read method.

You can read UDTs in if you define an equivalent struct to blit the data into. Arrays of UDTs also works. (see limitation below about UDTs with packed bools)


There is also a ```Server``` type that lets you recive msg instructions from the controller.  See "Server" in the examples folder.  It currently handles reads and writes of atomic data types (SINT, INT, DINT, REAL).  You could use this to create a "push" mechanism instead of having ot poll the controller for data changes.

### Limitations

You cannot write multiple items from an array at once yet, but you can do them piecewise if needed.

You can write to BOOL tags but NOT to bits of integers yet (ex: "MyBool" is OK, but "MyDint.3" is NOT).  You can read from either just fine.  I think there is a "write with mask" that I'll need to implement to do this.

If the UDT you're reading has bools packed in it, you'll need to use the ReadPacked() function instead of client.Read().  The plan is to eventually migrate this functionality to client.Read automatically.

The library currently only does large forward opens.  It also only does connected reads/writes.  At some point regular forward opens may be added.


No UDTs in the server yet.  This will eventually be implemented and that will greatly improve functionality.

## License

This project is licensed under the MIT license.

## Acknowledgements

* pylogix
* go-ethernet-ip
