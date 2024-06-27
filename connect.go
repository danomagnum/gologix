package gologix

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"
)

const (
	connSizeLargeDefault    = 4000   // default large connection size
	connSizeStandardDefault = 511    // default small connection size
	connSizeStandardMax     = 511    // maximum size of connection for standard
	portDefault             = 44818  // default CIP port
	vendorIdDefault         = 0x9999 // default vendor id. Used to prevent vendor ID conflicts
	socketTimeoutDefault    = time.Second * 10
	rpiDefault              = time.Millisecond * 2500
)

// Connect to the PLC.
func (client *Client) Connect() error {
	if client.disconnecting {
		client.SLogger.Debug("waiting for client to finish disconnecting before connecting")
		for client.disconnecting {
			time.Sleep(time.Millisecond * 10)
		}
	}
	if client.connected || client.connecting {
		return nil
	}
	client.connecting = true
	defer func() { client.connecting = false }()
	if client.SLogger != nil {
		client.SLogger = client.SLogger.With(slog.String("controllerIp", client.Controller.IpAddress))
	}
	if client.ConnectionSize == 0 {
		client.ConnectionSize = connSizeLargeDefault
	}

	if client.Controller.Port == 0 {
		client.Controller.Port = portDefault
	}
	if client.Controller.VendorId == 0 {
		client.Controller.VendorId = vendorIdDefault
	}
	if client.SocketTimeout == 0 {
		client.SocketTimeout = socketTimeoutDefault
	}
	if client.RPI == 0 {
		client.RPI = rpiDefault
	}
	client.sequenceNumber.Add(uint32(time.Now().UnixMilli()))

	// default path is back plane -> slot 0
	var err error
	if client.Controller.Path == nil {
		client.Controller.Path, err = Serialize(CIPPort{PortNo: 1}, cipAddress(0))
		if err != nil {
			msg := "cannot setup default path"
			client.SLogger.Error(msg, slog.Any("err", err))
			return fmt.Errorf("%s: %w", msg, err)
		}
	}

	if client.ioi_cache == nil {
		client.ioi_cache = make(map[string]*tagIOI)
	}

	address := fmt.Sprintf("%s:%v", client.Controller.IpAddress, client.Controller.Port)
	client.conn, err = net.DialTimeout("tcp", address, client.SocketTimeout)
	if err != nil {
		msg := "cannot connect to controller"
		client.SLogger.Error(msg, slog.Any("err", err))
		return fmt.Errorf("%s: %w", msg, err)
	}

	err = client.registerSession()
	if err != nil {
		return err
	}

	if client.ConnectionSize > connSizeStandardMax {
		item, err := client.newForwardOpenLarge()
		if err != nil {
			return err
		}
		err = client.forwardOpen(item)
		if err != nil {
			client.SLogger.Warn("large forward open failed. falling back to standard forward open", slog.Any("err", err))
			client.ConnectionSize = connSizeStandardDefault
		}
	}
	if client.ConnectionSize <= connSizeStandardMax {
		item, err := client.newForwardOpenStandard()
		if err != nil {
			return err
		}
		err = client.forwardOpen(item)
		if err != nil {
			client.SLogger.Error("unable to open connection", slog.Any("err", err))
			return err
		}
	}
	client.connected = true

	if client.KeepAliveAutoStart {
		go client.KeepAlive()
	}
	return nil
}

func (client *Client) Connected() bool {
	return client.connected
}

func (client *Client) registerSession() error {
	reg_msg := msgCIPRegister{
		ProtocolVersion: 1,
		OptionFlag:      0,
	}

	header, _, err := client.send_recv_data(cipCommandRegisterSession, reg_msg)
	if err != nil {
		msg := "cannot get connect response"
		client.SLogger.Error(msg, slog.Any("err", err))
		return fmt.Errorf("%s: %w", msg, err)
	}
	client.SessionHandle = header.SessionHandle
	client.SLogger.Info("Session connected", slog.Any("sessionHandle", client.SessionHandle))
	return nil
}

