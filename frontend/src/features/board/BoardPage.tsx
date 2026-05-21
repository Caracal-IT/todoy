import { useEffect, useMemo, useRef, useState } from 'react'
import './board.css'
import {
  createColumn,
  createTask,
  deleteColumn,
  deleteTask,
  fetchBoard,
  moveTask,
  reorderColumns,
  updateColumn,
  updateTask,
} from './api'
import { BoardColumn } from './components/BoardColumn'
import { StatsBar } from './components/StatsBar'
import { TaskModal } from './components/TaskModal'
import type { Board, ColumnDraft, Priority, Task, TaskDraft } from './types'

const initialColumnDraft: ColumnDraft = {
  name: '',
  color: '#8b5cf6',
}

const priorityTone: Record<Priority, string> = {
  low: 'var(--priority-low)',
  medium: 'var(--priority-medium)',
  high: 'var(--priority-high)',
  critical: 'var(--priority-critical)',
}

interface DragState {
  task: Task
  pointerX: number
  pointerY: number
  targetColumnId: number | null
}

export function BoardPage() {
  const [board, setBoard] = useState<Board | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [savingTask, setSavingTask] = useState(false)
  const [columnDraft, setColumnDraft] = useState<ColumnDraft>(initialColumnDraft)
  const [searchQuery, setSearchQuery] = useState('')
  const [priorityFilter, setPriorityFilter] = useState<'all' | Task['priority']>('all')
  const [modalState, setModalState] = useState<{
    task?: Task
    defaultColumnId?: number
  } | null>(null)
  const [dragState, setDragState] = useState<DragState | null>(null)
  const lockedScrollPositionRef = useRef({ x: 0, y: 0 })

  useEffect(() => {
    void loadBoard()
  }, [])

  useEffect(() => {
    const isDragging = dragState !== null
    const root = document.documentElement
    const body = document.body
    const handlePointerMove = (event: PointerEvent) => {
      event.preventDefault()

      const dropColumn = document
        .elementFromPoint(event.clientX, event.clientY)
        ?.closest<HTMLElement>('[data-drop-column-id]')
      const targetColumnId = dropColumn?.dataset.dropColumnId

      setDragState((current) =>
        current
          ? {
              ...current,
              pointerX: event.clientX,
              pointerY: event.clientY,
              targetColumnId: targetColumnId ? Number(targetColumnId) : null,
            }
          : null,
      )
    }
    const handlePointerUp = () => {
      if (!dragState) {
        return
      }

      const targetColumn = board?.columns.find((column) => column.id === dragState.targetColumnId)
      if (targetColumn && targetColumn.id !== dragState.task.columnId) {
        void handleMoveTask(dragState.task.id, targetColumn.id, targetColumn.tasks.length)
      } else {
        setDragState(null)
      }
    }
    const preventScroll = (event: Event) => {
      event.preventDefault()
      window.scrollTo(
        lockedScrollPositionRef.current.x,
        lockedScrollPositionRef.current.y,
      )
    }
    const preventScrollKeys = (event: KeyboardEvent) => {
      if (
        ['ArrowUp', 'ArrowDown', 'PageUp', 'PageDown', 'Home', 'End', ' ', 'Spacebar'].includes(
          event.key,
        )
      ) {
        event.preventDefault()
      }
    }

    root.classList.toggle('drag-active', isDragging)
    body.classList.toggle('drag-active', isDragging)

    if (isDragging) {
      lockedScrollPositionRef.current = {
        x: window.scrollX,
        y: window.scrollY,
      }

      body.style.position = 'fixed'
      body.style.top = `-${lockedScrollPositionRef.current.y}px`
      body.style.left = `-${lockedScrollPositionRef.current.x}px`
      body.style.right = '0'
      body.style.width = '100%'
      window.addEventListener('pointermove', handlePointerMove, { passive: false })
      window.addEventListener('pointerup', handlePointerUp)
      window.addEventListener('scroll', preventScroll, { passive: false })
      window.addEventListener('wheel', preventScroll, { passive: false })
      window.addEventListener('touchmove', preventScroll, { passive: false })
      window.addEventListener('keydown', preventScrollKeys)
    }

    return () => {
      window.removeEventListener('pointermove', handlePointerMove)
      window.removeEventListener('pointerup', handlePointerUp)
      window.removeEventListener('scroll', preventScroll)
      window.removeEventListener('wheel', preventScroll)
      window.removeEventListener('touchmove', preventScroll)
      window.removeEventListener('keydown', preventScrollKeys)
      root.classList.remove('drag-active')
      body.classList.remove('drag-active')
      body.style.position = ''
      body.style.top = ''
      body.style.left = ''
      body.style.right = ''
      body.style.width = ''

      if (!isDragging) {
        return
      }

      const restoreScrollPosition = () => {
        window.scrollTo(
          lockedScrollPositionRef.current.x,
          lockedScrollPositionRef.current.y,
        )
      }

      restoreScrollPosition()
      window.setTimeout(restoreScrollPosition, 0)
      window.requestAnimationFrame(() => {
        restoreScrollPosition()
        window.requestAnimationFrame(restoreScrollPosition)
      })
    }
  }, [board, dragState])

  const filteredBoard = useMemo(() => {
    if (!board) {
      return null
    }

    const query = searchQuery.trim().toLowerCase()
    return {
      ...board,
      columns: board.columns.map((column) => ({
        ...column,
        tasks: column.tasks.filter((task) => {
          const matchesQuery =
            query === '' ||
            task.title.toLowerCase().includes(query) ||
            task.description.toLowerCase().includes(query)
          const matchesPriority =
            priorityFilter === 'all' || task.priority === priorityFilter

          return matchesQuery && matchesPriority
        }),
      })),
    }
  }, [board, priorityFilter, searchQuery])

  async function loadBoard() {
    try {
      setLoading(true)
      const response = await fetchBoard()
      setBoard(response)
      setError(null)
    } catch (loadError) {
      setError(loadError instanceof Error ? loadError.message : 'Failed to load board')
    } finally {
      setLoading(false)
    }
  }

  async function handleCreateColumn(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()

    try {
      await createColumn(columnDraft)
      setColumnDraft(initialColumnDraft)
      await loadBoard()
      setError(null)
    } catch (createError) {
      setError(createError instanceof Error ? createError.message : 'Failed to create column')
    }
  }

  async function handleEditColumn(columnId: number, draft: ColumnDraft) {
    try {
      await updateColumn(columnId, draft)
      await loadBoard()
      setError(null)
    } catch (updateError) {
      setError(updateError instanceof Error ? updateError.message : 'Failed to update column')
    }
  }

  async function handleDeleteColumn(columnId: number) {
    try {
      await deleteColumn(columnId)
      await loadBoard()
      setError(null)
    } catch (deleteError) {
      setError(deleteError instanceof Error ? deleteError.message : 'Failed to delete column')
    }
  }

  async function handleDeleteTask(taskId: number) {
    try {
      await deleteTask(taskId)
      await loadBoard()
      setError(null)
    } catch (deleteError) {
      setError(deleteError instanceof Error ? deleteError.message : 'Failed to delete task')
    }
  }

  async function handleMoveColumn(columnId: number, direction: 'left' | 'right') {
    if (!board) {
      return
    }

    const currentIndex = board.columns.findIndex((column) => column.id === columnId)
    const targetIndex = direction === 'left' ? currentIndex - 1 : currentIndex + 1
    if (currentIndex === -1 || targetIndex < 0 || targetIndex >= board.columns.length) {
      return
    }

    const reordered = [...board.columns]
    ;[reordered[currentIndex], reordered[targetIndex]] = [reordered[targetIndex], reordered[currentIndex]]

    try {
      await reorderColumns(reordered.map((column) => column.id))
      await loadBoard()
      setError(null)
    } catch (reorderError) {
      setError(reorderError instanceof Error ? reorderError.message : 'Failed to reorder columns')
    }
  }

  async function handleMoveTask(taskId: number, targetColumnId: number, targetPosition: number) {
    try {
      await moveTask(taskId, targetColumnId, targetPosition)
      setDragState(null)
      await loadBoard()
      setError(null)
    } catch (moveError) {
      setError(moveError instanceof Error ? moveError.message : 'Failed to move task')
    }
  }

  function handleTaskPointerStart(task: Task, event: React.PointerEvent<HTMLElement>) {
    if (event.button !== 0) {
      return
    }

    const target = event.target as HTMLElement
    if (target.closest('button, input, textarea, select, a')) {
      return
    }

    event.preventDefault()
    setDragState({
      task,
      pointerX: event.clientX,
      pointerY: event.clientY,
      targetColumnId: task.columnId,
    })
  }

  async function handleTaskSubmit(draft: TaskDraft) {
    try {
      setSavingTask(true)

      if (modalState?.task) {
        await updateTask(modalState.task.id, {
          title: draft.title,
          description: draft.description,
          priority: draft.priority,
          dueDate: draft.dueDate,
        })

        if (modalState.task.columnId !== draft.columnId) {
          await moveTask(modalState.task.id, draft.columnId, Number.MAX_SAFE_INTEGER)
        }
      } else {
        await createTask(draft)
      }

      setModalState(null)
      await loadBoard()
      setError(null)
    } catch (saveError) {
      setError(saveError instanceof Error ? saveError.message : 'Failed to save task')
    } finally {
      setSavingTask(false)
    }
  }

  if (loading) {
    return (
      <main className="board-page">
        <section className="hero-panel">
          <p className="eyebrow">Todoy board</p>
          <h1>Loading your neon workspace...</h1>
        </section>
      </main>
    )
  }

  if (!board || !filteredBoard) {
    return (
      <main className="board-page">
        <section className="hero-panel">
          <p className="eyebrow">Todoy board</p>
          <h1>The board is not available.</h1>
          {error && <p className="error-banner">{error}</p>}
        </section>
      </main>
    )
  }

  return (
    <main className="board-page">
      <section className="hero-panel">
        <div>
          <p className="eyebrow">Project status</p>
          <h1>Kanban flow with a dark neon-purple edge.</h1>
          <p className="hero-copy">
            Track work, move cards across the flow, and keep high-priority tasks visible
            with a clean modern board.
          </p>
        </div>

        <form className="column-composer" onSubmit={handleCreateColumn}>
          <div className="column-composer__title">
            <h2>Add a new workflow stage</h2>
            <p>Create columns for each functional step in your delivery flow.</p>
          </div>

          <div className="column-composer__controls">
            <input
              value={columnDraft.name}
              onChange={(event) => setColumnDraft({ ...columnDraft, name: event.target.value })}
              placeholder="Blocked"
              required
            />
            <input
              type="color"
              value={columnDraft.color}
              onChange={(event) => setColumnDraft({ ...columnDraft, color: event.target.value })}
              aria-label="Column color"
            />
            <button type="submit" className="primary-button">
              Add column
            </button>
          </div>
        </form>
      </section>

      <StatsBar stats={board.stats} />

      <section className="toolbar">
        <label className="toolbar-field">
          Search
          <input
            value={searchQuery}
            onChange={(event) => setSearchQuery(event.target.value)}
            placeholder="Filter by title or description"
          />
        </label>

        <label className="toolbar-field">
          Priority
          <select
            value={priorityFilter}
            onChange={(event) =>
              setPriorityFilter(event.target.value as 'all' | Task['priority'])
            }
          >
            <option value="all">All priorities</option>
            <option value="low">Low</option>
            <option value="medium">Medium</option>
            <option value="high">High</option>
            <option value="critical">Critical</option>
          </select>
        </label>

        <button
          type="button"
          className="ghost-button"
          onClick={() => setModalState({ defaultColumnId: board.columns[0]?.id })}
        >
          + New task
        </button>
      </section>

      {error && <p className="error-banner">{error}</p>}

      <section className="columns-grid">
        {filteredBoard.columns.map((column, index) => (
          <BoardColumn
            activeDropColumnId={dragState?.targetColumnId ?? null}
            key={column.id}
            canMoveLeft={index > 0}
            canMoveRight={index < filteredBoard.columns.length - 1}
            column={column}
            draggedTaskId={dragState?.task.id ?? null}
            onAddTask={(columnId) => setModalState({ defaultColumnId: columnId })}
            onDeleteColumn={handleDeleteColumn}
            onDeleteTask={handleDeleteTask}
            onEditColumn={handleEditColumn}
            onEditTask={(task) => setModalState({ task })}
            onMoveColumn={handleMoveColumn}
            onMoveTask={handleMoveTask}
            onTaskPointerStart={handleTaskPointerStart}
          />
        ))}
      </section>

      {dragState && (
        <div
          className="drag-ghost"
          style={{
            left: `${dragState.pointerX + 18}px`,
            top: `${dragState.pointerY + 18}px`,
          }}
        >
          <article className="task-card task-card--ghost">
            <div className="task-card-top">
              <span
                className="priority-pill"
                style={{ backgroundColor: priorityTone[dragState.task.priority] }}
              >
                {dragState.task.priority}
              </span>
              {dragState.task.dueDate && <span className="due-pill">{dragState.task.dueDate}</span>}
            </div>
            <h4>{dragState.task.title}</h4>
            <p>{dragState.task.description || 'No extra notes yet.'}</p>
          </article>
        </div>
      )}

      {modalState && (
        <TaskModal
          columns={board.columns}
          defaultColumnId={modalState.defaultColumnId}
          initialTask={modalState.task}
          isSubmitting={savingTask}
          onClose={() => setModalState(null)}
          onSubmit={handleTaskSubmit}
        />
      )}
    </main>
  )
}
