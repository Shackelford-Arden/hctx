package cmd

import (
	"fmt"

	"github.com/Shackelford-Arden/hctx/cache"
	"github.com/urfave/cli/v2"
)

func Use(ctx *cli.Context) error {

	stackName := ctx.Args().First()

	selectedStack := AppConfig.GetStack(stackName)

	if selectedStack == nil {
		return fmt.Errorf("no stack named %s in config", stackName)
	}

	currentStack := AppConfig.GetCurrentStack()
	// Get current stacks tokens, if any and cache them
	fmt.Printf("Value of cache_auth: %v", AppConfig.CacheAuth)
	if currentStack != nil && AppConfig.CacheAuth == true {
		toCache := cache.GetCacheableValues()
		updateErr := AppCache.Update(currentStack.Name, toCache)
		if updateErr != nil {
			return fmt.Errorf("could not update cache for stack %s: %v", currentStack.Name, updateErr)
		}
	}

	// rehydrate env w/ new stack cache, if present
	newStackCache := AppCache.Get(selectedStack.Name)

	fmt.Print(selectedStack.Use(AppConfig.Shell, newStackCache, AppConfig.CacheAuth))

	return nil
}
