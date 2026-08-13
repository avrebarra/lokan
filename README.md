# lokan

Markdown task manager — kanban-focused. Go engine + React frontend in a single binary.

All tasks live in one git-friendly `.lokan/board.md` (YAML frontmatter blocks
grouped into Active/Archive sections). No database. Fork of `@onmyway133/nod`,
rebuilt in Go.

## Quick start

```sh
./runtask build          # one command: vite build → embed → go build → dist/lokan
./dist/lokan init        # create .lokan/ project (explicit — required first)
./dist/lokan create -t task "Fix counter race"    # new task
./dist/lokan list                                  # board summary
./dist/lokan edit task-1 --status in-progress     # update fields
./dist/lokan subtasks task-1                      # children of a task
./dist/lokan ui                                   # open the web UI (embedded)
```

All commands except `init`/`help` require an initialized project and error with
`not a lokan project — run lokan init` otherwise.

## Commands

| command         | description                                                           |
| --------------- | --------------------------------------------------------------------- |
| `init`          | create `.lokan/` project (config + tasks dir)                         |
| `create`        | new task — `--type` (task/bug/epic/subtask), `--priority`, `--parent` |
| `get <id>`      | full task (frontmatter + body)                                        |
| `list`          | all tasks as a table, filterable by `--status/--type/--priority`      |
| `edit <id>`     | `--status/--priority/--title/--parent` (empty string clears parent)   |
| `subtasks <id>` | direct children, indented                                             |
| `ui`            | serve the web UI (default port 7777)                                  |

IDs are type-prefixed with a shared counter (`epic-1`, `task-2`, `bug-3`).

## Web UI

`lokan ui` serves the embedded React app on `localhost:7777` and opens your
browser. Brutalist-terminal aesthetic (shiprank-derived): monochrome, mono
labels, rows with hairline separators, a single yellow accent on the
in-progress column. Dark mode follows your system with an explicit toggle.

Board: TODO / IN-PROGRESS / DONE columns (narrow screens stack vertically).
Click a row for detail (fields, notes, subtasks). `advance →` moves the task
through the status cycle. `+ new task` creates directly from the UI.

## API

JSON API for the UI, contract frozen in `docs/api.md`:

```
GET  /              embedded app
GET  /api/tasks     { tasks: [TaskSummary], root }
GET  /api/task/:id  { task: Task }
POST /api/create    { title, type, priority, parent? } → { task }
POST /api/update    { id, field, value } → { task }
POST /api/seed      { created }
```

## Development

```sh
./runtask web dev    # Vite dev server with HMR + mock API server
./runtask test       # engine Go tests
./runtask e2e        # full smoke: build + init + create + list + ui + API
```

Layout: `engine/` (Go: cobra CLI, core store/id/query, HTTP server, embedded
web) · `web/` (Vite + React + TS, design system in `docs/design/`). The build
copies `web/dist` into `engine/web/dist` and embeds it via `go:embed`.

## Design

`docs/design/tokens.md` is the frozen design spec (colors, type, components,
contrast rules). `docs/design/mockup.html` is the approved visual contract.
Deviations need approval — see `docs/api.md` for the API contract.
