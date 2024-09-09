package gologix

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type serverConnection struct {
	ID   uint16
	OT   uint32
	TO   uint32
	RPI  time.Duration
	Path []byte
	Open bool
}

type serverConnectionManager struct {
	Connections []*serverConnection
	Lock        sync.RWMutex
	Logger      *slog.Logger
}

func (cm *serverConnectionManager) Init(logger *slog.Logger) {
	cm.Connections = make([]*serverConnection, 0, 32)
	cm.Logger = logger
}
func (cm *serverConnectionManager) Add(conn *serverConnection) {
	cm.Logger.Info("New Managed Connection.", "conn", *conn)
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	cm.Connections = append(cm.Connections, conn)
}
func (cm *serverConnectionManager) GetByID(ID uint16) (*serverConnection, error) {
	cm.Lock.RLock()
	defer cm.Lock.RUnlock()
	for _, conn := range cm.Connections {
		if conn.ID == ID {
			return conn, nil
		}
	}
	return nil, fmt.Errorf("couldn't find connection %v by ID", ID)
}
func (cm *serverConnectionManager) GetByOT(OT uint32) (*serverConnection, error) {
	cm.Lock.RLock()
	defer cm.Lock.RUnlock()
	for _, conn := range cm.Connections {
		if conn.OT == OT {
			return conn, nil
		}
	}
	return nil, fmt.Errorf("couldn't find connection %v by OT", OT)
}
func (cm *serverConnectionManager) GetByTO(TO uint32) (*serverConnection, error) {
	cm.Lock.RLock()
	defer cm.Lock.RUnlock()
	for _, conn := range cm.Connections {
		if conn.TO == TO {
			return conn, nil
		}
	}
	return nil, fmt.Errorf("couldn't find connection %v by TO", TO)
}
func (cm *serverConnectionManager) CloseByID(ID uint16) error {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	for i, conn := range cm.Connections {
		if conn.ID == ID {
			conn.Open = false
			if len(cm.Connections) == 1 {
				cm.Connections = make([]*serverConnection, 0, 32)
				return nil
			}
			cm.Connections[i] = cm.Connections[len(cm.Connections)-1]
			cm.Connections = cm.Connections[:len(cm.Connections)-1]
			return nil
		}
	}
	return fmt.Errorf("couldn't find connection %v by ID", ID)
}
func (cm *serverConnectionManager) CloseByOT(OT uint32) error {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	for i, conn := range cm.Connections {
		if conn.OT == OT {
			conn.Open = false
			if len(cm.Connections) == 1 {
				cm.Connections = make([]*serverConnection, 0, 32)
				return nil
			}
			cm.Connections[i] = cm.Connections[len(cm.Connections)-1]
			cm.Connections = cm.Connections[:len(cm.Connections)-1]
			return nil
		}
	}
	return fmt.Errorf("couldn't find connection %v by OT", OT)
}
func (cm *serverConnectionManager) CloseByTO(TO uint32) error {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	for i, conn := range cm.Connections {
		if conn.TO == TO {
			conn.Open = false
			if len(cm.Connections) == 1 {
				cm.Connections = make([]*serverConnection, 0, 32)
				return nil
			}
			cm.Connections[i] = cm.Connections[len(cm.Connections)-1]
			cm.Connections = cm.Connections[:len(cm.Connections)-1]
			return nil
		}
	}
	return fmt.Errorf("couldn't find connection %v by TO", TO)
}
