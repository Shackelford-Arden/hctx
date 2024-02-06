package cmd

import (
	"github.com/urfave/cli/v2"
)

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
