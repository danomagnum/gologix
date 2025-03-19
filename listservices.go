package gologix

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func (client *Client) ListServices() ([]CIPListService, error) {
	client.Logger.Debug("listing services")

	_, data, err := client.send_recv_data(cipCommandListServices)
	if err != nil {
		return nil, err
	}

	items, err := readItems(data)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse items. %w", err)
	}

	services := make([]CIPListService, len(items))
	for i, item := range items {
		err := services[i].ParseFromBytes(item.Data)
		if err != nil {
			return nil, fmt.Errorf("problem reading list services response. %w", err)
		}
	}

	return services, nil
}

type CIPServiceTypeCode uint16

const (
	CIPServiceTypeCode_ListServices CIPServiceTypeCode = 0x0100
)

type CIPListService struct {
	EncapProtocolVersion uint16
	Capabilities         ServiceCapabilityFlags
	Name                 string
}

func (s *CIPListService) ParseFromBytes(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &s.EncapProtocolVersion)
	if err != nil {
		return fmt.Errorf("problem reading list services response. field EncapProtocolVersion %w", err)
	}

	err = binary.Read(buf, binary.LittleEndian, &s.Capabilities)
	if err != nil {
		return fmt.Errorf("problem reading list services response. field CapFlags %w", err)
	}

	// The name field is a 16 byte null terminated string
	name := make([]byte, 16)
	err = binary.Read(buf, binary.LittleEndian, &name)
	if err != nil {
		return fmt.Errorf("problem reading list services response. field Name %w", err)
	}
	// Remove any trailing NULL characters from the name
	s.Name = string(bytes.TrimRight(name, "\x00"))

	return nil
}
