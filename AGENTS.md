If the user asks you to rename the project, search and replace ALL instances of `metruzanca/go-cli` with `{username}/{repo}` (go.mod, imports, release config, README, AGENTS, etc.) and `GO_CLI` to `REPO` to rename enviroment variables. Ask the user if they want to use their github username + the name of the folder this template was cloned into. Then once renamed, delete this paragraph from the AGENTS.md.

## Project

Go CLI application built with [Cobra](https://github.com/spf13/cobra). Used as a starting template for new CLI projects.

## Conventions

Always-on rules that never get overridden:

- stdout is for results, errors go to stderr; logs go to a file, never stdout.
- Exit codes: `0` ok, `1` runtime error, `2` usage error.
- Cobra: use `RunE` and return errors; set `SilenceUsage: true`.
- Treat `huh.ErrUserAborted` as a silent cancel, never an error exit.

Optional: the full building-Go-CLIs conventions are packaged as an installable
skill for agents (configuration, output, dual-audience UX, testing, releases).
Install with `npx skills add metruzanca/template-go-cli` and update later with
`npx skills update building-clis-in-go`.

## Structure

- `main.go` — thin entrypoint; loads godotenv, inits debug logging, calls `cmd.Execute()`
- `cmd/` — Cobra command definitions; business logic lives in `internal/`, not here
- `internal/logging/` — env-gated debug logging to a file (off by default, `GO_CLI_DEBUG`)
- `internal/config/` — optional TOML config, created on first use (`GO_CLI_CONFIG_PATH` override); wired in `main.go` as commented-out code
- `.goreleaser.yaml` — release config (linux/darwin/windows, amd64/arm64) that injects version via ldflags
- `.github/workflows/release.yml` — tags a `v*` run GoReleaser and publish GitHub Releases
- `skills/building-clis-in-go/` — supplementary conventions skill, installable via `npx skills add`
