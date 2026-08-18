export type TaskType = 'task'

// statuses are configurable lanes resolved from the server at runtime
export type Status = string

export interface StatusDef {
  id: string
  archived: boolean
}

export interface TaskFrontmatter {
  id: string
  title: string
  status: Status
  created: string
  updated: string
  related?: string[]
  docs?: string[]
  tags?: string[]
}

export interface TaskSummary extends TaskFrontmatter {
  filePath: string
  lineStart: number
  lineEnd: number
}

export interface Task extends TaskSummary {
  body: string
}
