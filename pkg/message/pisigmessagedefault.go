package message

type PisigMessageDefault struct {
}

func (dpm *PisigMessageDefault) InitAllStatusCode() {

}

func (dpm *PisigMessageDefault) GetStatusCode(int) []byte {
	return []byte{}
}

func (dpm *PisigMessageDefault) HTTP404() []byte {
	return []byte(`{"data": null, "errors": [{"message": "PISIG404", "code": 404}]}`)
}

func (dpm *PisigMessageDefault) HTTP400() []byte {
	return []byte(`{"data": null, "errors": [{"message": "PISIG400", "code": 400}]}`)
}

func (dpm *PisigMessageDefault) HTTP500() []byte {
	return []byte(`{"data": null, "errors": [{"message": "PISIG500", "code": 500}]}`)
}
