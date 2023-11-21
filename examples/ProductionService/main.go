// This example shows a good way to continuously poll a PLC and get the data back into your main program without having to worry about
// reconnect logic, etc...
//
// Obviously this won't work for every scenario but it is a pretty solid foundation for simple things.
//
// The basic gist is that we create a channel and kick a goroutine off to pump data into that channel at a specified rate
// as the PLC is polled.  Then the main program can just receive on that channel to get fresh PLC data.
package main

import (
	"log"
	"time"

	"github.com/danomagnum/gologix"
)

type PLCPollData struct {
	Testtag1 int32   `gologix:"testtag1"`
	Testtag2 float32 `gologix:"testtag2"`
	Testdint int32   `gologix:"testdint"`
	Testint  int16   `gologix:"testint"`
}

func main() {
	// Here i'm using a struct and a multi-read, but you could use whatever makes sense for your application.
	c := make(chan PLCPollData)

	watchdog := time.NewTicker(time.Minute)

	go StartPLCComms(c)

	for {
		select {
		case newdata := <-c:
			// we got a new message so we'll reset the watchdog.
			watchdog.Reset(time.Minute)
			// now do whatever we want with the new data.
			log.Printf("multiread map has values %+v", newdata)
		case <-watchdog.C:
			// uh-oh, we didn't get new data from the PLC in longer than we were expecting.
			log.Printf("Didn't receive a message in too long!!!")

		}
	}
}

// this handles connecting and reconnecting to the PLC and then getting data forever.
// it runs as a goroutine in the background
// PollPLC could be inlined or made an inline anonymous function if desired.
func StartPLCComms(c chan PLCPollData) {
	for {
		PollPLC(c)
		time.Sleep(time.Second * 10)
		log.Print("Retrying connection in 10 seconds")
	}

}

// This is a separate function from prod_handler so i can defer() the Disconnect.  This makes it cleaner since now
// I don't have to worry about where we return from - we'll still clean up the connection resources properly.
func PollPLC(c chan PLCPollData) {

	// connect
	// Connect to the PLC.  You can hard code the address as shown or make it a parameter or something.
	client := gologix.NewClient("localhost")
	err := client.Connect()
	if err != nil {
		log.Printf("Error opening client: %v", err)
		return
	}
	defer client.Disconnect()

	// Set up the poll-rate for the PLC data read.  You can hard code this as shown or make it a parameter or something.
	pollrate := time.NewTicker(time.Second * 10)

	// loop forever until there is a problem.
	for {
		// set up the data as needed.
		m := PLCPollData{}

		// read the data
		err = client.ReadMulti(&m)
		if err != nil {
			log.Printf("error getting polled data: %v", err)
			return
		}

		// send it back to the main program.  One change that is sometimes helpful is to only send the data on the channel
		// if it has changed since the last poll.  (You'd have to do the watchdog in the main thread differently).
		c <- m

		// wait to do the read again until the next poll time
		<-pollrate.C
	}
}
