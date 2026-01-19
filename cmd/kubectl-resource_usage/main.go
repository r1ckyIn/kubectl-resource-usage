package main

import (
	"os"

	"github.com/r1ckyIn/kubectl-resource-usage/pkg/cmd"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-resource_usage", pflag.ExitOnError)
	pflag.CommandLine = flags

	streams := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	root := cmd.NewCmdResourceUsage(streams)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
