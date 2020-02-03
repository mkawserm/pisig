package cmd

import (
	"github.com/mkawserm/pisig/pkg/message"
)

type PisigCMDHook interface {
	AppName() string
	AppNameLong() string
	AppVersion() string
	AppAuthors() string
	AppDescription() string

	SetupCMD(pisigCMD *PisigCMD, pisigMessage message.PisigMessage)
}
