---
name: building-clis-in-go
description: Conventions for building Go CLIs — structure, configuration, logging, output discipline, testing, releases and dual-audience UX (humans and agents). Use when writing, extending, or reviewing a Go CLI.
---

# Building CLIs in Go

Conventions for building Go CLI applications. This repo is the reference
implementation: follow the patterns already in the codebase before inventing
new ones. Framework and library choices are recommendations (see "Recommended
packages"), not hard rules — the other conventions here apply regardless.

## Structure

- Thin `main.go` entrypoint: load godotenv, init debug logging, call `cmd.Execute()`.
- `cmd/` holds command definitions only — no business logic.
- Business logic lives in `internal/`, private to the binary.
- Release plumbing: `.goreleaser.yaml` (linux/darwin/windows, amd64/arm64, version via ldflags) and a `v*` GitHub Actions workflow that runs GoReleaser and publishes releases.

Build and verify with:

- `go build ./...`
- `go vet ./...`
- `go test ./...`
- `goreleaser check` and `goreleaser release --snapshot --clean` for release config.

## Recommended packages

- **CLI framework** — [spf13/cobra](https://github.com/spf13/cobra). Its command tree scales well: nesting subcommands per feature and registering them in `init()` keeps a CLI with lots of commands easy to maintain. If you do use it, define commands with `RunE` and return errors, and set `SilenceUsage: true` so runtime errors don't dump help text.
- **Configuration** — [BurntSushi/toml](https://github.com/BurntSushi/toml). Config file at `os.UserConfigDir()/<cli>/config.toml` (e.g. `~/.config/<cli>/config.toml`), generated on first use if it doesn't exist. Do NOT use the "warning: missing config, run init first" pattern. If an init command needs user input, still auto-run it on first use (detect missing config), in addition to exposing it as a command. Support an env var to override the config location (e.g. `GO_CLI_CONFIG_PATH`) so dev can point it at a repo-local file. Optional pattern: find-up (walk parent dirs) for project-local config, failing open to defaults — never warn about a missing config.
- **Environment overrides** — [joho/godotenv](https://github.com/joho/godotenv) to load env vars at dev time that override config options or config location. Optionally, load a local config file instead.
- **UI/pretty output** — [charmbracelet](https://github.com/charmbracelet) packages: `lipgloss`, `bubbletea`, `glow`, `bubbles` (common components), etc. `charm.land/lipgloss/v2` (with its `table` subpackage) is a valid alternative. Only style output when stdout is a terminal.
- **Prompts/forms** — [charmbracelet/huh](https://github.com/charmbracelet/huh). Treat `huh.ErrUserAborted` as a silent cancel — never an error exit.
- **README demos** — [charmbracelet/vhs](https://github.com/charmbracelet/vhs) to record terminal demos (.tape scripts → GIF/ASCIICast) for a pretty README, and [charmbracelet/vhs-action](https://github.com/charmbracelet/vhs-action) to render them in CI. Commit the generated GIFs so the README renders without tooling; deterministic via `Set FontSize`, fixed dimensions, and `Type @0ms`/`Sleep` for predictable pacing.

## Output & errors

- stdout is the pipe target — keep it clean. Errors go to stderr, never stdout.
- Only emit ANSI styling when stdout is a terminal (`term.IsTerminal` or charmbracelet's style detection); piped output stays plain and deterministic.
- Exit codes: `0` ok, `1` runtime error, `2` usage error.
- Commands return errors to a single handler instead of printing and exiting mid-flight, so runtime failures never dump usage/help text.
- See "UX: humans and agents" for how these rules serve interactive and automated users.

## UX: humans and agents

Design every command for two consumers at once: interactive humans and
automated agents/scripts. One command serves both — the same flags, the same
exit codes, no separate feature-parity binaries. Humans get the default
interactive experience; agents opt out explicitly with `--agent`.

### Both

- Every form field must also be a settable `--flag`, and forms are prefilled
  from flags already passed. The flag is the contract; the form is a
  convenience for humans.
- `--agent` is a root persistent flag that disables ALL interactivity for an
  invocation — even when stdin is a terminal. Never infer agent mode from a
  pipe or TTY; `--agent` is explicit, documented, and stable.
- Show form fields conditionally (`huh` field `Skip()`) instead of one giant
  form — only ask what still applies given earlier answers.
- Keep flags stable and additive. Agents and scripts encode flag names —
  rename via deprecation, never repurpose. Document exit codes and output so
  both audiences can rely on them.

### Humans

- Mutating commands confirm by default (`... [y/N]`); declining prints
  "Aborted." and exits cleanly. Style output when stdout is a terminal.
- Ship the full interactive experience: `huh` forms prefilled from flags,
  conditional fields, readable summaries, and examples in `--help`.
- Hide rarely-used utility commands (e.g. shell `completion`) from the help
  screen rather than removing them.

### Agents

- `--agent` guarantees stdin is never read: no forms, no y/N prompts, nothing
  blocks waiting for input.
- A missing required flag under `--agent` errors to stderr naming the field
  (`--name is required`), never launching a form.
- A mutating command under `--agent` without `--yes`/`--dry` errors with
  instructions: "pass --yes to apply or --dry to preview". Agent workflow:
  preview with `--dry`, get sign-off, apply with `--yes`.
- `--dry` prints what would happen and exits without changing anything; `--dry`
  takes priority over `--yes`.
- Deterministic output: sorted listing order, no ANSI when piped, and errors
  that lead with the fix (`no server configured; run 'disc server new' or use
  --server`).

## Logging

- Write logs to a file, never stdout — logs must never disrupt the CLI/TUI output.
- Gated by an env var and off by default; the same var can override the log path. Announce the path on stderr only when logging initializes.
- Location: in the repo during dev, and `os.UserConfigDir()/<cli>/` (e.g. `~/.config/<cli>/`) in production.

## Testing

- Table-driven tests next to source, in the same package.
- Golden-file helper (`golden(t, name, got)`) with a `-update` flag for byte-exact render output; fixture files for parser tests.
- Hermetic: `t.Setenv` for env-dependent code, `t.TempDir` for temp paths.

## Releases & versioning

- Release versions are tagged `v*` (e.g. `v0.1.0`)
- Version is baked in at build time via ldflags (`cmd.Version`) and surfaced through a built-in `--version` flag, not hardcoded
- Commit messages use conventional prefixes (`feat:`, `fix:`, `chore:`, `docs:`, `test:`, `ci:`) so the GoReleaser changelog groups cleanly
