package server

import (
	"context"
	"errors"
	"task/internal/api/response"
	"task/internal/db"
	"task/internal/middleware"
	"task/internal/services/accesscontrol/role/roleimpl"
	"task/internal/services/protocol/rest"
	"task/internal/services/user/userimpl"

	"github.com/gofiber/fiber/v2"
)

func healthCheck(db db.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var result int
		err := db.Get(context.Background(), &result, "SELECT 1")
		if err != nil {
			return errors.New("database unavailable")
		}
		return response.Ok(ctx, fiber.Map{
			"database": "available",
		})
	}
}

func (s *Server) SetupRoutes() {
	api := s.app.Group("/api")
	api.Get("/health", healthCheck(s.db))

	// User Routes
	user := userimpl.NewService(s.db, s.cfg)
	userHttp := rest.NewUserHandler(&user)

	api.Post("/users/register", userHttp.RegisterUser)
	api.Post("/users/login", userHttp.LoginUser)

	api.Use(middleware.JWTProtected(s.jwtSecret))
	api.Post("/users", userHttp.CreateUser)
	api.Get("/users", userHttp.SearchUser)
	api.Get("/users/:id", userHttp.GetUserByID)
	api.Put("/users/:id", userHttp.UpdateUser)
	api.Delete("/users/:id", userHttp.DeleteUser)

	// Role Routes
	role := roleimpl.NewService(s.db, s.cfg)
	roleHttp := rest.NewRoleHandler(&role)
	api.Post("/roles", roleHttp.CreateRole)
	api.Get("/roles/:id", roleHttp.GetRoleByID)
	api.Get("/roles", roleHttp.GetRoles)
	api.Delete("/roles/:id", roleHttp.DeleteRole)
}
