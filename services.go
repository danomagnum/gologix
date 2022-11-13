package main

type CIPService byte

const (
	CIPService_Read        CIPService = 0x4C
	CIPService_PartialRead CIPService = 0x52
	CIPService_Write       CIPService = 0x4D
	CIPService_ModWrite    CIPService = 0x4E
	CIPService_FragWrite   CIPService = 0x53
)
