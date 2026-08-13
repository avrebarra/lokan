export type TaskType = 'epic' | 'task' | 'subtask' | 'bug'

export type Status = 'todo' | 'in-progress' | 'backlog' | 'done' | 'cancelled'

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

export const STATUSES: Status[] = ['todo', 'in-progress', 'backlog', 'done', 'cancelled']

export const PRIORITIES: Priority[] = ['critical', 'high', 'medium', 'low']

export const TASK_TYPES: TaskType[] = ['epic', 'task', 'subtask', 'bug']
