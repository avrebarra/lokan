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

## Task frontmatter (YAML in `.lokan/tasks/*.md`)

| field    | type       | required | notes                       |
| -------- | ---------- | -------- | --------------------------- |
| id       | string     | yes      | unique, e.g. `task-01`      |
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

`Task` = frontmatter + `body` (raw markdown) + `filePath` (absolute path).
`TaskSummary` = frontmatter + `filePath` (no body).

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

Request body: `{ "id": string, "field": "status"|"priority"|"title", "value": string }`

Validation:

- `status` → must be in STATUSES, else 400 `{ "error": "Invalid status: <v>" }`
- `priority` → must be in PRIORITIES, else 400 `{ "error": "Invalid priority: <v>" }`
- `title` → free-form string
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
    tasks/
      task-01.md
      task-02.md
```

Task file: YAML frontmatter (fields above) + blank line + markdown body.
