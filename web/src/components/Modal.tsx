import { useEffect, useRef } from 'react'
import type { ReactNode } from 'react'

interface ModalProps {
  title: ReactNode
  onClose: () => void
  children: ReactNode
  footer?: ReactNode
  headerRight?: ReactNode
  maxWidth: string
  role?: 'dialog' | 'alertdialog'
  ariaLabel: string
  escapeDisabled?: boolean
  z?: string
  headerAlign?: 'center' | 'start'
}

// shared modal shell: overlay + backdrop + escape handler + header/footer slots
export default function Modal({
  title,
  onClose,
  children,
  footer,
  headerRight,
  maxWidth,
  role = 'dialog',
  ariaLabel,
  escapeDisabled = false,
  z = 'z-40',
  headerAlign = 'center',
}: ModalProps) {
  const onCloseRef = useRef(onClose)
  useEffect(() => {
    onCloseRef.current = onClose
  })

  // escape closes, unless disabled (e.g. unsaved edits pending confirmation)
  useEffect(() => {
    if (escapeDisabled) return
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onCloseRef.current()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [escapeDisabled])

  return (
    <div
      className={`fixed inset-0 ${z} flex items-center justify-center p-6`}
      role={role}
      aria-modal="true"
      aria-label={ariaLabel}
    >
      <div
        className="absolute inset-0 cursor-pointer bg-[color-mix(in_srgb,var(--bg)_72%,transparent)]"
        onClick={onClose}
      />
      <div
        className={`relative z-1 flex max-h-[min(92vh,840px)] flex-col border border-fg bg-bg shadow-[0_24px_80px_color-mix(in_srgb,var(--fg)_18%,transparent)] ${maxWidth}`}
      >
        <header
          className={`flex ${headerAlign === 'start' ? 'items-start' : 'items-center'} justify-between gap-6 border-b border-fg px-5 py-[18px]`}
        >
          <div className="min-w-0 flex-1">{title}</div>
          {headerRight}
        </header>
        {children}
        {footer}
      </div>
    </div>
  )
}
