package gologix

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type tagPartDescriptor struct {
	FullPath    string
	BasePath    string
	Array_Order []int
	BitNumber   int
	BitAccess   bool
}

var bit_access_regex, _ = regexp.Compile(`\.\d+$`)
var array_access_regex, _ = regexp.Compile(`\[([\d]|[,]|[\s])*\]$`)

func (tag *tagPartDescriptor) Parse(tagpath string) error {
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
			return fmt.Errorf("could't parse %v to a bit portion of tag. %w", bit_access_text, err)
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
					return fmt.Errorf("could't parse %v to an array position. %w", arr_access_text, err)
				}
			}

		} else {
			tag.Array_Order = make([]int, 1)
			tag.Array_Order[0], err = strconv.Atoi(arr_access_text)
			if err != nil {
				return fmt.Errorf("could't parse %v to an array position. %w", arr_access_text, err)
			}
		}
	}

	return nil

}

// parse the tag name into its base tag (remove array index or bit) and get the array index if it exists
func parse_tag_name(tagpath string) (tag tagPartDescriptor) {
	err := tag.Parse(tagpath)
	if err != nil {
		log.Printf("problem parsing path. %v", err)
		return
	}
	return

}

// Internal Object Identifier. Used to specify a tag name in the controller
// the Buffer has the CIP route for a tag path.
type tagIOI struct {
	Path        string
	Type        CIPType
	BitAccess   bool
	BitPosition int
	Buffer      []byte
}

func (ioi *tagIOI) Write(p []byte) (n int, err error) {
	ioi.Buffer = append(ioi.Buffer, p...)
	return len(p), nil
}

func (ioi *tagIOI) Bytes() []byte {
	return ioi.Buffer
}
func (ioi *tagIOI) Len() int {
	return len(ioi.Buffer)
}

// this is the default buffer size for tag IOI generation.
const defaultIOIBufferSize = 256

// The IOI is the tag name structure that CIP requires.  It's parsed out into tag length, tag name pairs with additional
// data on the backside to indicate what index is requested if needed.
func (client *Client) NewIOI(tagpath string, datatype CIPType) (ioi *tagIOI, err error) {
	if client.ioi_cache == nil {
		client.ioi_cache = make(map[string]*tagIOI)
	}
	ioi = new(tagIOI)
	// CIP doesn't care about case.  But we'll make it lowercase to match
	// the encodings shown in 1756-PM020H-EN-P
	tagpath = strings.ToLower(tagpath)
	tag_info, ok := client.KnownTags[tagpath]
	if ok {
		if tag_info.Info.Type != datatype && datatype != CIPTypeUnknown {
			err = fmt.Errorf("data type mismatch for IOI. %v was specified, but I have reason to believe that it's really %v", datatype, tag_info.Info.Type)
			return
		}
		if tag_info.Info.TypeInfo != 0 {
			ioi.Buffer = tag_info.Bytes()
			return
		}
	}
	extant, exists := client.ioi_cache[tagpath]
	if exists {
		ioi = extant
		return
	}
	tag_array := strings.Split(tagpath, ".")

	ioi.Path = tagpath
	ioi.Type = datatype
	// we'll build this byte structure up as we go.
	ioi.Buffer = make([]byte, 0, defaultIOIBufferSize)

	for _, tag_part := range tag_array {
		if strings.HasSuffix(tag_part, "]") {
			// part of an array
			start_index := strings.Index(tag_part, "[")
			ioi_part := marshalIOIPart(tag_part[0:start_index])
			_, err = ioi.Write(ioi_part)
			if err != nil {
				return ioi, fmt.Errorf("problem writing ioi part %w", err)

			}

			t := parse_tag_name(tag_part)

			for _, order_size := range t.Array_Order {
				if order_size < 256 {
					// byte, byte
					index_part := []byte{byte(cipElement_8bit), byte(order_size)}
					err = binary.Write(ioi, binary.LittleEndian, index_part)
					if err != nil {
						return nil, fmt.Errorf("problem reading index part. %w", err)
					}
				} else if order_size < 65536 {
					// uint16, uint16
					index_part := []uint16{uint16(cipElement_16bit), uint16(order_size)}
					err = binary.Write(ioi, binary.LittleEndian, index_part)
					if err != nil {
						return nil, fmt.Errorf("problem reading index part. %w", err)
					}
				} else {
					// uint16, uint32
					index_part0 := []uint16{uint16(cipElement_32bit)}
					err = binary.Write(ioi, binary.LittleEndian, index_part0)
					if err != nil {
						return nil, err
					}
					index_part1 := []uint32{uint32(order_size)}
					err = binary.Write(ioi, binary.LittleEndian, index_part1)
					if err != nil {
						return nil, err
					}
				}
			}

		} else {
			// not part of an array
			bit_access, err := strconv.Atoi(tag_part)
			if err == nil && bit_access <= 31 {
				// This is a bit access.
				// we won't do anything for now and will just parse the
				// bit out of the word when that time comes.
				ioi.BitAccess = true
				ioi.BitPosition = bit_access
				continue
			}
			ioi_part := marshalIOIPart(tag_part)
			_, err = ioi.Write(ioi_part)
			if err != nil {
				return nil, err
			}

		}

	}

	client.ioi_cache[tagpath] = ioi
	return
}

