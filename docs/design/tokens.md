# lokan — Design Tokens (hi-fi spec)

> Source of truth for the frontend build (G4). Derived from a brutalist-terminal
> aesthetic + lokan's kanban data model. Mockup:
> `docs/design/mockup.html`. Approved by user on 2026-08-13.

## 1. Design read

Kanban board dev tool for a solo developer. Brutalist-terminal language:
monochrome, sharp corners, 1px rules, zero decoration, mono typography for
labels. One yellow focal point reserved for the in-progress state. Feels like
a terminal that grew a UI — not a "product dashboard".

Dials: VARIANCE 5 · MOTION 2 · DENSITY 5.

## 2. Color tokens

Consumed via `web/src/tokens.css`; mapped into Tailwind utilities
(`bg-bg`, `text-fg`, ...) through `@theme` in `web/src/index.css`.

### Light (default)

| token      | value     | use                                 |
| ---------- | --------- | ----------------------------------- |
| `--bg`     | `#fff`    | page background                     |
| `--fg`     | `#000`    | text, strong borders                |
| `--border` | `#ebebeb` | card borders, hairlines             |
| `--muted`  | `#888`    | meta text, ids, secondary           |
| `--accent` | `#ffc800` | in-progress marker ONLY (sparingly) |
| `--zebra`  | `#f5f5f5` | subtle row hover / zebra            |

### Dark

| token      | value     |
| ---------- | --------- |
| `--bg`     | `#000`    |
| `--fg`     | `#fff`    |
| `--border` | `#1f1f1f` |
| `--muted`  | `#777`    |
| `--accent` | `#ffc800` |

Toggle: `data-theme="dark"` on `<html>`; default light, explicit toggle opts
into dark. (Stored preference wins when present.)

### Kan neutrals (additive, 2026-09-24 — task 01)

Imported from Kan `tooling/tailwind/web.ts` to support Kan-style cards/lists
without replacing brutalist core. Raw palette is **theme-agnostic**; semantic
aliases switch per `[data-theme]`.

**Raw palette** (`--light-*` / `--dark-*` in `tokens.css`):

| light token    | value                | dark token    | value     |
| -------------- | -------------------- | ------------- | --------- |
| `--light-50`   | `hsl(0deg 0% 98.8%)` | `--dark-50`   | `#161616` |
| `--light-100`  | `hsl(0deg 0% 97.3%)` | `--dark-100`  | `#1c1c1c` |
| `--light-200`  | `hsl(0deg 0% 95.3%)` | `--dark-200`  | `#232323` |
| `--light-300`  | `hsl(0deg 0% 92.9%)` | `--dark-300`  | `#282828` |
| `--light-400`  | `hsl(0deg 0% 91%)`   | `--dark-400`  | `#2e2e2e` |
| `--light-500`  | `hsl(0deg 0% 88.6%)` | `--dark-500`  | `#343434` |
| `--light-600`  | `hsl(0deg 0% 85.9%)` | `--dark-600`  | `#3e3e3e` |
| `--light-700`  | `hsl(0deg 0% 78%)`   | `--dark-700`  | `#505050` |
| `--light-800`  | `hsl(0deg 0% 56.1%)` | `--dark-800`  | `#707070` |
| `--light-900`  | `hsl(0deg 0% 52.2%)` | `--dark-900`  | `#7e7e7e` |
| `--light-950`  | `hsl(0deg 0% 43.5%)` | `--dark-950`  | `#bbb`    |
| `--light-1000` | `hsl(0deg 0% 9%)`    | `--dark-1000` | `#ededed` |

**Semantic aliases** (mapped in `tokens.css`, exposed as Tailwind `bg-card` etc via `@theme` in `index.css`):

| alias             | light value        | dark value        | use                                | tailwind                 |
| ----------------- | ------------------ | ----------------- | ---------------------------------- | ------------------------ |
| `--card-bg`       | `var(--light-50)`  | `var(--dark-200)` | card surface (Kan `Card`)          | `bg-card` `text-card`    |
| `--card-border`   | `var(--light-200)` | `var(--dark-200)` | card border                        | `border-card-border`     |
| `--list-bg`       | `var(--light-300)` | `var(--dark-100)` | column/list container (Kan `List`) | `bg-list`                |
| `--list-border`   | `var(--light-400)` | `var(--dark-300)` | list border                        | `border-list-border`     |
| `--surface-hover` | `var(--light-400)` | `var(--dark-300)` | hover on card/list                 | `hover:bg-surface-hover` |

