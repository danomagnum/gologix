package gologix

import (
	"encoding/binary"
	"fmt"
	"io"
)

type cipAddress byte

func (p cipAddress) Bytes() []byte {
	return []byte{byte(p)}
}

func (p cipAddress) Len() int {
	return 0
}

// Represents a CIP class attribute id - that is a specific attribute of a given class.
//
// If you're going to serialize this class to bytes for transimssion be sure to use one of the gologix
// serialization functions or call Bytes() to get the properly formatted data.
type CIPAttribute uint16

// Used to indicate how many bytes are used for the data. If they are more than 8 bits,
// they actually actually take n+1 bytes.  First byte after specifier is a 0
type cipAttributeType byte

const (
	cipAttribute_8bit  cipAttributeType = 0x30
	cipAttribute_16bit cipAttributeType = 0x31
)

func (p CIPAttribute) Bytes() []byte {
	if p < 256 {
		b := make([]byte, 2)
		b[0] = byte(cipAttribute_8bit)
		b[1] = byte(p)
		return b
	} else {

		b := make([]byte, 4)
		b[0] = byte(cipAttribute_16bit)
		binary.LittleEndian.PutUint16(b[2:], uint16(p))
		return b
	}
}
func (p *CIPAttribute) Read(r io.Reader) error {
	var size cipAttributeType
	binary.Read(r, binary.LittleEndian, &size)
	switch size {
	case cipAttribute_8bit:
		var val byte
		binary.Read(r, binary.LittleEndian, &val)
		*p = CIPAttribute(val)
		return nil
	case cipAttribute_16bit:
		binary.Read(r, binary.LittleEndian, p)
		return nil
	default:
		return fmt.Errorf("expected 0x30 or 0x31 but got class size of %x", size)
	}
}
func (p CIPAttribute) Len() int {
	if p < 256 {
		return 2
	}
	return 4
}

// Here are the objects

// used to indicate the array index of data out of an array.
//
// If you're going to serialize this class to bytes for transimssion be sure to use one of the gologix
// serialization functions or call Bytes() to get the properly formatted data.
type CIPElement uint32

type cipElementType byte

const (
	cipElement_8bit  cipElementType = 0x28
	cipElement_16bit cipElementType = 0x29
	cipElement_32bit cipElementType = 0x2A
)

func (p CIPElement) Bytes() []byte {
	if p < 256 {
		b := make([]byte, 2)
		b[0] = byte(cipElement_8bit)
		b[1] = byte(p)
		return b
	} else if p < 65536 {

		b := make([]byte, 4)
		b[0] = byte(cipElement_16bit)
		binary.LittleEndian.PutUint16(b[2:], uint16(p))
		return b
	} else {

		b := make([]byte, 6)
		b[0] = byte(cipElement_16bit)
		binary.LittleEndian.PutUint32(b[2:], uint32(p))
		return b
	}
}

func (p CIPElement) Len() int {
	if p < 256 {
		return 2
	} else if p < 65535 {
		return 4
	}
	return 6
}

type cipInstanceSize byte

const (
	cipInstance_8bit  cipInstanceSize = 0x24
	cipInstance_16bit cipInstanceSize = 0x25
	cipInstance_32bit cipInstanceSize = 0x26
)

// Represents a CIP class instance id - that is a specific instance of a given class.
//
// If you're going to serialize this class to bytes for transimssion be sure to use one of the gologix
// serialization functions or call Bytes() to get the properly formatted data.
type CIPInstance uint32

func (p CIPInstance) Bytes() []byte {
	if p < 256 {
		b := make([]byte, 2)
		b[0] = byte(cipInstance_8bit)
		b[1] = byte(p)
		return b
	} else if p <= 0xFFFF {

		b := make([]byte, 4)
		b[0] = byte(cipInstance_16bit)
		binary.LittleEndian.PutUint16(b[2:], uint16(p))
		return b
	}
	b := make([]byte, 6)
	b[0] = byte(cipInstance_32bit)
	binary.LittleEndian.PutUint32(b[2:], uint32(p))
	return b
}

