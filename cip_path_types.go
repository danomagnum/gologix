package gologix

import "encoding/binary"

type CIPAddress byte

func (p CIPAddress) Bytes() []byte {
	return []byte{byte(p)}
}

func (p CIPAddress) Len() int {
	return 0
}

// Here are the objects

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
func (p CIPAttribute) Len() int {
	if p < 256 {
		return 1
	}
	return 2
}

// currently unused
/*
const (
	cipAttribute_Data CIPAttribute = 0x03
)
*/

// Here are the objects

type CIPElement uint32

// Used to indicate how many bytes are used for the data. If they are more than 8 bits,
// they actually actually take n+1 bytes.  First byte after specifier is a 0
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

// Used to indicate how many bytes are used for the data. If they are more than 8 bits,
// they actually actually take n+1 bytes.  First byte after specifier is a 0
type cipInstanceSize byte

const (
	cipInstance_8bit  cipInstanceSize = 0x24
	cipInstance_16bit cipInstanceSize = 0x25
)

type CIPInstance uint16

func (p CIPInstance) Bytes() []byte {
	if p < 256 {
		b := make([]byte, 2)
		b[0] = byte(cipInstance_8bit)
		b[1] = byte(p)
		return b
	} else {

		b := make([]byte, 4)
		b[0] = byte(cipInstance_16bit)
		binary.LittleEndian.PutUint16(b[2:], uint16(p))
		return b
	}
}
func (p CIPInstance) Len() int {
	if p < 256 {
		return 2
	}
	return 4
}

type JustBytes []byte

func (p JustBytes) Bytes() []byte {
	if len(p) == 1 {
		b := make([]byte, len(p)+1)
		b[0] = byte(cipInstance_8bit)
		copy(b[1:], p)
		return b
	} else {
		b := make([]byte, len(p)+2)
		b[0] = byte(cipInstance_16bit)
		copy(b[2:], p)
		return b
	}

}

// Here are the objects

type CIPObject uint16

// Used to indicate how many bytes are used for the data. If they are more than 8 bits,
// they actually actually take n+1 bytes.  First byte after specifier is a 0
type CIPClassSize byte

const (
	cipClass_8bit  CIPClassSize = 0x20
	cipClass_16bit CIPClassSize = 0x21
)

