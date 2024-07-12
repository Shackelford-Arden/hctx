package cmd

func unsetTokens(shell string) string {

	output := ""

	switch shell {
	case "bash":
		fallthrough
	case "zsh":
		fallthrough
	default:
		output = `
unset NOMAD_TOKEN
unset CONSUL_TOKEN
unset VAULT_TOKEN
`
	}

	return output
}
