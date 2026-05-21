package httpapi

import (
	"github.com/Caracal-IT/todoy/internal/kanban"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// NewApp creates the Fiber application and registers REST routes.
func NewApp(service *kanban.Service) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Todoy API",
	})

	app.Use(cors.New())

	handler := newHandler(service)
	api := app.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	api.Get("/board", handler.getBoard)
	api.Post("/columns", handler.createColumn)
	api.Patch("/columns/:id", handler.updateColumn)
	api.Delete("/columns/:id", handler.deleteColumn)
	api.Post("/columns/reorder", handler.reorderColumns)
	api.Post("/tasks", handler.createTask)
	api.Patch("/tasks/:id", handler.updateTask)
	api.Delete("/tasks/:id", handler.deleteTask)
	api.Post("/tasks/:id/move", handler.moveTask)

	return app
}
