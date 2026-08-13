import { useEffect } from 'react'
import type { Task, TaskSummary } from '../types'
import { nextStatus, priorityTag } from '../format'

interface Props {
  task: Task
  subtasks: TaskSummary[]
  onClose: () => void
  onAdvance: () => void
}

// shared utility strings for the bordered micro-tag + ghost button
const tagClass = 'border border-border px-[5px] py-px text-[9px] uppercase leading-[1.2]'
const buttonClass =
  'inline-flex min-h-8 items-center justify-center border border-fg bg-bg px-2.5 text-[11px] uppercase text-fg transition-colors duration-[120ms] hover:bg-fg hover:text-bg'

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
    <div
      className="fixed inset-0 z-40 flex items-center justify-center p-6"
      role="dialog"
      aria-modal="true"
      aria-label="task detail"
    >
      <div
        className="absolute inset-0 cursor-pointer bg-[color-mix(in_srgb,var(--bg)_72%,transparent)]"
        onClick={onClose}
      />
      <div className="relative z-1 flex max-h-[min(92vh,840px)] w-[min(100%,680px)] flex-col border border-fg bg-bg shadow-[0_24px_80px_color-mix(in_srgb,var(--fg)_18%,transparent)]">
        {/* header: title, id, tags */}
        <div className="flex items-start justify-between gap-6 border-b border-fg px-5 py-[18px]">
          <div>
            <h2 className="font-sans text-2xl font-normal leading-[1.15]">{task.title}</h2>
            <div className="mt-2 flex flex-wrap items-baseline gap-2">
              <span className="text-[11px] uppercase text-muted">{task.id}</span>
              <span className={crit ? `${tagClass} text-fg` : `${tagClass} text-muted`}>
                {priorityTag(task.priority)}
              </span>
              <span className={`${tagClass} text-muted`}>{task.type}</span>
            </div>
          </div>
          <button className={`${buttonClass} flex-none`} onClick={onClose}>
            × close
          </button>
        </div>
        {/* body: fields, notes, subtasks */}
        <div className="overflow-y-auto px-5 pb-5">
          <div className="mt-4 grid grid-cols-2 border-t border-border max-[900px]:grid-cols-1">
            <div className="flex items-baseline justify-between gap-4 border-b border-border py-2.5 text-[11px] uppercase">
              <span className="text-muted">status</span>
              <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                {task.status}
              </span>
            </div>
            <div className="flex items-baseline justify-between gap-4 border-b border-border py-2.5 text-[11px] uppercase">
              <span className="text-muted">type</span>
              <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                {task.type}
              </span>
            </div>
            <div className="flex items-baseline justify-between gap-4 border-b border-border py-2.5 text-[11px] uppercase">
              <span className="text-muted">priority</span>
              <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                {task.priority}
              </span>
            </div>
            <div className="flex items-baseline justify-between gap-4 border-b border-border py-2.5 text-[11px] uppercase">
              <span className="text-muted">parent</span>
              <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                {task.parent ?? '—'}
              </span>
            </div>
            <div className="flex items-baseline justify-between gap-4 border-b border-border py-2.5 text-[11px] uppercase">
              <span className="text-muted">created</span>
              <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                {task.created}
              </span>
            </div>
            <div className="flex items-baseline justify-between gap-4 border-b border-border py-2.5 text-[11px] uppercase">
              <span className="text-muted">updated</span>
              <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                {task.updated}
              </span>
            </div>
            <div className="flex items-baseline justify-between gap-4 border-b border-border py-2.5 text-[11px] uppercase">
              <span className="text-muted">tags</span>
              <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                {task.tags?.join(', ') ?? '—'}
              </span>
            </div>
            <div className="flex items-baseline justify-between gap-4 border-b border-border py-2.5 text-[11px] uppercase">
              <span className="text-muted">related</span>
              <span className="text-right text-xs normal-case [overflow-wrap:anywhere]">
                {task.related?.join(', ') ?? '—'}
              </span>
            </div>
          </div>
          {task.body && (
            <>
              <h3 className="mt-6 border-t border-fg pt-3 text-[11px] uppercase">notes</h3>
              <p className="mt-3.5 text-[13px] leading-[1.55]">{task.body}</p>
            </>
          )}
          {subtasks.length > 0 && (
            <>
              <h3 className="mt-6 border-t border-fg pt-3 text-[11px] uppercase">
                subtasks <span className="font-sans text-sm text-muted">{subtasks.length}</span>
              </h3>
              {subtasks.map((st) => (
                <div
                  className="flex items-baseline justify-between gap-4 border-b border-border py-[9px] text-xs"
                  key={st.id}
                >
                  <span>{st.title}</span>
                  <span className="text-[11px] text-muted">{st.id}</span>
                </div>
              ))}
            </>
          )}
        </div>
        {/* actions: advance + placeholders for edit/subtask */}
        <div className="flex flex-wrap gap-2 border-t border-border px-5 pb-5 pt-4">
          <button className={buttonClass} onClick={onAdvance}>
            advance → {next}
          </button>
          <button className={buttonClass}>edit</button>
          <button className={buttonClass}>+ subtask</button>
        </div>
      </div>
    </div>
  )
}
