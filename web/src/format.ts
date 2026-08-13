import type { Priority, Status } from './types'

const TAG: Record<Priority, string> = {
  critical: 'crit',
  high: 'high',
  medium: 'med',
  low: 'low',
}

export function priorityTag(p: Priority): string {
  return TAG[p]
}

export function nextStatus(s: Status): Status {
  switch (s) {
    case 'todo':
      return 'in-progress'
    case 'in-progress':
      return 'done'
    default:
      return 'todo'
  }
}
