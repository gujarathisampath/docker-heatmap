package services

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"strings"
	"time"

	"docker-heatmap/internal/models"
)

type HeatmapService struct {
	dockerService *DockerHubService
}

func NewHeatmapService() *HeatmapService {
	return &HeatmapService{
		dockerService: NewDockerHubService(),
	}
}

// SVGOptions represents customizable options for the SVG heatmap
type SVGOptions struct {
	Theme       string // Theme name or "custom"
	CellSize    int    // Size of each cell (default 11)
	CellRadius  int    // Border radius of cells (default 2)
	Days        int    // Number of days to show (default 365)
	HideLegend  bool   // Hide the legend
	HideTotal   bool   // Hide total count
	HideLabels  bool   // Hide month/day labels
	FontFamily  string // Custom font family
	CustomTitle string // Custom title instead of default

	// Custom colors (when theme is "custom")
	BgColor      string   // Background color
	TextColor    string   // Text color
	CustomColors []string // Level 0-4 colors
}

// Theme represents a color theme for the heatmap
type Theme struct {
	Name      string
	BgColor   string
	TextColor string
	Colors    []string // Level 0-4 colors
}

var Themes = map[string]Theme{
	// Primary themes
	"github": {
		Name:      "GitHub Dark",
		BgColor:   "transparent",
		TextColor: "#8b949e",
		Colors:    []string{"#161b22", "#0e4429", "#006d32", "#26a641", "#39d353"},
	},
	"github-light": {
		Name:      "GitHub Light",
		BgColor:   "#ffffff",
		TextColor: "#57606a",
		Colors:    []string{"#ebedf0", "#9be9a8", "#40c463", "#30a14e", "#216e39"},
	},
	"docker": {
		Name:      "Docker",
		BgColor:   "transparent",
		TextColor: "#0db7ed",
		Colors:    []string{"#1a2634", "#1a4971", "#1d6fa5", "#2496ed", "#6db3f2"},
	},

	// Popular editor themes
	"dracula": {
		Name:      "Dracula",
		BgColor:   "#282a36",
		TextColor: "#f8f8f2",
		Colors:    []string{"#44475a", "#6272a4", "#bd93f9", "#ff79c6", "#50fa7b"},
	},
	"nord": {
		Name:      "Nord",
		BgColor:   "transparent",
		TextColor: "#d8dee9",
		Colors:    []string{"#2e3440", "#3b4252", "#5e81ac", "#81a1c1", "#88c0d0"},
	},
	"monokai": {
		Name:      "Monokai",
		BgColor:   "transparent",
		TextColor: "#f8f8f2",
		Colors:    []string{"#272822", "#49483e", "#a6e22e", "#e6db74", "#f92672"},
	},
	"one-dark": {
		Name:      "One Dark",
		BgColor:   "transparent",
		TextColor: "#abb2bf",
		Colors:    []string{"#282c34", "#3e4451", "#61afef", "#98c379", "#e5c07b"},
	},
	"tokyo-night": {
		Name:      "Tokyo Night",
		BgColor:   "transparent",
		TextColor: "#a9b1d6",
		Colors:    []string{"#1a1b26", "#24283b", "#7aa2f7", "#bb9af7", "#73daca"},
	},
	"catppuccin": {
		Name:      "Catppuccin",
		BgColor:   "transparent",
		TextColor: "#cdd6f4",
		Colors:    []string{"#1e1e2e", "#313244", "#89b4fa", "#a6e3a1", "#f5c2e7"},
	},

	// Color themes
	"ocean": {
		Name:      "Ocean",
		BgColor:   "transparent",
		TextColor: "#6b8fa3",
		Colors:    []string{"#1a2332", "#1e4976", "#2171b5", "#4292c6", "#6baed6"},
	},
	"sunset": {
		Name:      "Sunset",
		BgColor:   "transparent",
		TextColor: "#b38867",
		Colors:    []string{"#2d1f1f", "#6b3030", "#b54040", "#e06050", "#ff8c66"},
	},
	"forest": {
		Name:      "Forest",
		BgColor:   "transparent",
		TextColor: "#7d9c7d",
		Colors:    []string{"#1a2e1a", "#2d4a2d", "#3d6b3d", "#4d8c4d", "#5dac5d"},
	},
	"purple": {
		Name:      "Purple",
		BgColor:   "transparent",
		TextColor: "#9d8abf",
		Colors:    []string{"#1a1a2e", "#2d2d5a", "#6b3fa0", "#9d4edd", "#c77dff"},
	},
	"rose": {
		Name:      "Rose",
		BgColor:   "transparent",
		TextColor: "#bf8a9d",
		Colors:    []string{"#2e1a24", "#5a2d42", "#a03f6b", "#dd4e9d", "#ff7dc7"},
	},

	// Minimal/Grayscale
	"minimal": {
		Name:      "Minimal",
		BgColor:   "transparent",
		TextColor: "#666666",
		Colors:    []string{"#f0f0f0", "#d4d4d4", "#a8a8a8", "#6b6b6b", "#333333"},
	},
	"minimal-dark": {
		Name:      "Minimal Dark",
		BgColor:   "transparent",
		TextColor: "#999999",
		Colors:    []string{"#1a1a1a", "#333333", "#4d4d4d", "#808080", "#b3b3b3"},
	},
}

