package httpapi

import (
	"errors"
	"strconv"

	"github.com/Caracal-IT/todoy/internal/kanban"
	"github.com/gofiber/fiber/v2"
)

type handler struct {
	service *kanban.Service
}

func newHandler(service *kanban.Service) *handler {
	return &handler{service: service}
}

func (h *handler) getBoard(c *fiber.Ctx) error {
	board, err := h.service.GetBoard(c.UserContext())
	if err != nil {
		return respondError(c, err)
	}

	return c.JSON(board)
}

func (h *handler) createColumn(c *fiber.Ctx) error {
	var input kanban.CreateColumnInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid column payload")
	}

	column, err := h.service.CreateColumn(c.UserContext(), input)
	if err != nil {
		return respondError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(column)
}

func (h *handler) updateColumn(c *fiber.Ctx) error {
	columnID, err := parseID(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid column id")
	}

	var input kanban.UpdateColumnInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid column payload")
	}

	column, err := h.service.UpdateColumn(c.UserContext(), columnID, input)
	if err != nil {
		return respondError(c, err)
	}

	return c.JSON(column)
}

func (h *handler) deleteColumn(c *fiber.Ctx) error {
	columnID, err := parseID(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid column id")
	}

	if err := h.service.DeleteColumn(c.UserContext(), columnID); err != nil {
		return respondError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *handler) reorderColumns(c *fiber.Ctx) error {
	var input kanban.ReorderColumnsInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid column reorder payload")
	}

	if err := h.service.ReorderColumns(c.UserContext(), input); err != nil {
		return respondError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *handler) createTask(c *fiber.Ctx) error {
	var input kanban.CreateTaskInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task payload")
	}

	task, err := h.service.CreateTask(c.UserContext(), input)
	if err != nil {
		return respondError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

func (h *handler) updateTask(c *fiber.Ctx) error {
	taskID, err := parseID(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task id")
	}

	var input kanban.UpdateTaskInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task payload")
	}

	task, err := h.service.UpdateTask(c.UserContext(), taskID, input)
	if err != nil {
		return respondError(c, err)
	}

	return c.JSON(task)
}

func (h *handler) deleteTask(c *fiber.Ctx) error {
	taskID, err := parseID(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task id")
	}

	if err := h.service.DeleteTask(c.UserContext(), taskID); err != nil {
		return respondError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *handler) moveTask(c *fiber.Ctx) error {
	taskID, err := parseID(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task id")
	}

	var input kanban.MoveTaskInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid move payload")
	}

	task, err := h.service.MoveTask(c.UserContext(), taskID, input)
	if err != nil {
		return respondError(c, err)
	}

	return c.JSON(task)
}

func parseID(raw string) (int64, error) {
	return strconv.ParseInt(raw, 10, 64)
}

func respondError(c *fiber.Ctx, err error) error {
	var validationErr kanban.ValidationError

	switch {
	case errors.As(err, &validationErr):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Message})
	case errors.Is(err, kanban.ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "resource not found"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
