export type Priority = 'low' | 'medium' | 'high' | 'critical'

export interface BoardStats {
  totalTasks: number
  completedTasks: number
  overdueTasks: number
  highPriorityTasks: number
}

export interface Task {
  id: number
  columnId: number
  title: string
  description: string
  priority: Priority
  dueDate: string
  position: number
  createdAt: string
  updatedAt: string
}

export interface Column {
  id: number
  name: string
  color: string
  orderIndex: number
  tasks: Task[]
}

export interface Board {
  columns: Column[]
  stats: BoardStats
  now: string
}

export interface TaskDraft {
  columnId: number
  title: string
  description: string
  priority: Priority
  dueDate: string
}

export interface ColumnDraft {
  name: string
  color: string
}

