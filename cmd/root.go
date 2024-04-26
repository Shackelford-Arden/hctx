package cmd

import (
	"fmt"
	"github.com/Shackelford-Arden/hctx/cache"
	"github.com/Shackelford-Arden/hctx/config"
	"github.com/urfave/cli/v2"
	"log/slog"
	"os"
)

var AppConfig *config.Config
var AppCache *cache.Cache

func ValidateConfig(ctx *cli.Context) error {

	userHome, homeErr := os.UserHomeDir()
	if homeErr != nil {
		fmt.Printf("failed to get user homedir: %s", homeErr)
		os.Exit(10)
	}

	configPath := fmt.Sprintf("%s/%s/%s", userHome, config.ConfigParentDir, config.ConfigDir)
	configFilePath := fmt.Sprintf("%s/%s/%s/%s", userHome, config.ConfigParentDir, config.ConfigDir, config.ConfigFileName)
	configOldPath := fmt.Sprintf("%s/%s/%s", userHome, config.ConfigParentDir, config.OldConfigFileName)

	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		// Create the directory
		err := os.Mkdir(configPath, 0744)
		if err != nil {
			return fmt.Errorf("failed to create %s: %s", configPath, err)
		}
	}

	oldConfig, _ := os.Stat(configOldPath)
	newConfig, newConfigStatErr := os.Stat(configFilePath)

	if oldConfig != nil && newConfig != nil {
		slog.InfoContext(ctx.Context, fmt.Sprintf("both %s and %s exist. Only using %s, please merge the config files then remove %s", configPath, configOldPath, configPath, configOldPath))
	}

	if oldConfig != nil && os.IsNotExist(newConfigStatErr) {

		// Copy old config to new config path
		copyErr := os.Rename(configOldPath, configFilePath)
		if copyErr != nil {
			return fmt.Errorf("failed to copy %s to %s: %s", configOldPath, configFilePath, copyErr)
		}
	}

	// Parse config
	cfg, cfgErr := config.NewConfig("")
	if cfgErr != nil {
		return cfgErr
	}

	// Get Cache
	cache, cacheErr := cache.NewCache("")
	if cacheErr != nil {
		return cacheErr
	}

	AppConfig = cfg
	AppCache = cache

	return nil
}

func App() (*cli.App, error) {

	app := &cli.App{
		Name:        "Hashi Context",
		HelpName:    "hctx",
		Description: "A CLI tool to help you manage your CLI life interacting with some of HashiCorp's products.",
		Authors: []*cli.Author{
			{
				Name:  "Arden Shackelford",
				Email: "arden@ardens.tech",
			},
		},
		Before: ValidateConfig,
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List configured stacks",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "detailed",
						Aliases: []string{"d"},
					},
				},
				Action: List,
			},
			{
				Name:    "use",
				Aliases: []string{"u"},
				Usage:   "Use the selected stack as your current stack.",
				Args:    true,
				Action:  Use,
			},
			{
				Name:    "unset",
				Aliases: []string{"un"},
				Usage:   "Cleans up all managed environment variables for the current stack.",
				Args:    false,
				Action:  Unset,
			},
			{
				Name:    "activate",
				Aliases: []string{"a"},
				Action:  Activate,
				Usage:   "Used to generate the appropriate shell scripts to set environment variables.",
				Args:    true,
			},
		},
	}

	return app, nil
}
