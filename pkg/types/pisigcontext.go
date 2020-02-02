package types

import (
	"sync"
)

type PisigContext struct {
	IsLive        bool
	CORSOptions   *CORSOptions
	PisigSettings *PisigSettings

	mRWLock *sync.RWMutex
}

func (pc *PisigContext) GetCORSOptions() *CORSOptions {
	return pc.CORSOptions
}

func (pc *PisigContext) GetPisigSettings() *PisigSettings {
	return pc.PisigSettings
}

func NewPisigContext() *PisigContext {
	return &PisigContext{
		mRWLock: &sync.RWMutex{},
	}
}
