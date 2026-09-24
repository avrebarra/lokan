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
title: Kan-Inspired Reskin — Lokan Board
counter: 12
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

### 3 — 03 — Board drag: HTML5 → @dnd-kit + collision

```lokan
id: "3"
title: '03 — Board drag: HTML5 → @dnd-kit + collision'
status: backlog
created: "2026-09-24"
updated: "2026-09-24"
tags:
    - kan-reskin
    - board
    - dnd
```

```markdown
# 03 — Board drag: HTML5 → @dnd-kit + collision

## Goal

Replace Lokan's native `text/x-lokan-task` drag (Board.tsx `handleDragOver 0.7*height` hack) with Kan's `DndContext + collision.ts`.

## Scope

- Vendor `views/board/dnd/collision.ts` (72 LOC `createBoardCollisionDetection`), `dnd/ids.ts`, `dnd/types.ts` → `web/src/lib/dnd/`.
- Refactor `Board.tsx`: wrap columns in `DndContext` (`PointerSensor`, `KeyboardSensor`), `SortableContext(horizontalListSortingStrategy)` for columns, `SortableContext(verticalListSortingStrategy)` inside each `Column`.
- Map `onMove(ids, status, beforeId)` to current `moveTask` / `moveMany` APIs (preserve `beforeId=''` append + `isMoveManyNoop` guards).
- Keep `selectedIds` multi-select + `marqueeLive` compat (Kan has no marquee — ensure marquee still highlights).
- Preserve `flashMoved` + `movedId` left-flash.

## Non-goals

- No visual reskin; keep `TaskRow` leaderboard for now.

## Acceptance

- [ ] Single + multi-select drag works cross-lane + reorder-in-lane, with insertion indicator (keep 2px `bg-fg` bar or Kan's overlay).
- [ ] Keyboard sortable (`sortableKeyboardCoordinates`) works.
- [ ] No regression on `marquee select` + `BulkBar`.

## Refs

- `/tmp/kan/apps/web/src/views/board/dnd/collision.ts`, `views/board/index.tsx` `DndContext`
- `web/src/components/Board.tsx`, `Column.tsx`

## Effort

M — 2d

## Work Log
```

### 4 — 04 — Horizontal scroll + drag-to-scroll

```lokan
id: "4"
title: 04 — Horizontal scroll + drag-to-scroll
status: backlog
created: "2026-09-24"
updated: "2026-09-24"
tags:
    - kan-reskin
    - board
    - ux
```

```markdown
# 04 — Horizontal scroll + drag-to-scroll

## Goal

Switch board layout from `grid repeat(n,1fr)` to Kan's horizontal scroll, add `useDragToScroll`.

## Scope

- Change `Board.tsx` container from `grid` to `flex overflow-x-auto snap-x snap-mandatory` with `min-w-[18rem] max-w-[18rem]` per `Column` (Kan `List.tsx` width).
- Vendor `hooks/useDragToScroll.ts` (92 LOC) → `web/src/lib/useDragToScroll.ts`; wire via `ref + onMouseDown` on board scroll container.
- Guard: ignore `[draggable],[data-board-draggable],button,input,a,[role=button]` like Kan — so dnd-kit drag not hijacked.
- Responsive: keep `<900px` stack fallback? Or keep horizontal even on mobile with `snap-start`. Decide.
- Add `scrollbar` styling (`scrollbar-track-rounded` etc) or keep Lokan hairlines.

## Acceptance

- [ ] Boards with 5+ lanes scroll horizontally; drag on background pans.
- [ ] Desktop + mobile (touch) not broken.
- [ ] No conflict with dnd-kit pointer sensor.

## Refs

- `/tmp/kan/apps/web/src/hooks/useDragToScroll.ts`
- `web/src/components/Board.tsx` (grid style)

## Effort

S — 1d

## Work Log
```

### 5 — 05 — Column → Kan List (rounded container + header)

```lokan
id: "5"
title: 05 — Column → Kan List (rounded container + header)
status: backlog
created: "2026-09-24"
updated: "2026-09-24"
tags:
    - kan-reskin
    - column
```

