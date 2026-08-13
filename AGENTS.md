# AGENTS — Lokan

Single-file kanban board. Markdown storage, Go engine, React UI. Built for
dual use — humans and AI agents operate the same board.

## Architecture in Brief

**The board file is the whole state.** `.lokan/board.md` holds every task as
a `<!-- lokan:<id> -->` block (YAML frontmatter + markdown body), grouped into
Active / Archive sections. Read it directly for full state.

**The Go engine is the source of truth.** Every mutation (create/edit/move/
clear/lane rename) rewrites `board.md` atomically under a lock. Agents never
hand-rewrite it.

**The UI is a viewer + editor over the engine.** `lokan ui` serves the React
app and the HTTP API; `lokan` alone is the CLI. Both implement the same
contract.

## Doc Index

| Doc | What's in it |
|---|---|
| [`docs/api.md`](./docs/api.md) | Frozen HTTP API + CLI contract — domain types, endpoints, output discipline |
| [`docs/guides.md`](./docs/guides.md) | How-to: human daily loop, roadmap modeling, agent workflow, AI+human collaboration, gotchas |
| [`docs/architecture.md`](./docs/architecture.md) | Stack, repo layout, storage model, build chain |
| [`docs/roadmap.md`](./docs/roadmap.md) | Future plans (phases, gated items) |
| [`docs/review.md`](./docs/review.md) | Design invariants + test coverage priorities |

**Keyfiles:** storage & parsing → `engine/internal/store/`; API → `engine/internal/server/`; UI → `web/src/`.

## Conventions

- **Agents mutate via CLI/API only** — `lokan create` / `lokan edit <id>`.
  Never write `board.md` by hand; read it freely.
- **Read lean state** via `lokan list --md` (one line per task, grouped by
  status). Full state = `board.md`.
- **Engine-owned fields:** `id`, `created`, `updated` — never invent or
  rewrite them. IDs are plain counters (`1`, `2`, …) and never reused.
- **Statuses are configurable lanes** from `.lokan/config.json` (`statuses`
  list); tasks in archived lanes live under `## Archive`.
- **Type change keeps the id.** Changing `type` never changes a task's id.
- **Build & test:** `./runtask test` (Go) · `./runtask web build` (TS) ·
  `./runtask build` (single binary) · `./runtask e2e` (smoke).
