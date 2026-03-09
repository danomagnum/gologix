package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"
	"strings"
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
	created    bool                 // true after successful Create
	// expansions tracks multi-element TagDef expansions from AddTagGroup.
	// Key: the base resolved tag name (e.g. "EXAMPLE[1,0]")
	// Value: the expanded individual tag names (e.g. ["EXAMPLE[1,0]", ..., "EXAMPLE[1,4]"])
	expansions map[string][]string
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
func (buf *DataTableBuffer) resolveTypeSize(tagName string, cipType CIPType) (uint16, error) {
	if cipType == CIPTypeSTRING || cipType == CIPTypeStruct {
		kt, ok := buf.client.KnownTags[strings.ToLower(tagName)]
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
	data := []byte{0x01, 0x00, 0x03, 0x00, 0x03, 0x00}

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
}

// DataTableAddTags sends a CIP service 0x4E to the specified datatable buffer instance,
// adding one or more tags in a single call. Each tag carries its own dataTypeSize and
// IOI EPath — mixed type sizes are supported.
//
// Wire format (validated against pcap):
//
//	[02 00 01 01]         hardcoded header
//	[1 byte: tagCount]    number of tags in this call
//	Repeated tagCount times:
//	  [2 bytes LE: dataTypeSize]  this tag's slot size in bytes
//	  [1 byte: ioiLenWords]       this tag's IOI length in 16-bit words
//	  [N bytes: IOI EPath]        this tag's path bytes
//
// The response contains a single uint16: a batch/entry index that increments
// per 0x4E call (1 for the first call, 2 for the second, etc.). This value
// is used as a group marker in the 0x4C read response.
func (client *Client) DataTableAddTags(instanceID uint16, entries []addTagEntry) (uint16, error) {
	err := client.checkConnection()
	if err != nil {
		return 0, fmt.Errorf("could not add tag to datatable buffer: %w", err)
	}

	path := dataTablePath(instanceID)

	// Estimate capacity: 5 header + per-tag (2 typeSize + 1 ioiLen + IOI bytes)
	capacity := 5
	for _, e := range entries {
		capacity += 3 + len(e.ioi)
	}
	msgData := make([]byte, 0, capacity)

	// Header: 02 00 01 01
	msgData = append(msgData, 0x02, 0x00, 0x01, 0x01)
	// Tag count
	msgData = append(msgData, byte(len(entries)))

	for _, e := range entries {
		// Per-tag: dataTypeSize LE16
		sizeBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(sizeBytes, e.dataTypeSize)
		msgData = append(msgData, sizeBytes...)
		// Per-tag: IOI length in 16-bit words
		msgData = append(msgData, byte(len(e.ioi)/2))
		// Per-tag: IOI EPath bytes
		msgData = append(msgData, e.ioi...)
	}

	// Service 0x4E — Add Tag when targeting class 0xB2 buffer instance.
	resp, err := client.GenericCIPMessage(CIPService(0x4E), path, msgData)
	if err != nil {
		return 0, fmt.Errorf("datatable add tag failed: %w", err)
	}

	batchIndex, err := resp.Uint16()
	if err != nil {
		return 0, fmt.Errorf("could not read batch index from add tag response: %w", err)
	}
	return batchIndex, nil
}

// DataTableAddTag sends a CIP service 0x4E to add a single tag to the datatable buffer.
// This is a convenience wrapper around DataTableAddTags for single-tag adds.
func (client *Client) DataTableAddTag(instanceID uint16, dataTypeSize uint16, tagIOIs []byte) (uint16, error) {
	return client.DataTableAddTags(instanceID, []addTagEntry{{dataTypeSize, tagIOIs}})
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

// AddTag adds a single tag to the datatable buffer by its symbolic name and CIP type.
// The tag is sent to the PLC via service 0x4E and the PLC-assigned batch index is recorded.
//
// The IOI is built using the client's standard IOI encoding, which supports both
// symbolic (0x91 segments) and instance-optimized (class 0x6B) addressing.
func (buf *DataTableBuffer) AddTag(tagName string, cipType CIPType) error {
	if !buf.created {
		return fmt.Errorf("datatable buffer not created")
	}

	ioi, err := buf.client.newIOI(tagName, cipType)
	if err != nil {
		return fmt.Errorf("could not create IOI for tag %q: %w", tagName, err)
	}

	typeSize, err := buf.resolveTypeSize(tagName, cipType)
	if err != nil {
		return fmt.Errorf("tag %q: %w", tagName, err)
	}

	batchID, err := buf.client.DataTableAddTag(buf.instanceID, typeSize, ioi.Bytes())
	if err != nil {
		return fmt.Errorf("could not add tag %q to buffer: %w", tagName, err)
	}

	buf.tags = append(buf.tags, DataTableBufferTag{
		Name:         tagName,
		CIPType:      cipType,
		DataTypeSize: typeSize,
		index:        uint16(len(buf.tags) + 1),
		ioi:          ioi.Bytes(),
		batchID:      batchID,
	})

	return nil
}

// AddTags adds multiple tags to the datatable buffer in a single 0x4E call.
// Each tag carries its own dataTypeSize and IOI EPath — mixed type sizes are
// supported in one call (confirmed via pcap).
//
// Because tags is a map, iteration order (and therefore the tag position in the
// buffer) is non-deterministic. Use AddTag in a loop if you need a specific ordering.
func (buf *DataTableBuffer) AddTags(tags map[string]CIPType) error {
	if !buf.created {
		return fmt.Errorf("datatable buffer not created")
	}

	type preparedTag struct {
		name     string
		cipType  CIPType
		typeSize uint16
		ioi      []byte
	}

	// Resolve IOI and type size for each tag.
	prepared := make([]preparedTag, 0, len(tags))
	for name, cipType := range tags {
		ioi, err := buf.client.newIOI(name, cipType)
		if err != nil {
			return fmt.Errorf("could not create IOI for tag %q: %w", name, err)
		}
		typeSize, err := buf.resolveTypeSize(name, cipType)
		if err != nil {
			return fmt.Errorf("tag %q: %w", name, err)
		}
		prepared = append(prepared, preparedTag{
			name:     name,
			cipType:  cipType,
			typeSize: typeSize,
			ioi:      ioi.Bytes(),
		})
	}

	// Build entries for a single DataTableAddTags call.
	entries := make([]addTagEntry, len(prepared))
	for i, t := range prepared {
		entries[i] = addTagEntry{
			dataTypeSize: t.typeSize,
			ioi:          t.ioi,
		}
	}

	batchID, err := buf.client.DataTableAddTags(buf.instanceID, entries)
	if err != nil {
		return fmt.Errorf("could not add tags to buffer: %w", err)
	}

	// Assign sequential indices and record the batch ID for read parsing.
	baseIndex := uint16(len(buf.tags) + 1)
	for i, t := range prepared {
		buf.tags = append(buf.tags, DataTableBufferTag{
			Name:         t.name,
			CIPType:      t.cipType,
			DataTypeSize: t.typeSize,
			index:        baseIndex + uint16(i),
			ioi:          t.ioi,
			batchID:      batchID,
		})
	}

	return nil
}

// ReadAll reads all tag values from the datatable buffer in a single CIP request
// (service 0x4C) and returns them as a map from tag name to the decoded Go value.
//
// # Response layout
//
// The PLC returns a byte stream grouped by 0x4E add-tag calls (batches).
// For each batch (in the order the 0x4E calls were made):
//
//	[2 bytes LE] batch/entry index (1, 2, ...)
//	For each tag in that batch (in the order they were added):
//	  [DataTypeSize bytes] tag value
//
// The 2-byte entry header appears once per 0x4E call, NOT per tag. When all
// tags are added in a single AddTags call, there is exactly one entry header
// followed by all tag values concatenated. When using AddTag individually,
// each tag gets its own entry header.
//
// Tags are sorted by index (insertion order) to match the PLC's output order.
// We track batch boundaries via each tag's batchID field, skipping the 2-byte
// entry header whenever we encounter a new batch.
//
// For STRING/Struct types the slot is a fixed size determined at AddTag time
// (e.g. 88 bytes for standard AB STRING), but the actual string content is
// variable-length within that slot: [4-byte uint32 LEN][LEN chars][padding].
// decodeDataTableValue reads LEN to extract only the meaningful characters.
func (buf *DataTableBuffer) ReadAll() (map[string]any, error) {
	if !buf.created {
		return nil, fmt.Errorf("datatable buffer not created")
	}
	if len(buf.tags) == 0 {
		return nil, fmt.Errorf("no tags in datatable buffer")
	}

	resp, err := buf.client.DataTableReadBuffer(buf.instanceID)
	if err != nil {
		return nil, err
	}

	// Sort tags by index so we walk the response in the order the PLC emits entries.
	sorted := make([]DataTableBufferTag, len(buf.tags))
	copy(sorted, buf.tags)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].index < sorted[j].index
	})

	result := make(map[string]any, len(sorted))
	currentBatch := uint16(0)
	for _, tag := range sorted {
		// When we hit a new batch, skip the 2-byte entry header.
		if tag.batchID != currentBatch {
			if _, err := resp.Uint16(); err != nil {
				return nil, fmt.Errorf("could not read batch header before tag %q: %w", tag.Name, err)
			}
			currentBatch = tag.batchID
		}

		// Read exactly DataTypeSize bytes — this is the fixed slot size for
		// this tag's type, matching what was sent in the 0x4E Add Tag request.
		if resp.Pos+int(tag.DataTypeSize) > len(resp.Data) {
			return nil, fmt.Errorf("not enough data for tag %q (index %d): need %d bytes at offset %d, have %d total",
				tag.Name, tag.index, tag.DataTypeSize, resp.Pos, len(resp.Data))
		}
		raw := make([]byte, tag.DataTypeSize)
		copy(raw, resp.Data[resp.Pos:resp.Pos+int(tag.DataTypeSize)])
		resp.Pos += int(tag.DataTypeSize)

		value, err := decodeDataTableValue(tag.CIPType, raw)
		if err != nil {
			return nil, fmt.Errorf("could not decode value for tag %q: %w", tag.Name, err)
		}
		result[tag.Name] = value
	}

	return result, nil
}

