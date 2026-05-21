import { useEffect, useState } from 'react'
import type { Column, Priority, Task, TaskDraft } from '../types'

interface TaskModalProps {
  columns: Column[]
  initialTask?: Task
  defaultColumnId?: number
  isSubmitting: boolean
  onClose: () => void
  onSubmit: (draft: TaskDraft) => Promise<void>
}

const priorities: Priority[] = ['low', 'medium', 'high', 'critical']

export function TaskModal({
  columns,
  initialTask,
  defaultColumnId,
  isSubmitting,
  onClose,
  onSubmit,
}: TaskModalProps) {
  const [draft, setDraft] = useState<TaskDraft>({
    columnId: defaultColumnId ?? columns[0]?.id ?? 0,
    title: initialTask?.title ?? '',
    description: initialTask?.description ?? '',
    priority: initialTask?.priority ?? 'medium',
    dueDate: initialTask?.dueDate ?? '',
  })

  useEffect(() => {
    setDraft({
      columnId: initialTask?.columnId ?? defaultColumnId ?? columns[0]?.id ?? 0,
      title: initialTask?.title ?? '',
      description: initialTask?.description ?? '',
      priority: initialTask?.priority ?? 'medium',
      dueDate: initialTask?.dueDate ?? '',
    })
  }, [columns, defaultColumnId, initialTask])

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    await onSubmit(draft)
  }

  return (
    <div className="modal-backdrop" role="presentation" onClick={onClose}>
      <div
        className="modal-card"
        role="dialog"
        aria-modal="true"
        aria-labelledby="task-modal-title"
        onClick={(event) => event.stopPropagation()}
      >
        <div className="modal-header">
          <div>
            <p className="eyebrow">{initialTask ? 'Task editing' : 'New task'}</p>
            <h2 id="task-modal-title">
              {initialTask ? 'Refine your card' : 'Create a new card'}
            </h2>
          </div>
          <button type="button" className="ghost-button" onClick={onClose}>
            Close
          </button>
        </div>

        <form className="task-form" onSubmit={handleSubmit}>
          <label>
            Title
            <input
              value={draft.title}
              onChange={(event) => setDraft({ ...draft, title: event.target.value })}
              placeholder="Ship board analytics"
              required
            />
          </label>

          <label>
            Description
            <textarea
              rows={4}
              value={draft.description}
              onChange={(event) => setDraft({ ...draft, description: event.target.value })}
              placeholder="Add the important implementation notes here."
            />
          </label>

          <div className="task-form-grid">
            <label>
              Column
              <select
                value={draft.columnId}
                onChange={(event) =>
                  setDraft({ ...draft, columnId: Number(event.target.value) })
                }
              >
                {columns.map((column) => (
                  <option key={column.id} value={column.id}>
                    {column.name}
                  </option>
                ))}
              </select>
            </label>

            <label>
              Priority
              <select
                value={draft.priority}
                onChange={(event) =>
                  setDraft({ ...draft, priority: event.target.value as Priority })
                }
              >
                {priorities.map((priority) => (
                  <option key={priority} value={priority}>
                    {priority}
                  </option>
                ))}
              </select>
            </label>
          </div>

          <label>
            Due date
            <input
              type="date"
              value={draft.dueDate}
              onChange={(event) => setDraft({ ...draft, dueDate: event.target.value })}
            />
          </label>

          <div className="modal-actions">
            <button type="button" className="ghost-button" onClick={onClose}>
              Cancel
            </button>
            <button type="submit" className="primary-button" disabled={isSubmitting}>
              {isSubmitting ? 'Saving...' : initialTask ? 'Save changes' : 'Create task'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

