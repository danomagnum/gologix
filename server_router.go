package gologix

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
)

// The Server uses PathRouter to resolve paths to tag providers
// Its use is similar to a mux in an HTTP server where you add endpoints with
// the .Handle() method.  Instead of an http path you use a CIP route byte slice and
// instead of a handler function you use an object that provides the TagProvider interface.
type PathRouter struct {
	Path map[string]CIPEndpoint
}

func NewRouter() *PathRouter {
	p := new(PathRouter)
	p.Path = make(map[string]CIPEndpoint)
	return p
}

func (router *PathRouter) Handle(path []byte, p CIPEndpoint) {
	if router.Path == nil {
		router.Path = make(map[string]CIPEndpoint)
	}
	router.Path[string(path)] = p
}

// find the tag provider for a given cip path
func (router *PathRouter) Resolve(path []byte) (CIPEndpoint, error) {
	tp, ok := router.Path[string(path)]
	if !ok {
		return nil, fmt.Errorf("path %v not recognized", path)
	}
	return tp, nil
}

// This interface specifies all the needed methods to handle incoming CIP messages.
// currently supports Class1 IO messages and Class3 tag read/write messages.
// if a type only handles some subset, it should return an error for those methods
type CIPEndpoint interface {
	// These functions are called when a cip service attempts to use the write or read services.
	TagRead(tag string, qty int16) (any, error)
	TagWrite(tag string, value any) error

	// IORead is called every time the RPI triggers for an Input (from the PLC's perspective) IO message.
	// It should return the serialized bytes to send to the controller.
	IORead() ([]byte, error)

	// IOWrite is called every time a class 1 IO message comes in.  The CIP items that came in with the message
	// are passed as arguments.  You should check that you have the correct number of items (should be 2?) and
	// that they are the correct type.
	//
	// items[1] has the actual write payload and it should be a connectedData item.
	// it contains the following in the data section which you can deserialize with items[1].deserialize(xxx):
	// SequenceCounter uint32
	// Header uint16
	// Payload [items[1].Header.Length - 6]byte
	IOWrite(items []CIPItem) error
}

// This is a generic tag provider that can handle bi-directional class 3 tag reads and writes.
// If a tag is written that does not exist, that will create it.
// if a tag is read that does not exist, that will result in an error
// it does not handle IO messages.
//
// The built-in MapTagProvider type also only provides rudimentary tag access.
// It doesn't support addressing arrays directly
//
//	It interprets "testtag[3]" as a tag with a literal "[3]" as a string at the end of the map key.
//	Same thing for nested UDT tags - it interprets the dots as literals in the map keys.
//
// So if you need those for testing you'll have to kind of fake out the tag mapping on the server end by
// creating an individual entry in the map for each nested tag with the key being the full access path.
type MapTagProvider struct {
	Mutex sync.Mutex
	Data  map[string]any
}

func (p *MapTagProvider) IORead() ([]byte, error) {
	return nil, errors.New("not implemented")
}
func (p *MapTagProvider) IOWrite(items []CIPItem) error {
	return errors.New("not implemented")
}

// this is a thread-safe way to get the value for a tag.
func (p *MapTagProvider) TagRead(tag string, qty int16) (any, error) {
	log.Printf("Trying to read %v from MapTagProvider", tag)
	tag = strings.ToLower(tag)
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	if p.Data == nil {
		p.Data = make(map[string]any)
	}

	val, ok := p.Data[tag]
	if !ok {
		return nil, fmt.Errorf("tag %v not in map", tag)
	}

	t := reflect.ValueOf(val)

	if t.Kind() == reflect.Slice {
		if int(qty) <= t.Len() {
			values := reflect.Indirect(t)
			v := values.Slice(0, int(qty))
			return v.Interface(), nil
		}
		return nil, fmt.Errorf("too many elements requested %v > %v", qty, t.Len())
	} else {
		log.Printf("not a slice")
	}

	return val, nil
}

// this is a thread-safe way to write a value to a tag.
func (p *MapTagProvider) TagWrite(tag string, value any) error {
	log.Printf("Trying to set %v=%v from MapTagProvider", tag, value)

	tag = strings.ToLower(tag)
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	if p.Data == nil {
		p.Data = make(map[string]any)
	}
	p.Data[tag] = value
	return nil
}
