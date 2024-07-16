package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestUse(t *testing.T) {
	app, _ := App()

	testTmpDir := t.TempDir()

	t.Run("TestErrOnMissingStack", func(t *testing.T) {

		tmpConfig, _ := os.CreateTemp(testTmpDir, "missing-stack-*.hcl")
		defer tmpConfig.Close()

		args := []string{"hctx", "--config", tmpConfig.Name(), "use", "test-01"}
		err := app.Run(args)
		if err.Error() != "no stack named test-01 in config" {
			t.Errorf("Error received was not expected %v", err)
		}
	})

	t.Run("TestValidStackOutput", func(t *testing.T) {

		tmpConfig, _ := os.CreateTemp(testTmpDir, "missing-stack-*.hcl")
		defer tmpConfig.Close()

		tmpConfig.WriteString(`
stack "test-01" {
  nomad {
    address = "http://localhost:4646"
  }
}
`)

		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("Could not create pipe: %v", err)
		}

		// Save the original stdout and stderr
		oldStdout := os.Stdout
		oldStderr := os.Stderr

		// Set the pipe's writer as stdout and stderr
		os.Stdout = w
		os.Stderr = w

		// Create a channel to signal when we're done reading the output
		done := make(chan bool)
		var output string

		// Start a goroutine to read from the pipe
		// This allows us to capture the output of the CLI
		// being called in the next section.
		go func() {
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output = buf.String()
			done <- true
		}()

		args := []string{"hctx", "--shell", "bash", "--config", tmpConfig.Name(), "use", "test-01"}
		err = app.Run(args)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Restore the original stdout and stderr
		os.Stdout = oldStdout
		os.Stderr = oldStderr

		// Close the writer side of the pipe
		w.Close()

		// Wait for the reading goroutine to finish
		<-done

		// Check for command execution error
		if err != nil {
			t.Errorf("Command execution failed: %v", err)
		}

		expectedContent := strings.TrimSpace(fmt.Sprintf(`
export NOMAD_ADDR=http://localhost:4646
`,
		))

		if !strings.Contains(output, expectedContent) {
			t.Errorf("Expected output to contain '%s', but got: %s", expectedContent, output)
		}
	},
	)
}
