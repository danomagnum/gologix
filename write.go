package gologix

func (plc *PLC) Write_single(tag string, value any) error {
	//service = 0x4D // CIPService_Write
	datatype := GoVarToCIPType(value)
	ioi := NewIOI(tag, datatype)
	plc.readSequencer += 1
	ioi_header := CIPIOIHeader{
		Sequence: plc.readSequencer,
		Service:  CIPService_Write,
		Size:     byte(len(ioi.Buffer) / 2),
	}
	ioi_footer := CIPWriteIOIFooter{
		DataType: uint16(datatype),
		Elements: 1,
	}

	reqitems := make([]CIPItem, 2)
	reqitems[0] = NewItem(CIPItem_ConnectionAddress, &plc.OTNetworkConnectionID)
	reqitems[1] = CIPItem{Header: CIPItemHeader{ID: CIPItem_ConnectedData}}
	reqitems[1].Marshal(ioi_header)
	reqitems[1].Marshal(ioi.Buffer)
	reqitems[1].Marshal(ioi_footer)
	reqitems[1].Marshal(value)

	err := plc.Send(CIPCommandSendUnitData, MarshalItems(reqitems))
	if err != nil {
		return err
	}

	hdr, data, err := plc.recv_data()
	if err != nil {
		return err
	}
	_ = hdr
	_ = data
	return err
}
