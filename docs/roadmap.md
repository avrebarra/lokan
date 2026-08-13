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

- [ ] **AI agent ergonomics** — a documented interface (or small integration)
      so agents can interact with the task files and data and operate the
      kanban easily (list, create, advance, update via stable commands/format)
- [ ] **AI-readable `list` output** — a compact, token-efficient output mode
      (e.g. `--json` or a minimal plain format) so agents can read the board
      without consuming large tables

## Phase 3 — Kanban depth

The board is read + status-advance only. Natural next steps.

- [ ] Task detail editing from the UI (edit fields in modal, not just advance)
- [ ] Backlog/cancelled columns or filtering UI
- [ ] Drag-and-drop column moves (deferred by design — status click is the
      interaction contract; revisit if it becomes annoying)
- [ ] Subtask creation from the UI

## Phase 4 — Storage & stack evolution

Assessments to run; each item is a proposal, not a commitment.

- [ ] **Decouple type from ID** — plain counter IDs (`1`, `2`) instead of
      type-prefixed (`epic-1`, `task-2`), so changing a task's type doesn't
      change its ID
- [x] **Single-file storage** — DONE (2026-08-13): storage is one
      `.lokan/board.md` (Active/Archive sections, `<!-- lokan:<id> -->`
      blocks) parsed and rewritten atomically, with this app as editor/viewer
- [ ] **Switch CLI framework** — from cobra to urfave/cli
- [ ] **Tailwind CSS** — adopt Tailwind for the frontend as much as possible
- [ ] **Protobuf** — assess Protobuf instead of HTTP RESTful for the API
