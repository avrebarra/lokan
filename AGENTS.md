# AGENTS — Lokan

Single-file kanban board. Markdown storage, Go engine, React UI. Built for
dual use — humans and AI agents operate the same board.

## Architecture in Brief

**The board file is the whole state.** A board is a markdown file whose first
block is a `<!-- lokan:config` block (title, counter, version, lanes),
followed by task blocks grouped into Active / Archive sections. Each task
opens with a `### <id> — <title>` heading, then a ```lokan code block
(YAML frontmatter) and a markdown body — the frontmatter is visible on
render, the banner and config stay comment-hidden. A banner comment at the
top explains the format to cold-start readers. Every command takes the
board as its first positional argument — there is no discovery or default
path. Read the board file directly for full state.

**The Go engine is the source of truth.** Every mutation (create/edit/move/
clear/lane rename) rewrites the board file atomically under a `<board>.lock`.
Agents never hand-rewrite it.

**The UI is a viewer + editor over the engine.** `lokan ui <file>`
serves the React app and the HTTP API; `lokan` alone is the CLI. Both
implement the same contract.

## Doc Index

| Doc                                              | What's in it                                                                                |
| ------------------------------------------------ | ------------------------------------------------------------------------------------------- |
| [`docs/api.md`](./docs/api.md)                   | Frozen HTTP API + CLI contract — domain types, endpoints, output discipline                 |
| [`docs/guides.md`](./docs/guides.md)             | How-to: human daily loop, roadmap modeling, agent workflow, AI+human collaboration, gotchas |
| [`docs/architecture.md`](./docs/architecture.md) | Stack, repo layout, storage model, build chain                                              |
| [`docs/roadmap.md`](./docs/roadmap.md)           | Future plans — a lokan board (flat tasks on lanes)                                          |
| [`docs/release.md`](./docs/release.md)           | Release routine — tagging, GoReleaser pipeline, install.sh                                  |

**Keyfiles:** storage & parsing → `engine/internal/store/`; API → `engine/internal/server/`; UI → `web/src/`.

## Conventions

- **Agents mutate via CLI/API only** — `lokan create <file> <title>` /
  `lokan edit <file> <id>`. Never write the board file by hand; read
  it freely.
- **Read lean state** via `lokan list --md <file>` (one line per
  task, grouped by status). Full state = the board file.
- **Engine-owned fields:** `id`, `created`, `updated` — never invent or
  rewrite them. IDs are plain counters (`1`, `2`, …) and never reused.
- **Statuses are configurable lanes** in the board's `<!-- lokan:config`
  block; tasks in archived lanes live under `## Archive`.
- **Every card is a plain task.** No types, priorities, or parent nesting.
- **Format before finishing** — run `./runtask format` on all git-tracked
  files touched this session (`prettier` for `.ts`/`.tsx`/`.md`, `gofmt`
  for Go); `./runtask lint` to check. Keeps formatting drift out across
  sessions.
- **Build & test:** `./runtask test` (Go) · `./runtask web build` (TS) ·
  `./runtask build` (single binary) · `./runtask e2e` (smoke).
