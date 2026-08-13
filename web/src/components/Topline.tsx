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
    <header className="topline">
      <div className="wordmark">lokan</div>
      <div className="topline-actions">
        <span className="meta-line">
          {taskCount} tasks · updated {time}
        </span>
        <button
          className="theme-toggle"
          onClick={onToggleTheme}
          aria-label={`switch to ${target} theme`}
        >
          {target}
        </button>
        <button className="button accent" onClick={onCreate}>
          + new task
        </button>
      </div>
    </header>
  )
}
