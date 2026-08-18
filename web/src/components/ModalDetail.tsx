import { useState } from 'react'
import type { UpdateTaskInput } from '../lib/api'
import type { StatusDef, Task, TaskSummary, Status, Priority, TaskType } from '../lib/types'
import { PRIORITIES, TASK_TYPES } from '../lib/types'
import { priorityTag } from '../lib/format'
import Modal from './Modal'
import { buttonClass, fieldClass } from '../lib/modal-classes'

export interface TaskFieldChange {
  field: UpdateTaskInput['field']
  value: string
}

interface Props {
  task: Task
  subtasks: TaskSummary[]
  statuses: StatusDef[]
  onClose: () => void
  onSave: (changes: TaskFieldChange[]) => void
  onAddSubtask?: () => void
}

// bordered micro-tag in the header row
const tagClass = 'border border-border px-[5px] py-px text-[9px] uppercase leading-[1.2]'

interface Draft {
  title: string
  status: Status
  priority: Priority
  type: TaskType
  parent: string
  tags: string
  body: string
}

export default function ModalDetail({
  task,
  subtasks,
  statuses,
  onClose,
  onSave,
  onAddSubtask,
}: Props) {
  // edit mode + draft form state
  const [editing, setEditing] = useState(false)
  const [draft, setDraft] = useState<Draft | null>(null)

  // derived display values for the header
  const crit = task.priority === 'critical'

  // enter edit mode, seeding the draft from the current task
  const startEdit = () => {
    setDraft({
      title: task.title,
      status: task.status,
      priority: task.priority,
      type: task.type,
      parent: task.parent ?? '',
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
    if (draft.priority !== task.priority) changes.push({ field: 'priority', value: draft.priority })
    if (draft.type !== task.type) changes.push({ field: 'type', value: draft.type })
    if (draft.parent !== (task.parent ?? '')) changes.push({ field: 'parent', value: draft.parent })
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

  // subtasks are only allowed under task/bug parents (see types.go AllowedParents)
  const canHaveSubtasks = task.type === 'task' || task.type === 'bug'

  return (
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
            <span className={crit ? `${tagClass} text-fg` : `${tagClass} text-muted`}>
              {priorityTag(task.priority)}
            </span>
            <span className={`${tagClass} text-muted`}>{task.type}</span>
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
              {canHaveSubtasks && onAddSubtask && (
                <button className={buttonClass} onClick={onAddSubtask}>
                  + subtask
                </button>
              )}
            </>
          )}
        </div>
      }
    >
      {/* body: fields, notes, subtasks */}
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
            <span className="w-20 flex-none text-muted">type</span>
            {editing ? (
              <div className="w-1/2 min-w-0">
                <select
                  value={draft?.type ?? task.type}
                  onChange={(e) => setDraft({ ...draft!, type: e.target.value as TaskType })}
                  aria-label="edit type"
                  className={fieldClass}
                >
                  {TASK_TYPES.map((t) => (
                    <option key={t} value={t}>
                      {t}
                    </option>
                  ))}
                </select>
              </div>
            ) : (
              <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                {task.type}
              </span>
            )}
          </div>
          <div className="flex items-baseline gap-6 border-b border-border py-2.5 text-[11px] uppercase">
            <span className="w-20 flex-none text-muted">priority</span>
            {editing ? (
              <div className="w-1/2 min-w-0">
                <select
                  value={draft?.priority ?? task.priority}
                  onChange={(e) => setDraft({ ...draft!, priority: e.target.value as Priority })}
                  aria-label="edit priority"
                  className={fieldClass}
                >
                  {PRIORITIES.map((p) => (
                    <option key={p} value={p}>
                      {p}
                    </option>
                  ))}
                </select>
              </div>
            ) : (
              <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                {task.priority}
              </span>
            )}
          </div>
          <div className="flex items-baseline gap-6 border-b border-border py-2.5 text-[11px] uppercase">
            <span className="w-20 flex-none text-muted">parent</span>
            {editing ? (
              <div className="w-1/2 min-w-0">
                <input
                  type="text"
                  value={draft?.parent ?? ''}
                  onChange={(e) => setDraft({ ...draft!, parent: e.target.value })}
                  aria-label="edit parent"
                  placeholder="1"
                  className={fieldClass}
                />
              </div>
            ) : (
              <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                {task.parent ?? '—'}
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
        {subtasks.length > 0 && (
          <>
            <h3 className="mt-6 border-t border-fg pt-3 text-[11px] uppercase">
              subtasks <span className="font-sans text-sm text-muted">{subtasks.length}</span>
            </h3>
            {subtasks.map((st) => (
              <div
                className="flex items-baseline justify-between gap-6 border-b border-border py-[9px] text-xs"
                key={st.id}
              >
                <span>{st.title}</span>
                <span className="text-[11px] text-muted">{st.id}</span>
              </div>
            ))}
          </>
        )}
      </div>
    </Modal>
  )
}
