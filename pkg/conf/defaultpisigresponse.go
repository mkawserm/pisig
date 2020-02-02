package conf

type DefaultPisigResponse struct {
}

func (dpr *DefaultPisigResponse) HTTP404() []byte {
	return []byte(`{"data": null, "error": [{"message": "PISIG404", "code": 404}]}`)
}

func (dpr *DefaultPisigResponse) HTTP400() []byte {
	return []byte(`{"data": null, "error": [{"message": "PISIG400", "code": 400}]}`)
}

func (dpr *DefaultPisigResponse) HTTP500() []byte {
	return []byte(`{"data": null, "error": [{"message": "PISIG500", "code": 500}]}`)
}