func (client *Client) KeepAlive() {
	if !client.KeepAliveAutoStart || client.SocketTimeout == 0 {
		return
	}
	if client.keepAliveRunning {
		err := errors.New("keepalive already running")
		client.SLogger.Warn(err.Error())
	}
	client.SLogger.Debug("starting keep alive")
	client.cancel_keepalive = make(chan struct{})
	client.keepAliveRunning = true
	defer func() { client.keepAliveRunning = false }()

	originalProps, err := client.GetAttrList(CipObject_ControllerInfo, 1, client.KeepAliveProps...)
	if err != nil {
		client.Logger.Printf("keepalive prop list failed. %w", err)
		client.SLogger.Error(
			"initial keep alive property get failed",
			slog.Any("client.KeepAliveProps", client.KeepAliveProps),
			slog.Any("err", err),
		)
		return
	}

	err = client.ListAllTags(0)
	if err != nil {
		client.Logger.Printf("keepalive list tags failed. %w", err)
		client.SLogger.Error("keepalive list tags failed", slog.Any("err", err))
		return
	}

	t := time.NewTicker(client.KeepAliveFrequency)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if !client.connected {
				client.SLogger.Warn("keepalive failed. not connected")
				return
			}

			newProps, err := client.GetAttrList(CipObject_ControllerInfo, 1, client.KeepAliveProps...)
			if err != nil {
				client.Logger.Printf("keepalive failed. %w", err)
				client.SLogger.Error("keepalive failed", slog.Any("err", err))
				client.Disconnect()
				return
			}
			if newProps != originalProps {
				client.Logger.Printf("controller change detected. re-analyzing types.\n Was: %+v\n  Is: %+v", originalProps, newProps)
				client.SLogger.Info(
					"controller change detected. re-analyzing types",
					slog.Any("originalProps", originalProps),
					slog.Any("newProps", newProps),
				)
				err := client.ListAllTags(0)
				if err != nil {
					client.Logger.Printf("keepalive list tags failed. %v", err)
					return
				}
				originalProps = newProps
			}

		case <-client.cancel_keepalive:
			client.KeepAliveAutoStart = false
			return
		}
	}
}

type msgPreItemData struct {
	Handle  uint32
	Timeout uint16
}

type msgCIPMessageRouterResponse struct {
	Service   CIPService
	Reserved  byte      // always 0
	Status    CIPStatus // result status
	StatusLen byte      // additional result word count - can be zero
}

type msgEIPForwardClose struct {
	Service                CIPService
	PathSize               byte
	ClassType              byte
	Class                  byte
	InstanceType           byte
	Instance               byte
	Priority               byte
	TimeoutTicks           byte
	ConnectionSerialNumber uint16
	VendorID               uint16
	OriginatorSerialNumber uint32
	ConnPathSize           byte
}

// message for opening either a large or standard connection
// in this message T is for target and O is for originator so
// TO is target -> originator and OT is originator -> target
type cipForwardOpen[T uint16 | uint32] struct {
	Service      CIPService
	PathSize     byte // length in words
	ClassType    cipClassSize
	Class        byte
	InstanceType cipInstanceSize
	Instance     byte

	Priority               byte // 0x0A means normal multiplier (about 1 second)
	TimeoutTicks           byte // number of "priority" ticks (Ex: 0x0E = 14 * Priority = ~1 sec => ~ 14 seconds.)
	OTConnectionID         uint32
	TOConnectionID         uint32
	ConnectionSerialNumber uint16
	VendorID               uint16
	OriginatorSerialNumber uint32
	Multiplier             uint32
	OtRpi                  uint32
	OTNetworkConnParams    T // uint16 if standard, uint32 if large
	ToRpi                  uint32
	TONetworkConnParams    T // uint16 if standard, uint32 if large
	TransportTrigger       byte
	ConnPathSize           byte
}