```markdown
# 05 — Column → Kan List (rounded container + header)

## Goal

Restyle `Column.tsx` from brutalist hairline header to Kan `List.tsx` card container.

## Scope

- New look: `rounded-md border border-light-400 bg-light-300 dark:bg-dark-100 py-2 pl-2 pr-1 mr-5 min-w-[18rem]` (adapt to tokens: `bg-zebra` → Kan `bg-light-300`).
- Header: editable `input` for lane name (Kan `useForm` inline edit → reuse Lokan `ModalConfig` rename flow or inline like Kan), `+` add card button (`Tooltip` disabled when no permission), `…` `Dropdown` (Add card / Delete list).
- Wire `List` `useSortable({id: listPublicId, data:{type:"LIST"}})` + `isDragging`→`cursor-grabbing` (already in kan).
- Empty state: keep `no tasks — create one` vs Kan empty `min-h-[2rem]`.
- Decide `Dropdown` vendor: reuse `lucide` `Ellipsis` vs `HiEllipsisHorizontal`.

## Non-goals

- No card styling here.

## Acceptance

- [ ] Column visually matches Kan screenshot (`screenshot.jpg` vs `web/src/components/Column.tsx` before).
- [ ] Drag column to reorder works (if sortable enabled) — or explicitly disabled if Lokan locks lane order to `statuses[]` config.
- [ ] Add-card + delete-list actions hooked (or stubbed with TODO).

## Refs

- `/tmp/kan/apps/web/src/views/board/components/List.tsx`
- `web/src/components/Column.tsx`, `ModalConfig.tsx`

## Effort

M — 1.5d

## Work Log
```

### 6 — 06 — TaskRow → Kan Card (rich card + meta footer)

```lokan
id: "6"
title: 06 — TaskRow → Kan Card (rich card + meta footer)
status: backlog
created: "2026-09-24"
updated: "2026-09-24"
tags:
    - kan-reskin
    - card
```

```markdown
# 06 — TaskRow → Kan Card (rich card + meta footer)

## Goal

Turn leaderboard `TaskRow` into Kan `Card` — biggest visual win.

## Scope

- New `Card.tsx` shell: `rounded-md border border-light-200 bg-light-50 px-3 py-2 text-sm dark:border-dark-200 dark:bg-dark-200`.
- Map Lokan `TaskSummary {id,title,tags,status,created,updated}` → Kan `Card props {title,ticketNumber,labels,members,checklists,description,comments,attachments,dueDate}`:
  - `tags → labels` (Badge + LabelIcon)
  - `id → ticketNumber` small `text-xs text-light-700`
  - `body → hasDescription` (`HiBars3BottomLeft` dot)
  - Skip members/checklists/attachments for v1 (empty), but wire slots so future `dueDate`/`checklist` props fit.
- Footer: icons `HiBars3BottomLeft`, `HiOutlineClock`+`date-fns format  do MMM`, `HiChatBubbleLeft`, `HiOutlinePaperClip`, plus `CircularProgress + completed/total` when checklist exists.
- Keep Lokan behaviors: click→`openTask`, checkbox multi-select, `copy brief` button, `moved` flash, `scopedIds` marquee highlight, `draggable`.
- Wrap with `SortableCard` (`useSortable` + `CSS.Transform`) — Kan `SortableCard.tsx` 75 LOC.

## Acceptance

- [ ] Card matches `Kan Card.tsx` visuals; tags render as `Badge` pills.
- [ ] Empty card (title only) renders clean (no footer).
- [ ] Drag handle + click vs checkbox logic preserved.

## Refs

- `/tmp/kan/apps/web/src/views/board/components/Card.tsx`, `SortableCard.tsx`, `Badge.tsx`, `LabelIcon.tsx`
- `web/src/components/TaskRow.tsx`, `lib/types.ts`

## Effort

M — 2d

## Work Log
```

### 7 — 07 — Primitives: Badge / Avatar / CircularProgress

```lokan
id: "7"
title: '07 — Primitives: Badge / Avatar / CircularProgress'
status: backlog
created: "2026-09-24"
updated: "2026-09-24"
tags:
    - kan-reskin
    - ui
```

