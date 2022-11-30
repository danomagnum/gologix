package gologix

import (
	"encoding/binary"
	"io"
	"log"
)

type CIPType byte

// Go native types that correspond to logix types
// I'm not sure whether having interface here makes sense.
// On the one hand, we need to support composite types, but on the other this lets it accept anything
// which doesn't seem right.
type GoLogixTypes interface {
	bool | byte | uint16 | int16 | uint32 | int32 | uint64 | int64 | float32 | float64 | string | interface{}
}

// return the CIPType that corresponds to go type T
func GoTypeToCIPType[T GoLogixTypes]() CIPType {
	var t T
	return GoVarToCIPType(t)
}

// return the CIPType that corresponds to go type of variable T
func GoVarToCIPType(T any) CIPType {
	switch T.(type) {
	case byte:
		return CIPTypeBOOL
	case uint16:
		return CIPTypeUINT
	case int16:
		return CIPTypeINT
	case uint32:
		return CIPTypeUDINT
	case int32:
		return CIPTypeDINT
	case uint64:
		return CIPTypeLWORD
	case int64:
		return CIPTypeLINT
	case float32:
		return CIPTypeREAL
	case float64:
		return CIPTypeLREAL
	case string:
		return CIPTypeSTRING
	case interface{}:
		return CIPTypeStruct
	}
	return CIPTypeUnknown
}

const (
	CIPTypeUnknown CIPType = 0x00
	CIPTypeStruct  CIPType = 0xA0 // also used for strings.  Not sure what's up with CIPTypeSTRING
	CIPTypeBOOL    CIPType = 0xC1
	CIPTypeBYTE    CIPType = 0xD1 // 8 bits packed into one byte
	CIPTypeSINT    CIPType = 0xC2
	CIPTypeINT     CIPType = 0xC3
	CIPTypeDINT    CIPType = 0xC4
	CIPTypeLINT    CIPType = 0xC5
	CIPTypeUSINT   CIPType = 0xC6
	CIPTypeUINT    CIPType = 0xC7
	CIPTypeUDINT   CIPType = 0xC8
	CIPTypeLWORD   CIPType = 0xC9
	CIPTypeREAL    CIPType = 0xCA
	CIPTypeLREAL   CIPType = 0xCB
	CIPTypeWORD    CIPType = 0xD2
	CIPTypeDWORD   CIPType = 0xD3

	// As far as I can tell CIPTypeSTRING isn't actually used in the controllers. Strings actually come
	// accross as 0xA0 = CIPTypeStruct.  In this library we're using this as kind of a flag to keep track of whether
	// a structure is a string or not.
	CIPTypeSTRING CIPType = 0xDA
)

// return the size in bytes of the data structure
func (c CIPType) Size() int {
	switch c {
	case CIPTypeUnknown:
		return 0
	case CIPTypeStruct:
		return 88
	case CIPTypeBOOL:
		return 1
	case CIPTypeBYTE:
		return 1
	case CIPTypeSINT:
		return 1
	case CIPTypeINT:
		return 2
	case CIPTypeDINT:
		return 4
	case CIPTypeLINT:
		return 8
	case CIPTypeUSINT:
		return 1
	case CIPTypeUINT:
		return 2
	case CIPTypeUDINT:
		return 4
	case CIPTypeLWORD:
		return 8
	case CIPTypeREAL:
		return 4
	case CIPTypeLREAL:
		return 8
	case CIPTypeWORD:
		return 2
	case CIPTypeDWORD:
		return 4
	case CIPTypeSTRING:
		return 1
	default:
		return 0
	}
}

// return a buffer that can hold the data structure
func (c CIPType) NewBuffer() *[]byte {
	buf := make([]byte, c.Size())
	return &buf
}

// human readable version of the cip type for printing.
func (c CIPType) String() string {
	switch c {
	case CIPTypeUnknown:
		return "0x00 - Unknown"
	case CIPTypeStruct:
		return "0xA0 - Struct"
	case CIPTypeBOOL:
		return "0xC1 - BOOL"
	case CIPTypeBYTE:
		return "0xD1 - BYTE"
	case CIPTypeSINT:
		return "0xC2 - SINT"
	case CIPTypeINT:
		return "0xC3 - INT"
	case CIPTypeDINT:
		return "0xC4 - DINT"
	case CIPTypeLINT:
		return "0xC5 - LINT"
	case CIPTypeUSINT:
		return "0xC6 - USINT"
	case CIPTypeUINT:
		return "0xC7 - UINT"
	case CIPTypeUDINT:
		return "0xC8 - UDINT"
	case CIPTypeLWORD:
		return "0xC9 - LWORD"
	case CIPTypeREAL:
		return "0xCA - REAL"
	case CIPTypeLREAL:
		return "0xCB - LREAL"
	case CIPTypeWORD:
		return "0xD2 - WORD"
	case CIPTypeDWORD:
		return "0xD3 - DWORD"
	case CIPTypeSTRING:
		return "0xDA - String"
	default:
		return "0 - Unknown"
	}
}

