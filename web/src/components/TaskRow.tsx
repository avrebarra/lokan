import type { TaskSummary } from '../types'
import { priorityTag } from '../format'

interface Props {
  task: TaskSummary
  subtaskCount: number
  onClick: () => void
}

export default function TaskRow({ task, subtaskCount, onClick }: Props) {
  return (
    <button
      className="group block w-full border-b border-border bg-bg py-[11px] text-left text-fg transition-colors duration-[120ms] hover:bg-zebra"
      onClick={onClick}
    >
      <div className="flex items-baseline justify-between gap-4">
        <div className="flex flex-wrap items-baseline gap-2">
          <span className="text-[11px] uppercase text-muted">{task.id}</span>
          <span
            className={
              task.priority === 'critical'
                ? 'border border-border px-[5px] py-px text-[9px] uppercase leading-[1.2] text-fg'
                : 'border border-border px-[5px] py-px text-[9px] uppercase leading-[1.2] text-muted'
            }
          >
            {priorityTag(task.priority)}
          </span>
        </div>
        {subtaskCount > 0 && (
          <span className="flex-none text-[11px] text-muted">[{subtaskCount}]</span>
        )}
      </div>
      <div className="mt-[3px] font-sans text-sm font-normal leading-[1.35] group-hover:underline">
        {task.title}
      </div>
    </button>
  )
}
