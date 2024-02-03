package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

// Activate provides the given shell with the commands to set environment variables
// The intent is that this is sourced/imported in your shell, and not called ad-hoc.
func Activate(ctx *cli.Context) error {

	execPath, _ := os.Executable()

	// This bash/zsh script is pretty much the same as what mise has
	// Reference: https://github.com/jdx/mise/blob/be34b768d9c09feda3c59d9a949a40609c294dcf/src/shell/zsh.rs#L17
	activateScript := fmt.Sprintf(
		`
hctx () {
  local command
  HCTX_PATH='%s'
  command="${1:-}"
  if [ "$#" = 0 ]
  then
    command $HCTX_PATH
    return
  fi
  shift
  case "$command" in
  (use|unset) if [[ ! " $@ " =~ " --help " ]] && [[ ! " $@ " =~ " -h " ]]
    then
      eval "$(command $HCTX_PATH "$command" "$@")"
      return $?
    fi ;;
  esac
  command $HCTX_PATH "$command" "$@"
}
`, execPath,
	)

	fmt.Print(activateScript)

	return nil
}