func (t CIPType) readValue(r io.Reader) any {
	return readValue(t, r)
}

// readValue reads one unit of cip data type t into the correct go type.
// To do this it reads the needed number of bytes from r.
// It returns the value as an any so the caller will have to do a cast to get it back
func readValue(t CIPType, r io.Reader) any {

	var value any
	var err error
	switch t {
	case CIPTypeUnknown:
		panic("Unknown type.")
	case CIPTypeStruct:
		panic("Struct!")
	case CIPTypeBOOL:
		var trueval bool
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeBYTE:
		var trueval byte
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeSINT:
		var trueval byte
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeINT:
		var trueval int16
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeDINT:
		var trueval int32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeLINT:
		var trueval int64
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeUSINT:
		var trueval uint8
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeUINT:
		var trueval uint16
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeUDINT:
		var trueval uint32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeLWORD:
		var trueval uint64
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeREAL:
		var trueval float32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeLREAL:
		var trueval float64
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeWORD:
		var trueval uint16
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeDWORD:
		var trueval uint32
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	case CIPTypeSTRING:
		var trueval [86]byte
		err = binary.Read(r, binary.LittleEndian, &trueval)
		value = trueval
	default:
		panic("Default type.")

	}
	if err != nil {
		log.Printf("Problem reading %s as one unit of %T. %v", t, value, err)
	}
	//log.Printf("type %v. value %v", t, value)
	return value
}

// These next few types are used to indicate how many
// bytes are used for the data.  If they are more than 8 bits,
// they actually actually takes n+1 bytes.  First byte is a 0

type CIPInstanceType byte

const (
	CIPInstance_8bit  CIPInstanceType = 0x24
	CIPInstance_16bit CIPInstanceType = 0x25
)

type CIPClassType byte

const (
	CIPClass_8bit  CIPClassType = 0x20
	CIPClass_16bit CIPClassType = 0x21
)

type CIPAttributeType byte

const (
	CIPAttribute_8bit  CIPAttributeType = 0x30
	CIPAttribute_16bit CIPAttributeType = 0x31
)

type CIPElementType byte

const (
	CIPElement_8bit  CIPElementType = 0x28
	CIPElement_16bit CIPElementType = 0x29
	CIPElement_32bit CIPElementType = 0x2A
)

type CIPInstance uint16

func (p CIPInstance) Bytes() []byte {
	if p < 256 {
		b := make([]byte, 2)
		b[0] = byte(CIPInstance_8bit)
		b[1] = byte(p)
		return b
	} else {

		b := make([]byte, 4)
		b[0] = byte(CIPInstance_16bit)
		binary.LittleEndian.PutUint16(b[2:], uint16(p))
		return b
	}
}

// Here are the objects

type CIPObject uint16

func (p CIPObject) Bytes() []byte {
	if p < 256 {
		b := make([]byte, 2)
		b[0] = byte(CIPClass_8bit)
		b[1] = byte(p)
		return b
	} else {

		b := make([]byte, 4)
		b[0] = byte(CIPClass_16bit)
		binary.LittleEndian.PutUint16(b[2:], uint16(p))
		return b
	}
}

