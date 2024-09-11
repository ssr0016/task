package server

import (
	"context"
	"errors"
	"task/internal/api/response"
	"task/internal/db"
	"task/internal/middleware"
	"task/internal/services/protocol/rest"
	"task/internal/services/user/userimpl"

	"github.com/gofiber/fiber/v2"
)

var (
	requireCreateUser = middleware.RequirePermission("create")
	requireReadUser   = middleware.RequirePermission("read")
	requireUpdateUser = middleware.RequirePermission("update")
	requireDeleteUser = middleware.RequirePermission("delete")

	superuser   = middleware.RequireRole("superuser")
	defaultuser = middleware.RequireRole("user")
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
	api.Post("/users", superuser, requireCreateUser, userHttp.CreateUser)
	api.Get("/users", superuser, defaultuser, requireReadUser, userHttp.SearchUser)
	api.Get("/users/:id", superuser, defaultuser, requireReadUser, userHttp.GetUserByID)
	api.Put("/users/:id", superuser, requireUpdateUser, userHttp.UpdateUser)
	api.Delete("/users/:id", superuser, requireDeleteUser, userHttp.DeleteUser)

}
