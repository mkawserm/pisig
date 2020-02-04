package event

import "net"

type WebSocketEvent struct {
	Conn    net.Conn
	OpCode  byte
	Message []byte
}
