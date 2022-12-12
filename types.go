package gologix

import (
	"encoding/binary"
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
	bool | byte | uint16 | int16 | uint32 | int32 | uint64 | int64 | float32 | float64 | string | interface{}
}
type ComparableGoLogixTypes interface {
	bool | byte | uint16 | int16 | uint32 | int32 | uint64 | int64 | float32 | float64
}

// return the CIPType that corresponds to go type of variable T
func GoVarToCIPType(T any) CIPType {
	switch T.(type) {
	case bool:
		return CIPTypeBOOL
	case byte:
		return CIPTypeBYTE
	case uint16:
		return CIPTypeUINT
	case int16:
		return CIPTypeINT
	case uint32:
		return CIPTypeUDINT
	case int32:
		return CIPTypeDINT
	case uint64:
		return CIPTypeLWORD
	case int64:
		return CIPTypeLINT
	case float32:
		return CIPTypeREAL
	case float64:
		return CIPTypeLREAL
	case string:
		return CIPTypeSTRING
	case interface{}:
		return CIPTypeStruct
	}
	return CIPTypeUnknown
}

const (
	CIPTypeUnknown CIPType = 0x00
	CIPTypeStruct  CIPType = 0xA0 // also used for strings.  Not sure what's up with CIPTypeSTRING
	CIPTypeBOOL    CIPType = 0xC1
	CIPTypeBYTE    CIPType = 0xD1 // 8 bits packed into one byte
	CIPTypeSINT    CIPType = 0xC2
	CIPTypeINT     CIPType = 0xC3
	CIPTypeDINT    CIPType = 0xC4
	CIPTypeLINT    CIPType = 0xC5
	CIPTypeUSINT   CIPType = 0xC6
	CIPTypeUINT    CIPType = 0xC7
	CIPTypeUDINT   CIPType = 0xC8
	CIPTypeLWORD   CIPType = 0xC9
	CIPTypeREAL    CIPType = 0xCA
	CIPTypeLREAL   CIPType = 0xCB
	CIPTypeWORD    CIPType = 0xD2
	CIPTypeDWORD   CIPType = 0xD3

	// As far as I can tell CIPTypeSTRING isn't actually used in the controllers. Strings actually come
	// accross as 0xA0 = CIPTypeStruct.  In this library we're using this as kind of a flag to keep track of whether
	// a structure is a string or not.
	CIPTypeSTRING CIPType = 0xDA
)

// return the size in bytes of the data structure
func (c CIPType) Size() int {
	switch c {
	case CIPTypeUnknown:
		return 0
	case CIPTypeStruct:
		return 88
	case CIPTypeBOOL:
		return 1
	case CIPTypeBYTE:
		return 1
	case CIPTypeSINT:
		return 1
	case CIPTypeINT:
		return 2
	case CIPTypeDINT:
		return 4
	case CIPTypeLINT:
		return 8
	case CIPTypeUSINT:
		return 1
	case CIPTypeUINT:
		return 2
	case CIPTypeUDINT:
		return 4
	case CIPTypeLWORD:
		return 8
	case CIPTypeREAL:
		return 4
	case CIPTypeLREAL:
		return 8
	case CIPTypeWORD:
		return 2
	case CIPTypeDWORD:
		return 4
	case CIPTypeSTRING:
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
	case CIPTypeUnknown:
		return "0x00 - Unknown"
	case CIPTypeStruct:
		return "0xA0 - Struct"
	case CIPTypeBOOL:
		return "0xC1 - BOOL"
	case CIPTypeBYTE:
		return "0xD1 - BYTE"
	case CIPTypeSINT:
		return "0xC2 - SINT"
	case CIPTypeINT:
		return "0xC3 - INT"
	case CIPTypeDINT:
		return "0xC4 - DINT"
	case CIPTypeLINT:
		return "0xC5 - LINT"
	case CIPTypeUSINT:
		return "0xC6 - USINT"
	case CIPTypeUINT:
		return "0xC7 - UINT"
	case CIPTypeUDINT:
		return "0xC8 - UDINT"
	case CIPTypeLWORD:
		return "0xC9 - LWORD"
	case CIPTypeREAL:
		return "0xCA - REAL"
	case CIPTypeLREAL:
		return "0xCB - LREAL"
	case CIPTypeWORD:
		return "0xD2 - WORD"
	case CIPTypeDWORD:
		return "0xD3 - DWORD"
	case CIPTypeSTRING:
		return "0xDA - String"
	default:
		return "0 - Unknown"
	}
}

func (t CIPType) readValue(r io.Reader) any {
	return readValue(t, r)
}