const (
	CIPObject_Assembly                     CIPObject = 0x04
	CIPObject_AckHandler                   CIPObject = 0x2B
	CIPObject_Symbol                       CIPObject = 0x6B
	CIPObject_Template                     CIPObject = 0x6C
	CIPObject_Connection                   CIPObject = 0x05
	CIPObject_ConnectionConfig             CIPObject = 0xF3
	CIPObject_ConnectionManager            CIPObject = 0x06
	CIPObject_File                         CIPObject = 0x37
	CIPObject_Identity                     CIPObject = 0x01
	CIPObject_MessageRouter                CIPObject = 0x02
	CIPObject_OriginatorConnList           CIPObject = 0x45
	CIPObject_Parameter                    CIPObject = 0x0F
	CIPObject_ParameterGroup               CIPObject = 0x10
	CIPObject_Port                         CIPObject = 0xF4
	CIPObject_Register                     CIPObject = 0x07
	CIPObject_Selection                    CIPObject = 0x2E
	CIPObject_Drive                        CIPObject = 0x2A
	CIPObject_AnalogGroup                  CIPObject = 0x22
	CIPObject_AnalogInputGroup             CIPObject = 0x20
	CIPObject_AnalogInputPoint             CIPObject = 0x0A
	CIPObject_AnalogOutputGroup            CIPObject = 0x21
	CIPObject_AnalogOutputPoint            CIPObject = 0x0B
	CIPObject_BaseEnergy                   CIPObject = 0x4E
	CIPObject_BlockSequencer               CIPObject = 0x26
	CIPObject_CommandBlock                 CIPObject = 0x27
	CIPObject_ControlSupervisor            CIPObject = 0x29
	CIPObject_DiscreteGroup                CIPObject = 0x1F
	CIPObject_DiscreteInputGroup           CIPObject = 0x1D
	CIPObject_DiscreteOutputGroup          CIPObject = 0x1E
	CIPObject_DiscreteInputPoint           CIPObject = 0x08
	CIPObject_DiscreteOutputPoint          CIPObject = 0x09
	CIPObject_ElectricalEnergy             CIPObject = 0x4F
	CIPObject_EventLog                     CIPObject = 0x41
	CIPObject_Group                        CIPObject = 0x12
	CIPObject_MotionAxis                   CIPObject = 0x42
	CIPObject_MotorData                    CIPObject = 0x28
	CIPObject_NonElectricalEnergy          CIPObject = 0x50
	CIPObject_Overload                     CIPObject = 0x2C
	CIPObject_PositionController           CIPObject = 0x25
	CIPObject_PositionControlSupervisor    CIPObject = 0x24
	CIPObject_PositionSensor               CIPObject = 0x23
	CIPObject_PowerCurtailment             CIPObject = 0x5C
	CIPObject_PowerManagement              CIPObject = 0x53
	CIPObject_PresenceSensing              CIPObject = 0x0E
	CIPObject_SAnalogActuator              CIPObject = 0x32
	CIPObject_SAnalogSensor                CIPObject = 0x31
	CIPObject_SDeviceSupervisor            CIPObject = 0x30
	CIPObject_SGasCalibration              CIPObject = 0x34
	CIPObject_SPartialPressure             CIPObject = 0x38
	CIPObject_SSensorCalibration           CIPObject = 0x40
	CIPObject_SSingleStageController       CIPObject = 0x33
	CIPObject_SafetyAnalogInputGroup       CIPObject = 0x4A
	CIPObject_SafetyAnalogInputPoint       CIPObject = 0x49
	CIPObject_SafetyDualChannelFeedback    CIPObject = 0x59
	CIPObject_SafetyFeedback               CIPObject = 0x5A
	CIPObject_SafetyDiscreteInputGroup     CIPObject = 0x3E
	CIPObject_SafetyDiscreteInputPoint     CIPObject = 0xeD
	CIPObject_SafetyDiscreteOutputGroup    CIPObject = 0x3C
	CIPObject_SafetyDiscreteOutputPoint    CIPObject = 0x3B
	CIPObject_SafetyDualChannelAnalogInput CIPObject = 0x4B
	CIPObject_SafetyDualChannelOutput      CIPObject = 0x3F
	CIPObject_SafetyLimitFunctions         CIPObject = 0x5B
	CIPObject_SafetyStopFunctions          CIPObject = 0x5A
	CIPObject_SafetySupervisor             CIPObject = 0x39
	CIPObject_SafetyValidator              CIPObject = 0x3A
	CIPObject_SoftStart                    CIPObject = 0x2D
	CIPObject_TargetConnectionList         CIPObject = 0x4D
	CIPObject_TimeSync                     CIPObject = 0x43
	CIPObject_TripPoint                    CIPObject = 0x35
	CIPObject_BaseSwitch                   CIPObject = 0x51
	CIPObject_CompoNetLink                 CIPObject = 0xF7
	CIPObject_CompoNetRepeater             CIPObject = 0xF8
	CIPObject_ControlNet                   CIPObject = 0xF0
	CIPObject_ControlNetKeeper             CIPObject = 0xF1
	CIPObject_ControlNetScheduling         CIPObject = 0xF2
	CIPObject_DLR                          CIPObject = 0x47
	CIPObject_DeviceNet                    CIPObject = 0x03
	CIPObject_EthernetLink                 CIPObject = 0xF6
	CIPObject_Modbus                       CIPObject = 0x44
	CIPObject_ModbusSerial                 CIPObject = 0x46
	CIPObject_ParallelRedundancyProtocol   CIPObject = 0x56
	CIPObject_PRPNodesTable                CIPObject = 0x57
	CIPObject_SERCOSIIILink                CIPObject = 0x4C
	CIPObject_SNMP                         CIPObject = 0x52
	CIPObject_QoS                          CIPObject = 0x48
	CIPObject_RSTPBridge                   CIPObject = 0x54
	CIPObject_RSTPPort                     CIPObject = 0x55
	CIPObject_TCPIP                        CIPObject = 0xF5
	CIPObject_PCCC                         CIPObject = 0x67
)

