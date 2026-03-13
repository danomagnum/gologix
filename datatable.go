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

// DataTableBufferTag tracks a single tag that has been added to a datatable buffer.
// The index is 1-based, sequential in the order tags were added across all AddTag/AddTags calls.
type DataTableBufferTag struct {
	Name         string  // symbolic tag name (or "" for instance-addressed tags)
	CIPType      CIPType // the CIP data type of the tag
	DataTypeSize uint16  // raw byte count used by the datatable protocol
	index        uint16  // 1-based sequential position in the buffer
	ioi          []byte  // raw IOI EPath bytes
	batchID      uint16  // 0x4E response value — groups tags for 0x4C read parsing
	synched      bool    // internal flag to track if this tag has been added to the buffer on the PLC side.
	Var          any     // optional field to store the value of the tag after reading from the buffer, not used internally by the library
}

// Index returns the 1-based PLC-assigned index for this tag in the buffer.
func (t DataTableBufferTag) Index() uint16 { return t.index }

// DataTableBuffer represents a CIP class 0xB2 datatable buffer instance on the PLC.
// It tracks the buffer's lifecycle state, the instance ID returned by Create,
// and the list of tags that have been added along with their type sizes.
//
// The buffer is read-only from the client side. Tags are associated with the buffer
// via AddTag/AddTags/AddTagGroup, and all values are read at once via ReadAll or
// ReadAllTyped. To write tag values, use the standard Client.Write() method directly.
//
// DataTableBuffer is NOT safe for concurrent use from multiple goroutines.
type DataTableBuffer struct {
	client     *Client
	instanceID uint16               // assigned by PLC during Create; use InstanceID() getter
	tags       []DataTableBufferTag // tags added to the buffer, in order; use Tags() getter
	// expansions tracks multi-element TagDef expansions from AddTagGroup.
	// Key: the base resolved tag name (e.g. "EXAMPLE[1,0]")
	// Value: the expanded individual tag names (e.g. ["EXAMPLE[1,0]", ..., "EXAMPLE[1,4]"])
	expansions map[string][]string

	created bool // true after successful Create
	synched bool // internal flag to track if tags have been added on this side but not yet sent to PLC via AddTag/AddTags
}

// InstanceID returns the PLC-assigned instance ID for this datatable buffer.
func (buf *DataTableBuffer) InstanceID() uint16 { return buf.instanceID }

// Tags returns a copy of the tags added to this buffer.
func (buf *DataTableBuffer) Tags() []DataTableBufferTag {
	out := make([]DataTableBufferTag, len(buf.tags))
	copy(out, buf.tags)
	return out
}

// dataTablePath builds the 6-byte CIP path for class 0xB2 with a 16-bit instance ID.
// Forces 16-bit instance encoding (0x25) to match the reverse-engineered protocol.
func dataTablePath(instanceID uint16) []byte {
	path := []byte{0x20, 0xB2, 0x25, 0x00, 0x00, 0x00}
	binary.LittleEndian.PutUint16(path[4:], instanceID)
	return path
}

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

// ---------------------------------------------------------------------------
// Low-level Client methods — one per CIP service
// ---------------------------------------------------------------------------

// CreateDataTableBuffer sends a CIP Create (0x08) service to class 0xB2 instance 0x0000,
// requesting the PLC to allocate a new datatable buffer.
//
// The request data is hardcoded to 01 00 03 00 03 00.
// Returns the instance ID assigned by the PLC.
func (client *Client) CreateDataTableBuffer() (uint16, error) {
	err := client.checkConnection()
	if err != nil {
		return 0, fmt.Errorf("could not create datatable buffer: %w", err)
	}

	path := dataTablePath(0x0000)
	//data := []byte{0x01, 0x00, 0x03, 0x00, 0x03, 0x00}
	data := []byte{0x02, 0x00, 0x08, 0x00, 0x00, 0x10, 0x00, 0x00, 0x03, 0x00, 0x01}
	//data := []byte{}
	// Request format from pcap: [02 00][08 00][bufsize:4LE][03 00][num_tags:1]

	resp, err := client.GenericCIPMessage(CIPService_Create, path, data)
	if err != nil {
		return 0, fmt.Errorf("datatable buffer create failed: %w", err)
	}

	instanceID, err := resp.Uint16()
	if err != nil {
		return 0, fmt.Errorf("could not read instance ID from create response: %w", err)
	}

	return instanceID, nil
}

