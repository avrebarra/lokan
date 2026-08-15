import type { ReactNode } from 'react'
import Modal from './Modal'
import { buttonClass, confirmClass } from './modal-classes'

interface Props {
  title: string
  message: ReactNode
  confirmLabel: string
  onConfirm: () => void
  onCancel: () => void
}

export default function ModalConfirm({ title, message, confirmLabel, onConfirm, onCancel }: Props) {
  return (
    <Modal
      title={title}
      onClose={onCancel}
      ariaLabel={title}
      maxWidth="w-[min(100%,420px)]"
      role="alertdialog"
      z="z-50"
      footer={
        <div className="flex justify-end gap-2 border-t border-border px-5 pb-5 pt-4">
          <button className={buttonClass} onClick={onCancel}>
            cancel
          </button>
          <button className={confirmClass} onClick={onConfirm} autoFocus>
            {confirmLabel}
          </button>
        </div>
      }
    >
      <div className="px-5 py-4 text-[13px] leading-[1.55]">{message}</div>
    </Modal>
  )
}
