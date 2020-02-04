package core

import (
	"github.com/golang/glog"
	"github.com/mkawserm/pisig/pkg/event"
)

type TopicDispatcher struct {
	Pisig              *Pisig
	TopicQueuePoolSize int
	TopicProducerQueue event.TopicQueue
	TopicQueuePool     event.TopicQueuePool
}

func NewTopicDispatcher(pisig *Pisig, topicQueue event.TopicQueue) *TopicDispatcher {
	topicQueuePool := make(event.TopicQueuePool, pisig.PisigSettings().TopicQueuePoolSize)

	return &TopicDispatcher{
		Pisig:              pisig,
		TopicQueuePoolSize: pisig.PisigSettings().TopicQueuePoolSize,
		TopicProducerQueue: topicQueue,
		TopicQueuePool:     topicQueuePool,
	}
}

func (td *TopicDispatcher) Run() {
	if glog.V(3) {
		glog.Infof("Running topic dispatcher")
	}

	for i := 0; i < td.TopicQueuePoolSize; i++ {
		worker := NewTopicProcessor(td.Pisig, td.TopicQueuePool)
		worker.Start()
	}

	go td.dispatch()

	if glog.V(3) {
		glog.Infof("Topic dispatcher started")
	}
}

func (td *TopicDispatcher) dispatch() {
	for {
		select {
		case topic := <-td.TopicProducerQueue:
			go func(topic event.Topic) {
				topicQueue := <-td.TopicQueuePool
				topicQueue <- topic
			}(topic)
		}
	}
}