type HeatmapConfig struct {
	CellSize   int
	CellMargin int
	CellRadius int
	Rows       int // Always 7 for days of week
	FontSize   int
	Colors     []string
	TextColor  string
	BgColor    string
	FontFamily string
}

// SVGData represents the data needed to render the SVG
type SVGData struct {
	Width        int
	Height       int
	Cells        []Cell
	MonthLabels  []MonthLabel
	DayLabels    []DayLabel
	Config       HeatmapConfig
	Username     string
	TotalCount   int
	HideLegend   bool
	HideTotal    bool
	HideLabels   bool
	CustomTitle  string
	LegendX      int
	LegendY      int
	FooterY      int
	CellsOffsetX int
}

type Cell struct {
	X      int
	Y      int
	Width  int
	Height int
	Radius int
	Color  string
	Date   string
	Count  int
}

type MonthLabel struct {
	X     int
	Y     int
	Label string
}

type DayLabel struct {
	X     int
	Y     int
	Label string
}

const svgTemplate = `<svg width="100%" height="auto" viewBox="0 0 {{.Width}} {{.Height}}" preserveAspectRatio="xMidYMid meet" xmlns="http://www.w3.org/2000/svg">
  <style>
    .day { shape-rendering: geometricPrecision; outline: 1px solid rgba(27, 31, 35, 0.06); outline-offset: -1px; }
    .month-label { font-size: {{.Config.FontSize}}px; fill: {{.Config.TextColor}}; font-family: {{.Config.FontFamily}}; }
    .day-label { font-size: 9px; fill: {{.Config.TextColor}}; font-family: {{.Config.FontFamily}}; }
    .title { font-size: 11px; fill: {{.Config.TextColor}}; font-family: {{.Config.FontFamily}}; font-weight: 600; }
    .legend-label { font-size: 9px; fill: {{.Config.TextColor}}; font-family: {{.Config.FontFamily}}; }
  </style>
  <rect width="{{.Width}}" height="{{.Height}}" fill="{{.Config.BgColor}}" rx="6"/>
  {{if not .HideLabels}}
  <!-- Month labels -->
  {{range .MonthLabels}}
  <text x="{{.X}}" y="{{.Y}}" class="month-label">{{.Label}}</text>
  {{end}}
  
  <!-- Day labels -->
  {{range .DayLabels}}
  <text x="{{.X}}" y="{{.Y}}" class="day-label">{{.Label}}</text>
  {{end}}
  {{end}}
  
  <!-- Activity cells -->
  <g transform="translate({{.CellsOffsetX}}, 25)">
    {{range .Cells}}
    <rect class="day" x="{{.X}}" y="{{.Y}}" width="{{.Width}}" height="{{.Height}}" fill="{{.Color}}" rx="{{.Radius}}">
      <title>{{.Date}}: {{.Count}} activities</title>
    </rect>
    {{end}}
  </g>
  {{if not .HideTotal}}
  <!-- Footer -->
  <text x="{{.CellsOffsetX}}" y="{{.FooterY}}" class="title">{{if .CustomTitle}}{{.CustomTitle}}{{else}}@{{.Username}} Docker Activity • {{.TotalCount}} total{{end}}</text>
  {{end}}
  {{if not .HideLegend}}
  <!-- Legend -->
  <g transform="translate({{.LegendX}}, {{.LegendY}})">
    <text x="-25" y="10" class="legend-label">Less</text>
    {{range $i, $color := .Config.Colors}}
    <rect x="{{multiply $i 14}}" y="0" width="11" height="11" fill="{{$color}}" rx="2"/>
    {{end}}
    <text x="75" y="10" class="legend-label">More</text>
  </g>
  {{end}}
</svg>`

