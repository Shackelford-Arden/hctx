# AGENTS.md

## Build & Development Commands

```bash
# Run all tests
task test
# or directly
go test ./...

# Run a single test
go test -v ./cmd/ -run TestUse
go test -v ./config/ -run TestConfig
go test -v ./cache/ -run TestCache

# Run tests with coverage
go test -cover ./...

# Format code
task fmt
# or
go fmt ./...

# Lint and check (mod tidy, vet, staticcheck)
task check

# Pre-commit checks (fmt + check + test)
task pre-commit
# or
task pc

# Build development binary
task build-dev

# Build release snapshot
task build-test
```

## Code Style Guidelines

### Imports
- Group imports: stdlib, external packages, internal packages
- Use explicit imports; avoid dot imports
- Standard order: `fmt`, `os`, then external (`github.com/...`), then internal (`github.com/Shackelford-Arden/hctx/...`)

### Formatting
- Run `go fmt ./...` before committing
- Use tabs for indentation
- Follow standard Go formatting conventions

### Naming Conventions
- Use camelCase for private functions/variables
- Use PascalCase for exported functions/types
- Environment variable constants use PascalCase with descriptive names (e.g., `NomadAddr`)
- Package names are lowercase, single words when possible

### Error Handling
- Return errors with context using `fmt.Errorf("context: %v", err)`
- Don't log errors in library code; return them to caller
- Use meaningful error messages that describe what failed
- Check errors immediately after calls; avoid nested if/else when possible

### Types & Structs
- Define HCL tags for config structs: `hcl:"field_name,optional"`
- Use pointer receivers for methods that modify state
- Use value receivers for read-only methods on small structs
- Export types only when needed by other packages

### Comments
- Document exported functions with godoc-style comments
- Keep comments concise; explain why, not what
- Use inline comments sparingly; prefer clear variable names

### Testing
- Use table-driven tests where appropriate
- Name test files with `_test.go` suffix in same package
- Use `t.TempDir()` for temporary test directories
- Capture stdout/stderr with `os.Pipe()` when testing CLI output
- Use subtests with `t.Run()` for better organization

### File Organization
- `cmd/` - CLI command handlers
- `config/` - Configuration parsing
- `models/` - Data models and business logic
- `types/` - Shared type definitions and constants
- `internal/` - Internal packages (shells, github)
- `cache/` - Caching functionality
- `build/` - Build version info

### CLI Commands (urfave/cli)
- Use `Aliases` for shorter command versions
- Define flags with clear names and descriptions
- Use `Before` middleware for config validation
- Return errors from command handlers; don't call `os.Exit`

### Config Handling
- Default config location: `~/.config/hctx/config.hcl`
- Create config dir/file if not exists
- Support both HCL and legacy `.hctx.hcl` (with migration warning)
