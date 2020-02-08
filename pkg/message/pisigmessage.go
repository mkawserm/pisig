package message

type PisigMessageCode int

type PisigMessage interface {
	Init()
	Get(codeNo PisigMessageCode) []byte

	HTTP200() []byte
	HTTP404() []byte
	HTTP400() []byte
	HTTP500() []byte
}
