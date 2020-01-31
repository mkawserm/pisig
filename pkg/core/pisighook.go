package core

type PisigHook interface {
	AppName() string
	AppVersion() string
	AppAuthors() string
}

type DefaultPisigHook struct {
}

func (dph *DefaultPisigHook) AppName() string {
	return ConstAppName
}

func (dph *DefaultPisigHook) AppVersion() string {
	return ConstAppVersion
}

func (dph *DefaultPisigHook) AppAuthors() string {
	return ConstAppAuthors
}
