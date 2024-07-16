package cmd

import (
	"fmt"

	"github.com/Shackelford-Arden/hctx/cache"
	"github.com/urfave/cli/v2"
)

// Unset Remove everything hctx configured in the environment variables
func Unset(ctx *cli.Context) error {

	currentStack := AppConfig.GetCurrentStack()
	configPath := ctx.String("config")

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

	fmt.Println(currentStack.Unset(AppConfig.Shell))

	return nil

}
