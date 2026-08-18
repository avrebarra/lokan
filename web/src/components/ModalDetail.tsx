import { useState } from 'react'
import type { UpdateTaskInput } from '../lib/api'
import type { StatusDef, Task, Status } from '../lib/types'
import Modal from './Modal'
import ModalConfirm from './ModalConfirm'
import { buttonClass, confirmClass, fieldClass } from '../lib/modal-classes'

export interface TaskFieldChange {
  field: UpdateTaskInput['field']
  value: string
}

interface Props {
  task: Task
  statuses: StatusDef[]
  onClose: () => void
  onSave: (changes: TaskFieldChange[]) => void
  onDelete: () => void
}

interface Draft {
  title: string
  status: Status
  tags: string
  body: string
}

export default function ModalDetail({ task, statuses, onClose, onSave, onDelete }: Props) {
  // edit mode + draft form state
  const [editing, setEditing] = useState(false)
  const [draft, setDraft] = useState<Draft | null>(null)
  const [confirmingDelete, setConfirmingDelete] = useState(false)

  // enter edit mode, seeding the draft from the current task
  const startEdit = () => {
    setDraft({
      title: task.title,
      status: task.status,
      tags: task.tags?.join(', ') ?? '',
      body: task.body,
    })
    setEditing(true)
  }

  // collect changed fields and hand them up, then exit edit mode
  const save = () => {
    if (!draft) return
    const changes: TaskFieldChange[] = []
    if (draft.title !== task.title) changes.push({ field: 'title', value: draft.title })
    if (draft.status !== task.status) changes.push({ field: 'status', value: draft.status })
    if (draft.tags !== (task.tags?.join(', ') ?? '')) {
      changes.push({ field: 'tags', value: draft.tags })
    }
    if (draft.body !== task.body) changes.push({ field: 'body', value: draft.body })
    if (changes.length > 0) onSave(changes)
    setEditing(false)
  }

  const cancel = () => {
    setEditing(false)
    setDraft(null)
  }

  return (
    <>
      <Modal
        title={
          <>
            {editing ? (
              <input
                type="text"
                value={draft?.title ?? ''}
                onChange={(e) => setDraft({ ...draft!, title: e.target.value })}
                aria-label="edit title"
                autoFocus
                className={fieldClass}
              />
            ) : (
              <h2 className="font-sans text-2xl font-normal leading-[1.15]">{task.title}</h2>
            )}
            <div className="mt-2 flex flex-wrap items-baseline gap-2">
              <span className="text-[11px] uppercase text-muted">{task.id}</span>
            </div>
          </>
        }
        onClose={onClose}
        ariaLabel="task detail"
        maxWidth="w-[min(100%,680px)]"
        headerAlign="start"
        headerRight={
          <button className={`${buttonClass} flex-none`} onClick={onClose}>
            × close
          </button>
        }
        footer={
          <div className="flex flex-wrap gap-2 border-t border-border px-5 pb-5 pt-4">
            {editing ? (
              <>
                <button className={buttonClass} onClick={save}>
                  save
                </button>
                <button className={buttonClass} onClick={cancel}>
                  cancel
                </button>
              </>
            ) : (
              <>
                <button className={buttonClass} onClick={startEdit}>
                  edit
                </button>
                <button className={confirmClass} onClick={() => setConfirmingDelete(true)}>
                  delete
                </button>
              </>
            )}
          </div>
        }
      >
        {/* body: fields and notes */}
        <div className="overflow-y-auto px-5 pb-5">
          <div className="mt-4 grid grid-cols-2 border-t border-border max-[900px]:grid-cols-1">
            <div className="flex items-baseline gap-6 border-b border-border py-2.5 text-[11px] uppercase">
              <span className="w-20 flex-none text-muted">status</span>
              {editing ? (
                <div className="w-1/2 min-w-0">
                  <select
                    value={draft?.status ?? task.status}
                    onChange={(e) => setDraft({ ...draft!, status: e.target.value as Status })}
                    aria-label="edit status"
                    className={fieldClass}
                  >
                    {statuses.map((s) => (
                      <option key={s.id} value={s.id}>
                        {s.id}
                      </option>
                    ))}
                  </select>
                </div>
              ) : (
                <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                  {task.status}
                </span>
              )}
            </div>
            <div className="flex items-baseline gap-6 border-b border-border py-2.5 text-[11px] uppercase">
              <span className="w-20 flex-none text-muted">created</span>
              <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                {task.created}
              </span>
            </div>
            <div className="flex items-baseline gap-6 border-b border-border py-2.5 text-[11px] uppercase">
              <span className="w-20 flex-none text-muted">updated</span>
              <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                {task.updated}
              </span>
            </div>
            <div className="flex items-baseline gap-6 border-b border-border py-2.5 text-[11px] uppercase">
              <span className="w-20 flex-none text-muted">tags</span>
              {editing ? (
                <div className="w-1/2 min-w-0">
                  <input
                    type="text"
                    value={draft?.tags ?? ''}
                    onChange={(e) => setDraft({ ...draft!, tags: e.target.value })}
                    aria-label="edit tags"
                    placeholder="alpha, beta"
                    className={fieldClass}
                  />
                </div>
              ) : (
                <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                  {task.tags?.join(', ') ?? '—'}
                </span>
              )}
            </div>
            <div className="flex items-baseline gap-6 border-b border-border py-2.5 text-[11px] uppercase">
              <span className="w-20 flex-none text-muted">related</span>
              <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                {task.related?.join(', ') ?? '—'}
              </span>
            </div>
          </div>
          {editing ? (
            <>
              <h3 className="mt-6 border-t border-fg pt-3 text-[11px] uppercase">notes</h3>
              <textarea
                value={draft?.body ?? ''}
                onChange={(e) => setDraft({ ...draft!, body: e.target.value })}
                aria-label="edit body"
                rows={6}
                className={`${fieldClass} mt-3.5 min-h-32 resize-y`}
              />
            </>
          ) : (
            task.body && (
              <>
                <h3 className="mt-6 border-t border-fg pt-3 text-[11px] uppercase">notes</h3>
                <p className="mt-3.5 whitespace-pre-wrap text-[13px] leading-[1.55]">{task.body}</p>
              </>
            )
          )}
        </div>
      </Modal>
      {confirmingDelete && (
        <ModalConfirm
          title="Delete task"
          message={`Delete task ${task.id} — ${task.title}? This cannot be undone.`}
          confirmLabel="delete"
          onConfirm={() => {
            setConfirmingDelete(false)
            onDelete()
          }}
          onCancel={() => setConfirmingDelete(false)}
        />
      )}
    </>
  )
}
