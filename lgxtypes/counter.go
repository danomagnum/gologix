package lgxtypes

import (
	"encoding/binary"
	"io"
)

type COUNTER struct {
	PRE int32
	ACC int32
	CU  bool // bit 31
	CD  bool // bit 30
	DN  bool // bit 29
	OV  bool // bit 28 - set if we wrap over 2,147,483,648 to -2,147,483,648
	UN  bool // bit 27 - set if we wrap over -2,147,483,648 to 2,147,483,648
}

func (t COUNTER) Pack(w io.Writer) (int, error) {

	var CtrlWord uint32
	if t.CU {
		CtrlWord |= 1 << 31
	}
	if t.CD {
		CtrlWord |= 1 << 30
	}
	if t.DN {
		CtrlWord |= 1 << 29
	}
	if t.OV {
		CtrlWord |= 1 << 28
	}
	if t.UN {
		CtrlWord |= 1 << 27
	}

	err := binary.Write(w, binary.LittleEndian, CtrlWord)
	if err != nil {
		return 0, err
	}

	err = binary.Write(w, binary.LittleEndian, t.PRE)
	if err != nil {
		return 4, err
	}
	err = binary.Write(w, binary.LittleEndian, t.ACC)
	if err != nil {
		return 8, err
	}

	return 12, nil
}

func (t *COUNTER) Unpack(r io.Reader) (int, error) {
	var CtrlWord uint32
	err := binary.Read(r, binary.LittleEndian, &CtrlWord)
	if err != nil {
		return 0, err
	}

	t.CU = CtrlWord&(1<<31) != 0
	t.CD = CtrlWord&(1<<30) != 0
	t.DN = CtrlWord&(1<<29) != 0
	t.OV = CtrlWord&(1<<28) != 0
	t.UN = CtrlWord&(1<<27) != 0

	err = binary.Read(r, binary.LittleEndian, &(t.PRE))
	if err != nil {
		return 4, err
	}

	err = binary.Read(r, binary.LittleEndian, &(t.ACC))
	if err != nil {
		return 8, err
	}

	return 12, nil
}

func (COUNTER) TypeAbbr() (string, uint16) {
	return "COUNTER,DINT,DINT,DINT", 0x0F82
}
