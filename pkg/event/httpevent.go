package event

type HTTPEvent struct {
	Headers map[string]string
	Body    []byte
}
