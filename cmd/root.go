package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/Shackelford-Arden/hctx/build"
	"github.com/Shackelford-Arden/hctx/cache"
	"github.com/Shackelford-Arden/hctx/config"
)

var AppConfig *config.Config
var AppCache *cache.Cache

func ValidateConfig(ctx context.Context, cmd *cli.Command) (context.Context, error) {

	userConfig := cmd.String("config")
	if userConfig == "" {

		userHome, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return ctx, fmt.Errorf("failed to get user homedir: %s", homeErr)
		}

		configPath := fmt.Sprintf("%s/%s/%s", userHome, config.ConfigParentDir, config.ConfigDir)
		configFilePath := fmt.Sprintf("%s/%s/%s/%s", userHome, config.ConfigParentDir, config.ConfigDir, config.ConfigFileName)
		configOldPath := fmt.Sprintf("%s/%s/%s", userHome, config.ConfigParentDir, config.OldConfigFileName)

		_, err := os.Stat(configPath)
		if os.IsNotExist(err) {
			// Create the directory
			err := os.Mkdir(configPath, 0744)
			if err != nil {
				return ctx, fmt.Errorf("failed to create %s: %s", configPath, err)
			}
		}

		oldConfig, _ := os.Stat(configOldPath)
		newConfig, newConfigStatErr := os.Stat(configFilePath)

		if oldConfig != nil && newConfig != nil {
			fmt.Printf("both %s and %s exist. Only using %s, please merge the config files then remove %s\n", configPath, configOldPath, configPath, configOldPath)
		}

		if oldConfig != nil && os.IsNotExist(newConfigStatErr) {

			// Copy old config to new config path
			copyErr := os.Rename(configOldPath, configFilePath)
			if copyErr != nil {
				return ctx, fmt.Errorf("failed to copy %s to %s: %s", configOldPath, configFilePath, copyErr)
			}
		}

		userConfig = configFilePath
	}

	// Parse config
	cfg, cfgErr := config.NewConfig(userConfig)
	if cfgErr != nil {
		return ctx, cfgErr
	}

	// Get Cache
	cacheItem, cacheErr := cache.NewCache("")
	if cacheErr != nil {
		return ctx, cacheErr
	}

	AppConfig = cfg
	AppCache = cacheItem

	return ctx, nil
}

func App() (*cli.Command, error) {

	app := &cli.Command{
		Name:        "Hashi Context",
		Usage:       "Managing your Hashi contexts with style!",
		Description: "A CLI tool to help you manage your CLI life interacting with some of HashiCorp's products.",
		Version:     fmt.Sprintf("%s - %s - built with %s on %s", build.Version, build.Commit, build.BuiltWith, build.Date),
		Before:      ValidateConfig,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to config file to use.",
				Hidden:  true,
			},
			&cli.StringFlag{
				Name:   "shell",
				Hidden: true,
			},
		},
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
				Action:  Use,
			},
			{
				Name:    "unset",
				Aliases: []string{"un"},
				Usage:   "Cleans up all managed environment variables for the current stack.",
				Action:  Unset,
			},
			{
				Name:    "activate",
				Aliases: []string{"a"},
				Action:  Activate,
				Usage:   "Used to generate the appropriate shell scripts to set environment variables.",
			},
			{
				Name:  "self",
				Usage: "Actions for interacting with hctx itself.",
				Commands: []*cli.Command{
					{
						Name:   "update",
						Usage:  "Will attempt to find the latest release and download it. Connectivity to Github is required!",
						Action: SelfUpdate,
					},
					{
						Name:   "show-path",
						Usage:  "Gives you the absolute path to the hctx binary.",
						Action: ShowPath,
					},
				},
			},
			{
				Name:    "cache",
				Aliases: []string{"c"},
				Usage:   "Interact with the cache.",
				Commands: []*cli.Command{
					{
						Name:    "show",
						Aliases: []string{"s"},
						Usage:   "ShowCache the current cache",
						Action:  ShowCache,
					},
					{
						Name:   "clear",
						Usage:  "Clears out the cache of all items.",
						Action: ClearCache,
					},
					{
						Name:   "clean",
						Usage:  "Checks all cached items and removes those that have expired.",
						Action: CleanCache,
					},
				},
			},
			{
				Name:   "version",
				Usage:  "Display the version of hctx",
				Action: ShowVersion,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "check-latest",
						Usage: "Will attempt to check Github and see what the latest version is",
						Value: false,
					},
				},
			},
		},
	}

	return app, nil
}
