package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

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
