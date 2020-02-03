package message

type PisigMessage interface {
	HTTP404() []byte
	HTTP400() []byte
	HTTP500() []byte
}
