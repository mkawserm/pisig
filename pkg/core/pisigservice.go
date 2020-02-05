package core

import (
	"github.com/mkawserm/pisig/pkg/event"
)

type PisigService interface {
	SetPisig(pisig *Pisig)

	GroupName() string
	ServiceName() string
	ServiceVersion() string
	ServiceAuthors() string

	Process(topic event.Topic, synchronous bool) (error, *event.Topic)
}
