package cmd

import (
	goFlag "flag"
	"fmt"
	"github.com/mkawserm/pisig/pkg/core"
	"github.com/mkawserm/pisig/pkg/message"
	"github.com/spf13/cobra"
	pFlag "github.com/spf13/pflag"
	"os"
)

type PisigCMD struct {
	PisigCMDHook PisigCMDHook
	PisigMessage message.PisigMessage

	mRootCMD      *cobra.Command
	mRunSubCMD    *cobra.Command
	mCreateSubCMD *cobra.Command
}

func (pc *PisigCMD) AddCommand(cmds ...*cobra.Command) {
	pc.mRootCMD.AddCommand(cmds...)
}

func (pc *PisigCMD) AddRunCommand(cmds ...*cobra.Command) {
	pc.mRunSubCMD.AddCommand(cmds...)
}

func (pc *PisigCMD) AddCreateCommand(cmds ...*cobra.Command) {
	pc.mCreateSubCMD.AddCommand(cmds...)
}

func (pc *PisigCMD) Setup() {
	pc.mRootCMD = &cobra.Command{
		Use:   pc.PisigCMDHook.AppName(),
		Short: pc.PisigCMDHook.AppNameLong(),
		Long:  pc.PisigCMDHook.AppDescription(),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// For cobra + glog flags. Available to all sub commands.
			goFlag.Parse()
		},
	}

	pc.mRunSubCMD = &cobra.Command{
		Use:   "run",
		Short: "run any run command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Usage())
		},
	}

	pc.mCreateSubCMD = &cobra.Command{
		Use:   "create",
		Short: "run any create command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Usage())
		},
	}

	pc.mRootCMD.AddCommand(pc.mRunSubCMD)
	pc.mRootCMD.AddCommand(pc.mCreateSubCMD)
	pc.mRootCMD.AddCommand(getPisigSubCommand())

	// SETUP CUSTOM CMDS FROM HOOK
	pc.PisigCMDHook.SetupCMD(pc, pc.PisigMessage)

	pFlag.CommandLine.AddGoFlagSet(goFlag.CommandLine)
}

func (pc *PisigCMD) Execute() {
	if err := pc.mRootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getPisigSubCommand() *cobra.Command {
	pisigSubCommand := &cobra.Command{
		Use:   core.ConstAppName,
		Short: "pisig core",
		Long:  core.ConstAppDescription,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
		},
	}

	pisigVersionCommand := &cobra.Command{
		Use:   "version",
		Short: "print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(core.ConstAppVersion)
		},
	}

	pisigAuthorsCommand := &cobra.Command{
		Use:   "authors",
		Short: "print the authors",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(core.ConstAppAuthors)
		},
	}

	pisigSubCommand.AddCommand(pisigVersionCommand)
	pisigSubCommand.AddCommand(pisigAuthorsCommand)

	return pisigSubCommand
}
