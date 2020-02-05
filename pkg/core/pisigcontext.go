package core

import (
	"github.com/mkawserm/pisig/pkg/cache"
	"github.com/mkawserm/pisig/pkg/event"
	"github.com/mkawserm/pisig/pkg/message"
	"github.com/mkawserm/pisig/pkg/settings"
)

type PisigContext struct {
	PisigStore           *cache.PisigStore
	PisigServiceRegistry *PisigServiceRegistry

	CORSOptions   *CORSOptions
	PisigMessage  message.PisigMessage
	PisigSettings *settings.PisigSettings

	TopicProducerQueue event.TopicQueue
}

func (pc *PisigContext) GetCORSOptions() *CORSOptions {
	return pc.CORSOptions
}

func (pc *PisigContext) GetPisigSettings() *settings.PisigSettings {
	return pc.PisigSettings
}

func (pc *PisigContext) GetPisigStore() *cache.PisigStore {
	return pc.PisigStore
}

func (pc *PisigContext) GetPisigServiceRegistry() *PisigServiceRegistry {
	return pc.PisigServiceRegistry
}

func (pc *PisigContext) GetPisigMessage() message.PisigMessage {
	return pc.PisigMessage
}

func (pc *PisigContext) ProduceTopic(topic event.Topic) {
	pc.TopicProducerQueue <- topic
}

// Create new PisigContext
func NewPisigContext() *PisigContext {
	return &PisigContext{
		PisigStore:           cache.NewPisigStore(),
		PisigServiceRegistry: NewPisigServiceRegistry(),
	}
}
