import type { Status, TaskSummary } from '../types'
import Column from './Column'

interface Props {
  tasks: TaskSummary[]
  subtaskCount: Map<string, number>
  onSelect: (id: string) => void
  onMove: (id: string, status: Status, beforeId?: string) => void
}

interface ColumnSpec {
  label: string
  statuses: Status[]
  modifier?: string
}

// column layout: one column per status, ordered to read naturally
const COLUMNS: ColumnSpec[] = [
  { label: 'backlog', statuses: ['backlog'] },
  { label: 'todo', statuses: ['todo'] },
  { label: 'in-progress', statuses: ['in-progress'], modifier: 'in-progress' },
  { label: 'done', statuses: ['done'] },
  { label: 'cancelled', statuses: ['cancelled'] },
]

export default function Board({ tasks, subtaskCount, onSelect, onMove }: Props) {
  // subtasks live under their parents, not as board rows
  const visible = tasks.filter((t) => t.type !== 'subtask')
  return (
    <main className="grid grid-cols-5 items-start gap-6 max-[900px]:grid-cols-1 max-[900px]:gap-10">
      {COLUMNS.map((col) => (
        <Column
          key={col.label}
          label={col.label}
          modifier={col.modifier}
          status={col.statuses[0]}
          rows={visible.filter((t) => col.statuses.includes(t.status))}
          subtaskCount={subtaskCount}
          onSelect={onSelect}
          onMove={onMove}
        />
      ))}
    </main>
  )
}
