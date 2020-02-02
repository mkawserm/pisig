package types

import (
	"sync"
)

type PisigContext struct {
	CORSOptions   *CORSOptions
	PisigSettings *PisigSettings

	mContextStore map[string]interface{}
	mRWLock       *sync.RWMutex
}

func (pc *PisigContext) GetFromContextStore(key string) (interface{}, bool) {
	pc.mRWLock.RLock()
	defer pc.mRWLock.RUnlock()
	value, ok := pc.mContextStore[key]
	return value, ok
}

func (pc *PisigContext) AddToContextStore(key string, value interface{}) {
	pc.mRWLock.Lock()
	defer pc.mRWLock.Unlock()
	pc.mContextStore[key] = value
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
