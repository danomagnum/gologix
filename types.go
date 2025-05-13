package gologix

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type CIPType byte

// Go native types that correspond to logix types
// I'm not sure whether having interface here makes sense.
// On the one hand, we need to support composite types, but on the other this lets it accept anything
// which doesn't seem right.
type GoLogixTypes interface {
	bool | byte | int8 | uint16 | int16 | uint32 | int32 | uint64 | int64 | float32 | float64 | string
}

// return the CIPType that corresponds to go type of variable T
// also return the element count
func GoVarToCIPType(T any) (CIPType, int) {
	switch x := T.(type) {
	case bool:
		return CIPTypeBOOL, 1
	case byte:
		return CIPTypeBYTE, 1
	case int8:
		return CIPTypeSINT, 1
	case uint16:
		return CIPTypeUINT, 1
	case int16:
		return CIPTypeINT, 1
	case uint32:
		return CIPTypeUDINT, 1
	case int32:
		return CIPTypeDINT, 1
	case uint64:
		return CIPTypeLWORD, 1
	case int64:
		return CIPTypeLINT, 1
	case float32:
		return CIPTypeREAL, 1
	case float64:
		return CIPTypeLREAL, 1
	case string:
		return CIPTypeSTRING, 1
	case []bool:
		return CIPTypeBOOL, len(x)
	case []byte:
		return CIPTypeBYTE, len(x)
	case []int8:
		return CIPTypeSINT, len(x)
	case []uint16:
		return CIPTypeUINT, len(x)
	case []int16:
		return CIPTypeINT, len(x)
	case []uint32:
		return CIPTypeUDINT, len(x)
	case []int32:
		return CIPTypeDINT, len(x)
	case []uint64:
		return CIPTypeLWORD, len(x)
	case []int64:
		return CIPTypeLINT, len(x)
	case []float32:
		return CIPTypeREAL, len(x)
	case []float64:
		return CIPTypeLREAL, len(x)
	case []string:
		return CIPTypeSTRING, len(x)
	case interface{}:
		return CIPTypeStruct, 1
	}
	return CIPTypeUnknown, 1
}

const (
	CIPTypeUnknown         CIPType = 0x00
	CIPTypeStruct          CIPType = 0xA0 // also used for strings.  Not sure what's up with CIPTypeSTRING
	CIPTypeUTIME           CIPType = 0xC0
	CIPTypeBOOL            CIPType = 0xC1
	CIPTypeSINT            CIPType = 0xC2
	CIPTypeINT             CIPType = 0xC3
	CIPTypeDINT            CIPType = 0xC4
	CIPTypeLINT            CIPType = 0xC5
	CIPTypeUSINT           CIPType = 0xC6
	CIPTypeUINT            CIPType = 0xC7
	CIPTypeUDINT           CIPType = 0xC8
	CIPTypeULINT           CIPType = 0xC9
	CIPTypeREAL            CIPType = 0xCA
	CIPTypeLREAL           CIPType = 0xCB
	CIPTypeSTIME           CIPType = 0xCC
	CIPTypeDATE            CIPType = 0xCD
	CIPTypeTIMEOFDAY       CIPType = 0xCE
	CIPTypeDATETIME        CIPType = 0xCF
	CIPTypeSTRING_UNKNOWN  CIPType = 0xD0
	CIPTypeBYTE            CIPType = 0xD1 // 8 bits packed into one byte
	CIPTypeWORD            CIPType = 0xD2
	CIPTypeDWORD           CIPType = 0xD3
	CIPTypeLWORD           CIPType = 0xD4
	CIPTypeSTRING_UNKNOWN2 CIPType = 0xD5
	CIPTypeFTIME           CIPType = 0xD6
	CIPTypeLTIME           CIPType = 0xD7
	CIPTypeITIME           CIPType = 0xD8
	CIPTypeSTRING_UNKNOWN3 CIPType = 0xD9
	CIPTypeSTRING_SHORT    CIPType = 0xDA
	CIPTypeTIMEOFDAY2      CIPType = 0xDB
	CIPTypeEPATH           CIPType = 0xDC
	CIPTypeENGUNIT         CIPType = 0xDD

	//  Strings actually come accross as 0xA0 = CIPTypeStruct.
	//In this library we're using this as kind of a flag to keep track of whether
	// a structure is a normal logix string or not.
	CIPTypeSTRING CIPType = 0xFF
)

// return the size in bytes of the data structure
func (c CIPType) Size() int {
	switch c {
	case CIPTypeUnknown:
		return 0
	case CIPTypeStruct:
		return 88
	case CIPTypeUTIME:
		return 8
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
	case CIPTypeULINT:
		return 8
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
	case CIPTypeDATE:
		return 2
	case CIPTypeTIMEOFDAY:
		return 6
	case CIPTypeDATETIME:
		return 0 //?
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
		return "0xFF - (gologix specific) String"
	default:
		return fmt.Sprintf("0x%2x - Unknown", byte(c))
	}
}

func (t CIPType) IsAtomic() bool {
	v := byte(t)
	return v <= 254
}
func (t CIPType) readValue(r io.Reader) (any, error) {
	return readValue(t, r)
}

// readValue reads one unit of cip data type t into the correct go type.
// To do this it reads the needed number of bytes from r.
// It returns the value as an any so the caller will have to do a cast to get it back
func readValue(t CIPType, r io.Reader) (any, error) {

	var value any
	var err error
	switch t {
	case CIPTypeUnknown:
		return nil, fmt.Errorf("unknown type")
	case CIPTypeStruct:
		return nil, fmt.Errorf("don't know what to do with a struct")
	case CIPTypeBOOL:
		var trueval bool
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeBYTE:
		var trueval byte
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeSINT:
		var trueval int8
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
		return nil, fmt.Errorf("default (unknown) type %d", t)
		//panic(fmt.Sprintf("Default type %d", t))

	}
	if err != nil {
		return nil, fmt.Errorf("problem reading %s as one unit of %T. %w", t, value, err)
	}
	return value, nil
}

// after reading a value v from the controller, you can get a bit from it with
// getBit. bitpos must be between 0 and the length of the CIPType you read.
func getBit(t CIPType, v any, bitpos int) (bool, error) {

	var err error
	switch t {
	case CIPTypeUnknown:
		return false, errors.New("unknown type")
		//panic("Unknown type.")
	case CIPTypeStruct:
		return false, errors.New("got a struct - can't get a bit")
		//panic("Struct!")
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
		return false, errors.New("got an unknown type. don't know how to get bit")
		//panic("Default type.")

	}
	if err != nil {
		return false, err
	}
	return false, err
}
