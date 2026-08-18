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

### 38 — 2A — Shared `ui` daemon (one server, register boards)
```lokan
id: "38"
title: 2A — Shared `ui` daemon (one server, register boards)
status: backlog
created: "2026-08-17"
updated: "2026-08-17"
```

````markdown
ASSESSED (2026-08-18) — parked, revisit later. Idea: one `lokan ui` process; re-invoking `lokan ui <file>` registers a new board entry (filepath identifier) into the running server instead of spawning a new process. Open question: how to close — no browser-driven close signal exists; needs explicit `lokan ui close <file>` / `ui stop`, an idle TTL, or a UI close-button that unregisters. Plus daemon ownership (stale locks, orphans, registry location, control channel). Current auto-pick already solved crashes; this would fix process/tab sprawl. See handoff: docs/design/shared-ui-daemon-handoff.md
````

## Archive

