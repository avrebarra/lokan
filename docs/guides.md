# Guides — using lokan

Practical how-to for humans and AI agents. The API contract (`api.md`) is the
frozen source of truth for types/endpoints; this doc is the *workflow* layer.

## Human workflow (daily loop)

1. `lokan init` — create `.lokan/` (explicit, required once per project).
2. `lokan create -t task -p high "Do thing"` — add work.
3. `lokan list --md` — compact board view; `lokan list` for the UI-style table.
4. `lokan edit <id> --status in-progress` — advance; or `lokan ui` and click.
5. `lokan edit <id> --status done` — finishes the task; it auto-moves to Archive.
6. `lokan get <id>` — full task (frontmatter + body) when you need detail.

Filters: `--type/--status/--priority` on `list`. Hierarchy: `epic` →
`task`/`bug` → `subtask` (`lokan subtasks <id>` shows children).

## How to model a roadmap

- **Phase / theme = `epic`** — the container.
- **Work item = `task`** (or `bug` for defects).
- **Sub-step = `subtask`** (parent must be a `task`/`bug`).
- Status: `todo | in-progress | backlog | done | cancelled`.
- Cross-link code/docs with the `docs:` / `related:` / `tags:` fields.

This maps 1:1 to the board file, so the same board is readable by a human, the
UI, and an agent.

## AI agent workflow

Stable entry points (frozen in `api.md`):

- **Read full state:** open `.lokan/board.md` directly — it is the complete
  board as markdown.
- **Read lean state:** `lokan list --md` — one line per task, grouped by status.
- **Mutate:** `lokan create` / `lokan edit <id>` only.

Conventions an agent should follow:

- Treat `id`, `created`, `updated` as engine-owned. Never invent or rewrite them.
- Pick the next `todo` item; set `--status in-progress`; set `--status done`
  when finished.
- Use `tags`/`docs`/`related` to link work to code or docs.

## AI + human collaboration model

Both a human and an agent can operate the same board. The boundary that keeps
it safe:

- **Humans:** UI, CLI, or a text editor on `.lokan/board.md`.
- **Agents:** CLI/API only — **never hand-rewrite `board.md`**.

Why this holds: every engine mutation rewrites `board.md` atomically (temp file
+ rename) under a `board.md.lock` guard. A human editing the raw file in a
separate editor is outside that guard, so an agent must not also write by hand
concurrently — use the CLI and let the engine serialize. Conflicts surface as
exit code 1 with a stderr message; re-read and retry.

## Common gotchas

- **Auto-archive.** Tasks with status `done` or `cancelled` are moved into the
  `## Archive` section automatically on write. Don't be surprised when a task
  jumps sections after you flip its status.
- **`board.md.lock`.** Held only during an engine mutation. Don't delete it; if
  a crash leaves it stale, remove it manually and re-run the command.
- **Counter only increments.** IDs are plain counters (`1`, `2`, …) and never
  reused, even after a task is archived or deleted.
- **Type change keeps the ID.** Changing a task's `type` does not change its id.
- **Bad blocks are skipped with a warning.** If a `<!-- lokan:<id> -->` block
  fails to parse (broken YAML, missing marker), the engine skips it and prints
  `Warning: skipping invalid task block` to stderr. The task still won't appear
  on the board — keep `board.md` well-formed; prefer the CLI/UI over
  hand-editing.
- **Concurrent edits: git is the guard.** The engine lock only serializes lokan
  processes. A human raw-editing `board.md` while an agent runs the CLI is
  out-of-band; keep the repo committed and reconcile conflicts with git.
- **Explicit init.** Every command except `init`/`help` errors with
  `not a lokan project — run lokan init` when no `.lokan/` exists.
