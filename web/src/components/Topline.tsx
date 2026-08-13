interface Props {
  taskCount: number
  updatedAt: Date
  theme: 'light' | 'dark'
  onToggleTheme: () => void
  onCreate: () => void
}

export default function Topline({ taskCount, updatedAt, theme, onToggleTheme, onCreate }: Props) {
  const time = updatedAt.toTimeString().slice(0, 5)
  const target = theme === 'dark' ? 'light' : 'dark'
  return (
    <header className="mb-4 flex items-center justify-between gap-4 max-[900px]:flex-wrap">
      <div className="font-sans text-[22px] font-bold tracking-[-0.01em]">lokan</div>
      <div className="flex flex-wrap items-center gap-3">
        <span className="text-[11px] uppercase text-muted max-[900px]:hidden">
          {taskCount} tasks · updated {time}
        </span>
        <button
          className="px-0.5 text-[11px] uppercase text-muted transition-colors duration-[120ms] hover:text-fg"
          onClick={onToggleTheme}
          aria-label={`switch to ${target} theme`}
        >
          {target}
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
