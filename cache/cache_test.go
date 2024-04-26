package cache

import (
	"os"
	"testing"
)

func TestNewEmptyCache(t *testing.T) {

	cachePath := "testdata/empty-cache.json"

	cache, err := NewCache(cachePath)
	if err != nil {
		t.Fatal(err)
	}

	if cache == nil {
		t.Fatal("cache is nil when it shouldn't be")
	}
}

func TestCreateNonExistentCacheFile(t *testing.T) {

	tmpCachePath := "/tmp/missing-cache.json"
	defer os.Remove(tmpCachePath)

	cache, err := NewCache(tmpCachePath)
	if err != nil {
		t.Fatal(err)

	}

	// Validate permissions are set correctly
	cacheStat, err := os.Stat(tmpCachePath)
	if err != nil {
		t.Fatal(err)
	}

	if cacheStat.Mode().Perm() != FilePerms {
		t.Fatalf("cache file %s has an invalid permissions %d", tmpCachePath, cacheStat.Mode())
	}

	if cache == nil {
		t.Fatal("cache should not be nil")
	}
}

func TestValidCache(t *testing.T) {
	cachePath := "testdata/valid-cache.json"

	cache, err := NewCache(cachePath)
	if err != nil {
		t.Fatal(err)
	}

	cacheItem := cache.Get("test")

	if cacheItem == nil {
		t.Fatal("cache item is nil when it shouldn't be")
	}

	if cacheItem.NomadToken != "test-token" && cacheItem.ConsulToken != "" {
		t.Fatal("cache item is not valid")
	}
}

func TestMissingCacheItem(t *testing.T) {
	cachePath := "testdata/valid-cache.json"

	cache, err := NewCache(cachePath)
	if err != nil {
		t.Fatal(err)
	}

	cacheItem := cache.Get("fake-test")

	if cacheItem != nil {
		t.Fatal("cached item should be nil, as fake-test should be missing.")
	}
}
