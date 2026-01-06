package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	"github.com/npat-efault/crc16"
)

type cipPack struct {
}

type Packable interface {
	Pack(w io.Writer) (int, error)
}

type Unpackable interface {
	Unpack(r io.Reader) (n int, err error)
}

func (p cipPack) Align(t reflect.Type) int {
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

func (p cipPack) Order() binary.ByteOrder {
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
				return nil, fmt.Errorf("problem writing string header: %w", err)
			}
			if strlen%2 == 1 {
				strlen++
			}
			b2 := make([]byte, 84)
			copy(b2, serializable_str)
			err = binary.Write(b, binary.LittleEndian, b2)
			if err != nil {
				return nil, fmt.Errorf("problem writing string payload: %w", err)
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

// serialize data into w appropriately for CIP messaging
// obeys alignment and padding rules
func Pack(w io.Writer, data any) (int, error) {
	p := cipPack{}

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
			s, err = Pack(w, refVal.Field(i).Interface())
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

// deserialize data from r appropriately for CIP messaging
// obeys alignment and padding rules
func Unpack(r io.Reader, data any) (n int, err error) {
	p := cipPack{}

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
	k := refType.Kind()
	switch k {
	case reflect.Slice:
		// we have a slice of structs.  We need to unpack each one individually
		l := refVal.Len()
		for i := 0; i < l; i++ {
			s, err := Unpack(r, refVal.Index(i).Addr().Interface())
			if err != nil {
				return n, fmt.Errorf("problem unpacking slice element %d: %w", i, err)
			}
			n += s
			align := p.Align(refType.Elem())
			if s%align != 0 {
				s, err = r.Read(make([]byte, align-s%align))
				if err != nil {
					return n, fmt.Errorf("problem reading slice element padding: %w", err)
				}
				n += s
			}
		}
		return n, nil

	case reflect.Struct:
		// continue on

	}
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
			s, err = Unpack(r, val)
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
	size, err := Pack(buf, data)
	if err != nil {
		return data, err
	}

	b := make([]byte, size)
	err = client.Read(tag, &b)
	if err != nil {
		return data, fmt.Errorf("couldn't read %s as bytes. %w", tag, err)
	}
	_, err = Unpack(bytes.NewBuffer(b), &data)
	if err != nil {
		return data, fmt.Errorf("problem unpacking from buffer. %w", err)
	}

	return data, nil

}

type KnownType interface {
	TypeAbbr() (string, uint16)
}

// perform type encoding per TypeEncode_CIPRW.pdf from the rockwell site.  Also returns the abbreviated type ID
func TypeEncode(data any) (string, uint16, error) {
	// TODO: does this whole thing break if we have a struct with bools, what with their ZZZZZZ prefixed values and all?
	//       I suspect it does.  The UDT type definitions won gold at the bad idea olympics.

	Abbreviable, ok := data.(KnownType)
	if ok {
		s, t := Abbreviable.TypeAbbr()
		return s, t, nil
	}

	encoded := ""
	// bitpos and bitpack are for packing bits into bytes.  bitpos is the position in the byte and bitpack is the packed bits that
	// haven't been written to w yet.
	bitpos := 0

	// start reflecting and loop through the fields of the struct
	refType := reflect.TypeOf(data)
	refVal := reflect.ValueOf(data)

	// start with the structure name
	encoded = refType.Name()
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		k := refVal.Field(i).Kind()

		// there are two conditions where we pack bits.  Multiple bools in series or a bool array.
		switch k {
		case reflect.Array:
			// we have an array without "nopack".  Check if it is a bool array
			arr := refVal.Field(i)
			var at string
			if arr.Type().Elem().Kind() != reflect.Struct {
				at, _ = goTypeToLogixTypeName(field.Type.Elem())
			} else {
				at, _, _ = TypeEncode(arr.Index(0).Interface())
			}
			encoded += fmt.Sprintf(",%s[%d]", at, arr.Len())
			continue
		case reflect.Bool:
			// try to pack bools
			if bitpos == 0 {
				encoded += ",SINT"
			}
			bitpos++
			// when we have a full byte, flush it.
			if bitpos >= 8 {
				bitpos = 0

			}
			continue

		}

		// we don't have a packable bool.  First thing we need to do is check whether there are some packed bools that still need flushed out.
		bitpos = 0

		// finally, if the field is some sub-structure, recurse.  Otherwise we will write the data out
		if k != reflect.Struct {
			n, _ := goTypeToLogixTypeName(refVal.Field(i).Type())
			encoded += fmt.Sprintf(",%s", n)
		} else {
			str_text, _, err := TypeEncode(refVal.Field(i).Interface())
			if err != nil {
				return "", 0, fmt.Errorf("problem encoding sub-structure: %w", err)
			}
			encoded += fmt.Sprintf(",%s", str_text)
		}
	}

	crc := crc16.Checksum(crc_conf, []byte(encoded))

	return encoded, crc, nil
}

func goTypeToLogixTypeName(t reflect.Type) (string, error) {
	switch t.Name() {
	case "int8":
		return "SINT", nil
	case "int16":
		return "INT", nil
	case "int32":
		return "DINT", nil
	case "int64":
		return "LINT", nil
	case "uint8":
		return "SINT", nil
	case "uint16":
		return "UINT", nil
	case "uint32":
		return "UDINT", nil
	case "uint64":
		return "ULINT", nil
	case "float32":
		return "REAL", nil
	case "float64":
		return "LREAL", nil
	}
	return "", nil
}

// this is the confirguration for the CRC16 checksum used in the abbreviated type ID calculation for
// UDT types..
var crc_conf = &crc16.Conf{
	Poly:   0x8005,
	BitRev: true,
	BigEnd: false,
}
