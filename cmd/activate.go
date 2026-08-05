package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// Activate provides the given shell with the commands to set environment variables
// The intent is that this is sourced/imported in your shell, and not called ad-hoc.
func Activate(ctx *cli.Context) error {

	execPath, _ := os.Executable()

	fmt.Print(ActiveShell.Activate(execPath))

	return nil
}
