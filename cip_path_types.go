package gologix

import (
	"encoding/binary"
	"fmt"
	"io"
)

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
		binary.Read(r, binary.LittleEndian, p)
		return nil
	default:
		return fmt.Errorf("expected 0x24 or 0x25 but got class size of %x", size)
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

type CIPClass uint16

// Used to indicate how many bytes are used for the data. If they are more than 8 bits,
// they actually actually take n+1 bytes.  First byte after specifier is a 0
type CIPClassSize byte

const (
	cipClass_8bit  CIPClassSize = 0x20
	cipClass_16bit CIPClassSize = 0x21
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
	var classSize CIPClassSize
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
	CipObject_ControllerInfo               CIPClass = 0xAC // don't know the official name
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
