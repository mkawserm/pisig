package cmd

type PisigHook interface {
	AppName() string
	AppNameLong() string
	AppVersion() string
	AppAuthors() string
	AppDescription() string

	SetupCMD(pisigCMD *PisigCMD)
}
