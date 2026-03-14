package gologix

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
)

// dataTableTypeSize maps a CIPType to the raw byte count used by the datatable
// protocol's dataTypeSize field (LE16). This is a byte count, not a CIP type code.
//
// For most atomic types this matches CIPType.Size(). The exception is CIPTypeSTRING
// (gologix internal 0xFF) which returns Size()=1 but the datatable protocol expects
// 88 bytes (4-byte .LEN + 82-byte .DATA + 2 pad for standard AB STRING).
func dataTableTypeSize(t CIPType) (uint16, error) {
	switch t {
	case CIPTypeSTRING:
		return 88, nil
	case CIPTypeStruct:
		return 88, nil
	default:
		s := t.Size()
		if s == 0 {
			return 0, fmt.Errorf("unsupported CIP type %v for datatable buffer", t)
		}
		return uint16(s), nil
	}
}

// resolveTypeSize determines the datatable slot size for a tag. For STRING and
// Struct types it checks client.KnownTags for the tag's UDT descriptor, which
// carries the actual template data size (msgGetTemplateAttrListResponse.SizeBytes,
// attribute 5 of the CIP template object). This handles custom string types like
// STRING20 or STRING40 whose slot sizes differ from the standard 88-byte AB STRING.
//
// Falls back to dataTableTypeSize() if the tag is not in KnownTags or has no UDT.
func (client *Client) resolveTypeSize(tagName string, cipType CIPType) (uint16, error) {
	if cipType == CIPTypeSTRING || cipType == CIPTypeStruct {
		kt, ok := client.KnownTags[strings.ToLower(tagName)]
		if ok && kt.UDT != nil && kt.UDT.Info.SizeBytes > 0 {
			return uint16(kt.UDT.Info.SizeBytes), nil
		}
	}
	return dataTableTypeSize(cipType)
}

// decodeDataTableValue converts raw bytes from a datatable read response into
// the appropriate Go type based on the CIP type.
//
// For STRING types (88 bytes), the format is: uint32 length + char data + padding.
// For all other atomic types, this delegates to the existing readValue() function.
func decodeDataTableValue(cipType CIPType, raw []byte) (any, error) {
	if cipType == CIPTypeSTRING || cipType == CIPTypeStruct {
		if len(raw) < 4 {
			return nil, fmt.Errorf("STRING data too short: got %d bytes, need at least 4", len(raw))
		}
		strLen := binary.LittleEndian.Uint32(raw[0:4])
		if strLen > uint32(len(raw)-4) {
			strLen = uint32(len(raw) - 4)
		}
		if strLen > 82 {
			strLen = 82
		}
		return string(raw[4 : 4+strLen]), nil
	}
	r := bytes.NewReader(raw)
	return readValue(cipType, r)
}

type PLCStructTrend[T any] struct {
	instanceID CIPInstance
	client     *Client
	Prefix1    bool // true to use the single-tag timestamped mode.  false to use the multi-tag mode without timestamps.
}

type StructTrendSample[T any] struct {
	Timestamp time.Duration
	Data      T
}

func (trend *PLCStructTrend[T]) Path() (*bytes.Buffer, error) {
	path, err := Serialize(CipObject_DataTable, trend.instanceID)
	if err != nil {
		return nil, fmt.Errorf("could not serialize path for trend read: %w", err)
	}
	return path, nil
}