// ReadAllRaw reads all tag values from the buffer and returns them as a map
// from tag name to raw byte slices, without type decoding.
//
// The response is parsed identically to ReadAll (see its doc comment for the
// wire layout), but the raw DataTypeSize bytes for each tag are returned
// directly instead of being decoded into Go types.
func (buf *DataTableBuffer) ReadAllRaw() (map[string][]byte, error) {
	if !buf.created {
		return nil, fmt.Errorf("datatable buffer not created")
	}
	if len(buf.tags) == 0 {
		return nil, fmt.Errorf("no tags in datatable buffer")
	}

	resp, err := buf.client.DataTableReadBuffer(buf.instanceID)
	if err != nil {
		return nil, err
	}

	sorted := make([]DataTableBufferTag, len(buf.tags))
	copy(sorted, buf.tags)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].index < sorted[j].index
	})

	result := make(map[string][]byte, len(sorted))
	currentBatch := uint16(0)
	for _, tag := range sorted {
		// When we hit a new batch, skip the 2-byte entry header (see ReadAll).
		if tag.batchID != currentBatch {
			if _, err := resp.Uint16(); err != nil {
				return nil, fmt.Errorf("could not read batch header before tag %q: %w", tag.Name, err)
			}
			currentBatch = tag.batchID
		}

		if resp.Pos+int(tag.DataTypeSize) > len(resp.Data) {
			return nil, fmt.Errorf("not enough data for tag %q (index %d): need %d bytes at offset %d, have %d total",
				tag.Name, tag.index, tag.DataTypeSize, resp.Pos, len(resp.Data))
		}
		raw := make([]byte, tag.DataTypeSize)
		copy(raw, resp.Data[resp.Pos:resp.Pos+int(tag.DataTypeSize)])
		resp.Pos += int(tag.DataTypeSize)
		result[tag.Name] = raw
	}

	return result, nil
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
// TagGroup integration
// ---------------------------------------------------------------------------

