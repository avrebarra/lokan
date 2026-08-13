import { useEffect } from 'react'
import type { Task, TaskSummary } from '../types'
import { nextStatus, priorityTag } from '../format'

interface Props {
  task: Task
  subtasks: TaskSummary[]
  onClose: () => void
  onAdvance: () => void
}

export default function DetailModal({ task, subtasks, onClose, onAdvance }: Props) {
  // close on escape
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

  // derived display values for the header + advance button
  const crit = task.priority === 'critical'
  const next = nextStatus(task.status)

  return (
    <div className="detail-layer" role="dialog" aria-modal="true" aria-label="task detail">
      <div className="detail-backdrop" onClick={onClose} />
      <div className="detail-panel">
        {/* header: title, id, tags */}
        <div className="detail-head">
          <div>
            <h2>{task.title}</h2>
            <div className="detail-idline">
              <span className="detail-id">{task.id}</span>
              <span className={crit ? 'detail-tag crit' : 'detail-tag'}>
                {priorityTag(task.priority)}
              </span>
              <span className="detail-tag">{task.type}</span>
            </div>
          </div>
          <button className="button detail-close" onClick={onClose}>
            × close
          </button>
        </div>
        {/* body: fields, notes, subtasks */}
        <div className="detail-body">
          <div className="detail-fields">
            <div className="detail-field">
              <span className="detail-field-label">status</span>
              <span className="detail-field-value">{task.status}</span>
            </div>
            <div className="detail-field">
              <span className="detail-field-label">type</span>
              <span className="detail-field-value">{task.type}</span>
            </div>
            <div className="detail-field">
              <span className="detail-field-label">priority</span>
              <span className="detail-field-value">{task.priority}</span>
            </div>
            <div className="detail-field">
              <span className="detail-field-label">parent</span>
              <span className="detail-field-value">{task.parent ?? '—'}</span>
            </div>
            <div className="detail-field">
              <span className="detail-field-label">created</span>
              <span className="detail-field-value">{task.created}</span>
            </div>
            <div className="detail-field">
              <span className="detail-field-label">updated</span>
              <span className="detail-field-value">{task.updated}</span>
            </div>
            <div className="detail-field">
              <span className="detail-field-label">tags</span>
              <span className="detail-field-value">{task.tags?.join(', ') ?? '—'}</span>
            </div>
            <div className="detail-field">
              <span className="detail-field-label">related</span>
              <span className="detail-field-value">{task.related?.join(', ') ?? '—'}</span>
            </div>
          </div>
          {task.body && (
            <>
              <h3 className="detail-subhead">notes</h3>
              <p className="detail-prose">{task.body}</p>
            </>
          )}
          {subtasks.length > 0 && (
            <>
              <h3 className="detail-subhead">
                subtasks <span className="count">{subtasks.length}</span>
              </h3>
              {subtasks.map((st) => (
                <div className="detail-subtask" key={st.id}>
                  <span>{st.title}</span>
                  <span className="mark">{st.id}</span>
                </div>
              ))}
            </>
          )}
        </div>
        {/* actions: advance + placeholders for edit/subtask */}
        <div className="detail-actions">
          <button className="button" onClick={onAdvance}>
            advance → {next}
          </button>
          <button className="button">edit</button>
          <button className="button">+ subtask</button>
        </div>
      </div>
    </div>
  )
}
