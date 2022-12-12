package gologix

import "encoding/binary"

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
