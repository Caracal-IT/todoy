import type { Board, Column, ColumnDraft, Task, TaskDraft } from './types'

const API_BASE = '/api/v1'

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers ?? {}),
    },
    ...init,
  })

  if (!response.ok) {
    const data = (await response.json().catch(() => null)) as { error?: string } | null
    throw new Error(data?.error ?? 'Request failed')
  }

  if (response.status === 204) {
    return undefined as T
  }

  return (await response.json()) as T
}

export function fetchBoard(): Promise<Board> {
  return request<Board>('/board')
}

export function createColumn(payload: ColumnDraft): Promise<Column> {
  return request<Column>('/columns', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function updateColumn(columnId: number, payload: ColumnDraft): Promise<Column> {
  return request<Column>(`/columns/${columnId}`, {
    method: 'PATCH',
    body: JSON.stringify(payload),
  })
}

export function deleteColumn(columnId: number): Promise<void> {
  return request<void>(`/columns/${columnId}`, {
    method: 'DELETE',
  })
}

export function reorderColumns(columnIds: number[]): Promise<void> {
  return request<void>('/columns/reorder', {
    method: 'POST',
    body: JSON.stringify({ columnIds }),
  })
}

export function createTask(payload: TaskDraft): Promise<Task> {
  return request<Task>('/tasks', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function updateTask(taskId: number, payload: Omit<TaskDraft, 'columnId'>): Promise<Task> {
  return request<Task>(`/tasks/${taskId}`, {
    method: 'PATCH',
    body: JSON.stringify(payload),
  })
}

export function deleteTask(taskId: number): Promise<void> {
  return request<void>(`/tasks/${taskId}`, {
    method: 'DELETE',
  })
}

export function moveTask(
  taskId: number,
  targetColumnId: number,
  targetPosition: number,
): Promise<Task> {
  return request<Task>(`/tasks/${taskId}/move`, {
    method: 'POST',
    body: JSON.stringify({ targetColumnId, targetPosition }),
  })
}

