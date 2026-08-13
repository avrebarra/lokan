import type { TaskSummary } from '../types'
import { priorityTag } from '../format'

interface Props {
  task: TaskSummary
  subtaskCount: number
  onClick: () => void
}

export default function TaskRow({ task, subtaskCount, onClick }: Props) {
  return (
    <button className="row" onClick={onClick}>
      <div className="row-main">
        <div className="row-meta">
          <span className="row-id">{task.id}</span>
          <span className={task.priority === 'critical' ? 'row-tag crit' : 'row-tag'}>
            {priorityTag(task.priority)}
          </span>
        </div>
        {subtaskCount > 0 && <span className="row-count">[{subtaskCount}]</span>}
      </div>
      <div className="row-title">{task.title}</div>
    </button>
  )
}
