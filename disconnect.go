package gologix

import (
	"fmt"
	"log/slog"
	"time"
)

// You will want to defer this after a successful Connect() to make sure you free up the controller resources
// to disconnect we send two items - a null item and an unconnected data item for the unregister service
func (client *Client) Disconnect() error {
	if client.connecting {
		client.SLogger.Debug("waiting for client to finish connecting before disconnecting")
		for client.connecting {
			time.Sleep(time.Millisecond * 10)
		}
	}
	if !client.connected || client.disconnecting {
		return nil
	}
	client.disconnecting = true
	defer func() { client.disconnecting = false }()
	client.connected = false
	var err error
	client.SLogger.Info("starting disconnection")

	if client.keepAliveRunning {
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
		} else {
			_, err = client.parseResponse(&header, data)
			if err != nil {
				client.SLogger.Error(
					"error parsing disconnect response",
					slog.Any("err", err),
				)
			}
		}
	}

	err = client.conn.Close()
	if err != nil {
		client.SLogger.Error("error closing connection", slog.Any("err", err))
	}

	client.SLogger.Info("successfully disconnected from controller")
	return nil
}

// Cancels keepalive if KeepAliveAutoStart is false. Use force to cancel keepalive regardless.
// If forced, the keepalive will not resume unless the client is reconnected or KeepAlive is triggered
func (client *Client) KeepAliveCancel(force bool) error {
	if client.KeepAliveAutoStart && !force {
		return fmt.Errorf("unable to cancel keepalive due to AutoKeepAlive == true")
	}
	close(client.cancel_keepalive)
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
