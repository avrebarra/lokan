import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import {
  clearTasks,
  createTask,
  fetchTask,
  fetchTasks,
  moveTask,
  updateStatuses,
  updateTask,
} from './api'
import type { CreateTaskInput } from './api'
import type { Status, StatusDef, Task, TaskSummary, TaskType } from './types'
import Topline from './components/Topline'
import Board from './components/Board'
import DetailModal from './components/DetailModal'
import type { TaskFieldChange } from './components/DetailModal'
import CreateModal from './components/CreateModal'
import ConfigModal from './components/ConfigModal'

export default function App() {
  // board state: tasks, lanes, selection, modals, theme
  const [tasks, setTasks] = useState<TaskSummary[]>([])
  const [statuses, setStatuses] = useState<StatusDef[]>([])
  const [selected, setSelected] = useState<Task | null>(null)
  const [creating, setCreating] = useState<{ parent?: string; type?: TaskType } | null>(null)
  const [configOpen, setConfigOpen] = useState(false)
  const [theme, setTheme] = useState<'light' | 'dark'>(() =>
    document.documentElement.dataset.theme === 'dark' ? 'dark' : 'light',
  )
  const [updatedAt, setUpdatedAt] = useState<Date>(new Date())
  const [movedId, setMovedId] = useState<string | null>(null)
  const movedTimer = useRef<ReturnType<typeof setTimeout> | null>(null)

  // reload board and lanes from the api and stamp the refresh time
  const refresh = useCallback(async () => {
    const { tasks, statuses } = await fetchTasks()
    setTasks(tasks)
    setStatuses(statuses)
    setUpdatedAt(new Date())
  }, [])

  // initial load on mount
  useEffect(() => {
    void refresh()
  }, [refresh])

  // persist theme to the dom dataset + localStorage
  useEffect(() => {
    document.documentElement.dataset.theme = theme
    localStorage.setItem('lokan-theme', theme)
  }, [theme])

  // fetch full detail for the modal
  const openTask = async (id: string) => {
    const { task } = await fetchTask(id)
    setSelected(task)
  }

  const closeDetail = () => setSelected(null)

  // flash the left edge of a task after a successful lane move
  const flashMoved = (id: string) => {
    setMovedId(id)
    if (movedTimer.current) clearTimeout(movedTimer.current)
    movedTimer.current = setTimeout(() => setMovedId(null), 1000)
  }

  // move a task to a lane and/or position: dropping back into its own slot
  // (before itself, after itself, or at its own end) is a no-op
  const handleMove = async (id: string, status: Status, beforeId?: string) => {
    const task = tasks.find((t) => t.id === id)
    if (!task) return
    if (task.status === status) {
      const lane = tasks.filter((t) => t.status === status)
      const from = lane.findIndex((t) => t.id === id)
      if (beforeId === id || lane[from + 1]?.id === beforeId) return
      if (beforeId === '' && from === lane.length - 1) return
    }
    await moveTask(id, status, beforeId ?? '')
    await refresh()
    flashMoved(id)
  }

  // apply edited fields one at a time, close the modal, reload
  const handleSaveChanges = async (changes: TaskFieldChange[]) => {
    if (!selected) return
    for (const change of changes) {
      await updateTask(selected.id, change.field, change.value)
    }
    setSelected(null)
    await refresh()
  }

  // create a task, close the modal, reload
  const handleCreate = async (input: CreateTaskInput) => {
    await createTask(input)
    setCreating(null)
    await refresh()
  }

  // replace the lane set, close the modal, reload (moves + renames handled server-side)
  const handleSaveStatuses = async (next: StatusDef[]) => {
    await updateStatuses(next)
    setConfigOpen(false)
    await refresh()
  }

  // bulk delete: clear archived lanes or the whole board
  const handleClearArchived = async () => {
    await clearTasks('archived')
    await refresh()
  }

  const handleClearAll = async () => {
    await clearTasks('all')
    await refresh()
  }

  // count children per parent for the row badges
  const subtaskCount = useMemo(() => {
    const counts = new Map<string, number>()
    for (const t of tasks) {
      if (t.parent) counts.set(t.parent, (counts.get(t.parent) ?? 0) + 1)
    }
    return counts
  }, [tasks])

  // per-lane task counts for the config modal's delete confirmations
  const laneCounts = useMemo(() => {
    const counts: Record<string, number> = {}
    for (const t of tasks) {
      counts[t.status] = (counts[t.status] ?? 0) + 1
    }
    return counts
  }, [tasks])

  // children of the selected task for the modal list
  const selectedSubtasks = useMemo(
    () => (selected ? tasks.filter((t) => t.parent === selected.id) : []),
    [selected, tasks],
  )

  return (
    <div className="mx-auto w-full max-w-[1200px] px-6 pb-16 pt-8 max-[900px]:px-4 max-[900px]:pb-12 max-[900px]:pt-5">
      <Topline
        taskCount={tasks.filter((t) => t.type !== 'subtask').length}
        updatedAt={updatedAt}
        onCreate={() => setCreating({})}
        onOpenConfig={() => setConfigOpen(true)}
      />
      <Board
        statuses={statuses}
        tasks={tasks}
        subtaskCount={subtaskCount}
        movedId={movedId}
        onSelect={openTask}
        onMove={handleMove}
      />
      {selected && (
        <DetailModal
          task={selected}
          subtasks={selectedSubtasks}
          statuses={statuses}
          onClose={closeDetail}
          onSave={handleSaveChanges}
          onAddSubtask={() => setCreating({ parent: selected.id, type: 'subtask' })}
        />
      )}
      {creating && (
        <CreateModal
          onClose={() => setCreating(null)}
          onCreate={handleCreate}
          initialParent={creating.parent}
          initialType={creating.type}
        />
      )}
      {configOpen && (
        <ConfigModal
          statuses={statuses}
          laneCounts={laneCounts}
          theme={theme}
          onSetTheme={setTheme}
          onClose={() => setConfigOpen(false)}
          onSave={handleSaveStatuses}
          onClearArchived={handleClearArchived}
          onClearAll={handleClearAll}
        />
      )}
    </div>
  )
}
