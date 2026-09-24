import { useEffect } from 'react'

// true when the keydown originated inside an editable field — skip hotkeys
export function isTypingInInput(e: KeyboardEvent): boolean {
  const target = e.target as HTMLElement | null
  if (!target) return false
  const tag = target.tagName.toLowerCase()
  if (tag === 'input' || tag === 'textarea' || tag === 'select') return true
  if (target.isContentEditable) return true
  return false
}

// tiny hotkey helper — single key, no ShortcutTree; respects typing guard
export function useHotkey(key: string, handler: () => void, disabled = false) {
  useEffect(() => {
    if (disabled) return
    const onKey = (e: KeyboardEvent) => {
      if (isTypingInInput(e)) return
      if (e.key.toLowerCase() === key.toLowerCase()) {
        if (e.metaKey || e.ctrlKey || e.altKey) return
        e.preventDefault()
        handler()
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [key, handler, disabled])
}
