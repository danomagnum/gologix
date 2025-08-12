package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// ListIdentity queries the device for its identity information and capabilities.
//
// This function sends a List Identity command to retrieve comprehensive information
// about the connected device. This is essential for device discovery, identification,
// and determining device capabilities and configuration.
//
// Returns a listIdentityResponeBody structure containing:
//   - EncapProtocolVersion: EtherNet/IP encapsulation protocol version
//   - SocketAddress: Device network addressing information
//   - Vendor: Vendor ID (manufacturer identifier)
//   - DeviceType: Type of device (PLC, HMI, drive, etc.)
//   - ProductCode: Manufacturer-specific product identifier
//   - Revision: Device firmware/hardware revision information
//   - Status: Current device status flags
//   - SerialNumber: Unique device serial number
//   - ProductName: Human-readable product name
//   - State: Current device operational state
//
// Example:
//   identity, err := client.ListIdentity()
//   if err != nil {
//       log.Fatal(err)
//   }
//
//   fmt.Printf("Device: %s\n", identity.ProductName)
//   fmt.Printf("Vendor: %s\n", identity.Vendor)
//   fmt.Printf("Product Code: %d\n", identity.ProductCode)
//   fmt.Printf("Serial Number: 0x%08X\n", identity.SerialNumber)
//   fmt.Printf("Revision: %d.%d\n", identity.Revision.Major, identity.Revision.Minor)
//   fmt.Printf("Status: 0x%04X\n", identity.Status)
//
// This function is commonly used for:
//   - Device discovery and inventory
//   - Verifying device compatibility
//   - Diagnostic and troubleshooting information
//   - Network device scanning
//
// Note: The device must be connected before calling this function.
func (client *Client) ListIdentity() (*listIdentityResponeBody, error) {
	client.Logger.Debug("listing identity")

	_, data, err := client.send_recv_data(cipCommandListIdentity)
	if err != nil {
		return nil, err
	}

	items, err := readItems(data)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse items. %w", err)
	}

	if len(items) != 1 {
		return nil, fmt.Errorf("expected 1 item, got %d", len(items))
	}

	// The response only contains one item.
	response := listIdentityResponeBody{}
	err = response.ParseFromBytes(items[0].Data)
	if err != nil {
		return nil, fmt.Errorf("problem reading list identity response. %w", err)
	}

	return &response, nil
}

type listIdentityResponeBody struct {
	EncapProtocolVersion uint16
	SocketAddress        listIdentitySocketAddress
	Vendor               VendorId
	DeviceType           DeviceType
	ProductCode          uint16
	Revision             uint16
	Status               uint16
	SerialNumber         uint32
	ProductName          string
	State                uint8
}

type listIdentitySocketAddress struct {
	Family  uint16
	Port    uint16
	Address uint32
	Zero0   uint32
	Zero1   uint32
}

func (s *listIdentityResponeBody) ParseFromBytes(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &s.EncapProtocolVersion)
	if err != nil {
		return fmt.Errorf("problem reading list identity response. field EncapProtocolVersion %w", err)
	}

	// For some reason this field in particular is big endian.
	err = binary.Read(buf, binary.BigEndian, &s.SocketAddress)
	if err != nil {
		return fmt.Errorf("problem reading list identity response. field SocketAddress %w", err)
	}

	err = binary.Read(buf, binary.LittleEndian, &s.Vendor)
	if err != nil {
		return fmt.Errorf("problem reading list identity response. field Vendor %w", err)
	}
	err = binary.Read(buf, binary.LittleEndian, &s.DeviceType)
	if err != nil {
		return fmt.Errorf("problem reading list identity response. field DeviceType %w", err)
	}
	err = binary.Read(buf, binary.LittleEndian, &s.ProductCode)
	if err != nil {
		return fmt.Errorf("problem reading list identity response. field ProductCode %w", err)
	}
	err = binary.Read(buf, binary.LittleEndian, &s.Revision)
	if err != nil {
		return fmt.Errorf("problem reading list identity response. field Revision %w", err)
	}
	err = binary.Read(buf, binary.LittleEndian, &s.Status)
	if err != nil {
		return fmt.Errorf("problem reading list identity response. field Status %w", err)
	}
	err = binary.Read(buf, binary.LittleEndian, &s.SerialNumber)
	if err != nil {
		return fmt.Errorf("problem reading list identity response. field SerialNumber %w", err)
	}

	productNameLength := uint8(0)
	err = binary.Read(buf, binary.LittleEndian, &productNameLength)
	if err != nil {
		return fmt.Errorf("problem reading list identity response. field ProductNameLength %w", err)
	}

	productName := make([]byte, productNameLength)
	err = binary.Read(buf, binary.LittleEndian, &productName)
	if err != nil {
		return fmt.Errorf("problem reading list identity response. field ProductName %w", err)
	}
	s.ProductName = string(productName)

	err = binary.Read(buf, binary.LittleEndian, &s.State)
	if err != nil {
		return fmt.Errorf("problem reading list identity response. field State %w", err)
	}

	return nil
}