```markdown
# 07 — Primitives: Badge / Avatar / CircularProgress

## Goal

Vendor Kan's small UI primitives needed by Card/List; keep them isolated and reusable.

## Scope

- `Badge.tsx` (Kan 399B) → `web/src/components/Kan/Badge.tsx` (`value + iconLeft`, pill style).
- `Avatar.tsx` (1.6k) → initials-from-email fallback, `size sm/md`, `imageUrl` via `getAvatarUrl` (adapt to `id` initial for Lokan single-user).
- `CircularProgress.tsx` (1.4k) → `progress 0-100`, `size sm`.
- `LabelIcon.tsx` + `helpers.getInitialsFromName` + `labelColours` if needed for `colourCode`.
- Use `tailwind-merge` for class merging where Kan does.
- No Tippy/Tooltip here (task 8/9).

## Acceptance

- [ ] Each primitive renders in isolation (story-like preview or `Mockup` html).
- [ ] No `react-icons` dependency — map to `lucide-react` or keep `hi2` if small.

## Refs

- `/tmp/kan/apps/web/src/components/Badge.tsx`, `Avatar.tsx`, `CircularProgress.tsx`, `LabelIcon.tsx`
- `/tmp/kan/apps/web/src/utils/helpers.ts`, `labelColours.ts`

## Effort

S — 1d

## Work Log
```

### 8 — 08 — Chrome & modals polish + brutalist decision

```lokan
id: "8"
title: 08 — Chrome & modals polish + brutalist decision
status: backlog
created: "2026-09-24"
updated: "2026-09-24"
tags:
    - kan-reskin
    - ux
```

```markdown
# 08 — Chrome & modals polish + brutalist decision

## Goal

Decide how far to take the reskin beyond board; polish `Topline`, `ModalDetail`, `ModalConfig`, `BulkBar` either stay brutalist or adopt Kan neutrals.

## Scope

- `Topline.tsx` (1.6k): keep Lokan wordmark + `+ NEW TASK` accent vs Kan `Dashboard + WorkspaceMenu + SideNavigation`. Recommend keep Lokan chrome — less to port.
- `ModalDetail.tsx` (7.6k) vs Kan `views/card/index.tsx` (595 LOC, heavy): keep Lokan modal but adopt Kan card detail patterns if cheap (labels selector, member row).
- `ModalConfig` / lane rename: keep current but style with Kan `Button variants` + `Toggle`.
- `BulkBar` / `ModalConfirm`: ensure `rounded-md` vs sharp decision applied consistently.
- Write ADR: `brutalist shell + Kan board` hybrid is intentional; document deviation from `docs/design/tokens.md` radius/shadow rules.

## Acceptance

- [ ] ADR/note in `docs/design/tokens.md` or `AGENTS.md` explains hybrid choice.
- [ ] Modals render with new tokens, no visual clash (accent stays fill-only).

## Refs

- `web/src/components/Topline.tsx`, `ModalDetail.tsx`, `ModalConfig.tsx`, `BulkBar.tsx`
- `docs/design/tokens.md`, `docs/architecture.md`

## Effort

S — 1d

## Work Log
```

### 9 — 09 — Keyboard shortcuts lite (C, ?, Esc)

```lokan
id: "9"
title: 09 — Keyboard shortcuts lite (C, ?, Esc)
status: todo
created: "2026-09-24"
updated: "2026-09-24"
tags:
    - kan-reskin
    - a11y
```

```markdown
# 09 — Keyboard shortcuts lite (C, ?, Esc)

## Goal

Steal Kan's shortcut _idea_ without vendoring 17.9k tree.

## Scope

- Keep Lokan `Esc` clears selection; add `C` → `setCreating(true)` (Kan `createListShortcut` pattern `useKeyboardShortcut`).
- Optional: `?` → command palette stub (Kan `CommandPallette 7.6k` — skip full port, just `?` shows help modal listing shortcuts).
- Optional: `n` → new task (alt), `f` → focus filter, arrow keys via `KeyboardSensor sortableKeyboardCoordinates` already covered by dnd task.
- Implement via `useEffect keydown` in `App.tsx` or tiny `hooks/useHotkey.ts` (10 LOC), not Kan's `ShortcutTree`.
- Add `Tooltip` content hints: `C Create new card` like Kan `createListShortcutTooltipContent`.

## Acceptance

- [ ] `C` opens create modal from board; `Esc` still clears selection; `?` opens help (or deferred).
- [ ] No focus trap/conflict with `Editor`/inputs.

## Refs

- `/tmp/kan/apps/web/src/providers/keyboard-shortcuts.tsx`, `hooks/useEventListener.ts`
- `web/src/App.tsx` Esc handler

## Effort

XS — 0.5d

## Work Log
```

