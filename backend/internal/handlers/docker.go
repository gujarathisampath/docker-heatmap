package handlers

import (
	"context"
	"regexp"
	"time"

	"docker-heatmap/internal/middleware"
	"docker-heatmap/internal/services"

	"github.com/gofiber/fiber/v2"
)

// Docker username validation: 4-30 chars, alphanumeric with allowed special chars
var dockerUsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{2,29}$`)

type DockerHandler struct {
	dockerService *services.DockerHubService
}

func NewDockerHandler() *DockerHandler {
	return &DockerHandler{
		dockerService: services.NewDockerHubService(),
	}
}

type ConnectDockerRequest struct {
	DockerUsername string `json:"docker_username"`
	AccessToken    string `json:"access_token"`
}

// ConnectDocker connects a Docker Hub account
func (h *DockerHandler) ConnectDocker(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req ConnectDockerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.DockerUsername == "" || req.AccessToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Docker username and access token are required",
		})
	}

	// Security: Validate username format
	if !dockerUsernameRegex.MatchString(req.DockerUsername) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Docker username format",
		})
	}

	// Security: Validate token length (Docker PATs are typically 36+ chars)
	if len(req.AccessToken) < 10 || len(req.AccessToken) > 500 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid access token length",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	account, err := h.dockerService.ConnectAccount(ctx, user.ID, req.DockerUsername, req.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Docker account connected successfully",
		"account": fiber.Map{
			"id":              account.ID,
			"docker_username": account.DockerUsername,
			"is_active":       account.IsActive,
		},
	})
}

// GetDockerAccount returns the user's connected Docker account
func (h *DockerHandler) GetDockerAccount(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	account, err := h.dockerService.GetDockerAccount(user.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No Docker account connected",
		})
	}

	return c.JSON(fiber.Map{
		"account": fiber.Map{
			"id":               account.ID,
			"docker_username":  account.DockerUsername,
			"is_active":        account.IsActive,
			"auto_refresh":     account.AutoRefresh,
			"last_sync_at":     account.LastSyncAt,
			"last_sync_error":  account.LastSyncError,
			"sync_in_progress": account.SyncInProgress,
		},
	})
}

// DisconnectDocker removes the Docker Hub account connection
func (h *DockerHandler) DisconnectDocker(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	account, err := h.dockerService.GetDockerAccount(user.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No Docker account connected",
		})
	}

	if err := h.dockerService.DisconnectAccount(user.ID, account.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to disconnect account",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Docker account disconnected successfully",
	})
}

// SyncDockerActivity triggers a manual sync of Docker activity
func (h *DockerHandler) SyncDockerActivity(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	account, err := h.dockerService.GetDockerAccount(user.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No Docker account connected",
		})
	}

	if account.SyncInProgress {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Sync already in progress",
		})
	}

	// Trigger sync in background
	go h.dockerService.SyncActivity(context.Background(), account.ID)

	return c.JSON(fiber.Map{
		"message": "Sync started",
	})
}
