import { useCallback, useEffect, useRef, useState } from 'react'
import type { Status, StatusDef, TaskSummary } from '../lib/types'
import Column from './Column'

interface Props {
  statuses: StatusDef[]
  tasks: TaskSummary[]
  movedId: string | null
  selectedIds: Set<string>
  onSelect: (id: string) => void
  onMove: (ids: string[], status: Status, beforeId?: string) => void
  onToggleSelect: (id: string) => void
  onMarqueeSelect: (ids: string[]) => void
}

interface Marquee {
  x1: number
  y1: number
  x2: number
  y2: number
}

// true when two rectangles overlap (marquee box vs a row's box)
const rectsOverlap = (a: { x: number; y: number; w: number; h: number }, b: DOMRect) =>
  a.x < b.right && a.x + a.w > b.left && a.y < b.bottom && a.y + a.h > b.top

// one column per configured lane, in config order; the in-progress lane keeps
// the accent top bar
export default function Board({
  statuses,
  tasks,
  movedId,
  selectedIds,
  onSelect,
  onMove,
  onToggleSelect,
  onMarqueeSelect,
}: Props) {
  const [marquee, setMarquee] = useState<Marquee | null>(null)
  const marqueeStart = useRef<{ x: number; y: number } | null>(null)
  // rows inside the live marquee box — state, so rows re-render and highlight
  // while the rectangle is being dragged, not just on release
  const [marqueeLive, setMarqueeLive] = useState<Set<string>>(new Set())
  const onMarqueeSelectRef = useRef(onMarqueeSelect)
  useEffect(() => {
    onMarqueeSelectRef.current = onMarqueeSelect
  })

  // start a rubber-band selection anywhere — page background, board gaps,
  // below the board — but never on interactive elements or an open modal
  const startMarquee = useCallback((e: PointerEvent) => {
    const target = e.target as HTMLElement
    if (target.closest('button, input, [data-row], [role="dialog"], [role="alertdialog"]')) return
    marqueeStart.current = { x: e.clientX, y: e.clientY }
    setMarquee({ x1: e.clientX, y1: e.clientY, x2: e.clientX, y2: e.clientY })
    setMarqueeLive(new Set())
  }, [])

  // track the box and live-highlight the rows it intersects
  const moveMarquee = useCallback((e: PointerEvent) => {
    if (!marqueeStart.current) return
    setMarquee((m) => (m ? { ...m, x2: e.clientX, y2: e.clientY } : m))
    const box = {
      x: Math.min(marqueeStart.current.x, e.clientX),
      y: Math.min(marqueeStart.current.y, e.clientY),
      w: Math.abs(e.clientX - marqueeStart.current.x),
      h: Math.abs(e.clientY - marqueeStart.current.y),
    }
    const next = new Set<string>()
    document.querySelectorAll<HTMLElement>('[data-row]').forEach((el) => {
      const id = el.dataset.rowId
      if (id && rectsOverlap(box, el.getBoundingClientRect())) next.add(id)
    })
    setMarqueeLive(next)
  }, [])

  // commit the boxed rows computed fresh from the final pointer position —
  // reading state here would be stale (useCallback closes over the first render)
  const endMarquee = useCallback((e: PointerEvent) => {
    if (!marqueeStart.current) return
    const start = marqueeStart.current
    const dx = Math.abs(e.clientX - start.x)
    const dy = Math.abs(e.clientY - start.y)
    if (dx + dy >= 6) {
      const box = {
        x: Math.min(start.x, e.clientX),
        y: Math.min(start.y, e.clientY),
        w: Math.abs(e.clientX - start.x),
        h: Math.abs(e.clientY - start.y),
      }
      const ids: string[] = []
      document.querySelectorAll<HTMLElement>('[data-row]').forEach((el) => {
        const id = el.dataset.rowId
        if (id && rectsOverlap(box, el.getBoundingClientRect())) ids.push(id)
      })
      onMarqueeSelectRef.current(ids)
    }
    marqueeStart.current = null
    setMarqueeLive(new Set())
    setMarquee(null)
  }, [])

  // document-level listeners so the marquee can start outside the board
  useEffect(() => {
    document.addEventListener('pointerdown', startMarquee)
    document.addEventListener('pointermove', moveMarquee)
    document.addEventListener('pointerup', endMarquee)
    document.addEventListener('pointercancel', endMarquee)
    return () => {
      document.removeEventListener('pointerdown', startMarquee)
      document.removeEventListener('pointermove', moveMarquee)
      document.removeEventListener('pointerup', endMarquee)
      document.removeEventListener('pointercancel', endMarquee)
    }
  }, [startMarquee, moveMarquee, endMarquee])

  return (
    <>
      <main
        className={`grid select-none items-start gap-6 max-[900px]:grid-cols-1 max-[900px]:gap-10`}
        style={{ gridTemplateColumns: `repeat(${statuses.length}, minmax(0, 1fr))` }}
      >
        {statuses.map((col) => (
          <Column
            key={col.id}
            label={col.id}
            modifier={col.id === 'in-progress' ? 'in-progress' : undefined}
            status={col.id}
            rows={tasks.filter((t) => t.status === col.id)}
            movedId={movedId}
            selectedIds={selectedIds}
            scopedIds={marqueeLive}
            onSelect={onSelect}
            onMove={onMove}
            onToggleSelect={onToggleSelect}
          />
        ))}
      </main>
      {marquee && (
        <div
          className="pointer-events-none fixed z-20 border border-accent bg-accent/10"
          style={{
            left: Math.min(marquee.x1, marquee.x2),
            top: Math.min(marquee.y1, marquee.y2),
            width: Math.abs(marquee.x2 - marquee.x1),
            height: Math.abs(marquee.y2 - marquee.y1),
          }}
        />
      )}
    </>
  )
}
