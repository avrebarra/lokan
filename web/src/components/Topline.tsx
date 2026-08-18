interface Props {
  taskCount: number
  updatedAt: Date
  boardPath: string
  boardRoot: string
  onCreate: () => void
  onOpenConfig: () => void
}

export default function Topline({ taskCount, updatedAt, boardPath, boardRoot, onCreate, onOpenConfig }: Props) {
  const time = updatedAt.toTimeString().slice(0, 5)
  return (
    <header className="mb-4 flex items-center justify-between gap-4 max-[900px]:flex-wrap">
      <div className="min-w-0">
        <div className="font-sans text-[22px] font-bold tracking-[-0.01em]">lokan</div>
        {boardPath && (
          <div className="mt-0.5 truncate text-[11px] text-muted" title={boardPath}>
            {boardPath}
            {boardRoot ? ` (${boardRoot}/)` : ''}
          </div>
        )}
      </div>
      <div className="flex flex-wrap items-center gap-3">
        <span className="text-[11px] uppercase text-muted max-[900px]:hidden">
          {taskCount} tasks · updated {time}
        </span>
        <button
          className="inline-flex min-h-8 items-center justify-center border border-border px-2.5 text-[11px] uppercase text-muted transition-colors duration-[120ms] hover:border-fg hover:text-fg"
          onClick={onOpenConfig}
          aria-label="open board config"
        >
          config
        </button>
        <button
          className="inline-flex min-h-8 items-center justify-center border border-accent bg-accent px-2.5 text-[11px] font-medium uppercase text-black transition-colors duration-[120ms] hover:border-fg hover:bg-fg hover:text-bg"
          onClick={onCreate}
        >
          + new task
        </button>
      </div>
    </header>
  )
}
