// ABOUTME: Lipgloss styles for the Turtle TUI
// ABOUTME: Defines colors, borders, and visual theming for the learning app

package tui

import "github.com/charmbracelet/lipgloss"

// Colors - turtle-inspired palette.
var (
	ColorPrimary   = lipgloss.Color("#4CAF50") // Turtle green
	ColorSecondary = lipgloss.Color("#8BC34A") // Light green
	ColorAccent    = lipgloss.Color("#FFC107") // Golden (shell)
	ColorDanger    = lipgloss.Color("#F44336") // Red for errors/hearts
	ColorMuted     = lipgloss.Color("#9E9E9E") // Gray
	ColorBg        = lipgloss.Color("#1E1E1E") // Dark background
	ColorBgLight   = lipgloss.Color("#2D2D2D") // Slightly lighter bg
)

// Base styles.
var (
	BaseStyle = lipgloss.NewStyle().
			Background(ColorBg)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary)

	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	AccentStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	DangerStyle = lipgloss.NewStyle().
			Foreground(ColorDanger)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary)
)

// Component styles.
var (
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorMuted).
			Padding(1, 2)

	HighlightBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(1, 2)

	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	MenuItemSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(ColorAccent).
				Bold(true)

	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary)

	ProgressBarEmptyStyle = lipgloss.NewStyle().
				Foreground(ColorMuted)
)

// Header/footer.
var (
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Background(ColorBgLight).
			Padding(0, 1).
			Width(80)

	FooterStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Background(ColorBgLight).
			Padding(0, 1).
			Width(80)
)

// Stats display.
var (
	StatLabelStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Width(12)

	StatValueStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true)

	XPStyle = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true)

	StreakStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF9800")).
			Bold(true)

	HeartStyle = lipgloss.NewStyle().
			Foreground(ColorDanger)
)

// Terminal simulator styles.
var (
	TerminalStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorMuted).
			Padding(0, 1).
			Background(lipgloss.Color("#0D0D0D"))

	PromptStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	CommandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	OutputStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	ErrorOutputStyle = lipgloss.NewStyle().
				Foreground(ColorDanger)
)

// Helper to create a progress bar.
func ProgressBar(width int, percent float64) string {
	filled := int(float64(width) * percent)
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += ProgressBarStyle.Render("█")
		} else {
			bar += ProgressBarEmptyStyle.Render("░")
		}
	}
	return bar
}
