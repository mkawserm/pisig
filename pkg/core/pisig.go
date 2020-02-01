package core

import (
	"github.com/mkawserm/pisig/pkg/types"
	"net/http"
)

type Pisig struct {
	mServerMux *http.ServeMux

	mEventPool     *EventPool
	mCORSOptions   *types.CORSOptions
	mPisigContext  *types.PisigContext
	mPisigSettings *types.PisigSettings
}

func (p *Pisig) CORSOptions() *types.CORSOptions {
	return p.mCORSOptions
}

func (p *Pisig) PisigContext() *types.PisigContext {
	return p.mPisigContext
}

func (p *Pisig) PisigSettings() *types.PisigSettings {
	return p.mPisigSettings
}

func NewPisig(corsOptions *types.CORSOptions, pisigSettings *types.PisigSettings) *Pisig {

	serverMux := &http.ServeMux{}

	eventPool, err := NewEventPool(
		pisigSettings.EventPoolQueueSize,
		pisigSettings.EventPoolWaitingTime,
		nil,
		nil)

	if err != nil {
		panic(err)
		return nil
	}

	pisigContext := &types.PisigContext{}

	return &Pisig{
		mServerMux:     serverMux,
		mEventPool:     eventPool,
		mCORSOptions:   corsOptions,
		mPisigContext:  pisigContext,
		mPisigSettings: pisigSettings,
	}
}

func NewDefaultPisig() *Pisig {
	pisig := NewPisig(
		&types.CORSOptions{
			AllowAllOrigins:  true,
			AllowCredentials: true,
			AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
			AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization"},
		},
		types.NewDefaultPisigSettings(),
	)

	return pisig
}
