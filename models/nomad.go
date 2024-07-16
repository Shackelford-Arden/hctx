package models

type NomadConfig struct {
	Address   string `hcl:"address,optional"`
	Namespace string `hcl:"namespace,optional"`
}

var NomadAddr = "NOMAD_ADDR"
var NomadNamespace = "NOMAD_NAMESPACE"
var NomadToken = "NOMAD_TOKEN"

// Use provides commands to set appropriate Nomad environment variables.
func (n *NomadConfig) Use(shell string, token string) []string {
	var envCommands []string

	if n.Address != "" {
		envCommands = append(envCommands, genUseCommands(shell, NomadAddr, n.Address))
	}

	if n.Namespace != "" {
		envCommands = append(envCommands, genUseCommands(shell, NomadNamespace, n.Namespace))
	}

	if token != "" {
		envCommands = append(envCommands, genUseCommands(shell, NomadToken, token))
	}

	return envCommands
}

// Unset Provides commands to unset the Nomad environment variables for the given stack
func (n *NomadConfig) Unset(shell string) []string {

	var unsetCommands []string

	if n.Address != "" {
		unsetCommands = append(unsetCommands, genUnsetCommands(shell, NomadAddr))
	}

	if n.Namespace != "" {
		unsetCommands = append(unsetCommands, genUnsetCommands(shell, NomadNamespace))
	}

	unsetCommands = append(unsetCommands, genUnsetCommands(shell, NomadToken))

	return unsetCommands
}
