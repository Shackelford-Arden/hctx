package models

type ConsulConfig struct {
	Address   string `hcl:"address,optional"`
	Namespace string `hcl:"namespace,optional"`
}

var ConsulAddr = "CONSUL_HTTP_ADDR"
var ConsulNamespace = "CONSUL_NAMESPACE"
var ConsulToken = "CONSUL_HTTP_TOKEN"

// Use provides commands to set appropriate Consul environment variables.
func (n *ConsulConfig) Use(token string) map[string]string {
	envVars := map[string]string{}

	if n.Address != "" {
		envVars[ConsulAddr] = n.Address
	}

	if n.Namespace != "" {
		envVars[ConsulNamespace] = n.Namespace
	}

	if token != "" {
		envVars[ConsulToken] = token
	}

	return envVars
}

// Unset Provides commands to unset the Consul environment variables for the given stack
func (n *ConsulConfig) Unset() []string {

	envVarNames := []string{ConsulToken}

	if n.Address != "" {
		envVarNames = append(envVarNames, ConsulAddr)
	}

	if n.Namespace != "" {
		envVarNames = append(envVarNames, ConsulNamespace)
	}

	return envVarNames
}
