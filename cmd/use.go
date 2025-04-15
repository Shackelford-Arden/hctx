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
	currentTokens := cache.GetCacheableValues()

	// If caching is enabled, cache current stack tokens.
	if currentStack != nil && AppConfig.CacheAuth {
		currentCache := AppCache.GetStack(currentStack.Name)

		if currentCache == nil {
			currentCache = &models.StackCache{}
		}

		if currentTokens.NomadToken != "" {
			validToken := validNomadToken(selectedStack.Nomad.Address, currentTokens.NomadToken)
			if validToken {
				currentCache.NomadToken = currentTokens.NomadToken
			}
		}

		if currentTokens.ConsulToken != "" {
			validToken := validConsulToken(selectedStack.Consul.Address, currentTokens.ConsulToken)
			if validToken {
				currentCache.ConsulToken = currentTokens.ConsulToken
			}
		}
		updateErr := AppCache.Update(currentStack.Name, currentTokens)
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

	// Attempt to use the active token on the newly selected stack
	if selectedStack.Nomad != nil && currentTokens.NomadToken != "" && AppConfig.ShareNomadToken {

		// This handles:
		// 1) If the user has caching disabled (which is the default)
		// 2) If there was no cache for the stack
		if currentStackCache == nil {
			currentStackCache = &models.StackCache{}
		}

		// Check and see if the current token is valid on the new stack
		validTokenOnSelected := validNomadToken(selectedStack.Nomad.Address, currentTokens.NomadToken)
		if validTokenOnSelected {
			currentStackCache.NomadToken = currentTokens.NomadToken
		}
	}

	useOut := ActiveShell.UseOutput(selectedStack.Use(currentStackCache, AppConfig.CacheAuth))

	fmt.Println(useOut)

	return nil
}
