# gologix

gologix is a communication driver that lets you easily read/write values from tags in Rockwell Automation ControlLogix, and CompactLogix PLC's over Ethernet I/P using GO.  Only PLC's that are programmed with RSLogix5000/Studio5000, models like PLC5, SLC, MicroLogix, or Micro800 are *not* supported.  They use a different protocol, which I have no plans to support at this time.

It is modeled after pylogix with changes to make it usable in go.

### Your First Program:

There are a few examples in the examples folder, here is an abriged version of the SimpleRead example. See the actual example for description of what is going on.

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

You can read multiple tags at once by defining a struct with each field tagged with `gologix:"tagname"`.  see MultiRead in the examples directory.

To read multiple items from an array, pass a slice to the Read method.

You can read UDTs in if you define an equivalent struct to blit the data into. Arrays of UDTs also works.

### Limitations

You cannot write multiple items from an array at once, or write a whole UDT in the current version, but you can do them piecewise if needed.

You can write to BOOL tags but NOT to bits of integers yet (ex: "MyBool" is OK, but "MyDint.3" is NOT).  You can read from either just fine.

If the UDT you're reading has bools packed in it, you'll need to use the ReadPacked() function instead of client.Read().  The plan is to eventually migrate this functionality to client.Read automatically.

## License

This project is licensed under the MIT license.

## Acknowledgements

* pylogix
* go-ethernet-ip
