package event

const WebSocketEvent = "WebSocketEvent"

type WebSocketEventData struct {
	SocketId int
	OpCode   byte
	Message  []byte
}

const WebSocketConnEvent = "WebSocketConnEvent"

type WebSocketConnEventData struct {
	SocketId int
}
