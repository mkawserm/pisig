package event

import "reflect"

type Protocol int

const (
	Unknown   Protocol = 0
	TCP       Protocol = 1
	UDP       Protocol = 2
	gRPC      Protocol = 3
	WebSocket Protocol = 4
	HTTP      Protocol = 5
)

func (protocol Protocol) String() string {
	names := [...]string{
		"Unknown",
		"TCP",
		"UDP",
		"gRPC",
		"WebSocket",
		"HTTP",
	}

	if protocol < TCP || protocol > HTTP {
		return "Unknown"
	}
	return names[protocol]
}

type Topic struct {
	Name string
	Key  []byte
	Data interface{}
}

func (t *Topic) DataType() string {
	return reflect.TypeOf(t.Data).String()
}
