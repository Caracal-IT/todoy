package kanban

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Repository manages SQLite persistence for the Kanban board.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// EnsureSchema creates the required database schema.
func (r *Repository) EnsureSchema(ctx context.Context) error {
	schema := `
	CREATE TABLE IF NOT EXISTS columns (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		color TEXT NOT NULL,
		order_index INTEGER NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		column_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		priority TEXT NOT NULL,
		due_date TEXT NULL,
		position INTEGER NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		FOREIGN KEY(column_id) REFERENCES columns(id) ON DELETE CASCADE
	);`

	if _, err := r.db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	return nil
}

// SeedDefaults inserts the standard workflow columns for a new board.
func (r *Repository) SeedDefaults(ctx context.Context) error {
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM columns`).Scan(&count); err != nil {
		return fmt.Errorf("count columns: %w", err)
	}

	if count > 0 {
		return nil
	}

	defaultColumns := []CreateColumnInput{
		{Name: "Backlog", Color: "#7c3aed"},
		{Name: "In Progress", Color: "#9333ea"},
		{Name: "Review", Color: "#a855f7"},
		{Name: "Done", Color: "#6d28d9"},
	}

	for _, column := range defaultColumns {
		if _, err := r.CreateColumn(ctx, column); err != nil {
			return fmt.Errorf("seed column %s: %w", column.Name, err)
		}
	}

	return nil
}

// GetBoard returns all columns and tasks sorted for UI rendering.
func (r *Repository) GetBoard(ctx context.Context) (Board, error) {
	columns, err := r.listColumns(ctx)
	if err != nil {
		return Board{}, err
	}

	columnIndex := make(map[int64]int, len(columns))
	for index, column := range columns {
		columns[index].Tasks = []Task{}
		columnIndex[column.ID] = index
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, column_id, title, description, priority, COALESCE(due_date, ''), position, created_at, updated_at
		FROM tasks
		ORDER BY column_id, position, id`)
	if err != nil {
		return Board{}, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	stats := BoardStats{}
	today := time.Now().UTC().Truncate(24 * time.Hour)

	for rows.Next() {
		task, scanErr := scanTask(rows)
		if scanErr != nil {
			return Board{}, scanErr
		}

		stats.TotalTasks++
		if task.Priority == "high" || task.Priority == "critical" {
			stats.HighPriorityTasks++
		}

		isDoneColumn := false
		columnPosition, exists := columnIndex[task.ColumnID]
		if exists {
			isDoneColumn = strings.EqualFold(columns[columnPosition].Name, "Done")
			if dueDate, ok := task.parsedDueDate(); ok && dueDate.Before(today) && !isDoneColumn {
				stats.OverdueTasks++
			}

			if isDoneColumn {
				stats.CompletedTasks++
			}

			columns[columnPosition].Tasks = append(columns[columnPosition].Tasks, task)
		}
	}

	if err := rows.Err(); err != nil {
		return Board{}, fmt.Errorf("iterate tasks: %w", err)
	}

	return Board{
		Columns: columns,
		Stats:   stats,
		Now:     time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// CreateColumn inserts a new column at the end of the board.
func (r *Repository) CreateColumn(ctx context.Context, input CreateColumnInput) (Column, error) {
	var column Column

	err := r.withTx(ctx, func(tx *sql.Tx) error {
		var orderIndex int
		if err := tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(order_index) + 1, 0) FROM columns`).Scan(&orderIndex); err != nil {
			return fmt.Errorf("get next column order: %w", err)
		}

		now := nowRFC3339()
		result, err := tx.ExecContext(ctx, `
			INSERT INTO columns (name, color, order_index, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?)`,
			input.Name, input.Color, orderIndex, now, now)
		if err != nil {
			return fmt.Errorf("insert column: %w", err)
		}

		columnID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("read new column id: %w", err)
		}

		column = Column{
			ID:         columnID,
			Name:       input.Name,
			Color:      input.Color,
			OrderIndex: orderIndex,
			Tasks:      []Task{},
		}

		return nil
	})
	if err != nil {
		return Column{}, err
	}

	return column, nil
}

// UpdateColumn updates a column.
func (r *Repository) UpdateColumn(ctx context.Context, id int64, input UpdateColumnInput) (Column, error) {
	var column Column

	err := r.withTx(ctx, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, `
			UPDATE columns
			SET name = ?, color = ?, updated_at = ?
			WHERE id = ?`,
			input.Name, input.Color, nowRFC3339(), id)
		if err != nil {
			return fmt.Errorf("update column: %w", err)
		}

		affectedRows, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("read affected column rows: %w", err)
		}

		if affectedRows == 0 {
			return ErrNotFound
		}

		column, err = r.getColumnTx(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return Column{}, err
	}

	return column, nil
}

// DeleteColumn removes a column and its tasks, then normalizes board order.
func (r *Repository) DeleteColumn(ctx context.Context, id int64) error {
	return r.withTx(ctx, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, `DELETE FROM columns WHERE id = ?`, id)
		if err != nil {
			return fmt.Errorf("delete column: %w", err)
		}

		affectedRows, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("read affected deleted columns: %w", err)
		}

		if affectedRows == 0 {
			return ErrNotFound
		}

		return r.resequenceColumns(ctx, tx)
	})
}

// ReorderColumns updates the board column order based on a provided list.
func (r *Repository) ReorderColumns(ctx context.Context, input ReorderColumnsInput) error {
	return r.withTx(ctx, func(tx *sql.Tx) error {
		columnIDs, err := r.listColumnIDsTx(ctx, tx)
		if err != nil {
			return err
		}

		if len(columnIDs) != len(input.ColumnIDs) {
			return ValidationError{Message: "columnIds must include every existing column exactly once"}
		}

		existing := make(map[int64]struct{}, len(columnIDs))
		for _, id := range columnIDs {
			existing[id] = struct{}{}
		}

		for index, id := range input.ColumnIDs {
			if _, ok := existing[id]; !ok {
				return ValidationError{Message: "columnIds contains an unknown column id"}
			}

			if _, err := tx.ExecContext(ctx, `
				UPDATE columns
				SET order_index = ?, updated_at = ?
				WHERE id = ?`,
				index, nowRFC3339(), id); err != nil {
				return fmt.Errorf("update column order: %w", err)
			}
		}

		return nil
	})
}

// CreateTask inserts a new task at the end of a column.
func (r *Repository) CreateTask(ctx context.Context, input CreateTaskInput) (Task, error) {
	var task Task

	err := r.withTx(ctx, func(tx *sql.Tx) error {
		if err := r.requireColumnTx(ctx, tx, input.ColumnID); err != nil {
			return err
		}

		var nextPosition int
		if err := tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(position) + 1, 0) FROM tasks WHERE column_id = ?`, input.ColumnID).Scan(&nextPosition); err != nil {
			return fmt.Errorf("get next task position: %w", err)
		}

		now := nowRFC3339()
		result, err := tx.ExecContext(ctx, `
			INSERT INTO tasks (column_id, title, description, priority, due_date, position, created_at, updated_at)
			VALUES (?, ?, ?, ?, NULLIF(?, ''), ?, ?, ?)`,
			input.ColumnID, input.Title, input.Description, input.Priority, input.DueDate, nextPosition, now, now)
		if err != nil {
			return fmt.Errorf("insert task: %w", err)
		}

		taskID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("read new task id: %w", err)
		}

		task, err = r.getTaskTx(ctx, tx, taskID)
		return err
	})
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

