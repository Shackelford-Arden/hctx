package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

// Activate provides the given shell with the commands to set environment variables
// The intent is that this is sourced/imported in your shell, and not called ad-hoc.
func Activate(_ context.Context, cmd *cli.Command) error {

	execPath, _ := os.Executable()

	sh := resolveShell(cmd)

	fmt.Print(sh.Activate(execPath))

	return nil
}
