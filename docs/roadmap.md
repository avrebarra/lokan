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
      in `docs/api.md` — full state via `.lokan/board.md`, mutations via the
      stable CLI (create/edit), output discipline (stdout/stderr, exit 0/1)
- [x] **AI-readable `list` output** — DONE (2026-08-13): `lokan list --md`
      emits compact markdown (status groups, one line per task) — markdown
      chosen over JSON after assessment: LLMs read it more token-efficiently
      and the board file is already markdown

## Phase 3 — Kanban depth

The board is read + status-advance only. Natural next steps.

- [x] Task detail editing from the UI (edit fields in modal, not just advance) — DONE (2026-08-13)
- [x] Backlog/cancelled columns or filtering UI — DONE (2026-08-13)
- [ ] Drag-and-drop column moves — requested (2026-08-13): dragging tasks
      between lanes; status click stays as the interaction contract, drag is
      the gesture on top
- [x] Subtask creation from the UI — DONE (2026-08-13)
- [ ] **Config page** — customize the available lanes (add/rename/remove
      statuses) and bulk operations: clear archived, clear all tickets
- [ ] Default to light mode — currently dark-first/follows system; flip the
      default to light with dark as opt-in

## Phase 4 — Storage & stack evolution

Assessments to run; each item is a proposal, not a commitment.

- [x] **Decouple type from ID** — DONE (2026-08-13): plain counter IDs
      (`1`, `2`) instead of type-prefixed (`epic-1`, `task-2`), so changing
      a task's type doesn't change its ID
- [x] **Single-file storage** — DONE (2026-08-13): storage is one
      `.lokan/board.md` (Active/Archive sections, `<!-- lokan:<id> -->`
      blocks) parsed and rewritten atomically, with this app as editor/viewer
- [x] **Switch CLI framework** — DONE (2026-08-13): urfave/cli v2 adopted —
      all commands ported (init create get list edit subtasks ui), cobra
      dropped, CLI output/exit-code contract preserved
- [x] **Tailwind CSS** — DONE (2026-08-13): Tailwind v4 adopted — tokens
      mapped into `@theme`, all components converted to utilities,
      `styles.css` removed
- [ ] **Protobuf** — assess Protobuf instead of HTTP RESTful for the API
