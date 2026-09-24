# NOTICE — Third-party code

Lokan is MIT-licensed (see `LICENSE`). The following vendored or derived
code is licensed under AGPLv3 (Kan — https://github.com/kan/kan).

## Vendored files (planned / upcoming)

| File                             | Source                                                                                   | License | Notes                                                      |
| -------------------------------- | ---------------------------------------------------------------------------------------- | ------- | ---------------------------------------------------------- |
| `web/src/lib/dnd/collision.ts`   | `kan/apps/web/src/views/board/dnd/collision.ts` (72 LOC `createBoardCollisionDetection`) | AGPLv3  | Will be vendored in task 03 — keep original AGPLv3 header. |
| `web/src/lib/dnd/ids.ts`         | `kan/apps/web/src/views/board/dnd/ids.ts`                                                | AGPLv3  | Helper for list-body IDs (`getListBodyId`).                |
| `web/src/lib/dnd/types.ts`       | `kan/apps/web/src/views/board/dnd/types.ts`                                              | AGPLv3  | Drag data types (`LIST`, `LIST_BODY`, `CARD`).             |
| `web/src/lib/useDragToScroll.ts` | `kan/apps/web/src/hooks/useDragToScroll.ts` (92 LOC)                                     | AGPLv3  | Horizontal drag-to-scroll — vendored in task 04.           |

## Attribution

Kan is Copyright (c) Kan contributors, licensed under the GNU Affero
General Public License v3.0. See `https://github.com/kan/kan/blob/main/LICENSE`
for the full text. Lokan preserves the original file headers when vendoring
and lists the files here.

## Dependency audit (task 02 — 2026-09-24)

Minimal vendor set for Kan board mechanics:

| Package              | Version          | Reason                                                               | Status                                                                             |
| -------------------- | ---------------- | -------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `@dnd-kit/core`      | `6.3.1`          | Drag core — `DndContext`, sensors, collision                         | **installed** (pinned)                                                             |
| `@dnd-kit/sortable`  | `10.0.0`         | Sortable contexts — columns + cards                                  | **installed** (pinned)                                                             |
| `@dnd-kit/utilities` | `3.2.2`          | `CSS.Transform` / `CSS.Translate` helpers                            | **installed** (pinned, transitive of core — pinned explicitly for reproducibility) |
| `tailwind-merge`     | `2.5.2`          | `twMerge` for class merging in `Badge` / `CircularProgress` / `List` | **installed** (pinned)                                                             |
| `lucide-react`       | `^1.31.0`        | Icon set — keep, map `Hi*` (react-icons/hi2) to lucide equivalents   | **kept** (already in deps)                                                         |
| `react-icons`        | `^5.5.0` (Kan)   | 88 MB unpacked — heavy, duplicates lucide                            | **rejected**                                                                       |
| `framer-motion`      | `^12.26.2` (Kan) | Only used in Kan `BoardsList` + onboarding — not needed for board    | **rejected**                                                                       |
| `date-fns`           | `^4.1.0` (Kan)   | Date formatting for due dates — defer to task 06                     | **deferred**                                                                       |
| `tippy.js`           | (Kan)            | Tooltips                                                             | **rejected** — use native title / minimal tooltip                                  |

Bundle impact (2026-09-24, `vite build` baseline before/after install, no imports yet):

- `vite build` `dist/assets/index-*.js` **170.71 kB** (gzip 53.98 kB) — **unchanged** after install (deps tree-shaken until imported).
- Estimated unpacked sizes: `@dnd-kit/core` 1.1 MB, `sortable` 234 kB, `tailwind-merge` 1.1 MB, `react-icons` 88 MB (avoided), `framer-motion` 4.9 MB (avoided).
- `dist/lokan` single binary ~13 MB — **unchanged** (web/dist hashed assets still embedded via `go:embed`).
- `runtask build` (`go build -o dist/lokan ./cmd/lokan`) green.

Future imports will be measured per task (03/04/05/06) with actual `import` usage.
