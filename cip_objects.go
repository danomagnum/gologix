package gologix

import "encoding/binary"

// Here are the objects

type CIPObject uint16

// Used to indicate how many bytes are used for the data. If they are more than 8 bits,
// they actually actually take n+1 bytes.  First byte after specifier is a 0
type CIPClassSize byte

const (
	CIPClass_8bit  CIPClassSize = 0x20
	CIPClass_16bit CIPClassSize = 0x21
)

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
	CIPObject_ControllerInfo               CIPObject = 0xAC // don't know the official name
)