func (client *Client) newForwardOpenLarge() (CIPItem, error) {
	item := CIPItem{Header: cipItemHeader{ID: cipItem_UnconnectedData}}
	if client.ConnectionSize == 0 {
		client.ConnectionSize = connSizeLargeDefault
	}
	if client.ConnectionSize <= connSizeStandardMax {
		client.SLogger.Info(
			"The size could be a standard connection",
			slog.Any("standardMaxSize", connSizeStandardMax),
			slog.Any("size", client.ConnectionSize),
		)
	}

	path, err := Serialize(
		client.Controller.Path,
		CipObject_MessageRouter,
		CIPInstance(1),
	)
	if err != nil {
		return item, fmt.Errorf("couldn't build path. %w", err)
	}

	client.ConnectionSerialNumber = uint16(client.sequenceNumber.Add(1))
	const (
		redundantOwner     uint32 = 0 // 0 = no-redundant, 1 = redundant
		connectionType     uint32 = 2 // 0 = null, 1 = multicast, 2 = point to point, 3 = reserved
		priority           uint32 = 0 // 0 = low, 1 = high, 2 = scheduled, 3 = urgent
		connectionSizeType uint32 = 1 // 1 = variable, 0 = fixed
	)
	connectionParameters := uint32(
		redundantOwner<<31 |
			connectionType<<29 |
			priority<<26 |
			connectionSizeType<<25 |
			uint32(client.ConnectionSize),
	)

	var msg cipForwardOpen[uint32]
	msg.Service = CIPService_LargeForwardOpen
	// this next section is the path
	msg.PathSize = 0x02
	msg.ClassType = cipClass_8bit
	msg.Class = byte(CipObject_ConnectionManager)
	msg.InstanceType = cipInstance_8bit
	msg.Instance = 0x01
	// this next section is the path
	msg.Priority = 0x0A
	msg.TimeoutTicks = 0x0E
	msg.OTConnectionID = client.sequenceNumber.Add(1) // pyLogix always uses 0x20000002
	msg.TOConnectionID = client.sequenceNumber.Add(1)
	msg.ConnectionSerialNumber = client.ConnectionSerialNumber
	msg.VendorID = client.VendorId
	msg.OriginatorSerialNumber = client.SerialNumber
	msg.Multiplier = 0x03
	msg.OtRpi = uint32(client.RPI / time.Microsecond)
	msg.OTNetworkConnParams = connectionParameters
	msg.ToRpi = uint32(client.RPI / time.Microsecond)
	msg.TONetworkConnParams = connectionParameters
	msg.TransportTrigger = 0xA3
	msg.ConnPathSize = byte(path.Len() / 2)

	item.Serialize(msg)
	item.Serialize(path.Bytes())
	return item, nil
}

type msgCIPRegister struct {
	ProtocolVersion uint16
	OptionFlag      uint16
}

func (client *Client) newForwardOpenStandard() (CIPItem, error) {
	if client.ConnectionSize == 0 {
		client.ConnectionSize = connSizeStandardDefault
	}
	if client.ConnectionSize > connSizeStandardMax {
		client.SLogger.Warn(
			"connection size too large. resetting to max size",
			slog.Any("oldConnectionSize", client.ConnectionSize),
			slog.Any("newConnectionSize", connSizeStandardMax),
		)
		client.ConnectionSize = connSizeStandardMax
	}
	item := CIPItem{Header: cipItemHeader{ID: cipItem_UnconnectedData}}

	path, err := Serialize(
		client.Controller.Path,
		CipObject_MessageRouter,
		CIPInstance(1),
	)
	if err != nil {
		return item, fmt.Errorf("couldn't build path. %w", err)
	}

	client.ConnectionSerialNumber = uint16(client.sequenceNumber.Add(1))
	const (
		redundantOwner     uint16 = 0 // 0 = no-redundant, 1 = redundant
		connectionType     uint16 = 2 // 0 = null, 1 = multicast, 2 = point to point, 3 = reserved
		priority           uint16 = 0 // 0 = low, 1 = high, 2 = scheduled, 3 = urgent
		connectionSizeType uint16 = 1 // 1 = variable, 0 = fixed
	)
	connectionParameters := uint16(
		redundantOwner<<15 |
			connectionType<<13 |
			priority<<10 |
			connectionSizeType<<9 |
			client.ConnectionSize,
	)

	var msg cipForwardOpen[uint16]
	msg.Service = CIPService_ForwardOpen
	// this next section is the path
	msg.PathSize = 0x02
	msg.ClassType = cipClass_8bit
	msg.Class = byte(CipObject_ConnectionManager)
	msg.InstanceType = cipInstance_8bit
	msg.Instance = 0x01
	// end of path
	msg.Priority = 0x07
	msg.TimeoutTicks = 0xE9
	msg.OTConnectionID = client.sequenceNumber.Add(1)
	msg.TOConnectionID = client.sequenceNumber.Add(1)
	msg.ConnectionSerialNumber = client.ConnectionSerialNumber
	msg.VendorID = client.VendorId
	msg.OriginatorSerialNumber = client.SerialNumber
	msg.Multiplier = 0x00
	msg.OtRpi = uint32(client.RPI / time.Microsecond)
	msg.OTNetworkConnParams = connectionParameters
	msg.ToRpi = uint32(client.RPI / time.Microsecond)
	msg.TONetworkConnParams = connectionParameters
	msg.TransportTrigger = 0xA3
	msg.ConnPathSize = byte(path.Len() / 2)
	item.Serialize(msg)
	item.Serialize(path.Bytes())

	return item, nil
}

