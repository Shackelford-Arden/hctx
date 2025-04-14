package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/Shackelford-Arden/hctx/cache"
	"github.com/Shackelford-Arden/hctx/models"
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

	// If caching is enabled, cache current stack tokens.
	if currentStack != nil && AppConfig.CacheAuth {
		currentCache := AppCache.GetStack(currentStack.Name)

		if currentCache == nil {
			currentCache = &models.StackCache{}
		}

		toCache := cache.GetCacheableValues()
		if toCache.NomadToken != "" {
			validToken := validNomadToken(selectedStack.Nomad.Address, toCache.NomadToken)
			if validToken {
				currentCache.NomadToken = toCache.NomadToken
			}
		}

		if toCache.ConsulToken != "" {
			validToken := validConsulToken(selectedStack.Consul.Address, toCache.ConsulToken)
			if validToken {
				currentCache.ConsulToken = toCache.ConsulToken
			}
		}
		updateErr := AppCache.Update(currentStack.Name, toCache)
		if updateErr != nil {
			return fmt.Errorf("could not update cache for stack %s: %v", currentStack.Name, updateErr)
		}
		_ = AppCache.Save("")
	}

	currentStackCache := AppCache.GetStack(selectedStack.Name)

	// Pull in cache of selected stack, clearing out
	// any that are invalid.
	if AppConfig.CacheAuth && currentStackCache != nil {

		cleanCache := models.StackCache{}

		// Validate if tokens have expired.
		cleanCache.NomadToken = currentStackCache.NomadToken
		if currentStackCache.NomadToken != "" && selectedStack.Nomad != nil {
			nomToken := validNomadToken(selectedStack.Nomad.Address, currentStackCache.NomadToken)
			if !nomToken {
				// Remove expired token from cache
				cleanCache.NomadToken = ""
			}
		}

		cleanCache.ConsulToken = currentStackCache.ConsulToken
		if currentStackCache.ConsulToken != "" && selectedStack.Consul != nil {
			conToken := validConsulToken(selectedStack.Consul.Address, currentStackCache.ConsulToken)
			if !conToken {
				// Remove expired token from cache
				cleanCache.ConsulToken = ""
			}
		}

		err := AppCache.Update(selectedStack.Name, cleanCache)
		if err != nil {
			return fmt.Errorf("could not update cache for stack %s: %v", selectedStack.Name, err)
		}

		// Set this so that correct values are pulled in
		// later when cached values are pulled in, if enabled
		currentStackCache = &cleanCache
	}

	useOut := ActiveShell.UseOutput(selectedStack.Use(currentStackCache, AppConfig.CacheAuth))

	fmt.Println(useOut)

	return nil
}
