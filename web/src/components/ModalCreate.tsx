import { useState } from 'react'
import type { CreateTaskInput } from '../api'
import type { Priority, TaskType } from '../types'
import Modal from './Modal'
import { buttonClass, fieldClass } from './modal-classes'

interface Props {
  onClose: () => void
  onCreate: (input: CreateTaskInput) => void
  initialParent?: string
  initialType?: TaskType
}

export default function ModalCreate({ onClose, onCreate, initialParent, initialType }: Props) {
  // form state for the new task fields
  const [title, setTitle] = useState('')
  const [type, setType] = useState<TaskType>(initialType ?? 'task')
  const [priority, setPriority] = useState<Priority>('medium')
  const [parent, setParent] = useState(initialParent ?? '')

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
    <Modal
      title={<h2 className="font-sans text-xl font-normal leading-[1.15]">new task</h2>}
      onClose={onClose}
      ariaLabel="create task"
      maxWidth="w-[min(100%,440px)]"
      headerRight={
        <button className={`${buttonClass} flex-none`} onClick={onClose}>
          × close
        </button>
      }
    >
      <form onSubmit={submit}>
        <div className="overflow-y-auto p-5">
          <div className="mb-4 flex flex-col gap-1.5">
            <label htmlFor="create-title" className="text-[11px] uppercase text-muted">
              title
            </label>
            <input
              id="create-title"
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              autoFocus
              required
              className={fieldClass}
            />
          </div>
          <div className="mb-4 flex flex-col gap-1.5">
            <label htmlFor="create-type" className="text-[11px] uppercase text-muted">
              type
            </label>
            <select
              id="create-type"
              value={type}
              onChange={(e) => setType(e.target.value as TaskType)}
              className={fieldClass}
            >
              <option value="task">task</option>
              <option value="bug">bug</option>
              <option value="epic">epic</option>
              <option value="subtask">subtask</option>
            </select>
          </div>
          <div className="mb-4 flex flex-col gap-1.5">
            <label htmlFor="create-priority" className="text-[11px] uppercase text-muted">
              priority
            </label>
            <select
              id="create-priority"
              value={priority}
              onChange={(e) => setPriority(e.target.value as Priority)}
              className={fieldClass}
            >
              <option value="critical">critical</option>
              <option value="high">high</option>
              <option value="medium">medium</option>
              <option value="low">low</option>
            </select>
          </div>
          <div className="mb-4 flex flex-col gap-1.5">
            <label htmlFor="create-parent" className="text-[11px] uppercase text-muted">
              parent (optional)
            </label>
            <input
              id="create-parent"
              type="text"
              value={parent}
              onChange={(e) => setParent(e.target.value)}
              placeholder="1"
              className={fieldClass}
            />
          </div>
        </div>
        <div className="flex flex-wrap gap-2 border-t border-border px-5 pb-5 pt-4">
          <button className={buttonClass} type="submit">
            create
          </button>
          <button className={buttonClass} type="button" onClick={onClose}>
            cancel
          </button>
        </div>
      </form>
    </Modal>
  )
}
