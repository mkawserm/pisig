package cmd

import (
	"github.com/mkawserm/pisig/pkg/core"
	"github.com/spf13/cobra"
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
	serverCMD := &cobra.Command{
		Use:   "server",
		Short: "Run pisig server",
		Run: func(cmd *cobra.Command, args []string) {
			pisig := core.NewDefaultPisig()
			pisig.Run()
		},
	}
	pisigCMD.AddRunCommand(serverCMD)
}
