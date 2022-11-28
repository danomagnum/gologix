package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
)

func (plc *PLC) Read_Single(tag string) []byte {
	return nil
}

func (plc *PLC) read_single(tag string, datatype CIPType, elements uint16) (any, error) {
	ioi := BuildIOI(tag, datatype)
	// you have to change this read sequencer every time you make a new tag request.  If you don't, you
	// won't get an error but it will return the last value you requested again.
	// You don't have to keep incrementing it.  just going back and forth between 1 and 0 works OK.
	plc.readSequencer += 1

	ioi_header := CIPIOIHeader{
		Sequence: plc.readSequencer,
		Service:  CIPService_Read,
		Size:     byte(len(ioi.Buffer) / 2),
	}
	ioi_footer := CIPIOIFooter{
		Elements: 1,
	}

	reqitems := make([]CIPItem, 2)
	reqitems[0] = NewItem(CIPItem_ConnectionAddress, &plc.OTNetworkConnectionID)

	// right now I'm putting the IOI data into the cip Item, but I suspect it might actually be that the readsequencer is
	// the item's data and the service code actually starts the next portion of the message.  But the item's header length reflects
	// the total data so maybe not.
	reqitems[1] = CIPItem{Header: CIPItemHeader{ID: CIPItem_ConnectedData}}
	reqitems[1].Marshal(ioi_header)
	reqitems[1].Marshal(ioi.Buffer)
	reqitems[1].Marshal(ioi_footer)

	plc.Send(CIPCommandSendUnitData, BuildItemsBytes(reqitems))
	hdr, data, err := plc.recv_data()
	if err != nil {
		return nil, err
	}
	_ = hdr

	read_result_header := CIPReadResultHeader{}
	err = binary.Read(data, binary.LittleEndian, &read_result_header)
	if err != nil {
		log.Printf("Problem reading read result header. %v", err)
	}
	items, err := ReadItems(data)
	if err != nil {
		log.Printf("Problem reading items. %v", err)
		return 0, nil
	}
	if len(items) != 2 {
		return 0, fmt.Errorf("wrong Number of Items. Expected 2 but got %v", len(items))
	}
	var hdr2 CIPReadResultData
	err = items[1].Unmarshal(&hdr2)
	if err != nil {
		return 0, fmt.Errorf("problem reading item 2's header. %w", err)
	}

	value := readValue(hdr2.Type, &items[1])
	_ = value
	return value, nil
}

func Read[T GoLogixTypes](plc *PLC, tag string) (T, error) {
	var t T
	//fmt.Printf("reading type %T", t)
	ct := GoVarToCIPType(t)
	val, err := plc.read_single(tag, ct, 1)
	if err != nil {
		return t, err
	}
	cast, ok := val.(T)
	if !ok {
		return t, errors.New("couldn't convert to correct type")
	}
	return cast, nil

}

func (plc *PLC) read_multi(tags []string, datatype CIPType, elements uint16) (any, error) {
	// first generate IOIs for each tag
	qty := len(tags)
	iois := make([]*IOI, qty)
	for i, tag := range tags {
		iois[i] = BuildIOI(tag, datatype)
	}
	// you have to change this read sequencer every time you make a new tag request.  If you don't, you
	// won't get an error but it will return the last value you requested again.
	// You don't have to keep incrementing it.  just going back and forth between 1 and 0 works OK.
	plc.readSequencer += 1

	reqitems := make([]CIPItem, 2)
	reqitems[0] = NewItem(CIPItem_ConnectionAddress, &plc.OTNetworkConnectionID)

	ioi_header := CIPMultiServiceHeader{
		Sequence:     plc.readSequencer,
		Service:      CIPService_MultipleService,
		PathSize:     2,
		Path:         [4]byte{0x20, 0x02, 0x24, 0x01},
		ServiceCount: uint16(qty),
	}

	b := bytes.Buffer{}
	// we now have to build up the jump table for each IOI.
	jump_table := make([]uint16, qty)
	jump_start := 2 + qty*2 // 2 bytes + 2 bytes per jump entry
	for i := 0; i < qty; i++ {
		jump_table[i] = uint16(jump_start + b.Len())
		ioi := iois[i]
		h := CIPMultiIOIHeader{
			Service: CIPService_Read,
			Size:    byte(len(ioi.Buffer) / 2),
		}
		f := CIPIOIFooter{
			Elements: elements,
		}
		binary.Write(&b, binary.LittleEndian, h)
		b.Write(ioi.Buffer)
		binary.Write(&b, binary.LittleEndian, f)
	}

	// right now I'm putting the IOI data into the cip Item, but I suspect it might actually be that the readsequencer is
	// the item's data and the service code actually starts the next portion of the message.  But the item's header length reflects
	// the total data so maybe not.
	reqitems[1] = CIPItem{Header: CIPItemHeader{ID: CIPItem_ConnectedData}}
	reqitems[1].Marshal(ioi_header)
	reqitems[1].Marshal(jump_table)
	reqitems[1].Marshal(b.Bytes())

	plc.Send(CIPCommandSendUnitData, BuildItemsBytes(reqitems))
	hdr, data, err := plc.recv_data()
	return nil, nil
	if err != nil {
		return nil, err
	}
	_ = hdr

	read_result_header := CIPReadResultHeader{}
	err = binary.Read(data, binary.LittleEndian, &read_result_header)
	if err != nil {
		log.Printf("Problem reading read result header. %v", err)
	}
	items, err := ReadItems(data)
	if err != nil {
		log.Printf("Problem reading items. %v", err)
		return 0, nil
	}
	if len(items) != 2 {
		return 0, fmt.Errorf("wrong Number of Items. Expected 2 but got %v", len(items))
	}
	var hdr2 CIPReadResultData
	err = items[1].Unmarshal(&hdr2)
	if err != nil {
		return 0, fmt.Errorf("problem reading item 2's header. %w", err)
	}

	value := readValue(hdr2.Type, &items[1])
	_ = value
	return value, nil
}
