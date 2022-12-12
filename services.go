package gologix

import (
	"encoding/binary"
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

const (
	// cip common services
	cipService_GetAttributeAll        CIPService = 0x01
	cipService_SetAttributeAll        CIPService = 0x02
	cipService_GetAttributeList       CIPService = 0x03
	cipService_SetAttributeList       CIPService = 0x04
	cipService_Reset                  CIPService = 0x05
	cipService_Start                  CIPService = 0x06
	cipService_Stop                   CIPService = 0x07
	cipService_Create                 CIPService = 0x08
	cipService_Delete                 CIPService = 0x09
	cipService_MultipleService        CIPService = 0x0A
	cipService_ApplyAttributes        CIPService = 0x0D
	cipService_GetAttributeSingle     CIPService = 0x0E
	cipService_SetAttributeSingle     CIPService = 0x10
	cipService_FindNextObjectInstance CIPService = 0x11
	cipService_Restore                CIPService = 0x15
	cipService_Save                   CIPService = 0x16
	cipService_NOP                    CIPService = 0x17
	cipService_GetMember              CIPService = 0x18
	cipService_SetMember              CIPService = 0x19
	cipService_InsertMember           CIPService = 0x1A
	cipService_RemoveMember           CIPService = 0x1B
	cipService_GroupSync              CIPService = 0x1C
	cipService_GetMemberList          CIPService = 0x1D
	// cip object services
	cipService_Read                     CIPService = 0x4C // in OpENer this is called "Eth Link Get And Clear" for some reason
	cipService_Write                    CIPService = 0x4D
	cipService_ForwardClose             CIPService = 0x4E // Also seen this called "read modify write" //cipService_ReadModWrite     cipService = 0x4E // Read Modify Write
	cipService_GetConnectionOwner       CIPService = 0x5A
	cipService_LargeForwardOpen         CIPService = 0x5B
	cipService_FragRead                 CIPService = 0x52 // Fragmented Read
	cipService_FragWrite                CIPService = 0x53 // Fragmented Write
	cipService_GetInstanceAttributeList CIPService = 0x55
	cipService_GetConnectionData        CIPService = 0x57
)

type CIPCommand byte

const (
	cipCommand_NOP              CIPCommand = 0x00
	cipCommandListServices      CIPCommand = 0x04
	cipPCCCConnectedExplicit    CIPCommand = 0x0A
	cipPCCCUnconnectedExplicit  CIPCommand = 0x0B
	cipCommandListIdentity      CIPCommand = 0x63
	cipCommandListInterfaces    CIPCommand = 0x64
	cipCommandRegisterSession   CIPCommand = 0x65
	cipCommandUnRegisterSession CIPCommand = 0x66
	cipCommandSendRRData        CIPCommand = 0x6F
	cipCommandSendUnitData      CIPCommand = 0x70
	cipCommandIndicateStatus    CIPCommand = 0x72
	cipCommandCancel            CIPCommand = 0x73
)

type CIPClass byte

const (
	cipClass_Identiy        CIPClass = 0x01
	cipClass_AssemblyObject CIPClass = 0x04
)

func SizeOf(strs ...any) int {
	t := 0 // total
	for _, str := range strs {
		t += binary.Size(str)
	}
	return t
}
