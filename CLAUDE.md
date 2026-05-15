# Claude Code Entry Point for beads-lite

## Read First

- **Project scope**: [README.md](README.md) — beads-lite drops Dolt, external
  trackers, multi-repo hydration, and language packages.
- **Workflow and safety**: [AGENTS.md](AGENTS.md)
- **Architecture orientation**: [docs/CLAUDE.md](docs/CLAUDE.md)

## Current Ground Rules

- Run `bd prime` before doing tracked work.
- Build with `make build` (uses the `sqlite_lite` tag, `CGO_ENABLED=0`).
- The on-disk database is a single SQLite file under `.beads/`. There is no
  Dolt push/pull, no federation, and no external tracker sync — do not add
  doc references to those features.
- The internal/storage/dolt package is a SQLite-only shim retained for source
  compatibility with the rest of the cmd/bd surface; nothing inside it talks
  to Dolt.
