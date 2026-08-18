<!--
This board is a lokan kanban / roadmap — created and managed by lokan,
a single-file markdown task tool (CLI + web UI).

File format: markdown with a lokan:config block (board title, counter,
lanes) and task blocks — each task opens with a "### <id> — <title>"
heading, a lokan code fence (YAML frontmatter), and the markdown body in
its own code fence, so raw and rendered views show the same thing.

Prefer the lokan tool (CLI or UI) for edits — hand-editing is possible
but the engine rewrites this file atomically on every change.

Tool:        https://github.com/avrebarra/lokan
Reference:   https://github.com/avrebarra/lokan/blob/main/docs/guides.md
-->

<!-- lokan:config
title: Roadmap Board
counter: 48
version: "1"
statuses:
    - id: backlog
    - id: todo
    - id: in-progress
    - id: done
      archived: true
    - id: cancelled
      archived: true
-->

## Active

### 37 — Phase 7 — UI ergonomics & distribution

```lokan
id: "37"
title: Phase 7 — UI ergonomics & distribution
status: backlog
created: "2026-08-17"
updated: "2026-08-18"
```

```markdown
Proposals assessed (2026-08-18). Each item is a proposal, not a commitment.
```

### 38 — 2A — Shared `ui` daemon (one server, register boards)

```lokan
id: "38"
title: 2A — Shared `ui` daemon (one server, register boards)
status: backlog
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
ASSESSED (2026-08-18) — parked, revisit later. Idea: one `lokan ui` process; re-invoking `lokan ui <file>` registers a new board entry (filepath identifier) into the running server instead of spawning a new process. Open question: how to close — no browser-driven close signal exists; needs explicit `lokan ui close <file>` / `ui stop`, an idle TTL, or a UI close-button that unregisters. Plus daemon ownership (stale locks, orphans, registry location, control channel). Current auto-pick already solved crashes; this would fix process/tab sprawl. See handoff: docs/design/shared-ui-daemon-handoff.md
```

### 48 — Remove types, priorities, and hierarchy — every card is a plain task

```lokan
id: "48"
title: Remove types, priorities, and hierarchy — every card is a plain task
status: backlog
created: "2026-08-18"
updated: "2026-08-18"
```

```markdown
# Remove types, priorities, and hierarchy

## Notes

- Every card is a plain task: no epic/task/subtask/bug types, no priority, no parent nesting.
- Old boards keep parsing — legacy frontmatter fields (type, priority, parent) are tolerated on read and dropped on the next write.
- Supersedes task 43 (subtask visibility) — cancelled.

## Work Log
```

## Archive

### 43 — Subtask visibility & navigation in UI

```lokan
id: "43"
title: Subtask visibility & navigation in UI
status: cancelled
created: "2026-08-18"
updated: "2026-08-18"
```

```markdown
From Annotaat review (2026-08-18). Three related UI gaps around subtask UX:

1. **Visual differentiation** — tasks and subtasks look identical in card/list view; need a clear visual cue (indent, icon, border treatment, or type badge) to tell them apart at a glance.
2. **Clickable subtasks in task detail** — when viewing a parent task's detail, subtask items should be clickable links that open the subtask's own detail view (drill-down navigation).
3. **Subtask status in task detail** — the parent task detail should surface each subtask's current status (lane) so progress is visible without leaving the parent.

Design tension: subtask depth is currently one level only (parent→child), so a simple indent + type badge approach is sufficient; deeper nesting would need a tree view.
```

### 47 — Agent convention: prettier lint all git-tracked files

```lokan
id: "47"
title: 'Agent convention: prettier lint all git-tracked files'
status: done
created: "2026-08-18"
updated: "2026-08-18"
```

```markdown
From Annotaat review (2026-08-18). Agent convention: always prettier-lint all git-tracked files, including markdown. Adds a note to AGENTS.md so future agent sessions auto-format on save/commit. Covers `.ts`, `.tsx`, `.go`, `.md`, and any other tracked files. Prevents formatting drift across sessions.
```

### 46 — Move modal-classes.ts to lib/

```lokan
id: "46"
title: Move modal-classes.ts to lib/
status: done
created: "2026-08-18"
updated: "2026-08-18"
```

```markdown
From Annotaat review (2026-08-18). `modal-classes.ts` exports `buttonClass`, `confirmClass`, `fieldClass` — shared Tailwind utility strings used across modals. It sits in `components/` but isn't a component; move to `lib/modal-classes.ts` (or `lib/classes.ts`) to match its actual role as a shared constant file.
```

