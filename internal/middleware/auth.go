package middleware

import (
	"context"
	"fmt"
	"strings"
	"task/pkg/util/jwt"
	"time"

	"task/internal/identity/accesscontrol"
	"task/internal/identity/monitoringactivities"
	"task/internal/identity/user"

	"github.com/gofiber/fiber/v2"
)

// Middleware to check if the user has a valid JWT
func JWTProtected(secret string, service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or malformed JWT",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token format",
			})
		}

		tokenStr := authHeader[len("Bearer "):]

		claims, err := jwt.ValidateToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired JWT",
			})
		}

		// Check if the token is blacklisted
		isBlacklisted, err := service.IsTokenBlacklisted(context.Background(), tokenStr)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error while checking token",
			})
		}

		if isBlacklisted {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token is blacklisted",
			})
		}

		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// Middleware to check if the user has the required role
func RequireRole(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string) // Get the user's role from context

		// Check if the user's role is in the list of allowed roles
		roleAllowed := false
		for _, allowedRole := range requiredRoles {
			if role == allowedRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied. Insufficient permissions.",
			})
		}

		return c.Next()
	}
}

// RequirePermission checks if the user has the required permission
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)

		if !accesscontrol.HasPermission(role, permission) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You do not have permission to access this resource",
			})
		}

		if !accesscontrol.HasTaskPermission(role, permission) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You do not have permission to access this resource",
			})
		}

		return c.Next()
	}
}

// ActivityLoggingMiddleware logs the activity of the user
func NewActivityLoggingMiddleware(service monitoringactivities.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve userID from context, which is expected to be a string
		userIDStr, ok := c.Locals("userID").(string)
		if !ok {
			fmt.Println("userID not found or not a string")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "UserID not found or not a string",
			})
		}

		// Create an activity log command
		activityLogCommand := &monitoringactivities.CreateActivityLogCommand{
			UserID:    userIDStr, // UserID is a string now
			Activity:  c.Method(),
			Action:    c.Path(),
			Resource:  c.OriginalURL(),
			Details:   "",                              // Add any additional details if available
			CreatedAt: time.Now().Format(time.RFC3339), // Set current time in RFC3339 format
		}

		// Log activity
		err := service.LogActivity(c.Context(), activityLogCommand)
		if err != nil {
			fmt.Printf("Error logging activity: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to log activity",
			})
		}

		// Continue with the next middleware/handler
		return c.Next()
	}
}
