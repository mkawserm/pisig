package main

import (
	"github.com/mkawserm/pisig/pkg/cmd"
	"github.com/mkawserm/pisig/pkg/message"
)

func main() {
	pisigCMD := cmd.PisigCMD{
		PisigCMDHook: &cmd.PisigCMDHookDefault{},
		PisigMessage: &message.PisigMessageDefault{},
	}

	pisigCMD.Setup()
	pisigCMD.Execute()
}