### 45 — Extract server handler functions to handlers.go

```lokan
id: "45"
title: Extract server handler functions to handlers.go
status: done
created: "2026-08-18"
updated: "2026-08-18"
```

```markdown
From Annotaat review (2026-08-18). `server.go` mixes handler logic with helpers (`writeJSON`, `writeError`). Extract all `handle*` functions and response helpers into a dedicated `handlers.go` file so `server.go` stays focused on setup/routing/lifecycle. Pure code organization — no behavior change.
```

### 44 — Board filepath display in UI header

```lokan
id: "44"
title: Board filepath display in UI header
status: done
created: "2026-08-18"
updated: "2026-08-18"
```

```markdown
From Annotaat review (2026-08-18). Show the opened board's filepath in two places:

1. **Browser `<title>`** — currently just `lokan — board`; should include the filename (leaf) up to the first directory that has `.git`. E.g. `lokan — roadmap.md (lokan/)` or `lokan — docs/roadmap.md`.
2. **UI header** — below or beside the board title in the web UI, show the relative filepath so the user always knows which board file they're looking at.

The engine already receives the board path; the API just needs to pass it through (e.g. a `board_path` field in the config/status response) and the UI renders it.
```

### 40 — UI: multi-select (marquee) + bulk actions + drag multiple cards

```lokan
id: "40"
title: 'UI: multi-select (marquee) + bulk actions + drag multiple cards'
status: done
created: "2026-08-17"
updated: "2026-08-18"
```

```markdown
DONE (2026-08-18): multi-select (marquee from empty board space), checkboxes appear once a selection exists, sticky bottom BulkBar (delete / archive / move-to-lane / clear selection), and group drag (selected cards move together preserving selection order). New engine surface: POST /api/delete {ids} — store DeleteTasks, all-or-nothing, 404 on missing id — plus single-task delete from the detail modal. Click contract: selection active → row click toggles, double-click opens detail. Built in workpool/multiselect.
```

### 41 — Board format: human-readable raw header (title section or frontmatter)

```lokan
id: "41"
title: 'Board format: human-readable raw header (title section or frontmatter)'
status: done
created: "2026-08-17"
updated: "2026-08-18"
```