// UpdateTask edits a task.
func (r *Repository) UpdateTask(ctx context.Context, id int64, input UpdateTaskInput) (Task, error) {
	var task Task

	err := r.withTx(ctx, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, `
			UPDATE tasks
			SET title = ?, description = ?, priority = ?, due_date = NULLIF(?, ''), updated_at = ?
			WHERE id = ?`,
			input.Title, input.Description, input.Priority, input.DueDate, nowRFC3339(), id)
		if err != nil {
			return fmt.Errorf("update task: %w", err)
		}

		affectedRows, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("read affected task rows: %w", err)
		}

		if affectedRows == 0 {
			return ErrNotFound
		}

		task, err = r.getTaskTx(ctx, tx, id)
		return err
	})
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

// DeleteTask removes a task and resequences its column.
func (r *Repository) DeleteTask(ctx context.Context, id int64) error {
	return r.withTx(ctx, func(tx *sql.Tx) error {
		task, err := r.getTaskTx(ctx, tx, id)
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, `DELETE FROM tasks WHERE id = ?`, id); err != nil {
			return fmt.Errorf("delete task: %w", err)
		}

		return r.resequenceTasksInColumn(ctx, tx, task.ColumnID)
	})
}

// MoveTask repositions a task within a column or across columns.
func (r *Repository) MoveTask(ctx context.Context, id int64, input MoveTaskInput) (Task, error) {
	var task Task

	err := r.withTx(ctx, func(tx *sql.Tx) error {
		existingTask, err := r.getTaskTx(ctx, tx, id)
		if err != nil {
			return err
		}

		if err := r.requireColumnTx(ctx, tx, input.TargetColumnID); err != nil {
			return err
		}

		sourceIDs, err := r.listTaskIDsTx(ctx, tx, existingTask.ColumnID)
		if err != nil {
			return err
		}

		targetIDs := sourceIDs
		if existingTask.ColumnID != input.TargetColumnID {
			targetIDs, err = r.listTaskIDsTx(ctx, tx, input.TargetColumnID)
			if err != nil {
				return err
			}
		}

		sourceIDs = removeID(sourceIDs, id)
		if existingTask.ColumnID == input.TargetColumnID {
			targetIDs = sourceIDs
		}

		targetPosition := input.TargetPosition
		if targetPosition > len(targetIDs) {
			targetPosition = len(targetIDs)
		}

		targetIDs = insertID(targetIDs, id, targetPosition)
		now := nowRFC3339()

		if existingTask.ColumnID == input.TargetColumnID {
			if err := r.applyTaskOrder(ctx, tx, input.TargetColumnID, targetIDs, now); err != nil {
				return err
			}
		} else {
			if err := r.applyTaskOrder(ctx, tx, existingTask.ColumnID, sourceIDs, now); err != nil {
				return err
			}

			if err := r.applyTaskOrder(ctx, tx, input.TargetColumnID, targetIDs, now); err != nil {
				return err
			}
		}

		task, err = r.getTaskTx(ctx, tx, id)
		return err
	})
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (r *Repository) withTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) listColumns(ctx context.Context) ([]Column, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, color, order_index FROM columns ORDER BY order_index, id`)
	if err != nil {
		return nil, fmt.Errorf("query columns: %w", err)
	}
	defer rows.Close()

	var columns []Column
	for rows.Next() {
		var column Column
		if err := rows.Scan(&column.ID, &column.Name, &column.Color, &column.OrderIndex); err != nil {
			return nil, fmt.Errorf("scan column: %w", err)
		}

		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate columns: %w", err)
	}

	return columns, nil
}

func (r *Repository) getColumnTx(ctx context.Context, tx *sql.Tx, id int64) (Column, error) {
	var column Column
	err := tx.QueryRowContext(ctx, `
		SELECT id, name, color, order_index
		FROM columns
		WHERE id = ?`, id).Scan(&column.ID, &column.Name, &column.Color, &column.OrderIndex)
	if errors.Is(err, sql.ErrNoRows) {
		return Column{}, ErrNotFound
	}

	if err != nil {
		return Column{}, fmt.Errorf("get column: %w", err)
	}

	column.Tasks = []Task{}
	return column, nil
}

func (r *Repository) getTaskTx(ctx context.Context, tx *sql.Tx, id int64) (Task, error) {
	row := tx.QueryRowContext(ctx, `
		SELECT id, column_id, title, description, priority, COALESCE(due_date, ''), position, created_at, updated_at
		FROM tasks
		WHERE id = ?`, id)

	task, err := scanTask(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrNotFound
	}

	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (r *Repository) listColumnIDsTx(ctx context.Context, tx *sql.Tx) ([]int64, error) {
	rows, err := tx.QueryContext(ctx, `SELECT id FROM columns ORDER BY order_index, id`)
	if err != nil {
		return nil, fmt.Errorf("list column ids: %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan column id: %w", err)
		}

		ids = append(ids, id)
	}

	return ids, rows.Err()
}

func (r *Repository) listTaskIDsTx(ctx context.Context, tx *sql.Tx, columnID int64) ([]int64, error) {
	rows, err := tx.QueryContext(ctx, `SELECT id FROM tasks WHERE column_id = ? ORDER BY position, id`, columnID)
	if err != nil {
		return nil, fmt.Errorf("list task ids: %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan task id: %w", err)
		}

		ids = append(ids, id)
	}

	return ids, rows.Err()
}

func (r *Repository) requireColumnTx(ctx context.Context, tx *sql.Tx, id int64) error {
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM columns WHERE id = ?`, id).Scan(&count); err != nil {
		return fmt.Errorf("check column: %w", err)
	}

	if count == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) resequenceColumns(ctx context.Context, tx *sql.Tx) error {
	ids, err := r.listColumnIDsTx(ctx, tx)
	if err != nil {
		return err
	}

	for index, id := range ids {
		if _, err := tx.ExecContext(ctx, `UPDATE columns SET order_index = ?, updated_at = ? WHERE id = ?`, index, nowRFC3339(), id); err != nil {
			return fmt.Errorf("resequence columns: %w", err)
		}
	}

	return nil
}

