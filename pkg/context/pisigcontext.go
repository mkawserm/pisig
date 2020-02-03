package context

import (
	"github.com/mkawserm/pisig/pkg/cors"
	"github.com/mkawserm/pisig/pkg/settings"
	"github.com/mkawserm/pisig/pkg/variant"
)

type PisigContext struct {
	PisigStore    *variant.PisigStore
	CORSOptions   *cors.CORSOptions
	PisigSettings *settings.PisigSettings
}

func (pc *PisigContext) GetCORSOptions() *cors.CORSOptions {
	return pc.CORSOptions
}

func (pc *PisigContext) GetPisigSettings() *settings.PisigSettings {
	return pc.PisigSettings
}

func (pc *PisigContext) GetPisigStore() *variant.PisigStore {
	return pc.PisigStore
}

// Create new PisigContext
func NewPisigContext() *PisigContext {
	return &PisigContext{PisigStore: variant.NewPisigStore()}
}
