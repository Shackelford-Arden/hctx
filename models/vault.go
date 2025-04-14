package models

type VaultConfig struct {
	Address   string `hcl:"address,optional"`
	Namespace string `hcl:"namespace,optional"`
}

var VaultAddr = "VAULT_ADDR"
var VaultNamespace = "VAULT_NAMESPACE"

// Use provides commands to set appropriate Vault environment variables.
func (n *VaultConfig) Use() map[string]string {
	envVars := map[string]string{}

	if n.Address != "" {
		envVars[VaultAddr] = n.Address
	}

	if n.Namespace != "" {
		envVars[VaultNamespace] = n.Namespace
	}

	return envVars
}

// Unset Provides commands to unset the Vault environment variables for the given stack
func (n *VaultConfig) Unset() []string {

	envVarNames := []string{}

	if n.Address != "" {
		envVarNames = append(envVarNames, VaultAddr)
	}

	if n.Namespace != "" {
		envVarNames = append(envVarNames, VaultNamespace)
	}

	return envVarNames
}
