package ciptype

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
)

type CIPType byte

// Go native types that correspond to logix types
// I'm not sure whether having interface here makes sense.
// On the one hand, we need to support composite types, but on the other this lets it accept anything
// which doesn't seem right.
type GoLogixTypes interface {
	bool | byte | uint16 | int16 | uint32 | int32 | uint64 | int64 | float32 | float64 | string
}

// return the CIPType that corresponds to go type of variable T
func GoVarToCIPType(T any) CIPType {
	switch T.(type) {
	case bool:
		return BOOL
	case byte:
		return BYTE
	case uint16:
		return UINT
	case int16:
		return INT
	case uint32:
		return UDINT
	case int32:
		return DINT
	case uint64:
		return LWORD
	case int64:
		return LINT
	case float32:
		return REAL
	case float64:
		return LREAL
	case string:
		return STRING
	case []byte:
		return BYTE
	case []uint16:
		return UINT
	case []int16:
		return INT
	case []uint32:
		return UDINT
	case []int32:
		return DINT
	case []uint64:
		return LWORD
	case []int64:
		return LINT
	case []float32:
		return REAL
	case []float64:
		return LREAL
	case []string:
		return STRING
	case interface{}:
		return Struct
	}
	return Unknown
}

const (
	Unknown         CIPType = 0x00
	Struct          CIPType = 0xA0 // also used for strings.  Not sure what's up with CIPTypeSTRING
	UTIME           CIPType = 0xC0
	BOOL            CIPType = 0xC1
	SINT            CIPType = 0xC2
	INT             CIPType = 0xC3
	DINT            CIPType = 0xC4
	LINT            CIPType = 0xC5
	USINT           CIPType = 0xC6
	UINT            CIPType = 0xC7
	UDINT           CIPType = 0xC8
	ULINT           CIPType = 0xC9
	REAL            CIPType = 0xCA
	LREAL           CIPType = 0xCB
	STIME           CIPType = 0xCC
	DATE            CIPType = 0xCD
	TIMEOFDAY       CIPType = 0xCE
	DATETIME        CIPType = 0xCF
	STRING_UNKNOWN  CIPType = 0xD0
	BYTE            CIPType = 0xD1 // 8 bits packed into one byte
	WORD            CIPType = 0xD2
	DWORD           CIPType = 0xD3
	LWORD           CIPType = 0xD4
	STRING_UNKNOWN2 CIPType = 0xD5
	FTIME           CIPType = 0xD6
	LTIME           CIPType = 0xD7
	ITIME           CIPType = 0xD8
	STRING_UNKNOWN3 CIPType = 0xD9
	STRING_SHORT    CIPType = 0xDA
	TIMEOFDAY2      CIPType = 0xDB
	EPATH           CIPType = 0xDC
	ENGUNIT         CIPType = 0xDD

	//  Strings actually come accross as 0xA0 = CIPTypeStruct.
	//In this library we're using this as kind of a flag to keep track of whether
	// a structure is a normal logix string or not.
	STRING CIPType = 0xFF
)

// return the size in bytes of the data structure
func (c CIPType) Size() int {
	switch c {
	case Unknown:
		return 0
	case Struct:
		return 88
	case UTIME:
		return 8
	case BOOL:
		return 1
	case BYTE:
		return 1
	case SINT:
		return 1
	case INT:
		return 2
	case DINT:
		return 4
	case LINT:
		return 8
	case USINT:
		return 1
	case UINT:
		return 2
	case UDINT:
		return 4
	case ULINT:
		return 8
	case LWORD:
		return 8
	case REAL:
		return 4
	case LREAL:
		return 8
	case WORD:
		return 2
	case DWORD:
		return 4
	case DATE:
		return 2
	case TIMEOFDAY:
		return 6
	case DATETIME:
		return 0 //?
	case STRING:
		return 1
	default:
		return 0
	}
}

// return a buffer that can hold the data structure
func (c CIPType) NewBuffer() *[]byte {
	buf := make([]byte, c.Size())
	return &buf
}

// human readable version of the cip type for printing.
func (c CIPType) String() string {
	switch c {
	case Unknown:
		return "0x00 - Unknown"
	case Struct:
		return "0xA0 - Struct"
	case BOOL:
		return "0xC1 - BOOL"
	case BYTE:
		return "0xD1 - BYTE"
	case SINT:
		return "0xC2 - SINT"
	case INT:
		return "0xC3 - INT"
	case DINT:
		return "0xC4 - DINT"
	case LINT:
		return "0xC5 - LINT"
	case USINT:
		return "0xC6 - USINT"
	case UINT:
		return "0xC7 - UINT"
	case UDINT:
		return "0xC8 - UDINT"
	case LWORD:
		return "0xC9 - LWORD"
	case REAL:
		return "0xCA - REAL"
	case LREAL:
		return "0xCB - LREAL"
	case WORD:
		return "0xD2 - WORD"
	case DWORD:
		return "0xD3 - DWORD"
	case STRING:
		return "0xFF - (gologix specific) String"
	default:
		return fmt.Sprintf("0x%2x - Unknown", byte(c))
	}
}

