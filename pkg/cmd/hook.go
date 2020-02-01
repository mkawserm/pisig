package cmd

import (
	"github.com/mkawserm/pisig/pkg/core"
)

type DefaultPisigHook struct {
}

func (dph *DefaultPisigHook) AppName() string {
	return core.ConstAppName
}

func (dph *DefaultPisigHook) AppVersion() string {
	return core.ConstAppVersion
}

func (dph *DefaultPisigHook) AppAuthors() string {
	return core.ConstAppAuthors
}

func (dph *DefaultPisigHook) AppDescription() string {
	return core.ConstAppDescription
}

func (dph *DefaultPisigHook) AppNameLong() string {
	return core.ConstAppDescription
}

func (dph *DefaultPisigHook) SetupCMD(pisigCMD *PisigCMD) {

}

var Hook = &DefaultPisigHook{}
