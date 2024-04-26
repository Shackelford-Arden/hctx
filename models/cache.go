package models

type StackCache struct {
	NomadToken  string `json:"nomad-token,omitempty"`
	ConsulToken string `json:"consul-token,omitempty"`
}
