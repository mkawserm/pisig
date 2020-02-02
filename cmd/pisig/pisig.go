package main

import "github.com/mkawserm/pisig/pkg/cmd"

func main() {
	pisigCMD := cmd.PisigCMD{
		PisigHook: &cmd.DefaultPisigHook{},
	}

	pisigCMD.Setup()
	pisigCMD.Execute()
}
