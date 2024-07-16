package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Shackelford-Arden/hctx/cache"
	"github.com/Shackelford-Arden/hctx/models"
	"github.com/urfave/cli/v2"
)

// CleanCache checks all items in the cache file
// and removes any that have expired or are not valid,
// including stacks that are no longer in the config.
func CleanCache(ctx *cli.Context) error {

	newCache := cache.Cache{}
	currentCache, err := AppCache.Get()
	if err != nil {
		return fmt.Errorf("failed to get current cache: %s", err)
	}

	for stack, stackCache := range currentCache {

		// Check if cached stack is in config
		configStack := AppConfig.GetStack(stack)
		if configStack == nil {
			continue
		}

		cleanCache := models.StackCache{
			NomadToken:  stackCache.NomadToken,
			ConsulToken: stackCache.ConsulToken,
		}

		// Validate if tokens have expired.
		if stackCache.NomadToken != "" && configStack.Nomad != nil {
			nomToken := validNomadToken(configStack.Nomad.Address, stackCache.NomadToken)
			if !nomToken {
				// Remove expired token from cache
				cleanCache.NomadToken = ""
			}
		}

		if stackCache.ConsulToken != "" && configStack.Consul != nil {
			conToken := validNomadToken(configStack.Consul.Address, stackCache.ConsulToken)
			if !conToken {
				// Remove expired token from cache
				cleanCache.ConsulToken = ""
			}
		}

		newCache[stack] = cleanCache
	}

	AppCache = &newCache
	saveErr := AppCache.Save("")
	if saveErr != nil {
		return fmt.Errorf("failed to save cache: %s", saveErr)
	}

	return nil
}

// ClearCache removes all items from the cache file.
func ClearCache(ctx *cli.Context) error {

	AppCache = &cache.Cache{}
	saveErr := AppCache.Save("")
	if saveErr != nil {
		return fmt.Errorf("failed to save cleared cache: %s", saveErr)
	}

	return nil
}

// ShowCache shows the content of the cache file.
func ShowCache(ctx *cli.Context) error {

	currentCache, _ := AppCache.Get()
	fmtCache, fmtErr := json.MarshalIndent(currentCache, "", "  ")
	if fmtErr != nil {
		return fmt.Errorf("failed to read/format the cache: %s", fmtErr.Error())
	}

	fmt.Println(string(fmtCache))

	return nil
}
