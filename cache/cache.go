package cache

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Shackelford-Arden/hctx/config"
	"github.com/Shackelford-Arden/hctx/models"
)

const FilePerms = os.FileMode(0600)
const FileName = "cache.json"

func CachePath() (string, error) {

	// Get user homedir
	userHome, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", fmt.Errorf("failed to get user homedir: %s", homeErr)
	}

	return fmt.Sprintf("%s/%s/%s/%s", userHome, config.ConfigParentDir, config.ConfigDir, FileName), nil
}

type Cache map[string]models.StackCache

func NewCache(cachePath string) (*Cache, error) {

	if cachePath == "" {
		cp, err := CachePath()
		if err != nil {
			return nil, fmt.Errorf("failed getting cache path: %s", err.Error())
		}
		cachePath = cp
	}

	cachePathStat, err := os.Stat(cachePath)
	if os.IsNotExist(err) {
		cacheCreated, createErr := os.Create(cachePath)
		if createErr != nil {
			return nil, fmt.Errorf("failed to create %s: %s", cachePath, createErr)
		}

		emptyCache := []byte("{}")
		if _, err := cacheCreated.Write(emptyCache); err != nil {
			return nil, fmt.Errorf("failed to write empty cache to %s: %s", cachePath, err)
		}
	}

	// Set appropriate permissions
	if cachePathStat == nil {
		cachePathStat, _ = os.Stat(cachePath)
	}

	currentPerm := cachePathStat.Mode().Perm()
	if currentPerm != FilePerms {
		setPermErr := os.Chmod(cachePath, FilePerms)
		if setPermErr != nil {
			return nil, fmt.Errorf("failed to set permissions on %s: %s", cachePath, setPermErr)
		}
	}

	cacheFile, _ := os.ReadFile(cachePath)
	var cache *Cache

	cacheParseErr := json.Unmarshal(cacheFile, &cache)
	if cacheParseErr != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %s", cachePath, cacheParseErr)
	}

	return cache, nil
}

func (c *Cache) Update(stackName string, data models.StackCache) error {
	(*c)[stackName] = data
	saveErr := c.Save("")
	if saveErr != nil {
		return fmt.Errorf("failed to update cache: %s", saveErr.Error())
	}

	return nil
}

// GetStack retrieves the cache for the given stack.
func (c *Cache) GetStack(stackName string) *models.StackCache {
	var cacheStack *models.StackCache

	for name, cache := range *c {
		if name == stackName {
			cacheStack = &cache
			break
		}
	}

	return cacheStack
}

// Get retrieves the full cache blob, typically for displaying.
func (c *Cache) Get() (map[string]models.StackCache, error) {

	cache := map[string]models.StackCache{}

	for stackName, stackCache := range *c {
		cache[stackName] = stackCache
	}

	return cache, nil
}

func (c *Cache) Save(path string) error {

	cp := path

	if path == "" {
		cachePath, err := CachePath()
		if err != nil {
			return fmt.Errorf("failed getting cache path: %s", err.Error())
		}
		cp = cachePath
	}

	cacheData, err := json.MarshalIndent(*c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed marshalling cache data: %s", err.Error())
	}

	writeErr := os.WriteFile(cp, cacheData, FilePerms)
	if writeErr != nil {
		return fmt.Errorf("failed writing cache data to %s: %s", cp, writeErr)
	}

	return nil
}