// DeleteDataTableBuffer sends a CIP Delete (0x09) service to class 0xB2 with
// the specified instance ID, freeing the datatable buffer on the PLC.
func (client *Client) DeleteDataTableBuffer(instanceID uint16) error {
	err := client.checkConnection()
	if err != nil {
		return fmt.Errorf("could not delete datatable buffer: %w", err)
	}

	path := dataTablePath(instanceID)

	_, err = client.GenericCIPMessage(CIPService_Delete, path, []byte{})
	if err != nil {
		return fmt.Errorf("datatable buffer delete failed: %w", err)
	}
	return nil
}

// addTagEntry describes one tag to add in a multi-tag 0x4E call.
type addTagEntry struct {
	dataTypeSize uint16
	ioi          []byte
	index        int // the position of this tag in the buffer's tags slice, used to update batchID and synched after a successful add
}

// DataTableReadBuffer sends a CIP Read (0x4C) service to the specified datatable
// buffer instance, reading all tag values in a single response.
//
// The returned CIPItem contains: [tagCount LE16] [concatenated tag values...].
// The caller must know the order and sizes of tags to parse the response.
func (client *Client) DataTableReadBuffer(instanceID uint16) (*CIPItem, error) {
	err := client.checkConnection()
	if err != nil {
		return nil, fmt.Errorf("could not read datatable buffer: %w", err)
	}

	path := dataTablePath(instanceID)

	resp, err := client.GenericCIPMessage(CIPService_Read, path, []byte{})
	if err != nil {
		return nil, fmt.Errorf("datatable buffer read failed: %w", err)
	}
	return resp, nil
}

