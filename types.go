package main

type CIPType byte

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