// readValue reads one unit of cip data type t into the correct go type.
// To do this it reads the needed number of bytes from r.
// It returns the value as an any so the caller will have to do a cast to get it back
func readValue(t CIPType, r io.Reader) any {

	var value any
	var err error
	switch t {
	case CIPTypeUnknown:
		panic("Unknown type.")
	case CIPTypeStruct:
		panic("Struct!")
	case CIPTypeBOOL:
		var trueval bool
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeBYTE:
		var trueval byte
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeSINT:
		var trueval byte
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeINT:
		var trueval int16
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeDINT:
		var trueval int32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeLINT:
		var trueval int64
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeUSINT:
		var trueval uint8
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeUINT:
		var trueval uint16
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeUDINT:
		var trueval uint32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeLWORD:
		var trueval uint64
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeREAL:
		var trueval float32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeLREAL:
		var trueval float64
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeWORD:
		var trueval uint16
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeDWORD:
		var trueval uint32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeSTRING:
		var trueval [86]byte
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	default:
		panic("Default type.")

	}
	if err != nil {
		log.Printf("Problem reading %s as one unit of %T. %v", t, value, err)
	}
	//log.Printf("type %v. value %v", t, value)
	return value
}

// after reading a value v from the controller, you can get a bit from it with
// getBit. bitpos must be between 0 and the length of the CIPType you read.
func getBit(t CIPType, v any, bitpos int) (bool, error) {

	var err error
	switch t {
	case CIPTypeUnknown:
		panic("Unknown type.")
	case CIPTypeStruct:
		panic("Struct!")
	case CIPTypeBOOL:
		if bitpos == 0 {
			x, ok := v.(bool)
			if ok {
				return x, nil
			}
			err = fmt.Errorf("value was a bool, but bit %d was requested. must be 0 for bool", bitpos)
		}
	case CIPTypeBYTE:
		if bitpos >= 0 && bitpos < 8 {
			x, ok := v.(byte)
			if ok {
				mask := byte(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a byte, but bit %d was requested. must be 0-7 for byte", bitpos)
		}
	case CIPTypeSINT:
		if bitpos >= 0 && bitpos < 8 {
			x, ok := v.(byte)
			if ok {
				mask := byte(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a SINT, but bit %d was requested. must be 0-7 for SINT", bitpos)
		}
	case CIPTypeINT:
		if bitpos >= 0 && bitpos < 16 {
			x, ok := v.(int16)
			if ok {
				mask := int16(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was an INT, but bit %d was requested. must be 0-15 for INT", bitpos)
		}
	case CIPTypeDINT:
		if bitpos >= 0 && bitpos < 32 {
			x, ok := v.(int32)
			if ok {
				mask := int32(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a DINT, but bit %d was requested. must be 0-31 for DINT", bitpos)
		}
	case CIPTypeLINT:
		if bitpos >= 0 && bitpos < 64 {
			x, ok := v.(int64)
			if ok {
				mask := int64(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a LINT, but bit %d was requested. must be 0-63 for LINT", bitpos)
		}
	case CIPTypeUSINT:
		if bitpos >= 0 && bitpos < 8 {
			x, ok := v.(byte)
			if ok {
				mask := byte(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a USINT, but bit %d was requested. must be 0-7 for USINT", bitpos)
		}
	case CIPTypeUINT:
		if bitpos >= 0 && bitpos < 16 {
			x, ok := v.(uint16)
			if ok {
				mask := uint16(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was an UINT, but bit %d was requested. must be 0-15 for UINT", bitpos)
		}
	case CIPTypeUDINT:
		if bitpos >= 0 && bitpos < 32 {
			x, ok := v.(uint32)
			if ok {
				mask := uint32(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a UDINT, but bit %d was requested. must be 0-31 for UDINT", bitpos)
		}
	case CIPTypeLWORD:
		if bitpos >= 0 && bitpos < 64 {
			x, ok := v.(uint64)
			if ok {
				mask := uint64(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a LWORD, but bit %d was requested. must be 0-63 for LWORD", bitpos)
		}
	case CIPTypeREAL:
		err = fmt.Errorf("value was a REAL, not finding bit of real")
	case CIPTypeLREAL:
		err = fmt.Errorf("value was a LEAL, not finding bit of real")
	case CIPTypeWORD:
		if bitpos >= 0 && bitpos < 16 {
			x, ok := v.(uint16)
			if ok {
				mask := uint16(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was an WORD, but bit %d was requested. must be 0-15 for WORD", bitpos)
		}
	case CIPTypeDWORD:
		if bitpos >= 0 && bitpos < 32 {
			x, ok := v.(uint32)
			if ok {
				mask := uint32(1 << bitpos)
				masked := x & mask
				return masked != 0, nil
			}
			err = fmt.Errorf("value was a DWORD, but bit %d was requested. must be 0-31 for DWORD", bitpos)
		}
	case CIPTypeSTRING:
		err = fmt.Errorf("value was a STRING, not finding bit of string")
	default:
		panic("Default type.")

	}
	if err != nil {
		log.Printf("%v", err)
	}
	return false, err
}
