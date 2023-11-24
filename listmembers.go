package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// this is specifically the response for a GetAttrList service on a
// template object with requested attributes of 4,5,2,1
// where
//
//	4 = Size of the template in 32 bit words
//	5 = Size of the data in the template (when sent in a read response)
//	2 = Number of fields/members in the template
//	1 = The handle of the template. not sure what this is for yet
type msgGetTemplateAttrListResponse struct {
	SequenceCount   uint16
	Service         CIPService
	Reserved        byte
	Status          byte
	Status_extended byte
	Count           uint16

	// this is the size of the TEMPLATE data of the structure when read
	SizeWords_ID     uint16
	SizeWords_Status uint16
	SizeWords        uint32

	// this is the size of the DATA in the structure when read.
	SizeBytes_ID     uint16
	SizeBytes_Status uint16
	SizeBytes        uint32

	MemberCount_ID     uint16
	MemberCount_Status uint16
	MemberCount        uint16

	Handle_ID     uint16
	Handle_Status uint16
	Handle        uint16
}

func (client *Client) GetTemplateInstanceAttr(str_instance uint32) (msgGetTemplateAttrListResponse, error) {
	if verbose {
		client.Logger.Printf("list members for %v", str_instance)
	}

	// have to start at 1.
	if str_instance == 0 {
		str_instance = 1
	}

	reqitems := make([]CIPItem, 2)
	//reqitems[0] = cipItem{Header: cipItemHeader{ID: cipItem_Null}}
	reqitems[0] = NewItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	p, err := Serialize(
		CipObject_Template, CIPInstance(str_instance),
		//cipObject_Symbol, cipInstance(start_instance),
	)
	if err != nil {
		return msgGetTemplateAttrListResponse{}, fmt.Errorf("couldn't build path. %w", err)
	}

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: uint16(sequencer()),
		Service:       CIPService_GetAttributeList,
		PathLength:    byte(p.Len() / 2),
	}

	reqitems[1] = NewItem(cipItem_ConnectedData, readmsg)
	reqitems[1].Serialize(p.Bytes())
	number_of_attr_to_receive := 4
	attr_Size_32bitWords := 4
	attr_Size_Bytes := 5
	attr_MemberCount := 2
	attr_symbol_type := 1
	reqitems[1].Serialize([5]uint16{
		uint16(number_of_attr_to_receive),
		uint16(attr_Size_32bitWords),
		uint16(attr_Size_Bytes),
		uint16(attr_MemberCount),
		uint16(attr_symbol_type),
	})
	reqitems[1].Serialize(byte(1))
	reqitems[1].Serialize(byte(0))
	reqitems[1].Serialize(uint16(1))

	itemdata, err := SerializeItems(reqitems)
	if err != nil {
		return msgGetTemplateAttrListResponse{}, fmt.Errorf("problem serializing item data: %w", err)
	}
	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, itemdata)
	if err != nil {
		return msgGetTemplateAttrListResponse{}, err
	}
	_ = hdr
	_ = data
	//data_hdr := ListInstanceHeader{}
	//binary.Read(data, binary.LittleEndian, &data_hdr)

	// first six bytes are zero.
	padding := make([]byte, 6)
	_, err = data.Read(padding)
	if err != nil {
		return msgGetTemplateAttrListResponse{}, fmt.Errorf("problem reading padding. %w", err)
	}

	resp_items, err := ReadItems(data)
	if err != nil {
		return msgGetTemplateAttrListResponse{}, fmt.Errorf("couldn't parse items. %w", err)
	}

	// get ready to read tag info from item 1 data
	data2 := bytes.NewBuffer(resp_items[1].Data)

	result := msgGetTemplateAttrListResponse{}
	err = binary.Read(data2, binary.LittleEndian, &result)
	if err != nil {
		return result, fmt.Errorf("problem reading result. %w", err)
	}
	if verbose {
		client.Logger.Printf("Result: %+v\n\n", result)
	}

	return result, nil
}

type msgMemberInfoHdr struct {
	SequenceCount uint16
	Service       CIPService
	Reserved      byte
	Status        uint16
}
type msgMemberInfo struct {
	Info   uint16
	Type   uint16
	Offset uint32
}

func (m msgMemberInfo) CIPType() CIPType {
	return CIPType(m.Type & 0x00FF)
}

