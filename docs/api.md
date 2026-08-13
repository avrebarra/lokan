# Lokan API Contract (frozen)

> Source of truth for engine/server implementation and frontend types. Do not
> drift — any change requires updating this doc first.

## Domain types

### TaskType

```
epic | task | subtask | bug
```

### Status

Lanes are **configurable** — see [`POST /api/config/statuses`](#post-apiconfigstatuses).
The default lane set (also what an unconfigured project reports):

```
backlog | todo | in-progress | done | cancelled
```

`done` and `cancelled` are archived by default. Status validation on every
endpoint is against the project's **configured** lanes, not this default.

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

`Task` = frontmatter + `body` (raw markdown) + `filePath` (virtual task path) +
`lineStart`/`lineEnd` (1-based line range of the block in `board.md`).
`TaskSummary` = frontmatter + `filePath` + `lineStart`/`lineEnd` (no body).

> `filePath` is a **virtual path** — all tasks live in one file
> (`.lokan/board.md`) as `<!-- lokan:<id> -->`-delimited blocks. The reported
> path (`<root>/.lokan/tasks/<id>.md`) is stable per task and is how the
> engine addresses a block; it is not a real file on disk.
>
> `lineStart`/`lineEnd` address the block's location in `.lokan/board.md`,
> e.g. `lineStart: 5, lineEnd: 20` = "`.lokan/board.md` lines 5–20". An agent
> can use the id + line range to locate a task directly in the board file.

## Endpoints

Base: served by the embedded web server (default port 7777).

### GET /

Static HTML app (embedded dist).

### GET /api/tasks

```json
{
  "tasks": [TaskSummary...],
  "statuses": [{ "id": Status, "archived": boolean }...],
  "root": "/abs/path/to/project"
}
```

`statuses` is the effective lane set in board order (defaults when the
project has no configured lanes). The UI renders one column per lane in this
order.

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

Success: `{ "task": Task }` (status set to the first non-archived lane, id
from counter)
Failure: 500 `{ "error": "<message>" }`

### POST /api/update

Request body: `{ "id": string, "field": "status"|"priority"|"title"|"type"|"parent"|"tags"|"body", "value": string }`

Validation:

- `status` → must be a configured lane, else 400 `{ "error": "Invalid status: <v>" }`
- `priority` → must be in PRIORITIES, else 400 `{ "error": "Invalid priority: <v>" }`
- `title` → free-form string
- `type` → must be in TASK_TYPES, else 400 `{ "error": "Invalid type: <v>" }`
- `parent` → free-form string (empty clears)
- `tags` → comma-separated string; split on `,`, trimmed, empties dropped
- `body` → free-form string (normalized to end with a trailing newline on write)
- other field → 400 `{ "error": "Unknown field: <f>" }`

Success: `{ "task": Task }`
Failure: 500 `{ "error": "<message>" }`

### POST /api/move

Request body: `{ "id": string, "status": Status, "beforeId"?: string }`

Moves a task to another lane and/or position within a lane. The moved task
lands directly **before** `beforeId`; omit or empty `beforeId` to append at
the end of the lane. Backed by a physical block reorder of `board.md`
(preserves lane order, no position field).

Validation:

- `status` must be a configured lane, else 400 `{ "error": "Invalid status: <v>" }`
- `id` must exist, else 404 `{ "error": "task not found: <id>" }`
- `beforeId`, when given, must exist and already be in the target lane,
  else 400 `{ "error": "beforeId must be a task in the target lane" }`

Success: `{ "task": Task }`
Failure: 500 `{ "error": "<message>" }`

### POST /api/config/statuses

Request body: `{ "statuses": [{ "id": Status, "archived": boolean }...] }`

Replaces the project's lane set. The payload is the full new ordered list;
the server diffs it against the current lanes:

- **added** lanes (id not in current set) are appended
- **renamed** lanes (same position, id in neither list) rewrite the stored
  status of every task in that lane in `board.md`
- **removed** lanes move their tasks to the leftmost remaining lane, then
  drop the lane

Validation:

- empty list → 400 `{ "error": "At least one status is required" }`
- empty or duplicate id → 400 `{ "error": "<message>" }`

Success: `{ "statuses": [...], "moved": <tasks relocated/rewritten> }`
Failure: 500 `{ "error": "<message>" }`

### POST /api/clear

Request body: `{ "scope": "archived" | "all" }`

Bulk delete. `archived` deletes every task in an archived lane; `all` deletes
every task on the board.

Validation:

- other scope → 400 `{ "error": "Invalid scope: must be \"archived\" or \"all\"" }`

Success: `{ "deleted": <number of tasks deleted> }`
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
  group per configured lane (in lane order) with one
  `- <id> [<priority>] <title>` line per task. The
  `--type/--status/--priority` filters apply as normal.

Mutations use the regular CLI (stable commands): `lokan create "<title>"`
(`--type/--priority/--parent/--tag`) and `lokan edit <id>`
(`--status/--priority/--title/--parent`).

Output discipline for agents:

- stdout carries the result; stderr carries errors only
- exit code 0 on success, 1 on any failure (missing project, not found, validation)
- human `list`/`get` output remains available and is agent-readable too

## Agent write contract

The interface above is read/write-safe only when agents follow this contract:

- **Read** the board by opening `.lokan/board.md` or running `lokan list --md`.
  Both always reflect current state.
- **Mutate only via CLI/API.** Never hand-rewrite `.lokan/board.md`. The engine
  owns the file format — markers, block ordering, `## Archive` placement, and
  the `updated` timestamp — and rewrites it atomically under `board.md.lock`.
- **Treat these fields as engine-owned:** `id`, `created`, `updated`. Do not set
  or edit them; the engine assigns/refreshes them on every write.
- **Surfaces to expect:**
  - success → exit code 0, result on stdout
  - failure (missing project, task not found, validation) → exit code 1, message
    on stderr. Re-read state and retry; do not blindly re-issue.
- **Concurrency:** the lock guards engine-vs-engine writes, not an external text
  editor. Use the CLI so the engine serializes; don't mutate concurrently with a
  human hand-editing `board.md`.

## Embedding

The built Vite app (from `web/`) is copied to `engine/web/dist/` by the build
and embedded via `//go:embed all:dist` (package `engine/web`). `GET /` serves
`dist/index.html`; `GET /assets/*` serves the bundled JS/CSS. A placeholder
`dist/index.html` is committed so `go build` works before the frontend build.

## File layout on disk

```
<project>/
  .lokan/
    config.json       { "counter": number, "version": string, "statuses": [{ "id", "archived" }...] }
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

`## Archive` holds tasks whose lane has `archived: true`; everything else
renders under `## Active` (the default set archives `done` and `cancelled`).
The engine rewrites the file atomically (temp + rename) under a
`board.md.lock` guard on every mutation.