// GenerateSVG generates an SVG heatmap with default options
func (s *HeatmapService) GenerateSVG(dockerUsername string, days int) ([]byte, error) {
	return s.GenerateSVGWithOptions(dockerUsername, SVGOptions{
		Theme: "github",
		Days:  days,
	})
}

// GenerateSVGWithOptions generates an SVG heatmap with custom options
func (s *HeatmapService) GenerateSVGWithOptions(dockerUsername string, opts SVGOptions) ([]byte, error) {
	// Set defaults
	if opts.Days <= 0 || opts.Days > 365 {
		opts.Days = 365
	}
	if opts.CellSize <= 0 {
		opts.CellSize = 11
	}
	if opts.CellSize > 20 {
		opts.CellSize = 20
	}
	if opts.CellRadius < 0 {
		opts.CellRadius = 2
	}
	if opts.Theme == "" {
		opts.Theme = "github"
	}
	if opts.FontFamily == "" {
		opts.FontFamily = "-apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif"
	}

	// Get theme or use custom colors
	var bgColor, textColor string
	var colors []string

	if opts.Theme == "custom" && len(opts.CustomColors) == 5 {
		bgColor = opts.BgColor
		if bgColor == "" {
			bgColor = "transparent"
		}
		textColor = opts.TextColor
		if textColor == "" {
			textColor = "#8b949e"
		}
		colors = opts.CustomColors
	} else {
		theme, ok := Themes[opts.Theme]
		if !ok {
			theme = Themes["github"]
		}
		bgColor = theme.BgColor
		textColor = theme.TextColor
		colors = theme.Colors
	}

	// Get activity data
	activities, err := s.dockerService.GetActivitySummary(dockerUsername, opts.Days)
	if err != nil {
		return nil, err
	}

	// Calculate dimensions
	cellMargin := 3
	cellTotal := opts.CellSize + cellMargin
	numWeeks := (opts.Days + 6) / 7

	leftMargin := 40
	if opts.HideLabels {
		leftMargin = 10
	}

	// Calculate cells area dimensions
	cellsWidth := numWeeks * cellTotal
	cellsHeight := 7 * cellTotal

	// Calculate total width
	width := leftMargin + cellsWidth + 20

	// Calculate height based on what's shown
	topMargin := 25
	bottomMargin := 10
	if !opts.HideTotal || !opts.HideLegend {
		bottomMargin = 30
	}
	height := topMargin + cellsHeight + bottomMargin

	// Build config
	config := HeatmapConfig{
		CellSize:   opts.CellSize,
		CellMargin: cellMargin,
		CellRadius: opts.CellRadius,
		Rows:       7,
		FontSize:   10,
		Colors:     colors,
		TextColor:  textColor,
		BgColor:    bgColor,
		FontFamily: opts.FontFamily,
	}

	// Create cells
	cells := make([]Cell, 0, len(activities))
	totalCount := 0

	startDate := time.Now().AddDate(0, 0, -opts.Days+1)
	// Align to start of week (Sunday)
	for startDate.Weekday() != time.Sunday {
		startDate = startDate.AddDate(0, 0, -1)
	}

	activityMap := make(map[string]models.ActivitySummary)
	for _, a := range activities {
		activityMap[a.Date] = a
		totalCount += a.TotalCount
	}

	currentDate := startDate
	col := 0
	today := time.Now()
	for !currentDate.After(today) {
		row := int(currentDate.Weekday())
		dateStr := currentDate.Format("2006-01-02")

		activity := activityMap[dateStr]
		color := config.Colors[activity.Level]

		cells = append(cells, Cell{
			X:      col * cellTotal,
			Y:      row * cellTotal,
			Width:  opts.CellSize,
			Height: opts.CellSize,
			Radius: opts.CellRadius,
			Color:  color,
			Date:   currentDate.Format("Jan 2, 2006"),
			Count:  activity.TotalCount,
		})

		if currentDate.Weekday() == time.Saturday {
			col++
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Create month labels
	monthLabels := make([]MonthLabel, 0)
	if !opts.HideLabels {
		currentMonth := startDate.Month()
		for i := 0; i < numWeeks; i++ {
			checkDate := startDate.AddDate(0, 0, i*7)
			if checkDate.Month() != currentMonth || i == 0 {
				currentMonth = checkDate.Month()
				monthLabels = append(monthLabels, MonthLabel{
					X:     leftMargin + (i * cellTotal),
					Y:     15,
					Label: checkDate.Format("Jan"),
				})
			}
		}
	}

	// Create day labels
	var dayLabels []DayLabel
	if !opts.HideLabels {
		dayLabels = []DayLabel{
			{X: 5, Y: 25 + (1 * cellTotal) + 8, Label: "Mon"},
			{X: 5, Y: 25 + (3 * cellTotal) + 8, Label: "Wed"},
			{X: 5, Y: 25 + (5 * cellTotal) + 8, Label: "Fri"},
		}
	}

	// Calculate footer and legend positions
	footerY := topMargin + cellsHeight + 18
	legendY := topMargin + cellsHeight + 5
	legendX := width - 120

	// Security: Escape user-provided content to prevent XSS in SVG
	safeUsername := html.EscapeString(dockerUsername)
	safeCustomTitle := html.EscapeString(opts.CustomTitle)

	data := SVGData{
		Width:        width,
		Height:       height,
		Cells:        cells,
		MonthLabels:  monthLabels,
		DayLabels:    dayLabels,
		Config:       config,
		Username:     safeUsername,
		TotalCount:   totalCount,
		HideLegend:   opts.HideLegend,
		HideTotal:    opts.HideTotal,
		HideLabels:   opts.HideLabels,
		CustomTitle:  safeCustomTitle,
		LegendX:      legendX,
		LegendY:      legendY,
		FooterY:      footerY,
		CellsOffsetX: leftMargin,
	}

	// Create template with helper functions
	funcMap := template.FuncMap{
		"subtract": func(a, b int) int { return a - b },
		"multiply": func(a, b int) int { return a * b },
	}

	tmpl, err := template.New("heatmap").Funcs(funcMap).Parse(svgTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

// GetAvailableThemes returns all available theme names
func GetAvailableThemes() []string {
	themes := make([]string, 0, len(Themes))
	for name := range Themes {
		themes = append(themes, name)
	}
	return themes
}

// ParseSVGOptionsFromQuery parses SVG options from query parameters
func ParseSVGOptionsFromQuery(params map[string]string) SVGOptions {
	opts := SVGOptions{
		Theme:      "github",
		Days:       365,
		CellSize:   11,
		CellRadius: 2,
	}

	if v, ok := params["theme"]; ok {
		opts.Theme = strings.ToLower(v)
	}
	if v, ok := params["days"]; ok {
		fmt.Sscanf(v, "%d", &opts.Days)
	}
	if v, ok := params["cell_size"]; ok {
		fmt.Sscanf(v, "%d", &opts.CellSize)
	}
	if v, ok := params["radius"]; ok {
		fmt.Sscanf(v, "%d", &opts.CellRadius)
	}
	if v, ok := params["hide_legend"]; ok && (v == "true" || v == "1") {
		opts.HideLegend = true
	}
	if v, ok := params["hide_total"]; ok && (v == "true" || v == "1") {
		opts.HideTotal = true
	}
	if v, ok := params["hide_labels"]; ok && (v == "true" || v == "1") {
		opts.HideLabels = true
	}
	if v, ok := params["title"]; ok {
		opts.CustomTitle = v
	}

	// Custom colors support
	if v, ok := params["bg_color"]; ok {
		opts.BgColor = v
	}
	if v, ok := params["text_color"]; ok {
		opts.TextColor = v
	}
	// Custom level colors: color0, color1, color2, color3, color4
	customColors := make([]string, 0, 5)
	for i := 0; i < 5; i++ {
		if v, ok := params[fmt.Sprintf("color%d", i)]; ok {
			customColors = append(customColors, v)
		}
	}
	if len(customColors) == 5 {
		opts.CustomColors = customColors
		opts.Theme = "custom"
	}

	return opts
}
