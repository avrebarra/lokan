import { useCallback, useEffect, useMemo, useState } from 'react'
import { createTask, fetchTask, fetchTasks, updateTask } from './api'
import type { CreateTaskInput } from './api'
import type { Task, TaskSummary } from './types'
import { nextStatus } from './format'
import Topline from './components/Topline'
import Board from './components/Board'
import DetailModal from './components/DetailModal'
import type { TaskFieldChange } from './components/DetailModal'
import CreateModal from './components/CreateModal'

export default function App() {
  // board state: tasks, selection, modals, theme
  const [tasks, setTasks] = useState<TaskSummary[]>([])
  const [selected, setSelected] = useState<Task | null>(null)
  const [creating, setCreating] = useState(false)
  const [theme, setTheme] = useState<'light' | 'dark'>(() =>
    document.documentElement.dataset.theme === 'dark' ? 'dark' : 'light',
  )
  const [updatedAt, setUpdatedAt] = useState<Date>(new Date())

  // reload board from the api and stamp the refresh time
  const refresh = useCallback(async () => {
    const { tasks } = await fetchTasks()
    setTasks(tasks)
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

  const toggleTheme = () => setTheme((t) => (t === 'dark' ? 'light' : 'dark'))

  // fetch full detail for the modal
  const openTask = async (id: string) => {
    const { task } = await fetchTask(id)
    setSelected(task)
  }

  const closeDetail = () => setSelected(null)

  // advance selected task's status, close the modal, reload
  const handleAdvance = async () => {
    if (!selected) return
    await updateTask(selected.id, 'status', nextStatus(selected.status))
    setSelected(null)
    await refresh()
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
    setCreating(false)
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
        theme={theme}
        onToggleTheme={toggleTheme}
        onCreate={() => setCreating(true)}
      />
      <Board tasks={tasks} subtaskCount={subtaskCount} onSelect={openTask} />
      {selected && (
        <DetailModal
          task={selected}
          subtasks={selectedSubtasks}
          onClose={closeDetail}
          onAdvance={handleAdvance}
          onSave={handleSaveChanges}
        />
      )}
      {creating && <CreateModal onClose={() => setCreating(false)} onCreate={handleCreate} />}
    </div>
  )
}
