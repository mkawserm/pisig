package message

type PisigMessage interface {
	InitAllStatusCode()
	GetStatusCode(codeNo int) []byte

	HTTP200() []byte
	HTTP404() []byte
	HTTP400() []byte
	HTTP500() []byte
}
