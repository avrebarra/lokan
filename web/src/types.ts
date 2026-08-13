export type TaskType = 'epic' | 'task' | 'subtask' | 'bug'

// statuses are configurable lanes resolved from the server at runtime
export type Status = string

export interface StatusDef {
  id: string
  archived: boolean
}

export type Priority = 'critical' | 'high' | 'medium' | 'low'

export interface TaskFrontmatter {
  id: string
  title: string
  type: TaskType
  status: Status
  priority: Priority
  parent?: string
  related?: string[]
  docs?: string[]
  tags?: string[]
  created: string
  updated: string
}

export interface TaskSummary extends TaskFrontmatter {
  filePath: string
}

export interface Task extends TaskSummary {
  body: string
}

export const PRIORITIES: Priority[] = ['critical', 'high', 'medium', 'low']

export const TASK_TYPES: TaskType[] = ['epic', 'task', 'subtask', 'bug']
