package core

import "net/http"

type Pisig struct {
	Host string
	Port string

	EnableTLS bool   //read only
	CertFile  string //read only
	KeyFile   string //read only

	EventPoolQueueSize   int
	EventPoolWaitingTime int

	mServerMux *http.ServeMux
}
