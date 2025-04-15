package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Shackelford-Arden/hctx/internal/shells"
)

// Activate provides the given shell with the commands to set environment variables
// The intent is that this is sourced/imported in your shell, and not called ad-hoc.
func Activate(ctx *cli.Context) error {

	execPath, _ := os.Executable()

	shell := ctx.String("shell")
	if shell == "" {
		shell = AppConfig.Shell
	}

	var sh shells.Shell

	switch shell {
	case "nushell":
		sh = &shells.Nushell{}

	default:
		// This bash/zsh script is pretty much the same as what mise has
		// Reference: https://github.com/jdx/mise/blob/be34b768d9c09feda3c59d9a949a40609c294dcf/src/shell/zsh.rs#L17
		sh = &shells.Bash{}
	}

	fmt.Print(sh.Activate(execPath))

	return nil
}
