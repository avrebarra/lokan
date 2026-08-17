<!--
This board is a lokan kanban / roadmap — created and managed by lokan,
a single-file markdown task tool (CLI + web UI).

File format: markdown with a lokan:config block and task blocks marked
lokan:<id> (YAML frontmatter + markdown body). All engine markup is
comment-wrapped, so rendered markdown shows only the human-readable part.

Prefer the lokan tool (CLI or UI) for edits — hand-editing is possible
but the engine rewrites this file atomically on every change.

Tool:        https://github.com/avrebarra/lokan
Reference:   https://github.com/avrebarra/lokan/blob/main/docs/guides.md
-->

<!-- lokan:config
counter: 42
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

# Roadmap Board

## Active

<!-- lokan:37
id: "37"
title: Phase 7 — UI ergonomics & distribution
type: epic
status: backlog
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
-->

Proposals assessed (2026-08-18). Each item is a proposal, not a commitment.

<!-- lokan:38
id: "38"
title: 2A — Shared `ui` daemon (one server, register boards)
type: task
status: backlog
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "37"
-->

ASSESSED (2026-08-18) — parked, revisit later. Idea: one `lokan ui` process; re-invoking `lokan ui <file>` registers a new board entry (filepath identifier) into the running server instead of spawning a new process. Open question: how to close — no browser-driven close signal exists; needs explicit `lokan ui close <file>` / `ui stop`, an idle TTL, or a UI close-button that unregisters. Plus daemon ownership (stale locks, orphans, registry location, control channel). Current auto-pick already solved crashes; this would fix process/tab sprawl. See handoff: docs/design/shared-ui-daemon-handoff.md

<!-- lokan:40
id: "40"
title: 'UI: multi-select (marquee) + bulk actions + drag multiple cards'
type: task
status: backlog
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "37"
-->

ASSESSED (2026-08-18) — from Annotaat review. Multi-select via marquee drag (image-editor style); once a selection exists, checkboxes appear on cards; sticky floating bulk-action bar at bottom (delete / archive / etc.); selected cards can be dragged together to move multiple at once. Complements single-card drag (task 13). Open questions: marquee start zone, drag-vs-marquee gesture conflict, checkbox-only entry, bulk action scope.

<!-- lokan:41
id: "41"
title: 'Board format: human-readable raw header (title section or frontmatter)'
type: task
status: backlog
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "37"
-->

ASSESSED (2026-08-18) — from Annotaat review of this board: the `<!-- lokan:<id> -->` comment-wrapped header is unreadable in raw markdown; rendered view hides the markup but the raw doc reads incoherently. Proposals: human title section per item (`## 37 — Phase 7 …` before each block) or real YAML frontmatter. Tension: task 23 deliberately comment-wraps so markup stays invisible when rendered; frontmatter would re-expose it. Title section keeps both clean. Related engine fix landed 2026-08-18: boards may be titled anything — heading preserved across rewrites.

## Archive

<!-- lokan:42
id: "42"
title: runtask install — build + put latest binary on PATH
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "37"
-->
DONE (2026-08-18): `./runtask install` implemented — full build then copies dist/lokan to ~/.local/bin/lokan (chmod +x). Fixed binary now on PATH; local testing after engine changes needs no release/push. Installed while fixing the heading-preservation bug (boards may be titled anything).

<!-- lokan:33
id: "33"
title: Phase 6 — Multi-board & dev ergonomics
type: epic
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
-->

Assessments from annotation pass (2026-08-15). Each item is a proposal, not a commitment.

<!-- lokan:24
id: "24"
title: Phase 5 — Dual-use hardening (AI + human)
type: epic
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
-->

Assessment (2026-08-13): lokan already works for humans and agents, but the gap is documentation plus a few safety/ergonomics gaps that block _clean_ dual use. Easy docs items are done; implementation items are parked for later.

<!-- lokan:17
id: "17"
title: Phase 4 — Storage & stack evolution
type: epic
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
-->

Assessments to run; each item is a proposal, not a commitment.

<!-- lokan:10
id: "10"
title: Phase 3 — Kanban depth
type: epic
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
-->

The board is read + status-advance only. Natural next steps.

<!-- lokan:7
id: "7"
title: Phase 2 — AI agent integration
type: epic
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
-->

The task files and CLI are meant to be operated by AI agents, not just humans.

<!-- lokan:1
id: "1"
title: Phase 1 — Ship the binary
type: epic
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
-->

