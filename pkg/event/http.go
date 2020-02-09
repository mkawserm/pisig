package event

const HTTPEvent = "HTTPEvent"

type HTTPEventData struct {
	Headers map[string]string
	Body    []byte
}
