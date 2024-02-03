package models

import "testing"

func TestNomadAddressExportBash(t *testing.T) {

	nCfg := NomadConfig{
		Address: "http://localhost:4646",
	}

	nExport := nCfg.Use("bash")

	expectedResult := "export NOMAD_ADDR=http://localhost:4646\n"
	if nExport[0] != expectedResult {
		t.Errorf("Expected: '%s' | Receieved: '%s'", expectedResult, nExport[0])
	}
}

func TestNomadAddressExportZsh(t *testing.T) {

	nCfg := NomadConfig{
		Address: "http://localhost:4646",
	}

	nExport := nCfg.Use("zsh")

	expectedResult := "export NOMAD_ADDR=http://localhost:4646\n"
	if nExport[0] != expectedResult {
		t.Errorf("Expected: '%s' | Receieved: '%s'", expectedResult, nExport[0])
	}
}
