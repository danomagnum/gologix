package main

type EIPHeader struct {
	Command       uint16
	Length        uint16
	SessionHandle uint32
	Status        uint32
	Context       uint64 // 8 bytes you can do whatever you want with. They'll be echoed back.
	Options       uint32
}

type CIPMessage_Register struct {
	ProtocolVersion uint16
	OptionFlag      uint16
}
type CIPMessage_UnRegister struct {
	Service                CIPService
	CipPathSize            byte
	ClassType              CIPClassType
	Class                  byte
	InstanceType           CIPInstanceType
	Instance               byte
	Priority               byte
	TimeoutTicks           byte
	ConnectionSerialNumber uint16
	VendorID               uint16
	OriginatorSerialNumber uint32
	PathSize               uint16
	Path                   [6]byte
}

type CIPIOIHeader struct {
	Service CIPService
	Size    byte
}
type CIPIOIFooter struct {
	Elements uint16
	Offset   uint32
}

type CIPReadResultHeader struct {
	InterfaceHandle uint32
	Timeout         uint16
}

// This should be everything before the actual result value data
// so you can read this off the buffer and be in the correct position to
// read the actual value as the type indicated by Type
type CIPReadResultData struct {
	SequenceCounter uint16
	Service         CIPService
	Status          [3]byte
	Type            CIPType
	Unknown         byte
}
