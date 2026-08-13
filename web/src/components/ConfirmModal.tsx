import { useEffect } from 'react'
import type { ReactNode } from 'react'

interface Props {
  title: string
  message: ReactNode
  confirmLabel: string
  onConfirm: () => void
  onCancel: () => void
}

// shared button styles, matching the other modals
const buttonClass =
  'inline-flex min-h-8 items-center justify-center border border-fg bg-bg px-2.5 text-[11px] uppercase text-fg transition-colors duration-[120ms] hover:bg-fg hover:text-bg'
// confirm inverts: solid fg on bg — reads as the primary/destructive action
const confirmClass =
  'inline-flex min-h-8 items-center justify-center border border-fg bg-fg px-2.5 text-[11px] uppercase text-bg transition-colors duration-[120ms] hover:bg-bg hover:text-fg'

export default function ConfirmModal({ title, message, confirmLabel, onConfirm, onCancel }: Props) {
  // escape cancels, same as clicking away
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onCancel()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onCancel])

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-6"
      role="alertdialog"
      aria-modal="true"
      aria-label={title}
    >
      <div
        className="absolute inset-0 cursor-pointer bg-[color-mix(in_srgb,var(--bg)_72%,transparent)]"
        onClick={onCancel}
      />
      <div className="relative z-1 w-[min(100%,420px)] border border-fg bg-bg shadow-[0_24px_80px_color-mix(in_srgb,var(--fg)_18%,transparent)]">
        <div className="border-b border-fg px-5 py-[18px]">
          <h2 className="font-sans text-xl font-normal leading-[1.15]">{title}</h2>
        </div>
        <div className="px-5 py-4 text-[13px] leading-[1.55]">{message}</div>
        <div className="flex justify-end gap-2 border-t border-border px-5 pb-5 pt-4">
          <button className={buttonClass} onClick={onCancel}>
            cancel
          </button>
          <button className={confirmClass} onClick={onConfirm} autoFocus>
            {confirmLabel}
          </button>
        </div>
      </div>
    </div>
  )
}
