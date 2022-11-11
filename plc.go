package main

import "time"

type PLC struct {
	IPAddress     string
	ProcessorSlot int
	SocketTimeout time.Duration
	// Route
	conn Connection
}

func (plc *PLC) Read_Bytes(tag string, count int) []byte {
	return nil

}

type Connection struct {
	Size int // 508 is the default
}

func NewConnection() (conn Connection) {
	conn.Size = 508
	return
}
