package kanban

import "time"

// Board represents the full Kanban board response.
type Board struct {
	Columns []Column   `json:"columns"`
	Stats   BoardStats `json:"stats"`
	Now     string     `json:"now"`
}

// BoardStats summarizes the state of the board.
type BoardStats struct {
	TotalTasks        int `json:"totalTasks"`
	CompletedTasks    int `json:"completedTasks"`
	OverdueTasks      int `json:"overdueTasks"`
	HighPriorityTasks int `json:"highPriorityTasks"`
}

// Column represents a workflow lane on the board.
type Column struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	OrderIndex int    `json:"orderIndex"`
	Tasks      []Task `json:"tasks"`
}

// Task represents a card on the board.
type Task struct {
	ID          int64  `json:"id"`
	ColumnID    int64  `json:"columnId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	DueDate     string `json:"dueDate"`
	Position    int    `json:"position"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// CreateColumnInput contains the required data for a new column.
type CreateColumnInput struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// UpdateColumnInput contains editable column fields.
type UpdateColumnInput struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// ReorderColumnsInput contains ordered column identifiers.
type ReorderColumnsInput struct {
	ColumnIDs []int64 `json:"columnIds"`
}

// CreateTaskInput contains the required data for a new task.
type CreateTaskInput struct {
	ColumnID    int64  `json:"columnId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	DueDate     string `json:"dueDate"`
}

// UpdateTaskInput contains editable task fields.
type UpdateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	DueDate     string `json:"dueDate"`
}

// MoveTaskInput contains the target location for a task move.
type MoveTaskInput struct {
	TargetColumnID int64 `json:"targetColumnId"`
	TargetPosition int   `json:"targetPosition"`
}

func (t Task) parsedDueDate() (time.Time, bool) {
	if t.DueDate == "" {
		return time.Time{}, false
	}

	dueDate, err := time.Parse("2006-01-02", t.DueDate)
	if err != nil {
		return time.Time{}, false
	}

	return dueDate, true
}
