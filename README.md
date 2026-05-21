# Todoy

A modern Kanban board built with Go, Fiber, SQLite, and React.

Feature documentation lives in `docs\features`. The current feature spec is `docs\features\kanban-board.md`.

## Summary

- REST API built with Go and Fiber
- SQLite persistence with seeded default columns
- React frontend with a dark neon-purple theme
- Kanban essentials: create, edit, delete, move, and reorder columns and tasks

## User Expectations

- Manage work on a single board without authentication
- Create and edit tasks with priorities, descriptions, and due dates
- Move tasks between columns and reorder workflow stages
- Use a clean modern UI with fast feedback and clear status visibility

## Acceptance Criteria

- The API exposes REST endpoints for board, column, and task operations
- The application persists board state in SQLite
- The frontend supports the normal Kanban flow: add, edit, delete, move, and reorder tasks and columns
- The UI uses a dark theme with neon purple accents and responsive layouts

## Project Status

In progress

## Run locally

### API

```powershell
go run .\cmd\api
```

### Frontend

```powershell
Set-Location .\frontend
npm install
npm run dev
```

The Vite development server proxies `/api` requests to `http://localhost:8080`.
