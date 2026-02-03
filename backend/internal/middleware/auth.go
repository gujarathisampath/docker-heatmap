package middleware

import (
	"strings"

	"docker-heatmap/internal/database"
	"docker-heatmap/internal/models"
	"docker-heatmap/internal/utils"

	"github.com/gofiber/fiber/v2"
)

const (
	UserContextKey = "user"
)

// AuthMiddleware validates JWT tokens and adds user to context
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		tokenString := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Fetch user from database
		var user models.User
		if err := database.DB.First(&user, claims.UserID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		// Add user to context
		c.Locals(UserContextKey, &user)

		return c.Next()
	}
}

// OptionalAuthMiddleware tries to authenticate but doesn't require it
func OptionalAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Next()
		}

		var user models.User
		if err := database.DB.First(&user, claims.UserID).Error; err != nil {
			return c.Next()
		}

		c.Locals(UserContextKey, &user)
		return c.Next()
	}
}

// GetUserFromContext retrieves the authenticated user from context
func GetUserFromContext(c *fiber.Ctx) *models.User {
	user, ok := c.Locals(UserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
