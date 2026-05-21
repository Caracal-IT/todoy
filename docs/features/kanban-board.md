# Kanban Board

## Project Status

Complete

## Summary

The Kanban Board feature provides a modern task-management workspace built with a Go Fiber REST API, SQLite persistence, and a React frontend. It supports the core Kanban workflow for planning, tracking, and completing work without authentication.

## User Expectations

- [x] Users can manage work on a single shared board without logging in.
- [x] Users can create, edit, delete, move, and reorder tasks across workflow columns.
- [x] Users can create, edit, delete, and reorder columns to match their process.
- [x] Users can set task descriptions, priorities, and due dates.
- [x] Users can quickly understand board progress through summary statistics.
- [x] Users can work in a modern dark UI with neon-purple styling and responsive layouts.
- [x] Users can drag and drop tasks without the page scrolling during the interaction.
- [x] Users can clearly distinguish cards through stronger contrast, glow, and elevated styling.

## Acceptance Criteria

- [x] The backend exposes REST endpoints for board, column, and task operations.
- [x] The application persists board state in SQLite.
- [x] Default columns are created automatically for a new board.
- [x] Users can create, edit, delete, and reorder columns.
- [x] Users can create, edit, delete, move, and reorder tasks.
- [x] Users can filter tasks and search by task content.
- [x] The UI uses a dark theme with neon-purple accents.
- [x] The feature works without authentication.
- [x] Dragging a task does not trigger page scrolling while drag-and-drop is active.
- [x] Task cards have clear visual contrast and stand out from the board background.

## Implementation Checklist

- [x] Go Fiber REST API for board, column, and task operations
- [x] SQLite persistence with seeded default columns
- [x] Create, edit, delete, move, and reorder columns
- [x] Create, edit, delete, move, and reorder tasks
- [x] Search and priority filtering
- [x] Dark neon-purple React UI
- [x] No-auth single-board workflow
- [x] Custom pointer-driven drag-and-drop that avoids browser auto-scroll behavior
- [x] Enhanced card styling with stronger borders, shadow, glow, and hover states
- [x] Backend tests and benchmark coverage

## Functional Scope

### Board management

- Load the full board with columns, tasks, and summary statistics.
- Seed standard starting columns for first use.
- Show totals for overall tasks, completed tasks, overdue tasks, and high-priority tasks.

### Column management

- Create new columns.
- Rename existing columns.
- Change column colors.
- Delete columns and their tasks.
- Reorder columns left and right.

### Task management

- Create new tasks in any column.
- Edit task title, description, priority, due date, and target column.
- Delete tasks.
- Move tasks across columns.
- Reorder tasks within a column.
- Support pointer-driven drag-and-drop movement across columns and manual movement controls.
- Prevent page scrolling while task drag-and-drop is active.

### Filtering and usability

- Search tasks by title or description.
- Filter tasks by priority.
- Display clear empty states and actionable controls.
- Render cards with stronger contrast, glow, hover lift, and a clearer dragged state.

## Technical Notes

- Backend: Go, Fiber, SQLite
- Frontend: React, TypeScript, Vite
- API base path: `/api/v1`
- Persistence file: `todoy.db`
- UX refinements: custom pointer-driven drag interactions, no page scrolling during task movement, and higher-contrast task card presentation
