import type { BoardStats } from '../types'

interface StatsBarProps {
  stats: BoardStats
}

export function StatsBar({ stats }: StatsBarProps) {
  const cards = [
    { label: 'Total cards', value: stats.totalTasks },
    { label: 'Completed', value: stats.completedTasks },
    { label: 'High priority', value: stats.highPriorityTasks },
    { label: 'Overdue', value: stats.overdueTasks },
  ]

  return (
    <div className="stats-grid">
      {cards.map((card) => (
        <article key={card.label} className="stat-card">
          <p>{card.label}</p>
          <strong>{card.value}</strong>
        </article>
      ))}
    </div>
  )
}

