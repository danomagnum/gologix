package gologix

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
