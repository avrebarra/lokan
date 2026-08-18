import { useState } from 'react'
import type { DragEvent } from 'react'
import { Check, Copy } from 'lucide-react'
import type { TaskSummary } from '../types'
import { fetchTask } from '../api'
import { priorityTag } from '../format'

interface Props {
  task: TaskSummary
  subtaskCount: number
  moved: boolean
  onClick: () => void
}

export default function TaskRow({ task, subtaskCount, moved, onClick }: Props) {
  const isSubtask = task.type === 'subtask'
  const [copied, setCopied] = useState(false)

  // start a lane move: carry the task id to the drop target
  const onDragStart = (e: DragEvent<HTMLButtonElement>) => {
    e.dataTransfer.setData('text/x-lokan-task', task.id)
    e.dataTransfer.effectAllowed = 'move'
  }

  // compose a complete, pasteable brief for an AI agent working on this task:
  // id, key fields, the board file line range, and the full body
  const handleCopy = async () => {
    const { task: full } = await fetchTask(task.id)
    const parts = [
      `# ${full.id} — ${full.title}`,
      `Type: ${full.type} | Status: ${full.status} | Priority: ${full.priority}`,
      `Location: .lokan/board.md lines ${full.lineStart}-${full.lineEnd} (task id: ${full.id})`,
    ]
    if (full.tags?.length) parts.push(`Tags: ${full.tags.join(', ')}`)
    parts.push('', '## Body', full.body.trim(), '', '## Work on this')
    parts.push(
      `Complete "${full.title}" (id ${full.id}) — task block is in .lokan/board.md lines ${full.lineStart}-${full.lineEnd}.`,
    )
    await navigator.clipboard.writeText(parts.join('\n'))
    setCopied(true)
    setTimeout(() => setCopied(false), 1500)
  }

  return (
    <div
      className={`group relative w-full select-none bg-bg text-fg transition-colors duration-[120ms] hover:bg-zebra ${
        moved ? 'animate-left-flash' : ''
      }`}
    >
      <button
        draggable
        data-row
        onDragStart={onDragStart}
        className={`block w-full border-b border-border py-[11px] text-left hover:border-l-[3px] hover:border-l-accent ${
          isSubtask ? 'pl-6' : 'pl-2.5'
        }`}
        onClick={onClick}
      >
        <div className="flex items-baseline justify-between gap-4">
          <div className="flex flex-wrap items-baseline gap-2">
            <span
              className={
                isSubtask ? 'text-[10px] uppercase text-muted' : 'text-[11px] uppercase text-muted'
              }
            >
              {task.id}
            </span>
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
        </div>
        <div
          className={`mt-[3px] font-sans font-normal leading-[1.35] group-hover:underline ${
            isSubtask ? 'text-[13px]' : 'text-sm'
          }`}
        >
          {isSubtask ? `— ${task.title}` : task.title}
          {subtaskCount > 0 && (
            <span className="ml-2 text-[11px] text-muted">[{subtaskCount}]</span>
          )}
        </div>
      </button>
      <button
        onClick={() => void handleCopy()}
        className="absolute right-2 top-2 hidden p-1.5 text-muted transition-colors hover:text-fg group-hover:block"
        title="copy a pasteable brief for an AI agent"
      >
        {copied ? <Check size={18} /> : <Copy size={18} />}
      </button>
    </div>
  )
}
