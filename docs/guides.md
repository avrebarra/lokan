# Guides — using lokan

Practical how-to for humans and AI agents. The API contract (`api.md`) is the
frozen source of truth for types/endpoints; this doc is the _workflow_ layer.

## Reading a board file (cold-start)

Found a `board.md` with no lokan knowledge configured? The file is
self-describing — everything you need is inside it:

- It opens with a banner comment explaining what the file is (a lokan
  kanban/roadmap), its format, and where to read the reference. The board
  file's raw markdown IS the whole state — read it directly.
- After the banner, a `<!-- lokan:config ... -->` block holds the engine
  config (board title, counter, format version, lane definitions) — comment
  hidden on render. Each task opens with a `### <id> — <title>` heading,
  then a ```lokan code block: YAML frontmatter (`id`, `title`,
  `status`, `related?`, `docs?`, `tags?`, `created`,
  `updated`), then the markdown body (`# Title`, `## Notes`,
  `## Work Log`) in its own ````markdown fence. Both fences are visible in
  rendered output, so the rendered view shows the task exactly as the raw
  file does.
- Blocks group under `## Active` (open statuses) and `## Archive` (lanes with
  `archived: true`). The board title lives inside the config block, not as a
  visible heading. Legacy boards (comment-wrapped task blocks, bare `---`
  fences) still parse; their `# Heading` migrates into the config title on
  the first rewrite.
- **To read:** the board file (full state) or `lokan list --md <board>` (one
  line per task). **To change:** the `lokan` CLI only — never hand-rewrite
  the file; every mutation rewrites it atomically under a `<board>.lock`.
  `id`/`created`/`updated` are engine-owned.

## Human workflow (daily loop)

1. `lokan init <board>` — create a board (one self-contained markdown
   file; explicit, required once). `docs/board.md` is the conventional spot.
2. `lokan create "Do thing"` — add work (optionally `--notes "<markdown>"` to seed the body).
3. `lokan list --md` — compact board view; `lokan list` for the UI-style table.
4. `lokan edit <id> --status in-progress` — advance; or `lokan ui` and click.
5. `lokan edit <id> --status done` — finishes the task; it auto-moves to Archive.
6. `lokan get <id>` — full task (frontmatter + body) when you need detail.

Every command takes the board as its first positional argument. Filters:
`--status` and `--tag` on `list`. Every card is a plain task — there are no
types, priorities, or parent nesting.

## How to model a roadmap

- **Every card is a task.** No epic/task/subtask/bug types, no priority, no
  hierarchy — the board is flat.
- **Phase / theme = a lane (status).** Rename or add lanes to model stages
  (`backlog`, `todo`, `in-progress`, `done`, `cancelled` by default).
- Cross-link code/docs with the `docs:` / `related:` / `tags:` fields.

This maps 1:1 to the board file, so the same board is readable by a human, the
UI, and an agent.

## AI agent workflow

Stable entry points (frozen in `api.md`):

- **Read full state:** open the board file directly (your board path) —
  it is the complete
  board as markdown.
- **Read lean state:** `lokan list --md` — one line per task, grouped by status.
- **Mutate:** `lokan create` / `lokan edit <id>` only. Seed a task's notes at
  create time with `lokan create <board> "<title>" --notes "<markdown>"` — the
  notes land at the top of the task body, before `## Notes`.

Conventions an agent should follow:

- Treat `id`, `created`, `updated` as engine-owned. Never invent or rewrite them.
- Pick the next `todo` item; set `--status in-progress`; set `--status done`
  when finished.
- Use `tags`/`docs`/`related` to link work to code or docs.

## AI + human collaboration model

Both a human and an agent can operate the same board. The boundary that keeps
it safe:

- **Humans:** UI, CLI, or a text editor on the board file.
- **Agents:** CLI/API only — **never hand-rewrite the board file**.

Why this holds: every engine mutation rewrites the board file atomically
(temp + rename) under a `<board>.lock` guard. A human editing the raw file in a
separate editor is outside that guard, so an agent must not also write by hand
concurrently — use the CLI and let the engine serialize. Conflicts surface as
exit code 1 with a stderr message; re-read and retry.

## Common gotchas

- **Auto-archive.** Tasks with status `done` or `cancelled` are moved into the
  `## Archive` section automatically on write. Don't be surprised when a task
  jumps sections after you flip its status.
- **`<board>.lock`.** Held only during an engine mutation. Don't delete it; if
  a crash leaves it stale, remove it manually and re-run the command.
- **Counter only increments.** IDs are plain counters (`1`, `2`, …) and never
  reused, even after a task is archived or deleted.
- **Bad blocks are skipped with a warning.** If a task block (```lokan fence
  or a legacy `<!-- lokan:<id>` comment) fails to parse (broken YAML, missing
  marker), the engine skips it and prints
  `Warning: skipping invalid task block` to stderr. The task still won't appear
  on the board — keep the board file well-formed; prefer the CLI/UI over
  hand-editing.
- **Concurrent edits: git is the guard.** The engine lock only serializes lokan
  processes. A human raw-editing the board file while an agent runs the CLI is
  out-of-band; keep the repo committed and reconcile conflicts with git.
- **Explicit init.** Every command except `init`/`help` errors with
  `not a lokan board — run lokan init <file>` when the path points at
  a missing file or one without the config marker.
