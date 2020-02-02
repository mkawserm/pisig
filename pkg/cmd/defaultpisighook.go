package cmd

import (
	"github.com/mkawserm/pisig/pkg/conf"
	"github.com/mkawserm/pisig/pkg/core"
	"github.com/mkawserm/pisig/pkg/variant"
	"github.com/mkawserm/pisig/pkg/view"
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

func (dph *DefaultPisigHook) SetupCMD(pisigCMD *PisigCMD, pisigResponse conf.PisigResponse) {
	serverCMD := &cobra.Command{
		Use:   "server",
		Short: "Run pisig server",
		Run: func(cmd *cobra.Command, args []string) {
			pisig := core.NewPisigSimple(
				&variant.CORSOptions{
					AllowAllOrigins:  true,
					AllowCredentials: true,
					AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
					AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization"},
				},
				variant.NewDefaultPisigSettings(),
				pisigResponse,
			)

			pisig.AddView("/", &view.ErrorView{})

			pisig.Run()
		},
	}

	pisigCMD.AddRunCommand(serverCMD)
}
