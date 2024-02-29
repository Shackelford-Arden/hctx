package cmd

import (
	"fmt"
	"github.com/Shackelford-Arden/hctx/models"
	"github.com/Shackelford-Arden/hctx/types"
	"github.com/urfave/cli/v2"
	"os"
)

// Unset Remove everything hctx configured in the environment variables
func Unset(ctx *cli.Context) error {

	currStack := os.Getenv(types.StackNameEnv)

	if currStack == "" {
		return nil
	}

	// Parse config
	cfg, cfgErr := models.NewConfig("")
	if cfgErr != nil {
		return cfgErr
	}

	stack := cfg.StackExists(currStack)
	if stack == nil {
		return fmt.Errorf("stack %s doesn't exist, no action taken", currStack)
	}

	fmt.Println(stack.Unset(cfg.Shell))

	return nil

}
