package cmd

import (
	"fmt"
	"github.com/Shackelford-Arden/hctx/cache"
	"github.com/urfave/cli/v2"
)

// Unset Remove everything hctx configured in the environment variables
func Unset(ctx *cli.Context) error {

	currentStack := AppConfig.GetCurrentStack()

	if currentStack == nil {
		return nil
	}

	// Get current stacks tokens, if any and cache them
	toCache := cache.GetCacheableValues()
	updateErr := AppCache.Update(currentStack.Name, toCache)
	if updateErr != nil {
		return fmt.Errorf("could not update cache for stack %s: %v", currentStack.Name, updateErr)
	}

	fmt.Println(currentStack.Unset(AppConfig.Shell))

	return nil

}
