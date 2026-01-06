package l5x

import (
	"fmt"
	"strconv"
)

func L5xTypeToGoType(typestr string, valuestr string) (any, error) {
	switch typestr {
	case "REAL":
		value, _ := strconv.ParseFloat(valuestr, 64)
		return float64(value), nil
	case "DINT":
		value, _ := strconv.ParseInt(valuestr, 10, 32)
		return int32(value), nil
	case "BOOL", "BIT":
		value, _ := strconv.ParseBool(valuestr)
		return bool(value), nil
	case "INT":
		value, _ := strconv.ParseInt(valuestr, 10, 16)
		return int16(value), nil
	case "STRING":
		return valuestr, nil
	case "SINT":
		value, _ := strconv.ParseInt(valuestr, 10, 8)
		return int8(value), nil
	case "LINT":
		value, _ := strconv.ParseInt(valuestr, 10, 64)
		return int64(value), nil
	case "BYTE":
		value, _ := strconv.ParseInt(valuestr, 10, 8)
		return uint8(value), nil
	case "WORD":
		value, _ := strconv.ParseInt(valuestr, 10, 16)
		return uint16(value), nil
	case "DWORD":
		value, _ := strconv.ParseInt(valuestr, 10, 32)
		return uint32(value), nil
	case "LWORD":
		value, _ := strconv.ParseInt(valuestr, 10, 64)
		return uint64(value), nil

	default:
		return fmt.Sprintf("Unknown type: %s", typestr), nil
	}
}