func (p *CIPInstance) Read(r io.Reader) error {
	var size cipInstanceSize
	binary.Read(r, binary.LittleEndian, &size)
	switch size {
	case cipInstance_8bit:
		var val byte
		binary.Read(r, binary.LittleEndian, &val)
		*p = CIPInstance(val)
		return nil
	case cipInstance_16bit:
		var val uint16
		binary.Read(r, binary.LittleEndian, &val)
		*p = CIPInstance(val)
		return nil
	case cipInstance_32bit:
		binary.Read(r, binary.LittleEndian, p)
		return nil
	default:
		return fmt.Errorf("expected 0x24 or 0x25 but got class size of %x", size)
	}
}
func (p CIPInstance) Len() int {
	if p < 256 {
		return 2
	} else if p <= 0xFFFF {
		return 4
	}
	return 6
}

// Represents a CIP class / object type.
//
// All cip class types are numbered.  Some predefined well-known classes are availabe as constants
// with the prefix of CIPObject
//
// If you're going to serialize this class to bytes for transimssion be sure to use one of the gologix
// serialization functions or call Bytes() to get the properly formatted data.
type CIPClass uint16

// Used to indicate how many bytes are used for the data. If they are more than 8 bits,
// they actually actually take n+1 bytes.  First byte after specifier is a 0
type cipClassSize byte

const (
	cipClass_8bit  cipClassSize = 0x20
	cipClass_16bit cipClassSize = 0x21
)

func (p CIPClass) Bytes() []byte {
	if p < 256 {
		b := make([]byte, 2)
		b[0] = byte(cipClass_8bit)
		b[1] = byte(p)
		return b
	} else {

		b := make([]byte, 4)
		b[0] = byte(cipClass_16bit)
		binary.LittleEndian.PutUint16(b[2:], uint16(p))
		return b
	}
}

func (p *CIPClass) Read(r io.Reader) error {
	var classSize cipClassSize
	binary.Read(r, binary.LittleEndian, &classSize)
	switch classSize {
	case cipClass_8bit:
		var val byte
		binary.Read(r, binary.LittleEndian, &val)
		*p = CIPClass(val)
		return nil
	case cipClass_16bit:
		binary.Read(r, binary.LittleEndian, p)
		return nil
	default:
		return fmt.Errorf("expected 0x20 or 0x21 but got class size of %x", classSize)
	}
}

func (p CIPClass) Len() int {
	if p < 256 {
		return 2
	}
	return 4
}

