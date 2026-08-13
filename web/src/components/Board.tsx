import type { Status, StatusDef, TaskSummary } from '../types'
import Column from './Column'

interface Props {
  statuses: StatusDef[]
  tasks: TaskSummary[]
  subtaskCount: Map<string, number>
  onSelect: (id: string) => void
  onMove: (id: string, status: Status, beforeId?: string) => void
}

// one column per configured lane, in config order; the in-progress lane keeps
// the accent top bar
export default function Board({ statuses, tasks, subtaskCount, onSelect, onMove }: Props) {
  // subtasks live under their parents, not as board rows
  const visible = tasks.filter((t) => t.type !== 'subtask')
  return (
    <main
      className={`grid items-start gap-6 max-[900px]:grid-cols-1 max-[900px]:gap-10`}
      style={{ gridTemplateColumns: `repeat(${statuses.length}, minmax(0, 1fr))` }}
    >
      {statuses.map((col) => (
        <Column
          key={col.id}
          label={col.id}
          modifier={col.id === 'in-progress' ? 'in-progress' : undefined}
          status={col.id}
          rows={visible.filter((t) => t.status === col.id)}
          subtaskCount={subtaskCount}
          onSelect={onSelect}
          onMove={onMove}
        />
      ))}
    </main>
  )
}