### 10 — 10 — Dark mode & responsive parity

```lokan
id: "10"
title: 10 — Dark mode & responsive parity
status: backlog
created: "2026-09-24"
updated: "2026-09-24"
tags:
    - kan-reskin
    - theme
```

```markdown
# 10 — Dark mode & responsive parity

## Goal

Ensure Kan neutrals have dark variants (`dark:bg-dark-*`) and `<900px` still usable.

## Scope

- Map Kan `dark:` classes: `dark:border-dark-200 dark:bg-dark-200 dark:text-dark-1000 dark:hover:bg-dark-300` etc — via `tokens.css` `[data-theme='dark']` vars (`--bg:#000 --fg:#fff --border:#1f1f1f --muted:#777 --zebra:#111`) already exists; add Kan-specific dark vars.
- Test `document.documentElement.dataset.theme` toggle across Board/List/Card.
- Responsive: validate `max-[900px]:grid-cols-1` → new `flex` + horizontal scroll still works on narrow; decide to keep stack or keep scroll+ snap.
- Add `globals.css` `html {font-size:14px}` vs Lokan current — align.

## Acceptance

- [ ] Light + dark screenshots reviewed (no yellow-on-white contrast regression).
- [ ] Mobile (<900px) board usable (scroll or stack, no cutoff).

## Refs

- `web/src/tokens.css` `[data-theme='dark']`, `web/src/index.css`
- `/tmp/kan/apps/web/src/styles/globals.css`

## Effort

S — 1d

## Work Log
```

### 11 — 11 — Build, e2e & single-binary smoke

```lokan
id: "11"
title: 11 — Build, e2e & single-binary smoke
status: backlog
created: "2026-09-24"
updated: "2026-09-24"
tags:
    - kan-reskin
    - build
```

```markdown
# 11 — Build, e2e & single-binary smoke

## Goal

Prove `runtask build` still produces a single `dist/lokan` with embedded `web/dist` and e2e smoke passes.

## Scope

- Run `./runtask web build` → `vite build` (check hashed assets), `cp -r web/dist engine/web/dist`, `go build -o dist/lokan ./cmd/lokan`.
- Run `./runtask test` (Go), `./runtask e2e` (init, create, list, ui, API), plus manual: `dist/lokan ui <board>` drag & persist reload.
- Check `runtask lint` / `runtask format` (prettier + gofmt) — new `Badge/Avatar/CircularProgress` files follow repo style (lowercase section comments, no retro-add).
- Guard `engine/web/dist/.gitkeep` still works for pre-build `go build`.

## Acceptance

- [ ] `runtask build && runtask test && runtask e2e` green on fresh clone.
- [ ] Dragged card order survives `fetchTasks` refresh + file reload (`board.md` lineStart/End stable).

## Refs

- `runtask`, `docs/architecture.md` Build Chain
- `engine/web/dist/.gitkeep`

## Effort

S — 0.5d

## Work Log
```

### 12 — 12 — Docs, screenshot & PR

```lokan
id: "12"
title: 12 — Docs, screenshot & PR
status: backlog
created: "2026-09-24"
updated: "2026-09-24"
tags:
    - kan-reskin
    - docs
```

```markdown
# 12 — Docs, screenshot & PR

## Goal

Ship the reskin as a clean branch + docs, ready for PR.

## Scope

- Update `docs/design/tokens.md` + `docs/architecture.md` (new deps, board layout flex vs grid, dnd rationale).
- Update `README.md` screenshot (replace/append Kan-inspired board capture).
- Run `./runtask format` on all touched `web/src/**/*.{ts,tsx,css,md}`.
- Branch hygiene per `AGENTS.md` finalize: propose `feat/kan-inspired-board`, logical commits (feat, fix, docs), exclude `CONTEXT.md`, verify `diff <old> <new> -- . ':(exclude)CONTEXT.md'` empty except board `done` diff.
- PR description template: `## 📝 Summary` (ticket/context) → `## 🛠️ Changes` → `## 🧪 Verification & Checklist`.

## Non-goals

- No merge into main from clone — hub store + delete old only.

## Acceptance

- [ ] `docs/design/mockup.html` or new `mockups/kan-lokan.html` reflects final board.
- [ ] PR ready, zero `workpool` in name, no ticket numbers unless given.

