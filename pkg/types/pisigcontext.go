package types

type PisigContext struct {
	CORSOptions   *CORSOptions
	PisigSettings *PisigSettings
	PisigStore    *PisigStore
}

func (pc *PisigContext) GetCORSOptions() *CORSOptions {
	return pc.CORSOptions
}

func (pc *PisigContext) GetPisigSettings() *PisigSettings {
	return pc.PisigSettings
}

func (pc *PisigContext) GetPisigStore() *PisigStore {
	return pc.PisigStore
}

// Create new PisigContext
func NewPisigContext() *PisigContext {
	return &PisigContext{PisigStore: NewPisigStore()}
}
