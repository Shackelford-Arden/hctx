package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestActivateOutput(t *testing.T) {
	app, _ := App()

	t.Run("TestBashOutput", func(t *testing.T) {

		execPath, _ := os.Executable()
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

		args := []string{"hctx", "--shell", "bash", "activate"}
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
hctx () {
  local command
  HCTX_PATH='%s'
  command="${1:-}"
  if [ "$#" = 0 ]
  then
    command $HCTX_PATH
    return
  fi
  shift
  case "$command" in
  (use|u|unset|un) if [[ ! " $@ " =~ " --help " ]] && [[ ! " $@ " =~ " -h " ]]
    then
      eval "$(command $HCTX_PATH "$command" "$@")"
      return $?
    fi ;;
  esac
  command $HCTX_PATH "$command" "$@"
}
`, execPath,
		))

		if !strings.Contains(output, expectedContent) {
			t.Errorf("Expected output to contain '%s', but got: %s", expectedContent, output)
		}
	},
	)
}
