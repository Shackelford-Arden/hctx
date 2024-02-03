package cmd

import (
	"fmt"
	"github.com/Shackelford-Arden/hctx/models"
	"github.com/urfave/cli/v2"
	"strings"
)

func Use(ctx *cli.Context) error {

	// Bail if more than one stack is provided.
	if ctx.Args().Len() > 1 {
		return fmt.Errorf("only 1 stack can be used at a time. You provided %d: %s", ctx.Args().Len(), strings.Join(ctx.Args().Slice(), ","))
	}

	stackName := ctx.Args().Get(0)
	var selectedStack *models.Stack

	// Parse config
	cfg, cfgErr := models.NewConfig()
	if cfgErr != nil {
		return cfgErr
	}

	for _, stack := range cfg.Stacks {
		if stack.Name == stackName {
			selectedStack = &stack
			break
		}
	}

	if selectedStack == nil {
		return fmt.Errorf("no stack named %s", stackName)
	}

	fmt.Print(selectedStack.Use(cfg.Shell))

	return nil
}
