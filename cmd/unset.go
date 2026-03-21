package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/Shackelford-Arden/hctx/cache"
)

// Unset Remove everything hctx configured in the environment variables
func Unset(_ context.Context, cmd *cli.Command) error {

	currentStack := AppConfig.GetCurrentStack()
	configPath := cmd.String("config")

	if currentStack == nil {
		return nil
	}

	// Get current stacks tokens, if any and cache them
	toCache := cache.GetCacheableValues()
	if AppConfig.CacheAuth {
		updateErr := AppCache.Update(currentStack.Name, toCache)
		if updateErr != nil {
			return fmt.Errorf("could not update cache for stack %s: %v", currentStack.Name, updateErr)
		}

		saveErr := AppCache.Save(configPath)
		if saveErr != nil {
			return fmt.Errorf("could not save cache for stack %s: %v", currentStack.Name, saveErr)
		}
	}

	fmt.Println(resolveShell(cmd).UnsetOutput(currentStack.Unset()))

	return nil

}
