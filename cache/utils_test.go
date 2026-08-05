package cache

import (
	"testing"

	"github.com/Shackelford-Arden/hctx/models"
)

func TestGetCacheableValues(t *testing.T) {

	tests := []struct {
		name           string
		nomadToken     string
		consulToken    string
		expectedNomad  string
		expectedConsul string
	}{
		{
			name:           "both tokens set",
			nomadToken:     "nomad-tok",
			consulToken:    "consul-tok",
			expectedNomad:  "nomad-tok",
			expectedConsul: "consul-tok",
		},
		{
			name:           "only consul token set",
			nomadToken:     "",
			consulToken:    "consul-tok",
			expectedNomad:  "",
			expectedConsul: "consul-tok",
		},
		{
			name:           "only nomad token set",
			nomadToken:     "nomad-tok",
			consulToken:    "",
			expectedNomad:  "nomad-tok",
			expectedConsul: "",
		},
		{
			name:           "neither token set",
			nomadToken:     "",
			consulToken:    "",
			expectedNomad:  "",
			expectedConsul: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(models.NomadToken, tt.nomadToken)
			t.Setenv(models.ConsulToken, tt.consulToken)

			got := GetCacheableValues()

			if got.NomadToken != tt.expectedNomad {
				t.Errorf("NomadToken = %q, want %q", got.NomadToken, tt.expectedNomad)
			}

			if got.ConsulToken != tt.expectedConsul {
				t.Errorf("ConsulToken = %q, want %q", got.ConsulToken, tt.expectedConsul)
			}
		})
	}
}
