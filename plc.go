package main

import (
	"net"
	"sync"
	"time"
)

var ioi_cache map[string]*IOI

type PLC struct {
	IPAddress     string
	ProcessorSlot int
	SocketTimeout time.Duration
	readSequencer uint16
	// Route

	Size                   int // 508 is the default
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