// A samplerate of 0 means you will only get a single sample every time you read.
func NewStructTrend[T any](client *Client, sampleRate time.Duration, bufferSize uint16, prefix bool, args ...any) (*PLCStructTrend[T], error) {
	var trend = new(PLCStructTrend[T])
	trend.client = client
	trend.Prefix1 = prefix

	var str T

	// Reflect over struct fields and figure out everything we need to add tags for them
	// (names, types, IOIs) before we create the trend instance on the PLC
	Typ := reflect.TypeOf(str)
	if Typ.Kind() == reflect.Ptr {
		Typ = Typ.Elem()
	}
	v := reflect.ValueOf(str)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	vf := reflect.VisibleFields(Typ)
	tagList := make([]string, 0, len(vf))
	tagTypes := make([]CIPType, 0, len(vf))

	for i := range vf {
		field := vf[i]
		tagPath, ok := field.Tag.Lookup("gologix")
		if !ok || tagPath == "" {
			continue
		}
		if args != nil {
			tagPath = formatName(tagPath, args...)
		}
		fieldVal := v.Field(i)
		ct, elements := GoVarToCIPType(fieldVal.Interface())

		// Handle slices/arrays (multi-element tags)
		if elements > 1 {
			expanded := expandArrayTag(tagPath, elements)
			for _, name := range expanded {
				tagList = append(tagList, name)
				tagTypes = append(tagTypes, ct)
				// For slices, set Var to the element pointer
			}
		} else {
			tagList = append(tagList, tagPath)
			tagTypes = append(tagTypes, ct)
		}
	}

	iois := make([][]byte, len(tagList))
	for i, tagName := range tagList {
		ioi, err := trend.client.newIOI(tagName, CIPTypeDINT) // TODO: determine type from struct field
		if err != nil {
			return nil, fmt.Errorf("could not create IOI for tag %q: %w", tagName, err)
		}
		iois[i] = ioi.Bytes()
	}

	// ------------
	// Create the trend instance data table
	// ------------
	path, err := Serialize(CipObject_DataTable, CIPInstance(0))
	if err != nil {
		return nil, fmt.Errorf("could not serialize path for trend creation: %w", err)
	}

	data, err := attrValueBytes(
		AttributeValue{Attribute: 8, Value: uint32(bufferSize)}, // Buffer size in bytes.
		//AttributeValue{Attribute: 3, Value: uint8(1)},           // Buffer size in bytes.
	)
	if err != nil {
		return nil, fmt.Errorf("could not serialize attribute values: %w", err)
	}

	resp, err := client.GenericCIPMessage(CIPService_Create, path.Bytes(), data.Bytes())
	if err != nil {
		return nil, fmt.Errorf("datatable buffer create failed: %w", err)
	}

	instanceID, err := resp.Uint16()
	if err != nil {
		return nil, fmt.Errorf("could not read instance ID from create response: %w", err)
	}
	trend.instanceID = CIPInstance(instanceID)

	path, err = trend.Path()
	if err != nil {
		return nil, fmt.Errorf("could not serialize path with trend id %d: %w", trend.instanceID, err)
	}

	// ------------
	// Set up some additional attributes on the trend
	// ------------

	// Service 0x04 = SetAttributeList
	err = trend.client.SetAttrList(CipObject_DataTable, CIPInstance(trend.instanceID),
		AttributeValue{Attribute: 1, Value: uint32(sampleRate.Microseconds())}, // sample rate in microseconds
		AttributeValue{Attribute: 5, Value: uint8(0)},                          // state (0=stopped)
	)
	if err != nil {
		return nil, fmt.Errorf("SetTrendAttrs CIP request failed: %w", err)
	}

	// ------------
	// check some stuff
	// ------------

	resp, err = client.GetAttrList(CipObject_DataTable, trend.instanceID, 1, 3, 5, 6, 7, 8, 10)
	if err != nil {
		return nil, fmt.Errorf("could not read back trend attributes: %w", err)
	}
	_ = resp

	// ------------
	// Add tags to the trend
	// ------------

	// Estimate capacity: 5 header + per-tag (2 typeSize + 1 ioiLen + IOI bytes)
	capacity := 5
	for _, e := range iois {
		capacity += 3 + len(e)
	}
	msgData := make([]byte, 0, capacity)

	// Header: 02 00 01 01
	if trend.Prefix1 {
		msgData = append(msgData, 0x01, 0x00, 0x01) // - this prefix gives you a timestamp but only allows one tag.  No data type size is required.
	} else {
		msgData = append(msgData, 0x02, 0x00, 0x01, 0x01) // - this prefix allows multiple tags but you don't get a timestamp.  You also have to include each data type size before the IOI.
	}
	// Tag count
	msgData = append(msgData, byte(len(iois)))

	for i, e := range iois {
		// Per-tag: dataTypeSize LE16
		sizeBytes := make([]byte, 2)
		s, err := client.resolveTypeSize(tagList[i], tagTypes[i])
		if err != nil {
			return nil, fmt.Errorf("could not resolve type size for tag %q: %w", tagList[i], err)
		}
		if !trend.Prefix1 {
			binary.LittleEndian.PutUint16(sizeBytes, s)
			msgData = append(msgData, sizeBytes...)
		}
		// Per-tag: IOI length in 16-bit words
		msgData = append(msgData, byte(len(e)/2))
		// Per-tag: IOI EPath bytes
		msgData = append(msgData, e...)
	}

	msgData = append(msgData, []byte{0xFF, 0xFF, 0xFF, 0xFF}...)

	// Service 0x4E — Add Tag when targeting class 0xB2 buffer instance.
	resp, err = client.GenericCIPMessage(CIPService(0x4E), path.Bytes(), msgData)
	if err != nil {
		return nil, fmt.Errorf("datatable add tag failed: %w", err)
	}

	batchIndex, err := resp.Uint16()
	if err != nil {
		return nil, fmt.Errorf("could not read batch index from add tag response: %w", err)
	}
	if batchIndex != 1 {
		return nil, fmt.Errorf("unexpected batch index in add tag response: got %d, want 1", batchIndex)
	}

	return trend, nil
}

