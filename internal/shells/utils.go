package shells

import (
	"os"
)

// DiscoverShell attempts to find out what the shell is for the
// current session.
func DiscoverShell() Shell {

	// Pull in environment variables that help indicate what shell
	// is currently active.

	// Using NU_VERSION to check for Nu
	// Nu does _not_ set SHELL
	// Nu may also be initiated from another shell, which would
	// leave SHELL set to the "old" shell.
	nuVersion := os.Getenv("NU_VERSION")
	if nuVersion != "" {
		return &Nushell{}
	}
	// SHELL is set by a number of shells, which is the
	// path to the shell binary itself.
	//shellEnvVar := os.Getenv("SHELL")

	// uncomment these when needing to add support
	// for more shells, like fish.
	//splitName := strings.Split(shellEnvVar, "/")
	//shellName := splitName[len(splitName)-1]

	return &Bash{}
}
