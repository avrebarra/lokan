import { useState } from 'react'
import type { StatusDef } from '../lib/types'
import Modal from './Modal'
import { buttonClass, confirmClass, fieldClass } from '../lib/modal-classes'
import ModalConfirm from './ModalConfirm'

interface Props {
  statuses: StatusDef[]
  laneCounts: Record<string, number>
  theme: 'light' | 'dark'
  onSetTheme: (theme: 'light' | 'dark') => void
  onClose: () => void
  onSave: (statuses: StatusDef[]) => void
  onClearArchived: () => void
  onClearAll: () => void
}

interface Row {
  key: string
  id: string
  // archived is preserved from the current config — lanes seeded as archived
  // (done/cancelled) stay archived through renames; new lanes are active
  archived: boolean
}

// pending confirmation: close-without-save or a destructive action
type Confirm =
  | { kind: 'close' }
  | { kind: 'remove'; row: Row }
  | { kind: 'clear-archived' }
  | { kind: 'clear-all' }

// primary action: accent fill, same as the top-level CTA
const primaryClass =
  'inline-flex min-h-8 items-center justify-center border border-accent bg-accent px-2.5 text-[11px] font-medium uppercase text-black transition-colors duration-[120ms] hover:border-fg hover:bg-fg hover:text-bg'
// destructive action: red border + text, red fill on hover
const dangerClass =
  'inline-flex min-h-8 items-center justify-center border border-danger px-2.5 text-[11px] uppercase text-danger transition-colors duration-[120ms] hover:bg-danger hover:text-bg'

function uid(): string {
  return Math.random().toString(36).slice(2)
}

// lanes in saveable form, for dirty detection against the server state
function serialize(rows: Row[]): StatusDef[] {
  return rows.map((r) => ({ id: r.id.trim(), archived: r.archived }))
}

