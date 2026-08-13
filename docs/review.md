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
- **O(1) id lookup.** `FindByID` resolves the file by id-prefix glob
  (`task-2-*.md`) instead of scanning and parsing the whole tasks dir. Don't
  reintroduce full scans on the hot path.
- **Filename counter derived from the file, not user input.** Title edits
  re-parse the counter from the existing basename, so `task-NaN-title.md`
  corruption is impossible.
- **One file per task.** Git-diffable and editor-editable is the point.
  Consolidating into a database or a single blob would kill the
  human-readable benefit. If storage is revisited, it must preserve
  git-friendliness.

## Test Coverage Priorities

The core layers (id/store/query) have solid happy-path coverage. The gaps
that matter most, in order:

1. **Serialize/deserialize round-trip** — every task read and write goes
   through `format.go`. A round-trip test (parse → serialize → parse equals
   original) catches silent corruption that unit tests of individual fields
   miss.
2. **`renameTask`** — writes a new file then deletes the old. Partial failure
   (write succeeds, unlink fails) leaves duplicates. Cover: old file removed,
   new file correct, collision behavior.
3. **Parent type validation** — creating a `subtask` under an `epic` must
   fail. This is the core guard against an invalid task graph.
4. **Slug edge cases** — empty title, all-symbol title (`!!!`), exactly-50
   chars, unicode. All-symbol input can produce `task-1-.md`, which is
   invalid on some filesystems.
5. **Cycle detection** — buildTree on circular parent references returns
   cleanly instead of recursing forever.

## Architecture Decisions That Are Fine As-Is

- **One file per task** — see invariants above.
- **No caching layer** — correct at this scale; the O(1) lookup buys the
  headroom first.
- **Boring CLI framework, stdlib HTTP** — the right call. Don't swap for
  custom parsing or a heavy web framework without a concrete win.
- **Test structure** — pure unit tests for query/id, integration tests with
  temp dirs for store. Keep this split.
