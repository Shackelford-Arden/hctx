package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Shackelford-Arden/hctx/build"
	"github.com/Shackelford-Arden/hctx/cache"
	"github.com/Shackelford-Arden/hctx/config"
)

var AppConfig *config.Config
var AppCache *cache.Cache

func ValidateConfig(ctx *cli.Context) error {

	userConfig := ctx.String("config")
	if userConfig == "" {

		userHome, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return fmt.Errorf("failed to get user homedir: %s", homeErr)
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
			fmt.Println(fmt.Sprintf("both %s and %s exist. Only using %s, please merge the config files then remove %s", configPath, configOldPath, configPath, configOldPath))
		}

		if oldConfig != nil && os.IsNotExist(newConfigStatErr) {

			// Copy old config to new config path
			copyErr := os.Rename(configOldPath, configFilePath)
			if copyErr != nil {
				return fmt.Errorf("failed to copy %s to %s: %s", configOldPath, configFilePath, copyErr)
			}
		}

		userConfig = configFilePath
	}

	// Parse config
	cfg, cfgErr := config.NewConfig(userConfig)
	if cfgErr != nil {
		return cfgErr
	}

	// Get Cache
	cacheItem, cacheErr := cache.NewCache("")
	if cacheErr != nil {
		return cacheErr
	}

	AppConfig = cfg
	AppCache = cacheItem

	return nil
}

func App() (*cli.App, error) {

	app := &cli.App{
		Name:        "Hashi Context",
		Usage:       "Managing your Hashi contexts with style!",
		HelpName:    "hctx",
		Description: "A CLI tool to help you manage your CLI life interacting with some of HashiCorp's products.",
		Version:     fmt.Sprintf("%s - %s - built with %s on %s", build.Version, build.Commit, build.BuiltWith, build.Date),
		Authors: []*cli.Author{
			{
				Name:  "Arden Shackelford",
				Email: "arden@ardens.tech",
			},
		},
		Before: ValidateConfig,
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
			{
				Name:    "cache",
				Aliases: []string{"c"},
				Usage:   "Interact with the cache.",
				Subcommands: []*cli.Command{
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
				Name:  "self",
				Usage: "Actions for interacting with hctx itself.",
				Subcommands: []*cli.Command{
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
