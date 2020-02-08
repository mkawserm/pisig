package message

type PisigMessageDefault struct {
}

func (dpm *PisigMessageDefault) Init() {

}

func (dpm *PisigMessageDefault) Get(PisigMessageCode) []byte {
	return []byte{}
}

func (dpm *PisigMessageDefault) HTTP200() []byte {
	return []byte(`{"data": null, "errors": [{"message": "PISIG200", "code": 200}]}`)
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
