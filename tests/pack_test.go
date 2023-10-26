package gologix_tests

import (
	"bytes"
	"testing"

	"github.com/danomagnum/gologix"
)

func TestPack(t *testing.T) {
	type S2 struct {
		Flag1    bool
		U32      uint32
		PostFlag bool
	}
	type S struct {
		Flag0  bool
		Flag1  bool `pack:"nopack"`
		Flag2  bool
		Flag3  bool
		Flag4  bool
		Flag5  bool
		Flag6  bool
		Flag7  bool
		Flag8  bool
		Flag9  bool
		Flag10 bool
		Flag11 bool
		U32    uint32
		Sub1   S2
		Flags  [16]bool //`pack:"nopack"`
		Sub2   S2
	}
	s := S{}
	s.Flag1 = true
	s.Flag3 = true
	s.Flag9 = true
	s.Flag10 = true
	s.Flag11 = true
	//s.Flag7 = true
	//s.Flag8 = true
	s.U32 = 0xFFFF_FFFF
	s.Sub1.Flag1 = true
	s.Sub1.U32 = 0xEEEE_EEEE
	s.Flags[0] = true
	s.Flags[15] = true
	s.Sub2.Flag1 = true
	s.Sub2.U32 = 0xDDDD_DDDD
	b := bytes.Buffer{}
	//binary.Write(&b, binary.LittleEndian, s)

	gologix.Pack(&b, gologix.CIPPack{}, s)

	have := b.Bytes()

	want := []byte{
		0,      // flag0
		1,      // flag1 gets packed on its own because of nopack
		130, 3, // flag2+ all get packed into combined words
		0xFF, 0xFF, 0xFF, 0xFF, // U32
		1, 0, 0, 0, // single flag1 in S2 gets buffered to 4 bytes
		0xEE, 0xEE, 0xEE, 0xEE, // U32 in Sub1
		0,         // post flag in sub1 //TODO: does this need to pad?
		1, 128, 0, // 16 bool array padded out to 4 byte boundry
		1, 0, 0, 0, // single flag1 in S2 buffered to 4 bytes
		0xDD, 0xDD, 0xDD, 0xDD, // U32 in Sub2
		0, // last postflag in sub2
	}

	if !check_bytes(have, want) {
		t.Errorf("ResultMismatch.\n Have %v\n Want %v\n", have, want)
	}

	have2 := S{}
	_, err := gologix.Unpack(bytes.NewBuffer(b.Bytes()), gologix.CIPPack{}, &have2)
	if err != nil {
		t.Errorf("problem unpacking bytes. %v", err)
	}
	if have2 != s {
		t.Errorf("ResultMismatch.\n Have \n%v\n Want \n%v\n", have2, s)

	}
}

func TestPack2(t *testing.T) {
	type S2 struct {
		Flag1    bool
		U32      uint32
		PostFlag bool
		LongInt  uint64
	}
	type S struct {
		Flag0  bool
		Flag1  bool `pack:"nopack"`
		Flag2  bool
		Flag3  bool
		Flag4  bool
		Flag5  bool
		Flag6  bool
		Flag7  bool
		Flag8  bool
		Flag9  bool
		Flag10 bool
		Flag11 bool
		U32    uint32
		Sub1   S2
		Flags  [16]bool //`pack:"nopack"`
		Sub2   S2
	}
	s := S{}
	s.Flag1 = true
	s.Flag3 = true
	s.Flag9 = true
	s.Flag10 = true
	s.Flag11 = true
	//s.Flag7 = true
	//s.Flag8 = true
	s.U32 = 0xFFFF_FFFF
	s.Sub1.Flag1 = true
	s.Sub1.U32 = 0xEEEE_EEEE
	s.Sub1.LongInt = 0xCCCC_CCCC_CCCC_CCCC
	s.Flags[0] = true
	s.Flags[15] = true
	s.Sub2.Flag1 = true
	s.Sub2.U32 = 0xDDDD_DDDD
	s.Sub2.LongInt = 0xBBBB_BBBB_BBBB_BBBB
	b := bytes.Buffer{}
	//binary.Write(&b, binary.LittleEndian, s)

	gologix.Pack(&b, gologix.CIPPack{}, s)

	have := b.Bytes()

	want := []byte{
		0,      // flag0
		1,      // flag1 gets packed on its own because of nopack
		130, 3, // flag2+ all get packed into combined words
		0xFF, 0xFF, 0xFF, 0xFF, // U32
		1, 0, 0, 0, // start of Sub1 on byte 8. single flag1 in S2 gets buffered to 4 bytes
		0xEE, 0xEE, 0xEE, 0xEE, // U32 in Sub1
		0,                   // post flag in sub1 //TODO: does this need to pad?
		0, 0, 0, 0, 0, 0, 0, // padding for lint
		0xCC, 0xCC, 0xCC, 0xCC, 0xCC, 0xCC, 0xCC, 0xCC, // lint
		1, 128, 0, 0, 0, 0, 0, 0, // 16 bool array padded out to 8 byte boundry for sub2
		1, 0, 0, 0, // single flag1 in S2 buffered to 4 bytes
		0xDD, 0xDD, 0xDD, 0xDD, // U32 in Sub2
		0,                   // last postflag in sub2
		0, 0, 0, 0, 0, 0, 0, // padding for lint in sub2
		0xBB, 0xBB, 0xBB, 0xBB, 0xBB, 0xBB, 0xBB, 0xBB, // lint in sub2
	}

	if !check_bytes(have, want) {
		t.Errorf("ResultMismatch.\n Have %v\n Want %v\n", have, want)
	}

	have2 := S{}
	_, err := gologix.Unpack(bytes.NewBuffer(b.Bytes()), gologix.CIPPack{}, &have2)
	if err != nil {
		t.Errorf("problem unpacking bytes. %v", err)
	}
	if have2 != s {
		t.Errorf("ResultMismatch.\n Have \n%v\n Want \n%v\n", have2, s)

	}
}
