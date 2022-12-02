package gologix

import (
	"bytes"
	"net"
	"sync"
	"time"
)

type Client struct {
	IPAddress      string
	Path           *bytes.Buffer
	SocketTimeout  time.Duration
	sequencerValue uint16

	KnownTags map[string]KnownTag

	Mutex                  sync.Mutex
	Conn                   net.Conn
	SessionHandle          uint32
	OTNetworkConnectionID  uint32
	HeaderSequenceCounter  uint16
	Connected              bool
	ConnectionSize         int
	ConnectionSerialNumber uint16
	Context                uint64 // fun fact - rockwell PLCs don't mind being rickrolled.
}

func (client *Client) Sequencer() uint16 {
	client.sequencerValue++
	return client.sequencerValue
}

type KnownTag struct {
	Name        string
	Type        CIPType
	Class       CIPClass
	Instance    CIPInstance
	Array_Order []int
}
