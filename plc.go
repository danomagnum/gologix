package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"time"
)

type PLC struct {
	IPAddress     string
	ProcessorSlot int
	SocketTimeout time.Duration
	readSequencer uint16
	// Route
	conn Connection
}

func (plc *PLC) Read_Single(tag string) []byte {
	return nil
}

func (plc *PLC) Connect() error {
	return plc.conn.Connect(plc.IPAddress)
}

func Read[T GoLogixTypes](plc *PLC, tag string) (T, error) {
	var t T
	fmt.Printf("reading type %T", t)
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

func (plc *PLC) read_single(tag string, datatype CIPType, elements uint16) (any, error) {
	ioi := BuildIOI(tag, datatype)
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
	reqitems[0] = NewItem(CIPItem_ConnectionAddress, &plc.conn.OTNetworkConnectionID)
	reqitems[1] = CIPItem{Header: CIPItemHeader{ID: CIPItem_ConnectedData}}
	reqitems[1].Marshal(ioi_header)
	reqitems[1].Marshal(ioi.Buffer)
	reqitems[1].Marshal(ioi_footer)

	plc.conn.Send(CIPCommandSendUnitData, BuildItemsBytes(reqitems))
	hdr, data, err := plc.conn.recv_data()
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
