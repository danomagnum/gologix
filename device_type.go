package gologix

type DeviceType uint16

func (v DeviceType) Name() string {
	if name, ok := cipDeviceTypeNames[v]; ok {
		return name
	}

	return "Unknown"
}

// List of Device Types from the CIP vendor ID list.
// Obtained from the CIP Wireshark dissector on 2025-03-18.
// Available at https://fossies.org/linux/wireshark/epan/dissectors/packet-cip.c
var cipDeviceTypeNames = map[DeviceType]string{
	0x00: "Generic Device (deprecated)",
	0x02: "AC Drive",
	0x03: "Motor Overload",
	0x04: "Limit Switch",
	0x05: "Inductive Proximity Switch",
	0x06: "Photoelectric Sensor",
	0x07: "General Purpose Discrete I/O",
	0x09: "Resolver",
	0x0C: "Communications Adapter",
	0x0E: "Programmable Logic Controller",
	0x10: "Position Controller",
	0x13: "DC Drive",
	0x15: "Contactor",
	0x16: "Motor Starter",
	0x17: "Soft Start",
	0x18: "Human-Machine Interface",
	0x1A: "Mass Flow Controller",
	0x1B: "Pneumatic Valve",
	0x1C: "Vacuum Pressure Gauge",
	0x1D: "Process Control Value",
	0x1E: "Residual Gas Analyzer",
	0x1F: "DC Power Generator",
	0x20: "RF Power Generator",
	0x21: "Turbomolecular Vacuum Pump",
	0x22: "Encoder",
	0x23: "Safety Discrete I/O Device",
	0x24: "Fluid Flow Controller",
	0x25: "CIP Motion Drive",
	0x26: "CompoNet Repeater",
	0x27: "Mass Flow Controller, Enhanced",
	0x28: "CIP Modbus Device",
	0x29: "CIP Modbus Translator",
	0x2A: "Safety Analog I/O Device",
	0x2B: "Generic Device (keyable)",
	0x2C: "Managed Ethernet Switch",
	0x2D: "CIP Motion Safety Drive Device",
	0x2E: "Safety Drive Device",
	0x2F: "CIP Motion Encoder",
	0x30: "CIP Motion Converter",
	0x31: "CIP Motion I/O",
	0x32: "ControlNet Physical Layer Component",
	0x33: "Circuit Breaker",
	0x34: "HART Device",
	0x35: "CIP-HART Translator",
	0xC8: "Embedded Component",
}
