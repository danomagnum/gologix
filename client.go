package gologix

import (
	"bytes"
	"log"
	"net"
	"sync"
	"time"
)

type Client struct {
	// ip address for connecting to the PLC comms module.
	IPAddress string

	// the port is always 44818
	Port string // = ":44818"

	// vendor ID should be provided by ODVA.  Since we probably don't have an
	// official ID, we'll use the hex'd version of the founding of america.
	VendorID uint16 //= 0x1776

	// serial number for the device.
	SerialNumber uint32

	// path to the controller as a byte slice.
	// The data in the path should be similar to how you set it up in a msg instruction.
	// ex: 1, 0 where 1 -> backlane, 0 -> slot 0, etc...
	// but it has to be formatted properly as bytes (there are header bytes, etc for each portion of the path)
	// you can use the Serialize function to generate this or the GeneratePath function if it's a simpler path
	Path *bytes.Buffer

	// Used for the keepalive messages.
	SocketTimeout time.Duration

	RPI time.Duration

	// you have to change this read sequencer every time you make a new tag request.  If you don't, you
	// won't get an error but it will return the last value you requested again.
	// You don't even have to keep incrementing it.  just going back and forth between 1 and 0 works OK.
	// Use Sequencer() instead of accessing this directly to achieve that.
	sequencerValue uint16

	// this keeps track of what tags are in the controller.
	// it maps tag names to a struct which has, among other things, the intance ID and class
	// which can be used to read the tag more efficiently than sending the ascii tag name to the
	// controller.  If you don't want to use this, set SocketTimeout to 0 and never call ListAllTags
	KnownTags map[string]KnownTag

	mutex                  sync.Mutex
	conn                   net.Conn
	SessionHandle          uint32
	OTNetworkConnectionID  uint32
	HeaderSequenceCounter  uint16
	Connected              bool
	ConnectionSize         int
	ConnectionSerialNumber uint16
	Context                uint64 // fun fact - rockwell PLCs don't mind being rickrolled.

	cancel_keepalive chan struct{}
}

func (client *Client) Sequencer() uint16 {
	client.sequencerValue++
	return client.sequencerValue
}

// create a client with reasonable defaults for the given ip address.
// Default path is backplane, slot 0
func NewClient(ip string) *Client {
	// default path is backplane -> slot 0
	p, err := ParsePath("1,0")
	if ioi_cache == nil {
		ioi_cache = make(map[string]*tagIOI)
	}
	if err != nil {
		log.Panicf("this should not have failed since the path is hardcoded.  problem with path. %v", err)
	}
	return &Client{
		IPAddress:      ip,
		ConnectionSize: 4000,
		Path:           p,
		Port:           ":44818",
		VendorID:       0x1776,
		SerialNumber:   42,
		RPI:            time.Millisecond * 2500,
	}

}

type KnownTag struct {
	Name        string
	Info        TagInfo
	Instance    CIPInstance
	Array_Order []int
	UDT         *UDTDescriptor
}

func (t KnownTag) Bytes() []byte {
	ins := CIPInstance(t.Instance)
	b := bytes.Buffer{}
	b.Write(cipObject_Symbol.Bytes()) // 0x20 0x6B
	b.Write(ins.Bytes())
	return b.Bytes()
}

func (t KnownTag) Len() int {
	return len(t.Bytes())
}
