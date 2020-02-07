package message

type PisigMessage interface {
	Init()
	Get(codeNo int) []byte
	Set(codeNo int, data []byte) bool

	HTTP200() []byte
	HTTP404() []byte
	HTTP400() []byte
	HTTP500() []byte
}
