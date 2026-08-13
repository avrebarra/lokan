# lokan

Markdown task manager — kanban-focused. Go engine + React frontend in a single binary.

The whole project lives in one git-friendly markdown board file: a
`<!-- lokan:config -->` block (counter, version, lanes) plus task blocks. Every
command targets a board explicitly with `--board <file>` — no discovery, no
default path. No database. Fork of `@onmyway133/nod`, rebuilt in Go.

## Quick start

```sh
./runtask build          # one command: vite build → embed → go build → dist/lokan
./dist/lokan init --board docs/board.md   # create a board (explicit — required first)
./dist/lokan create --board docs/board.md -t task "Fix counter race"   # new task
./dist/lokan list --board docs/board.md                                # board summary
./dist/lokan edit --board docs/board.md 1 --status in-progress  # update fields
./dist/lokan subtasks --board docs/board.md 1                   # children of a task
./dist/lokan ui --board docs/board.md                            # open the web UI
```

All commands except `init`/`help` take a required `--board <file>` and error
when the file is missing or lacks the lokan config marker.

## Commands

| command         | description                                                           |
| --------------- | --------------------------------------------------------------------- |
| `init`          | create a board — `--board <file>` (required)                          |
| `create`        | new task — `--type` (task/bug/epic/subtask), `--priority`, `--parent` |
| `get <id>`      | full task (frontmatter + body)                                        |
| `list`          | all tasks as a table, filterable by `--status/--type/--priority`      |
| `edit <id>`     | `--status/--priority/--title/--parent` (empty string clears parent)   |
| `subtasks <id>` | direct children, indented                                             |
| `ui`            | serve the web UI (default port 7777)                                  |

IDs are plain counter values (`1`, `2`, `3`), shared across all types.

## Web UI

`lokan ui` serves the embedded React app on `localhost:7777` and opens your
browser. Brutalist-terminal aesthetic: monochrome, mono
labels, rows with hairline separators, a single yellow accent on the
in-progress column. Light is the default; dark is opt-in via the toggle.

Board: one column per configured lane (narrow screens stack vertically).
Click a row for detail (fields, notes, subtasks); drag rows between lanes.
`+ new task` creates directly from the UI; `config` edits lanes (add/rename/
remove, archived flag) and bulk-clears archived or all tasks.

## API

JSON API for the UI, contract frozen in `docs/api.md`:

```
GET  /              embedded app
GET  /api/tasks     { tasks: [TaskSummary], statuses: [StatusDef], root }
GET  /api/task/:id  { task: Task }
POST /api/create    { title, type, priority, parent? } → { task }
POST /api/update    { id, field, value } → { task }
POST /api/move      { id, status, beforeId? } → { task }
POST /api/config/statuses  { statuses: [{ id, archived }...] } → { statuses, moved }
POST /api/clear     { scope: "archived"|"all" } → { deleted }
POST /api/seed      { created }
```

## Development

```sh
./runtask web dev    # Vite dev server with HMR + mock API server
./runtask test       # engine Go tests
./runtask e2e        # full smoke: build + init + create + list + ui + API
```

Layout: `engine/` (Go: urfave/cli, core store/id/query, HTTP server, embedded
web) · `web/` (Vite + React + TS, design system in `docs/design/`). The build
copies `web/dist` into `engine/web/dist` and embeds it via `go:embed`.

## Design

`docs/design/tokens.md` is the frozen design spec (colors, type, components,
contrast rules). `docs/design/mockup.html` is the approved visual contract.
Deviations need approval — see `docs/api.md` for the API contract.
