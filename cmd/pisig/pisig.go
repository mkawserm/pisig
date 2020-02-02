package main

import (
	"github.com/mkawserm/pisig/pkg/cmd"
	"github.com/mkawserm/pisig/pkg/conf"
)

func main() {
	pisigCMD := cmd.PisigCMD{
		PisigCMDHook:  &cmd.DefaultPisigCMDHook{},
		PisigResponse: &conf.DefaultPisigResponse{},
	}

	pisigCMD.Setup()
	pisigCMD.Execute()
}
