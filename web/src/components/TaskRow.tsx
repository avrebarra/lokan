import { useRef, useState } from 'react'
import type { DragEvent } from 'react'
import { Check, Copy } from 'lucide-react'
import type { TaskSummary } from '../lib/types'
import { fetchTask } from '../lib/api'
import { priorityTag } from '../lib/format'

interface Props {
  task: TaskSummary
  subtaskCount: number
  moved: boolean
  selectedIds: Set<string>
  scopedIds: Set<string>
  onClick: () => void
  onToggleSelect: () => void
}

export default function TaskRow({
  task,
  subtaskCount,
  moved,
  selectedIds,
  scopedIds,
  onClick,
  onToggleSelect,
}: Props) {
  const isSubtask = task.type === 'subtask'
  const [copied, setCopied] = useState(false)
  const lastClick = useRef(0)
  const selected = selectedIds.has(task.id)
  const scoped = scopedIds.has(task.id)
  const selectionActive = selectedIds.size > 0
  // live marquee preview uses the hover tint; committed selection the accent
  const rowBg = selected ? 'bg-accent/10' : scoped ? 'bg-zebra' : ''

  // start a lane move: dragging a selected card carries the whole selection
  // (JSON ids), otherwise just this task
  const onDragStart = (e: DragEvent<HTMLButtonElement>) => {
    const ids = selectionActive && selected ? [...selectedIds] : [task.id]
    e.dataTransfer.setData('text/x-lokan-task', JSON.stringify(ids))
    e.dataTransfer.effectAllowed = 'move'
  }

  // with a selection or live marquee active, a row click toggles membership
  // and a second click within 300ms opens the detail; otherwise it opens
  const handleClick = () => {
    if (!selectionActive) {
      onClick()
      return
    }
    const now = Date.now()
    if (lastClick.current && now - lastClick.current < 300) {
      lastClick.current = 0
      onClick()
      return
    }
    lastClick.current = now
    onToggleSelect()
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
      } ${rowBg}`}
    >
      {(selectionActive) && (
        <input
          type="checkbox"
          checked={selected || scoped}
          onChange={onToggleSelect}
          aria-label={`select ${task.title}`}
          className="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 cursor-pointer accent-accent"
        />
      )}
      <button
        draggable
        data-row
        data-row-id={task.id}
        onDragStart={onDragStart}
        onClick={handleClick}
        className={`block w-full border-b border-border py-[11px] text-left hover:border-l-[3px] hover:border-l-accent ${
          isSubtask ? 'pl-6' : 'pl-2.5'
        } ${selectionActive ? 'pl-9' : ''}`}
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
