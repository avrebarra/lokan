# Roadmap

Prioritized by dependency — items within a phase are roughly ordered.

## Phase 1 — Ship the binary

The rebuild is functionally complete. Make it distributable.

- [x] One-command build (`./runtask build`) → single binary with embedded UI
- [x] E2E smoke (`./runtask e2e`) covering CLI + API + embedded UI
- [x] Finalize workpool history — done on main, remaining old branches
      in the local hub are throwaway
- [x] Decide distribution: `go install`, GitHub releases, or plain binary copy
- [x] Install docs (how to get `lokan` on PATH)

## Phase 2 — AI agent integration

The task files and CLI are meant to be operated by AI agents, not just humans.

- [x] **AI agent ergonomics** — DONE (2026-08-13): agent interface documented
      in `docs/api.md` — full state via `docs/board.md`, mutations via the
      stable CLI (create/edit), output discipline (stdout/stderr, exit 0/1)
- [x] **AI-readable `list` output** — DONE (2026-08-13): `lokan list --md`
      emits compact markdown (status groups, one line per task) — markdown
      chosen over JSON after assessment: LLMs read it more token-efficiently
      and the board file is already markdown

## Phase 3 — Kanban depth

The board is read + status-advance only. Natural next steps.

- [x] Task detail editing from the UI (edit fields in modal, not just advance) — DONE (2026-08-13)
- [x] Backlog/cancelled columns or filtering UI — DONE (2026-08-13)
- [x] Drag-and-drop column moves — DONE (2026-08-13): dragging tasks
      between lanes; status click stays as the interaction contract, drag is
      the gesture on top
- [x] Subtask creation from the UI — DONE (2026-08-13)
- [x] **Config page** — customize the available lanes (add/rename/remove
      statuses) and bulk operations: clear archived, clear all tickets —
      DONE (2026-08-13): lanes live in the board's config block as an ordered
      `statuses` list (id + archived flag); config modal in the UI
      (add/rename/remove + archived toggle); renames rewrite `board.md`,
      removed lanes move tasks to the leftmost lane; `POST /api/clear` +
      `lokan clear --archived/--all` for bulk deletes
- [x] Default to light mode — DONE (2026-08-13): new sessions resolve to
      light regardless of system preference; dark is opt-in via the
      theme toggle (stored preference wins when present)

## Phase 4 — Storage & stack evolution

Assessments to run; each item is a proposal, not a commitment.

- [x] **Decouple type from ID** — DONE (2026-08-13): plain counter IDs
      (`1`, `2`) instead of type-prefixed (`epic-1`, `task-2`), so changing
      a task's type doesn't change its ID
- [x] **Single-file storage** — DONE (2026-08-13): storage is one
      `docs/board.md` (Active/Archive sections, `<!-- lokan:<id> -->`
      blocks) parsed and rewritten atomically, with this app as editor/viewer
- [x] **Switch CLI framework** — DONE (2026-08-13): urfave/cli v2 adopted —
      all commands ported (init create get list edit subtasks ui), cobra
      dropped, CLI output/exit-code contract preserved
- [x] **Tailwind CSS** — DONE (2026-08-13): Tailwind v4 adopted — tokens
      mapped into `@theme`, all components converted to utilities,
      `styles.css` removed
- [ ] **Protobuf** — assess Protobuf instead of HTTP RESTful for the API

## Phase 5 — Dual-use hardening (AI + human)

Assessment (2026-08-13): lokan already works for humans and agents, but the gap
is documentation plus a few safety/ergonomics gaps that block *clean* dual use.
Easy docs items are done; implementation items are parked for later.

### Docs (done)

- [x] **`docs/guides.md`** — DONE (2026-08-13): human daily loop, roadmap
      modeling (phases=epic, items=task), agent conventions, AI+human
      collaboration model, and common gotchas (auto-archive, lock, counter,
      type-keeps-id, silent block-skip)
- [x] **Agent write contract in `api.md`** — DONE (2026-08-13): agents read via
      `board.md`/`list --md`, mutate only via CLI/API, never hand-rewrite
      `board.md`; `id`/`created`/`updated` engine-owned; exit 0/1 discipline
- [x] **Auto-archive + gotchas documented** — DONE (2026-08-13): see `guides.md`
- [x] **`guides.md` registered in `docs/README.md`** — DONE (2026-08-13)

### Implementation (later)

- [x] **G1 — Configurable board path** — DONE (2026-08-13): every command
      targets a board explicitly via `--board <file>` — no discovery, no
      default path. A board is self-contained: a `<!-- lokan:config -->`
      block (counter, version, statuses) sits at the top, and any markdown
      file with that marker can be one. `lokan init --board <file>` creates a
      fresh board — a single `roadmap.md` can be managed by the tool itself.
- [x] **G7 — Parser warnings** — DONE (2026-08-13): already implemented —
      `parseBoard` logs `Warning: skipping invalid task block` to stderr
      (`engine/internal/store/format.go`); docs gotcha corrected to match.
- [ ] **G10 — Surface `related`/`docs`/`tags`** — fields already modeled and
      shown in `lokan get`, but `list` (table + `--md`) hides them and no
      `--tag` filter is exposed (query layer already supports tags). Add a
      `--tag` filter and show tags in `--md` output so roadmap cross-links are
      visible at a glance.
- [ ] **G11 — Dogfood the roadmap** — `docs/roadmap.md` is a hand-written
      checklist, not a lokan board. G1 landed (board is configurable), so
      author it in board format and let lokan manage its own roadmap.
