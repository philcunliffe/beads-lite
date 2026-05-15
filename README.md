# beads-lite

**SQLite-only fork of [beads](https://github.com/gastownhall/beads) for solo
projects and Gas City supervisor integration.**

beads-lite is a dependency-graph issue tracker meant to be embedded in a
single repository or wired into a Gas City supervisor as a workspace beads
backend. Compared to the upstream project it drops:

- Dolt (no embedded server, no federation, no `bd dolt …`)
- External tracker sync (Jira, Linear, Notion, Azure DevOps, GitHub, GitLab)
- Multi-repo hydration and remote-cache routing
- Standalone language SDKs (npm package, MCP/plugin distribution)
- The marketing website and the long upstream changelog

## What still works

- Issue CRUD: `bd init`, `bd create`, `bd list`, `bd show`, `bd update`,
  `bd close`, `bd reopen`, `bd dep`, `bd label`, `bd note`, `bd query`,
  `bd ready`
- Routing-aware writes via `bd config` (`routing.default`,
  `routing.maintainer`, `routing.contributor`)
- Molecules (`bd mol …`), gates, batch mode, hooks, audit, JSON output
- Messaging primitives (`bd message …`, `bd mail …`) used by Gas City rigs
- `bd prime`, `bd remember`, `bd doctor` (lite check set), `bd bootstrap`

## Build & install

```sh
make build         # produces ./bd-lite (sqlite_lite tag, CGO disabled)
make install       # copies bd-lite to $HOME/.local/bin
```

A `bd` shim at `$HOME/.local/bin/bd` is expected to forward to `bd-lite`. The
shim is *not* shipped from this repo — Gas City installs it from the pack.

## Usage

```sh
cd my-project
bd init --prefix my       # create .beads/ with SQLite database
bd create "first issue" --type task
bd list --json
```

## Layout

- `cmd/bd/` — CLI entry point (uses the `sqlite_lite` build tag)
- `internal/storage/sqlite/` — the only storage backend
- `internal/storage/dolt/` — compatibility shim that wraps the SQLite store
  (kept so existing call sites compile; no Dolt code is executed)
- `internal/molecules/`, `internal/routing/`, `internal/formula/`,
  `internal/templates/` — workflow plumbing reused from upstream

## License

Apache-2.0 (see `LICENSE`).
