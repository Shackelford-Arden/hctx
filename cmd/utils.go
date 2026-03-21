package cmd

import (
	consul "github.com/hashicorp/consul/api"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/urfave/cli/v3"

	"github.com/Shackelford-Arden/hctx/internal/shells"
)

// validNomadToken checks if the given token is still valid.
// Docs: https://developer.hashicorp.com/nomad/api-docs/acl/tokens#read-self-token
func validNomadToken(addr string, token string) bool {

	nom, _ := nomad.NewClient(&nomad.Config{
		Address: addr,
	})

	nomToken, _, err := nom.ACLTokens().Self(&nomad.QueryOptions{AuthToken: token})
	if err != nil {
		// TODO Consider logging some things to a log file somewhere?
		// Preferring to avoid logging text
		return false
	}

	// Likely means the API return a 404 on the token,
	// indicating it has expired.
	if nomToken == nil {
		return false
	}

	return true
}

// validConsulToken checks to see if the given token is valid anymore.
// Docs: https://developer.hashicorp.com/consul/api-docs/acl/tokens#read-self-token
func validConsulToken(addr string, token string) bool {

	con, _ := consul.NewClient(&consul.Config{
		Address: addr,
	})

	consulToken, _, err := con.ACL().TokenReadSelf(&consul.QueryOptions{Token: token})
	if err != nil {
		// TODO Consider logging some things to a log file somewhere?
		// Preferring to avoid logging text
		return false
	}

	// Likely means the API return a 404 on the token,
	// indicating it has expired.
	if consulToken == nil {
		return false
	}

	return true
}

// resolveShell determines the shell to use based on the --shell CLI flag,
// falling back to the config's auto-discovered shell value.
func resolveShell(cmd *cli.Command) shells.Shell {
	name := cmd.String("shell")
	if name == "" {
		name = AppConfig.Shell
	}
	switch name {
	case "nushell":
		return &shells.Nushell{}
	default:
		return &shells.Bash{}
	}
}