const (
	CipObject_Identity                     CIPClass = 0x01
	CipObject_MessageRouter                CIPClass = 0x02
	CipObject_DeviceNet                    CIPClass = 0x03
	CipObject_Assembly                     CIPClass = 0x04
	CipObject_Connection                   CIPClass = 0x05
	CipObject_ConnectionManager            CIPClass = 0x06
	CipObject_Register                     CIPClass = 0x07
	CipObject_DiscreteInputPoint           CIPClass = 0x08
	CipObject_DiscreteOutputPoint          CIPClass = 0x09
	CipObject_AnalogInputPoint             CIPClass = 0x0A
	CipObject_AnalogOutputPoint            CIPClass = 0x0B
	CipObject_PresenceSensing              CIPClass = 0x0E
	CipObject_Parameter                    CIPClass = 0x0F
	CipObject_ParameterGroup               CIPClass = 0x10
	CipObject_Group                        CIPClass = 0x12
	CipObject_DiscreteInputGroup           CIPClass = 0x1D
	CipObject_DiscreteOutputGroup          CIPClass = 0x1E
	CipObject_DiscreteGroup                CIPClass = 0x1F
	CipObject_AnalogInputGroup             CIPClass = 0x20
	CipObject_AnalogOutputGroup            CIPClass = 0x21
	CipObject_AnalogGroup                  CIPClass = 0x22
	CipObject_PositionSensor               CIPClass = 0x23
	CipObject_PositionControlSupervisor    CIPClass = 0x24
	CipObject_PositionController           CIPClass = 0x25
	CipObject_BlockSequencer               CIPClass = 0x26
	CipObject_CommandBlock                 CIPClass = 0x27
	CipObject_MotorData                    CIPClass = 0x28
	CipObject_ControlSupervisor            CIPClass = 0x29
	CipObject_Drive                        CIPClass = 0x2A
	CipObject_AckHandler                   CIPClass = 0x2B
	CipObject_Overload                     CIPClass = 0x2C
	CipObject_SoftStart                    CIPClass = 0x2D
	CipObject_Selection                    CIPClass = 0x2E
	CipObject_SDeviceSupervisor            CIPClass = 0x30
	CipObject_SAnalogSensor                CIPClass = 0x31
	CipObject_SAnalogActuator              CIPClass = 0x32
	CipObject_SSingleStageController       CIPClass = 0x33
	CipObject_SGasCalibration              CIPClass = 0x34
	CipObject_TripPoint                    CIPClass = 0x35
	CipObject_File                         CIPClass = 0x37
	CipObject_Symbol                       CIPClass = 0x6B
	CipObject_Template                     CIPClass = 0x6C
	CipObject_ConnectionConfig             CIPClass = 0xF3
	CipObject_OriginatorConnList           CIPClass = 0x45
	CipObject_Port                         CIPClass = 0xF4
	CipObject_BaseEnergy                   CIPClass = 0x4E
	CipObject_ElectricalEnergy             CIPClass = 0x4F
	CipObject_EventLog                     CIPClass = 0x41
	CipObject_MotionAxis                   CIPClass = 0x42
	CipObject_NonElectricalEnergy          CIPClass = 0x50
	CipObject_PowerCurtailment             CIPClass = 0x5C
	CipObject_PowerManagement              CIPClass = 0x53
	CipObject_SPartialPressure             CIPClass = 0x38
	CipObject_SSensorCalibration           CIPClass = 0x40
	CipObject_SafetyAnalogInputGroup       CIPClass = 0x4A
	CipObject_SafetyAnalogInputPoint       CIPClass = 0x49
	CipObject_SafetyDualChannelFeedback    CIPClass = 0x59
	CipObject_SafetyFeedback               CIPClass = 0x5A
	CipObject_SafetyDiscreteInputGroup     CIPClass = 0x3E
	CipObject_SafetyDiscreteInputPoint     CIPClass = 0xeD
	CipObject_SafetyDiscreteOutputGroup    CIPClass = 0x3C
	CipObject_SafetyDiscreteOutputPoint    CIPClass = 0x3B
	CipObject_SafetyDualChannelAnalogInput CIPClass = 0x4B
	CipObject_SafetyDualChannelOutput      CIPClass = 0x3F
	CipObject_SafetyLimitFunctions         CIPClass = 0x5B
	CipObject_SafetyStopFunctions          CIPClass = 0x5A
	CipObject_SafetySupervisor             CIPClass = 0x39
	CipObject_SafetyValidator              CIPClass = 0x3A
	CipObject_TargetConnectionList         CIPClass = 0x4D
	CipObject_TimeSync                     CIPClass = 0x43
	CipObject_BaseSwitch                   CIPClass = 0x51
	CipObject_CompoNetLink                 CIPClass = 0xF7
	CipObject_CompoNetRepeater             CIPClass = 0xF8
	CipObject_ControlNet                   CIPClass = 0xF0
	CipObject_ControlNetKeeper             CIPClass = 0xF1
	CipObject_ControlNetScheduling         CIPClass = 0xF2
	CipObject_DLR                          CIPClass = 0x47
	CipObject_EthernetLink                 CIPClass = 0xF6
	CipObject_Modbus                       CIPClass = 0x44
	CipObject_ModbusSerial                 CIPClass = 0x46
	CipObject_ParallelRedundancyProtocol   CIPClass = 0x56
	CipObject_PRPNodesTable                CIPClass = 0x57
	CipObject_SERCOSIIILink                CIPClass = 0x4C
	CipObject_SNMP                         CIPClass = 0x52
	CipObject_QoS                          CIPClass = 0x48
	CipObject_RSTPBridge                   CIPClass = 0x54
	CipObject_RSTPPort                     CIPClass = 0x55
	CipObject_TCPIP                        CIPClass = 0xF5
	CipObject_PCCC                         CIPClass = 0x67
	CipObject_Programs                     CIPClass = 0x68
	CipObject_TIME                         CIPClass = 0x8B
	CipObject_ControllerInfo               CIPClass = 0xAC // don't know the official name
	CipObject_RunMode                      CIPClass = 0x8E
)