func (trend *PLCStructTrend[T]) StartTrend() error {
	path, err := trend.Path()
	if err != nil {
		return fmt.Errorf("could not serialize path for trend start: %w", err)
	}
	_, err = trend.client.GenericCIPMessage(CIPService_Start, path.Bytes(), nil)
	return err
}

func (trend *PLCStructTrend[T]) StopTrend() error {
	path, err := trend.Path()
	if err != nil {
		return fmt.Errorf("could not serialize path for trend stop: %w", err)
	}
	_, err = trend.client.GenericCIPMessage(CIPService_Stop, path.Bytes(), nil)
	return err
}

func (trend *PLCStructTrend[T]) continueRead(path []byte) (*CIPItem, error) {
	//start := make([]byte, 2)
	//binary.LittleEndian.PutUint16(start, uint16(offset))
	return trend.client.GenericCIPMessage(CIPService_Read, path, nil)
}

func (trend *PLCStructTrend[T]) ReadAll() ([]StructTrendSample[T], error) {
	err := trend.client.checkConnection()
	if err != nil {
		return nil, fmt.Errorf("could not read datatable buffer: %w", err)
	}

	path, err := trend.Path()
	if err != nil {
		return nil, fmt.Errorf("could not serialize path for trend read: %w", err)
	}

	resp, err := trend.client.GenericCIPMessage(CIPService_Read, path.Bytes(), []byte{})
	if err != nil {
		statErr, ok := err.(CIPStatusError)
		if ok {
			if statErr.Status == CIPStatus_PartialTransfer {
				fmt.Printf("GOTTA GET THAT DATA BRO!!")
				resp2, err := trend.continueRead(path.Bytes())
				if err != nil {
					return nil, fmt.Errorf("error during fragmented read: %w", err)
				}
				resp.Data = append(resp.Data, resp2.Data[resp2.Pos:]...)
			}
		} else {

			return nil, fmt.Errorf("datatable buffer read failed: %w", err)
		}
	}
	results := make([]StructTrendSample[T], 0)
	for {
		// Read until we run out of data or hit an error. Each read should give us a full struct's worth of data.
		str, err := decodeStructTrendValue[T](resp, trend.Prefix1)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break // no more complete structs to read
			}
			return nil, fmt.Errorf("could not decode trend value: %w", err)
		}
		results = append(results, str)
	}

	return results, nil
}

func decodeStructTrendValue[T any](resp *CIPItem, timestamp bool) (StructTrendSample[T], error) {
	var str T
	var result StructTrendSample[T]
	Typ := reflect.TypeOf(str)
	if Typ.Kind() == reflect.Ptr {
		Typ = Typ.Elem()
	}
	v := reflect.ValueOf(&str).Elem()

	vf := reflect.VisibleFields(Typ)

	dtIndex, err := resp.Uint16()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return result, io.EOF // signal that we've read all available data
		}
		return result, fmt.Errorf("could not read data type index from response: %w", err)
	}
	if dtIndex != 1 {
		return result, fmt.Errorf("unexpected index in response: got %d, want 1", dtIndex)
	}

	if timestamp {
		ts, err := resp.Uint32()
		if err != nil {
			return result, fmt.Errorf("could not read timestamp from response: %w", err)
		}
		result.Timestamp = time.Duration(ts) * time.Microsecond
	}

	for i := range vf {
		field := vf[i]
		tagPath, ok := field.Tag.Lookup("gologix")
		if !ok || tagPath == "" {
			continue
		}
		ct, elements := GoVarToCIPType(v.Field(i).Interface())
		if elements > 1 {
			for j := 0; j < int(elements); j++ {
				value, err := decodeDataTableValue( /*
						ts, err := resp.Uint32()
						if err != nil {
							return result, fmt.Errorf("could not read timestamp from response: %w", err)
						}
						result.Timestamp = time.Unix(int64(ts)*1000/128, 0)
					*/ct, resp.Data[resp.Pos:])
				if err != nil {
					return result, fmt.Errorf("could not decode value for field %q element %d: %w", field.Name, j, err)
				}
				v.Field(i).Set(reflect.Append(v.Field(i), reflect.ValueOf(value)))
				resp.Pos += int(ct.Size())
			}
		} else {
			value, err := decodeDataTableValue(ct, resp.Data[resp.Pos:])
			if err != nil {
				return result, fmt.Errorf("could not decode value for field %q: %w", field.Name, err)
			}
			v.Field(i).Set(reflect.ValueOf(value))
			resp.Pos += int(ct.Size())
		}
	}

	result.Data = str
	return result, nil
}
