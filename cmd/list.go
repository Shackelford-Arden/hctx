package cmd

import (
	"fmt"
	"os"

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

	currStack := os.Getenv("HCTX_STACK_NAME")

	fmt.Println("Stacks:")
	for _, stack := range cfg.Stacks {
		var indicator string
		if stack.Name == currStack || stack.Alias == currStack {
			indicator = "*"
		}
		fmt.Printf("  %s %s\n", stack.Name, indicator)

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
