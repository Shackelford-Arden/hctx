package models

type ConsulConfig struct {
	Address   string `hcl:"address,optional"`
	Namespace string `hcl:"namespace,optional"`
}

var ConsulAddr = "CONSUL_HTTP_ADDR"
var ConsulNamespace = "CONSUL_NAMESPACE"

// Use provides commands to set appropriate Consul environment variables.
func (n *ConsulConfig) Use(shell string) []string {
	var envCommands []string

	if n.Address != "" {
		envCommands = append(envCommands, genUseCommands(shell, ConsulAddr, n.Address))
	}

	if n.Namespace != "" {
		envCommands = append(envCommands, genUseCommands(shell, ConsulNamespace, n.Namespace))
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

	return unsetCommands
}
