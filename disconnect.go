package gologix

import (
	"fmt"
	"log/slog"
)

func (client *Client) startDisconnect() error {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	switch client.connStatus {
	case connectionStatusDisconnected:
		return fmt.Errorf("client is already disconnected")
	case connectionStatusConnecting:
		return fmt.Errorf("client is still connecting, cannot disconnect")
	case connectionStatusConnected:
		// continue to disconnect
		client.connStatus = connectionStatusDisconnecting
		return nil
	case connectionStatusDisconnecting:
		return fmt.Errorf("client is already disconnecting")
	default:
		// Don't know what to do with this status - but continuing the disconnection is probably the best option
		client.connStatus = connectionStatusDisconnecting
		return nil
	}
}

// Disconnect gracefully closes the CIP connection to the PLC and releases controller resources.
//
// Always call Disconnect() after a successful Connect() to ensure proper cleanup.
// It's recommended to use defer for this:
//
//	client := gologix.NewClient("192.168.1.100")
//	err := client.Connect()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Disconnect()  // Ensures cleanup even if errors occur
//
// If you don't call Disconnect(), the PLC connection may remain allocated until
// it times out (typically about 2 minutes), which can prevent immediate reconnection
// and may consume PLC connection resources.
//
// Disconnect can be called multiple times safely - subsequent calls after the first
// successful disconnect will return an error but won't cause issues.
//
// The function attempts to send a proper disconnection message to the PLC, but will
// force the local connection closed even if the PLC communication fails, ensuring
// local resources are always cleaned up.
//
// Returns an error if already disconnected or if there are issues during the
// disconnection sequence, but the connection will be closed regardless.
func (client *Client) Disconnect() error {
	err := client.startDisconnect()
	if err != nil {
		return err
	}
	client.Logger.Info("starting disconnection")

	// No matter what happens, when we finish here, we'll consider the connection disconnected.
	defer func() {
		client.mutex.Lock()
		client.connStatus = connectionStatusDisconnected
		client.mutex.Unlock()
	}()

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
		client.Logger.Error("Error serializing path", slog.Any("err", err))
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

	err = items[1].Serialize(msg)
	if err != nil {
		return fmt.Errorf("error serializing disconnect msg: %w", err)
	}
	err = items[1].Serialize(path)
	if err != nil {
		return fmt.Errorf("error serializing disconnect path: %w", err)
	}

	itemData, err := serializeItems(items)
	if err != nil {
		client.Logger.Error(
			"unable to serialize itemData. Forcing connection closed",
			slog.Any("err", err),
		)
	} else {
		header, data, err := client.send_recv_data(cipCommandSendRRData, itemData)
		if err != nil {
			client.Logger.Error(
				"error sending disconnect request",
				slog.Any("err", err),
			)
		} else {
			_, err = client.parseResponse(&header, data)
			if err != nil {
				client.Logger.Error(
					"error parsing disconnect response",
					slog.Any("err", err),
				)
			}
		}
	}

	err = client.conn.Close()
	if err != nil {
		client.Logger.Error("error closing connection", slog.Any("err", err))
	}

	client.Logger.Info("successfully disconnected from controller")
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
