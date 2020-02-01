package core

import "net/http"

type Pisig struct {
	mHost string
	mPort string

	mEnableTLS bool
	mCertFile  string
	mKeyFile   string

	mEventPoolQueueSize   int
	mEventPoolWaitingTime int

	mEventPool    *EventPool
	mServerMux    *http.ServeMux
	mPisigContext *PisigContext
}
