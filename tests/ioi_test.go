package gologix_tests

import (
	"fmt"
	"testing"

	"github.com/danomagnum/gologix"
	"github.com/danomagnum/gologix/ciptype"
)

// these tests came from the tag names in 1756-PM020H-EN-P
// it only tests the request path portion of each tag addressing example
// it also only tests symbolic paths.
func TestIOI(t *testing.T) {
	var tests = []struct {
		path string
		t    ciptype.CIPType
		want []byte
	}{
		{
			"profile[0,1,257]",
			ciptype.DINT,
			[]byte{
				0x91, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x00, // symbolic segment for "profile"
				0x28, 0x00, // member segment for 0
				0x28, 0x01, // member segment for 1
				0x29, 0x00, 0x01, 0x01, // member segment for 257
			},
		},
		{
			"profile[1,2,258]",
			ciptype.DINT,
			[]byte{
				0x91, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x00, // symbolic segment for "profile"
				0x28, 0x01, // member segment for 1
				0x28, 0x02, // member segment for 2
				0x29, 0x00, 0x02, 0x01, // member segment for 258
			},
		},
		{
			"profile[300,2,258]",
			ciptype.DINT,
			[]byte{
				0x91, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x00, // symbolic segment for "profile"
				0x29, 0x00, 0x2c, 0x01, // member segment for 300
				0x28, 0x02, // member segment for 2
				0x29, 0x00, 0x02, 0x01, // member segment for 258
			},
		},
		{
			"dwell3.acc",
			ciptype.DINT,
			[]byte{
				0x91, 0x06, 0x64, 0x77, 0x65, 0x6C, 0x6C, 0x33, // symbolic segment for "dwell3"
				0x91, 0x03, 0x61, 0x63, 0x63, 0x00, // member segment for ACC
			},
		},
		{
			"struct3.today.rate",
			ciptype.Struct,
			[]byte{
				0x91, 0x07, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x33, 0x00, // symbolic segment for "struct3"
				0x91, 0x05, 0x74, 0x6F, 0x64, 0x61, 0x79, 0x00, // symbolic segment for today
				0x91, 0x04, 0x72, 0x61, 0x74, 0x65, // symbolic segment for rate
			},
		},
		{
			"my2dstruct4[1].today.hourlycount[3]",
			ciptype.INT,
			[]byte{
				0x91, 0x0B, 0x6d, 0x79, 0x32, 0x64, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x34, 0x00, // symbolic segment for my2dstruct4
				0x28, 0x01, // index 1
				0x91, 0x05, 0x74, 0x6F, 0x64, 0x61, 0x79, 0x00, //today
				0x91, 0x0B, 0x68, 0x6F, 0x75, 0x72, 0x6C, 0x79, 0x63, 0x6F, 0x75, 0x6E, 0x74, 0x00, // hourlycount
				0x28, 0x03, // index 3
			},
		},
		{
			"My2DstRucT4[1].ToDaY.hoURLycOuNt[3]",
			ciptype.INT,
			[]byte{
				0x91, 0x0B, 0x6d, 0x79, 0x32, 0x64, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x34, 0x00, // symbolic segment for my2dstruct4
				0x28, 0x01, // index 1
				0x91, 0x05, 0x74, 0x6F, 0x64, 0x61, 0x79, 0x00, //today
				0x91, 0x0B, 0x68, 0x6F, 0x75, 0x72, 0x6C, 0x79, 0x63, 0x6F, 0x75, 0x6E, 0x74, 0x00, // hourlycount
				0x28, 0x03, // index 3
			},
		},
	}
	client := gologix.Client{}

	for _, tt := range tests {

		testname := fmt.Sprintf("tag: %s", tt.path)
		t.Run(testname, func(t *testing.T) {
			res, err := client.NewIOI(tt.path, tt.t)
			if err != nil {
				t.Errorf("IOI Generation error. %v", err)
			}
			if !check_bytes(res.Buffer, tt.want) {
				t.Errorf("Wrong Value for result.  \nWanted %v. \nGot    %v", to_hex(tt.want), to_hex(res.Buffer))
			}
		})
	}

}

func to_hex(b []byte) []string {
	out := make([]string, len(b))

	for i, v := range b {
		out[i] = fmt.Sprintf("% X", v)
	}
	return out

}

func TestIOIToBytesAndBackAgain(t *testing.T) {
	tests := []struct {
		Tag  string
		Type ciptype.CIPType
	}{
		{"test", ciptype.DINT},
		{"test[2]", ciptype.DINT},
		{"test[2,3]", ciptype.DINT},
		{"test[3000,3]", ciptype.DINT},
		{"test.tester", ciptype.DINT},
		{"test[2,3].tester", ciptype.DINT},
	}
	client := gologix.Client{}

	for _, tt := range tests {

		testname := fmt.Sprintf("tag: %s", tt.Tag)
		t.Run(testname, func(t *testing.T) {
			res, err := client.NewIOI(tt.Tag, tt.Type)
			if err != nil {
				t.Errorf("IOI Generation error. %v", err)
			}
			item := gologix.NewItem(0, res)
			path, err := gologix.GetTagFromPath(&item)
			if err != nil {
				t.Errorf("problem parsing path from byte item")
			}
			if path != tt.Tag {
				t.Errorf("Wrong Value for result.  \nWanted %v. \nGot    %v", tt.Tag, path)
			}
		})
	}

}
