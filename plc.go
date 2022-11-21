package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

type PLC struct {
	IPAddress     string
	ProcessorSlot int
	SocketTimeout time.Duration
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

	ioi_header := CIPIOIHeader{
		Service: CIPService_FragRead,
		Size:    byte(len(ioi.Buffer) / 2),
	}
	ioi_footer := CIPIOIFooter{
		Elements: 1,
		Offset:   0,
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

	read_result_header := CIPReadResultHeader2{}
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

type CIPReadResultHeader2 struct {
	InterfaceHandle uint32
	Timeout         uint16
}
type CIPReadResultData struct {
	SequenceCounter uint16
	Service         CIPService
	Status          [3]byte
	Type            CIPType
	Unknown         byte
}

// readValue reads one unit of cip data type t into the correct go type.
// To do this it reads the needed number of bytes from r.
// It returns the value as an any so the caller will have to do a cast to get it back
func readValue(t CIPType, r io.Reader) any {

	var value any
	var err error
	switch t {
	case CIPTypeUnknown:
		panic("Unknown type.")
	case CIPTypeStruct:
		panic("Struct!")
	case CIPTypeBOOL:
		var trueval bool
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeSINT:
		var trueval byte
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeINT:
		var trueval int16
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeDINT:
		var trueval int32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeLINT:
		var trueval int64
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeUSINT:
		var trueval uint8
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeUINT:
		var trueval uint16
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeUDINT:
		var trueval uint32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeLWORD:
		var trueval uint64
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeREAL:
		var trueval float32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeLREAL:
		var trueval float64
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeDWORD:
		var trueval uint32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeSTRING:
		var trueval [86]byte
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	default:
		panic("Default type.")

	}
	if err != nil {
		log.Printf("Problem reading %s as one unit of %T. %v", t, value, err)
	}
	log.Printf("type %v. value %v", t, value)
	return value
}
