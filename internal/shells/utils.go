package shells

// FromName maps a shell name (e.g. "nushell", "bash") to its concrete
// Shell implementation. Unrecognized or empty names fall back to Bash.
func FromName(name string) Shell {
	switch name {
	case "nushell":
		return &Nushell{}
	default:
		return &Bash{}
	}
}
