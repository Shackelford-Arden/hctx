package shells

import (
	"encoding/json"
	"fmt"
)

type Nushell struct {
}

func (s *Nushell) Activate(binPath string) string {
	return fmt.Sprintf(`
# --env is needed because we're modifying environment variables that we want persistent
# --wrapped is needed to make sure thing flags like '--help' are passed to hctx, since
#   Nu tries to be helpful with the custom command and provide auto-generated help docs.
#   Which, if for whatever reason I wanted to down the road, 'help hctx' in Nu still outputs
#   those generated docs.
def --env --wrapped hctx [
    ...rest
    ] {

  # Need to know the full path to the hctx binary first
  let hctx_path = "%s"

  # If no command, get hctx to run so the default --help text
  # is provided to the user.
  if ($rest | is-empty) {
	^$hctx_path
	return
  }

  let command = $rest.0
  # For later: I'm not sure why I have to use 'slice' here instead of 'range'
  let rem_args = if ($rest | length) > 1 { $rest | slice 1.. } else { [] }

  if ($command in ["use", "u", "unset", "un"]) and (not ($rem_args | any { |arg| $arg in ["--help", "-h"] })) {
    let env_out = (do {^$hctx_path $command ...$rem_args } | complete)

    if $env_out.exit_code == 0 {

         if ($command in ["unset", "un"]) and ($env_out.stdout != "") {
            # Iterate over the array given from JSON
            let unsetVars = ($env_out.stdout | from json)
            hide-env --ignore-errors ...$unsetVars
           return
         }

        # For some reason, I need to have this come _after_ the "unset" logic block
        # as the " | from json" seems to fail... even though we shouldn't hit this
         if ($command in ["use", "u"]) {
            $env_out.stdout | from json | load-env
            return
          }
    } else {
		echo $env_out.stdout
    }
  } else {
    ^$hctx_path $command ...$rem_args
  }
}
`, binPath,
	)
}

func (s *Nushell) UseOutput(envVars map[string]string) string {
	// For Nu, we end up having to deal with things a little differently,
	// so for `use` we need to output to JSON. The script from Activate
	// will parse the JSON and save the values as environment variables.

	unsetOut, _ := json.MarshalIndent(envVars, "", "  ")
	return string(unsetOut)
}

func (s *Nushell) UnsetOutput(envVars []string) string {

	// For Nu, we're going to again output as JSON so that the script
	// from Activate can easily iterate on the values and run the correct
	// commands.

	unsetOut, _ := json.MarshalIndent(envVars, "", "  ")
	return string(unsetOut)
}
