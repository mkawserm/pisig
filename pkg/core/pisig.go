package core

import (
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/glog"
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

func (p *Pisig) PisigContext() *PisigContext {
	return p.mPisigContext
}

// Send message to all active websocket connection
func (p *Pisig) SendMessageToAll(message []byte) {
	for _, conn := range p.mEPool.GetConnectionSlice() {
		err := wsutil.WriteServerText(conn, message)
		if err != nil {
			glog.Infof("%v\n", err)
		}
	}
}

// Send message to the specific user active websocket connection
func (p *Pisig) SendMessageToUser(uniqueId string, message []byte) {
	if p.mPisigContext.OnlineUserStore == nil {
		return
	}

	socketIdList := p.mPisigContext.OnlineUserStore.GetSocketIdListFromUniqueId(uniqueId)

	for _, socketId := range socketIdList {
		conn, ok := p.mEPool.GetConnection(socketId)
		if conn == nil || !ok {
			continue
		}
		err := wsutil.WriteServerText(conn, message)
		if err != nil {
			glog.Infof("%v\n", err)
		}
	}
}

// Send message to the specific group of active websocket connection
func (p *Pisig) SendMessageToGroup(groupId string, message []byte) {
	if p.mPisigContext.OnlineUserStore == nil {
		return
	}

	socketIdList := p.mPisigContext.OnlineUserStore.GetSocketIdListFromGroupId(groupId)

	for _, socketId := range socketIdList {
		conn, ok := p.mEPool.GetConnection(socketId)
		if conn == nil || !ok {
			continue
		}
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
	if p.mPisigContext.PisigSettings.EnableTLS {
		if glog.V(1) {
			glog.Infoln("Server is listening at: https://" + p.mPisigContext.PisigSettings.Host +
				":" + p.mPisigContext.PisigSettings.Port)
		}
		err := http.ListenAndServeTLS(p.mPisigContext.PisigSettings.Host+":"+p.mPisigContext.PisigSettings.Port,
			p.mPisigContext.PisigSettings.CertFile,
			p.mPisigContext.PisigSettings.KeyFile, p.mServerMux)
		if err != nil {
			glog.Errorln("Server error: ", err)
		}
	} else {
		if glog.V(1) {
			glog.Infoln("Server is listening at: http://" + p.mPisigContext.PisigSettings.Host +
				":" + p.mPisigContext.PisigSettings.Port)
		}
		err := http.ListenAndServe(p.mPisigContext.PisigSettings.Host+
			":"+p.mPisigContext.PisigSettings.Port, p.mServerMux)
		if err != nil {
			glog.Errorln("Server error: ", err)
		}
	}
}

func (p *Pisig) AddOnlineUser(conn net.Conn, uniqueId string, groupId string, data interface{}) bool {
	if p.mPisigContext.OnlineUserStore == nil {
		return false
	}
	socketId := WebsocketFileDescriptor(conn)
	return p.mPisigContext.OnlineUserStore.AddUser(uniqueId, groupId, socketId, data)
}

func (p *Pisig) GetTotalOnlineWebSocketConnection() int {
	return p.mEPool.TotalActiveConnections()
}

func (p *Pisig) GetTotalOnlineWSConnection() int {
	return p.mEPool.TotalActiveConnections()
}

func (p *Pisig) GetOnlineSocketIdList() []int {
	return p.mEPool.GetConnectionIdSlice()
}

func (p *Pisig) AddWSConnection(conn net.Conn) bool {
	return p.AddWebSocketConnection(conn)
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

func (p *Pisig) hookRemoveWebSocketConnection(conn net.Conn) error {
	if glog.V(3) {
		glog.Infof("Removing connection\n")
	}

	// CLEAN UP RESOURCES
	if p.mPisigContext.OnlineUserStore != nil {
		socketId := WebsocketFileDescriptor(conn)
		p.mPisigContext.OnlineUserStore.RemoveUser(socketId)
	}

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

	var pisigSettings *settings.PisigSettings
	var corsOptions *CORSOptions
	var pisigMessage message.PisigMessage
	var onlineUserStore storage.OnlineUserStore
	var pisigContext *PisigContext
	var serverMux *http.ServeMux

	for _, val := range args {
		if glog.V(3) {
			glog.Infof("Pisig arg type: %T\n", val)
		}
		switch val.(type) {

		case *PisigContext:
			pisigContext = val.(*PisigContext)
		case *settings.PisigSettings:
			pisigSettings = val.(*settings.PisigSettings)
		case *CORSOptions:
			corsOptions = val.(*CORSOptions)
		case message.PisigMessage:
			pisigMessage = val.(message.PisigMessage)
		case storage.OnlineUserStore:
			onlineUserStore = val.(storage.OnlineUserStore)
		case *http.ServeMux:
			serverMux = val.(*http.ServeMux)
		default:
			break
		}
	}

	if serverMux == nil {
		pisig.mServerMux = &http.ServeMux{}
	} else {
		pisig.mServerMux = serverMux
	}

	if pisigContext != nil {
		pisig.mPisigContext = pisigContext
		pisig.mPisigContext.TopicProducerQueue = make(event.TopicQueue,
			pisigContext.GetPisigSettings().TopicQueueSize)
		pisig.mTopicDispatcher = NewTopicDispatcher(pisig, pisig.mPisigContext.TopicProducerQueue)

		ePool, err := NewEPool(pisig.mPisigContext, pisig.hookRemoveWebSocketConnection)

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

		ePool, err := NewEPool(pisig.mPisigContext, pisig.hookRemoveWebSocketConnection)

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
