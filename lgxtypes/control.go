package lgxtypes

import (
	"encoding/binary"
	"io"
)

type CONTROL struct {
	LEN int32
	POS int32
	EN  bool // bit 31
	EU  bool // bit 30
	DN  bool // bit 29
	EM  bool // bit 28
	ER  bool // bit 27
	UL  bool // bit 26
	IN  bool // bit 25
	FD  bool // bit 24
}

func (t CONTROL) Pack(w io.Writer) (int, error) {

	var CtrlWord uint32
	if t.EN {
		CtrlWord |= 1 << 31
	}
	if t.EU {
		CtrlWord |= 1 << 30
	}
	if t.DN {
		CtrlWord |= 1 << 29
	}
	if t.EM {
		CtrlWord |= 1 << 28
	}
	if t.ER {
		CtrlWord |= 1 << 27
	}
	if t.UL {
		CtrlWord |= 1 << 26
	}
	if t.IN {
		CtrlWord |= 1 << 25
	}
	if t.FD {
		CtrlWord |= 1 << 24
	}

	err := binary.Write(w, binary.LittleEndian, CtrlWord)
	if err != nil {
		return 0, err
	}

	err = binary.Write(w, binary.LittleEndian, t.LEN)
	if err != nil {
		return 4, err
	}
	err = binary.Write(w, binary.LittleEndian, t.POS)
	if err != nil {
		return 8, err
	}

	return 12, nil
}

func (t *CONTROL) Unpack(r io.Reader) (int, error) {
	var CtrlWord uint32
	err := binary.Read(r, binary.LittleEndian, &CtrlWord)
	if err != nil {
		return 0, err
	}

	t.EN = CtrlWord&(1<<31) != 0
	t.EU = CtrlWord&(1<<30) != 0
	t.DN = CtrlWord&(1<<29) != 0
	t.EM = CtrlWord&(1<<28) != 0
	t.ER = CtrlWord&(1<<27) != 0
	t.UL = CtrlWord&(1<<26) != 0
	t.IN = CtrlWord&(1<<25) != 0
	t.FD = CtrlWord&(1<<24) != 0

	err = binary.Read(r, binary.LittleEndian, &(t.LEN))
	if err != nil {
		return 4, err
	}

	err = binary.Read(r, binary.LittleEndian, &(t.POS))
	if err != nil {
		return 8, err
	}

	return 12, nil
}
