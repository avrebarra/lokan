import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import {
  clearTasks,
  createTask,
  deleteTasks,
  fetchTask,
  fetchTasks,
  moveTask,
  updateStatuses,
  updateTask,
} from './lib/api'
import type { CreateTaskInput } from './lib/api'
import type { Status, StatusDef, Task, TaskSummary } from './lib/types'
import Topline from './components/Topline'
import Board from './components/Board'
import BulkBar from './components/BulkBar'
import ModalDetail from './components/ModalDetail'
import type { TaskFieldChange } from './components/ModalDetail'
import ModalCreate from './components/ModalCreate'
import ModalConfig from './components/ModalConfig'

export default function App() {
  // board state: tasks, lanes, selection, modals, theme
  const [tasks, setTasks] = useState<TaskSummary[]>([])
  const [statuses, setStatuses] = useState<StatusDef[]>([])
  const [selected, setSelected] = useState<Task | null>(null)
  const [creating, setCreating] = useState(false)
  const [configOpen, setConfigOpen] = useState(false)
  const [theme, setTheme] = useState<'light' | 'dark'>(() =>
    document.documentElement.dataset.theme === 'dark' ? 'dark' : 'light',
  )
  const [updatedAt, setUpdatedAt] = useState<Date>(new Date())
  const [movedId, setMovedId] = useState<string | null>(null)
  const movedTimer = useRef<ReturnType<typeof setTimeout> | null>(null)
  // multi-select: ids currently selected for bulk actions / group drag
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set())
  const [boardPath, setBoardPath] = useState('')
  const [boardRoot, setBoardRoot] = useState('')

  // reload board and lanes from the api and stamp the refresh time
  const refresh = useCallback(async () => {
    const { tasks, statuses, board_path, board_root } = await fetchTasks()
    setTasks(tasks)
    setStatuses(statuses)
    setBoardPath(board_path)
    setBoardRoot(board_root)
    setUpdatedAt(new Date())
  }, [])

  // show the opened board filename in the browser tab
  useEffect(() => {
    if (!boardPath) return
    const name = boardPath.split('/').pop() ?? boardPath
    document.title = `lokan — ${name}`
  }, [boardPath])

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

  // toggle a task's membership in the multi-select set
  const toggleSelected = (id: string) => {
    setSelectedIds((prev) => {
      const next = new Set(prev)
      if (next.has(id)) {
        next.delete(id)
      } else {
        next.add(id)
      }
      return next
    })
  }

  const clearSelection = () => setSelectedIds(new Set())

  // a finished marquee adds the boxed rows to the existing selection
  const handleMarqueeSelect = (ids: string[]) =>
    setSelectedIds((prev) => new Set([...prev, ...ids]))

  // escape clears the selection
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') clearSelection()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [])

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

  // dispatch a drop: single card keeps the original path (with its no-op
  // guards and flash), a group uses the ordered bulk move
  const handleRowMove = (ids: string[], status: Status, beforeId?: string) => {
    if (ids.length === 1) {
      void handleMove(ids[0], status, beforeId)
    } else {
      void handleMoveGroup(ids, status, beforeId ?? '')
    }
  }

  const handleMoveGroup = async (ids: string[], status: Status, beforeId: string) => {
    await moveMany(ids, status, beforeId)
    flashMoved(ids[0])
  }

  // true when the group already sits where the move would put it — skip the
  // rewrite so `updated` timestamps don't churn for nothing
  const isMoveManyNoop = (ids: string[], status: Status, beforeId: string) => {
    if (beforeId && ids.includes(beforeId)) return false
    const lane = tasks.filter((t) => t.status === status)
    const laneIds = lane.map((t) => t.id)
    const first = laneIds.indexOf(ids[0])
    if (first < 0) return false
    if (beforeId) return laneIds[first + ids.length] === beforeId
    return first + ids.length === laneIds.length
  }

  // move a group of tasks: forward order for append (keeps selection order),
  // reverse for beforeId-insertion (each lands before the anchor)
  const moveMany = async (ids: string[], status: Status, beforeId = '') => {
    if (isMoveManyNoop(ids, status, beforeId)) return
    const order = beforeId ? [...ids].reverse() : ids
    for (const id of order) {
      await moveTask(id, status, beforeId)
    }
    await refresh()
  }

  // bulk actions from the selection bar
  const handleDeleteSelected = async () => {
    const ids = [...selectedIds]
    clearSelection()
    await deleteTasks(ids)
    await refresh()
  }

  const handleArchiveSelected = async () => {
    const target = statuses.find((s) => s.archived)
    if (!target) return
    const ids = [...selectedIds]
    clearSelection()
    await moveMany(ids, target.id)
    await refresh()
  }

  const handleMoveSelected = async (status: Status) => {
    const ids = [...selectedIds]
    clearSelection()
    await moveMany(ids, status)
    await refresh()
  }

  // apply edited fields one at a time, close the modal, reload
  const handleDeleteTask = async () => {
    const id = selected?.id
    setSelected(null)
    if (!id) return
    await deleteTasks([id])
    await refresh()
  }

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
    setCreating(false)
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

  // per-lane task counts for the config modal's delete confirmations
  const laneCounts = useMemo(() => {
    const counts: Record<string, number> = {}
    for (const t of tasks) {
      counts[t.status] = (counts[t.status] ?? 0) + 1
    }
    return counts
  }, [tasks])

  return (
    <div className="mx-auto w-full max-w-[1200px] px-6 pb-16 pt-8 max-[900px]:px-4 max-[900px]:pb-12 max-[900px]:pt-5">
      <Topline
        taskCount={tasks.length}
        updatedAt={updatedAt}
        boardPath={boardPath}
        boardRoot={boardRoot}
        onCreate={() => setCreating(true)}
        onOpenConfig={() => setConfigOpen(true)}
      />
      <Board
        statuses={statuses}
        tasks={tasks}
        movedId={movedId}
        selectedIds={selectedIds}
        onSelect={openTask}
        onMove={handleRowMove}
        onToggleSelect={toggleSelected}
        onMarqueeSelect={handleMarqueeSelect}
      />
      {selectedIds.size > 0 && (
        <BulkBar
          count={selectedIds.size}
          lanes={statuses}
          hasArchived={statuses.some((s) => s.archived)}
          onDelete={() => void handleDeleteSelected()}
          onArchive={() => void handleArchiveSelected()}
          onMoveTo={(status) => void handleMoveSelected(status)}
          onClear={clearSelection}
        />
      )}
      {selected && (
        <ModalDetail
          task={selected}
          statuses={statuses}
          onClose={closeDetail}
          onSave={handleSaveChanges}
          onDelete={() => void handleDeleteTask()}
        />
      )}
      {creating && <ModalCreate onClose={() => setCreating(false)} onCreate={handleCreate} />}
      {configOpen && (
        <ModalConfig
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