// AddTagGroup adds all tags from a TagGroup to this datatable buffer.
// Args substitute {0}, {1}, ... placeholders in tag names.
//
// Multi-element tags (Elements > 1) are expanded into individual buffer tags.
// For example, a TagDef with Name="EXAMPLE[{0},0]", Type=CIPTypeDINT,
// Elements=5 and args=[1] adds 5 individual tags: "EXAMPLE[1,0]" through
// "EXAMPLE[1,4]".
//
// All tags are batched into a single AddTags call for efficiency (the protocol
// allows mixed dataTypeSizes in one 0x4E service request).
//
// Can be called multiple times with different args to build up a large buffer:
//
//	buf.AddTagGroup(inspPointTags, 1)  // inspection point 1
//	buf.AddTagGroup(inspPointTags, 2)  // inspection point 2
//	// buf now has all tags for both points
func (buf *DataTableBuffer) AddTagGroup(group *TagGroup, args ...any) error {
	if !buf.created {
		return fmt.Errorf("datatable buffer not created")
	}

	// Collect all resolved tag names with their types for batched AddTags.
	tagMap := make(map[string]CIPType)

	for _, def := range group.defs {
		resolvedName := formatName(def.Name, args...)

		if def.Elements > 1 {
			expanded := expandArrayTag(resolvedName, def.Elements)
			buf.expansions[resolvedName] = expanded
			for _, name := range expanded {
				tagMap[name] = def.Type
			}
		} else {
			buf.expansions[resolvedName] = []string{resolvedName}
			tagMap[resolvedName] = def.Type
		}
	}

	return buf.AddTags(tagMap)
}

