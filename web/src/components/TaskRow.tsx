import type { DragEvent } from 'react'
import type { TaskSummary } from '../types'
import { priorityTag } from '../format'

interface Props {
  task: TaskSummary
  subtaskCount: number
  onClick: () => void
}

export default function TaskRow({ task, subtaskCount, onClick }: Props) {
  const isSubtask = task.type === 'subtask'

  // start a lane move: carry the task id to the drop target
  const onDragStart = (e: DragEvent<HTMLButtonElement>) => {
    e.dataTransfer.setData('text/x-lokan-task', task.id)
    e.dataTransfer.effectAllowed = 'move'
  }

  return (
    <button
      draggable
      onDragStart={onDragStart}
      className={`group block w-full select-none border-b border-border bg-bg py-[11px] text-left text-fg transition-colors duration-[120ms] hover:bg-zebra hover:border-l-[3px] hover:border-l-accent ${
        isSubtask ? 'pl-6' : 'pl-2.5'
      }`}
      onClick={onClick}
    >
      <div className="flex items-baseline justify-between gap-4">
        <div className="flex flex-wrap items-baseline gap-2">
          <span className={isSubtask ? 'text-[10px] uppercase text-muted' : 'text-[11px] uppercase text-muted'}>{task.id}</span>
          <span
            className={
              !isSubtask && task.priority === 'critical'
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
      <div
        className={`mt-[3px] font-sans font-normal leading-[1.35] group-hover:underline ${
          isSubtask ? 'text-[13px]' : 'text-sm'
        }`}
      >
        {isSubtask ? `— ${task.title}` : task.title}
      </div>
    </button>
  )
}
