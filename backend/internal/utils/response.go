package utils

import "github.com/gofiber/fiber/v2"

// APIResponse represents a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// SuccessResponse returns a standardized success response
func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// SuccessWithMessage returns a success response with a message
func SuccessWithMessage(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// ErrorResponse returns a standardized error response
func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}

// ValidationError returns a 400 Bad Request with validation error message
func ValidationError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message)
}

// UnauthorizedError returns a 401 Unauthorized response
func UnauthorizedError(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	return ErrorResponse(c, fiber.StatusUnauthorized, message)
}

// NotFoundError returns a 404 Not Found response
func NotFoundError(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Resource not found"
	}
	return ErrorResponse(c, fiber.StatusNotFound, message)
}

// InternalError returns a 500 Internal Server Error response
func InternalError(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "An internal error occurred"
	}
	return ErrorResponse(c, fiber.StatusInternalServerError, message)
}
