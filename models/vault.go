package models

type VaultConfig struct {
	Address   string `hcl:"address,optional"`
	Namespace string `hcl:"namespace,optional"`
}

var VaultAddr = "VAULT_ADDR"
var VaultNamespace = "VAULT_NAMESPACE"

// Use provides commands to set appropriate Vault environment variables.
func (n *VaultConfig) Use(shell string) []string {
	var envCommands []string

	if n.Address != "" {
		envCommands = append(envCommands, genUseCommands(shell, VaultAddr, n.Address))
	}

	if n.Namespace != "" {
		envCommands = append(envCommands, genUseCommands(shell, VaultNamespace, n.Namespace))
	}

	return envCommands
}

// Unset Provides commands to unset the Vault environment variables for the given stack
func (n *VaultConfig) Unset(shell string) []string {

	var unsetCommands []string

	if n.Address != "" {
		unsetCommands = append(unsetCommands, genUnsetCommands(shell, VaultAddr))
	}

	if n.Namespace != "" {
		unsetCommands = append(unsetCommands, genUnsetCommands(shell, VaultNamespace))
	}

	return unsetCommands
}
