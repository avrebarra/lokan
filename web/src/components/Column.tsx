import { useState } from 'react'
import { useDroppable } from '@dnd-kit/core'
import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { Ellipsis, Plus } from 'lucide-react'
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable'
import type { Status, TaskSummary } from '../lib/types'
import { getListBodyId } from '../lib/dnd/ids'
import type { ListBodyDragData, ListDragData } from '../lib/dnd/types'
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
  onToggleSelect: (id: string) => void
  onAddCard?: (status: Status) => void
}

export default function Column({
  label,
  status,
  rows,
  movedId,
  selectedIds,
  scopedIds,
  onSelect,
  onToggleSelect,
  onAddCard,
}: Props) {
  const [menuOpen, setMenuOpen] = useState(false)

  // sortable scaffold for lane — disabled for now: lane order is locked to
  // statuses[] config order (Kan allows drag-to-reorder, Lokan does not yet).
  // keep hook so styling/behaviour matches Kan List.tsx; enable by flipping
  // disabled.draggable to true and handling LIST drop in Board.tsx.
  const {
    attributes,
    listeners,
    setNodeRef: setSortableRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: status,
    data: { type: 'LIST' } satisfies ListDragData,
    disabled: { draggable: true, droppable: true },
  })

  const sortableStyle: React.CSSProperties = {
    transform: CSS.Transform.toString(transform),
    transition,
  }

  const { setNodeRef: setDroppableRef, isOver } = useDroppable({
    id: getListBodyId(status),
    data: { type: 'LIST_BODY', listPublicId: status } satisfies ListBodyDragData,
  })

  const rowIds = rows.map((r) => r.id)

  return (
    <div
      data-board-draggable
      ref={setSortableRef}
      style={sortableStyle}
      {...attributes}
      {...listeners}
      className={`flex h-fit min-w-[18rem] max-w-[18rem] flex-col rounded-md border border-list-border bg-list py-2 pl-2 pr-1 snap-start dark:border-dark-300 dark:bg-dark-100 ${isDragging ? 'cursor-grabbing opacity-60' : ''} ${isOver ? 'ring-1 ring-border' : ''}`}
      aria-label={label}
    >
      {/* header — Kan List header: inline name + + add + … menu */}
      <div className="mb-2 flex items-center justify-between gap-1">
        <div className="min-w-0 flex-1">
          <input
            type="text"
            value={label}
            readOnly
            aria-label={`List ${label}`}
            className="w-full border-0 bg-transparent px-3 pt-1 text-sm font-medium text-fg focus:outline-none focus:ring-0"
          />
        </div>
        <div className="flex shrink-0 items-center">
          <button
            onClick={() => onAddCard?.(status)}
            aria-label={`Add card to ${label}`}
            className="mx-1 inline-flex h-7 w-7 items-center justify-center rounded-md text-muted hover:bg-surface-hover hover:text-fg disabled:cursor-not-allowed disabled:opacity-60"
          >
            <Plus className="h-[18px] w-[18px]" aria-hidden="true" />
          </button>
          <div className="relative mr-1">
            <button
              onClick={() => setMenuOpen((v) => !v)}
              aria-label={`List options for ${label}`}
              aria-expanded={menuOpen}
              className="flex h-7 w-7 items-center justify-center rounded-[5px] hover:bg-light-200 dark:hover:bg-dark-200"
            >
              <Ellipsis className="h-5 w-5 text-muted" aria-hidden="true" />
            </button>
            {menuOpen && (
              <>
                <button
                  className="fixed inset-0 z-10 cursor-default"
                  aria-hidden="true"
                  onClick={() => setMenuOpen(false)}
                  tabIndex={-1}
                />
                <div className="absolute right-0 z-20 mt-1 w-44 origin-top-right rounded-md border border-light-200 bg-white p-1 shadow-lg ring-1 ring-black/5 dark:border-dark-400 dark:bg-dark-300">
                  <button
                    onClick={() => {
                      setMenuOpen(false)
                      onAddCard?.(status)
                    }}
                    className="flex w-full items-center gap-2 rounded-[5px] px-2.5 py-1.5 text-left text-sm text-fg hover:bg-light-200 dark:hover:bg-dark-400"
                  >
                    <Plus className="h-[16px] w-[16px]" aria-hidden="true" />
                    Add card
                  </button>
                  <div className="my-1 border-t border-border" />
                  <span className="block px-2.5 py-1 text-[11px] uppercase text-muted">
                    Delete via board config
                  </span>
                </div>
              </>
            )}
          </div>
        </div>
        <span className="mr-1 shrink-0 font-sans text-xs text-muted tabular-nums">
          {rows.length}
        </span>
      </div>

      {/* body — droppable + sortable cards */}
      <div ref={setDroppableRef} className={`relative min-h-[2rem] ${isOver ? 'bg-zebra/50' : ''}`}>
        {rows.length === 0 ? (
          <div className="px-3 py-3 text-[11px] uppercase text-muted">no tasks — create one</div>
        ) : (
          <SortableContext items={rowIds} strategy={verticalListSortingStrategy}>
            <div className="flex flex-col gap-2 pr-1">
              {rows.map((row) => (
                <TaskRow
                  key={row.id}
                  task={row}
                  listId={status}
                  moved={row.id === movedId}
                  selectedIds={selectedIds}
                  scopedIds={scopedIds}
                  onClick={() => onSelect(row.id)}
                  onToggleSelect={() => onToggleSelect(row.id)}
                />
              ))}
            </div>
          </SortableContext>
        )}
      </div>
    </div>
  )
}
