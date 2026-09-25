import { useDroppable } from '@dnd-kit/core'
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable'
import type { Status, TaskSummary } from '../lib/types'
import { getListBodyId } from '../lib/dnd/ids'
import type { ListBodyDragData } from '../lib/dnd/types'
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
  onToggleSelect,
}: Props) {
  // accent bar on the in-progress column head
  const headClass = modifier
    ? 'mb-1 flex items-baseline justify-between border-t-8 border-accent pl-2.5 pt-3 text-[13px] font-normal uppercase text-fg'
    : 'mb-1 flex items-baseline justify-between border-t-8 border-fg pl-2.5 pt-3 text-[13px] font-normal uppercase'

  const { setNodeRef, isOver } = useDroppable({
    id: getListBodyId(status),
    data: { type: 'LIST_BODY', listPublicId: status } satisfies ListBodyDragData,
  })

  const rowIds = rows.map((r) => r.id)

  return (
    <section className="flex flex-col overflow-hidden" aria-label={label}>
      <h2 className={headClass}>
        {label}
        <span className="font-sans text-muted">{rows.length}</span>
      </h2>
      <div ref={setNodeRef} className={`relative min-h-12 ${isOver ? 'bg-zebra' : ''}`}>
        {rows.length === 0 ? (
          <div className="py-3.5 text-[11px] uppercase text-muted">
            no tasks — create one with lokan create
          </div>
        ) : (
          <SortableContext items={rowIds} strategy={verticalListSortingStrategy}>
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
          </SortableContext>
        )}
      </div>
    </section>
  )
}
