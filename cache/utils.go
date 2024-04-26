package cache

import (
	"github.com/Shackelford-Arden/hctx/models"
	"os"
)

func GetCacheableValues() models.StackCache {
	var currentValues models.StackCache

	nt := os.Getenv(models.NomadToken)
	if nt != "" {
		currentValues.NomadToken = nt
	}

	ct := os.Getenv(models.ConsulToken)
	if nt != "" {
		currentValues.ConsulToken = ct
	}

	return currentValues
}
