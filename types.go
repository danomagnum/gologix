package main

type CIPType byte

// Go native types that correspond to logix types
type GoLogixTypes interface {
	bool | byte | uint16 | int16 | uint32 | int32 | uint64 | int64 | float32 | float64
}

func GoTypeToLogixType(T any) CIPType {
	switch T.(type) {
	case byte:
		return CIPTypeBOOL
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
	}
	return CIPTypeUnknown
}

const (
	CIPTypeUnknown CIPType = 0x00
	CIPTypeStruct  CIPType = 0xA0
	CIPTypeBOOL    CIPType = 0xC1
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
	CIPTypeDWORD   CIPType = 0xD3
	CIPTypeSTRING  CIPType = 0xDA
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

// return the size in bytes of the data structure
func (c CIPType) String() string {
	switch c {
	case CIPTypeUnknown:
		return "0 - Unknown"
	case CIPTypeStruct:
		return "88 - Struct"
	case CIPTypeBOOL:
		return "1 - BOOL"
	case CIPTypeSINT:
		return "1 - SINT"
	case CIPTypeINT:
		return "2 - INT"
	case CIPTypeDINT:
		return "4 - DINT"
	case CIPTypeLINT:
		return "8 - LINT"
	case CIPTypeUSINT:
		return "1 - USINT"
	case CIPTypeUINT:
		return "2 - UINT"
	case CIPTypeUDINT:
		return "4 - UDINT"
	case CIPTypeLWORD:
		return "8 - LWORD"
	case CIPTypeREAL:
		return "4 - REAL"
	case CIPTypeLREAL:
		return "8 - LREAL"
	case CIPTypeDWORD:
		return "4 - DWORD"
	case CIPTypeSTRING:
		return "1 - String"
	default:
		return "0 - Unknown"
	}
}
