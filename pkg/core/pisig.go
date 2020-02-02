package core

import (
	"github.com/golang/glog"
	"github.com/mkawserm/pisig/pkg/variant"
	"net/http"
)

type Pisig struct {
	mServerMux    *http.ServeMux
	mEventPool    *EventPool
	mPisigContext *variant.PisigContext

	mMiddlewareViewList []HTTPMiddlewareView
}

func (p *Pisig) CORSOptions() *variant.CORSOptions {
	return p.mPisigContext.GetCORSOptions()
}

func (p *Pisig) PisigContext() *variant.PisigContext {
	return p.mPisigContext
}

func (p *Pisig) PisigSettings() *variant.PisigSettings {
	return p.mPisigContext.GetPisigSettings()
}

func (p *Pisig) PisigStore() *variant.PisigStore {
	return p.mPisigContext.GetPisigStore()
}

func (p *Pisig) AddView(urlPattern string, view HTTPView) {
	p.mServerMux.HandleFunc(urlPattern, view.Process(p))
}

func (p *Pisig) AddMiddlewareView(middlewareView HTTPMiddlewareView) {
	p.mMiddlewareViewList = append(p.mMiddlewareViewList, middlewareView)
}

func (p *Pisig) MiddlewareViewList() []HTTPMiddlewareView {
	return p.mMiddlewareViewList
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
	if p.PisigSettings().EnableTLS {
		if glog.V(1) {
			glog.Infoln("Server is listening at: https://" + p.PisigSettings().Host + ":" + p.PisigSettings().Port)
		}
		err := http.ListenAndServeTLS(p.PisigSettings().Host+":"+p.PisigSettings().Port,
			p.PisigSettings().CertFile,
			p.PisigSettings().KeyFile, p.mServerMux)
		if err != nil {
			glog.Errorln("Server error: ", err)
		}
	} else {
		if glog.V(1) {
			glog.Infoln("Server is listening at: http://" + p.PisigSettings().Host + ":" + p.PisigSettings().Port)
		}
		err := http.ListenAndServe(p.PisigSettings().Host+":"+p.PisigSettings().Port, p.mServerMux)
		if err != nil {
			glog.Errorln("Server error: ", err)
		}
	}
}

func NewPisig(args ...interface{}) *Pisig {
	pisig := &Pisig{}
	pisig.mPisigContext = nil
	pisig.mEventPool = nil
	pisig.mServerMux = &http.ServeMux{}

	var pisigSettings *variant.PisigSettings
	var corsOptions *variant.CORSOptions
	var pisigContext *variant.PisigContext

	for _, val := range args {
		switch val.(type) {

		case *variant.PisigContext:
			pisigContext = val.(*variant.PisigContext)
		case *variant.PisigSettings:
			pisigSettings = val.(*variant.PisigSettings)
		case *variant.CORSOptions:
			corsOptions = val.(*variant.CORSOptions)
		default:
			break
		}
	}

	if pisigContext != nil {
		pisig.mPisigContext = pisigContext
		eventPool, err := NewEventPool(
			pisig.PisigSettings().EventPoolQueueSize,
			pisig.PisigSettings().EventPoolWaitingTime,
			nil,
			nil,
		)

		if err != nil {
			panic(err)
			return nil
		}
		pisig.mEventPool = eventPool
		return pisig
	}

	if pisigSettings != nil && corsOptions != nil {

		eventPool, err := NewEventPool(
			pisigSettings.EventPoolQueueSize,
			pisigSettings.EventPoolWaitingTime,
			nil,
			nil)

		if err != nil {
			panic(err)
			return nil
		}

		pisigContext := variant.NewPisigContext()
		pisigContext.CORSOptions = corsOptions
		pisigContext.PisigSettings = pisigSettings

		pisig.mEventPool = eventPool
		pisig.mPisigContext = pisigContext
		return pisig
	}

	glog.Errorln("Failed to create new pisig instance.")
	return nil
}

func NewPisigSimple(corsOptions *variant.CORSOptions, pisigSettings *variant.PisigSettings) *Pisig {
	return NewPisig(corsOptions, pisigSettings)
}

func NewDefaultPisig() *Pisig {
	pisig := NewPisigSimple(
		&variant.CORSOptions{
			AllowAllOrigins:  true,
			AllowCredentials: true,
			AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
			AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization"},
		},
		variant.NewDefaultPisigSettings(),
	)

	return pisig
}
