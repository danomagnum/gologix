package gologix

import (
	"errors"
	"fmt"
)

// to disconect we send two items - a null item and an unconnected data item for the unregister service
func (client *Client) Disconnect() error {
	if !client.Connected {
		return nil
	}
	var err error

	items := make([]cipItem, 2)
	items[0] = cipItem{} // null item

	reg_msg := msgCIPMessage_UnRegister{
		Service:                cipService_ForwardClose,
		CipPathSize:            0x02,
		ClassType:              cipClass_8bit,
		Class:                  0x06,
		InstanceType:           cipInstance_8bit,
		Instance:               0x01,
		Priority:               0x0A,
		TimeoutTicks:           0x0E,
		ConnectionSerialNumber: client.ConnectionSerialNumber,
		VendorID:               client.VendorID,
		OriginatorSerialNumber: client.serialNumber,
		PathSize:               3,                                           // 16 bit words
		Path:                   [6]byte{0x01, 0x00, 0x20, 0x02, 0x24, 0x01}, // TODO: generate paths automatically
	}

	items[1] = NewItem(cipItem_UnconnectedData, reg_msg)

	err = client.send(cipCommandSendRRData, SerializeItems(items)) // 0x65 is register session
	if err != nil {
		err2 := client.disconnect()
		return fmt.Errorf("couldn't send unconnect req %w: %v", err, err2)
	}
	client.disconnect()
	return nil

}

// module internal disconnect that closes the connection and cancels the watchdog/keepalive
func (client *Client) disconnect() error {
	if !client.Connected {
		return errors.New("already disconnected")
	}
	client.Connected = false

	// this will kill the keepalive goroutine
	close(client.cancel_keepalive)

	err := client.conn.Close()
	if err != nil {
		err = fmt.Errorf("error closing connection: %w", err)
	}
	return err
}

type msgCIPMessage_UnRegister struct {
	Service                CIPService
	CipPathSize            byte
	ClassType              CIPClassSize
	Class                  byte
	InstanceType           cipInstanceSize
	Instance               byte
	Priority               byte
	TimeoutTicks           byte
	ConnectionSerialNumber uint16
	VendorID               uint16
	OriginatorSerialNumber uint32
	PathSize               uint16
	Path                   [6]byte
}
