package lgxtypes

type STRING struct {
	Len  int32
	Data [82]int8
}

func (STRING) TypeAbbr() (string, uint16) {
	return "STRING,DINT,SINT[82]", 0x0FCE
}
