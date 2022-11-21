package main

import "encoding/binary"

type CIPService byte

/*
 * Documentation:
 *
 *  https://literature.rockwellautomation.com/idc/groups/literature/documents/pm/1756-pm020_-en-p.pdf
 *  https://www.odva.org/wp-content/uploads/2020/06/PUB00123R1_Common-Industrial_Protocol_and_Family_of_CIP_Networks.pdf
 *  https://scadahacker.com/library/Documents/ICS_Protocols/Rockwell%20-%20Communicating%20with%20RA%20Products%20Using%20EtherNetIP%20Explicit%20Messaging.pdf
 *  http://iatips.com/digiwiki/quick_eip_demo.pdf
 *  https://github.com/EIPStackGroup/OpENer/blob/master/source/src/cip/ciptypes.h
 *
 */

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
	CIPService_ApplyAttributes        CIPService = 0x0D
	CIPService_GetAttributeSingle     CIPService = 0x0E
	CIPService_SetAttributeSingle     CIPService = 0x10
	CIPService_FindNextObjectInstance CIPService = 0x11
	CIPService_Restore                CIPService = 0x15
	CIPService_Save                   CIPService = 0x16
	CIPService_NOP                    CIPService = 0x17
	CIPService_GetMember              CIPService = 0x18
	CIPService_SetMember              CIPService = 0x19
	CIPService_InsertMember           CIPService = 0x1A
	CIPService_RemoveMember           CIPService = 0x1B
	CIPService_GroupSync              CIPService = 0x1C
	CIPService_GetMemberList          CIPService = 0x1D
	// cip object services
	CIPService_Read               CIPService = 0x4C // in OpENer this is called "Eth Link Get And Clear" for some reason
	CIPService_Write              CIPService = 0x4D
	CIPService_ForwardClose       CIPService = 0x4E // Also seen this called "read modify write" //CIPService_ReadModWrite     CIPService = 0x4E // Read Modify Write
	CIPService_GetConnectionOwner CIPService = 0x5A
	CIPService_LargeForwardOpen   CIPService = 0x5B
	CIPService_FragRead           CIPService = 0x52 // Fragmented Read
	CIPService_FragWrite          CIPService = 0x53 // Fragmented Write
	CIPService_GetConnectionData  CIPService = 0x57
)

type CIPCommand byte

const (
	CIPCommand_NOP              CIPCommand = 0x00
	CIPCommandListServices      CIPCommand = 0x04
	CIPPCCCConnectedExplicit    CIPCommand = 0x0A
	CIPPCCCUnconnectedExplicit  CIPCommand = 0x0B
	CIPCommandListIdentity      CIPCommand = 0x63
	CIPCommandListInterfaces    CIPCommand = 0x64
	CIPCommandRegisterSession   CIPCommand = 0x65
	CIPCommandUnRegisterSession CIPCommand = 0x66
	CIPCommandSendRRData        CIPCommand = 0x6F
	CIPCommandSendUnitData      CIPCommand = 0x70
	CIPCommandIndicateStatus    CIPCommand = 0x72
	CIPCommandCancel            CIPCommand = 0x73
)

type CIPClass byte

const (
	CIPClass_Identiy        CIPClass = 0x01
	CIPClass_AssemblyObject CIPClass = 0x04
)

type CIPAttribute byte

const (
	CIPAttribute_Data CIPAttribute = 0x03
)

func SizeOf(strs ...any) int {
	t := 0 // total
	for _, str := range strs {
		t += binary.Size(str)
	}
	return t
}
