import Modal from './Modal'
import { buttonClass } from '../lib/modal-classes'

interface Props {
  onClose: () => void
}

// help modal listing available keyboard shortcuts
export default function ModalHelp({ onClose }: Props) {
  return (
    <Modal
      title={<h2 className="font-sans text-xl font-normal leading-[1.15]">keyboard shortcuts</h2>}
      onClose={onClose}
      ariaLabel="keyboard shortcuts help"
      maxWidth="w-[min(100%,440px)]"
      headerRight={
        <button className={`${buttonClass} flex-none`} onClick={onClose}>
          × close
        </button>
      }
    >
      <div className="overflow-y-auto p-5">
        <table className="w-full text-left text-sm">
          <thead>
            <tr className="border-b border-border text-[11px] uppercase text-muted">
              <th className="pb-2 font-normal">shortcut</th>
              <th className="pb-2 font-normal">action</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            <tr>
              <td className="py-2">
                <kbd className="rounded border border-border bg-zebra px-1.5 py-0.5 font-mono text-xs">
                  C
                </kbd>
              </td>
              <td className="py-2 text-muted">Create new task</td>
            </tr>
            <tr>
              <td className="py-2">
                <kbd className="rounded border border-border bg-zebra px-1.5 py-0.5 font-mono text-xs">
                  ?
                </kbd>
              </td>
              <td className="py-2 text-muted">Open this help</td>
            </tr>
            <tr>
              <td className="py-2">
                <kbd className="rounded border border-border bg-zebra px-1.5 py-0.5 font-mono text-xs">
                  Esc
                </kbd>
              </td>
              <td className="py-2 text-muted">Clear selection / close modal</td>
            </tr>
          </tbody>
        </table>
        <p className="mt-4 text-xs leading-relaxed text-muted">
          Shortcuts are disabled while typing in inputs. Press{' '}
          <kbd className="rounded border border-border bg-zebra px-1 py-0.5 font-mono text-[11px]">
            ?
          </kbd>{' '}
          again or{' '}
          <kbd className="rounded border border-border bg-zebra px-1 py-0.5 font-mono text-[11px]">
            Esc
          </kbd>{' '}
          to close.
        </p>
      </div>
    </Modal>
  )
}
