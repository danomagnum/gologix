package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
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

func Pack(w io.Writer, p Packing, data any) int {

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
							_, err := w.Write([]byte{bitpack})
							if err != nil {
								//TODO: make this function return error?
								log.Printf("problem writing bitpack to buffer. %v", err)
							}
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
					_, err := w.Write([]byte{bitpack})
					if err != nil {
						log.Printf("problem writing bitpacked byte. %v", err)
					}
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
			_, err := w.Write([]byte{bitpack})
			if err != nil {
				log.Printf("problem writing bitpacked byte. %v", err)
			}
			bitpos = 0
			bitpack = 0
			pos += 1
		}

		// make sure we are writing the new data for this field to the properly aligned byte
		rem := a - (pos % a)
		if rem < a && rem > 0 {
			// need paddding bits
			pad := make([]byte, rem)
			_, err := w.Write(pad)
			if err != nil {
				log.Printf("problem writing pad to buffer. %v", err)
			}
			pos += rem
		}

		// finally, if the field is some sub-structure, recurse.  Otherwise we will write the data out
		if k != reflect.Struct {
			err := binary.Write(w, p.Order(), refVal.Field(i).Interface())
			if err != nil {
				log.Printf("problem reading sub-structure. %v", err)
			}

		} else {
			s = Pack(w, p, refVal.Field(i).Interface())
		}
		pos += s
	}
	// Last thing we need to do is check whether there are some packed bools that still need flushed out.
	if bitpos > 0 {
		// we have at least one bit that needs flushed.
		_, err := w.Write([]byte{bitpack})
		if err != nil {
			log.Printf("problem flushing bitpack, %v", err)
		}
		pos += 1
	}

	return pos
}

func Unpack(r io.Reader, p Packing, data any) (n int, err error) {

	// bitpos and bitpack are for packing bits into bytes.  bitpos is the position in the byte and bitpack is the packed bits that
	// haven't been written to w yet.
	bitpos := 0
	bitpack := byte(0)

	// start reflecting and loop through the fields of the struct
	refVal := reflect.ValueOf(data)
	if refVal.Kind() == reflect.Ptr {
		refVal = reflect.ValueOf(data).Elem()
	}
	refType := refVal.Type()
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		a := p.Align(field.Type)
		t := field.Tag.Get("pack")
		s := int(field.Type.Size())
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
						if bitpos == 0 {
							br := []byte{0}
							_, err = r.Read(br)
							if err != nil {
								return n, fmt.Errorf("problem reading bool. %w", err)
							}
							bitpack = br[0]
							n += 1
						}
						val := bitpack & (1 << bitpos)
						bval := val != 0
						arr.Index(ai).SetBool(bval)
						bitpos++
						// when we have a full byte, flush it.
						if bitpos >= 8 {
							bitpos = 0
						}
					}
					continue

				}
			case reflect.Bool:
				// try to pack bools
				if bitpos == 0 {
					br := []byte{0}
					_, err = r.Read(br)
					if err != nil {
						return n, fmt.Errorf("problem reading packed bool. %w", err)
					}
					bitpack = br[0]
					n += 1
				}
				val := bitpack & (1 << bitpos)
				bval := val != 0
				refVal.Field(i).SetBool(bval)
				bitpos++
				// when we have a full byte, flush it.
				if bitpos >= 8 {
					bitpos = 0
				}
				continue

			}

		}

		// we don't have a packable bool.  First thing we need to do is check whether there are some packed bools that still need flushed out.
		if bitpos > 0 {
			// we have at least one bit that needs flushed.
			bitpos = 0
		}

		// make sure we are writing the new data for this field to the properly aligned byte
		rem := a - (n % a)
		if rem < a && rem > 0 {
			// need paddding bits
			pad := make([]byte, rem)
			_, err = r.Read(pad)
			if err != nil {
				return
			}
			n += rem
		}

		// finally, if the field is some sub-structure, recurse.  Otherwise we will write the data out
		if k != reflect.Struct {
			//binary.Read(r, p.Order(), refVal.Field(i).Interface())
			err = binary.Read(r, p.Order(), refVal.Field(i).Addr().Interface())
			if err != nil {
				return
			}
		} else {
			val := refVal.Field(i).Addr().Interface()
			s, err = Unpack(r, p, val)
			if err != nil {
				return
			}
		}
		n += s
	}
	// Last thing we need to do is check whether there are some packed bools that still need flushed out.
	return
}

func ReadPacked[T any](client *Client, tag string) (T, error) {
	var data T
	buf := new(bytes.Buffer)
	size := Pack(buf, CIPPack{}, data)

	b := make([]byte, size)
	err := client.Read(tag, &b)
	if err != nil {
		return data, fmt.Errorf("couldn't read %s as bytes. %w", tag, err)
	}
	_, err = Unpack(bytes.NewBuffer(b), CIPPack{}, &data)
	if err != nil {
		return data, fmt.Errorf("problem unpacking from buffer. %w", err)
	}

	return data, nil

}
