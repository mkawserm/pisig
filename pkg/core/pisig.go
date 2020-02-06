package core

import (
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/glog"
	"github.com/mkawserm/pisig/pkg/cache"
	"github.com/mkawserm/pisig/pkg/event"
	"github.com/mkawserm/pisig/pkg/message"
	"github.com/mkawserm/pisig/pkg/settings"
	"github.com/mkawserm/pisig/pkg/storage"
	"net"
	"net/http"
)

type Pisig struct {
	mServerMux       *http.ServeMux
	mEPool           *EPool
	mPisigContext    *PisigContext
	mTopicDispatcher *TopicDispatcher

	mMiddlewareViewList []HTTPMiddlewareView
}

func (p *Pisig) CORSOptions() *CORSOptions {
	return p.mPisigContext.GetCORSOptions()
}

func (p *Pisig) PisigContext() *PisigContext {
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

func (p *Pisig) SendMessageToAll(message []byte) {
	for _, conn := range p.mEPool.GetConnectionSlice() {
		err := wsutil.WriteServerText(conn, message)
		if err != nil {
			glog.Infof("%v\n", err)
		}
	}
}

func (p *Pisig) ProduceTopic(topic event.Topic) {
	p.mPisigContext.TopicProducerQueue <- topic
}

func (p *Pisig) GetTopicListenerList(topicName string) []interface{} {
	return p.mPisigContext.GetPisigServiceRegistry().GetTopicListenerList(topicName)
}

func (p *Pisig) AddService(topicNameList []string, pisigService PisigService) bool {
	added, err := p.mPisigContext.PisigServiceRegistry.AddService(pisigService)
	if err != nil {
		glog.Errorf("Error: %v", err)
		return false
	}

	if added {
		for _, topicName := range topicNameList {
			p.mPisigContext.PisigServiceRegistry.AddTopicListener(topicName, pisigService)
		}
		pisigService.SetPisig(p)
		return true
	}

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
	var corsOptions *CORSOptions
	var pisigMessage message.PisigMessage
	var onlineUserStore storage.OnlineUserStore
	var pisigContext *PisigContext

	for _, val := range args {
		if glog.V(3) {
			glog.Infof("Pisig arg type: %T\n", val)
		}
		switch val.(type) {

		case *PisigContext:
			pisigContext = val.(*PisigContext)
			break
		case *settings.PisigSettings:
			pisigSettings = val.(*settings.PisigSettings)
		case *CORSOptions:
			corsOptions = val.(*CORSOptions)
		case message.PisigMessage:
			pisigMessage = val.(message.PisigMessage)
		case storage.OnlineUserStore:
			onlineUserStore = val.(storage.OnlineUserStore)
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

		pisigContext := NewPisigContext()
		pisigContext.CORSOptions = corsOptions
		pisigContext.PisigSettings = pisigSettings
		pisigContext.OnlineUserStore = onlineUserStore

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

func NewPisigSimple(corsOptions *CORSOptions,
	pisigSettings *settings.PisigSettings,
	pisigResponse message.PisigMessage,
	onlineUserStore storage.OnlineUserStore) *Pisig {
	return NewPisig(corsOptions, pisigSettings, pisigResponse, onlineUserStore)
}
