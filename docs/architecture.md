# Architecture — lokan

## Contents

- [What it is](#what-it-is)
- [Stack](#stack)
- [Repository Layout](#repository-layout)
- [Storage Model](#storage-model)
- [Build Chain & Embedding](#build-chain--embedding)
- [HTTP API](#http-api)
- [Design Language](#design-language)

## What it is

Lokan is a markdown task manager with a kanban-focused CLI and a web UI, all
in one static binary. The whole project lives in one git-friendly markdown
board file — a config block plus task blocks. No project directory, no config
sidecar.
(Active/Archive sections) — no database.

## Stack

| Layer    | Technology                                              | Why                                                          |
| -------- | ------------------------------------------------------- | ------------------------------------------------------------ |
| Engine   | Go 1.26, stdlib `net/http`, urfave/cli                  | Single static binary, trivial distribution                   |
| Storage  | Single markdown board file + YAML frontmatter (yaml.v3) | Git-diffable, editor-editable, human-readable                |
| Frontend | Vite + React 18 + TypeScript + Tailwind CSS             | Design tokens mapped into Tailwind `@theme`, utility styling |
| Fonts    | Geist Sans + Geist Mono (Google Fonts)                  | terminal design heritage                                     |
| Build    | `runtask` (root) → `go:embed`                           | One command produces the final binary                        |

Stack evolution is tracked in [`roadmap.md`](./roadmap.md) — a lokan board
(phases are epics, items are tasks; read with `lokan list --md docs/roadmap.md`).

## Repository Layout

```
lokan/
  runtask                 # task runner: build / test / e2e / dev web / dev engine
  dist/                   # built binary: dist/lokan (gitignored)
  docs/                   # architecture, API contract, roadmap, design
  engine/                 # Go module github.com/avrebarra/lokan
    cmd/lokan/            # urfave/cli: init create get list edit subtasks clear ui
    internal/
      types/              # enums, Task/TaskSummary/TaskFrontmatter, ALLOWED_PARENTS
      store/              # board document parse/serialize, load/find/write/create + input validation (lock + atomic rewrite)
      query/              # filters, children, sortByPriority
      server/             # HTTP handlers + seed data
    web/                  # embed package: dist/ (built frontend, committed placeholder)
  web/                    # Vite React frontend
    src/                  # components, tokens.css, index.css (tailwind entry), api client
```

## Storage Model

```
<project>/
  docs/
    board.md              whole project: config block + all task blocks (Active/Archive)
```

- **One board file, one block per task:** every task is a
  `<!-- lokan:<id>`-delimited block (YAML frontmatter — `id, title, type,
status, priority, parent?, related?, docs?, tags?, created, updated` — plus
  markdown body), grouped into `## Active` and `## Archive` (done/cancelled).
  The marker line opens one HTML comment that hides all engine markup when the
  markdown renders (e.g. GitHub); the engine reads the raw file, so parsing is
  unaffected. Older boards with bare `---` fences or self-closed markers
  still parse. A descriptive banner comment opens the file — what lokan is,
  the file format, and the reference — so cold-start readers understand it
  without lokan knowledge.
- **Configurable lanes:** the board's statuses live in the board's
  `<!-- lokan:config` block as an ordered `statuses` array (`id` +
  `archived` flag). Unconfigured projects
  use the built-in defaults (`backlog, todo, in-progress, done, cancelled`,
  last two archived). All status validation and the Active/Archive split are
  driven by this list, so custom lanes round-trip through parse and
  serialize. Renaming a lane rewrites the board file; removing one moves its
  tasks to the leftmost remaining lane.
- Mutations rewrite the document atomically (temp + rename) under a
  `<board>.lock` guard, so concurrent writers cannot lose updates.
- IDs are **plain counter values** (`1`, `2`, `3`), shared across all types.
  The counter lives in the board's config block.
- **Explicit board targeting (DECIDED 2026-08-13):** every command takes the
  board as its required first positional argument. There is no discovery and
   no default path — a markdown file is a board only when its first block is
   the `<!-- lokan:config` marker. `lokan init <file>` creates one;
  every other command errors when the file is missing or lacks the marker.

## Build Chain & Embedding

```
./runtask build
  1. cd web && npm run build        → web/dist/ (index.html + hashed assets)
  2. rm -rf engine/web/dist && cp -r web/dist engine/web/dist
  3. cd engine && go build -o dist/lokan ./cmd/lokan
     → go:embed all:dist (package engine/web) → single binary in dist/
```

- `engine/web/dist/.gitkeep` is the committed placeholder so `go build`
  works before the frontend has ever been built (go:embed fails on empty
  dirs — `.gitkeep` keeps the dir non-empty). Real build assets are
  gitignored; `runtask build` regenerates them.
- The binary lives in `dist/lokan` and is gitignored.

## HTTP API

Frozen contract in [`api.md`](./api.md):

```
GET  /              embedded app (dist/index.html)
GET  /assets/*      bundled JS/CSS
GET  /api/tasks     { tasks: [TaskSummary], statuses: [StatusDef], root }
GET  /api/task/:id  { task: Task }
POST /api/create    { title, type, priority, parent? } → { task }
POST /api/update    { id, field, value } → { task }
POST /api/move      { id, status, beforeId? } → { task }
POST /api/config/statuses  { statuses: [{ id, archived }...] } → { statuses, moved }
POST /api/clear     { scope: "archived"|"all" } → { deleted }
POST /api/seed      { created }
```

Validation is enforced server-side (enum membership, required title) — the
UI never trusts itself. Unknown routes 404; error bodies are always
`{ "error": "<message>" }`.

## Design Language

Brutalist-terminal: monochrome, sharp corners (radius 0),
no shadows except the detail modal, Geist Mono uppercase labels + Geist Sans
titles, tasks as leaderboard **rows** with hairline separators.

**Contrast rule (DECIDED 2026-08-13):** accent `#ffc800` is fill-only, never
text. Exactly two accent uses: the in-progress column's 2px top bar and the
primary CTA button. Yellow text on white fails WCAG (~1.9:1).

Full spec: [`design/tokens.md`](./design/tokens.md). Visual
contract: [`design/mockup.html`](./design/mockup.html).
