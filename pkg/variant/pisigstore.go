package variant

import (
	"strings"
	"sync"
)

type PisigStore struct {
	mStore  map[string]interface{}
	mRWLock *sync.RWMutex
}

func (ps *PisigStore) Get(key string) (interface{}, bool) {
	ps.mRWLock.RLock()
	defer ps.mRWLock.RUnlock()
	value, ok := ps.mStore[strings.ToLower(key)]
	return value, ok
}

func (ps *PisigStore) Set(key string, value interface{}) {
	ps.mRWLock.Lock()
	defer ps.mRWLock.Unlock()
	ps.mStore[strings.ToLower(key)] = value
}

func (ps *PisigStore) IsSet(key string) bool {
	_, ok := ps.Get(key)
	return ok
}

func (ps *PisigStore) GetByte(key string) byte {
	val, ok := ps.Get(key)
	if ok {
		return val.(byte)
	}
	return byte(0)
}

func (ps *PisigStore) GetInt(key string) int {
	val, ok := ps.Get(key)
	if ok {
		return val.(int)
	}
	return 0
}

func (ps *PisigStore) GetInt64(key string) int64 {
	val, ok := ps.Get(key)
	if ok {
		return val.(int64)
	}
	return 0
}

func (ps *PisigStore) GetUInt(key string) uint {
	val, ok := ps.Get(key)
	if ok {
		return val.(uint)
	}
	return 0
}

func (ps *PisigStore) GetUInt64(key string) uint64 {
	val, ok := ps.Get(key)
	if ok {
		return val.(uint64)
	}
	return 0
}

func (ps *PisigStore) GetString(key string) string {
	val, ok := ps.Get(key)
	if ok {
		return val.(string)
	}
	return ""
}

func (ps *PisigStore) GetBool(key string) bool {
	val, ok := ps.Get(key)
	if ok {
		return val.(bool)
	}
	return false
}

func (ps *PisigStore) GetByteSlice(key string) []byte {
	val, ok := ps.Get(key)
	if ok {
		return val.([]byte)
	}
	return nil
}

func (ps *PisigStore) GetMap(key string) map[string]interface{} {
	val, ok := ps.Get(key)
	if ok {
		return val.(map[string]interface{})
	}
	return nil
}

func (ps *PisigStore) GetStringMap(key string) map[string]string {
	val, ok := ps.Get(key)
	if ok {
		return val.(map[string]string)
	}
	return nil
}

func (ps *PisigStore) GetByteSliceMap(key string) map[string][]byte {
	val, ok := ps.Get(key)
	if ok {
		return val.(map[string][]byte)
	}
	return nil
}

// Create new PisigStore
func NewPisigStore() *PisigStore {
	return &PisigStore{
		mStore:  make(map[string]interface{}),
		mRWLock: &sync.RWMutex{},
	}
}
