package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

func multi_to_dict(data any) (map[string]interface{}, error) {
	//TODO: handle nested structs and arrays
	// convert the struct to a dict of FieldName: FieldValue
	d := make(map[string]interface{})

	vs := reflect.ValueOf(data)
	fs := reflect.TypeOf(data)
	for i := 0; i < vs.NumField(); i++ {
		v := vs.Field(i)
		f := fs.Field(i)
		fulltag := f.Tag.Get("gologix")
		switch v.Kind() {

		}
		switch v.Kind() {
		case reflect.Struct:
			d2, err := udt_to_dict(fulltag, v.Interface())
			if err != nil {
				return nil, fmt.Errorf("problem parsing %s. %w", fulltag, err)
			}
			for k := range d2 {
				d[k] = d2[k]
			}
		case reflect.Array:
			//TODO
		default:
			d[fulltag] = v.Interface()
		}

	}

	return d, nil

}

func udt_to_dict(tag string, data any) (map[string]interface{}, error) {
	//TODO: handle nested structs and arrays
	// convert the struct to a dict of FieldName: FieldValue
	d := make(map[string]interface{})

	vs := reflect.ValueOf(data)
	fs := reflect.TypeOf(data)
	for i := 0; i < vs.NumField(); i++ {
		v := vs.Field(i)
		f := fs.Field(i)
		fulltag := fmt.Sprintf("%s.%s", tag, f.Name)
		switch v.Kind() {

		}
		switch v.Kind() {
		case reflect.Struct:
			d2, err := udt_to_dict(fulltag, v.Interface())
			if err != nil {
				return nil, fmt.Errorf("problem parsing %s. %w", fulltag, err)
			}
			for k := range d2 {
				d[k] = d2[k]
			}
		case reflect.Array:
			//TODO
		default:
			d[fulltag] = v.Interface()
		}

	}

	return d, nil

}

func (client *Client) writeDict(tag_str map[string]interface{}) error {

	// build the tag list from the structure
	tags := make([]string, 0)
	types := make([]CIPType, 0)
	for k := range tag_str {
		ct := GoVarToCIPType(tag_str[k])
		types = append(types, ct)
		tags = append(tags, k)
	}

	// first generate IOIs for each tag
	qty := len(tags)
	iois := make([]*tagIOI, qty)
	for i, tag := range tags {
		var err error
		iois[i], err = client.newIOI(tag, types[i])
		if err != nil {
			return err
		}
	}

	reqitems := make([]cipItem, 2)
	reqitems[0] = NewItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	ioi_header := msgCIPConnectedMultiServiceReq{
		Sequence:     client.Sequencer(),
		Service:      cipService_MultipleService,
		PathSize:     2,
		Path:         [4]byte{0x20, 0x02, 0x24, 0x01},
		ServiceCount: uint16(qty),
	}

	b := bytes.Buffer{}
	// we now have to build up the jump table for each IOI.
	// and pack all the IOIs together into b
	jump_table := make([]uint16, qty)
	jump_start := 2 + qty*2 // 2 bytes + 2 bytes per jump entry
	for i := 0; i < qty; i++ {
		jump_table[i] = uint16(jump_start + b.Len())
		ioi := iois[i]
		h := msgCIPMultiIOIHeader{
			Service: cipService_Write,
			Size:    byte(len(ioi.Buffer) / 2),
		}
		f := msgCIPWriteIOIFooter{
			DataType: uint16(types[i]),
			Elements: 1,
		}
		//f := msgCIPIOIFooter{
		//Elements: 1,
		//}
		binary.Write(&b, binary.LittleEndian, h)
		b.Write(ioi.Buffer)
		binary.Write(&b, binary.LittleEndian, f)
		binary.Write(&b, binary.LittleEndian, tag_str[tags[i]])
	}

	// right now I'm putting the IOI data into the cip Item, but I suspect it might actually be that the readsequencer is
	// the item's data and the service code actually starts the next portion of the message.  But the item's header length reflects
	// the total data so maybe not.
	reqitems[1] = cipItem{Header: cipItemHeader{ID: cipItem_ConnectedData}}
	reqitems[1].Marshal(ioi_header)
	reqitems[1].Marshal(jump_table)
	reqitems[1].Marshal(b.Bytes())

	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, MarshalItems(reqitems))
	if err != nil {
		return err
	}
	_ = hdr
	_ = data

	return nil
}
