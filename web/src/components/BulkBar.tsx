import { useState } from 'react'
import type { Status, StatusDef } from '../lib/types'
import Modal from './Modal'
import ModalConfirm from './ModalConfirm'
import { buttonClass, confirmClass } from '../lib/modal-classes'

interface Props {
  count: number
  lanes: StatusDef[]
  hasArchived: boolean
  onDelete: () => void
  onArchive: () => void
  onMoveTo: (status: Status) => void
  onClear: () => void
}

// sticky bottom bar with bulk actions for the current multi-select
export default function BulkBar({
  count,
  lanes,
  hasArchived,
  onDelete,
  onArchive,
  onMoveTo,
  onClear,
}: Props) {
  const [confirming, setConfirming] = useState(false)
  const [moveOpen, setMoveOpen] = useState(false)

  return (
    <>
      <div className="fixed bottom-4 left-1/2 z-30 -translate-x-1/2">
        <div className="flex items-center gap-3 border border-border bg-bg px-4 py-2.5 shadow-lg">
          <span className="text-[11px] uppercase text-muted">{count} selected</span>
          <span className="h-4 w-px bg-border" aria-hidden />
          <button className={buttonClass} onClick={() => setMoveOpen(true)}>
            move
          </button>
          <button
            className={buttonClass}
            onClick={onArchive}
            disabled={!hasArchived}
            title={hasArchived ? 'archive selected' : 'no archived lane configured'}
          >
            archive
          </button>
          <button className={confirmClass} onClick={() => setConfirming(true)}>
            delete
          </button>
          <button
            className="ml-2 text-[15px] leading-none text-muted hover:text-fg"
            onClick={onClear}
            title="clear selection"
            aria-label="clear selection"
          >
            ×
          </button>
        </div>
      </div>
      {moveOpen && (
        <Modal
          title={`Move ${count} task${count === 1 ? '' : 's'}`}
          onClose={() => setMoveOpen(false)}
          ariaLabel="move selected tasks"
          maxWidth="w-[min(100%,340px)]"
          footer={
            <div className="flex justify-end border-t border-border px-5 pb-5 pt-4">
              <button className={buttonClass} onClick={() => setMoveOpen(false)}>
                cancel
              </button>
            </div>
          }
        >
          <div className="px-5 py-4">
            <p className="mb-3 text-[11px] uppercase text-muted">choose a lane</p>
            <div className="flex flex-col gap-1.5">
              {lanes.map((l) => (
                <button
                  key={l.id}
                  className={`${buttonClass} w-full`}
                  onClick={() => {
                    setMoveOpen(false)
                    onMoveTo(l.id)
                  }}
                >
                  {l.id}
                </button>
              ))}
            </div>
          </div>
        </Modal>
      )}
      {confirming && (
        <ModalConfirm
          title="Delete tasks"
          message={`Delete ${count} selected task${count === 1 ? '' : 's'}? This cannot be undone.`}
          confirmLabel="delete"
          onConfirm={() => {
            setConfirming(false)
            onDelete()
          }}
          onCancel={() => setConfirming(false)}
        />
      )}
    </>
  )
}
