package models

import (
	"fmt"
	"github.com/Shackelford-Arden/hctx/types"
	"testing"
)

// region Bash/ZSH

func TestNewConfigUseFullNomad(t *testing.T) {

	cfg, err := NewConfig("testdata/fullNomad.hcl")
	if err != nil {
		t.Fatalf("failed to read config: %s", err)
	}

	stackName := "just_nomad"
	useOut := cfg.StackExists(stackName).Use("bash")

	expectedOutput := fmt.Sprintf(`
export %s='%s'
export %s=http://localhost:4646
export %s=default`,
		types.StackNameEnv, stackName, NomadAddr, NomadNamespace)

	if useOut != expectedOutput {
		t.Fatalf("\nExpected: %s\nActual: %s", expectedOutput, useOut)
	}

}

func TestNewConfigUseFullConsul(t *testing.T) {

	cfg, err := NewConfig("testdata/fullConsul.hcl")
	if err != nil {
		t.Fatalf("failed to read config: %s", err)
	}
	stackName := "just_consul"
	useOut := cfg.StackExists(stackName).Use("bash")

	expectedOutput := fmt.Sprintf(`
export %s='%s'
export %s=http://localhost:8500
export %s=default`,
		types.StackNameEnv, stackName, ConsulAddr, ConsulNamespace)

	if useOut != expectedOutput {
		t.Fatalf("\nExpected: %s\nActual: %s", expectedOutput, useOut)
	}

}

func TestNewConfigUseFullVault(t *testing.T) {

	cfg, err := NewConfig("testdata/fullVault.hcl")
	if err != nil {
		t.Fatalf("failed to read config: %s", err)
	}

	stackName := "just_vault"
	useOut := cfg.StackExists(stackName).Use("bash")

	expectedOutput := fmt.Sprintf(`
export %s='%s'
export %s=http://localhost:8200
export %s=default`,
		types.StackNameEnv, stackName, VaultAddr, VaultNamespace)

	if useOut != expectedOutput {
		t.Fatalf("\nExpected: %s\nActual: %s", expectedOutput, useOut)
	}

}

func TestNewConfigUnsetNomad(t *testing.T) {

	cfg, err := NewConfig("testdata/fullNomad.hcl")
	if err != nil {
		t.Errorf("failed to read config: %s", err)
	}

	stackName := "just_nomad"
	unsetOut := cfg.StackExists(stackName).Unset("bash")

	expectedOutput := fmt.Sprintf(`
unset %s
unset %s
unset %s`,
		types.StackNameEnv, NomadAddr, NomadNamespace)

	if unsetOut != expectedOutput {
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

	existingStack := testConfig.StackExists("test-stack")

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

	existingStack := testConfig.StackExists("Test Stack")

	if existingStack == nil {
		t.Fatalf("Expected stack to exist, but it does not.")
	}
}
