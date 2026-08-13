import { useEffect, useState } from 'react'
import type { CreateTaskInput } from '../api'
import type { Priority, TaskType } from '../types'

interface Props {
  onClose: () => void
  onCreate: (input: CreateTaskInput) => void
}

export default function CreateModal({ onClose, onCreate }: Props) {
  // form state for the new task fields
  const [title, setTitle] = useState('')
  const [type, setType] = useState<TaskType>('task')
  const [priority, setPriority] = useState<Priority>('medium')
  const [parent, setParent] = useState('')

  // close on escape
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

  // build the create payload from the form, ignoring empty title
  const submit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!title.trim()) return
    onCreate({
      title: title.trim(),
      type,
      priority,
      parent: parent.trim() || undefined,
    })
  }

  return (
    <div className="detail-layer" role="dialog" aria-modal="true" aria-label="create task">
      <div className="detail-backdrop" onClick={onClose} />
      <div className="create-panel">
        <div className="create-head">
          <h2>new task</h2>
          <button className="button detail-close" onClick={onClose}>
            × close
          </button>
        </div>
        <form onSubmit={submit}>
          <div className="create-body">
            <div className="create-field">
              <label htmlFor="create-title">title</label>
              <input
                id="create-title"
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                autoFocus
                required
              />
            </div>
            <div className="create-field">
              <label htmlFor="create-type">type</label>
              <select
                id="create-type"
                value={type}
                onChange={(e) => setType(e.target.value as TaskType)}
              >
                <option value="task">task</option>
                <option value="bug">bug</option>
                <option value="epic">epic</option>
              </select>
            </div>
            <div className="create-field">
              <label htmlFor="create-priority">priority</label>
              <select
                id="create-priority"
                value={priority}
                onChange={(e) => setPriority(e.target.value as Priority)}
              >
                <option value="critical">critical</option>
                <option value="high">high</option>
                <option value="medium">medium</option>
                <option value="low">low</option>
              </select>
            </div>
            <div className="create-field">
              <label htmlFor="create-parent">parent (optional)</label>
              <input
                id="create-parent"
                type="text"
                value={parent}
                onChange={(e) => setParent(e.target.value)}
                placeholder="task-01"
              />
            </div>
          </div>
          <div className="create-actions">
            <button className="button" type="submit">
              create
            </button>
            <button className="button" type="button" onClick={onClose}>
              cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
