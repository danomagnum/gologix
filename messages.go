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

type CIPMessage_IOI struct {
	Header CIPCommonPacketConnected
	IOI    IOI
}

type CIPMessage_ReadModifyWriteHeader struct {
	ServiceCode CIPService
	IOIString   CIPMessage_IOI
	MaskSize    byte // mask size in bytes.  Msut be 1,2,4,8, or 12
}

type CIPCommonPacketUnconnected struct {
	Count            uint16 // always 2
	ItemType_Null    uint16 // always 0
	ItemDataLen_Null uint16 // always 0
	ItemType         uint16 // 0xB2
	ItemDataLen      uint16
	// variable ItemData follows
}

type CIPCommonPacketConnected struct {
	InterfaceHandle uint32
	Timeout         uint16
	ItemCount       uint16
	Item1ID         uint16
	Item1Length     uint16
	Item1           uint32
	Item2ID         uint16
	Item2Length     uint16
	Sequence        uint16
}

type CIPIOIHeader struct {
	Service CIPService
	Size    byte
}
type CIPIOIFooter struct {
	Elements uint16
	Offset   uint32
}
