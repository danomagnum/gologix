package eipcommand

type CIPCommand uint16

const (
	NOP                     CIPCommand = 0x00
	ListServices            CIPCommand = 0x04
	PCCCConnectedExplicit   CIPCommand = 0x0A
	PCCCUnconnectedExplicit CIPCommand = 0x0B
	ListIdentity            CIPCommand = 0x63
	ListInterfaces          CIPCommand = 0x64
	SendUnregistered        CIPCommand = 0x52
	RegisterSession         CIPCommand = 0x65
	UnRegisterSession       CIPCommand = 0x66
	SendRRData              CIPCommand = 0x6F
	SendUnitData            CIPCommand = 0x70
	IndicateStatus          CIPCommand = 0x72
	Cancel                  CIPCommand = 0x73
)
