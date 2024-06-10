package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/Shackelford-Arden/hctx/models"
	"github.com/Shackelford-Arden/hctx/types"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

const ConfigParentDir = ".config"
const ConfigDir = "hctx"
const ConfigFileName = "config.hcl"
const OldConfigFileName = ".hctx.hcl"

type Config struct {
	Stacks    []Stack `hcl:"stack,block"`
	Shell     string  `hcl:"shell,optional"`
	CacheAuth bool    `hcl:"cache_auth,optional"`
}

func NewConfig(cp string) (*Config, error) {

	var configPath = cp

	if cp == "" {
		// Get user homedir
		userHome, homeErr := os.UserHomeDir()
		if homeErr != nil {
			fmt.Printf("failed to get user homedir: %s", homeErr)
			os.Exit(10)
		}

		configPath = fmt.Sprintf("%s/%s/%s/%s", userHome, ConfigParentDir, ConfigDir, ConfigFileName)
	}

	// Check if there is a config file
	_, statErr := os.Stat(configPath)
	if statErr != nil {
		// Create the file
		fmt.Printf("Didn't find a config file at %s, so I created one!\n", configPath)
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

// Map Provides the stacks in a map for easier use in some use cases
func (c *Config) Map() map[string]Stack {

	stacks := make(map[string]Stack)

	for _, stack := range c.Stacks {
		stacks[stack.Name] = stack
		stacks[stack.Alias] = stack
	}

	return stacks
}

// GetStack Checks if the stack exists in the current configuration
func (c *Config) GetStack(name string) *Stack {

	if stack, stackExists := c.Map()[name]; stackExists {
		return &stack
	}

	return nil
}

// GetCurrentStack Checks current environment variable(s) to identify current stack, if any
func (c *Config) GetCurrentStack() *Stack {
	currentStackName := os.Getenv(types.StackNameEnv)
	if currentStackName != "" {
		for _, stack := range c.Stacks {
			if stack.Name == currentStackName || stack.Alias == currentStackName {
				return &stack
			}
		}
	}

	return nil
}

type Stack struct {
	Name   string               `hcl:",label"`
	Alias  string               `hcl:"alias,optional"`
	Nomad  *models.NomadConfig  `hcl:"nomad,block"`
	Consul *models.ConsulConfig `hcl:"consul,block"`
	Vault  *models.VaultConfig  `hcl:"vault,block"`
}

// Use provides commands to set appropriate environment variables
func (s *Stack) Use(shell string, cache *models.StackCache, useCache bool) string {
	// Include Stack Name as an environment variable
	// Allow the Alias name to show in the environment variable
	stackName := s.Name
	if s.Alias != "" {
		stackName = s.Alias
	}

	var nomadToken string
	var consulToken string

	if cache != nil && useCache {
		nomadToken = cache.NomadToken
		consulToken = cache.ConsulToken
	}

	var exportCommands = []string{fmt.Sprintf("\nexport %s='%s'", types.StackNameEnv, stackName)}

	if s.Nomad != nil {
		exportCommands = append(exportCommands, s.Nomad.Use(shell, nomadToken)...)
	}

	if s.Consul != nil {
		exportCommands = append(exportCommands, s.Consul.Use(shell, consulToken)...)
	}

	if s.Vault != nil {
		exportCommands = append(exportCommands, s.Vault.Use(shell)...)
	}

	var exportCommand = strings.Join(exportCommands, "\n")

	return exportCommand
}

// Unset Provides shell commands to unset environment variables
func (s *Stack) Unset(shell string) string {

	// Remove Stack environment variables
	var unsetCommands = []string{fmt.Sprintf("\nunset %s", types.StackNameEnv)}

	if s.Nomad != nil {
		unsetCommands = append(unsetCommands, s.Nomad.Unset(shell)...)
	}

	if s.Consul != nil {
		unsetCommands = append(unsetCommands, s.Consul.Unset(shell)...)
	}

	if s.Vault != nil {
		unsetCommands = append(unsetCommands, s.Vault.Unset(shell)...)
	}

	var unsetCommand = strings.Join(unsetCommands, "\n")

	return unsetCommand
}