The rebuild is functionally complete. Make it distributable.

<!-- lokan:39
id: "39"
title: 2B — npx install instead of Go
type: task
status: cancelled
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "37"
-->

CANCELLED (2026-08-18): a Go binary can't ride npm directly — an npm wrapper needs per-platform optionalDependencies (esbuild-style) or a release downloader, adding a second release pipeline for marginal convenience. `go install` + install.sh already covers the audience. Stay Go.

<!-- lokan:36
id: "36"
title: Extract a shared Modal shell
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "33"
-->

DONE (2026-08-15): shared `Modal.tsx` (overlay + escape + header/footer slots, `escapeDisabled`, `z`, `role`, `maxWidth`, `ariaLabel`) plus `modal-classes.ts` consolidating the duplicated `buttonClass`/`fieldClass`/`confirmClass` strings; all four modals refactored onto it. The `ModalXX` rename was skipped per plan.

<!-- lokan:35
id: "35"
title: Consolidate dev commands
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "33"
-->

DONE (2026-08-15): `./runtask preview` dropped; `dev web` (renamed from `web dev`, Vite + mock API) and `dev engine` (real binary against a tmp board + demo data, no rebuild, clean error when `dist/lokan` is missing) added; README command table and `docs/architecture.md` updated. Also fixed the preview seed-check bug that never seeded (`<!-- lokan: -->` always matched the config block).

<!-- lokan:34
id: "34"
title: Multi-board UI without port crashes
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "33"
-->

DONE (2026-08-15): `lokan ui` defaults to port 17762; when the default port is taken a free port is auto-picked (printed in the URL), so multiple boards can be viewed side-by-side. An explicit `--port` is a hard requirement and fails with a clear error if already in use. `ui` also auto-opens the browser (`--no-browser` to skip).

<!-- lokan:32
id: "32"
title: G11 — Dogfood the roadmap
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "24"
-->

DONE (2026-08-18): `docs/roadmap.md` is now a lokan board — phases are epics, items are tasks, DONE annotations live in task bodies, and the roadmap is managed through lokan itself (this migration)

<!-- lokan:31
id: "31"
title: G10 — Surface `related`/`docs`/`tags`
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "24"
-->

DONE (2026-08-15): `--tag` filter added to `list` (comma-separated, AND semantics via the existing query layer); tags shown in `--md` output (`(tags: a,b)`) and a TAGS column in the table when present; `docs/api.md` lean-view contract updated. (The `related`/`docs` fields remain surfaced via `lokan get` only — no list exposure yet.)

<!-- lokan:30
id: "30"
title: G7 — Parser warnings
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "24"
-->

DONE (2026-08-13): already implemented — `parseBoard` logs `Warning: skipping invalid task block` to stderr (`engine/internal/store/format.go`); docs gotcha corrected to match.

<!-- lokan:29
id: "29"
title: G1 — Configurable board path
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "24"
-->

DONE (2026-08-13): every command takes the board as its first positional argument — no discovery, no default path. A board is self-contained: a `<!-- lokan:config` block (counter, version, statuses) sits at the top, and any markdown file with that marker can be one. `lokan init <file>` creates a fresh board — a single `roadmap.md` can be managed by the tool itself.

<!-- lokan:28
id: "28"
title: '`guides.md` registered in `docs/README.md`'
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "24"
-->

DONE (2026-08-13)

<!-- lokan:27
id: "27"
title: Auto-archive + gotchas documented
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "24"
-->

DONE (2026-08-13): see `guides.md`

<!-- lokan:26
id: "26"
title: Agent write contract in `api.md`
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "24"
-->

DONE (2026-08-13): agents read via `board.md`/`list --md`, mutate only via CLI/API, never hand-rewrite `board.md`; `id`/`created`/`updated` engine-owned; exit 0/1 discipline

<!-- lokan:25
id: "25"
title: docs/guides.md
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "24"
-->

DONE (2026-08-13): human daily loop, roadmap modeling (phases=epic, items=task), agent conventions, AI+human collaboration model, and common gotchas (auto-archive, lock, counter, type-keeps-id, silent block-skip)

<!-- lokan:23
id: "23"
title: Hide engine markup when rendered
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "17"
-->

DONE (2026-08-15): each block's marker opens one HTML comment (`<!-- lokan:<id>` … `-->`) so lokan markup is invisible in rendered markdown (GitHub) while staying parseable from the raw file; older bare-`---` / self-closed-marker boards still parse. Boards open with a descriptive banner comment (what lokan is, the format, and the docs reference) so cold-start readers can get oriented without lokan knowledge. YAML is fenceless (no `---` delimiters) so prettier 2+/3+ leave the comment blocks untouched — verified byte-identical on both formatter versions