Mapping to brutalist core (for reference): `light-300 (~#ededed)` ≈ `--zebra #f5f5f5`,
`light-400 (~#e8e8e8)` ≈ `--border #ebebeb` — close but kept distinct so brutalist
shell is untouched. Tailwind exposes both sets: `bg-zebra` (brutalist) and `bg-list` (Kan).

Tw v4 `@theme` mapping: `--color-card: var(--card-bg)` etc in `index.css` — enables
`bg-card`, `border-list-border`, `bg-light-200`, `dark:bg-dark-200` utilities.

## 3. Typography

| role         | family     | size | weight | transform | tracking  |
| ------------ | ---------- | ---- | ------ | --------- | --------- |
| wordmark     | Geist Sans | 22px | 700    | —         | `-0.01em` |
| column head  | Geist Mono | 11px | 400    | uppercase | `0`       |
| section head | Geist Mono | 11px | 400    | uppercase | `0`       |
| card title   | Geist Sans | 14px | 400    | —         | `0`       |
| card meta/id | Geist Mono | 11px | 400    | uppercase | `0`       |
| tag / badge  | Geist Mono | 9px  | 400    | uppercase | `0`       |
| button       | Geist Mono | 11px | 400    | uppercase | `0`       |
| empty state  | Geist Mono | 11px | 400    | uppercase | `0`       |

Fallbacks: Geist Sans → `Helvetica, Arial, sans-serif`; Geist Mono →
`"Courier New", monospace`. **Load from Google Fonts** (via `next/font/google`):
`https://fonts.googleapis.com/css2?family=Geist:wght@400..700&family=Geist+Mono:wght@400..500&display=swap`
(preconnect to fonts.googleapis.com + fonts.gstatic.com). Do NOT use jsdelivr
geist package — Google Fonts is the canonical source.

## 4. Shape & rules

- **radius: 0 default, 6px for Kan surfaces (DECIDED 2026-09-24 — task 01).** Brutalist shell (topline, buttons, modals) stays sharp `0`. Kan cards/lists use `--radius-card: 6px` (`rounded-md` in Kan) — exposed as `--radius-card` / `rounded-card` via `@theme`. Pills use `--radius-pill: 9999px` for `Badge`. No other radii.
- **1px solid borders**, color `--border` for hairlines, `--fg` for
  structural lines (topline, section tops, primary buttons).
- **No box-shadows** (except detail modal — see 7).
- **Rows, not cards (DECIDED 2026-08-13):** tasks render as
  leaderboard `.row` entries — `border-bottom: 1px solid var(--border)`,
  `padding: 11px 0`, hover = zebra bg + title underline. NO boxed cards.
- **Contrast rule (DECIDED 2026-08-13):** `--accent` `#ffc800` is **fill-only
  — never text.** Yellow text on white = 1.9:1, fails WCAG. Allowed uses:
  (1) in-progress column header `border-top: 2px solid var(--accent)` bar
  (header TEXT stays `--fg` black — the bar is fill-only, not colored text), (2) primary CTA button bg with `#000` text
  (black-on-yellow ≈ 14:1 ✓).

## 5. Components

### Topline (sticky)

- Bottom border `1px solid var(--fg)`, `padding: 14px 0`.
- Left: wordmark `lokan` (Geist Sans, bold, `-0.01em`).
- Right: actions — `+ NEW TASK` button, subtle `theme` toggle (text,
  not icon), meta line `N TASKS · updated HH:MM`.

### Column

- Header: uppercase mono status label + count, `border-top: 1px solid var(--fg)`.
- **In-progress header: `border-top: 2px solid var(--accent)` bar; TEXT stays
  `--fg`** (contrast rule — yellow is fill-only).
- **Accent scope (DECIDED 2026-08-13):** exactly two allowed accent uses — (1)
  in-progress column header bar, (2) primary CTA button (`+ NEW TASK`).
  Matches the dual use — accent bar + primary button, both fill-only. Everything
  else stays monochrome.
- Body: `display: flex; flex-direction: column`.

