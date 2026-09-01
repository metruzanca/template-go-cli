Search and replace ALL instances of `metruzanca/go-cli` with `{username}/{repo}` (go.mod, imports, release config, README, AGENTS, etc.). Ask the user if they want to use their github username + the name of the folder this template was cloned into. Then once renamed, delete this paragraph

## Project

Go CLI application built with [Cobra](https://github.com/spf13/cobra). Used as a starting template for new CLI projects.

## Structure

- `main.go` — entrypoint; defines `version`, `commit`, `date` vars injected via ldflags
- `cmd/` — Cobra command definitions
- `.goreleaser.yaml` — release config (linux/darwin/windows, amd64/arm64)
- `.github/workflows/release.yml` — tags a `v*` run GoReleaser and publish GitHub Releases

## Commands

- `go build ./...` — build
- `go vet ./...` — vet
- `goreleaser check` — validate release config
- `goreleaser release --snapshot --clean` — local test release

## Recommended packages

- **Configuration** — [BurntSushi/toml](https://github.com/BurntSushi/toml). Config file: `~/.config/go-cli/config.toml`, generated on first use if it doesn't exist. Do NOT use the "warning: missing config, run init first" pattern. If an init command needs user input, still auto-run it on first use (detect missing config), in addition to exposing it as a command.
- **Environment overrides** — [joho/godotenv](https://github.com/joho/godotenv) to load env vars at dev time that override config options. Optionally, load a local config file instead.
- **UI/pretty output** — [charmbracelet](https://github.com/charmbracelet) packages: `lipgloss`, `bubbletea`, `glow`, `bubbles` (common components), etc.
- **Prompts/forms** — [charmbracelet/huh](https://github.com/charmbracelet/huh)

## Logging

- Write logs to a file so they never disrupt the CLI/TUI output.
- Location: in the repo during dev, and `~/.config/go-cli/` in production.

## Conventions

- Release versions are tagged `v*` (e.g. `v0.1.0`)
- Binary version info is baked in at build time via ldflags, not hardcoded
