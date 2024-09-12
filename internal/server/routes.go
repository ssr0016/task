package server

import (
	"context"
	"errors"
	"task/internal/api/response"
	"task/internal/db"
	"task/internal/middleware"
	"task/internal/services/department/departmentimpl"
	"task/internal/services/protocol/rest"
	"task/internal/services/user/userimpl"

	"github.com/gofiber/fiber/v2"
)

var (
	requireCreateUser = middleware.RequirePermission("create")
	requireReadUser   = middleware.RequirePermission("read")
	requireUpdateUser = middleware.RequirePermission("update")
	requireDeleteUser = middleware.RequirePermission("delete")

	reqOnlyBySuperuser      = middleware.RequireRole("superuser")
	reqBothUserAndSuperuser = middleware.RequireRole("user", "superuser")
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
	userHttp := rest.NewUserHandler(user)

	api.Post("/users/register", userHttp.RegisterUser)
	api.Post("/users/login", userHttp.LoginUser)

	api.Use(middleware.JWTProtected(s.jwtSecret, user))
	api.Post("/users", reqOnlyBySuperuser, requireCreateUser, userHttp.CreateUser)
	api.Get("/users", reqBothUserAndSuperuser, requireReadUser, userHttp.SearchUser)
	api.Get("/users/:id", reqBothUserAndSuperuser, requireReadUser, userHttp.GetUserByID)
	api.Put("/users/:id", reqOnlyBySuperuser, requireUpdateUser, userHttp.UpdateUser)
	api.Delete("/users/:id", reqOnlyBySuperuser, requireDeleteUser, userHttp.DeleteUser)

	// Logout
	api.Post("/users/logout", reqBothUserAndSuperuser, userHttp.LogoutUser)

	// Department Routes
	department := departmentimpl.NewService(s.db, s.cfg)
	departmentHttp := rest.NewDepartmentHandler(department)

	api.Post("/departments", reqOnlyBySuperuser, requireCreateUser, departmentHttp.CreateDepartment)
	api.Get("/departments", reqBothUserAndSuperuser, requireReadUser, departmentHttp.SearchDepartment)
	api.Get("/departments/:id", reqBothUserAndSuperuser, requireReadUser, departmentHttp.GetDepartmentByID)
	api.Put("/departments/:id", reqOnlyBySuperuser, requireUpdateUser, departmentHttp.UpdateDepartment)
	api.Delete("/departments/:id", reqOnlyBySuperuser, requireDeleteUser, departmentHttp.DeleteDepartment)

	api.Post("/users/assigned/departments", reqOnlyBySuperuser, requireUpdateUser, departmentHttp.AssignUserToDepartment)
	api.Get("/users/assigned/:id/departments", reqBothUserAndSuperuser, requireReadUser, departmentHttp.GetUsersByDepartment)
	api.Delete("/users/assigned/:id/departments", reqOnlyBySuperuser, requireUpdateUser, departmentHttp.RemoveUserFromDepartment)
}
