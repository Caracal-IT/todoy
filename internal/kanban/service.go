package kanban

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrNotFound indicates a requested resource does not exist.
	ErrNotFound = errors.New("resource not found")
)

// ValidationError indicates a client-side validation problem.
type ValidationError struct {
	Message string
}

// Error returns the error message.
func (e ValidationError) Error() string {
	return e.Message
}

// Service coordinates Kanban validation and persistence.
type Service struct {
	repository *Repository
}

// NewService creates a Kanban service.
func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

// Bootstrap ensures the schema exists and the default board is seeded.
func (s *Service) Bootstrap(ctx context.Context) error {
	if err := s.repository.EnsureSchema(ctx); err != nil {
		return fmt.Errorf("ensure schema: %w", err)
	}

	if err := s.repository.SeedDefaults(ctx); err != nil {
		return fmt.Errorf("seed defaults: %w", err)
	}

	return nil
}

// GetBoard returns the full board projection.
func (s *Service) GetBoard(ctx context.Context) (Board, error) {
	return s.repository.GetBoard(ctx)
}

// CreateColumn adds a column to the workflow.
func (s *Service) CreateColumn(ctx context.Context, input CreateColumnInput) (Column, error) {
	columnName := strings.TrimSpace(input.Name)
	if columnName == "" {
		return Column{}, ValidationError{Message: "column name is required"}
	}

	columnColor := strings.TrimSpace(input.Color)
	if columnColor == "" {
		columnColor = "#8b5cf6"
	}

	return s.repository.CreateColumn(ctx, CreateColumnInput{
		Name:  columnName,
		Color: columnColor,
	})
}

// UpdateColumn changes a column name or color.
func (s *Service) UpdateColumn(ctx context.Context, id int64, input UpdateColumnInput) (Column, error) {
	columnName := strings.TrimSpace(input.Name)
	if columnName == "" {
		return Column{}, ValidationError{Message: "column name is required"}
	}

	columnColor := strings.TrimSpace(input.Color)
	if columnColor == "" {
		return Column{}, ValidationError{Message: "column color is required"}
	}

	return s.repository.UpdateColumn(ctx, id, UpdateColumnInput{
		Name:  columnName,
		Color: columnColor,
	})
}

// DeleteColumn removes a column and its tasks.
func (s *Service) DeleteColumn(ctx context.Context, id int64) error {
	return s.repository.DeleteColumn(ctx, id)
}

// ReorderColumns updates the workflow order.
func (s *Service) ReorderColumns(ctx context.Context, input ReorderColumnsInput) error {
	if len(input.ColumnIDs) == 0 {
		return ValidationError{Message: "columnIds must not be empty"}
	}

	return s.repository.ReorderColumns(ctx, input)
}

// CreateTask adds a new task card.
func (s *Service) CreateTask(ctx context.Context, input CreateTaskInput) (Task, error) {
	if err := validateTaskFields(input.Title, input.Priority, input.DueDate); err != nil {
		return Task{}, err
	}

	return s.repository.CreateTask(ctx, CreateTaskInput{
		ColumnID:    input.ColumnID,
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Priority:    normalizePriority(input.Priority),
		DueDate:     strings.TrimSpace(input.DueDate),
	})
}

// UpdateTask edits the content of an existing task card.
func (s *Service) UpdateTask(ctx context.Context, id int64, input UpdateTaskInput) (Task, error) {
	if err := validateTaskFields(input.Title, input.Priority, input.DueDate); err != nil {
		return Task{}, err
	}

	return s.repository.UpdateTask(ctx, id, UpdateTaskInput{
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Priority:    normalizePriority(input.Priority),
		DueDate:     strings.TrimSpace(input.DueDate),
	})
}

// DeleteTask removes a task from the board.
func (s *Service) DeleteTask(ctx context.Context, id int64) error {
	return s.repository.DeleteTask(ctx, id)
}

// MoveTask changes a task position or column.
func (s *Service) MoveTask(ctx context.Context, id int64, input MoveTaskInput) (Task, error) {
	if input.TargetColumnID <= 0 {
		return Task{}, ValidationError{Message: "targetColumnId must be greater than zero"}
	}

	if input.TargetPosition < 0 {
		return Task{}, ValidationError{Message: "targetPosition must not be negative"}
	}

	return s.repository.MoveTask(ctx, id, input)
}

func validateTaskFields(title string, priority string, dueDate string) error {
	taskTitle := strings.TrimSpace(title)
	if taskTitle == "" {
		return ValidationError{Message: "task title is required"}
	}

	if !isValidPriority(normalizePriority(priority)) {
		return ValidationError{Message: "priority must be low, medium, high, or critical"}
	}

	trimmedDueDate := strings.TrimSpace(dueDate)
	if trimmedDueDate == "" {
		return nil
	}

	if _, err := time.Parse("2006-01-02", trimmedDueDate); err != nil {
		return ValidationError{Message: "dueDate must use YYYY-MM-DD format"}
	}

	return nil
}

func normalizePriority(priority string) string {
	normalized := strings.ToLower(strings.TrimSpace(priority))
	if normalized == "" {
		return "medium"
	}

	return normalized
}

func isValidPriority(priority string) bool {
	switch priority {
	case "low", "medium", "high", "critical":
		return true
	default:
		return false
	}
}
