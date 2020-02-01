package main

import "github.com/mkawserm/pisig/pkg/cmd"

func main() {
	pisigCMD := cmd.PisigCMD{}
	pisigCMD.Setup()
	pisigCMD.Execute()
}
