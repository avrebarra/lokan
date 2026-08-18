import { useState } from 'react'
import type { CreateTaskInput } from '../lib/api'
import Modal from './Modal'
import { buttonClass, fieldClass } from '../lib/modal-classes'

interface Props {
  onClose: () => void
  onCreate: (input: CreateTaskInput) => void
}

export default function ModalCreate({ onClose, onCreate }: Props) {
  // form state for the new task title
  const [title, setTitle] = useState('')

  // build the create payload from the form, ignoring empty title
  const submit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!title.trim()) return
    onCreate({ title: title.trim() })
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
