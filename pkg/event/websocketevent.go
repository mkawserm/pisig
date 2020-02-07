package event

import "net"

const WebSocketEventString = "WebSocketEvent"

type WebSocketEvent struct {
	Conn    net.Conn
	OpCode  byte
	Message []byte
}
