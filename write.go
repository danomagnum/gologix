package gologix

import "fmt"

func (client *Client) Write_single(tag string, value any) error {
	//service = 0x4D // CIPService_Write
	datatype := GoVarToCIPType(value)
	ioi, err := client.newIOI(tag, datatype)
	if err != nil {
		return fmt.Errorf("problem generating IOI. %w", err)
	}
	ioi_header := msgCIPIOIHeader{
		Sequence: client.Sequencer(),
		Service:  CIPService_Write,
		Size:     byte(len(ioi.Buffer) / 2),
	}
	ioi_footer := msgCIPWriteIOIFooter{
		DataType: uint16(datatype),
		Elements: 1,
	}

	reqitems := make([]CIPItem, 2)
	reqitems[0] = NewItem(CIPItem_ConnectionAddress, &client.OTNetworkConnectionID)
	reqitems[1] = CIPItem{Header: CIPItemHeader{ID: CIPItem_ConnectedData}}
	reqitems[1].Marshal(ioi_header)
	reqitems[1].Marshal(ioi.Buffer)
	reqitems[1].Marshal(ioi_footer)
	reqitems[1].Marshal(value)

	hdr, data, err := client.send_recv_data(CIPCommandSendUnitData, MarshalItems(reqitems))
	if err != nil {
		return err
	}
	_ = hdr
	_ = data
	return err
}
