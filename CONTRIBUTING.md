# Contributing

## Tooling

This project uses [`mise`](https://mise.jdx.dev) to manage tool versions, pinned in `mise.toml`. Run `mise install` to
set up everything needed for local development.

Tools managed by `mise`:

* `go`
* `changie`
  * Used for managing the generation of the Changelog.
* `task` - [Taskfile](https://taskfile.dev)
  * Conceptually similar to how Makefiles are used in many projects
* `golangci-lint`
  * Used for linting, run via `task check`
* `goreleaser`
  * Used for building/publishing releases

## Pull Requests

Nothing too crazy here:

* CI/tests are expected to pass. New tests are not always required, but are preferred and may be requested.
* Use `changie` to provide changelog items