func (p CIPObject) Bytes() []byte {
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
func (p CIPObject) Len() int {
	if p < 256 {
		return 2
	}
	return 4
}

const (
	cipObject_Assembly                     CIPObject = 0x04
	cipObject_AckHandler                   CIPObject = 0x2B
	cipObject_Symbol                       CIPObject = 0x6B
	cipObject_Template                     CIPObject = 0x6C
	cipObject_Connection                   CIPObject = 0x05
	cipObject_ConnectionConfig             CIPObject = 0xF3
	cipObject_ConnectionManager            CIPObject = 0x06
	cipObject_File                         CIPObject = 0x37
	cipObject_Identity                     CIPObject = 0x01
	cipObject_MessageRouter                CIPObject = 0x02
	cipObject_OriginatorConnList           CIPObject = 0x45
	cipObject_Parameter                    CIPObject = 0x0F
	cipObject_ParameterGroup               CIPObject = 0x10
	cipObject_Port                         CIPObject = 0xF4
	cipObject_Register                     CIPObject = 0x07
	cipObject_Selection                    CIPObject = 0x2E
	cipObject_Drive                        CIPObject = 0x2A
	cipObject_AnalogGroup                  CIPObject = 0x22
	cipObject_AnalogInputGroup             CIPObject = 0x20
	cipObject_AnalogInputPoint             CIPObject = 0x0A
	cipObject_AnalogOutputGroup            CIPObject = 0x21
	cipObject_AnalogOutputPoint            CIPObject = 0x0B
	cipObject_BaseEnergy                   CIPObject = 0x4E
	cipObject_BlockSequencer               CIPObject = 0x26
	cipObject_CommandBlock                 CIPObject = 0x27
	cipObject_ControlSupervisor            CIPObject = 0x29
	cipObject_DiscreteGroup                CIPObject = 0x1F
	cipObject_DiscreteInputGroup           CIPObject = 0x1D
	cipObject_DiscreteOutputGroup          CIPObject = 0x1E
	cipObject_DiscreteInputPoint           CIPObject = 0x08
	cipObject_DiscreteOutputPoint          CIPObject = 0x09
	cipObject_ElectricalEnergy             CIPObject = 0x4F
	cipObject_EventLog                     CIPObject = 0x41
	cipObject_Group                        CIPObject = 0x12
	cipObject_MotionAxis                   CIPObject = 0x42
	cipObject_MotorData                    CIPObject = 0x28
	cipObject_NonElectricalEnergy          CIPObject = 0x50
	cipObject_Overload                     CIPObject = 0x2C
	cipObject_PositionController           CIPObject = 0x25
	cipObject_PositionControlSupervisor    CIPObject = 0x24
	cipObject_PositionSensor               CIPObject = 0x23
	cipObject_PowerCurtailment             CIPObject = 0x5C
	cipObject_PowerManagement              CIPObject = 0x53
	cipObject_PresenceSensing              CIPObject = 0x0E
	cipObject_SAnalogActuator              CIPObject = 0x32
	cipObject_SAnalogSensor                CIPObject = 0x31
	cipObject_SDeviceSupervisor            CIPObject = 0x30
	cipObject_SGasCalibration              CIPObject = 0x34
	cipObject_SPartialPressure             CIPObject = 0x38
	cipObject_SSensorCalibration           CIPObject = 0x40
	cipObject_SSingleStageController       CIPObject = 0x33
	cipObject_SafetyAnalogInputGroup       CIPObject = 0x4A
	cipObject_SafetyAnalogInputPoint       CIPObject = 0x49
	cipObject_SafetyDualChannelFeedback    CIPObject = 0x59
	cipObject_SafetyFeedback               CIPObject = 0x5A
	cipObject_SafetyDiscreteInputGroup     CIPObject = 0x3E
	cipObject_SafetyDiscreteInputPoint     CIPObject = 0xeD
	cipObject_SafetyDiscreteOutputGroup    CIPObject = 0x3C
	cipObject_SafetyDiscreteOutputPoint    CIPObject = 0x3B
	cipObject_SafetyDualChannelAnalogInput CIPObject = 0x4B
	cipObject_SafetyDualChannelOutput      CIPObject = 0x3F
	cipObject_SafetyLimitFunctions         CIPObject = 0x5B
	cipObject_SafetyStopFunctions          CIPObject = 0x5A
	cipObject_SafetySupervisor             CIPObject = 0x39
	cipObject_SafetyValidator              CIPObject = 0x3A
	cipObject_SoftStart                    CIPObject = 0x2D
	cipObject_TargetConnectionList         CIPObject = 0x4D
	cipObject_TimeSync                     CIPObject = 0x43
	cipObject_TripPoint                    CIPObject = 0x35
	cipObject_BaseSwitch                   CIPObject = 0x51
	cipObject_CompoNetLink                 CIPObject = 0xF7
	cipObject_CompoNetRepeater             CIPObject = 0xF8
	cipObject_ControlNet                   CIPObject = 0xF0
	cipObject_ControlNetKeeper             CIPObject = 0xF1
	cipObject_ControlNetScheduling         CIPObject = 0xF2
	cipObject_DLR                          CIPObject = 0x47
	cipObject_DeviceNet                    CIPObject = 0x03
	cipObject_EthernetLink                 CIPObject = 0xF6
	cipObject_Modbus                       CIPObject = 0x44
	cipObject_ModbusSerial                 CIPObject = 0x46
	cipObject_ParallelRedundancyProtocol   CIPObject = 0x56
	cipObject_PRPNodesTable                CIPObject = 0x57
	cipObject_SERCOSIIILink                CIPObject = 0x4C
	cipObject_SNMP                         CIPObject = 0x52
	cipObject_QoS                          CIPObject = 0x48
	cipObject_RSTPBridge                   CIPObject = 0x54
	cipObject_RSTPPort                     CIPObject = 0x55
	cipObject_TCPIP                        CIPObject = 0xF5
	cipObject_PCCC                         CIPObject = 0x67
	cipObject_ControllerInfo               CIPObject = 0xAC // don't know the official name
)

// Here are predefined profiles
// "Any device that does not fall into the scope of one of the specialized
//  Device Profiles must use the Generic Device profile (0x2B) or a vendor-specific profile"
// commented out because they are currently unused
/*
type CIPDeviceProfile byte

const (
	cipDevice_ACDrive                          CIPDeviceProfile = 0x02
	cipDevice_CIPModbusDevice                  CIPDeviceProfile = 0x28
	cipDevice_CIPModbusTranslator              CIPDeviceProfile = 0x29
	cipDevice_CIPMotionDrive                   CIPDeviceProfile = 0x25
	cipDevice_CIPMotionEncoder                 CIPDeviceProfile = 0x2F
	cipDevice_CIPMotionIO                      CIPDeviceProfile = 0x31
	cipDevice_CIPMotionSafetyDrive             CIPDeviceProfile = 0x2D
	cipDevice_CommsAdapter                     CIPDeviceProfile = 0x0C
	cipDevice_CompoNetRepeater                 CIPDeviceProfile = 0x26
	cipDevice_Contactor                        CIPDeviceProfile = 0x15
	cipDevice_ControlNetPhysicalLayerComponent CIPDeviceProfile = 0x32
	cipDevice_DCDrive                          CIPDeviceProfile = 0x13
	cipDevice_DCPowerGenerator                 CIPDeviceProfile = 0x1F
	cipDevice_EmbeddedComponent                CIPDeviceProfile = 0xC8
	cipDevice_Encoder                          CIPDeviceProfile = 0x22
	cipDevice_EnhancedMassFlowController       CIPDeviceProfile = 0x27
	cipDevice_FluidFlowController              CIPDeviceProfile = 0x24
	cipDevice_DiscreteIO                       CIPDeviceProfile = 0x07
	cipDevice_Generic                          CIPDeviceProfile = 0x2B
	cipDevice_HMI                              CIPDeviceProfile = 0x18
	cipDevice_ProxSwitch                       CIPDeviceProfile = 0x05
	cipDevice_LimitSwitch                      CIPDeviceProfile = 0x04
	cipDevice_ManagedEthernetSwitch            CIPDeviceProfile = 0x2C
	cipDevice_MassFlowController               CIPDeviceProfile = 0x2C
	cipDevice_MotorOverload                    CIPDeviceProfile = 0x03
	cipDevice_MotorStarter                     CIPDeviceProfile = 0x16
	cipDevice_Photoeye                         CIPDeviceProfile = 0x06
	cipDevice_PneumaticValve                   CIPDeviceProfile = 0x1B
	cipDevice_PositionController               CIPDeviceProfile = 0x10
	cipDevice_ProcessControlValve              CIPDeviceProfile = 0x1D
	cipDevice_PLC                              CIPDeviceProfile = 0x0E
	cipDevice_ResidualGasAnalyzer              CIPDeviceProfile = 0x1E
	cipDevice_Resolver                         CIPDeviceProfile = 0x09
	cipDevice_RFPowerGenerator                 CIPDeviceProfile = 0x20
	cipDevice_SafetyAnalogIO                   CIPDeviceProfile = 0x2A
	cipDevice_SafetyDrive                      CIPDeviceProfile = 0x2E
	cipDevice_SafetyDiscreteIO                 CIPDeviceProfile = 0x23
	cipDevice_SoftStart                        CIPDeviceProfile = 0x17
	cipDevice_VacuumPump                       CIPDeviceProfile = 0x21
	cipDevice_PressureGauge                    CIPDeviceProfile = 0x1C
)

*/
