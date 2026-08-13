import type { Status, Task, TaskSummary, Priority, TaskType } from './types'

// fetch wrapper: json in/out, error body extracted as Error
async function req<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    headers: { 'Content-Type': 'application/json' },
    ...init,
  })
  // surface the server error message instead of a generic failure
  if (!res.ok) {
    const body = await res.json().catch(() => null)
    throw new Error(body?.error ?? `request failed: ${res.status}`)
  }
  return res.json() as Promise<T>
}

export interface FetchTasksResult {
  tasks: TaskSummary[]
  root: string
}

export interface FetchTaskResult {
  task: Task
}

export interface UpdateTaskResult {
  task: Task
}

export async function fetchTasks(): Promise<FetchTasksResult> {
  return req<FetchTasksResult>('/api/tasks')
}

export async function fetchTask(id: string): Promise<FetchTaskResult> {
  return req<FetchTaskResult>(`/api/task/${id}`)
}

export interface UpdateTaskInput {
  id: string
  field: 'status' | 'priority' | 'title' | 'parent' | 'tags' | 'type' | 'body'
  value: string
}

export async function updateTask(
  id: string,
  field: UpdateTaskInput['field'],
  value: string,
): Promise<UpdateTaskResult> {
  return req<UpdateTaskResult>('/api/update', {
    method: 'POST',
    body: JSON.stringify({ id, field, value }),
  })
}

// move a task to a lane and/or position: beforeId anchors the landing spot,
// empty appends at the end of the lane
export async function moveTask(
  id: string,
  status: Status,
  beforeId: string,
): Promise<UpdateTaskResult> {
  return req<UpdateTaskResult>('/api/move', {
    method: 'POST',
    body: JSON.stringify({ id, status, beforeId }),
  })
}

export interface CreateTaskInput {
  title: string
  type: TaskType
  priority: Priority
  parent?: string
}

export async function createTask(input: CreateTaskInput): Promise<UpdateTaskResult> {
  return req<UpdateTaskResult>('/api/create', {
    method: 'POST',
    body: JSON.stringify(input),
  })
}
