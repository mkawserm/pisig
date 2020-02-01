package types

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
	value, ok := ps.mSettings[key]
	return value, ok
}

func (ps *PisigSettings) Set(key string, value interface{}) {
	ps.mSettings[key] = value
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
