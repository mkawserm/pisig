package cmd

import "github.com/mkawserm/pisig/pkg/conf"

type PisigHook interface {
	AppName() string
	AppNameLong() string
	AppVersion() string
	AppAuthors() string
	AppDescription() string

	SetupCMD(pisigCMD *PisigCMD, pisigResponse conf.PisigResponse)
}
