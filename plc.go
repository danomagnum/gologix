package gologix

import (
	"net"
	"sync"
	"time"
)

type Client struct {
	IPAddress     string
	ProcessorSlot int
	SocketTimeout time.Duration
	readSequencer uint16

	KnownTags map[string]KnownTag

	// Route

	Mutex                  sync.Mutex
	Conn                   net.Conn
	SessionHandle          uint32
	OTNetworkConnectionID  uint32
	SequenceCounter        uint16
	Connected              bool
	ConnectionSize         int
	ConnectionSerialNumber uint16
	Context                uint64 // fun fact - rockwell PLCs don't mind being rickrolled.
}

type KnownTag struct {
	Name        string
	Type        CIPType
	Class       CIPClass
	Instance    CIPInstance
	Array_Order []int
}
