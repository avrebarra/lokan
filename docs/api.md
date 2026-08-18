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

## Task frontmatter (```lokan blocks in `docs/board.md`)

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
> (`docs/board.md`) as `### <id> — <title>`-headed blocks whose frontmatter
> sits in a visible ```lokan fence (legacy `<!-- lokan:<id>` comments still
> parse). The reported path (`<board>#<id>`) is stable per task and is how
> the engine addresses a block; it is not a real file on disk.
>
> `lineStart`/`lineEnd` address the block's location in `.lokan/board.md`,
> e.g. `lineStart: 5, lineEnd: 20` = "`.lokan/board.md` lines 5–20". An agent
> can use the id + line range to locate a task directly in the board file.

## Endpoints

Base: served by the embedded web server (default port 17762; if taken, a free port is auto-picked and printed). An explicit `--port` fails if already in use.

### GET /

Static HTML app (embedded dist).

### GET /api/tasks

```json
{
  "tasks": [TaskSummary...],
  "statuses": [{ "id": Status, "archived": boolean }...],
  "root": "/abs/path/to/board.md",
  "board_path": "docs/board.md",
  "board_root": "lokan"
}
```

`statuses` is the effective lane set in board order (defaults when the
project has no configured lanes). The UI renders one column per lane in this
order.

`root` is the board path the server was started with. `board_path` is that
path made relative to the nearest git root (walking up from the board file
for a `.git` entry), so the UI can show the leaf and parent dirs down to the
repo; `board_root` is that git root's directory name. When no git root is
found, `board_path` falls back to `root` and `board_root` is empty.

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
the end of the lane. Backed by a physical block reorder of `docs/board.md`
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
  status of every task in that lane in `docs/board.md`
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

- **Full state:** read the board file directly (your board path) — it
  is the complete board (config, all fields, bodies, Active/Archive
  sections) as markdown.
- **Lean board view:** `lokan list --md <file>` prints a compact
  markdown summary —
  a `# Board — <n> active, <n> archived` header, then one `## <status>`
  group per configured lane (in lane order) with one
  `- <id> [<priority>] <title>` line per task, plus `(tags: a,b)` when
  the task carries tags. The `--type/--status/--priority/--tag` filters
  apply as normal; `--tag` accepts comma-separated or repeated values
  (AND semantics).

Mutations use the regular CLI (stable commands): `lokan create "<title>"`
(`--type/--priority/--parent/--tag`) and `lokan edit <id>`
(`--status/--priority/--title/--parent`).

Output discipline for agents:

- stdout carries the result; stderr carries errors only
- exit code 0 on success, 1 on any failure (missing board path, not found, validation)
- human `list`/`get` output remains available and is agent-readable too

## Agent write contract

The interface above is read/write-safe only when agents follow this contract:

- **Read** the board by opening your board file or running
  `lokan list --md <file>`.
  Both always reflect current state.
- **Mutate only via CLI/API.** Never hand-rewrite the board file (your
  board path). The engine
  owns the file format — markers, block ordering, `## Archive` placement, and
  the `updated` timestamp — and rewrites it atomically under `<board>.lock`.
- **Treat these fields as engine-owned:** `id`, `created`, `updated`. Do not set
  or edit them; the engine assigns/refreshes them on every write.
- **Surfaces to expect:**
  - success → exit code 0, result on stdout
  - failure (missing board, task not found, validation) → exit code 1, message
    on stderr. Re-read state and retry; do not blindly re-issue.
- **Concurrency:** the lock guards engine-vs-engine writes, not an external text
  editor. Use the CLI so the engine serializes; don't mutate concurrently with a
  human hand-editing the board file.

## Embedding

The built Vite app (from `web/`) is copied to `engine/web/dist/` by the build
and embedded via `//go:embed all:dist` (package `engine/web`). `GET /` serves
`dist/index.html`; `GET /assets/*` serves the bundled JS/CSS. A placeholder
`dist/index.html` is committed so `go build` works before the frontend build.

## File layout on disk

```
<project>/
  docs/
    board.md          a board: config block + all task blocks
```

A board is a markdown file that opens with a descriptive banner comment
(what lokan is, the file format, and a reference to the docs), followed by
the `<!-- lokan:config` block (YAML: `title`, `counter`, `version`,
`statuses`) and the task blocks. Every command takes the board as its first
positional argument — there is no discovery and no default path.
`lokan init <file>` creates a fresh board; any markdown file can become one.

The board file holds every task as a `### <id> — <title>` heading plus a
```lokan fence (YAML frontmatter fields above + markdown body), grouped
into two sections. The fence keeps the frontmatter visible in rendered
output (e.g. GitHub) while the raw file stays parseable; the banner and
config stay comment-hidden. Older boards (comment-wrapped task blocks with
bare `---` fences or a self-closed marker line) are still read:

```
<!--
This board is a lokan kanban / roadmap — created and managed by lokan...
Reference:   https://github.com/avrebarra/lokan/blob/main/docs/guides.md
-->

<!-- lokan:config
counter: 0
version: "1"
statuses:
    - id: backlog
-->

# Lokan Board

## Active

<!-- lokan:1
id: "1"
...
-->

## Archive

<!-- lokan:2
...
-->
```

`## Archive` holds tasks whose lane has `archived: true`; everything else
renders under `## Active` (the default set archives `done` and `cancelled`).
The engine rewrites the file atomically (temp + rename) under a
`<board>.lock` guard on every mutation.
