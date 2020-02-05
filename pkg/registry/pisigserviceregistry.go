package registry

import (
	"sync"
)

type PisigServiceRegistry struct {
	mStore  map[string]interface{}
	mRWLock *sync.RWMutex
}

func NewPisigServiceRegistry() *PisigServiceRegistry {
	return &PisigServiceRegistry{
		mStore:  make(map[string]interface{}),
		mRWLock: &sync.RWMutex{},
	}
}
