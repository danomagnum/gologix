package gologix

import (
	"bytes"
	"log"
	"log/slog"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// this is an interface that matches the parts of the default go library log.Logger
// that are used in this library.  You can pass an instance of it to a Client to redirect the
// logs.  Or you can set it to nil to not log
type Logger interface {
	Printf(format string, v ...any)
}

type sLogger interface {
	Printf(format string, v ...any)
	Debug(format string, v ...any)
}

// you have to change this read sequencer every time you make a new tag request.  If you don't, you
// won't get an error but it will return the last value you requested again.
// You don't even have to keep incrementing it.  just going back and forth between 1 and 0 works OK.
// Use sequencer() instead of accessing this directly to achieve that.
var sequenceValue uint32 = 0

func init() {
	rand.Seed(time.Now().UnixMicro())
	sequenceValue = rand.Uint32()
}
func sequencer() uint32 {
	return atomic.AddUint32(&sequenceValue, 1)
}

// Defines parameters from the host CIP device that the client will connect to
type Controller struct {
	IpAddress string
	Port      uint   // Default CIP port = 44818
	VendorId  uint16 // VendorID from ODVA. Default = 0x9999 to prevent conflicts with existing vendors

	// path to the controller as a byte slice.
	// typically 1, 0 where 1 is the back plane and 0 is the slot
	// The data in the path should be similar to how you set it up in a msg instruction.
	// ex: 1, 0 where 1 -> back plane, 0 -> slot 0, etc...
	// but it has to be formatted properly as bytes (there are header bytes, etc for each portion of the path)
	// you can use the Serialize function to generate this or the GeneratePath function if it's a simpler path
	Path *bytes.Buffer
}

// Client is the main class for reading and writing tags in the PLC.
// You probably want to create a new client using NewClient() instead of instantiating
// the struct directly.
type Client struct {
	Controller Controller

	SerialNumber uint32 // serial number for the client
	VendorId     uint16 // vendor id for the client as determined from ODVA

	// Used for the keepalive messages.
	SocketTimeout time.Duration

	// Set to true to allow auto-connects on reads and writes without having to call Connect() yourself.
	AutoConnect bool

	// Monitor status of PLC and keeps connection alive
	KeepAlive          bool
	KeepAliveProps     []CIPAttribute
	KeepAliveFrequency time.Duration
	KeepAliveRunning   bool

	RPI time.Duration // Request Packet Interval

	// this keeps track of what tags are in the controller.
	// it maps tag names to a struct which has, among other things, the instance ID and class
	// which can be used to read the tag more efficiently than sending the ascii tag name to the
	// controller.  If you don't want to use this, set SocketTimeout to 0 and never call ListAllTags
	KnownTags map[string]KnownTag

	mutex                  sync.Mutex
	conn                   net.Conn
	SessionHandle          uint32
	OTNetworkConnectionID  uint32
	HeaderSequenceCounter  uint16
	Connected              bool
	ConnectionSize         uint16
	ConnectionSerialNumber uint16
	Context                uint64 // fun fact - rockwell PLCs don't mind being rick rolled.
	sequenceNumber         atomic.Uint32

	cancel_keepalive chan struct{}

	// this just lets us not have to re-process tag strings.
	ioi_cache map[string]*tagIOI

	// Replace this to capture logs
	Logger  Logger
	SLogger *slog.Logger
}

// Create a client with reasonable defaults for the given ip address.
//
// Before using the client, you will probably want to call Connect().
// After connecting be sure to call disconnect() when you are done with the client.  Probably a good place for a defer.
//
// Default path is back plane, slot 0.  For devices that aren't in a rack and aren't control or compact logix,
// such as the micro800 series or io modules, etc...  you probably want to change the path to []byte{}
// after creating the client with this function.
func NewClient(ip string) *Client {
	// default path is back plane -> slot 0
	path, err := ParsePath("1,0")
	if err != nil {
		log.Panicf("this should not have failed since the path is hardcoded.  problem with path. %v", err)
	}
	controller := Controller{
		IpAddress: ip,
		Port:      portDefault,
		Path:      path,
	}
	return &Client{
		Controller:         controller,
		VendorId:           vendorIdDefault,
		ConnectionSize:     connSizeLargeDefault,
		AutoConnect:        true,
		KeepAlive:          true,
		KeepAliveFrequency: time.Second * 30,
		KeepAliveProps:     []CIPAttribute{1, 2, 3, 4, 10},
		RPI:                rpiDefault,
		SocketTimeout:      socketTimeoutDefault,
		KnownTags:          make(map[string]KnownTag),
		ioi_cache:          make(map[string]*tagIOI),
		Logger:             log.Default(),
		SLogger:            slog.Default(),
	}

}

// This type documents a tag once it is returned with a list call.
type KnownTag struct {
	Name        string
	Info        TagInfo
	Instance    CIPInstance
	Array_Order []int
	UDT         *UDTDescriptor
	Parent      *KnownTag
}

func (t KnownTag) Bytes() []byte {
	if t.Parent == nil {
		ins := CIPInstance(t.Instance)
		b := bytes.Buffer{}
		b.Write(CipObject_Symbol.Bytes()) // 0x20 0x6B
		b.Write(ins.Bytes())
		return b.Bytes()
	} else {
		ins := CIPInstance(t.Instance)
		b := bytes.NewBuffer(t.Parent.Bytes())
		b.Write(CipObject_Symbol.Bytes()) // 0x20 0x6B
		b.Write(ins.Bytes())
		return b.Bytes()

	}
}

func (t KnownTag) Len() int {
	return len(t.Bytes())
}
