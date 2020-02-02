package types

import "strings"

type PisigSettings struct {
	Host string
	Port string

	EnableTLS bool
	CertFile  string
	KeyFile   string

	EventPoolQueueSize   int
	EventPoolWaitingTime int

	mSettings map[string]interface{}
}

func (ps *PisigSettings) Get(key string) (interface{}, bool) {
	value, ok := ps.mSettings[strings.ToLower(key)]
	return value, ok
}

func (ps *PisigSettings) Set(key string, value interface{}) {
	ps.mSettings[strings.ToLower(key)] = value
}

func (ps *PisigSettings) IsSet(key string) bool {
	_, ok := ps.Get(key)
	return ok
}

func (ps *PisigSettings) GetByte(key string) byte {
	val, ok := ps.Get(key)
	if ok {
		return val.(byte)
	}
	return byte(0)
}

func (ps *PisigSettings) GetInt(key string) int {
	val, ok := ps.Get(key)
	if ok {
		return val.(int)
	}
	return 0
}

func (ps *PisigSettings) GetInt64(key string) int64 {
	val, ok := ps.Get(key)
	if ok {
		return val.(int64)
	}
	return 0
}

func (ps *PisigSettings) GetUInt(key string) uint {
	val, ok := ps.Get(key)
	if ok {
		return val.(uint)
	}
	return 0
}

func (ps *PisigSettings) GetUInt64(key string) uint64 {
	val, ok := ps.Get(key)
	if ok {
		return val.(uint64)
	}
	return 0
}

func (ps *PisigSettings) GetString(key string) string {
	val, ok := ps.Get(key)
	if ok {
		return val.(string)
	}
	return ""
}

func (ps *PisigSettings) GetBool(key string) bool {
	val, ok := ps.Get(key)
	if ok {
		return val.(bool)
	}
	return false
}

func (ps *PisigSettings) GetByteSlice(key string) []byte {
	val, ok := ps.Get(key)
	if ok {
		return val.([]byte)
	}
	return nil
}

func (ps *PisigSettings) GetMap(key string) map[string]interface{} {
	val, ok := ps.Get(key)
	if ok {
		return val.(map[string]interface{})
	}
	return nil
}

func (ps *PisigSettings) GetStringMap(key string) map[string]string {
	val, ok := ps.Get(key)
	if ok {
		return val.(map[string]string)
	}
	return nil
}

func (ps *PisigSettings) GetByteSliceMap(key string) map[string][]byte {
	val, ok := ps.Get(key)
	if ok {
		return val.(map[string][]byte)
	}
	return nil
}

// Create new PisigSettings
func NewPisigSettings(host string,
	port string,
	enableTLS bool,
	certFile string,
	keyFile string,
	eventPoolQueueSize int,
	eventPoolWaitTime int) *PisigSettings {

	return &PisigSettings{
		Host:                 host,
		Port:                 port,
		EnableTLS:            enableTLS,
		CertFile:             certFile,
		KeyFile:              keyFile,
		EventPoolQueueSize:   eventPoolQueueSize,
		EventPoolWaitingTime: eventPoolWaitTime,
		mSettings:            make(map[string]interface{}),
	}
}

// Create new default PisigSettings
func NewDefaultPisigSettings() *PisigSettings {
	return NewPisigSettings("0.0.0.0",
		"8080",
		false,
		"",
		"",
		100,
		100)
}
