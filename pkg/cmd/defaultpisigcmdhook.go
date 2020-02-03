package cmd

import (
	"github.com/golang/glog"
	"github.com/mkawserm/pisig/pkg/core"
	"github.com/mkawserm/pisig/pkg/cors"
	"github.com/mkawserm/pisig/pkg/message"
	"github.com/mkawserm/pisig/pkg/settings"
	"github.com/mkawserm/pisig/pkg/view"
	"github.com/spf13/cobra"
)

type DefaultPisigCMDHook struct {
}

func (dph *DefaultPisigCMDHook) AppName() string {
	return core.ConstAppName
}

func (dph *DefaultPisigCMDHook) AppVersion() string {
	return core.ConstAppVersion
}

func (dph *DefaultPisigCMDHook) AppAuthors() string {
	return core.ConstAppAuthors
}

func (dph *DefaultPisigCMDHook) AppDescription() string {
	return core.ConstAppDescription
}

func (dph *DefaultPisigCMDHook) AppNameLong() string {
	return core.ConstAppDescription
}

func (dph *DefaultPisigCMDHook) SetupCMD(pisigCMD *PisigCMD, pisigMessage message.PisigMessage) {
	serverCMD := &cobra.Command{
		Use:   "server",
		Short: "Run pisig server",
		Run: func(cmd *cobra.Command, args []string) {
			pisig := core.NewPisigSimple(
				&cors.CORSOptions{
					AllowAllOrigins:  true,
					AllowCredentials: true,
					AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
					AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization"},
				},
				settings.NewDefaultPisigSettings(),
				pisigMessage,
			)
			if glog.V(3) {
				glog.Infof("Registering all views")
			}
			pisig.AddView("/", &view.ErrorView{})

			if glog.V(3) {
				glog.Infof("Running Pisig")
			}
			pisig.Run()
		},
	}

	pisigCMD.AddRunCommand(serverCMD)
}
