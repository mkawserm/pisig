package event

const HTTPEventString = "HTTPEvent"

type HTTPEvent struct {
	Headers map[string]string
	Body    []byte
}