// Here are predefined profiles
// "Any device that does not fall into the scope of one of the specialized
//  Device Profiles must use the Generic Device profile (0x2B) or a vendor-specific profile"

type CIPDeviceProfile byte

const (
	CIPDevice_ACDrive                          CIPDeviceProfile = 0x02
	CIPDevice_CIPModbusDevice                  CIPDeviceProfile = 0x28
	CIPDevice_CIPModbusTranslator              CIPDeviceProfile = 0x29
	CIPDevice_CIPMotionDrive                   CIPDeviceProfile = 0x25
	CIPDevice_CIPMotionEncoder                 CIPDeviceProfile = 0x2F
	CIPDevice_CIPMotionIO                      CIPDeviceProfile = 0x31
	CIPDevice_CIPMotionSafetyDrive             CIPDeviceProfile = 0x2D
	CIPDevice_CommsAdapter                     CIPDeviceProfile = 0x0C
	CIPDevice_CompoNetRepeater                 CIPDeviceProfile = 0x26
	CIPDevice_Contactor                        CIPDeviceProfile = 0x15
	CIPDevice_ControlNetPhysicalLayerComponent CIPDeviceProfile = 0x32
	CIPDevice_DCDrive                          CIPDeviceProfile = 0x13
	CIPDevice_DCPowerGenerator                 CIPDeviceProfile = 0x1F
	CIPDevice_EmbeddedComponent                CIPDeviceProfile = 0xC8
	CIPDevice_Encoder                          CIPDeviceProfile = 0x22
	CIPDevice_EnhancedMassFlowController       CIPDeviceProfile = 0x27
	CIPDevice_FluidFlowController              CIPDeviceProfile = 0x24
	CIPDevice_DiscreteIO                       CIPDeviceProfile = 0x07
	CIPDevice_Generic                          CIPDeviceProfile = 0x2B
	CIPDevice_HMI                              CIPDeviceProfile = 0x18
	CIPDevice_ProxSwitch                       CIPDeviceProfile = 0x05
	CIPDevice_LimitSwitch                      CIPDeviceProfile = 0x04
	CIPDevice_ManagedEthernetSwitch            CIPDeviceProfile = 0x2C
	CIPDevice_MassFlowController               CIPDeviceProfile = 0x2C
	CIPDevice_MotorOverload                    CIPDeviceProfile = 0x03
	CIPDevice_MotorStarter                     CIPDeviceProfile = 0x16
	CIPDevice_Photoeye                         CIPDeviceProfile = 0x06
	CIPDevice_PneumaticValve                   CIPDeviceProfile = 0x1B
	CIPDevice_PositionController               CIPDeviceProfile = 0x10
	CIPDevice_ProcessControlValve              CIPDeviceProfile = 0x1D
	CIPDevice_PLC                              CIPDeviceProfile = 0x0E
	CIPDevice_ResidualGasAnalyzer              CIPDeviceProfile = 0x1E
	CIPDevice_Resolver                         CIPDeviceProfile = 0x09
	CIPDevice_RFPowerGenerator                 CIPDeviceProfile = 0x20
	CIPDevice_SafetyAnalogIO                   CIPDeviceProfile = 0x2A
	CIPDevice_SafetyDrive                      CIPDeviceProfile = 0x2E
	CIPDevice_SafetyDiscreteIO                 CIPDeviceProfile = 0x23
	CIPDevice_SoftStart                        CIPDeviceProfile = 0x17
	CIPDevice_VacuumPump                       CIPDeviceProfile = 0x21
	CIPDevice_PressureGauge                    CIPDeviceProfile = 0x1C
)
