package service

import (
	"github.com/mkawserm/pisig/pkg/event"
)

type PisigService interface {
	SetSettings(settings map[string]interface{}) (error, bool)
	UpdateSettings(settings map[string]interface{}) (error, bool)
	SetTopicProducerHandler(func(topic event.Topic))

	GroupName() string
	ServiceName() string
	ServiceVersion() string
	ServiceAuthors() string

	Process(topic event.Topic, synchronous bool) (error, *event.Topic)
}
