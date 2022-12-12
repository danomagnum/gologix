package gologix

// todo: move sequence to a different struct and combine msgCIPIOIHeader and CIPMultiIOIHeader
type msgCIPIOIHeader struct {
	Sequence uint16
	Service  CIPService
	Size     byte
}

type msgCIPMultiIOIHeader struct {
	Service CIPService
	Size    byte
}

// this is the generic connected message.
// it goes into an item (always item[1]?) and is followed up with
// a valid path.  The item specifies the CIPService that goes with the message
type msgCIPConnectedServiceReq struct {
	SequenceCount uint16
	Service       CIPService
	PathLength    byte
}

type msgCIPConnectedMultiServiceReq struct {
	Sequence     uint16
	Service      CIPService
	PathSize     byte
	Path         [4]byte
	ServiceCount uint16
}

type msgCIPWriteIOIFooter struct {
	DataType uint16
	Elements uint16
}
type msgCIPIOIFooter struct {
	Elements uint16
}

type msgCIPResultHeader struct {
	InterfaceHandle uint32
	Timeout         uint16
}

// This should be everything before the actual result value data
// so you can read this off the buffer and be in the correct position to
// read the actual value as the type indicated by Type
type msgCIPReadResultData struct {
	SequenceCounter uint16
	Service         CIPService
	Status          [3]byte
	Type            CIPType
	Unknown         byte
}