<!-- lokan:22
id: "22"
title: Protobuf
type: task
status: cancelled
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "17"
-->

CANCELLED (2026-08-18): assessed and rejected — single-user localhost REST is fine; gRPC+Protobuf adds codegen + a schema for zero benefit

<!-- lokan:21
id: "21"
title: Tailwind CSS
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "17"
-->

DONE (2026-08-13): Tailwind v4 adopted — tokens mapped into `@theme`, all components converted to utilities, `styles.css` removed

<!-- lokan:20
id: "20"
title: Switch CLI framework
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "17"
-->

DONE (2026-08-13): urfave/cli v2 adopted — all commands ported (init create get list edit subtasks ui), cobra dropped, CLI output/exit-code contract preserved

<!-- lokan:19
id: "19"
title: Single-file storage
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "17"
-->

DONE (2026-08-13): storage is one `docs/board.md` (Active/Archive sections, `<!-- lokan:<id> -->` blocks) parsed and rewritten atomically, with this app as editor/viewer

<!-- lokan:18
id: "18"
title: Decouple type from ID
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "17"
-->

DONE (2026-08-13): plain counter IDs (`1`, `2`) instead of type-prefixed (`epic-1`, `task-2`), so changing a task's type doesn't change its ID

<!-- lokan:16
id: "16"
title: Default to light mode
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "10"
-->

DONE (2026-08-13): new sessions resolve to light regardless of system preference; dark is opt-in via the theme toggle (stored preference wins when present)

<!-- lokan:15
id: "15"
title: Config page
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "10"
-->

DONE (2026-08-13): lanes live in the board's config block as an ordered `statuses` list (id + archived flag); config modal in the UI (add/rename/remove + archived toggle); renames rewrite `board.md`, removed lanes move tasks to the leftmost lane; `POST /api/clear` + `lokan clear --archived/--all` for bulk deletes

<!-- lokan:14
id: "14"
title: Subtask creation from the UI
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "10"
-->

DONE (2026-08-13)

<!-- lokan:13
id: "13"
title: Drag-and-drop column moves
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "10"
-->

DONE (2026-08-13): dragging tasks between lanes; status click stays as the interaction contract, drag is the gesture on top

<!-- lokan:12
id: "12"
title: Backlog/cancelled columns or filtering UI
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "10"
-->

DONE (2026-08-13)

<!-- lokan:11
id: "11"
title: Task detail editing from the UI (edit fields in modal, not just advance)
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "10"
-->

DONE (2026-08-13)

<!-- lokan:9
id: "9"
title: AI-readable `list` output
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "7"
-->

DONE (2026-08-13): `lokan list --md` emits compact markdown (status groups, one line per task) — markdown chosen over JSON after assessment: LLMs read it more token-efficiently and the board file is already markdown

<!-- lokan:8
id: "8"
title: AI agent ergonomics
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "7"
-->

DONE (2026-08-13): agent interface documented in `docs/api.md` — full state via `docs/board.md`, mutations via the stable CLI (create/edit), output discipline (stdout/stderr, exit 0/1)

<!-- lokan:6
id: "6"
title: Install docs (how to get `lokan` on PATH)
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "1"
-->

# Install docs (how to get `lokan` on PATH)

## Notes

## Work Log

<!-- lokan:5
id: "5"
title: 'Decide distribution: `go install`, GitHub releases, or plain binary copy'
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "1"
-->

# Decide distribution: `go install`, GitHub releases, or plain binary copy

## Notes

## Work Log

<!-- lokan:4
id: "4"
title: Finalize workpool history
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "1"
-->

done on main, remaining old branches in the local hub are throwaway

<!-- lokan:3
id: "3"
title: E2E smoke (`./runtask e2e`) covering CLI + API + embedded UI
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "1"
-->

# E2E smoke (`./runtask e2e`) covering CLI + API + embedded UI

## Notes

## Work Log

<!-- lokan:2
id: "2"
title: One-command build (`./runtask build`) → single binary with embedded UI
type: task
status: done
priority: medium
created: "2026-08-17"
updated: "2026-08-17"
parent: "1"
-->

# One-command build (`./runtask build`) → single binary with embedded UI

## Notes

## Work Log