// DataTableRemoveTag sends a CIP service 0x4F to remove a tag from the datatable
// buffer by its 1-based index.
func (client *Client) DataTableRemoveTag(instanceID uint16, tagIndex uint16) error {
	err := client.checkConnection()
	if err != nil {
		return fmt.Errorf("could not remove tag from datatable buffer: %w", err)
	}

	path := dataTablePath(instanceID)

	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, tagIndex)

	_, err = client.GenericCIPMessage(CIPService_RemoveTagFromBuffer, path, data)
	if err != nil {
		return fmt.Errorf("datatable remove tag failed: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// High-level API
// ---------------------------------------------------------------------------

// NewDataTableBuffer creates a new datatable buffer on the PLC by sending a CIP
// Create (0x08) service to class 0xB2. The returned DataTableBuffer tracks the
// buffer instance and all tags added to it.
//
// Call Close() when finished to delete the buffer and free PLC resources.
//
// Example:
//
//	buf, err := client.NewDataTableBuffer()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer buf.Close()
//
//	buf.AddTag("MyDINT", gologix.CIPTypeDINT)
//	values, err := buf.ReadAll()
func (client *Client) NewDataTableBuffer() (*DataTableBuffer, error) {
	instanceID, err := client.CreateDataTableBuffer()
	if err != nil {
		return nil, err
	}
	return &DataTableBuffer{
		client:     client,
		instanceID: instanceID,
		tags:       make([]DataTableBufferTag, 0),
		created:    true,
		expansions: make(map[string][]string),
	}, nil
}

func setGoVar(dest any, value any) error {
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr {
		return fmt.Errorf("Var field must be a pointer to set value, got %T", dest)
	}
	if !destVal.Elem().CanSet() {
		return fmt.Errorf("cannot set value on Var field of type %T", dest)
	}
	valueVal := reflect.ValueOf(value)
	if !valueVal.Type().AssignableTo(destVal.Elem().Type()) {
		return fmt.Errorf("cannot assign value of type %T to Var field of type %T", value, dest)
	}
	destVal.Elem().Set(valueVal)
	return nil
}

// RemoveTag removes a tag from the buffer by its name. Sends service 0x4F with
// the tag's recorded index and updates the internal tag list.
func (buf *DataTableBuffer) RemoveTag(tagName string) error {
	if !buf.created {
		return fmt.Errorf("datatable buffer not created")
	}

	idx := -1
	for i, tag := range buf.tags {
		if tag.Name == tagName {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("tag %q not found in buffer", tagName)
	}

	err := buf.client.DataTableRemoveTag(buf.instanceID, buf.tags[idx].index)
	if err != nil {
		return fmt.Errorf("could not remove tag %q from buffer: %w", tagName, err)
	}

	buf.tags = append(buf.tags[:idx], buf.tags[idx+1:]...)
	return nil
}

// Close deletes the datatable buffer on the PLC, freeing resources.
// After Close returns, the buffer instance is invalid and should not be reused.
// Close is idempotent — calling it on an already-closed buffer is a no-op.
func (buf *DataTableBuffer) Close() error {
	if !buf.created {
		return nil
	}
	err := buf.client.DeleteDataTableBuffer(buf.instanceID)
	buf.created = false
	if err != nil {
		return fmt.Errorf("could not delete datatable buffer: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Struct Tag integration
// ---------------------------------------------------------------------------

// AddTaggedStruct adds all tags from a tagged go struct (with `gologix` struct tags)
// to this datatable buffer. Args substitute {0}, {1}, ... placeholders in tag names.
//
// Automatically sets references to struct fields in the Var field of each buffer tag,
// so that after reading the struct fields will be populated with the read values.
//
// Can be called multiple times with different args to build up a large buffer:
//
// Example:
//
//	    type MyTags struct {
//	    IntTag    int16     `gologix:"Machine{0}TestInt"`
//	    RealTag   float32   `gologix:"Machine{0}TestReal"`
//	    ArrayTag  []int32   `gologix:"Machine{0}TestDintArr[2]"`  // Read 5 elements starting at index 2
//	    }
//	    var tags1 MyTags
//	    var tags2 MyTags
//		buf.AddTaggedStruct(&tags1, 1)  // inspection point 1
//		buf.AddTaggedStruct(&tags2, 2)  // inspection point 2
//		// buf now has all tags for both structs
func (trend *PLCStructTrend[T]) AddTaggedStruct(args ...any) error {
	var str T

	// Reflect over struct fields
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
		_, elements := GoVarToCIPType(fieldVal.Interface())

		// Handle slices/arrays (multi-element tags)
		if elements > 1 {
			expanded := expandArrayTag(tagPath, elements)
			for _, name := range expanded {
				tagList = append(tagList, name)
				// For slices, set Var to the element pointer
			}
		} else {
			tagList = append(tagList, tagPath)
		}
	}

	iois := make([][]byte, len(tagList))
	for i, tagName := range tagList {
		ioi, err := trend.client.newIOI(tagName, CIPTypeDINT) // TODO: determine type from struct field
		if err != nil {
			return fmt.Errorf("could not create IOI for tag %q: %w", tagName, err)
		}
		iois[i] = ioi.Bytes()
	}

	return nil
}

type PLCStructTrend[T any] struct {
	instanceID uint16
	client     *Client
}

func (trend *PLCStructTrend[T]) Path() (*bytes.Buffer, error) {
	path, err := Serialize(CipObject_DataTable, CIPInstance(trend.instanceID))
	if err != nil {
		return nil, fmt.Errorf("could not serialize path for trend read: %w", err)
	}
	return path, nil
}

type attribute struct {
	id   uint16
	data any
}

func attrList(attrs ...attribute) []byte {
	var data bytes.Buffer

	err := binary.Write(&data, binary.LittleEndian, uint16(len(attrs)))
	if err != nil {
		return nil
	}
	for _, a := range attrs {
		err = binary.Write(&data, binary.LittleEndian, a.id)
		if err != nil {
			return nil
		}
		err = binary.Write(&data, binary.LittleEndian, a.data)
		if err != nil {
			return nil
		}
	}
	return data.Bytes()
}

// A samplerate of 0 means you will only get a single sample every time you read.
func NewStructTrend[T any](client *Client, sampleRate time.Duration, bufferSize uint16, args ...any) (*PLCStructTrend[T], error) {
	var trend = new(PLCStructTrend[T])
	trend.client = client

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

	// TODO: only count fields with gologix struct tags
	fieldCount := uint16(len(tagList))

	// ------------
	// Create the trend instance data table
	// ------------
	path, err := Serialize(CipObject_DataTable, CIPInstance(0))
	if err != nil {
		return nil, fmt.Errorf("could not serialize path for trend creation: %w", err)
	}
	data := attrList(
		attribute{id: 8, data: uint32(bufferSize)}, // buffer size
		attribute{id: 3, data: uint32(fieldCount)}, // number of tags in the trend
	)

	resp, err := client.GenericCIPMessage(CIPService_Create, path.Bytes(), data)
	if err != nil {
		return nil, fmt.Errorf("datatable buffer create failed: %w", err)
	}

	trend.instanceID, err = resp.Uint16()
	if err != nil {
		return nil, fmt.Errorf("could not read instance ID from create response: %w", err)
	}

	path, err = trend.Path()
	if err != nil {
		return nil, fmt.Errorf("could not serialize path with trend id %d: %w", trend.instanceID, err)
	}

	// ------------
	// Set up some additional attributes on the trend
	// ------------
	data = attrList(
		attribute{id: 1, data: uint32(sampleRate.Microseconds())}, // sample rate in microseconds
		attribute{id: 5, data: uint8(0)},                          // state (0=stopped)
	)

	// Service 0x04 = SetAttributeList
	_, err = trend.client.GenericCIPMessage(CIPService(0x04), path.Bytes(), data)
	if err != nil {
		return nil, fmt.Errorf("SetTrendAttrs CIP request failed: %w", err)
	}

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
	msgData = append(msgData, 0x02, 0x00, 0x01, 0x01)
	// Tag count
	msgData = append(msgData, byte(len(iois)))

	for i, e := range iois {
		// Per-tag: dataTypeSize LE16
		sizeBytes := make([]byte, 2)
		s, err := client.resolveTypeSize(tagList[i], tagTypes[i])
		if err != nil {
			return nil, fmt.Errorf("could not resolve type size for tag %q: %w", tagList[i], err)
		}
		binary.LittleEndian.PutUint16(sizeBytes, s)
		msgData = append(msgData, sizeBytes...)
		// Per-tag: IOI length in 16-bit words
		msgData = append(msgData, byte(len(e)/2))
		// Per-tag: IOI EPath bytes
		msgData = append(msgData, e...)
	}

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
	path := dataTablePath(trend.instanceID)
	_, err := trend.client.GenericCIPMessage(CIPService_Start, path, nil)
	return err
}

func (trend *PLCStructTrend[T]) StopTrend() error {
	path := dataTablePath(trend.instanceID)
	_, err := trend.client.GenericCIPMessage(CIPService_Stop, path, nil)
	return err
}

func (trend *PLCStructTrend[T]) ReadAll() ([]T, error) {
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
		return nil, fmt.Errorf("datatable buffer read failed: %w", err)
	}
	results := make([]T, 0)
	for {
		// Read until we run out of data or hit an error. Each read should give us a full struct's worth of data.
		str, err := decodeStructTrendValue[T](resp)
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

func decodeStructTrendValue[T any](resp *CIPItem) (T, error) {
	var str T
	Typ := reflect.TypeOf(str)
	if Typ.Kind() == reflect.Ptr {
		Typ = Typ.Elem()
	}
	v := reflect.ValueOf(&str).Elem()

	vf := reflect.VisibleFields(Typ)

	for i := range vf {
		field := vf[i]
		tagPath, ok := field.Tag.Lookup("gologix")
		if !ok || tagPath == "" {
			continue
		}
		ct, elements := GoVarToCIPType(v.Field(i).Interface())
		if elements > 1 {
			for j := 0; j < int(elements); j++ {
				value, err := decodeDataTableValue(ct, resp.Data[resp.Pos:])
				if err != nil {
					return str, fmt.Errorf("could not decode value for field %q element %d: %w", field.Name, j, err)
				}
				v.Field(i).Set(reflect.Append(v.Field(i), reflect.ValueOf(value)))
				resp.Pos += int(ct.Size())
			}
		} else {
			value, err := decodeDataTableValue(ct, resp.Data[resp.Pos:])
			if err != nil {
				return str, fmt.Errorf("could not decode value for field %q: %w", field.Name, err)
			}
			v.Field(i).Set(reflect.ValueOf(value))
			resp.Pos += int(ct.Size())
		}
	}
	return str, nil
}