// ReadAllTyped reads all tag values from the buffer in a single CIP request
// and returns a TagGroupResult with typed accessors.
//
// Multi-element tags that were added via AddTagGroup are automatically
// re-collapsed from individual values into slices. For example, if
// "EXAMPLE[1,0]" was expanded into 5 individual tags during AddTagGroup,
// the result will contain "EXAMPLE[1,0]" → []any{val0, val1, ..., val4}.
//
// Tags added via the plain AddTag/AddTags methods (not through a TagGroup)
// are included as-is in the result.
func (buf *DataTableBuffer) ReadAllTyped() (*TagGroupResult, error) {
	rawValues, err := buf.ReadAll()
	if err != nil {
		return nil, err
	}

	// If no TagGroup expansions were recorded, just wrap the raw values.
	if len(buf.expansions) == 0 {
		return &TagGroupResult{values: rawValues}, nil
	}

	// Build the result, collapsing expanded multi-element tags back into slices.
	result := make(map[string]any, len(rawValues))

	// Track which raw keys have been consumed by expansions.
	consumed := make(map[string]bool)

	for baseName, expanded := range buf.expansions {
		if len(expanded) > 1 {
			// Multi-element: collapse into a slice.
			slice := make([]any, len(expanded))
			for i, name := range expanded {
				v, ok := rawValues[name]
				if !ok {
					return nil, fmt.Errorf("expanded tag %q not found in buffer response", name)
				}
				slice[i] = v
				consumed[name] = true
			}
			result[baseName] = slice
		} else if len(expanded) == 1 {
			// Scalar from TagGroup: pass through directly.
			v, ok := rawValues[expanded[0]]
			if !ok {
				return nil, fmt.Errorf("tag %q not found in buffer response", expanded[0])
			}
			result[baseName] = v
			consumed[expanded[0]] = true
		}
	}

	// Include any tags that were added via plain AddTag (not through a TagGroup).
	for name, val := range rawValues {
		if !consumed[name] {
			result[name] = val
		}
	}

	return &TagGroupResult{values: result}, nil
}
