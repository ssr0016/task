package rest

import (
	"task/internal/api/errors"
	"task/internal/api/response"
	"task/internal/identity/task"

	"github.com/gofiber/fiber/v2"
)

type taskHandler struct {
	s task.Service
}

func NewTaskHandler(s task.Service) *taskHandler {
	return &taskHandler{
		s: s,
	}
}

func (h *taskHandler) CreateTask(ctx *fiber.Ctx) error {
	var cmd task.CreateTaskCommand

	if err := ctx.BodyParser(&cmd); err != nil {
		return errors.ErrorBadRequest(err)
	}

	if err := cmd.Validate(); err != nil {
		return errors.ErrorBadRequest(err)
	}

	if err := h.s.CreateTask(ctx.Context(), &cmd); err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Created(ctx, fiber.Map{
		"task created successfully!": cmd,
	})
}

func (h *taskHandler) UpdateTask(ctx *fiber.Ctx) error {
	var cmd task.UpdateTaskCommand

	if err := ctx.BodyParser(&cmd); err != nil {
		return errors.ErrorBadRequest(err)
	}

	if err := cmd.Validate(); err != nil {
		return errors.ErrorBadRequest(err)
	}

	if err := h.s.UpdateTask(ctx.Context(), &cmd); err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"task updated successfully!": cmd,
	})
}

func (h *taskHandler) GetTaskByID(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	result, err := h.s.GetTaskByID(ctx.Context(), id)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"task": result,
	})
}

func (h *taskHandler) DeleteTask(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	if err := h.s.DeleteTask(ctx.Context(), id); err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"task deleted successfully!": id,
	})
}

func (h *taskHandler) SearchTask(ctx *fiber.Ctx) error {
	var query task.SearchTaskQuery

	if err := ctx.QueryParser(&query); err != nil {
		return errors.ErrorBadRequest(err)
	}

	result, err := h.s.SearchTask(ctx.Context(), &query)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"tasks": result,
	})
}
