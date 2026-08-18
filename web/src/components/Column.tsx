import { useRef, useState } from 'react'
import type { DragEvent } from 'react'
import type { Status, TaskSummary } from '../lib/types'
import TaskRow from './TaskRow'

interface Props {
  label: string
  modifier?: string
  status: Status
  rows: TaskSummary[]
  movedId: string | null
  selectedIds: Set<string>
  scopedIds: Set<string>
  onSelect: (id: string) => void
  onMove: (ids: string[], status: Status, beforeId?: string) => void
  onToggleSelect: (id: string) => void
}

interface Insertion {
  index: number
  top: number
}

export default function Column({
  label,
  modifier,
  status,
  rows,
  movedId,
  selectedIds,
  scopedIds,
  onSelect,
  onMove,
  onToggleSelect,
}: Props) {
  const [dragOver, setDragOver] = useState(false)
  const [insert, setInsert] = useState<Insertion | null>(null)
  const rowsRef = useRef<HTMLDivElement>(null)
  // accent bar on the in-progress column head
  const headClass = modifier
    ? 'mb-1 flex items-baseline justify-between border-t-8 border-accent pl-2.5 pt-3 text-[13px] font-normal uppercase text-fg'
    : 'mb-1 flex items-baseline justify-between border-t-8 border-fg pl-2.5 pt-3 text-[13px] font-normal uppercase'

  // allow the drop and track the insertion point from the pointer vs rows
  const handleDragOver = (e: DragEvent) => {
    e.preventDefault()
    e.dataTransfer.dropEffect = 'move'
    const el = rowsRef.current
    if (!el) return
    setDragOver(true)
    const containerTop = el.getBoundingClientRect().top
    let index = rows.length
    let top = el.offsetHeight
    // row buttons only — cards carry a data-row marker so the copy button
    // and other affordances don't shift the insertion index
    const rowButtons = el.querySelectorAll('button[data-row]')
    for (let i = 0; i < rowButtons.length; i++) {
      const r = rowButtons[i].getBoundingClientRect()
      if (e.clientY < r.top + r.height * 0.7) {
        index = i
        top = r.top - containerTop
        break
      }
    }
    setInsert({ index, top })
  }

  // clear the highlight only when leaving the column, not its children
  const handleDragLeave = (e: DragEvent) => {
    if (!e.currentTarget.contains(e.relatedTarget as Node)) {
      setDragOver(false)
      setInsert(null)
    }
  }

  // accept a lane move: read the carried task ids (JSON array from a group
  // drag, or a bare id from an older payload) and hand them up
  const handleDrop = (e: DragEvent) => {
    e.preventDefault()
    setDragOver(false)
    const raw = e.dataTransfer.getData('text/x-lokan-task')
    let ids: string[] = []
    try {
      const parsed = JSON.parse(raw)
      if (Array.isArray(parsed)) ids = parsed
    } catch {
      if (raw) ids = [raw]
    }
    const beforeId = insert ? (rows[insert.index]?.id ?? '') : ''
    setInsert(null)
    if (ids.length > 0) onMove(ids, status, beforeId)
  }

  return (
    <section
      className="flex flex-col overflow-hidden"
      aria-label={label}
      onDragEnter={handleDragOver}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
    >
      <h2 className={headClass}>
        {label}
        <span className="font-sans text-muted">{rows.length}</span>
      </h2>
      <div ref={rowsRef} className={`relative min-h-12 ${dragOver ? 'bg-zebra' : ''}`}>
        {dragOver && insert && (
          <div
            className="pointer-events-none absolute inset-x-0 z-10 h-[2px] bg-fg"
            style={{ top: insert.top }}
          />
        )}
        {rows.length === 0 ? (
          <div className="py-3.5 text-[11px] uppercase text-muted">
            no tasks — create one with lokan create
          </div>
        ) : (
          rows.map((row) => (
            <TaskRow
              key={row.id}
              task={row}
              moved={row.id === movedId}
              selectedIds={selectedIds}
              scopedIds={scopedIds}
              onClick={() => onSelect(row.id)}
              onToggleSelect={() => onToggleSelect(row.id)}
            />
          ))
        )}
      </div>
    </section>
  )
}
