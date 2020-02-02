package cmd

import "github.com/mkawserm/pisig/pkg/conf"

type PisigCMDHook interface {
	AppName() string
	AppNameLong() string
	AppVersion() string
	AppAuthors() string
	AppDescription() string

	SetupCMD(pisigCMD *PisigCMD, pisigResponse conf.PisigResponse)
}
