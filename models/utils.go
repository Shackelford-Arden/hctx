package models

import "fmt"

func genUseCommands(shell string, varName string, varValue string) string {

	var envCmd string

	switch shell {
	case "bash", "zsh":

		envCmd = fmt.Sprintf("export %s=%s", varName, varValue)

	default:
		envCmd = fmt.Sprintf("export %s=%s", varName, varValue)
	}

	return fmt.Sprintf("%s\n", envCmd)

}

func genUnsetCommands(shell string, varName string) string {
	var unuseCmd string

	switch shell {
	case "base", "zsh":
		unuseCmd = fmt.Sprintf("unset %s", varName)

	default:
		unuseCmd = fmt.Sprintf("unset %s", varName)
	}

	return unuseCmd

}
