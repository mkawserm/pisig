package core

import "sync"

type PisigContext struct {
	mIsLive bool

	mRWLock *sync.RWMutex
}