func (client *Client) forwardOpen(forwardOpenMsg CIPItem) error {
	reqItems := make([]CIPItem, 2)
	reqItems[0] = CIPItem{Header: cipItemHeader{ID: cipItem_Null}}
	reqItems[1] = forwardOpenMsg
	itemData, err := serializeItems(reqItems)
	if err != nil {
		client.SLogger.Error("error serializing items", slog.Any("err", err))
		return err
	}

	header, data, err := client.send_recv_data(cipCommandSendRRData, itemData)
	if err != nil {
		client.SLogger.Error("error sending data", slog.Any("err", err))
		return err
	}

	items, err := client.parseResponse(&header, data)
	if err != nil {
		client.SLogger.Error("error parsing response", slog.Any("err", err))
		return err
	}

	respContent := msgCipForwardOpenReply{}
	err = items[1].DeSerialize(&respContent)
	if err != nil {
		client.SLogger.Error("error deserializing forward open response", slog.Any("err", err))
		return fmt.Errorf("error deserializing forward open response content. %w", err)
	}

	client.OTNetworkConnectionID = respContent.OtNetworkConnectionId

	client.SLogger.Info(
		"successfully opened connection",
		slog.Any("ConnectionSize", uint32(client.ConnectionSize)),
		slog.Any("OTNetworkConnectionId", respContent.OtNetworkConnectionId),
	)

	return nil
}

type msgCipForwardOpenReply struct {
	OtNetworkConnectionId  uint32
	TOConnectionId         uint32
	ConnectionSerialNumber uint16
	OriginatorVendorId     uint16
	OriginatorSerialNumber uint32
	OTApiNs                uint32
	TOApiNs                uint32
	ApplicationReply       uint8
	Reserved               uint8
}

type msgEIPForwardOpen_Standard_Reply struct {
	Service                CIPService
	Reserved               byte
	Status                 CIPStatus
	StatusLen              byte
	OTConnectionID         uint32
	TOConnectionID         uint32
	ConnectionSerialNumber uint16
	VendorID               uint16
	OriginatorSerialNumber uint32
	OTApi                  uint32
	TOApi                  uint32
	ReplySize              byte
	Reserved2              byte
}

func (client *Client) checkConnection() error {
	if !client.connected {
		if client.AutoConnect {
			err := client.Connect()
			if err != nil {
				return fmt.Errorf("not connected and connect attempt failed: %w", err)
			}
		} else {
			return fmt.Errorf("not connected and AutoConnect not enabled")
		}
	}
	return nil
}

func (client *Client) parseResponse(header *eipHeader, data *bytes.Buffer) ([]CIPItem, error) {
	if header.Status != 0 {
		return nil, fmt.Errorf("forward open failed. status: %v", header.Status)
	}

	preItem := msgPreItemData{}
	err := binary.Read(data, binary.LittleEndian, &preItem)
	if err != nil {
		return nil, fmt.Errorf("problem reading items header from forward open request. %w", err)
	}

	items, err := readItems(data)
	if err != nil {
		return nil, fmt.Errorf("problem reading items from forward open request. %w", err)
	}

	respHeader := msgCIPMessageRouterResponse{}
	err = items[1].DeSerialize(&respHeader)
	if err != nil {
		return nil, fmt.Errorf("error deserializing forward open response header. %w", err)
	}

	extended_status := make([]byte, respHeader.StatusLen*2)
	if respHeader.StatusLen != 0 {
		err = items[1].DeSerialize(&extended_status)
		if err != nil {
			return nil, fmt.Errorf("error deserializing forward open response header extended status. %w", err)
		}
	}
	if respHeader.Status != 0 {
		errMsg := "bad status on response"
		client.SLogger.Error(errMsg,
			slog.Any("status", respHeader.Status),
			slog.String("statusDesc", respHeader.Status.String()),
		)
		return nil, fmt.Errorf("%s status: %v", errMsg, respHeader.Status)
	}
	return items, nil
}
