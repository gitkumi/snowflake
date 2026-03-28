# Repository Guidelines

## Project Structure & Module Organization
Snowflake is a Go CLI generator; `main.go` wires the Cobra tree while `cmd/cli/main.go` exposes the packaged entry point. Command implementations live in `internal/command/{generate,run,tui,version}` and scaffolding assets reside under `internal/initialize` with reusable fragments in `internal/initialize/template`. Tests mirror their packages inside `internal/...`, and build outputs land in `bin/` or `dist/`. Target Go 1.25.0 per `go.mod` and `.tool-versions`.

## Build, Test, and Development Commands
Key targets: `make run` executes the CLI; `make build` drops `bin/main`, while `make build-all` cross-compiles into `dist/` with checksums. Use `make test` (or `make test.coverage` for HTML output via `gotestsum`) before committing. `make audit` runs `go mod tidy -diff`, vet, staticcheck, and govulncheck—CI enforces it on every push. `make tidy` refreshes modules and formatting.

## Coding Style & Naming Conventions
Stick to idiomatic Go formatting (`gofmt`, tabs, camelCase) and run formatters before committing. Exported identifiers need doc comments and packages stay short nouns (`internal/command/run`). Keep new implementations under `internal/...` unless they are public, follow the existing Cobra file split of `command.go` plus verb-specific files, and run `make audit` before pushing to catch lint and security issues.

## Testing Guidelines
Place `_test.go` files next to the code they cover (e.g., `internal/initialize/initialize_test.go`) and use `TestXxx` naming with table-driven cases where branching exists. Execute `make test` routinely; inspect `coverage.html` from `make test.coverage` when touching scaffolding logic. Prefer shared fixtures in `internal/initialize/template` over duplicating test data.

## Commit & Pull Request Guidelines
Write short, imperative commit subjects (<60 chars) and append issue refs like `(#54)` when relevant; avoid WIP chains by squashing locally. PR descriptions should state motivation, summarize changes, list tests run (`make audit`, `make test`), and flag follow-up work. Include screenshots or recordings when modifying TUI flows and wait for green CI before requesting review.

## Tooling & Environment Setup
Install Go 1.25.0 via `asdf` (see `.tool-versions`). Ensure `sqlc`, `templ`, and `gotestsum` are available—mirror CI with `go install github.com/...`. Configure your editor to run `gofmt` on save and surface `go vet` diagnostics.
