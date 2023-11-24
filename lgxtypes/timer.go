package lgxtypes

import (
	"encoding/binary"
	"io"
)

type TIMER struct {
	PRE int32
	ACC int32
	EN  bool // bit 31
	TT  bool // bit 30
	DN  bool // bit 29

	// These bits were added in anticipation of using timers with SFCs (Sequential Function Charts).
	// However, at this time SFCs do not use timer structures, so these bits are not used and are currently undefined.
	FS bool // bit 28 Unused
	LS bool // bit 27 Unused
	OV bool // bit 26 Unused
	ER bool // bit 25 Unused
}

func (t TIMER) Pack(w io.Writer) (int, error) {

	var CtrlWord uint32
	if t.EN {
		CtrlWord |= 1 << 31
	}
	if t.TT {
		CtrlWord |= 1 << 30
	}
	if t.DN {
		CtrlWord |= 1 << 29
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

func (t *TIMER) Unpack(r io.Reader) (int, error) {
	var CtrlWord uint32
	err := binary.Read(r, binary.LittleEndian, &CtrlWord)
	if err != nil {
		return 0, err
	}

	t.EN = CtrlWord&(1<<31) != 0
	t.TT = CtrlWord&(1<<30) != 0
	t.DN = CtrlWord&(1<<29) != 0

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
