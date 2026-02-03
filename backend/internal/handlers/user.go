package handlers

import (
	"docker-heatmap/internal/database"
	"docker-heatmap/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

type UpdateProfileRequest struct {
	Name          string `json:"name"`
	Bio           string `json:"bio"`
	PublicProfile *bool  `json:"public_profile"`
}

// GetProfile returns the current user's profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

// UpdateProfile updates the current user's profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.PublicProfile != nil {
		user.PublicProfile = *req.PublicProfile
	}

	if err := database.DB.Save(user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

// GetEmbedCode returns embed code snippets for the user's heatmap
func (h *UserHandler) GetEmbedCode(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	baseURL := c.BaseURL()
	dockerUsername := c.Query("docker_username")

	if dockerUsername == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Docker username required",
		})
	}

	svgURL := baseURL + "/api/heatmap/" + dockerUsername + ".svg"
	jsonURL := baseURL + "/api/activity/" + dockerUsername + ".json"

	return c.JSON(fiber.Map{
		"svg_url":   svgURL,
		"json_url":  jsonURL,
		"markdown":  "![Docker Activity](" + svgURL + ")",
		"html":      `<img src="` + svgURL + `" alt="Docker Activity Heatmap" />`,
		"html_link": `<a href="` + baseURL + `/profile/` + dockerUsername + `"><img src="` + svgURL + `" alt="Docker Activity Heatmap" /></a>`,
	})
}
