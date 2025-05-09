package config

import (
	"reflect"
	"testing"

	"github.com/Shackelford-Arden/hctx/models"
	"github.com/Shackelford-Arden/hctx/types"
)

// region Bash/ZSH

func TestNewConfigUseFullNomad(t *testing.T) {

	cfg, err := NewConfig("testdata/fullNomad.hcl")
	if err != nil {
		t.Fatalf("failed to read config: %s", err)
	}

	stackName := "just_nomad"
	useOut := cfg.GetStack(stackName).Use(nil, false)
	expectedOutput := map[string]string{
		types.StackNameEnv:    stackName,
		models.NomadAddr:      "http://localhost:4646",
		models.NomadNamespace: "default",
	}

	if !reflect.DeepEqual(useOut, expectedOutput) {
		t.Fatalf("\nExpected: %s\nActual: %s", expectedOutput, useOut)
	}
}

func TestNewConfigUseFullConsul(t *testing.T) {

	cfg, err := NewConfig("testdata/fullConsul.hcl")
	if err != nil {
		t.Fatalf("failed to read config: %s", err)
	}
	stackName := "just_consul"
	useOut := cfg.GetStack(stackName).Use(nil, false)

	expectedOutput := map[string]string{
		types.StackNameEnv:     stackName,
		models.ConsulAddr:      "http://localhost:8500",
		models.ConsulNamespace: "default",
	}

	if !reflect.DeepEqual(useOut, expectedOutput) {
		t.Fatalf("\nExpected: %s\nActual: %s", expectedOutput, useOut)
	}

}

func TestNewConfigUseFullVault(t *testing.T) {

	cfg, err := NewConfig("testdata/fullVault.hcl")
	if err != nil {
		t.Fatalf("failed to read config: %s", err)
	}

	stackName := "just_vault"
	useOut := cfg.GetStack(stackName).Use(nil, false)

	expectedOutput := map[string]string{
		types.StackNameEnv:    stackName,
		models.VaultAddr:      "http://localhost:8200",
		models.VaultNamespace: "default",
	}

	if !reflect.DeepEqual(useOut, expectedOutput) {
		t.Fatalf("\nExpected: %s\nActual: %s", expectedOutput, useOut)
	}

}

func TestNewConfigUnsetNomad(t *testing.T) {

	cfg, err := NewConfig("testdata/fullNomad.hcl")
	if err != nil {
		t.Errorf("failed to read config: %s", err)
	}

	stackName := "just_nomad"
	unsetOut := cfg.GetStack(stackName).Unset()

	expectedOutput := []string{types.StackNameEnv, models.NomadToken, models.NomadAddr, models.NomadNamespace}

	if !reflect.DeepEqual(unsetOut, expectedOutput) {
		t.Fatalf("\nExpected: %s\nActual: %s", expectedOutput, unsetOut)
	}

}

// endregion

func TestStackExistsUsingStackName(t *testing.T) {

	// Create test stack
	stack := Stack{
		Name:  "test-stack",
		Alias: "Test Stack",
	}

	testConfig := Config{
		Stacks: []Stack{stack},
	}

	existingStack := testConfig.GetStack("test-stack")

	if existingStack == nil {
		t.Fatalf("Expected stack to exist, but it does not.")
	}
}

func TestStackExistsUsingAlias(t *testing.T) {

	// Create test stack
	stack := Stack{
		Name:  "test-stack",
		Alias: "Test Stack",
	}

	testConfig := Config{
		Stacks: []Stack{stack},
	}

	existingStack := testConfig.GetStack("Test Stack")

	if existingStack == nil {
		t.Fatalf("Expected stack to exist, but it does not.")
	}
}
