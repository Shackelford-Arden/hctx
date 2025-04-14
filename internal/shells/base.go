package shells

type Shell interface {
	// Activate returns the script/function definition that is used
	// as a wrapper around hctx to ensure that the shell is handling
	// the environment variables, since we can't manipulate them directly
	// via Go.
	Activate(string) string

	// UseOutput generates a string of text that the script from Activate parses
	// to correctly set environment variables with their correct values.
	UseOutput(envVars map[string]string) string

	// UnsetOutput generates a string of text that the script from Activate parses
	// to correctly unset/delete/remove environment variables.
	UnsetOutput(envVars []string) string
}
