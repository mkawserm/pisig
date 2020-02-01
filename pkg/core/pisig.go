package core

import (
	"github.com/golang/glog"
	"github.com/mkawserm/pisig/pkg/types"
	"net/http"
)

type Pisig struct {
	mServerMux *http.ServeMux

	mEventPool     *EventPool
	mCORSOptions   *types.CORSOptions
	mPisigContext  *types.PisigContext
	mPisigSettings *types.PisigSettings
}

func (p *Pisig) CORSOptions() *types.CORSOptions {
	return p.mCORSOptions
}

func (p *Pisig) PisigContext() *types.PisigContext {
	return p.mPisigContext
}

func (p *Pisig) PisigSettings() *types.PisigSettings {
	return p.mPisigSettings
}

func (p *Pisig) Run() {
	if glog.V(1) {
		glog.Infof("Server starting...\n")
	}

	p.mEventPool.Setup()
	go p.mEventPool.RunMainEventLoop()
	p.runServer()

	if glog.V(1) {
		glog.Infof("Server exited gracefully.\n")
	}
}

func (p *Pisig) runServer() {
	if p.mPisigSettings.EnableTLS {
		if glog.V(1) {
			glog.Infoln("Server is listening at: https://" + p.mPisigSettings.Host + ":" + p.mPisigSettings.Port)
		}
		err := http.ListenAndServeTLS(p.mPisigSettings.Host+":"+p.mPisigSettings.Port,
			p.mPisigSettings.CertFile,
			p.mPisigSettings.KeyFile, p.mServerMux)
		if err != nil {
			glog.Errorln("Server error: ", err)
		}
	} else {
		if glog.V(1) {
			glog.Infoln("Server is listening at: http://" + p.mPisigSettings.Host + ":" + p.mPisigSettings.Port)
		}
		err := http.ListenAndServe(p.mPisigSettings.Host+":"+p.mPisigSettings.Port, p.mServerMux)
		if err != nil {
			glog.Errorln("Server error: ", err)
		}
	}
}

func NewPisig(corsOptions *types.CORSOptions, pisigSettings *types.PisigSettings) *Pisig {

	serverMux := &http.ServeMux{}

	eventPool, err := NewEventPool(
		pisigSettings.EventPoolQueueSize,
		pisigSettings.EventPoolWaitingTime,
		nil,
		nil)

	if err != nil {
		panic(err)
		return nil
	}

	pisigContext := &types.PisigContext{}

	return &Pisig{
		mServerMux:     serverMux,
		mEventPool:     eventPool,
		mCORSOptions:   corsOptions,
		mPisigContext:  pisigContext,
		mPisigSettings: pisigSettings,
	}
}

func NewDefaultPisig() *Pisig {
	pisig := NewPisig(
		&types.CORSOptions{
			AllowAllOrigins:  true,
			AllowCredentials: true,
			AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
			AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization"},
		},
		types.NewDefaultPisigSettings(),
	)

	return pisig
}
