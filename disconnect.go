package gologix

import (
	"fmt"
	"log/slog"
)

// You will want to defer this after a successful Connect() to make sure you free up the controller resources
// to disconnect we send two items - a null item and an unconnected data item for the unregister service
func (client *Client) Disconnect() error {
	if !client.Connected {
		return nil
	}
	var err error
	client.SLogger.Info("starting disconnection")

	client.Connected = false
	if client.KeepAliveRunning {
		close(client.cancel_keepalive)
	}

	items := make([]CIPItem, 2)
	items[0] = CIPItem{} // null item
	items[1] = CIPItem{Header: cipItemHeader{ID: cipItem_UnconnectedData}}

	path, err := Serialize(
		client.Controller.Path,
		CipObject_MessageRouter,
		CIPInstance(1),
	)
	if err != nil {
		client.SLogger.Error("Error serializing path", slog.Any("err", err))
		return fmt.Errorf("error serializing path: %w", err)
	}

	msg := msgCipUnRegister{
		Service:                CIPService_ForwardClose,
		CipPathSize:            0x02,
		ClassType:              cipClass_8bit,
		Class:                  0x06,
		InstanceType:           cipInstance_8bit,
		Instance:               0x01,
		Priority:               0x0A,
		TimeoutTicks:           0x0E,
		ConnectionSerialNumber: client.ConnectionSerialNumber,
		VendorID:               client.VendorId,
		OriginatorSerialNumber: client.SerialNumber,
		PathSize:               byte(path.Len() / 2),
		Reserved:               0x00,
	}

	items[1].Serialize(msg)
	items[1].Serialize(path)

	itemData, err := serializeItems(items)
	if err != nil {
		client.SLogger.Error(
			"unable to serialize itemData. Forcing connection closed",
			slog.Any("err", err),
		)
	} else {
		header, data, err := client.send_recv_data(cipCommandSendRRData, itemData)
		if err != nil {
			client.SLogger.Error(
				"error sending disconnect request",
				slog.Any("err", err),
			)
		}

		_, err = client.parseResponse(&header, data)
		if err != nil {
			client.SLogger.Error(
				"error parsing disconnect response",
				slog.Any("err", err),
			)
		}
	}

	err = client.conn.Close()
	if err != nil {
		client.SLogger.Error("error closing connection", slog.Any("err", err))
	}

	client.SLogger.Info("successfully disconnected from controller")
	return nil
}

type msgCipUnRegister struct {
	Service                CIPService
	CipPathSize            byte
	ClassType              cipClassSize
	Class                  byte
	InstanceType           cipInstanceSize
	Instance               byte
	Priority               byte
	TimeoutTicks           byte
	ConnectionSerialNumber uint16
	VendorID               uint16
	OriginatorSerialNumber uint32
	PathSize               uint8
	Reserved               byte // Always 0x00
}
