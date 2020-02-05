package context

import (
	"github.com/mkawserm/pisig/pkg/cache"
	"github.com/mkawserm/pisig/pkg/cors"
	"github.com/mkawserm/pisig/pkg/event"
	"github.com/mkawserm/pisig/pkg/message"
	"github.com/mkawserm/pisig/pkg/registry"
	"github.com/mkawserm/pisig/pkg/settings"
)

type PisigContext struct {
	PisigStore           *cache.PisigStore
	PisigServiceRegistry *registry.PisigServiceRegistry

	CORSOptions   *cors.CORSOptions
	PisigMessage  message.PisigMessage
	PisigSettings *settings.PisigSettings

	TopicProducerQueue event.TopicQueue
}

func (pc *PisigContext) GetCORSOptions() *cors.CORSOptions {
	return pc.CORSOptions
}

func (pc *PisigContext) GetPisigSettings() *settings.PisigSettings {
	return pc.PisigSettings
}

func (pc *PisigContext) GetPisigStore() *cache.PisigStore {
	return pc.PisigStore
}

func (pc *PisigContext) GetPisigServiceRegistry() *registry.PisigServiceRegistry {
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
		PisigServiceRegistry: registry.NewPisigServiceRegistry(),
	}
}
