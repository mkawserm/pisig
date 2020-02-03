package message

type DefaultPisigMessage struct {
}

func (dpm *DefaultPisigMessage) InitAllStatusCode() {

}

func (dpm *DefaultPisigMessage) GetStatusCode(int) []byte {
	return []byte{}
}

func (dpm *DefaultPisigMessage) HTTP404() []byte {
	return []byte(`{"data": null, "errors": [{"message": "PISIG404", "code": 404}]}`)
}

func (dpm *DefaultPisigMessage) HTTP400() []byte {
	return []byte(`{"data": null, "errors": [{"message": "PISIG400", "code": 400}]}`)
}

func (dpm *DefaultPisigMessage) HTTP500() []byte {
	return []byte(`{"data": null, "errors": [{"message": "PISIG500", "code": 500}]}`)
}
