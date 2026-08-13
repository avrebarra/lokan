# Lokan API Contract (frozen)

> Source of truth for engine/server implementation and frontend types. Do not
> drift — any change requires updating this doc first.

## Domain types

### TaskType

```
epic | task | subtask | bug
```

### Status

```
todo | in-progress | backlog | done | cancelled
```

### Priority

```
critical | high | medium | low
```

### ALLOWED_PARENTS (parent validation)

| type    | allowed parents |
| ------- | --------------- |
| epic    | (none)          |
| task    | epic            |
| subtask | task, bug       |
| bug     | epic, task      |

## Task frontmatter (YAML blocks in `.lokan/board.md`)

| field    | type       | required | notes                       |
| -------- | ---------- | -------- | --------------------------- |
| id       | string     | yes      | unique, e.g. `1`            |
| title    | string     | yes      |                             |
| type     | TaskType   | yes      |                             |
| status   | Status     | yes      |                             |
| priority | Priority   | yes      |                             |
| parent   | string     | no       | id of parent task           |
| related  | string[]   | no       |                             |
| docs     | string[]   | no       |                             |
| tags     | string[]   | no       |                             |
| created  | YYYY-MM-DD | yes      | set once on create          |
| updated  | YYYY-MM-DD | yes      | auto-updated on every write |

`Task` = frontmatter + `body` (raw markdown) + `filePath` (virtual task path).
`TaskSummary` = frontmatter + `filePath` (no body).

> `filePath` is a **virtual path** — all tasks live in one file
> (`.lokan/board.md`) as `<!-- lokan:<id> -->`-delimited blocks. The reported
> path (`<root>/.lokan/tasks/<id>.md`) is stable per task and is how the
> engine addresses a block; it is not a real file on disk.

## Endpoints

Base: served by the embedded web server (default port 7777).

### GET /

Static HTML app (embedded dist).

### GET /api/tasks

```json
{ "tasks": [TaskSummary...], "root": "/abs/path/to/project" }
```

### GET /api/task/:id

```json
{ "task": Task }
```

404: `{ "error": "<message>" }`

### POST /api/create

Request body: `{ "title": string, "type": TaskType, "priority": Priority, "parent"?: string }`

Validation:

- `title` missing/empty → 400 `{ "error": "Missing title" }`
- `type` not in TASK_TYPES → 400 `{ "error": "Invalid type: <v>" }`
- `priority` not in PRIORITIES → 400 `{ "error": "Invalid priority: <v>" }`

Success: `{ "task": Task }` (status set to `todo`, id from counter)
Failure: 500 `{ "error": "<message>" }`

### POST /api/update

Request body: `{ "id": string, "field": "status"|"priority"|"title"|"type"|"parent"|"tags"|"body", "value": string }`

Validation:

- `status` → must be in STATUSES, else 400 `{ "error": "Invalid status: <v>" }`
- `priority` → must be in PRIORITIES, else 400 `{ "error": "Invalid priority: <v>" }`
- `title` → free-form string
- `type` → must be in TASK_TYPES, else 400 `{ "error": "Invalid type: <v>" }`
- `parent` → free-form string (empty clears)
- `tags` → comma-separated string; split on `,`, trimmed, empties dropped
- `body` → free-form string (normalized to end with a trailing newline on write)
- other field → 400 `{ "error": "Unknown field: <f>" }`

Success: `{ "task": Task }`
Failure: 500 `{ "error": "<message>" }`

### POST /api/seed

```json
{ "created": <number of tasks created> }
```

### Errors

- 400 → validation
- 404 → task not found / unknown route
- 500 → internal failure

All error bodies: `{ "error": "<message>" }`

## Agent interface

The board and CLI are designed to be operated by AI agents. Two stable
entry points, both frozen:

- **Full state:** read `.lokan/board.md` directly — it is the complete board
  (all fields, bodies, Active/Archive sections) as markdown.
- **Lean board view:** `lokan list --md` prints a compact markdown summary —
  a `# Board — <n> active, <n> archived` header, then one `## <status>`
  group per status with one `- <id> [<priority>] <title>` line per task.
  The `--type/--status/--priority` filters apply as normal.

Mutations use the regular CLI (stable commands): `lokan create "<title>"`
(`--type/--priority/--parent/--tag`) and `lokan edit <id>`
(`--status/--priority/--title/--parent`).

Output discipline for agents:

- stdout carries the result; stderr carries errors only
- exit code 0 on success, 1 on any failure (missing project, not found, validation)
- human `list`/`get` output remains available and is agent-readable too

## Embedding

The built Vite app (from `web/`) is copied to `engine/web/dist/` by the build
and embedded via `//go:embed all:dist` (package `engine/web`). `GET /` serves
`dist/index.html`; `GET /assets/*` serves the bundled JS/CSS. A placeholder
`dist/index.html` is committed so `go build` works before the frontend build.

## File layout on disk

```
<project>/
  .lokan/
    config.json       { "counter": number, "version": string }
    board.md          single file: all task blocks
```

`board.md` holds every task as a `<!-- lokan:<id> -->`-delimited block (YAML
frontmatter fields above + blank line + markdown body), grouped into two
sections:

```
# Lokan Board

## Active

<!-- lokan:1 -->
---
id: "1"
...
---

## Archive

<!-- lokan:2 -->
...
```

`## Archive` holds tasks whose status is `done` or `cancelled`; everything
else renders under `## Active`. The engine rewrites the file atomically
(temp + rename) under a `board.md.lock` guard on every mutation.
