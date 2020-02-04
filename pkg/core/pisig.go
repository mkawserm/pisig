package core

import (
	"github.com/golang/glog"
	"github.com/mkawserm/pisig/pkg/cache"
	"github.com/mkawserm/pisig/pkg/context"
	"github.com/mkawserm/pisig/pkg/cors"
	"github.com/mkawserm/pisig/pkg/event"
	"github.com/mkawserm/pisig/pkg/message"
	"github.com/mkawserm/pisig/pkg/service"
	"github.com/mkawserm/pisig/pkg/settings"
	"net"
	"net/http"
)

type Pisig struct {
	mServerMux       *http.ServeMux
	mEPool           *EPool
	mPisigContext    *context.PisigContext
	mTopicDispatcher *TopicDispatcher

	mMiddlewareViewList []HTTPMiddlewareView
}

func (p *Pisig) CORSOptions() *cors.CORSOptions {
	return p.mPisigContext.GetCORSOptions()
}

func (p *Pisig) PisigContext() *context.PisigContext {
	return p.mPisigContext
}

func (p *Pisig) PisigMessage() message.PisigMessage {
	return p.mPisigContext.GetPisigMessage()
}

func (p *Pisig) PisigSettings() *settings.PisigSettings {
	return p.mPisigContext.GetPisigSettings()
}

func (p *Pisig) PisigStore() *cache.PisigStore {
	return p.mPisigContext.GetPisigStore()
}

func (p *Pisig) ProduceTopic(topic event.Topic) {
	p.mPisigContext.TopicProducerQueue <- topic
}

func (p *Pisig) AddConsumer(topicName string, service service.PisigService) bool {

	return false
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

func (p *Pisig) RunTopicDispatcher() {
	p.mTopicDispatcher.Run()
}

func (p *Pisig) RunHTTPServer() {
	if glog.V(1) {
		glog.Infof("Server starting...\n")
	}

	p.mEPool.Setup()
	go p.mEPool.RunMainEventLoop()
	p.runHTTPServer()

	if glog.V(1) {
		glog.Infof("Server exited gracefully.\n")
	}
}

func (p *Pisig) runHTTPServer() {
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

func (p *Pisig) AddWebSocketConnection(conn net.Conn) bool {
	if glog.V(3) {
		glog.Infof("Adding connection\n")
	}

	err := p.mEPool.AddConnection(conn)

	if err != nil {
		glog.Errorf("Failed to add connection - %v \n", err)
		return false
	}

	if glog.V(3) {
		glog.Infof("Connection added\n")
	}
	return true
}

func (p *Pisig) RemoveWebSocketConnection(conn net.Conn) bool {
	if glog.V(3) {
		glog.Infof("Removing connection\n")
	}

	//err := p.mEPool.RemoveConnection(conn)
	//
	//if err != nil {
	//	glog.Errorf("Failed to remove connection - %v \n",err)
	//	return false
	//}
	//
	//if glog.V(3) {
	//	glog.Infof("Connection removed \n")
	//}

	return true
}

func (p *Pisig) hookRemoveConnection(conn net.Conn) error {
	if glog.V(3) {
		glog.Infof("Removing connection\n")
	}

	// CLEAN UP RESOURCES

	err := p.mEPool.RemoveConnection(conn)

	if err != nil {
		glog.Errorf("Failed to remove connection - %v \n", err)
	}

	if glog.V(3) {
		glog.Infof("Connection removed\n")
	}

	return err
}

func NewPisig(args ...interface{}) *Pisig {
	if glog.V(3) {
		glog.Infof("Creating new Pisig instance\n")
	}
	pisig := &Pisig{}
	pisig.mPisigContext = nil
	pisig.mEPool = nil
	pisig.mServerMux = &http.ServeMux{}

	var pisigSettings *settings.PisigSettings
	var corsOptions *cors.CORSOptions
	var pisigContext *context.PisigContext
	var pisigMessage message.PisigMessage

	for _, val := range args {
		if glog.V(3) {
			glog.Infof("Pisig arg type: %T\n", val)
		}
		switch val.(type) {

		case *context.PisigContext:
			pisigContext = val.(*context.PisigContext)
		case *settings.PisigSettings:
			pisigSettings = val.(*settings.PisigSettings)
		case *cors.CORSOptions:
			corsOptions = val.(*cors.CORSOptions)
		case message.PisigMessage:
			pisigMessage = val.(message.PisigMessage)
		default:
			break
		}
	}

	if pisigContext != nil {
		pisig.mPisigContext = pisigContext
		pisig.mPisigContext.TopicProducerQueue = make(event.TopicQueue,
			pisigContext.GetPisigSettings().TopicQueueSize)
		pisig.mTopicDispatcher = NewTopicDispatcher(pisig, pisig.mPisigContext.TopicProducerQueue)

		ePool, err := NewEPool(pisig.mPisigContext, pisig.hookRemoveConnection)

		if err != nil {
			panic(err)
			return nil
		}
		pisig.mEPool = ePool

		if glog.V(3) {
			glog.Infof("New Pisig instance created")
		}
		return pisig
	}

	if pisigSettings != nil && corsOptions != nil && pisigMessage != nil {

		pisigContext := context.NewPisigContext()
		pisigContext.CORSOptions = corsOptions
		pisigContext.PisigSettings = pisigSettings

		pisigContext.PisigMessage = pisigMessage
		pisig.mPisigContext = pisigContext

		ePool, err := NewEPool(pisig.mPisigContext, pisig.hookRemoveConnection)

		if err != nil {
			panic(err)
			return nil
		}
		pisig.mEPool = ePool

		pisig.mPisigContext.TopicProducerQueue = make(event.TopicQueue,
			pisigContext.GetPisigSettings().TopicQueueSize)
		pisig.mTopicDispatcher = NewTopicDispatcher(pisig, pisig.mPisigContext.TopicProducerQueue)

		if glog.V(3) {
			glog.Infof("New Pisig instance created")
		}
		return pisig
	}

	glog.Errorln("Failed to create new Pisig instance.")
	return nil
}

func NewPisigSimple(corsOptions *cors.CORSOptions,
	pisigSettings *settings.PisigSettings,
	pisigResponse message.PisigMessage) *Pisig {
	return NewPisig(corsOptions, pisigSettings, pisigResponse)
}
