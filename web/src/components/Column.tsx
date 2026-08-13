import type { TaskSummary } from '../types'
import TaskRow from './TaskRow'

interface Props {
  label: string
  modifier?: string
  rows: TaskSummary[]
  subtaskCount: Map<string, number>
  onSelect: (id: string) => void
}

export default function Column({ label, modifier, rows, subtaskCount, onSelect }: Props) {
  // accent bar on the in-progress column head
  const headClass = modifier
    ? 'mb-1 flex items-baseline justify-between border-t-2 border-accent pt-3 text-[11px] font-normal uppercase text-fg'
    : 'mb-1 flex items-baseline justify-between border-t border-fg pt-3 text-[11px] font-normal uppercase'
  return (
    <section className="flex flex-col" aria-label={label}>
      <h2 className={headClass}>
        {label} <span className="font-sans text-sm text-muted">{rows.length}</span>
      </h2>
      {rows.length === 0 ? (
        <div className="py-3.5 text-[11px] uppercase text-muted">
          no tasks — create one with lokan create
        </div>
      ) : (
        rows.map((row) => (
          <TaskRow
            key={row.id}
            task={row}
            subtaskCount={subtaskCount.get(row.id) ?? 0}
            onClick={() => onSelect(row.id)}
          />
        ))
      )}
    </section>
  )
}
