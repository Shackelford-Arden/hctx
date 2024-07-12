package cmd

import (
	"fmt"

	"github.com/Shackelford-Arden/hctx/cache"
	"github.com/Shackelford-Arden/hctx/models"
	"github.com/urfave/cli/v2"
)

// Use sets the appropriate environment variables to use
// the selected Stack.
// If caching is disabled, this will also unset any
// existing env vars for tokens (ie NOMAD_TOKEN).
func Use(ctx *cli.Context) error {

	stackName := ctx.Args().First()

	selectedStack := AppConfig.GetStack(stackName)

	if selectedStack == nil {
		return fmt.Errorf("no stack named %s in config", stackName)
	}

	currentStack := AppConfig.GetCurrentStack()
	// If caching is enabled, set to cache.
	if currentStack != nil && AppConfig.CacheAuth {
		toCache := cache.GetCacheableValues()
		updateErr := AppCache.Update(currentStack.Name, toCache)
		if updateErr != nil {
			return fmt.Errorf("could not update cache for stack %s: %v", currentStack.Name, updateErr)
		}
	}

	useOut := unsetTokens(AppConfig.Shell)
	var stackCache *models.StackCache
	if AppConfig.CacheAuth {
		stackCache = AppCache.Get(selectedStack.Name)
	}

	useOut += selectedStack.Use(AppConfig.Shell, stackCache, AppConfig.CacheAuth)

	fmt.Println(useOut)

	return nil
}
