package gologix

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type Packing interface {
	Align(t reflect.Type) int
	Order() binary.ByteOrder
}

type CIPPack struct {
}

func (p CIPPack) Align(t reflect.Type) int {
	a := t.Align()
	return a
}

func (p CIPPack) Order() binary.ByteOrder {
	return binary.LittleEndian
}

func pack(w io.Writer, p Packing, data any) int {

	// keep track of how many bytes we've written.  This is so we can correct field alignment with padding bytes if needed
	pos := 0

	// bitpos and bitpack are for packing bits into bytes.  bitpos is the position in the byte and bitpack is the packed bits that
	// haven't been written to w yet.
	bitpos := 0
	bitpack := byte(0)

	// start reflecting and loop through the fields of the struct
	refType := reflect.TypeOf(data)
	refVal := reflect.ValueOf(data)
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		a := p.Align(field.Type)
		t := field.Tag.Get("pack")
		s := int(field.Type.Size())
		fmt.Printf("pos: %d, name: %s, align: %d \n", pos, field.Name, a)
		k := refVal.Field(i).Kind()

		// if there isn't a nopack tag on the field, we need to check for bools that need combined into bytes
		if t != "nopack" {
			// there are two conditions where we pack bits.  Multiple bools in series or a bool array.
			switch k {
			case reflect.Array:
				// we have an array without "nopack".  Check if it is a bool array
				arr := refVal.Field(i)
				at := arr.Type().Elem()
				if at.Kind() == reflect.Bool {
					l := arr.Len()
					for ai := 0; ai < l; ai++ {
						bval := arr.Index(ai).Bool()
						ival := byte(0)
						if bval {
							ival = 1
						}
						bitpack = bitpack | (ival << bitpos)
						bitpos++
						// when we have a full byte, flush it.
						if bitpos >= 8 {
							w.Write([]byte{bitpack})
							bitpos = 0
							bitpack = 0
							pos += 1
						}
					}
					continue

				}
			case reflect.Bool:
				// try to pack bools
				bval := refVal.Field(i).Bool()
				ival := byte(0)
				if bval {
					ival = 1
				}
				bitpack = bitpack | (ival << bitpos)
				bitpos++
				// when we have a full byte, flush it.
				if bitpos >= 8 {
					w.Write([]byte{bitpack})
					bitpos = 0
					bitpack = 0
					pos += 1
				}
				continue

			}

		}

		// we don't have a packable bool.  First thing we need to do is check whether there are some packed bools that still need flushed out.
		if bitpos > 0 {
			// we have at least one bit that needs flushed.
			w.Write([]byte{bitpack})
			bitpos = 0
			bitpack = 0
			pos += 1
			fmt.Printf("flushed bit packed byte. pos is now %d\n", pos)
		}

		// make sure we are writing the new data for this field to the properly aligned byte
		rem := a - (pos % a)
		if rem < a && rem > 0 {
			// need paddding bits
			fmt.Printf("pad: %d", rem)
			pad := make([]byte, rem)
			w.Write(pad)
			pos += rem
		}

		// finally, if the field is some sub-structure, recurse.  Otherwise we will write the data out
		if k != reflect.Struct {
			binary.Write(w, p.Order(), refVal.Field(i).Interface())
		} else {
			s = pack(w, p, refVal.Field(i).Interface())
		}
		pos += s
	}
	// Last thing we need to do is check whether there are some packed bools that still need flushed out.
	if bitpos > 0 {
		// we have at least one bit that needs flushed.
		w.Write([]byte{bitpack})
		bitpos = 0
		bitpack = 0
		pos += 1
	}

	fmt.Printf("wrote struct of size %d\n", pos)
	return pos
}
