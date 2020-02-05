package service

import (
	"github.com/mkawserm/pisig/pkg/context"
	"github.com/mkawserm/pisig/pkg/event"
)

type PisigService interface {
	Setup(pisigContext *context.PisigContext)

	// SetSettings(settings map[string]interface{}) (error, bool)
	// UpdateSettings(settings map[string]interface{}) (error, bool)

	GroupName() string
	ServiceName() string
	ServiceVersion() string
	ServiceAuthors() string

	Process(topic event.Topic, synchronous bool) (error, *event.Topic)
}