### Row (task entry — leaderboard pattern, NO boxes)

```
task-05                        ← mono meta: id
Add cycle detection to buildTree         ← Geist Sans title
───────────────────────────────          ← border-bottom: 1px solid var(--border)
```

- `border-bottom: 1px solid var(--border)`, `padding: 11px 0`,
  `background: var(--bg)` — rows touch, separated by hairline only.
- `.row-main`: `flex; align-items: baseline; justify-content: space-between`
  (meta left, count right).
- Hover: `background: var(--zebra)`; title underline.
- Row click → detail modal (G4).

### Detail modal (task detail data — DECIDED 2026-08-13)

Clicking a row opens a centered modal (share-layer pattern):

- Backdrop: `color-mix(in srgb, var(--bg) 72%, transparent)`.
- Panel: `1px solid var(--fg)`, `max-width: 680px`, `max-height: min(92vh, 840px)`,
  the ONLY allowed shadow: `0 24px 80px color-mix(in srgb, var(--fg) 18%, transparent)`.
- Head: title (Geist Sans 24px) + id line + `× close` button.
- Fields grid (2 cols, 1 col narrow): **status, created, updated, tags,
  related** — label `--muted` mono 11px uppercase left, value
  right-aligned 12px, `border-bottom` hairlines.
- `notes` subhead (border-top `--fg` section-head) → body prose 13px.
- Footer actions: `edit`, `delete` (buttons). Lane moves are drag-and-drop
  on the board; the edit form has a status select.

### Status cycle (interaction contract)

Lane moves are drag-and-drop (also editable via the status select in the
detail modal). Lanes are configurable in the config modal — the board renders
one column per configured lane in order, archived lanes feed the Archive
section and bulk clear.

### Buttons

- `.button`: `1px solid var(--fg)`, radius 0, mono 11px uppercase,
  `padding: 0 10px; min-height: 32px`; hover = **invert** (bg→fg, fg→bg).
- `.button.accent` (primary CTA): `background: var(--accent)`,
  `border-color: var(--accent)`, `color: #000`; hover invert to `--fg`.
- **Kan variants (DECIDED 2026-09-24 — task 01):** Kan `Button` has `primary / secondary / danger / ghost` with `rounded-md shadow-sm` (`primary = bg-light-1000 dark:bg-dark-1000`, `secondary = border-light-600 bg-light-50`). Lokan keeps **invert** hover as the brand (brutalist), but tokens now expose the neutrals so a future `web/src/lib/modal-classes.ts` can map `secondary → border border-card-border bg-card` and `danger → bg-danger` without new CSS. No `shadow-sm` — Lokan stays flat; radius only on Kan cards/lists (`--radius-card`), not buttons (buttons stay `0`).

### Empty state

Mono uppercase `--muted`: `no tasks — create one with lokan create`.

## 6. Layout

### Wide (≥ 900px)

- `.wrap`: `max-width: 1200px; margin: 0 auto; padding: 32px 24px 64px`.
- Board: `display: grid; grid-template-columns: repeat(3, 1fr); gap: 24px`.

### Narrow (< 900px) — single-column stack

- Board becomes `grid-template-columns: 1fr` — columns **stack vertically**,
  one per section with its own border-top header. No horizontal scroll.
- Topline actions collapse: `+ NEW TASK` stays, meta line hides or wraps.
- Card meta row wraps (`flex-wrap`).

## 7. Motion (minimal)

- Hover transitions only: `background 120ms`, `border-color 120ms`, opacity.
- Modal (detail view): backdrop
  `color-mix(in srgb, var(--bg) 72%, transparent)`, panel
  `1px solid var(--fg)`, the ONLY allowed shadow:
  `box-shadow: 0 24px 80px color-mix(in srgb, var(--fg) 18%, transparent)`.
- Respect `prefers-reduced-motion`.

## 8. Content rules

- All labels UPPERCASE mono (status, buttons, tags, ids).
- Title case never forced — task titles are user text, left as-is.
- No icons/emoji. Text glyphs only (`+`, `→`).
- Numbers: Geist Sans (counts in column heads).

## 9. Grid of elements to copy into G4

`docs/design/mockup.html` is the visual contract — G4 must reproduce it with
React components, same tokens, same structure. Deviations need approval.