export default function ModalConfig({
  statuses,
  laneCounts,
  theme,
  onSetTheme,
  onClose,
  onSave,
  onClearArchived,
  onClearAll,
}: Props) {
  // editable copy of the lane set, keyed for stable rows
  const [rows, setRows] = useState<Row[]>(() =>
    statuses.map((s) => ({ key: uid(), id: s.id, archived: s.archived })),
  )
  const [error, setError] = useState('')
  const [confirm, setConfirm] = useState<Confirm | null>(null)

  // unsaved edits block a plain close
  const dirty = JSON.stringify(serialize(rows)) !== JSON.stringify(statuses)

  // close only when clean; otherwise confirm discarding first
  const requestClose = () => {
    if (dirty) {
      setConfirm({ kind: 'close' })
    } else {
      onClose()
    }
  }

  // update a row field in place
  const patch = (key: string, changes: Partial<Row>) => {
    setRows((prev) => prev.map((r) => (r.key === key ? { ...r, ...changes } : r)))
  }

  // append an empty lane; ids validate on save
  const addLane = () => setRows((prev) => [...prev, { key: uid(), id: '', archived: false }])

  // validate and hand the trimmed lane set up
  const save = () => {
    const ids = rows.map((r) => r.id.trim())
    if (ids.some((id) => id === '')) {
      setError('Every lane needs a name.')
      return
    }
    if (new Set(ids).size !== ids.length) {
      setError('Lane names must be unique.')
      return
    }
    setError('')
    onSave(serialize(rows))
  }

  // counts for the confirmation copy
  const totalTasks = Object.values(laneCounts).reduce((a, b) => a + b, 0)
  const archivedCount = rows.reduce((sum, r) => sum + (r.archived ? (laneCounts[r.id] ?? 0) : 0), 0)
  const confirmCopy = (() => {
    switch (confirm?.kind) {
      case 'close':
        return {
          title: 'Discard changes?',
          message: 'You have unsaved lane changes. Discard them?',
          confirmLabel: 'discard',
        }
      case 'remove': {
        const n = laneCounts[confirm.row.id] ?? 0
        return {
          title: `Remove lane "${confirm.row.id}"?`,
          message:
            n > 0
              ? `${n} task${n === 1 ? '' : 's'} will move to the first lane.`
              : 'This lane is empty.',
          confirmLabel: 'remove',
        }
      }
      case 'clear-archived':
        return {
          title: 'Clear archived?',
          message: `Delete ${archivedCount} archived task${archivedCount === 1 ? '' : 's'}? This cannot be undone.`,
          confirmLabel: 'delete',
        }
      case 'clear-all':
        return {
          title: 'Clear all?',
          message: `Delete all ${totalTasks} task${totalTasks === 1 ? '' : 's'}? This cannot be undone.`,
          confirmLabel: 'delete',
        }
      default:
        return null
    }
  })()

  return (
    <Modal
      title={
        <div>
          <h2 className="font-sans text-2xl font-normal leading-[1.15]">config</h2>
          <p className="mt-1 text-[11px] uppercase text-muted">lanes · bulk ops</p>
        </div>
      }
      onClose={requestClose}
      ariaLabel="board config"
      maxWidth="w-[min(100%,680px)]"
      escapeDisabled={confirm !== null}
      footer={
        <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border px-5 pb-5 pt-4">
          <div className="flex flex-wrap gap-2">
            <button className={primaryClass} onClick={save}>
              save settings
            </button>
            <button className={buttonClass} onClick={requestClose}>
              cancel
            </button>
          </div>
          <div className="flex flex-wrap gap-2">
            <button className={dangerClass} onClick={() => setConfirm({ kind: 'clear-archived' })}>
              clear archived
            </button>
            <button className={dangerClass} onClick={() => setConfirm({ kind: 'clear-all' })}>
              clear all
            </button>
          </div>
        </div>
      }
    >
      {/* body: theme, then lanes */}
      <div className="overflow-y-auto px-5 pb-5">
        <div className="mt-4 border-b border-fg pb-2 text-[11px] uppercase text-muted">theme</div>
        <div className="flex gap-2 py-2.5">
          {(['light', 'dark'] as const).map((t) => (
            <button
              key={t}
              className={theme === t ? confirmClass : buttonClass}
              onClick={() => onSetTheme(t)}
              aria-pressed={theme === t}
            >
              {t}
            </button>
          ))}
        </div>
        <div className="mt-4 border-b border-fg pb-2 text-[11px] uppercase text-muted">lanes</div>
        {rows.map((row) => (
          <div key={row.key} className="flex items-center gap-3 py-1">
            <div className="min-w-0 flex-1">
              <input
                type="text"
                value={row.id}
                onChange={(e) => patch(row.key, { id: e.target.value })}
                aria-label="lane name"
                className={fieldClass}
              />
            </div>
            <button
              className={`${buttonClass} flex-none`}
              onClick={() => setConfirm({ kind: 'remove', row })}
            >
              remove
            </button>
          </div>
        ))}
        <div className="mt-3 flex items-center justify-between">
          <button className={buttonClass} onClick={addLane}>
            + add lane
          </button>
          {error && <span className="text-[11px] uppercase text-fg">{error}</span>}
        </div>
      </div>

      {/* confirmation layer above the config panel */}
      {confirmCopy && confirm && (
        <ModalConfirm
          title={confirmCopy.title}
          message={confirmCopy.message}
          confirmLabel={confirmCopy.confirmLabel}
          onCancel={() => setConfirm(null)}
          onConfirm={() => {
            setConfirm(null)
            switch (confirm.kind) {
              case 'close':
                onClose()
                break
              case 'remove':
                setRows((prev) => prev.filter((r) => r.key !== confirm.row.key))
                break
              case 'clear-archived':
                onClearArchived()
                break
              case 'clear-all':
                onClearAll()
                break
            }
          }}
        />
      )}
    </Modal>
  )
}
