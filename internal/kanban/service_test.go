package kanban

import (
	"context"
	"path/filepath"
	"testing"

	sqlitestore "github.com/Caracal-IT/todoy/internal/platform/sqlite"
)

func TestServiceSeedsDefaultBoard(t *testing.T) {
	t.Parallel()

	service := newTestService(t)

	board, err := service.GetBoard(context.Background())
	if err != nil {
		t.Fatalf("get board: %v", err)
	}

	if len(board.Columns) != 4 {
		t.Fatalf("expected 4 default columns, got %d", len(board.Columns))
	}

	if board.Columns[0].Name != "Backlog" || board.Columns[3].Name != "Done" {
		t.Fatalf("unexpected default columns: %#v", board.Columns)
	}
}

func TestServiceTaskLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	service := newTestService(t)

	board, err := service.GetBoard(ctx)
	if err != nil {
		t.Fatalf("get board: %v", err)
	}

	task, err := service.CreateTask(ctx, CreateTaskInput{
		ColumnID:    board.Columns[0].ID,
		Title:       "Ship onboarding flow",
		Description: "Finalize the no-auth onboarding board copy.",
		Priority:    "high",
		DueDate:     "2026-06-01",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	updatedTask, err := service.UpdateTask(ctx, task.ID, UpdateTaskInput{
		Title:       "Ship onboarding board flow",
		Description: "Finalize copy and UI spacing.",
		Priority:    "critical",
		DueDate:     "2026-06-03",
	})
	if err != nil {
		t.Fatalf("update task: %v", err)
	}

	if updatedTask.Priority != "critical" {
		t.Fatalf("expected priority to update, got %s", updatedTask.Priority)
	}

	if err := service.DeleteTask(ctx, task.ID); err != nil {
		t.Fatalf("delete task: %v", err)
	}

	afterDeleteBoard, err := service.GetBoard(ctx)
	if err != nil {
		t.Fatalf("get board after delete: %v", err)
	}

	if len(afterDeleteBoard.Columns[0].Tasks) != 0 {
		t.Fatalf("expected task to be deleted, got %d tasks", len(afterDeleteBoard.Columns[0].Tasks))
	}
}

func TestServiceMoveTaskAndReorderColumns(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	service := newTestService(t)

	board, err := service.GetBoard(ctx)
	if err != nil {
		t.Fatalf("get board: %v", err)
	}

	firstColumn := board.Columns[0]
	secondColumn := board.Columns[1]

	firstTask, err := service.CreateTask(ctx, CreateTaskInput{
		ColumnID: firstColumn.ID,
		Title:    "Design quick add",
		Priority: "medium",
	})
	if err != nil {
		t.Fatalf("create first task: %v", err)
	}

	secondTask, err := service.CreateTask(ctx, CreateTaskInput{
		ColumnID: firstColumn.ID,
		Title:    "Build task editor",
		Priority: "low",
	})
	if err != nil {
		t.Fatalf("create second task: %v", err)
	}

	movedTask, err := service.MoveTask(ctx, firstTask.ID, MoveTaskInput{
		TargetColumnID: secondColumn.ID,
		TargetPosition: 0,
	})
	if err != nil {
		t.Fatalf("move task: %v", err)
	}

	if movedTask.ColumnID != secondColumn.ID {
		t.Fatalf("expected moved task column %d, got %d", secondColumn.ID, movedTask.ColumnID)
	}

	if err := service.ReorderColumns(ctx, ReorderColumnsInput{
		ColumnIDs: []int64{board.Columns[1].ID, board.Columns[0].ID, board.Columns[2].ID, board.Columns[3].ID},
	}); err != nil {
		t.Fatalf("reorder columns: %v", err)
	}

	reorderedBoard, err := service.GetBoard(ctx)
	if err != nil {
		t.Fatalf("get reordered board: %v", err)
	}

	if reorderedBoard.Columns[0].ID != board.Columns[1].ID {
		t.Fatalf("expected column reorder to persist")
	}

	if len(reorderedBoard.Columns[1].Tasks) != 1 || reorderedBoard.Columns[1].Tasks[0].ID != secondTask.ID {
		t.Fatalf("expected source column to keep remaining task after move")
	}
}

func BenchmarkServiceGetBoard(b *testing.B) {
	service := newTestService(b)
	ctx := context.Background()

	board, err := service.GetBoard(ctx)
	if err != nil {
		b.Fatalf("seed board: %v", err)
	}

	for index := 0; index < 100; index++ {
		if _, err := service.CreateTask(ctx, CreateTaskInput{
			ColumnID: board.Columns[index%len(board.Columns)].ID,
			Title:    "Benchmark task",
			Priority: "medium",
		}); err != nil {
			b.Fatalf("seed task: %v", err)
		}
	}

	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		if _, err := service.GetBoard(ctx); err != nil {
			b.Fatalf("get board: %v", err)
		}
	}
}

func newTestService(tb testing.TB) *Service {
	tb.Helper()

	databasePath := filepath.Join(tb.TempDir(), "test.db")
	db, err := sqlitestore.Open(databasePath)
	if err != nil {
		tb.Fatalf("open test database: %v", err)
	}

	tb.Cleanup(func() {
		db.Close()
	})

	repository := NewRepository(db)
	service := NewService(repository)
	if err := service.Bootstrap(context.Background()); err != nil {
		tb.Fatalf("bootstrap test service: %v", err)
	}

	return service
}
