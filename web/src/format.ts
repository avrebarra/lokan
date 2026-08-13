import type { Priority } from './types'

const TAG: Record<Priority, string> = {
  critical: 'crit',
  high: 'high',
  medium: 'med',
  low: 'low',
}

export function priorityTag(p: Priority): string {
  return TAG[p]
}
