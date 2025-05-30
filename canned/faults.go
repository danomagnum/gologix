package canned

import (
	"fmt"

	"github.com/danomagnum/gologix"
)

type FaultType uint32

type FaultEvent struct {
	Timestamp  uint64
	FaultClass uint16
	FaultCode  uint16
	Detail     [30]byte
}

func (f FaultEvent) Description() FaultDescription {
	t, ok := FaultDescriptions[int(f.FaultClass)]
	if !ok {
		return FaultDescription{
			Type:             int(f.FaultClass),
			Code:             int(f.FaultCode),
			Display:          fmt.Sprintf("Unknown fault T%02d:C%02d", f.FaultClass, f.FaultCode),
			Description:      string(f.Detail[:]),
			CorrectiveAction: "Unknown",
			TypeName:         "Unknown",
		}
	}

	desc, ok := t[int(f.FaultCode)]
	if !ok {
		return FaultDescription{
			Type:             int(f.FaultClass),
			Code:             int(f.FaultCode),
			Display:          fmt.Sprintf("Unknown fault T%02d:C%02d", f.FaultClass, f.FaultCode),
			Description:      string(f.Detail[:]),
			CorrectiveAction: "Unknown",
			TypeName:         "Unknown",
		}
	}
	return desc

}
func (f FaultEvent) String() string {
	if f.FaultClass == 0 && f.FaultCode == 0 {
		return "No Fault"
	}
	return f.Description().Display
}

type FaultSummary struct {
	MajorType  FaultType
	MinorType  FaultType
	MajorCount uint16
	MinorCount uint16
	Events     [3]FaultEvent // Last 3 major faults.
}

const (
	FaultType_None           FaultType = 1 << iota // bit 1
	FaultType_Powerup                              // bit 2
	FaultType_IO                                   // bit 3
	FaultType_Program                              // bit 4
	_                                              // 5
	FaultType_Watchdog                             // bit 6
	FaultType_NVMemory                             // bit 7
	FaultType_ModeChange                           // bit 8
	FaultType_SerialPort                           // bit 9
	FaultType_EnergyStorage                        // bit 10
	FaultType_Motion                               // bit 11
	FaultType_Redundancy                           // bit 12
	FaultType_RTC                                  // bit 13
	FaultType_NonRecoverable                       // bit 14
	_                                              // 15
	FaultType_Communication                        // bit 16
	FaultType_Diagnostics                          // bit 17
	FaultType_CIPMotion                            // bit 18
	FaultType_Ethernet                             // bit 19
	FaultType_License                              // bit 20
	FaultType_Alarm                                // bit 21
	FaultType_OPCUA                                // bit 22
)

func (f FaultType) PowerupFault() bool {
	return f&FaultType_Powerup != 0
}

func (f FaultType) IOFault() bool {
	return f&FaultType_IO != 0
}

func (f FaultType) ProgramFault() bool {
	return f&FaultType_Program != 0
}

func (f FaultType) WatchdogFault() bool {
	return f&FaultType_Watchdog != 0
}

func (f FaultType) NVMemoryFault() bool {
	return f&FaultType_NVMemory != 0
}

func (f FaultType) ModeChangeFault() bool {
	return f&FaultType_ModeChange != 0
}

func (f FaultType) SerialPortFault() bool {
	return f&FaultType_SerialPort != 0
}

func (f FaultType) EnergyStorageFault() bool {
	return f&FaultType_EnergyStorage != 0
}

func (f FaultType) MotionFault() bool {
	return f&FaultType_Motion != 0
}

func (f FaultType) RedundancyFault() bool {
	return f&FaultType_Redundancy != 0
}

func (f FaultType) RTCFault() bool {
	return f&FaultType_RTC != 0
}

func (f FaultType) NonRecoverableFault() bool {
	return f&FaultType_NonRecoverable != 0
}

func (f FaultType) CommunicationFault() bool {
	return f&FaultType_Communication != 0
}

func (f FaultType) DiagnosticsFault() bool {
	return f&FaultType_Diagnostics != 0
}

func (f FaultType) CIPMotionFault() bool {
	return f&FaultType_CIPMotion != 0
}

func (f FaultType) EthernetFault() bool {
	return f&FaultType_Ethernet != 0
}

func (f FaultType) LicenseFault() bool {
	return f&FaultType_License != 0
}

func (f FaultType) AlarmFault() bool {
	return f&FaultType_Alarm != 0
}

func (f FaultType) OPCUAFault() bool {
	return f&FaultType_OPCUA != 0
}

func GetFaults(client *gologix.Client) (FaultSummary, error) {
	path, _ := gologix.Serialize(gologix.CIPClass(0x73), gologix.CIPInstance(0x1))
	item, err := client.GenericCIPMessage(gologix.CIPService_GetAttributeAll, path.Bytes(), nil)

	if err != nil {
		return FaultSummary{}, err
	}

	if item == nil {
		return FaultSummary{}, fmt.Errorf("item is nil")
	}

	response := FaultSummary{}

	err = item.DeSerialize(&response)
	if err != nil {
		return FaultSummary{}, fmt.Errorf("could not deserialize response: %v", err)
	}
	return response, nil
}
