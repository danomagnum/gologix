package main

type CIPService byte

/*
 * Documentation:
 *
 *  https://literature.rockwellautomation.com/idc/groups/literature/documents/pm/1756-pm020_-en-p.pdf
 *  https://www.odva.org/wp-content/uploads/2020/06/PUB00123R1_Common-Industrial_Protocol_and_Family_of_CIP_Networks.pdf
 *  https://scadahacker.com/library/Documents/ICS_Protocols/Rockwell%20-%20Communicating%20with%20RA%20Products%20Using%20EtherNetIP%20Explicit%20Messaging.pdf
 *  http://iatips.com/digiwiki/quick_eip_demo.pdf
 *
 */

const (
	CIPService_Read            CIPService = 0x4C
	CIPService_FragRead        CIPService = 0x52 // Fragmented Read
	CIPService_Write           CIPService = 0x4D
	CIPService_ReadModWrite    CIPService = 0x4E // Read Modify Write
	CIPService_FragWrite       CIPService = 0x53 // Fragmented Write
	CIPService_MultipleService CIPService = 0x0A
)

type CIPCommand byte

const (
	CIPCommand_NOP              CIPCommand = 0x00
	CIPCommandListServices      CIPCommand = 0x04
	CIPCommandListIdentity      CIPCommand = 0x63
	CIPCommandListInterfaces    CIPCommand = 0x64
	CIPCommandRegisterSession   CIPCommand = 0x65
	CIPCommandUnRegisterSession CIPCommand = 0x66
	CIPCommandSendRRData        CIPCommand = 0x6F
	CIPCommandSendUnitData      CIPCommand = 0x70
	CIPCommandIndicateStatus    CIPCommand = 0x72
	CIPCommandCancel            CIPCommand = 0x73
	CIPPCCCConnectedExplicit    CIPCommand = 0x0A
	CIPPCCCUnconnectedExplicit  CIPCommand = 0x0B
)
