package rest

import (
	"task/internal/api/errors"
	"task/internal/api/response"
	"task/internal/services/user"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	s user.Service
}

func NewUserHandler(s user.Service) *userHandler {
	return &userHandler{
		s: s,
	}
}

func (h *userHandler) CreateUser(ctx *fiber.Ctx) error {
	var cmd user.CreateUserCommand

	if err := ctx.BodyParser(&cmd); err != nil {
		return errors.ErrorBadRequest(err)
	}

	if err := cmd.Validate(); err != nil {
		return errors.ErrorBadRequest(err)
	}

	if err := h.s.CreateUser(ctx.Context(), &cmd); err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Created(ctx, fiber.Map{
		"user": cmd,
	})
}

func (h *userHandler) GetUserByID(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	result, err := h.s.GetUserByID(ctx.Context(), id)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"user": result,
	})
}

func (h *userHandler) UpdateUser(ctx *fiber.Ctx) error {
	var cmd user.UpdateUserCommand

	if err := ctx.BodyParser(&cmd); err != nil {
		return errors.ErrorBadRequest(err)
	}

	err := h.s.UpdateUser(ctx.Context(), &cmd)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"user updated": cmd,
	})
}

func (h *userHandler) SearchUser(ctx *fiber.Ctx) error {
	var query user.SearchUserQuery

	if err := ctx.QueryParser(&query); err != nil {
		return errors.ErrorBadRequest(err)
	}

	result, err := h.s.SearchUser(ctx.Context(), &query)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"users": result,
	})
}

func (h *userHandler) DeleteUser(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	if err := h.s.DeleteUser(ctx.Context(), id); err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"message": "user deleted successfully!",
	})
}
