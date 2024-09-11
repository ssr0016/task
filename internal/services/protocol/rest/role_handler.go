package rest

import (
	"task/internal/api/errors"
	"task/internal/api/response"
	"task/internal/services/accesscontrol/role"

	"github.com/gofiber/fiber/v2"
)

type roleHandler struct {
	s role.Service
}

func NewRoleHandler(s role.Service) *roleHandler {
	return &roleHandler{
		s: s,
	}
}

func (h *roleHandler) CreateRole(ctx *fiber.Ctx) error {
	var cmd role.CreateRoleCommand

	if err := ctx.BodyParser(&cmd); err != nil {
		return errors.ErrorBadRequest(err)
	}

	if err := cmd.Validate(); err != nil {
		return errors.ErrorBadRequest(err)
	}

	if err := h.s.CreateRole(ctx.Context(), &cmd); err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Created(ctx, fiber.Map{
		"role": cmd,
	})
}

func (h *roleHandler) GetRoleByID(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	result, err := h.s.GetRoleByID(ctx.Context(), id)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"role": result,
	})
}

func (h *roleHandler) GetRoles(ctx *fiber.Ctx) error {
	result, err := h.s.GetRoles(ctx.Context())
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"roles": result,
	})
}

func (h *roleHandler) DeleteRole(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	err := h.s.DeleteRole(ctx.Context(), id)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"role deleted": id,
	})
}