func (client *Client) ListMembers(str_instance uint32) (UDTDescriptor, error) {
	if verbose {
		client.Logger.Printf("list members for %v", str_instance)
	}

	template_info, err := client.GetTemplateInstanceAttr(str_instance)

	if err != nil {
		return UDTDescriptor{}, fmt.Errorf("couldn't get template info. %w", err)
	}

	reqitems := make([]CIPItem, 2)
	//reqitems[0] = cipItem{Header: cipItemHeader{ID: cipItem_Null}}
	reqitems[0] = NewItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	p, err := Serialize(
		CipObject_Template, CIPInstance(str_instance),
		//cipObject_Symbol, cipInstance(start_instance),
	)
	if err != nil {
		return UDTDescriptor{}, fmt.Errorf("couldn't build path. %w", err)
	}

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: uint16(sequencer()),
		Service:       CIPService_Read,
		PathLength:    byte(p.Len() / 2),
	}

	reqitems[1] = NewItem(cipItem_ConnectedData, readmsg)
	reqitems[1].Serialize(p.Bytes())
	start_offset := uint32(0)
	read_length := uint16(template_info.SizeWords*4 - 23)
	reqitems[1].Serialize(start_offset)
	reqitems[1].Serialize(read_length)
	reqitems[1].Serialize(byte(1))
	reqitems[1].Serialize(byte(0))
	reqitems[1].Serialize(uint16(1))

	itemdata, err := SerializeItems(reqitems)
	if err != nil {
		return UDTDescriptor{}, fmt.Errorf("problem serializing item data: %w", err)
	}
	hdr, data, err := client.send_recv_data(cipCommandSendUnitData, itemdata)
	if err != nil {
		return UDTDescriptor{}, err
	}
	_ = hdr
	_ = data
	//data_hdr := ListInstanceHeader{}
	//binary.Read(data, binary.LittleEndian, &data_hdr)

	// first six bytes are zero.
	padding := make([]byte, 6)
	_, err = data.Read(padding)
	if err != nil {
		return UDTDescriptor{}, fmt.Errorf("couldn't read padding. %w", err)
	}

	resp_items, err := ReadItems(data)
	if err != nil {
		return UDTDescriptor{}, fmt.Errorf("couldn't parse items. %w", err)
	}

	// get ready to read tag info from item 1 data
	data2 := bytes.NewBuffer(resp_items[1].Data)

	mihdr := msgMemberInfoHdr{}
	err = binary.Read(data2, binary.LittleEndian, &mihdr)
	if err != nil {
		return UDTDescriptor{}, fmt.Errorf("couldn't read member info header. %w", err)
	}

	memberInfos := make([]msgMemberInfo, template_info.MemberCount)
	err = binary.Read(data2, binary.LittleEndian, &memberInfos)
	if err != nil {
		return UDTDescriptor{}, fmt.Errorf("couldn't read memberinfos. %w", err)
	}
	if verbose {
		client.Logger.Printf("Hdr: %+v\nResult: %+v\n\n", mihdr, memberInfos)
	}

	descriptor := UDTDescriptor{}
	descriptor.Info = template_info
	descriptor.Instance_ID = str_instance
	descriptor.Members = make([]UDTMemberDescriptor, template_info.MemberCount)

	struct_name, err := data2.ReadString(0x3B)
	if err != nil {
		return UDTDescriptor{}, fmt.Errorf("couldn't read struct name. %w", err)
	}
	struct_name = struct_name[:len(struct_name)-1]
	descriptor.Name = struct_name

	_, err = data2.ReadString(0x00)
	if err != nil {
		return UDTDescriptor{}, fmt.Errorf("couldn't read unknown data. %w", err)
	}

	for i := 0; i < int(template_info.MemberCount); i++ {

		fieldname, err := data2.ReadString(0x00)
		if err != nil {
			return UDTDescriptor{}, fmt.Errorf("couldn't read field name. %w", err)
		}
		fieldname = fieldname[:len(fieldname)-1]

		descriptor.Members[i].Name = fieldname
		descriptor.Members[i].Info = memberInfos[i]
	}

	return descriptor, nil
}

// full descriptor of a struct in the controller.
// could be a UDT or a builtin struct like a TON
type UDTDescriptor struct {
	Instance_ID uint32
	Name        string
	Info        msgGetTemplateAttrListResponse
	Members     []UDTMemberDescriptor
}

// This function is experimental and not accurate.  I suspect it is accurate only if the last field in the
// udt is a simple atomic type (int, real, dint, etc...).  Use at your own risk.
func (u UDTDescriptor) Size() int {
	maxsize := uint32(0)
	for i := range u.Members {
		m := u.Members[i]
		end := m.Info.Offset + uint32(m.Info.CIPType().Size())
		if end > maxsize {
			maxsize = end
		}
	}
	return int(maxsize)
}

type UDTMemberDescriptor struct {
	Name string
	Info msgMemberInfo
}
