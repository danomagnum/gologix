package main

import (
	"encoding/binary"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type TagPartDescriptor struct {
	FullPath    string
	BasePath    string
	Array_Order []int
	BitNumber   int
	BitAccess   bool
}

var bit_access_regex, _ = regexp.Compile(`\.\d+$`)
var array_access_regex, _ = regexp.Compile(`\[([\d]|[,]|[\s])*\]$`)

func (tag *TagPartDescriptor) Parse(tagpath string) error {
	var err error
	tag.FullPath = tagpath
	tag.BasePath = tagpath

	// Check if tag is accessing a bit in the data
	bitpos := bit_access_regex.FindStringIndex(tagpath)
	if bitpos == nil {
		tag.BitAccess = false
	} else {
		tag.BitAccess = true
		bit_access_text := tagpath[bitpos[0]+1 : bitpos[1]]
		tag.BasePath = strings.ReplaceAll(tag.BasePath, bit_access_text, "")
		tag.BitNumber, err = strconv.Atoi(bit_access_text)
		if err != nil {
			return fmt.Errorf("could't parse %v to a bit portion of tag. %v", bit_access_text, err)
		}
	}

	// check if tag is accessing an array
	arrpos := array_access_regex.FindStringIndex(tagpath)
	if arrpos == nil {
		tag.Array_Order = nil
	} else {
		arr_access_text := tagpath[arrpos[0]+1 : arrpos[1]-1]
		tag.BasePath = strings.ReplaceAll(tag.BasePath, arr_access_text, "")
		if strings.Contains(arr_access_text, ",") {
			parts := strings.Split(arr_access_text, ",")
			tag.Array_Order = make([]int, len(parts))
			for i, part := range parts {
				tag.Array_Order[i], err = strconv.Atoi(part)
				if err != nil {
					return fmt.Errorf("could't parse %v to an array position. %v", arr_access_text, err)
				}
			}

		} else {
			tag.Array_Order = make([]int, 1)
			tag.Array_Order[0], err = strconv.Atoi(arr_access_text)
			if err != nil {
				return fmt.Errorf("could't parse %v to an array position. %v", arr_access_text, err)
			}
		}
	}

	return nil

}

// parse the tag name into its base tag (remove array index or bit) and get the array index if it exists
func parse_tag_name(tagpath string) (tag TagPartDescriptor) {
	tag.Parse(tagpath)
	return

}

type IOI struct {
	Path   string
	Type   CIPType
	Buffer []byte
}

func (ioi *IOI) Write(p []byte) (n int, err error) {
	ioi.Buffer = append(ioi.Buffer, p...)
	return len(p), nil
}

// this is the default buffer size for tag IOI generation.
const DEFAULT_BUFFER_SIZE = 256

// The IOI is the tag name structure that CIP requires.  It's parsed out into tag length, tag name pairs with additional
// data on the backside to indicate what index is requested if needed.
func BuildIOI(tagpath string, datatype CIPType) (ioi *IOI) {
	tag_array := strings.Split(tagpath, ".")

	ioi = new(IOI)
	ioi.Path = tagpath
	ioi.Type = datatype
	// we'll build this byte structure up as we go.
	ioi.Buffer = make([]byte, 0, DEFAULT_BUFFER_SIZE)

	for _, tag_part := range tag_array {
		if strings.HasSuffix(tag_part, "]") {
			// part of an array
			start_index := strings.Index(tag_part, "[")
			ioi_part := buildIOI_Part(tag_part[0:start_index])
			ioi.Write(ioi_part)

			t := parse_tag_name(tag_part)

			for _, order_size := range t.Array_Order {
				if order_size < 256 {
					// byte, byte
					index_part := []byte{0x28, byte(order_size)}
					binary.Write(ioi, binary.LittleEndian, index_part)
				} else if order_size < 65536 {
					// uint16, uint16
					index_part := []uint16{0x29, uint16(order_size)}
					binary.Write(ioi, binary.LittleEndian, index_part)
				} else {
					// uint16, uint32
					index_part0 := []uint16{0x2A}
					binary.Write(ioi, binary.LittleEndian, index_part0)
					index_part1 := []uint32{uint32(order_size)}
					binary.Write(ioi, binary.LittleEndian, index_part1)
				}
			}

		} else {
			// not part of an array
			bit_access, err := strconv.Atoi(tag_part)
			if err == nil && bit_access <= 31 {
				// This is a bit access.
				// we won't do anything for now and will just parse the
				// bit out of the word when taht time comes.
				continue
			}
			ioi_part := buildIOI_Part(tag_part)
			ioi.Write(ioi_part)

		}

	}

	return ioi
}

func buildIOI_Part(tagpath string) []byte {
	t := parse_tag_name(tagpath)
	tag_size := len(t.BasePath)

	tag_name_header := [2]byte{0x91, byte(tag_size)}
	tag_name_msg := append(tag_name_header[:], []byte(t.BasePath)...)
	// has to be an even number of bytes.
	if tag_size%2 == 1 {
		tag_name_msg = append(tag_name_msg, []byte{0x00}...)
	}
	return tag_name_msg
}