func (r *Repository) resequenceTasksInColumn(ctx context.Context, tx *sql.Tx, columnID int64) error {
	ids, err := r.listTaskIDsTx(ctx, tx, columnID)
	if err != nil {
		return err
	}

	return r.applyTaskOrder(ctx, tx, columnID, ids, nowRFC3339())
}

func (r *Repository) applyTaskOrder(ctx context.Context, tx *sql.Tx, columnID int64, ids []int64, updatedAt string) error {
	for index, id := range ids {
		if _, err := tx.ExecContext(ctx, `
			UPDATE tasks
			SET column_id = ?, position = ?, updated_at = ?
			WHERE id = ?`,
			columnID, index, updatedAt, id); err != nil {
			return fmt.Errorf("apply task order: %w", err)
		}
	}

	return nil
}

func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func scanTask(scanner interface {
	Scan(dest ...any) error
}) (Task, error) {
	var task Task
	err := scanner.Scan(
		&task.ID,
		&task.ColumnID,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.DueDate,
		&task.Position,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return Task{}, fmt.Errorf("scan task: %w", err)
	}

	return task, nil
}

func removeID(ids []int64, target int64) []int64 {
	filtered := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id != target {
			filtered = append(filtered, id)
		}
	}

	return filtered
}

func insertID(ids []int64, target int64, index int) []int64 {
	if index < 0 {
		index = 0
	}

	if index > len(ids) {
		index = len(ids)
	}

	result := make([]int64, 0, len(ids)+1)
	result = append(result, ids[:index]...)
	result = append(result, target)
	result = append(result, ids[index:]...)
	return result
}
