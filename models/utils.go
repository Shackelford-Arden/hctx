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

	return envCmd

}

func genUnsetCommands(shell string, varName string) string {
	var unsetCmd string

	switch shell {
	case "bash", "zsh":
		unsetCmd = fmt.Sprintf("unset %s", varName)

	default:
		unsetCmd = fmt.Sprintf("unset %s", varName)
	}

	return unsetCmd

}
