# Handoff — 2A: Shared `ui` daemon (one server, register boards)

**Status:** ASSESSED (2026-08-18) → **parked**. Backlog item in `docs/roadmap.md` (Phase 7, task 38). Proposal, not a commitment. A future agent can pick this up from here.

## The idea (as the user phrased it)

> Making the process shared between all `lokan ui`? So one server active — when it's invoked again it's just adding a new entry to that server with a filepath identifier. But I don't know how to close it — how to make it close if we're not opening anything anymore.

Concretely: `lokan ui <file>` when a server is already running should **not** spawn a second process. It registers `<file>` with the running server (a board entry identified by its filepath) and opens the browser at that board's URL. One process, N boards.

## The pain it solves

- Today: N boards = N processes + N ports (auto-picked) + N browser tabs.
- Current auto-pick (already shipped) killed the _crash_ problem; this would fix the _sprawl_: one process, one management surface, a natural landing page listing all open boards.

## Why it's hard — the close problem

There is **no browser-driven close signal**. When the user closes a tab, nothing notifies the server. Options:

1. **Explicit close surface** — `lokan ui close <file>` (unregister) and/or `lokan ui stop` (shutdown daemon). Predictable, scriptable, agent-friendly. This is the `id`/`created`/`updated`-style contract lokan already likes: explicit commands, no magic.
2. **Idle TTL** — daemon exits after N minutes with no registered boards. Fuzzy; kills long-lived boards the user is just not looking at.
3. **UI close-button per board** — unregister from inside the web UI; daemon exits when the registry is empty. Needs the UI to know the daemon registry.
4. **Detect tab close via heartbeat** — fragile, needs websocket/SSE; reject.

Lean: **explicit commands (`ui close <file>`, `ui stop`) + auto-exit when registry is empty.**

## Design sketch (if picked up)

- **Registry:** a pid/lockfile + JSON registry (e.g. `~/.lokan/daemon.json` or alongside the board) mapping filepath → URL/tab.
- **Control channel:** second invocation detects the running daemon (pidfile or a fixed control port), POSTs the board path to it, daemon opens the browser to the right URL. No control channel reachable → spawn the daemon.
- **Routing:** single HTTP server; each board served at `/b/<encoded-path>` or via a board param, keeping the embedded UI contract (`/api/tasks` etc. per board context).
- **Concurrency:** per-board `<board>.lock` already serializes engine writes — unchanged.
- **Port:** daemon owns one port (default 17762, auto-pick fallback stays for the first instance only).

## Open questions

- Registry location + format (lockfile vs JSON file vs unix socket).
- How the UI knows the daemon exists (landing page needs the registry).
- Whether `--port` semantics stay (explicit fail) or the daemon owns the port.
- Stale-pid recovery after crash/kill -9.

## Suggested skills for the picking-up agent

- `git-workpool` — implement in an isolated workpool branch (this repo's convention for feature work).
- `verification-before-completion` — verify with live two-instance smoke before claiming done.
- `ponytail` — check YAGNI first: the minimal landing page (list recently-opened boards, no daemon) may deliver 80% at 20% cost.
- `clarifier` — the close-lifecycle decision (explicit commands vs TTL) is a user-facing choice; ask before building.

## Not this

- **2B (npx install)** — CANCELLED (2026-08-18): npm wrapper = second release pipeline for marginal convenience; stay Go.
