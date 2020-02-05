package core

import (
	"github.com/golang/glog"
	"github.com/mkawserm/pisig/pkg/event"
)

type TopicProcessor struct {
	Pisig          *Pisig
	quit           chan bool
	TopicQueue     event.TopicQueue
	TopicQueuePool event.TopicQueuePool
}

func NewTopicProcessor(pisig *Pisig, topicQueuePool event.TopicQueuePool) TopicProcessor {
	return TopicProcessor{
		Pisig:          pisig,
		TopicQueuePool: topicQueuePool,
		TopicQueue:     make(event.TopicQueue),
		quit:           make(chan bool),
	}
}

func (tp TopicProcessor) Start() {
	go func() {
		//if glog.V(3) {
		//	glog.Infof("Topic processor start goroutine started")
		//}

		for {
			tp.TopicQueuePool <- tp.TopicQueue
			select {
			case topic := <-tp.TopicQueue:
				// we have received a topic do something with it.
				if glog.V(3) {
					glog.Infof("Distributing topic to the topic listener\n")
				}

				topicListenerList := tp.Pisig.GetTopicListenerList(topic.Name)
				for i := range topicListenerList {
					topicListener := topicListenerList[i]

					pisigService, isPisigService := topicListener.(PisigService)
					if isPisigService {
						err, _ := pisigService.Process(topic, false)
						if err != nil {
							glog.Errorf("Error: %v\n", err)
						}
					}
				}

			case <-tp.quit:
				return
			}
		}

		//if glog.V(3) {
		//	glog.Infof("Topic processor start goroutine finished")
		//}
	}()
}

// Stop topic processor
func (tp TopicProcessor) Stop() {
	go func() {
		//if glog.V(3) {
		//	glog.Infof("Topic processor stop goroutine started")
		//}

		tp.quit <- true

		//if glog.V(3) {
		//	glog.Infof("Topic processor stop goroutine finished")
		//}
	}()
}
