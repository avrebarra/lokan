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
  const className = modifier ? `column ${modifier}` : 'column'
  return (
    <section className={className} aria-label={label}>
      <h2 className="column-head">
        {label} <span className="count">{rows.length}</span>
      </h2>
      {rows.length === 0 ? (
        <div className="empty">no tasks — create one with lokan create</div>
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
