package cmd

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/mkawserm/pisig/pkg/core"
	"github.com/mkawserm/pisig/pkg/cors"
	"github.com/mkawserm/pisig/pkg/message"
	"github.com/mkawserm/pisig/pkg/settings"
	"github.com/mkawserm/pisig/pkg/view"
	"github.com/spf13/cobra"
)

type PisigCMDHookDefault struct {
}

func (dph *PisigCMDHookDefault) AppName() string {
	return core.ConstAppName
}

func (dph *PisigCMDHookDefault) AppVersion() string {
	return core.ConstAppVersion
}

func (dph *PisigCMDHookDefault) AppAuthors() string {
	return core.ConstAppAuthors
}

func (dph *PisigCMDHookDefault) AppDescription() string {
	return core.ConstAppDescription
}

func (dph *PisigCMDHookDefault) AppNameLong() string {
	return core.ConstAppDescription
}

func (dph *PisigCMDHookDefault) SetupCMD(pisigCMD *PisigCMD, pisigMessage message.PisigMessage) {
	serverCMD := &cobra.Command{
		Use:   "server",
		Short: "Run pisig server",
		Run: func(cmd *cobra.Command, args []string) {

			corsOptions := &cors.CORSOptions{
				AllowAllOrigins:  true,
				AllowCredentials: true,
				AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
				AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization"},
			}

			pisigSettings := settings.NewDefaultPisigSettings()

			pisig := core.NewPisigSimple(corsOptions, pisigSettings, pisigMessage)

			if glog.V(3) {
				glog.Infof("Registering all views")
			}
			pisig.AddView("/ws", &view.WebSocketView{})
			pisig.AddView("/", &view.ErrorView{})

			if glog.V(3) {
				glog.Infof("Running Pisig")
			}
			pisig.Run()
		},
	}

	pisigCMD.AddRunCommand(serverCMD)
}

func (dph *PisigCMDHookDefault) ProcessShellCMD(string) {

}

func (dph *PisigCMDHookDefault) ShellNewLinePrefix(inputCounter int) string {
	return fmt.Sprintf("%s%s%d%s%s ",
		dph.AppName(),
		"[",
		inputCounter,
		"]",
		"$",
	)
}
