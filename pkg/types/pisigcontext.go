package types

import "sync"

type PisigContext struct {
	mIsLive bool

	mRWLock *sync.RWMutex
}

func NewPisigContext() *PisigContext {
	return &PisigContext{
		mIsLive: false,
		mRWLock: &sync.RWMutex{},
	}
}
