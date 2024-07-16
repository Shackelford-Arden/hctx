package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Shackelford-Arden/hctx/types"
)

func TestListOutput(t *testing.T) {
	app, _ := App()

	testTmpDir := t.TempDir()

	t.Run("TestIndicatorIncluded", func(t *testing.T) {

		tmpConfig, _ := os.CreateTemp(testTmpDir, "indicator-*.hcl")
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

		// Set the environment variable to the stack name
		t.Setenv(types.StackNameEnv, "test-01")

		args := []string{"hctx", "--config", tmpConfig.Name(), "list"}
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

		expectedContent := strings.TrimSpace(`
Stacks:
  test-01 *
`,
		)

		if !strings.Contains(output, expectedContent) {
			t.Errorf("Expected output to contain '%s', but got: %s", expectedContent, output)
		}
	},
	)

	t.Run("TestNoStackSelected", func(t *testing.T) {

		tmpConfig, _ := os.CreateTemp(testTmpDir, "indicator-*.hcl")
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

		args := []string{"hctx", "--config", tmpConfig.Name(), "list"}
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

		expectedContent := strings.TrimSpace(`
Stacks:
  test-01
`,
		)

		if !strings.Contains(output, expectedContent) {
			t.Errorf("Expected output to contain '%s', but got: %s", expectedContent, output)
		}
	},
	)
}
