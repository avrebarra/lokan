import type { Status, TaskSummary } from '../types'
import Column from './Column'

interface Props {
  tasks: TaskSummary[]
  subtaskCount: Map<string, number>
  onSelect: (id: string) => void
}

interface ColumnSpec {
  label: string
  statuses: Status[]
  modifier?: string
}

// column layout: statuses bucketed per column, modifier drives styling
const COLUMNS: ColumnSpec[] = [
  { label: 'todo', statuses: ['todo', 'backlog', 'cancelled'] },
  { label: 'in-progress', statuses: ['in-progress'], modifier: 'in-progress' },
  { label: 'done', statuses: ['done'] },
]

export default function Board({ tasks, subtaskCount, onSelect }: Props) {
  // subtasks live under their parents, not as board rows
  const visible = tasks.filter((t) => t.type !== 'subtask')
  return (
    <main className="board">
      {COLUMNS.map((col) => (
        <Column
          key={col.label}
          label={col.label}
          modifier={col.modifier}
          rows={visible.filter((t) => col.statuses.includes(t.status))}
          subtaskCount={subtaskCount}
          onSelect={onSelect}
        />
      ))}
    </main>
  )
}
