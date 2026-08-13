# Docs

Deep reference for lokan. One line per area — open what you need.

| Area                                         | What's in it                                                      |
| -------------------------------------------- | ----------------------------------------------------------------- |
| [`architecture.md`](./architecture.md)       | Stack, repo layout, storage model, build chain, design language   |
| [`api.md`](./api.md)                         | Frozen HTTP API contract — domain types, endpoints, errors        |
| [`guides.md`](./guides.md)                   | How-to for humans + agents — daily loop, roadmap modeling, AI+human collaboration, gotchas |
| [`roadmap.md`](./roadmap.md)                 | Future plans                                                      |
| [`review.md`](./review.md)                   | Design invariants + test coverage priorities from the code review |
| [`design/tokens.md`](./design/tokens.md)     | Frozen design spec — colors, type, components, contrast rules     |
| [`design/mockup.html`](./design/mockup.html) | Approved visual contract (hi-fi mockup, open in browser)          |
| [`../README.md`](../README.md)               | Quick start + command reference (user-facing)                     |

## Conventions

- **The API contract (`api.md`) and design spec (`design/tokens.md`) are
  frozen.** Any drift requires updating the doc first, then the code.
- **Docs match the current state of main.** When the code moves, the docs move
  with it — stale docs are a bug.
- Decisions worth remembering are recorded inline with `(DECIDED <date>)` in
  the doc that governs them (e.g. explicit-init, yellow fill-only, rows-not-cards).