// from https://rockwellautomation.custhelp.com/ci/okcsFattach/get/114390_5

type CIPStatus byte

const (
	CIPStatus_OK                             CIPStatus = 0x00
	CIPStatus_ConnectionFailure              CIPStatus = 0x01
	CIPStatus_ResourceUnavailable            CIPStatus = 0x02
	CIPStatus_InvalidParameterValue          CIPStatus = 0x03
	CIPStatus_PathSegmentError               CIPStatus = 0x04
	CIPStatus_PathDestinationUnknown         CIPStatus = 0x05
	CIPStatus_PartialTransfer                CIPStatus = 0x06
	CIPStatus_ConnectionLost                 CIPStatus = 0x07
	CIPStatus_ServiceNotSupported            CIPStatus = 0x08
	CIPStatus_InvalidAttributeValue          CIPStatus = 0x09
	CIPStatus_AttributeListError             CIPStatus = 0x0A
	CIPStatus_AlreadyInRequestedMode         CIPStatus = 0x0B
	CIPStatus_ObjectStateConflict            CIPStatus = 0x0C
	CIPStatus_ObjectAlreadyExists            CIPStatus = 0x0D
	CIPStatus_AttributeNotSettable           CIPStatus = 0x0E
	CIPStatus_PrivilegeViolation             CIPStatus = 0x0F
	CIPStatus_DeviceStateConflict            CIPStatus = 0x10
	CIPStatus_ReplyDataTooLarge              CIPStatus = 0x11
	CIPStatus_FragmentationOfMessage         CIPStatus = 0x12
	CIPStatus_NotEnoughData                  CIPStatus = 0x13
	CIPStatus_AttributeNotSupported          CIPStatus = 0x14
	CIPStatus_TooMuchData                    CIPStatus = 0x15
	CIPStatus_ObjectDoesNotExist             CIPStatus = 0x16
	CIPStatus_ServiceFragmentation           CIPStatus = 0x17
	CIPStatus_NoStoredAttributeData          CIPStatus = 0x18
	CIPStatus_StoreOperationFailure          CIPStatus = 0x19
	CIPStatus_RoutingFailureReqTooLarge      CIPStatus = 0x1A
	CIPStatus_RoutingFailureRespTooLarge     CIPStatus = 0x1B
	CIPStatus_MissingAttributeListEntry      CIPStatus = 0x1C
	CIPStatus_InvalidAttributeValueList      CIPStatus = 0x1D
	CIPStatus_EmbeddedServiceError           CIPStatus = 0x1E
	CIPStatus_VendorSpecificError            CIPStatus = 0x1F
	CIPStatus_InvalidParameter               CIPStatus = 0x20
	CIPStatus_WriteOnceValueOrMedium         CIPStatus = 0x21
	CIPStatus_InvalidReplyReceived           CIPStatus = 0x22
	CIPStatus_BufferOverflow                 CIPStatus = 0x23
	CIPStatus_MessageFormatError             CIPStatus = 0x24
	CIPStatus_KeyFailure                     CIPStatus = 0x25
	CIPStatus_PathSizeInvalid                CIPStatus = 0x26
	CIPStatus_UnexpectedAttribInList         CIPStatus = 0x27
	CIPStatus_InvalidMemberID                CIPStatus = 0x28
	CIPStatus_MemberNotSettable              CIPStatus = 0x29
	CIPStatus_Group2OnlyServerGeneralFailure CIPStatus = 0x2A
	CIPStatus_UnknownModbusError             CIPStatus = 0x2B
	CIPStatus_AttributeNotGettable           CIPStatus = 0x2C
)

