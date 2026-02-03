package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"docker-heatmap/internal/services"

	"github.com/gofiber/fiber/v2"
)

type HeatmapHandler struct {
	heatmapService *services.HeatmapService
	dockerService  *services.DockerHubService
}

func NewHeatmapHandler() *HeatmapHandler {
	return &HeatmapHandler{
		heatmapService: services.NewHeatmapService(),
		dockerService:  services.NewDockerHubService(),
	}
}

// GetHeatmapSVG returns the heatmap as an SVG image with customization options
// Query params:
//   - days: number of days (1-365, default 365)
//   - theme: color theme (github, docker, dracula, nord, etc.) or "custom"
//   - cell_size: size of each cell (5-20, default 11)
//   - radius: border radius of cells (0-10, default 2)
//   - hide_legend: hide the color legend (true/false)
//   - hide_total: hide the total count (true/false)
//   - hide_labels: hide month/day labels (true/false)
//   - title: custom title text
//   - bg_color: custom background color (hex without #)
//   - text_color: custom text color (hex without #)
//   - color0-color4: custom level colors (hex without #)
func (h *HeatmapHandler) GetHeatmapSVG(c *fiber.Ctx) error {
	username := c.Params("username")

	// Remove .svg extension if present
	username = strings.TrimSuffix(username, ".svg")

	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username is required",
		})
	}

	// Parse options from query params
	opts := services.SVGOptions{
		Theme:       c.Query("theme", "github"),
		Days:        365,
		CellSize:    11,
		CellRadius:  2,
		HideLegend:  c.Query("hide_legend") == "true" || c.Query("hide_legend") == "1",
		HideTotal:   c.Query("hide_total") == "true" || c.Query("hide_total") == "1",
		HideLabels:  c.Query("hide_labels") == "true" || c.Query("hide_labels") == "1",
		CustomTitle: c.Query("title"),
	}

	// Parse numeric options with validation
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 365 {
			opts.Days = parsed
		}
	}

	if cs := c.Query("cell_size"); cs != "" {
		if parsed, err := strconv.Atoi(cs); err == nil && parsed >= 5 && parsed <= 20 {
			opts.CellSize = parsed
		}
	}

	if r := c.Query("radius"); r != "" {
		if parsed, err := strconv.Atoi(r); err == nil && parsed >= 0 && parsed <= 10 {
			opts.CellRadius = parsed
		}
	}

	// Parse custom colors
	if bg := c.Query("bg_color"); bg != "" {
		opts.BgColor = parseHexColor(bg)
	}
	if txt := c.Query("text_color"); txt != "" {
		opts.TextColor = parseHexColor(txt)
	}

	// Custom level colors
	customColors := make([]string, 0, 5)
	for i := 0; i < 5; i++ {
		if clr := c.Query(fmt.Sprintf("color%d", i)); clr != "" {
			customColors = append(customColors, parseHexColor(clr))
		}
	}
	if len(customColors) == 5 {
		opts.CustomColors = customColors
		opts.Theme = "custom"
	}

	svg, err := h.heatmapService.GenerateSVGWithOptions(username, opts)
	if err != nil {
		if err == services.ErrDockerAccountNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found or no Docker account connected",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate heatmap",
		})
	}

	c.Set("Content-Type", "image/svg+xml")
	c.Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	return c.Send(svg)
}

// parseHexColor ensures color has # prefix
func parseHexColor(color string) string {
	color = strings.TrimSpace(color)
	if color == "" {
		return ""
	}
	if !strings.HasPrefix(color, "#") {
		return "#" + color
	}
	return color
}

// GetAvailableThemes returns all available SVG themes with details
func (h *HeatmapHandler) GetAvailableThemes(c *fiber.Ctx) error {
	themes := make([]fiber.Map, 0)

	// Define order for themes
	order := []string{
		"github", "github-light", "docker",
		"dracula", "nord", "monokai", "one-dark", "tokyo-night", "catppuccin",
		"ocean", "sunset", "forest", "purple", "rose",
		"minimal", "minimal-dark",
	}

	for _, name := range order {
		if theme, ok := services.Themes[name]; ok {
			themes = append(themes, fiber.Map{
				"id":         name,
				"name":       theme.Name,
				"bg_color":   theme.BgColor,
				"text_color": theme.TextColor,
				"colors":     theme.Colors,
			})
		}
	}

	return c.JSON(fiber.Map{
		"themes": themes,
		"customization": fiber.Map{
			"description": "You can also create custom themes using query parameters",
			"params": fiber.Map{
				"bg_color":   "Background color (hex without #)",
				"text_color": "Text color (hex without #)",
				"color0":     "Level 0 (no activity) color",
				"color1":     "Level 1 (low) color",
				"color2":     "Level 2 (medium) color",
				"color3":     "Level 3 (high) color",
				"color4":     "Level 4 (max) color",
			},
			"example": "/api/heatmap/username.svg?theme=custom&bg_color=1a1a2e&color0=16213e&color1=0f3460&color2=533483&color3=e94560&color4=ff6b6b",
		},
	})
}

// GetActivityJSON returns activity data as JSON
func (h *HeatmapHandler) GetActivityJSON(c *fiber.Ctx) error {
	username := c.Params("username")

	// Remove .json extension if present
	username = strings.TrimSuffix(username, ".json")

	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username is required",
		})
	}

	// Get days parameter (default 365)
	days := 365
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 365 {
			days = parsed
		}
	}

	activities, err := h.dockerService.GetActivitySummary(username, days)
	if err != nil {
		if err == services.ErrDockerAccountNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found or no Docker account connected",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch activity",
		})
	}

	// Calculate totals
	var totalActivities, totalPushes, totalPulls, totalBuilds int
	for _, a := range activities {
		totalActivities += a.TotalCount
		totalPushes += a.Pushes
		totalPulls += a.Pulls
		totalBuilds += a.Builds
	}

	c.Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	return c.JSON(fiber.Map{
		"username": username,
		"days":     days,
		"totals": fiber.Map{
			"activities": totalActivities,
			"pushes":     totalPushes,
			"pulls":      totalPulls,
			"builds":     totalBuilds,
		},
		"activity": activities,
	})
}

// GetProfilePage returns profile data for public profile page
func (h *HeatmapHandler) GetProfilePage(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username is required",
		})
	}

	// Get user by Docker username
	account, err := h.dockerService.GetDockerAccountByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get user info
	user, err := services.GetUserByID(account.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Check if profile is public
	if !user.PublicProfile {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Profile is private",
		})
	}

	// Get activity summary
	activities, _ := h.dockerService.GetActivitySummary(username, 365)

	var totalActivities int
	for _, a := range activities {
		totalActivities += a.TotalCount
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"github_username": user.GitHubUsername,
			"name":            user.Name,
			"avatar_url":      user.AvatarURL,
			"bio":             user.Bio,
		},
		"docker": fiber.Map{
			"username":     account.DockerUsername,
			"last_sync_at": account.LastSyncAt,
		},
		"stats": fiber.Map{
			"total_activities": totalActivities,
		},
		"available_themes": services.GetAvailableThemes(),
	})
}
