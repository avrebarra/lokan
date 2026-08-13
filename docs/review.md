# Code Review Notes

**Date:** 2026-08-13
**Source:** Prior code reviews of the task manager (architecture + test coverage).
**Scope:** Design invariants the rebuild must keep, and test gaps worth guarding.

---

## Design Invariants (must not regress)

These decisions were identified during review and are baked into the current
implementation. Treat them as contract:

- **Atomic counter allocation.** `id.NextCounter` uses an exclusive lock file
  around the read-modify-write of `config.json`. Concurrent `create`
  invocations (CI, AI agents running in parallel) must never hand out
  duplicate counters — duplicate IDs mean silent file overwrite with no
  recovery path.
- **Server-side validation.** Every write path (CLI + HTTP) validates enum
  values against the allowed sets before touching disk. An invalid value
  would otherwise be written to frontmatter and the file would silently
  disappear from queries.
- **Cycle-safe tree building.** `buildTree` tracks ids on the current path, so
  a corrupt parent graph (A→B→A) cannot infinite-loop. Human/AI-edited
  frontmatter is an expected input.
- **O(N) id lookup on the board document.** `FindByID` scans `.lokan/board.md`
  blocks (there are no per-task files). Keep the scan linear and cheap — no
  regex-over-filepath tricks; a full board scan is fine at kanban scale.
- **ID immutable on title edits.** Title edits update the task block in
  place — the id always comes from frontmatter, never from a filename, so
  `task-NaN-title.md` corruption is impossible by construction.
- **One board file, git-friendly.** Git-diffable and editor-editable is the
  point. The single `.lokan/board.md` preserves that: `<!-- lokan:<id> -->`
  blocks, deterministic Active/Archive grouping, atomic temp+rename rewrites.
  If storage is revisited, it must preserve git-friendliness.

## Test Coverage Priorities

The core layers (id/store/query) have solid happy-path coverage. The gaps
that matter most, in order:

1. **Serialize/deserialize round-trip** — every task read and write goes
   through `format.go` (board document parse + per-block serialize). A
   round-trip test (parse → serialize → parse equals original) catches silent
   corruption that unit tests of individual fields miss.
2. **Board rewrite atomicity** — every mutation rewrites the whole document.
   Cover: lock-guarded concurrent creates never lose updates, temp+rename
   leaves no partial writes, invalid blocks are skipped with a warning.
3. **Parent type validation** — creating a `subtask` under an `epic` must
   fail. This is the core guard against an invalid task graph.
4. **Archive grouping** — `done`/`cancelled` tasks render under `## Archive`,
   everything else under `## Active`; reopening a task moves it back.
5. **Cycle detection** — buildTree on circular parent references returns
   cleanly instead of recursing forever.

## Architecture Decisions That Are Fine As-Is

- **One board file** — see invariants above.
- **No caching layer** — correct at this scale; the board scan stays cheap.
- **Boring CLI framework, stdlib HTTP** — the right call. Don't swap for
  custom parsing or a heavy web framework without a concrete win.
- **Test structure** — pure unit tests for query/id, integration tests with
  temp dirs for store. Keep this split.
