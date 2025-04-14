package models

type NomadConfig struct {
	Address   string `hcl:"address,optional"`
	Namespace string `hcl:"namespace,optional"`
}

var NomadAddr = "NOMAD_ADDR"
var NomadNamespace = "NOMAD_NAMESPACE"
var NomadToken = "NOMAD_TOKEN"

// Use provides a map of environment variables and their expected values that will
// be used by a shells.Shell to generate appropriate commands to set the environment
// variables.
func (n *NomadConfig) Use(token string) map[string]string {
	envVars := map[string]string{}

	if n.Address != "" {
		envVars[NomadAddr] = n.Address
	}

	if n.Namespace != "" {
		envVars[NomadNamespace] = n.Namespace
	}

	if token != "" {
		envVars[NomadToken] = token
	}

	return envVars
}

// Unset provides a list of environment variables that should be unset by the shells.Shell.
func (n *NomadConfig) Unset() []string {

	envVarNames := []string{NomadToken}

	if n.Address != "" {
		envVarNames = append(envVarNames, NomadAddr)
	}

	if n.Namespace != "" {
		envVarNames = append(envVarNames, NomadNamespace)
	}

	return envVarNames
}