func marshalIOIPart(tagpath string) []byte {
	t := parse_tag_name(tagpath)
	tag_size := len(t.BasePath)
	need_extend := false
	if tag_size%2 == 1 {
		need_extend = true
		//tag_size += 1
	}

	tag_name_header := [2]byte{byte(SegmentTypeExtendedSymbolic), byte(tag_size)}
	tag_name_msg := append(tag_name_header[:], []byte(t.BasePath)...)
	// has to be an even number of bytes.
	if need_extend {
		tag_name_msg = append(tag_name_msg, []byte{0x00}...)
	}
	return tag_name_msg
}

// these next functions are for reversing the bytes back to a tag string
func getAsciiTagPart(item *CIPItem) (string, error) {
	var tag_len byte
	err := item.DeSerialize(&tag_len)
	if err != nil {
		return "", fmt.Errorf("problem getting tag len. %w", err)
	}
	b := make([]byte, tag_len)
	err = item.DeSerialize(&b)
	if err != nil {
		return "", fmt.Errorf("problem reading tag path. %w", err)
	}
	if tag_len%2 == 1 {
		var pad byte
		err = item.DeSerialize(&pad)
		if err != nil {
			return "", fmt.Errorf("problem reading pad byte. %w", err)
		}
	}

	tag_str := string(b)
	return tag_str, nil
}
func getTagFromPath(item *CIPItem) (string, error) {

	tag_str := ""

morepath:
	for {
		// we haven't read all the tag path info.
		var tag_path_type byte
		err := item.DeSerialize(&tag_path_type)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return tag_str, nil

			}
			return "", fmt.Errorf("couldn't get path part type: %w", err)
		}
		switch tag_path_type {
		case 0x28:
			// one byte index
			var array_index byte
			err = item.DeSerialize(&array_index)
			if err != nil {
				return "", fmt.Errorf("couldn't get array index: %w", err)
			}
			if tag_str[len(tag_str)-1] == ']' {
				tag_str = fmt.Sprintf("%s,%d]", tag_str[:len(tag_str)-1], array_index)
			} else {
				tag_str = fmt.Sprintf("%s[%d]", tag_str, array_index)
			}
		case 0x29:
			// two byte index
			var pad byte
			err = item.DeSerialize(&pad)
			if err != nil {
				return "", fmt.Errorf("couldn't get padding: %w", err)
			}
			var array_index uint16
			err = item.DeSerialize(&array_index)
			if err != nil {
				return "", fmt.Errorf("couldn't get array index: %w", err)
			}
			if tag_str[len(tag_str)-1] == ']' {
				tag_str = fmt.Sprintf("%s,%d]", tag_str[:len(tag_str)-1], array_index)
			} else {
				tag_str = fmt.Sprintf("%s[%d]", tag_str, array_index)
			}
		case 0x91:
			// ascii portion of tag path
			s, err := getAsciiTagPart(item)
			if err != nil {
				return "", fmt.Errorf("problem in ascii tag part: %w", err)
			}
			if tag_str == "" {
				tag_str = s
			} else {
				tag_str = fmt.Sprintf("%s.%s", tag_str, s)
			}
		default:
			// this byte does not indicate the tag path is continuing.  go back by one in the item's data buffer to "unread" it.
			item.Pos--
			break morepath
		}

	}

	return tag_str, nil

}
