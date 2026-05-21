package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/Caracal-IT/todoy/internal/httpapi"
	"github.com/Caracal-IT/todoy/internal/kanban"
	sqlitestore "github.com/Caracal-IT/todoy/internal/platform/sqlite"
)

func main() {
	ctx := context.Background()
	databasePath := getenv("DATABASE_PATH", "todoy.db")

	if err := ensureDatabaseDirectory(databasePath); err != nil {
		log.Fatalf("prepare database path: %v", err)
	}

	db, err := sqlitestore.Open(databasePath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	repository := kanban.NewRepository(db)
	service := kanban.NewService(repository)
	if err := service.Bootstrap(ctx); err != nil {
		log.Fatalf("bootstrap service: %v", err)
	}

	app := httpapi.NewApp(service)
	address := ":" + getenv("PORT", "8080")

	log.Printf("API listening on %s", address)
	if err := app.Listen(address); err != nil {
		log.Fatalf("start server: %v", err)
	}
}

func ensureDatabaseDirectory(databasePath string) error {
	directory := filepath.Dir(databasePath)
	if directory == "." {
		return nil
	}

	return os.MkdirAll(directory, 0o755)
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
