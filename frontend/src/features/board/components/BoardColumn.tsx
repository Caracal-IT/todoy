import { useMemo, useState } from 'react'
import type { Column, ColumnDraft, Priority, Task } from '../types'

interface BoardColumnProps {
  activeDropColumnId: number | null
  column: Column
  canMoveLeft: boolean
  canMoveRight: boolean
  draggedTaskId: number | null
  onAddTask: (columnId: number) => void
  onDeleteColumn: (columnId: number) => Promise<void>
  onDeleteTask: (taskId: number) => Promise<void>
  onEditColumn: (columnId: number, draft: ColumnDraft) => Promise<void>
  onEditTask: (task: Task) => void
  onMoveColumn: (columnId: number, direction: 'left' | 'right') => Promise<void>
  onMoveTask: (taskId: number, targetColumnId: number, targetPosition: number) => Promise<void>
  onTaskPointerStart: (task: Task, event: React.PointerEvent<HTMLElement>) => void
}

const priorityTone: Record<Priority, string> = {
  low: 'var(--priority-low)',
  medium: 'var(--priority-medium)',
  high: 'var(--priority-high)',
  critical: 'var(--priority-critical)',
}

export function BoardColumn({
  activeDropColumnId,
  column,
  canMoveLeft,
  canMoveRight,
  draggedTaskId,
  onAddTask,
  onDeleteColumn,
  onDeleteTask,
  onEditColumn,
  onEditTask,
  onMoveColumn,
  onMoveTask,
  onTaskPointerStart,
}: BoardColumnProps) {
  const [isEditing, setIsEditing] = useState(false)
  const [draft, setDraft] = useState<ColumnDraft>({ name: column.name, color: column.color })

  const taskCountLabel = useMemo(() => `${column.tasks.length} cards`, [column.tasks.length])

  async function handleColumnSave(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    await onEditColumn(column.id, draft)
    setIsEditing(false)
  }

  async function handleDeleteColumn() {
    if (window.confirm(`Delete "${column.name}" and every task in it?`)) {
      await onDeleteColumn(column.id)
    }
  }

  return (
    <section
      className={`board-column${activeDropColumnId === column.id ? ' board-column--drop-active' : ''}`}
      data-drop-column-id={column.id}
    >
      <header className="column-header" style={{ ['--column-accent' as string]: column.color }}>
        {isEditing ? (
          <form className="column-edit-form" onSubmit={handleColumnSave}>
            <input
              value={draft.name}
              onChange={(event) => setDraft({ ...draft, name: event.target.value })}
              aria-label="Column name"
            />
            <div className="column-edit-actions">
              <input
                type="color"
                value={draft.color}
                onChange={(event) => setDraft({ ...draft, color: event.target.value })}
                aria-label="Column color"
              />
              <button type="submit" className="primary-button primary-button--small">
                Save
              </button>
              <button
                type="button"
                className="ghost-button ghost-button--small"
                onClick={() => {
                  setDraft({ name: column.name, color: column.color })
                  setIsEditing(false)
                }}
              >
                Cancel
              </button>
            </div>
          </form>
        ) : (
          <>
            <div>
              <div className="column-title-row">
                <span className="column-dot" />
                <h3>{column.name}</h3>
              </div>
              <p>{taskCountLabel}</p>
            </div>
            <div className="column-actions">
              <button
                type="button"
                className="ghost-button ghost-button--small"
                onClick={() => onMoveColumn(column.id, 'left')}
                disabled={!canMoveLeft}
              >
                Left
              </button>
              <button
                type="button"
                className="ghost-button ghost-button--small"
                onClick={() => onMoveColumn(column.id, 'right')}
                disabled={!canMoveRight}
              >
                Right
              </button>
              <button
                type="button"
                className="ghost-button ghost-button--small"
                onClick={() => setIsEditing(true)}
              >
                Edit
              </button>
              <button
                type="button"
                className="ghost-button ghost-button--small ghost-button--danger"
                onClick={handleDeleteColumn}
              >
                Delete
              </button>
            </div>
          </>
        )}
      </header>

      <div className="column-body">
        <button type="button" className="primary-button add-card-button" onClick={() => onAddTask(column.id)}>
          + Add task
        </button>

        <div className="task-list">
          {column.tasks.map((task, taskIndex) => (
            <article
              key={task.id}
              className={`task-card${draggedTaskId === task.id ? ' task-card--dragging' : ''}`}
              onPointerDown={(event) => onTaskPointerStart(task, event)}
            >
              <div className="task-card-top">
                <span className="priority-pill" style={{ backgroundColor: priorityTone[task.priority] }}>
                  {task.priority}
                </span>
                {task.dueDate && <span className="due-pill">{task.dueDate}</span>}
              </div>

              <h4>{task.title}</h4>
              <p>{task.description || 'No extra notes yet.'}</p>

              <div className="task-card-actions">
                <button
                  type="button"
                  className="ghost-button ghost-button--small"
                  onClick={() => onMoveTask(task.id, column.id, Math.max(taskIndex - 1, 0))}
                  disabled={taskIndex === 0}
                >
                  Up
                </button>
                <button
                  type="button"
                  className="ghost-button ghost-button--small"
                  onClick={() =>
                    onMoveTask(task.id, column.id, Math.min(taskIndex + 1, column.tasks.length - 1))
                  }
                  disabled={taskIndex === column.tasks.length - 1}
                >
                  Down
                </button>
                <button
                  type="button"
                  className="ghost-button ghost-button--small"
                  onClick={() => onEditTask(task)}
                >
                  Edit
                </button>
                <button
                  type="button"
                  className="ghost-button ghost-button--small ghost-button--danger"
                  onClick={() => void onDeleteTask(task.id)}
                >
                  Delete
                </button>
              </div>
            </article>
          ))}

          {column.tasks.length === 0 && (
            <div className="empty-column">
              <p>Drop a task here or create a new one.</p>
            </div>
          )}
        </div>
      </div>
    </section>
  )
}
