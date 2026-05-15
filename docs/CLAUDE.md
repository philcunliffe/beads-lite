# CLAUDE.md

This file guides Claude Code when working in beads-lite. For the full agent
checklist see [../AGENTS.md](../AGENTS.md).

## Project overview

beads-lite (command: `bd-lite`, fronted by a `bd` shim) is a single-binary
issue tracker for solo projects and Gas City supervisor integration. The
storage backend is a local SQLite file under `.beads/`. Nothing in this fork
talks to Dolt, external trackers, federation peers, or remote caches.

## Architecture

1. **Storage layer** (`internal/storage/`)
   - `sqlite/` — the only concrete storage backend (`modernc.org/sqlite`,
     pure-Go, no CGO).
   - `dolt/` — a thin compatibility shim retained so the rest of `cmd/bd`
     keeps compiling. All factory functions delegate to `sqlite/`.
   - `storage.go` — interface definitions consumed by `cmd/bd`.

2. **CLI layer** (`cmd/bd/`)
   - Cobra-based; one file per command. All commands support `--json`.
   - The `sqlite_lite` build tag selects the lite factories
     (`store_factory_sqlite_lite.go`, `init_sqlite_lite.go`, etc.).
   - Doctor checks specific to the removed features (Dolt server, federation,
     btrfs NoCOW, etc.) are no-op stubs that report "Not applicable in lite
     mode".

3. **Workflow helpers** (`internal/molecules/`, `internal/routing/`,
   `internal/formula/`, `internal/templates/`)
   - Molecule formulas (`bd mol`) drive Gas City convoys.
   - Routing is config-driven; multi-repo hydration and remote-URL targets
     were removed with the rest of the federation surface.

## Working in this repo

- `make build` produces `./bd-lite` with the `sqlite_lite` tag and
  `CGO_ENABLED=0`. Use this every time before claiming a fix.
- `make test` runs `go test -tags=sqlite_lite ./...`.
- The on-disk database is just a SQLite file; the schema lives in
  `internal/storage/sqlite/schema.go`.
- Many doctor checks are stubs. If you need to add a real check, put it in
  the doctor package and ensure it works against `*sqlite.Store`.
