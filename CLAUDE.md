# CLAUDE.md

Snowflake is a Go CLI that scaffolds opinionated Go web projects with selectable features.

## Commands

```
make test              # run all tests (uses gotestsum)
make test.coverage     # run tests with HTML coverage report
make build             # build to bin/main
make audit             # tidy check + vet + staticcheck + govulncheck
make tidy              # go mod tidy + go fmt
make run               # run the CLI
```

## Architecture

- `main.go` / `cmd/cli/main.go` — Cobra CLI wiring
- `internal/command/tui/` — interactive TUI (charmbracelet/huh)
- `internal/command/run/` — non-interactive CLI flags
- `internal/initialize/` — core generation logic
  - `types.go` — `Database`, `Queue`, `ContainerRuntime` enums
  - `initialize.go` — `Config` struct and `Run()` orchestrator
  - `project.go` — `Project` struct, file exclusions, file renames
  - `operations.go` — post-generation commands (go mod init, tidy, build, etc.)
  - `template/base/` — embedded Go template files (`.templ` suffix = snowflake template, stripped on output)
  - `template/fragments/database/` — per-database SQL migration/query fragments

## Adding a new feature

Follow the pattern of existing features (SMTP, Storage, Redis, Templ):

1. **Config + Project** — Add a `bool` field to `Config` in `initialize.go` and `Project` in `project.go`. Wire it in `NewProject`.
2. **TUI** — Add an option to the `huh.NewMultiSelect` in `internal/command/tui/command.go`.
3. **CLI** — Add a `--flag` in `internal/command/run/command.go`.
4. **Template files** — Create files under `template/base/` with `.templ` suffix. Use `{{- if .FeatureName }}` guards in shared templates (server.go, main.go, env.go, Makefile, .env, etc.).
5. **File exclusions** — Add a `FileExclusion` entry in `project.go` so files are skipped when the feature is disabled.
6. **Post-commands** — If the feature needs a build step (like `templ generate` or `sqlc generate`), add it to `operations.go`. Order matters: generation steps must run before `go mod tidy`.
7. **Tests** — Add `TestGenerateFeature` and `TestGenerateNoFeature` in `initialize_test.go`.

## Template naming

Snowflake uses `.templ` as its own template suffix, which gets stripped to produce the output filename. To generate a file that itself has a `.templ` extension (e.g. for the templ library), name it `.templ.templ`.

## Testing

Tests in `initialize_test.go` run the full generation pipeline (template rendering, go mod init, go mod tidy, build). They require `go`, `gofmt`, and `make` on PATH. Tests that enable database features also need `sqlc`. Tests that enable templ need the `templ` CLI.

All tests use `Quiet: true` and `Git: false` to keep them fast and side-effect-free.
