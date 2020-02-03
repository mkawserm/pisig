package context

import (
	"github.com/mkawserm/pisig/pkg/cache"
	"github.com/mkawserm/pisig/pkg/cors"
	"github.com/mkawserm/pisig/pkg/settings"
)

type PisigContext struct {
	PisigStore    *cache.PisigStore
	CORSOptions   *cors.CORSOptions
	PisigSettings *settings.PisigSettings
}

func (pc *PisigContext) GetCORSOptions() *cors.CORSOptions {
	return pc.CORSOptions
}

func (pc *PisigContext) GetPisigSettings() *settings.PisigSettings {
	return pc.PisigSettings
}

func (pc *PisigContext) GetPisigStore() *cache.PisigStore {
	return pc.PisigStore
}

// Create new PisigContext
func NewPisigContext() *PisigContext {
	return &PisigContext{PisigStore: cache.NewPisigStore()}
}
