package gologix

import (
	"bytes"
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

type Packable interface {
	Pack(w io.Writer) (int, error)
}

type Unpackable interface {
	Unpack(r io.Reader) (n int, err error)
}

func (p CIPPack) Align(t reflect.Type) int {
	// If a type is a struct we need to check the alignment of every field.
	// If any of the fields have an alignment of 8 (LINT, LREAL, etc...)
	// then the struct also has an alignment of 8.
	if t.Kind() != reflect.Struct {
		return t.Align()
	}
	if t.Kind() == reflect.Array {
		return p.Align(t.Elem())
	}

	a := 1
	for i := 0; i < t.NumField(); i++ {
		subfield_align := p.Align(t.Field(i).Type)
		if subfield_align > a {
			a = subfield_align
		}

	}
	return a
}

func (p CIPPack) Order() binary.ByteOrder {
	return binary.LittleEndian
}

type Serializable interface {
	Bytes() []byte
	Len() int
}

// given a list of structures, serialize them with the Bytes() method if available,
// otherwise serialize with binary.write()
func Serialize(strs ...any) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	for _, str := range strs {
		switch serializable_str := str.(type) {
		case string:
			strlen := uint32(len(serializable_str))
			err := binary.Write(b, binary.LittleEndian, strlen)
			if err != nil {
				return nil, fmt.Errorf("Problem writing string header: %v", err)
			}
			if strlen%2 == 1 {
				strlen++
			}
			b2 := make([]byte, 84)
			copy(b2, serializable_str)
			err = binary.Write(b, binary.LittleEndian, b2)
			if err != nil {
				return nil, fmt.Errorf("Problem writing string payload: %v", err)
			}
		case Serializable:
			// if the struct is serializable, we should use its Bytes() function to get its
			// representation instead of binary.Write
			_, err := b.Write(serializable_str.Bytes())
			if err != nil {
				return nil, err
			}

		case any:
			err := binary.Write(b, binary.LittleEndian, str)
			if err != nil {
				return nil, fmt.Errorf("problem writing str to buffer. %w", err)
			}
		}
	}
	return b, nil
}

func Pack(w io.Writer, p Packing, data any) (int, error) {

	switch d := data.(type) {
	case Packable:
		return d.Pack(w)
	case Serializable:
		n, err := w.Write(d.Bytes())
		if err != nil {
			return 0, nil
		}
		return n, nil
	}

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
								return pos, fmt.Errorf("problem writing bitpack to buffer. %v", err)
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
						return pos, fmt.Errorf("problem writing bitpacked byte. %v", err)
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
				return pos, fmt.Errorf("problem writing bitpacked byte. %v", err)
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
				return pos, fmt.Errorf("problem writing pad to buffer. %v", err)
			}
			pos += rem
		}

		// finally, if the field is some sub-structure, recurse.  Otherwise we will write the data out
		if k != reflect.Struct {
			err := binary.Write(w, p.Order(), refVal.Field(i).Interface())
			if err != nil {
				return pos, fmt.Errorf("problem reading sub-structure. %v", err)
			}

		} else {
			var err error
			s, err = Pack(w, p, refVal.Field(i).Interface())
			if err != nil {
				return pos, fmt.Errorf("problem packing interface: %w", err)
			}
		}
		pos += s
	}
	// Last thing we need to do is check whether there are some packed bools that still need flushed out.
	if bitpos > 0 {
		// we have at least one bit that needs flushed.
		_, err := w.Write([]byte{bitpack})
		if err != nil {
			return pos, fmt.Errorf("problem flushing bitpack, %v", err)
		}
		pos += 1
	}

	return pos, nil
}

func Unpack(r io.Reader, p Packing, data any) (n int, err error) {

	switch d := data.(type) {
	case Unpackable:
		return d.Unpack(r)
	}

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
	size, err := Pack(buf, CIPPack{}, data)
	if err != nil {
		return data, err
	}

	b := make([]byte, size)
	err = client.Read(tag, &b)
	if err != nil {
		return data, fmt.Errorf("couldn't read %s as bytes. %w", tag, err)
	}
	_, err = Unpack(bytes.NewBuffer(b), CIPPack{}, &data)
	if err != nil {
		return data, fmt.Errorf("problem unpacking from buffer. %w", err)
	}

	return data, nil

}
