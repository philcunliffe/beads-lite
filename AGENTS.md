# Agent Instructions for beads-lite

beads-lite is a SQLite-only fork of beads. The instructions below cover the
surface that actually exists in this repo. See [README.md](README.md) for the
list of features deliberately removed (Dolt, federation, external tracker
sync, multi-repo hydration, language SDKs).

## Build & test

```sh
make build         # ./bd-lite, sqlite_lite tag, CGO_ENABLED=0
make test          # go test ./... with the sqlite_lite tag
```

`make install` copies `bd-lite` to `$HOME/.local/bin`. A separate `bd` shim
in that directory is expected to forward to `bd-lite`; it is provided by the
Gas City pack, not by this repo.

## Working with issues

- `bd init --prefix <slug>` initializes `.beads/` with the SQLite database
  and a `metadata.json` recording the prefix.
- Standard CRUD: `bd create`, `bd list`, `bd show`, `bd update`, `bd close`,
  `bd reopen`, `bd dep`, `bd label`, `bd note`, `bd query`, `bd ready`.
- JSON output is supported on every reading command via `--json`.
- Routing-aware writes are configured via `bd config`
  (`routing.default`, `routing.maintainer`, `routing.contributor`).

## Project scope

This repo intentionally has no:

- Dolt embedded engine, Dolt server, `bd dolt …` commands, or remote sync
- Jira / Linear / Notion / Azure DevOps / GitHub / GitLab integrations
- Multi-repo hydration or remote-cache routing
- Language SDKs (npm package, MCP server, separate Claude plugin)

When adding a feature, ensure it can live entirely on top of the local
SQLite store; do not reintroduce the removed surface areas. The
`internal/storage/dolt` package is a SQLite-only compatibility shim; nothing
inside it talks to Dolt.
