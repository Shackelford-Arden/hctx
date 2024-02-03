package cmd

import (
	"fmt"
	"github.com/Shackelford-Arden/hctx/models"
	"github.com/urfave/cli/v2"
)

func List(ctx *cli.Context) error {

	// Parse config
	cfg, cfgErr := models.NewConfig()
	if cfgErr != nil {
		return cfgErr
	}

	if len(cfg.Stacks) == 0 {
		fmt.Fprintf(ctx.App.Writer, "No stacks!\n")
		return nil
	}

	fmt.Println("Stacks:")
	for _, stack := range cfg.Stacks {
		fmt.Printf("  %s\n", stack.Name)

		if !ctx.Bool("verbose") {
			continue
		}

		if stack.Nomad != nil {
			fmt.Printf("    Nomad: %s\n", stack.Nomad.Address)
		}
		if stack.Consul != nil {
			fmt.Printf("    Consul: %s\n", stack.Consul.Address)
		}
		if stack.Vault != nil {
			fmt.Printf("    Vault: %s\n", stack.Vault.Address)
		}
	}

	return nil
}