func (s CIPStatus) String() string {
	switch s {
	case CIPStatus_OK:
		return "OK"
	case CIPStatus_ConnectionFailure:
		return "ConnectionFailure"
	case CIPStatus_ResourceUnavailable:
		return "ResourceUnavailable"
	case CIPStatus_InvalidParameterValue:
		return "InvalidParameterValue"
	case CIPStatus_PathSegmentError:
		return "PathSegmentError"
	case CIPStatus_PathDestinationUnknown:
		return "PathDestinationUnknown"
	case CIPStatus_PartialTransfer:
		return "PartialTransfer"
	case CIPStatus_ConnectionLost:
		return "ConnectionLost"
	case CIPStatus_ServiceNotSupported:
		return "ServiceNotSupported"
	case CIPStatus_InvalidAttributeValue:
		return "InvalidAttributeValue"
	case CIPStatus_AttributeListError:
		return "AttributeListError"
	case CIPStatus_AlreadyInRequestedMode:
		return "AlreadyInRequestedMode"
	case CIPStatus_ObjectStateConflict:
		return "ObjectStateConflict"
	case CIPStatus_ObjectAlreadyExists:
		return "ObjectAlreadyExists"
	case CIPStatus_AttributeNotSettable:
		return "AttributeNotSettable"
	case CIPStatus_PrivilegeViolation:
		return "PrivilegeViolation"
	case CIPStatus_DeviceStateConflict:
		return "DeviceStateConflict"
	case CIPStatus_ReplyDataTooLarge:
		return "ReplyDataTooLarge"
	case CIPStatus_FragmentationOfMessage:
		return "FragmentationOfMessage"
	case CIPStatus_NotEnoughData:
		return "NotEnoughData"
	case CIPStatus_AttributeNotSupported:
		return "AttributeNotSupported"
	case CIPStatus_TooMuchData:
		return "TooMuchData"
	case CIPStatus_ObjectDoesNotExist:
		return "ObjectDoesNotExist"
	case CIPStatus_ServiceFragmentation:
		return "ServiceFragmentation"
	case CIPStatus_NoStoredAttributeData:
		return "NoStoredAttributeData"
	case CIPStatus_StoreOperationFailure:
		return "StoreOperationFailure"
	case CIPStatus_RoutingFailureReqTooLarge:
		return "RoutingFailureReqTooLarge"
	case CIPStatus_RoutingFailureRespTooLarge:
		return "RoutingFailureRespTooLarge"
	case CIPStatus_MissingAttributeListEntry:
		return "MissingAttributeListEntry"
	case CIPStatus_InvalidAttributeValueList:
		return "InvalidAttributeValueList"
	case CIPStatus_EmbeddedServiceError:
		return "EmbeddedServiceError"
	case CIPStatus_VendorSpecificError:
		return "VendorSpecificError"
	case CIPStatus_InvalidParameter:
		return "InvalidParameter"
	case CIPStatus_WriteOnceValueOrMedium:
		return "WriteOnceValueOrMedium"
	case CIPStatus_InvalidReplyReceived:
		return "InvalidReplyReceived"
	case CIPStatus_BufferOverflow:
		return "BufferOverflow"
	case CIPStatus_MessageFormatError:
		return "MessageFormatError"
	case CIPStatus_KeyFailure:
		return "KeyFailure"
	case CIPStatus_PathSizeInvalid:
		return "PathSizeInvalid"
	case CIPStatus_UnexpectedAttribInList:
		return "UnexpectedAttribInList"
	case CIPStatus_InvalidMemberID:
		return "InvalidMemberID"
	case CIPStatus_MemberNotSettable:
		return "MemberNotSettable"
	case CIPStatus_Group2OnlyServerGeneralFailure:
		return "Group2OnlyServerGeneralFailure"
	case CIPStatus_UnknownModbusError:
		return "UnknownModbusError"
	case CIPStatus_AttributeNotGettable:
		return "AttributeNotGettable"
	default:
		if s > 0x2C && s < 0xD0 {
			return fmt.Sprintf("Unknown Error: 0x%X (reserved by CIP for future extensions)", byte(s))
		}
		return fmt.Sprintf("Unknown CIPStatus: 0x%X (reserved for object class and service errors)", byte(s))
	}
}
