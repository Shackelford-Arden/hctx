package shells

import (
	"fmt"
)

type Bash struct {
}

func (s *Bash) Activate(binPath string) string {
	return fmt.Sprintf(`
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
  (use|u|unset|un) if [[ ! " $@ " =~ " --help " ]] && [[ ! " $@ " =~ " -h " ]]
    then

      USE_OUTPUT="$(command $HCTX_PATH "$command" "$@")"

      # Doing this if to avoid attempting to run the eval command
      # on somethinig that shouldn't be run through eval (like an error message!)
      if [[ $? -eq 0 ]]; then
        eval "${USE_OUTPUT}"
      else
        echo "${USE_OUTPUT}"
      fi

      return $?
    fi ;;
  esac
  command $HCTX_PATH "$command" "$@"
}
`, binPath,
	)
}

func (s *Bash) UseOutput(envVars map[string]string) string {

	var useOut string

	for key, value := range envVars {
		useOut += fmt.Sprintf(`export %s="%s"`, key, value) + "\n"
	}

	return useOut
}

func (s *Bash) UnsetOutput(envVars []string) string {

	var unsetOut string

	for _, value := range envVars {
		unsetOut += fmt.Sprintf(`unset %s`, value) + "\n"
	}

	return unsetOut
}
