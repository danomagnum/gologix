package main

type CIPMessage_Register struct {
	ProtocolVersion uint16
	OptionFlag      uint16
}
type CIPMessage_UnRegister struct {
	Service                CIPService
	CipPathSize            byte
	ClassType              byte
	Class                  byte
	InstanceType           byte
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
