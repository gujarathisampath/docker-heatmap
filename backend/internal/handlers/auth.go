package handlers

import (
	"context"
	"time"

	"docker-heatmap/internal/config"
	"docker-heatmap/internal/middleware"
	"docker-heatmap/internal/services"
	"docker-heatmap/internal/utils"

	"sync"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *services.GitHubAuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewGitHubAuthService(),
	}
}

// OAuthState stores temporary state for OAuth flow
var (
	oauthStates = make(map[string]time.Time)
	stateMutex  sync.Mutex
)

// InitiateGitHubAuth starts the GitHub OAuth flow
func (h *AuthHandler) InitiateGitHubAuth(c *fiber.Ctx) error {
	state, err := utils.GenerateStateToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate state",
		})
	}

	// Store state with expiry
	stateMutex.Lock()
	oauthStates[state] = time.Now().Add(10 * time.Minute)
	stateMutex.Unlock()

	// Clean old states
	go cleanupOAuthStates()

	authURL := h.authService.GetAuthURL(state)

	return c.JSON(fiber.Map{
		"auth_url": authURL,
	})
}

// GitHubCallback handles the OAuth callback
func (h *AuthHandler) GitHubCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		return c.Redirect(config.AppConfig.FrontendURL + "/auth/error?message=missing_params")
	}

	// Validate state
	stateMutex.Lock()
	expiry, exists := oauthStates[state]
	if exists {
		delete(oauthStates, state)
	}
	stateMutex.Unlock()

	if !exists || time.Now().After(expiry) {
		return c.Redirect(config.AppConfig.FrontendURL + "/auth/error?message=invalid_state")
	}

	// Exchange code for user
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := h.authService.ExchangeCode(ctx, code)
	if err != nil {
		return c.Redirect(config.AppConfig.FrontendURL + "/auth/error?message=auth_failed")
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, user.GitHubUsername)
	if err != nil {
		return c.Redirect(config.AppConfig.FrontendURL + "/auth/error?message=token_failed")
	}

	// Redirect to frontend with token
	return c.Redirect(config.AppConfig.FrontendURL + "/auth/callback?token=" + token)
}

// GetCurrentUser returns the authenticated user
func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
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

// Logout invalidates the current session
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Since we use stateless JWT, just return success
	// Client should delete the token
	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func cleanupOAuthStates() {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	now := time.Now()
	for state, expiry := range oauthStates {
		if now.After(expiry) {
			delete(oauthStates, state)
		}
	}
}
