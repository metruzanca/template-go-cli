# go-cli

A starting point for new Go CLI applications. This repo is a template: clone
it, rename it, and build your own CLI on top.

## Structure

- A single-entry binary with a thin entrypoint that wires up startup concerns
  and delegates to a [Cobra](https://github.com/spf13/cobra) command tree.
- **Command definitions** build the CLI surface with Cobra; they contain no
  business logic.
- **Internal packages** hold the real work: configuration, logging, and any
  domain logic. Everything under internal is private to the binary.
- **Release plumbing.** [GoReleaser](https://goreleaser.com) ships binaries
  for Linux, macOS, and Windows on both amd64 and arm64, and a
  GitHub Actions workflow publishes releases automatically from version tags.

## Features

- **Optionally configured.** Config files are TOML ([BurntSushi/toml](https://github.com/BurntSushi/toml)) and are created
  automatically the first time they're needed, with no "run init first" ceremony. If
  your CLI doesn't need config, the wiring is commented out and the feature
  simply isn't used.
- **Off-by-default debug logging.** Debug logs go to a file, never to the
  terminal, so they can't corrupt normal output. Enabled with an env var; the
  log path is overridable. During development the log lives in the repo; in
  production it follows your OS's standard config location.
- **Dev-friendly environment overrides.** A local `.env` is loaded at
  startup with [godotenv](https://github.com/joho/godotenv), letting you
  override config settings or file locations without touching the checked-in
  config.
- **Version baked in at build time.** The binary's version is injected via
  linker flags when it's compiled from a release tag, so it's always accurate
  and never hardcoded.
- **Polished output, when you're a terminal.** Optional styling via the
  [charmbracelet](https://charm.sh) stack, with a "clean by default" rule:
  piped/redirected output stays plain and deterministic, while interactive
  terminals get colors and styling.
- **Interactive prompts** for config and forms use [huh](https://github.com/charmbracelet/huh).
- **Hermetic, deterministic tests.** Tests use environment and temp-dir
  isolation so they don't depend on your machine, and golden-file comparison so
  rendered output is checked byte-for-byte.

## Conventions

- **stdout is for results, stderr is for errors.** The two are never mixed, so
  the CLI composes cleanly with pipes and scripts.
- **Exit codes are meaningful:** success, runtime error, and usage error are
  distinct.
- **Commands fail with returned errors**, letting the framework handle exit
  behavior uniformly.
- **Logs never touch the terminal**; that's what the log file is for.
- **Missing config is not an error.** Defaults are written on first use instead
  of prompting or warning.
- **Releases are version-tagged** and follow semantic versioning; the changelog
  is grouped from conventional commit prefixes.
- **Interactive prompts treat cancellation as a normal path**: backing out is
  not a failure.

## Getting started

1. Clone the repo.
2. Rename the module, imports, and release config for your project.
3. Add your command(s) under commands and your domain logic under internal.
4. Run the test/build commands to verify.
5. Tag a version to cut a release.