func (t CIPType) IsAtomic() bool {
	v := byte(t)
	return v <= 254
}
func (t CIPType) ReadValue(r io.Reader) (any, error) {
	return ReadValue(t, r)
}

// readValue reads one unit of cip data type t into the correct go type.
// To do this it reads the needed number of bytes from r.
// It returns the value as an any so the caller will have to do a cast to get it back
func ReadValue(t CIPType, r io.Reader) (any, error) {

	var value any
	var err error
	switch t {
	case Unknown:
		return nil, fmt.Errorf("unknown type")
	case Struct:
		return nil, fmt.Errorf("don't know what to do with a struct")
	case BOOL:
		var trueval bool
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case BYTE:
		var trueval byte
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case SINT:
		var trueval byte
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case INT:
		var trueval int16
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case DINT:
		var trueval int32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case LINT:
		var trueval int64
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case USINT:
		var trueval uint8
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case UINT:
		var trueval uint16
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case UDINT:
		var trueval uint32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case LWORD:
		var trueval uint64
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case REAL:
		var trueval float32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case LREAL:
		var trueval float64
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case WORD:
		var trueval uint16
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case DWORD:
		var trueval uint32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case STRING:
		var trueval [86]byte
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	default:
		return nil, fmt.Errorf("default (unknown) type %d", t)
		//panic(fmt.Sprintf("Default type %d", t))

	}
	if err != nil {
		log.Printf("Problem reading %s as one unit of %T. %v", t, value, err)
	}
	//log.Printf("type %v. value %v", t, value)
	return value, nil
}

// after reading a value v from the controller, you can get a bit from it with
// getBit. bitpos must be between 0 and the length of the CIPType you read.
func GetBit(t CIPType, v any, bitpos int) (bool, error) {

	var err error
	switch t {
	case Unknown:
		return false, errors.New("unknown type")
		//panic("Unknown type.")
	case Struct:
		return false, errors.New("got a struct - can't get a bit")
		//panic("Struct!")
	case BOOL:
		if bitpos == 0 {
			x, ok := v.(bool)
			if ok {
				return x, nil
			}
			err = fmt.Errorf("value was a bool, but bit %d was requested. must be 0 for bool", bitpos)
		}
	case BYTE:
		if bitpos >= 0 && bitpos < 8 {
			x, ok := v.(byte)
			if ok {
				mask := byte(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a byte, but bit %d was requested. must be 0-7 for byte", bitpos)
		}
	case SINT:
		if bitpos >= 0 && bitpos < 8 {
			x, ok := v.(byte)
			if ok {
				mask := byte(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a SINT, but bit %d was requested. must be 0-7 for SINT", bitpos)
		}
	case INT:
		if bitpos >= 0 && bitpos < 16 {
			x, ok := v.(int16)
			if ok {
				mask := int16(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was an INT, but bit %d was requested. must be 0-15 for INT", bitpos)
		}
	case DINT:
		if bitpos >= 0 && bitpos < 32 {
			x, ok := v.(int32)
			if ok {
				mask := int32(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a DINT, but bit %d was requested. must be 0-31 for DINT", bitpos)
		}
	case LINT:
		if bitpos >= 0 && bitpos < 64 {
			x, ok := v.(int64)
			if ok {
				mask := int64(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a LINT, but bit %d was requested. must be 0-63 for LINT", bitpos)
		}
	case USINT:
		if bitpos >= 0 && bitpos < 8 {
			x, ok := v.(byte)
			if ok {
				mask := byte(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a USINT, but bit %d was requested. must be 0-7 for USINT", bitpos)
		}
	case UINT:
		if bitpos >= 0 && bitpos < 16 {
			x, ok := v.(uint16)
			if ok {
				mask := uint16(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was an UINT, but bit %d was requested. must be 0-15 for UINT", bitpos)
		}
	case UDINT:
		if bitpos >= 0 && bitpos < 32 {
			x, ok := v.(uint32)
			if ok {
				mask := uint32(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a UDINT, but bit %d was requested. must be 0-31 for UDINT", bitpos)
		}
	case LWORD:
		if bitpos >= 0 && bitpos < 64 {
			x, ok := v.(uint64)
			if ok {
				mask := uint64(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a LWORD, but bit %d was requested. must be 0-63 for LWORD", bitpos)
		}
	case REAL:
		err = fmt.Errorf("value was a REAL, not finding bit of real")
	case LREAL:
		err = fmt.Errorf("value was a LEAL, not finding bit of real")
	case WORD:
		if bitpos >= 0 && bitpos < 16 {
			x, ok := v.(uint16)
			if ok {
				mask := uint16(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was an WORD, but bit %d was requested. must be 0-15 for WORD", bitpos)
		}
	case DWORD:
		if bitpos >= 0 && bitpos < 32 {
			x, ok := v.(uint32)
			if ok {
				mask := uint32(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a DWORD, but bit %d was requested. must be 0-31 for DWORD", bitpos)
		}
	case STRING:
		err = fmt.Errorf("value was a STRING, not finding bit of string")
	default:
		return false, errors.New("got an unknown type. don't know how to get bit")
		//panic("Default type.")

	}
	if err != nil {
		log.Printf("%v", err)
	}
	return false, err
}
