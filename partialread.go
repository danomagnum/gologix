package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// ReadList reads multiple tags with explicitly specified types and element counts in a single efficient request.
//
// This function allows reading tags when you know the tag names but want to specify types and quantities
// explicitly rather than relying on Go variable types. It's useful when working with dynamic tag lists
// or when you need to read partial arrays with specific element counts.
//
// Parameters:
//   - tagnames: Slice of tag names to read (case insensitive)
//   - types: Slice of example values with correct Go types for each tag
//   - elements: Slice of element counts to read for each tag (1 for scalars)
//
// The three slices must have the same length and correspond by index.
//
// Returns a slice of interface{} values in the same order as the input tags.
// You'll need to type assert the returned values to their expected types.
//
// Examples:
//   // Read scalar tags
//   tagnames := []string{"TestInt", "TestReal", "TestBool"}
//   types := []any{int16(0), float32(0), false}
//   elements := []int{1, 1, 1}
//   values, err := client.ReadList(tagnames, types, elements)
//   if err != nil {
//       log.Fatal(err)
//   }
//   intVal := values[0].(int16)
//   realVal := values[1].(float32)
//   boolVal := values[2].(bool)
//
//   // Read partial arrays with different element counts
//   tagnames := []string{"DintArray[0]", "RealArray[5]", "StringArray[0]"}
//   types := []any{[]int32{}, []float32{}, []string{}}
//   elements := []int{10, 5, 3}  // Read 10 DINTs, 5 REALs, 3 STRINGs
//   values, err := client.ReadList(tagnames, types, elements)
//
// For strongly-typed reading with structs, use ReadMulti instead.
// For map-based reading, use ReadMap.
// For single tags, use Read or Read_single.
//
// ReadList automatically handles message splitting for large requests to stay within connection limits.
func (client *Client) ReadList(tagnames []string, types []any, elements []int) ([]any, error) {
	err := client.checkConnection()
	if err != nil {
		return nil, fmt.Errorf("could not start list read: %w", err)
	}
	n := 0
	n_new := 0
	total := len(tagnames)
	results := make([]any, 0, total)
	msgs := 0

	tags := make([]tagDesc, total)

	for i := range tagnames {
		typ, _ := GoVarToCIPType(types[i])
		tags[i] = tagDesc{
			TagName:  tagnames[i],
			TagType:  typ,
			Elements: elements[i],
			Struct:   types[i],
		}
	}

	for n < total {
		msgs += 1
		n_new, err = client.countIOIsThatFit(tags[n:])
		if err != nil {
			return nil, err
		}
		subresults, err := client.readList(tags[n : n+n_new])
		n += n_new
		if err != nil {
			return nil, err
		}
		results = append(results, subresults...)

	}

	client.Logger.Debug("Multi Read", "messages", msgs, "tags", n)
	return results, nil
}

func (client *Client) countIOIsThatFit(tags []tagDesc) (int, error) {
	// first generate IOIs for each tag
	qty := len(tags)

	ioi_header := msgCIPConnectedMultiServiceReq{
		Sequence:     uint16(sequencer()),
		Service:      CIPService_MultipleService,
		PathSize:     2,
		Path:         [4]byte{0x20, 0x02, 0x24, 0x01},
		ServiceCount: uint16(qty),
	}

	mainhdr_size := binary.Size(ioi_header)
	ioihdr_size := binary.Size(msgCIPMultiIOIHeader{})
	ioiftr_size := binary.Size(msgCIPIOIFooter{})

	b := bytes.Buffer{}
	// we now have to build up the jump table for each IOI.
	// and pack all the IOIs together into b
	jump_table := make([]uint16, qty)

	// how many ioi's fit in the message
	n := 1

	response_size := 0

	for i, tag := range tags {
		ioi, err := client.newIOI(tag.TagName, tag.TagType)
		if err != nil {
			return 0, err
		}

		jump_table[i] = uint16(b.Len())

		h := msgCIPMultiIOIHeader{
			Service: CIPService_Read,
			Size:    byte(len(ioi.Buffer) / 2),
		}
		f := msgCIPIOIFooter{
			Elements: 1,
		}

		// Calculate the size of the data once we add this ioi to the list.
		newSize := mainhdr_size                                // length of the multi-read header
		newSize += 2 * n                                       // add in the jump table
		newSize += b.Len()                                     // everything we have so far
		newSize += ioihdr_size + len(ioi.Buffer) + ioiftr_size // the new ioi data
		//newSize += 4                                           // Fudge for alignment if needed

		response_size += tags[i].TagType.Size() * tags[i].Elements
		if newSize >= int(client.ConnectionSize) || response_size >= int(client.ConnectionSize) {
			// break before adding this ioi to the list since it will push us over.
			// we'll continue with n iois (n only increments after an IOI is added)
			break
		}

		err = binary.Write(&b, binary.LittleEndian, h)
		if err != nil {
			return 0, fmt.Errorf("problem writing cip IO header to buffer. %w", err)
		}
		b.Write(ioi.Buffer)
		err = binary.Write(&b, binary.LittleEndian, f)
		if err != nil {
			return 0, fmt.Errorf("problem writing ioi buffer to msg buffer. %w", err)
		}

		n = i + 1
	}

	client.Logger.Debug("Packed Efficiency", "tags", n, "bytes", client.ConnectionSize)

	return n, nil

}