## Refs

- `AGENTS.md` finalize workflow, `docs/design/tokens.md`, `runtask format`

## Effort

XS — 0.5d

## Work Log
```

## Archive

### 2 — 02 — Vendor deps & dnd audit (@dnd-kit)

```lokan
id: "2"
title: 02 — Vendor deps & dnd audit (@dnd-kit)
status: done
created: "2026-09-24"
updated: "2026-09-24"
tags:
    - kan-reskin
    - deps
```

```markdown
# 02 — Vendor deps & dnd audit (@dnd-kit)

## Goal

Decide minimal vendor set to steal Kan mechanics without bloat; keep `dist/lokan` single binary small.

## Scope

- Audit Kan deps: `@dnd-kit/core@6.3.1 + @dnd-kit/sortable@10 + @dnd-kit/utilities`, `react-icons/hi2` vs `lucide-react`, `tailwind-merge`, `date-fns`, `framer-motion`.
- Decision: install `@dnd-kit/*` + `tailwind-merge` (for `Badge` class merging) — keep `lucide-react`, drop `react-icons` (map `Hi*` to lucide equivalents), skip `framer-motion`/`tippy.js`.
- Check bundle impact: `vite build --report` before/after.
- Add to `web/package.json` with pin, run `npm install`, ensure `runtask build` still `go:embed` works (web/dist hashed assets).
- Note AGPLv3 attribution: add `NOTICE` or file header for vendored `collision.ts` / `useDragToScroll`.

## Non-goals

- No component wiring yet.

## Acceptance

- [ ] `package.json` diff minimal (≤3 new deps).
- [ ] `vite build` + `runtask build` green, `dist/lokan` size delta noted.
- [ ] License note added.

## Refs

- `/tmp/kan/apps/web/package.json` (`@dnd-kit/*`)
- `web/package.json`, `runtask`

## Effort

XS — 0.5d

## Work Log
```

### 1 — 01 — Design tokens & palette mapping (Tw v4)

```lokan
id: "1"
title: 01 — Design tokens & palette mapping (Tw v4)
status: done
created: "2026-09-24"
updated: "2026-09-24"
tags:
    - kan-reskin
    - design
```

```markdown
# 01 — Design tokens & palette mapping (Tw v4)

## Goal

Merge Kan's soft neutral palette (Radix slate + `light-50..1000 / dark-50..1000`) into Lokan's Tailwind v4 `@theme` without breaking `tokens.css` + brutalist spec (`docs/design/tokens.md`).

## Scope

- Audit `tokens.css` (`--bg/--fg/--border/--muted/--accent/--zebra`) vs Kan `globals.css` + `tooling/tailwind-config/web.ts` (Radix cyan/slate, `light-50..1000`, `dark-50..1000`, `gradient` bg).
- Decide: **keep brutalist monochrome as default**, add Kan neutrals as additive vars (`--card-bg`, `--list-bg`, `--surface-hover`) or map `light-300 → --zebra`, `light-400 → --border`.
- Draft `web/src/tokens-kan.css` or extend `tokens.css` with `@theme { --color-card: var(--list-bg) }` so `bg-card` works.
- Port `Kan Button variants` (primary/secondary/danger/ghost) to match `modal-classes.ts` invert logic vs Kan `rounded-md shadow-sm`.
- Fork decision: radius `0` vs Kan `rounded-md` — propose `radius: 6px` for cards/lists only, keep sharp elsewhere.
- Update `web/src/index.css` `@theme` mapping accordingly.

## Non-goals

- No component rewrites here; tokens only.
- No dark-mode polish (task 10).

## Acceptance

- [ ] `tokens.css` / `kan.css` diff reviewed; `index.css` `@theme` compiles under `vite build` (Tw v4).
- [ ] Light + dark vars render correctly (`[data-theme]` toggle still works).
- [ ] `docs/design/tokens.md` updated with new vars + contrast rule preserved (`--accent` fill-only).

## Refs

- `web/src/tokens.css`, `web/src/index.css`
- `/tmp/kan/apps/web/src/styles/globals.css`, `/tmp/kan/tooling/tailwind-config/web.ts`
- Kan `Button.tsx` variants

## Effort

S — 1d

## Work Log
```
