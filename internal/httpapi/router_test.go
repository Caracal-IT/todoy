package httpapi

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/Caracal-IT/todoy/internal/kanban"
	sqlitestore "github.com/Caracal-IT/todoy/internal/platform/sqlite"
	"github.com/gofiber/fiber/v2"
)

func TestBoardEndpointReturnsSeededColumns(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	request := httptest.NewRequest("GET", "/api/v1/board", nil)
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("send test request: %v", err)
	}

	if response.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}

	var board kanban.Board
	if err := json.NewDecoder(response.Body).Decode(&board); err != nil {
		t.Fatalf("decode board response: %v", err)
	}

	if len(board.Columns) != 4 {
		t.Fatalf("expected 4 columns, got %d", len(board.Columns))
	}
}

func newTestApp(tb testing.TB) *fiber.App {
	tb.Helper()

	databasePath := filepath.Join(tb.TempDir(), "api.db")
	db, err := sqlitestore.Open(databasePath)
	if err != nil {
		tb.Fatalf("open api test db: %v", err)
	}

	tb.Cleanup(func() {
		db.Close()
	})

	repository := kanban.NewRepository(db)
	service := kanban.NewService(repository)
	if err := service.Bootstrap(context.Background()); err != nil {
		tb.Fatalf("bootstrap service: %v", err)
	}

	return NewApp(service)
}
