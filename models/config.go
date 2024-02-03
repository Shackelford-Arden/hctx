package models

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"os"
	"strings"
)

type Config struct {
	Stacks []Stack `hcl:"stack,block"`
	Shell  string  `hcl:"shell,optional"`
}

type Stack struct {
	Name   string        `hcl:",label"`
	Nomad  *NomadConfig  `hcl:"nomad,block"`
	Consul *ConsulConfig `hcl:"consul,block"`
	Vault  *VaultConfig  `hcl:"vault,block"`
}

func NewConfig() (*Config, error) {
	// Get user homedir
	userHome, homeErr := os.UserHomeDir()
	if homeErr != nil {
		fmt.Printf("failed to get user homedir: %s", homeErr)
		os.Exit(10)
	}

	configPath := fmt.Sprintf("%s/.config/%s", userHome, ".hctx.hcl")

	// Check if there is a config file
	_, statErr := os.Stat(configPath)
	if statErr != nil {
		// Create the file
		fmt.Printf("Didn't find a config file at %s, so I created one!", configPath)
		newCfgFile, createErr := os.Create(configPath)
		if createErr != nil {
			return nil, createErr
		}

		defer newCfgFile.Close()
	}

	// Parse config
	var config Config

	if cfgErr := hclsimple.DecodeFile(configPath, nil, &config); cfgErr != nil {
		fmt.Printf("failed to decode config: %s", cfgErr)
		os.Exit(2)
	}

	if config.Shell == "" {
		config.Shell = "bash"
	}

	return &config, nil
}

// Use provides commands to set appropriate environment variables
func (s *Stack) Use(shell string) string {
	var exportCommand string

	// Include Stack Name as an environment variable
	exportCommand += strings.Join([]string{fmt.Sprintf("export HCTX_STACK_NAME=%s\n", s.Name)}, exportCommand)

	if s.Nomad != nil {
		exportCommand += strings.Join(s.Nomad.Use(shell), exportCommand)
	}

	if s.Consul != nil {
		exportCommand += strings.Join(s.Consul.Use(shell), exportCommand)
	}

	if s.Vault != nil {
		exportCommand += strings.Join(s.Vault.Use(shell), exportCommand)
	}

	return exportCommand
}

// Unset Provides shell commands to unset environment variables
func (s *Stack) Unset(shell string) string {

	var unsetCommand string

	// Remove Stack environment variables
	unsetCommand += strings.Join([]string{"unset HCTX_STACK_NAME\n"}, unsetCommand)

	if s.Nomad != nil {
		unsetCommand += strings.Join(s.Nomad.Unset(shell), unsetCommand)
	}

	if s.Consul != nil {
		unsetCommand += strings.Join(s.Consul.Unset(shell), unsetCommand)
	}

	if s.Vault != nil {
		unsetCommand += strings.Join(s.Vault.Unset(shell), unsetCommand)
	}

	return unsetCommand
}

// Map Provides the stacks in a map for easier use in some use cases
func (c *Config) Map() map[string]Stack {

	stacks := make(map[string]Stack)

	for _, stack := range c.Stacks {
		stacks[stack.Name] = stack
	}

	return stacks
}

// StackExists Checks if the stack exists in the current configuration
func (c *Config) StackExists(name string) *Stack {

	if stack, stackExists := c.Map()[name]; stackExists {
		return &stack
	}

	return nil
}
