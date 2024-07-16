package models

type ConsulConfig struct {
	Address   string `hcl:"address,optional"`
	Namespace string `hcl:"namespace,optional"`
}

var ConsulAddr = "CONSUL_HTTP_ADDR"
var ConsulNamespace = "CONSUL_NAMESPACE"
var ConsulToken = "CONSUL_HTTP_TOKEN"

// Use provides commands to set appropriate Consul environment variables.
func (n *ConsulConfig) Use(shell string, token string) []string {
	var envCommands []string

	if n.Address != "" {
		envCommands = append(envCommands, genUseCommands(shell, ConsulAddr, n.Address))
	}

	if n.Namespace != "" {
		envCommands = append(envCommands, genUseCommands(shell, ConsulNamespace, n.Namespace))
	}

	if token != "" {
		envCommands = append(envCommands, genUseCommands(shell, ConsulToken, token))
	}

	return envCommands
}

// Unset Provides commands to unset the Consul environment variables for the given stack
func (n *ConsulConfig) Unset(shell string) []string {

	var unsetCommands []string

	if n.Address != "" {
		unsetCommands = append(unsetCommands, genUnsetCommands(shell, ConsulAddr))
	}

	if n.Namespace != "" {
		unsetCommands = append(unsetCommands, genUnsetCommands(shell, ConsulNamespace))
	}

	unsetCommands = append(unsetCommands, genUnsetCommands(shell, ConsulToken))

	return unsetCommands
}
