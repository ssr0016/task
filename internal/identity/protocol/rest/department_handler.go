package rest

import (
	"strconv"
	"task/internal/api/errors"
	"task/internal/api/response"
	"task/internal/identity/department"

	"github.com/gofiber/fiber/v2"
)

type departmentHandler struct {
	s department.Service
}

func NewDepartmentHandler(s department.Service) *departmentHandler {
	return &departmentHandler{
		s: s,
	}
}

func (h *departmentHandler) CreateDepartment(ctx *fiber.Ctx) error {
	var cmd department.CreateDepartmentCommand

	if err := ctx.BodyParser(&cmd); err != nil {
		return errors.ErrorBadRequest(err)
	}

	if err := cmd.Validate(); err != nil {
		return errors.ErrorBadRequest(err)
	}

	if err := h.s.CreateDepartment(ctx.Context(), &cmd); err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Created(ctx, fiber.Map{
		"department created successfully!": cmd,
	})
}

func (h *departmentHandler) UpdateDepartment(ctx *fiber.Ctx) error {
	var cmd department.UpdateDepartmentCommand

	if err := ctx.BodyParser(&cmd); err != nil {
		return errors.ErrorBadRequest(err)
	}

	if err := cmd.Validate(); err != nil {
		return errors.ErrorBadRequest(err)
	}

	if err := h.s.UpdateDepartment(ctx.Context(), &cmd); err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"department updated successfully!": cmd,
	})
}

func (h *departmentHandler) GetDepartmentByID(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	result, err := h.s.GetDepartmentByID(ctx.Context(), id)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"department": result,
	})
}

func (h *departmentHandler) SearchDepartment(ctx *fiber.Ctx) error {
	var query department.SearchDepartmentQuery

	if err := ctx.QueryParser(&query); err != nil {
		return errors.ErrorBadRequest(err)
	}

	result, err := h.s.SearchDepartment(ctx.Context(), &query)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"departments": result,
	})
}

func (h *departmentHandler) DeleteDepartment(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	if err := h.s.DeleteDepartment(ctx.Context(), id); err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"department deleted successfully!": id,
	})
}

func (h *departmentHandler) AssignUserToDepartment(ctx *fiber.Ctx) error {
	var cmd department.AssignUserToDepartmentCommand

	if err := ctx.BodyParser(&cmd); err != nil {
		return errors.ErrorBadRequest(err)
	}

	if err := h.s.AssignUserToDepartment(ctx.Context(), &cmd); err != nil {
		return errors.ErrorInternalServerError(err)
	}
	return response.Ok(ctx, fiber.Map{
		"department assigned to user successfully!": cmd,
	})
}

func (h *departmentHandler) GetUsersByDepartment(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	result, err := h.s.GetUsersByDepartment(ctx.Context(), id)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"users": result,
	})
}

func (h *departmentHandler) RemoveUserFromDepartment(ctx *fiber.Ctx) error {
	userID, _ := strconv.Atoi(ctx.Params("id"))

	if err := h.s.RemoveUserFromDepartment(ctx.Context(), userID); err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"user removed from department successfully!": userID,
	})
}

func (h *departmentHandler) SearchAllUsersByDepartment(ctx *fiber.Ctx) error {
	var query department.SearchAllUsersByDepartmentQuery

	if err := ctx.QueryParser(&query); err != nil {
		return errors.ErrorBadRequest(err)
	}

	result, err := h.s.SearchAllUsersByDepartment(ctx.Context(), &query)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"users": result,
	})
}
