package gologix

import (
	"encoding/binary"
	"fmt"
)

// Represents a CIP service id
//
// If you're going to serialize this class to bytes for transimssion be sure to use one of the gologix
// serialization functions or call Bytes() to get the properly formatted data.
type CIPService byte

func (s CIPService) IsResponse() bool {
	// bit 8 of the service indicates whether it is a response service
	is_response := s & 0b10000000
	return is_response != 0
}
func (s CIPService) AsResponse() CIPService {
	// bit 8 of the service indicates whether it is a response service
	return s | 0b10000000
}
func (s CIPService) UnResponse() CIPService {
	// bit 8 of the service indicates whether it is a response service
	return s & 0b01111111
}
func (s CIPService) String() string {
	switch s {
	case CIPService_GetAttributeAll:
		return "cipService_GetAttributeAll"
	case CIPService_SetAttributeAll:
		return "cipService_SetAttributeAll"
	case CIPService_GetAttributeList:
		return "cipService_GetAttributeList"
	case CIPService_SetAttributeList:
		return "cipService_SetAttributeList"
	case CIPService_Reset:
		return "cipService_Reset"
	case CIPService_Start:
		return "cipService_Start"
	case CIPService_Stop:
		return "cipService_Stop"
	case CIPService_Create:
		return "cipService_Create"
	case CIPService_Delete:
		return "cipService_Delete"
	case CIPService_MultipleService:
		return "cipService_MultipleService"
	case CIPService_ApplyAttributes:
		return "cipService_ApplyAttributes"
	case CIPService_GetAttributeSingle:
		return "cipService_GetAttributeSingle"
	case CIPService_SetAttributeSingle:
		return "cipService_SetAttributeSingle"
	case CIPService_FindNextObjectInstance:
		return "cipService_FindNextObjectInstance"
	case CIPService_Restore:
		return "cipService_Restore"
	case CIPService_Save:
		return "cipService_Save"
	case CIPService_NOP:
		return "cipService_NOP"
	case CIPService_GetMember:
		return "cipService_GetMember"
	case CIPService_SetMember:
		return "cipService_SetMember"
	case CIPService_InsertMember:
		return "cipService_InsertMember"
	case CIPService_RemoveMember:
		return "cipService_RemoveMember"
	case CIPService_GroupSync:
		return "cipService_GroupSync"
	case CIPService_GetMemberList:
		return "cipService_GetMemberList"
	case CIPService_Read:
		return "cipService_Read"
	case CIPService_Write:
		return "cipService_Write"
	case CIPService_ForwardClose:
		return "cipService_ForwardClose"
	case CIPService_GetConnectionOwner:
		return "cipService_GetConnectionOwner"
	case CIPService_ForwardOpen:
		return "cipService_ForwardOpen"
	case CIPService_LargeForwardOpen:
		return "cipService_LargeForwardOpen"
	case CIPService_FragRead:
		return "cipService_FragRead"
	case CIPService_FragWrite:
		return "cipService_FragWrite"
	case CIPService_GetInstanceAttributeList:
		return "cipService_GetInstanceAttributeList"
	case CIPService_GetConnectionData:
		return "cipService_GetConnectionData"
	}
	return fmt.Sprintf("unknown service %d", s)

}

const (
	// cip common services
	CIPService_GetAttributeAll        CIPService = 0x01
	CIPService_SetAttributeAll        CIPService = 0x02
	CIPService_GetAttributeList       CIPService = 0x03
	CIPService_SetAttributeList       CIPService = 0x04
	CIPService_Reset                  CIPService = 0x05
	CIPService_Start                  CIPService = 0x06
	CIPService_Stop                   CIPService = 0x07
	CIPService_Create                 CIPService = 0x08
	CIPService_Delete                 CIPService = 0x09
	CIPService_MultipleService        CIPService = 0x0A
	_                                 CIPService = 0x0B // CIP Reserved
	_                                 CIPService = 0x0C // CIP Reserved
	CIPService_ApplyAttributes        CIPService = 0x0D
	CIPService_GetAttributeSingle     CIPService = 0x0E
	_                                 CIPService = 0x0F // CIP Reserved
	CIPService_SetAttributeSingle     CIPService = 0x10
	CIPService_FindNextObjectInstance CIPService = 0x11
	_                                 CIPService = 0x12 // CIP Reserved
	_                                 CIPService = 0x13 // CIP Reserved
	CIPService_ErrorResponse          CIPService = 0x14
	CIPService_Restore                CIPService = 0x15
	CIPService_Save                   CIPService = 0x16
	CIPService_NOP                    CIPService = 0x17
	CIPService_GetMember              CIPService = 0x18
	CIPService_SetMember              CIPService = 0x19
	CIPService_InsertMember           CIPService = 0x1A
	CIPService_RemoveMember           CIPService = 0x1B
	CIPService_GroupSync              CIPService = 0x1C
	CIPService_GetMemberList          CIPService = 0x1D
	_                                 CIPService = 0x1E // CIP Reserved
	_                                 CIPService = 0x1F
	_                                 CIPService = 0x31 // CIP Reserved
	// cip object services -- do these change per CIP object?
	CIPService_GetInstanceList          CIPService = 0x4B
	CIPService_Read                     CIPService = 0x4C // in OpENer this is called "Eth Link Get And Clear" for some reason
	CIPService_Write                    CIPService = 0x4D
	CIPService_ForwardClose             CIPService = 0x4E
	cipService_ReadModWrite             CIPService = 0x4E
	CIPService_GetConnectionOwner       CIPService = 0x5A
	CIPService_ForwardOpen              CIPService = 0x54
	CIPService_LargeForwardOpen         CIPService = 0x5B
	CIPService_FragRead                 CIPService = 0x52 // Fragmented Read
	CIPService_FragWrite                CIPService = 0x53 // Fragmented Write
	CIPService_GetInstanceAttributeList CIPService = 0x55
	CIPService_GetConnectionData        CIPService = 0x57
)

type CIPCommand uint16

const (
	cipCommand_NOP              CIPCommand = 0x00
	cipCommandListServices      CIPCommand = 0x04
	cipPCCCConnectedExplicit    CIPCommand = 0x0A
	cipPCCCUnconnectedExplicit  CIPCommand = 0x0B
	cipCommandListIdentity      CIPCommand = 0x63
	cipCommandListInterfaces    CIPCommand = 0x64
	cipCommandSendUnregistered  CIPCommand = 0x52
	cipCommandRegisterSession   CIPCommand = 0x65
	cipCommandUnRegisterSession CIPCommand = 0x66
	cipCommandSendRRData        CIPCommand = 0x6F
	cipCommandSendUnitData      CIPCommand = 0x70
	cipCommandIndicateStatus    CIPCommand = 0x72
	cipCommandCancel            CIPCommand = 0x73
)

type CIPTagInfo byte

// currently unused
/*
const (
	cipClass_Identiy        CIPClass = 0x01
	cipClass_AssemblyObject CIPClass = 0x04
)
*/
//

func SizeOf(strs ...any) int {
	t := 0 // total
	for _, str := range strs {
		t += binary.Size(str)
	}
	return t
}
