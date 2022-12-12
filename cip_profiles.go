package gologix

// Here are predefined profiles
// "Any device that does not fall into the scope of one of the specialized
//  Device Profiles must use the Generic Device profile (0x2B) or a vendor-specific profile"

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
