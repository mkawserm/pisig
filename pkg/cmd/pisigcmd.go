package cmd

import (
	"bufio"
	goFlag "flag"
	"fmt"
	"github.com/mkawserm/pisig/pkg/core"
	"github.com/mkawserm/pisig/pkg/message"
	"github.com/spf13/cobra"
	pFlag "github.com/spf13/pflag"
	"os"
	"strings"
)

type PisigCMD struct {
	CMDHook PisigCMDHook
	Message message.PisigMessage

	AllowPisigCMD bool

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
		Use:   pc.CMDHook.AppName(),
		Short: pc.CMDHook.AppNameLong(),
		Long:  pc.CMDHook.AppDescription(),
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
	pc.mRootCMD.AddCommand(pc.getShellSubCommand())

	if pc.AllowPisigCMD {
		pc.mRootCMD.AddCommand(pc.getPisigSubCommand())
	}

	// INIT ALL STATUS CODE
	pc.Message.InitAllStatusCode()

	// SETUP CUSTOM CMDS FROM HOOK
	pc.CMDHook.SetupCMD(pc, pc.Message)

	pFlag.CommandLine.AddGoFlagSet(goFlag.CommandLine)
}

func (pc *PisigCMD) Execute() {
	if err := pc.mRootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (pc *PisigCMD) getPisigSubCommand() *cobra.Command {
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

func (pc *PisigCMD) getShellSubCommand() *cobra.Command {
	shellSubCommand := &cobra.Command{
		Use:   "shell",
		Short: "command interpreter",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)
			inputCounter := 0
			for {
				inputCounter++
				fmt.Printf(pc.CMDHook.ShellNewLinePrefix(inputCounter))

				cmdString, err := reader.ReadString('\n')
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, err)
				}
				cmdString = strings.TrimSuffix(cmdString, "\n")

				switch cmdString {
				case "clear":
					fmt.Print("\x1b[H\x1b[2J")
				case "reset":
					fmt.Print("\x1b[H\x1b[2J")
					inputCounter = 0
				case "version":
					fmt.Println(pc.CMDHook.AppVersion())
				case "authors":
					fmt.Println(pc.CMDHook.AppAuthors())
				case "exit":
					os.Exit(0)
				default:
					pc.CMDHook.ProcessShellCMD(cmdString)
				}
			}

		},
	}
	return shellSubCommand
}
