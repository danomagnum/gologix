package cipservice

import (
	"encoding/binary"
	"fmt"
)

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
	case GetAttributeAll:
		return "cipService_GetAttributeAll"
	case SetAttributeAll:
		return "cipService_SetAttributeAll"
	case GetAttributeList:
		return "cipService_GetAttributeList"
	case SetAttributeList:
		return "cipService_SetAttributeList"
	case Reset:
		return "cipService_Reset"
	case Start:
		return "cipService_Start"
	case Stop:
		return "cipService_Stop"
	case Create:
		return "cipService_Create"
	case Delete:
		return "cipService_Delete"
	case MultipleService:
		return "cipService_MultipleService"
	case ApplyAttributes:
		return "cipService_ApplyAttributes"
	case GetAttributeSingle:
		return "cipService_GetAttributeSingle"
	case SetAttributeSingle:
		return "cipService_SetAttributeSingle"
	case FindNextObjectInstance:
		return "cipService_FindNextObjectInstance"
	case Restore:
		return "cipService_Restore"
	case Save:
		return "cipService_Save"
	case NOP:
		return "cipService_NOP"
	case GetMember:
		return "cipService_GetMember"
	case SetMember:
		return "cipService_SetMember"
	case InsertMember:
		return "cipService_InsertMember"
	case RemoveMember:
		return "cipService_RemoveMember"
	case GroupSync:
		return "cipService_GroupSync"
	case GetMemberList:
		return "cipService_GetMemberList"
	case Read:
		return "cipService_Read"
	case Write:
		return "cipService_Write"
	case ForwardClose:
		return "cipService_ForwardClose"
	case GetConnectionOwner:
		return "cipService_GetConnectionOwner"
	case ForwardOpen:
		return "cipService_ForwardOpen"
	case LargeForwardOpen:
		return "cipService_LargeForwardOpen"
	case FragRead:
		return "cipService_FragRead"
	case FragWrite:
		return "cipService_FragWrite"
	case GetInstanceAttributeList:
		return "cipService_GetInstanceAttributeList"
	case GetConnectionData:
		return "cipService_GetConnectionData"
	}
	return fmt.Sprintf("unknown service %d", s)

}

const (
	// cip common services
	GetAttributeAll        CIPService = 0x01
	SetAttributeAll        CIPService = 0x02
	GetAttributeList       CIPService = 0x03
	SetAttributeList       CIPService = 0x04
	Reset                  CIPService = 0x05
	Start                  CIPService = 0x06
	Stop                   CIPService = 0x07
	Create                 CIPService = 0x08
	Delete                 CIPService = 0x09
	MultipleService        CIPService = 0x0A
	ApplyAttributes        CIPService = 0x0D
	GetAttributeSingle     CIPService = 0x0E
	SetAttributeSingle     CIPService = 0x10
	FindNextObjectInstance CIPService = 0x11
	Restore                CIPService = 0x15
	Save                   CIPService = 0x16
	NOP                    CIPService = 0x17
	GetMember              CIPService = 0x18
	SetMember              CIPService = 0x19
	InsertMember           CIPService = 0x1A
	RemoveMember           CIPService = 0x1B
	GroupSync              CIPService = 0x1C
	GetMemberList          CIPService = 0x1D
	// cip object services
	Read                     CIPService = 0x4C // in OpENer this is called "Eth Link Get And Clear" for some reason
	Write                    CIPService = 0x4D
	ForwardClose             CIPService = 0x4E // Also seen this called "read modify write" //cipService_ReadModWrite     cipService = 0x4E // Read Modify Write
	GetConnectionOwner       CIPService = 0x5A
	ForwardOpen              CIPService = 0x54
	LargeForwardOpen         CIPService = 0x5B
	FragRead                 CIPService = 0x52 // Fragmented Read
	FragWrite                CIPService = 0x53 // Fragmented Write
	GetInstanceAttributeList CIPService = 0x55
	GetConnectionData        CIPService = 0x57
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
