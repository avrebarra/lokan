import { useEffect, useState } from 'react'
import type { CreateTaskInput } from '../api'
import type { Priority, TaskType } from '../types'

interface Props {
  onClose: () => void
  onCreate: (input: CreateTaskInput) => void
}

// shared utility strings for the ghost button + form controls
const buttonClass =
  'inline-flex min-h-8 items-center justify-center border border-fg bg-bg px-2.5 text-[11px] uppercase text-fg transition-colors duration-[120ms] hover:bg-fg hover:text-bg'
const fieldClass =
  'min-h-8 w-full border border-border bg-bg px-2.5 py-[9px] text-xs text-fg [border-radius:0] focus:border-fg focus:outline-none'

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
    <div
      className="fixed inset-0 z-40 flex items-center justify-center p-6"
      role="dialog"
      aria-modal="true"
      aria-label="create task"
    >
      <div
        className="absolute inset-0 cursor-pointer bg-[color-mix(in_srgb,var(--bg)_72%,transparent)]"
        onClick={onClose}
      />
      <div className="relative z-1 flex max-h-[min(92vh,840px)] w-[min(100%,440px)] flex-col border border-fg bg-bg shadow-[0_24px_80px_color-mix(in_srgb,var(--fg)_18%,transparent)]">
        <div className="flex items-center justify-between gap-6 border-b border-fg px-5 py-[18px]">
          <h2 className="font-sans text-xl font-normal leading-[1.15]">new task</h2>
          <button className={`${buttonClass} flex-none`} onClick={onClose}>
            × close
          </button>
        </div>
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
                placeholder="task-01"
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
      </div>
    </div>
  )
}