```markdown
ASSESSED (2026-08-18) — from Annotaat review of this board: the `<!-- lokan:<id> -->` comment-wrapped header is unreadable in raw markdown; rendered view hides the markup but the raw doc reads incoherently. Proposals: human title section per item (`## 37 — Phase 7 …` before each block) or real YAML frontmatter. Tension: task 23 deliberately comment-wraps so markup stays invisible when rendered; frontmatter would re-expose it. Title section keeps both clean. Related engine fix landed 2026-08-18: boards may be titled anything — heading preserved across rewrites.
```

### 42 — runtask install — build + put latest binary on PATH

```lokan
id: "42"
title: runtask install — build + put latest binary on PATH
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-18): `./runtask install` implemented — full build then copies dist/lokan to ~/.local/bin/lokan (chmod +x). Fixed binary now on PATH; local testing after engine changes needs no release/push. Installed while fixing the heading-preservation bug (boards may be titled anything).
```

### 33 — Phase 6 — Multi-board & dev ergonomics

```lokan
id: "33"
title: Phase 6 — Multi-board & dev ergonomics
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
Assessments from annotation pass (2026-08-15). Each item is a proposal, not a commitment.
```

### 24 — Phase 5 — Dual-use hardening (AI + human)

```lokan
id: "24"
title: Phase 5 — Dual-use hardening (AI + human)
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
Assessment (2026-08-13): lokan already works for humans and agents, but the gap is documentation plus a few safety/ergonomics gaps that block _clean_ dual use. Easy docs items are done; implementation items are parked for later.
```

### 17 — Phase 4 — Storage & stack evolution

```lokan
id: "17"
title: Phase 4 — Storage & stack evolution
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
Assessments to run; each item is a proposal, not a commitment.
```

### 10 — Phase 3 — Kanban depth

```lokan
id: "10"
title: Phase 3 — Kanban depth
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
The board is read + status-advance only. Natural next steps.
```

### 7 — Phase 2 — AI agent integration

```lokan
id: "7"
title: Phase 2 — AI agent integration
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
The task files and CLI are meant to be operated by AI agents, not just humans.
```

### 1 — Phase 1 — Ship the binary

```lokan
id: "1"
title: Phase 1 — Ship the binary
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
The rebuild is functionally complete. Make it distributable.
```

### 39 — 2B — npx install instead of Go

```lokan
id: "39"
title: 2B — npx install instead of Go
status: cancelled
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
CANCELLED (2026-08-18): a Go binary can't ride npm directly — an npm wrapper needs per-platform optionalDependencies (esbuild-style) or a release downloader, adding a second release pipeline for marginal convenience. `go install` + install.sh already covers the audience. Stay Go.
```

### 36 — Extract a shared Modal shell

```lokan
id: "36"
title: Extract a shared Modal shell
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-15): shared `Modal.tsx` (overlay + escape + header/footer slots, `escapeDisabled`, `z`, `role`, `maxWidth`, `ariaLabel`) plus `modal-classes.ts` consolidating the duplicated `buttonClass`/`fieldClass`/`confirmClass` strings; all four modals refactored onto it. The `ModalXX` rename was skipped per plan.
```

### 35 — Consolidate dev commands

```lokan
id: "35"
title: Consolidate dev commands
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-15): `./runtask preview` dropped; `dev web` (renamed from `web dev`, Vite + mock API) and `dev engine` (real binary against a tmp board + demo data, no rebuild, clean error when `dist/lokan` is missing) added; README command table and `docs/architecture.md` updated. Also fixed the preview seed-check bug that never seeded (`<!-- lokan: -->` always matched the config block).
```

### 34 — Multi-board UI without port crashes

```lokan
id: "34"
title: Multi-board UI without port crashes
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-15): `lokan ui` defaults to port 17762; when the default port is taken a free port is auto-picked (printed in the URL), so multiple boards can be viewed side-by-side. An explicit `--port` is a hard requirement and fails with a clear error if already in use. `ui` also auto-opens the browser (`--no-browser` to skip).
```

### 32 — G11 — Dogfood the roadmap

```lokan
id: "32"
title: G11 — Dogfood the roadmap
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-18): `docs/roadmap.md` is now a lokan board — phases are epics, items are tasks, DONE annotations live in task bodies, and the roadmap is managed through lokan itself (this migration)
```

### 31 — G10 — Surface `related`/`docs`/`tags`

```lokan
id: "31"
title: G10 — Surface `related`/`docs`/`tags`
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-15): `--tag` filter added to `list` (comma-separated, AND semantics via the existing query layer); tags shown in `--md` output (`(tags: a,b)`) and a TAGS column in the table when present; `docs/api.md` lean-view contract updated. (The `related`/`docs` fields remain surfaced via `lokan get` only — no list exposure yet.)
```

### 30 — G7 — Parser warnings

```lokan
id: "30"
title: G7 — Parser warnings
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13): already implemented — `parseBoard` logs `Warning: skipping invalid task block` to stderr (`engine/internal/store/format.go`); docs gotcha corrected to match.
```

### 29 — G1 — Configurable board path

```lokan
id: "29"
title: G1 — Configurable board path
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13): every command takes the board as its first positional argument — no discovery, no default path. A board is self-contained: a `<!-- lokan:config` block (counter, version, statuses) sits at the top, and any markdown file with that marker can be one. `lokan init <file>` creates a fresh board — a single `roadmap.md` can be managed by the tool itself.
```

### 28 — `guides.md` registered in `docs/README.md`

```lokan
id: "28"
title: '`guides.md` registered in `docs/README.md`'
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13)
```

### 27 — Auto-archive + gotchas documented

```lokan
id: "27"
title: Auto-archive + gotchas documented
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13): see `guides.md`
```

### 26 — Agent write contract in `api.md`

```lokan
id: "26"
title: Agent write contract in `api.md`
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13): agents read via `board.md`/`list --md`, mutate only via CLI/API, never hand-rewrite `board.md`; `id`/`created`/`updated` engine-owned; exit 0/1 discipline
```

### 25 — docs/guides.md

```lokan
id: "25"
title: docs/guides.md
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13): human daily loop, roadmap modeling (phases=epic, items=task), agent conventions, AI+human collaboration model, and common gotchas (auto-archive, lock, counter, type-keeps-id, silent block-skip)
```

### 23 — Hide engine markup when rendered

```lokan
id: "23"
title: Hide engine markup when rendered
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-15): each block's marker opens one HTML comment (`<!-- lokan:<id>` … `-->`) so lokan markup is invisible in rendered markdown (GitHub) while staying parseable from the raw file; older bare-`---` / self-closed-marker boards still parse. Boards open with a descriptive banner comment (what lokan is, the format, and the docs reference) so cold-start readers can get oriented without lokan knowledge. YAML is fenceless (no `---` delimiters) so prettier 2+/3+ leave the comment blocks untouched — verified byte-identical on both formatter versions
```

### 22 — Protobuf

```lokan
id: "22"
title: Protobuf
status: cancelled
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
CANCELLED (2026-08-18): assessed and rejected — single-user localhost REST is fine; gRPC+Protobuf adds codegen + a schema for zero benefit
```

### 21 — Tailwind CSS

```lokan
id: "21"
title: Tailwind CSS
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13): Tailwind v4 adopted — tokens mapped into `@theme`, all components converted to utilities, `styles.css` removed
```

### 20 — Switch CLI framework

```lokan
id: "20"
title: Switch CLI framework
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13): urfave/cli v2 adopted — all commands ported (init create get list edit subtasks ui), cobra dropped, CLI output/exit-code contract preserved
```

### 19 — Single-file storage

```lokan
id: "19"
title: Single-file storage
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13): storage is one `docs/board.md` (Active/Archive sections, `<!-- lokan:<id> -->` blocks) parsed and rewritten atomically, with this app as editor/viewer
```

### 18 — Decouple type from ID

```lokan
id: "18"
title: Decouple type from ID
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13): plain counter IDs (`1`, `2`) instead of type-prefixed (`epic-1`, `task-2`), so changing a task's type doesn't change its ID
```

### 16 — Default to light mode

```lokan
id: "16"
title: Default to light mode
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13): new sessions resolve to light regardless of system preference; dark is opt-in via the theme toggle (stored preference wins when present)
```

### 15 — Config page

```lokan
id: "15"
title: Config page
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13): lanes live in the board's config block as an ordered `statuses` list (id + archived flag); config modal in the UI (add/rename/remove + archived toggle); renames rewrite `board.md`, removed lanes move tasks to the leftmost lane; `POST /api/clear` + `lokan clear --archived/--all` for bulk deletes
```

### 14 — Subtask creation from the UI

```lokan
id: "14"
title: Subtask creation from the UI
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13)
```

### 13 — Drag-and-drop column moves

```lokan
id: "13"
title: Drag-and-drop column moves
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13): dragging tasks between lanes; status click stays as the interaction contract, drag is the gesture on top
```

### 12 — Backlog/cancelled columns or filtering UI

```lokan
id: "12"
title: Backlog/cancelled columns or filtering UI
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13)
```

### 11 — Task detail editing from the UI (edit fields in modal, not just advance)

```lokan
id: "11"
title: Task detail editing from the UI (edit fields in modal, not just advance)
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13)
```

### 9 — AI-readable `list` output

```lokan
id: "9"
title: AI-readable `list` output
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13): `lokan list --md` emits compact markdown (status groups, one line per task) — markdown chosen over JSON after assessment: LLMs read it more token-efficiently and the board file is already markdown
```

### 8 — AI agent ergonomics

```lokan
id: "8"
title: AI agent ergonomics
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
DONE (2026-08-13): agent interface documented in `docs/api.md` — full state via `docs/board.md`, mutations via the stable CLI (create/edit), output discipline (stdout/stderr, exit 0/1)
```

### 6 — Install docs (how to get `lokan` on PATH)

```lokan
id: "6"
title: Install docs (how to get `lokan` on PATH)
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
# Install docs (how to get `lokan` on PATH)

## Notes

## Work Log
```

### 5 — Decide distribution: `go install`, GitHub releases, or plain binary copy

```lokan
id: "5"
title: 'Decide distribution: `go install`, GitHub releases, or plain binary copy'
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
# Decide distribution: `go install`, GitHub releases, or plain binary copy

## Notes

## Work Log
```

### 4 — Finalize workpool history

```lokan
id: "4"
title: Finalize workpool history
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
done on main, remaining old branches in the local hub are throwaway
```

### 3 — E2E smoke (`./runtask e2e`) covering CLI + API + embedded UI

```lokan
id: "3"
title: E2E smoke (`./runtask e2e`) covering CLI + API + embedded UI
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
# E2E smoke (`./runtask e2e`) covering CLI + API + embedded UI

## Notes

## Work Log
```

### 2 — One-command build (`./runtask build`) → single binary with embedded UI

```lokan
id: "2"
title: One-command build (`./runtask build`) → single binary with embedded UI
status: done
created: "2026-08-17"
updated: "2026-08-17"
```

```markdown
# One-command build (`./runtask build`) → single binary with embedded UI

## Notes

## Work Log
```
