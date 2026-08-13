# lokan — Design Tokens (hi-fi spec)

> Source of truth for the frontend build (G4). Derived from shiprank
> (brutalist-terminal aesthetic) + lokan's kanban data model. Mockup:
> `web/design/mockup.html`. Approved by user on 2026-08-13.

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

| token            | value     | use                                    |
|------------------|-----------|----------------------------------------|
| `--bg`           | `#fff`    | page background                        |
| `--fg`           | `#000`    | text, strong borders                   |
| `--border`       | `#ebebeb` | card borders, hairlines                |
| `--muted`        | `#888`    | meta text, ids, secondary              |
| `--accent`       | `#ffc800` | in-progress marker ONLY (sparingly)    |
| `--zebra`        | `#f5f5f5` | subtle row hover / zebra               |

### Dark

| token      | value     |
|------------|-----------|
| `--bg`     | `#000`    |
| `--fg`     | `#fff`    |
| `--border` | `#1f1f1f` |
| `--muted`  | `#777`    |
| `--accent` | `#ffc800` |

Toggle: `data-theme="dark"` on `<html>`; respect `prefers-color-scheme`
as default, explicit toggle overrides. (No theme switcher in mockup — ships
in G4; default light + auto dark.)

## 3. Typography

| role          | family     | size  | weight | transform | tracking     |
|---------------|-----------|-------|--------|-----------|--------------|
| wordmark      | Geist Sans| 22px  | 700    | —         | `-0.01em`    |
| column head   | Geist Mono| 11px  | 400    | uppercase | `0`          |
| section head  | Geist Mono| 11px  | 400    | uppercase | `0`          |
| card title    | Geist Sans| 14px  | 400    | —         | `0`          |
| card meta/id  | Geist Mono| 11px  | 400    | uppercase | `0`          |
| tag / badge   | Geist Mono| 9px   | 400    | uppercase | `0`          |
| button        | Geist Mono| 11px  | 400    | uppercase | `0`          |
| empty state   | Geist Mono| 11px  | 400    | uppercase | `0`          |

Fallbacks: Geist Sans → `Helvetica, Arial, sans-serif`; Geist Mono →
`"Courier New", monospace`. **Load from Google Fonts** (same source shiprank
uses via `next/font/google`):
`https://fonts.googleapis.com/css2?family=Geist:wght@400..700&family=Geist+Mono:wght@400..500&display=swap`
(preconnect to fonts.googleapis.com + fonts.gstatic.com). Do NOT use jsdelivr
geist package — Google Fonts is the canonical source.

## 4. Shape & rules

- **radius: 0 everywhere.** No rounded corners, no pills.
- **1px solid borders**, color `--border` for hairlines, `--fg` for
  structural lines (topline, section tops, primary buttons).
- **No box-shadows** (except detail modal — see 7).
- **Rows, not cards (DECIDED 2026-08-13):** tasks render as shiprank
  leaderboard `.row` entries — `border-bottom: 1px solid var(--border)`,
  `padding: 11px 0`, hover = zebra bg + title underline. NO boxed cards.
- **Contrast rule (DECIDED 2026-08-13):** `--accent` `#ffc800` is **fill-only
  — never text.** Yellow text on white = 1.9:1, fails WCAG. Allowed uses:
  (1) in-progress column header `border-top: 2px solid var(--accent)` bar
  (header TEXT stays `--fg` black — mirrors shiprank's WED bar, which is a
  filled bar, not colored text), (2) primary CTA button bg with `#000` text
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
  Matches shiprank's dual use (WED bar + `.github-connect-button`). Everything
  else stays monochrome.
- Body: `display: flex; flex-direction: column`.

### Row (task entry — shiprank leaderboard pattern, NO boxes)
```
task-05 · MED                  [2]      ← mono meta: id · priority tag · subtask count
Add cycle detection to buildTree         ← Geist Sans title
───────────────────────────────          ← border-bottom: 1px solid var(--border)
```
- `border-bottom: 1px solid var(--border)`, `padding: 11px 0`,
  `background: var(--bg)` — rows touch, separated by hairline only.
- `.row-main`: `flex; align-items: baseline; justify-content: space-between`
  (meta left, count right).
- Hover: `background: var(--zebra)`; title underline.
- Priority tag: all `1px solid var(--border)`, CRIT differs by `color: var(--fg)`.
- Subtask count: `[N]` in `--muted` mono, only when `N > 0`, `flex: 0 0 auto`.
- Row click → detail modal (G4).

### Detail modal (task detail data — DECIDED 2026-08-13)
Clicking a row opens a centered modal (shiprank share-layer pattern):
- Backdrop: `color-mix(in srgb, var(--bg) 72%, transparent)`.
- Panel: `1px solid var(--fg)`, `max-width: 680px`, `max-height: min(92vh, 840px)`,
  the ONLY allowed shadow: `0 24px 80px color-mix(in srgb, var(--fg) 18%, transparent)`.
- Head: title (Geist Sans 24px) + id line (`task-01 · CRIT · bug`) + `× close` button.
- Fields grid (2 cols, 1 col narrow): **status, type, priority, parent, created,
  updated, tags, related** — label `--muted` mono 11px uppercase left, value
  right-aligned 12px, `border-bottom` hairlines.
- `notes` subhead (border-top `--fg` section-head) → body prose 13px.
- `subtasks` subhead → list rows: title + `--muted` mono id, hairline-separated.
- Footer actions: `advance → <next status>`, `edit`, `+ subtask` (buttons).

### Status cycle (interaction contract for G4)
Click advances: `todo → in-progress → done` (+ `backlog`, `cancelled` per
contract — board shows 3 core columns: TODO / IN-PROGRESS / DONE;
backlog/cancelled tasks render in TODO column with tag, or filtered —
**G4 decision, mockup shows core 3**).

### Buttons
- `.button`: `1px solid var(--fg)`, radius 0, mono 11px uppercase,
  `padding: 0 10px; min-height: 32px`; hover = **invert** (bg→fg, fg→bg).
- `.button.accent` (primary CTA): `background: var(--accent)`,
  `border-color: var(--accent)`, `color: #000`; hover invert to `--fg`.

### Empty state
Mono uppercase `--muted`: `no tasks — create one with lokan create`.

## 6. Layout

### Wide (≥ 900px)
- `.wrap`: `max-width: 1200px; margin: 0 auto; padding: 32px 24px 64px`.
- Board: `display: grid; grid-template-columns: repeat(3, 1fr); gap: 24px`.

### Narrow (< 900px) — shiprank single-column stack
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
- Numbers: Geist Sans (counts in column heads, subtask counts).

## 9. Grid of elements to copy into G4

`web/design/mockup.html` is the visual contract — G4 must reproduce it with
React components, same tokens, same structure. Deviations need approval.
