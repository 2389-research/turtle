// ABOUTME: Lipgloss styles for the Turtle TUI - "Phosphor Dreams" theme
// ABOUTME: Retro amber CRT aesthetic with modern game UI flourishes

package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Colors - "Phosphor Dreams" palette.
// Inspired by warm CRT phosphor glow with modern neon accents.
var (
	// Primary amber tones for warm phosphor glow.
	ColorAmber     = lipgloss.Color("#FFBF00") // Bright amber
	ColorAmberDim  = lipgloss.Color("#CC9900") // Muted amber
	ColorAmberGlow = lipgloss.Color("#FFD54F") // Glowing amber
	ColorCopper    = lipgloss.Color("#B87333") // Deep copper

	// Accent colors for neon game UI.
	ColorCyan    = lipgloss.Color("#00FFFF") // Electric cyan (success/XP)
	ColorCyanDim = lipgloss.Color("#00B8B8") // Muted cyan
	ColorMagenta = lipgloss.Color("#FF00FF") // Hot magenta (combos)
	ColorCoral   = lipgloss.Color("#FF6B6B") // Soft coral (errors)
	ColorLime    = lipgloss.Color("#ADFF2F") // Lime green (mastered)

	// Neutrals for deep space feel.
	ColorBgDeep  = lipgloss.Color("#0A0A0F") // Near black with blue tint
	ColorBgMid   = lipgloss.Color("#151520") // Card backgrounds
	ColorBgLight = lipgloss.Color("#1E1E2E") // Elevated surfaces
	ColorMuted   = lipgloss.Color("#5C5C6E") // Muted text
	ColorDim     = lipgloss.Color("#3A3A4A") // Very dim
	ColorText    = lipgloss.Color("#E8E8F0") // Primary text
)

// Decorative characters for borders and UI.
var (
	// Box drawing characters.
	BorderTL = "╭"
	BorderTR = "╮"
	BorderBL = "╰"
	BorderBR = "╯"
	BorderH  = "─"
	BorderV  = "│"

	// Decorative UI elements.
	Bullet      = "◆"
	BulletEmpty = "◇"
	Arrow       = "▸"
	ArrowLeft   = "◂"
	Star        = "★"
	StarEmpty   = "☆"
	Heart       = "♥"
	HeartEmpty  = "♡"
	Lightning   = "⚡"
	Fire        = "🔥"
	Turtle      = "🐢"

	// Progress bar characters for gradient feel.
	BarFull  = "█"
	BarHigh  = "▓"
	BarMid   = "▒"
	BarLow   = "░"
	BarEmpty = "·"

	// Mastery level indicators.
	MasteryNone     = "○"
	MasteryLearning = "◐"
	MasteryHalf     = "◑"
	MasteryFull     = "●"
	MasteryLocked   = "◌"

	// Pre-made divider strings.
	DividerLight  = "┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈"
	DividerMedium = "────────────────────────────────────────"
	DividerHeavy  = "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	DividerDot    = "• • • • • • • • • • • • • • • • • • • •"
)

// Custom border style for that retro feel.
var RetroDoubleBorder = lipgloss.Border{
	Top:         "═",
	Bottom:      "═",
	Left:        "║",
	Right:       "║",
	TopLeft:     "╔",
	TopRight:    "╗",
	BottomLeft:  "╚",
	BottomRight: "╝",
}

var RetroSingleBorder = lipgloss.Border{
	Top:         "─",
	Bottom:      "─",
	Left:        "│",
	Right:       "│",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "╰",
	BottomRight: "╯",
}

// Base styles for the application.
var (
	BaseStyle = lipgloss.NewStyle().
			Background(ColorBgDeep)

	// TitleStyle renders main titles with glowing amber.
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAmber).
			MarginBottom(1)

	// SubtitleStyle renders section headers with underline.
	SubtitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAmberDim).
			BorderStyle(lipgloss.Border{Bottom: "─"}).
			BorderForeground(ColorDim).
			BorderBottom(true).
			MarginBottom(1)

	// MutedStyle renders helper text in dim italic.
	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true)

	// AccentStyle renders highlighted text in cyan.
	AccentStyle = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Bold(true)

	// DangerStyle renders error states in coral.
	DangerStyle = lipgloss.NewStyle().
			Foreground(ColorCoral).
			Bold(true)

	// SuccessStyle renders success states in lime.
	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorLime).
			Bold(true)

	// TextStyle renders primary content text.
	TextStyle = lipgloss.NewStyle().
			Foreground(ColorText)
)

// Component styles for UI elements.
var (
	// BoxStyle renders standard content containers.
	BoxStyle = lipgloss.NewStyle().
			Border(RetroSingleBorder).
			BorderForeground(ColorDim).
			Padding(1, 2).
			MarginTop(1)

	// HighlightBoxStyle renders highlighted/focused containers.
	HighlightBoxStyle = lipgloss.NewStyle().
				Border(RetroDoubleBorder).
				BorderForeground(ColorAmber).
				Padding(1, 2).
				MarginTop(1)

	// GlowBoxStyle renders glowing containers for important info.
	GlowBoxStyle = lipgloss.NewStyle().
			Border(RetroDoubleBorder).
			BorderForeground(ColorCyan).
			Padding(1, 2).
			MarginTop(1)

	// MenuItemStyle renders unselected menu items.
	MenuItemStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			PaddingLeft(2)

	// MenuItemSelectedStyle renders the currently selected menu item.
	MenuItemSelectedStyle = lipgloss.NewStyle().
				Foreground(ColorAmberGlow).
				Bold(true).
				PaddingLeft(1).
				SetString(Arrow + " ")

	// ProgressBarStyle renders filled progress segments.
	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(ColorCyan)

	// ProgressBarEmptyStyle renders empty progress segments.
	ProgressBarEmptyStyle = lipgloss.NewStyle().
				Foreground(ColorDim)
)

// Header/footer.
var (
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAmber).
			Background(ColorBgLight).
			Padding(0, 2).
			MarginBottom(1)

	FooterStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Background(ColorBgMid).
			Padding(0, 2).
			MarginTop(1).
			Italic(true)
)

// Stats display styles for game metrics.
var (
	// StatLabelStyle renders stat labels in muted text.
	StatLabelStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Width(14)

	// StatValueStyle renders stat values in bold.
	StatValueStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Bold(true)

	// XPStyle renders XP counter in electric cyan glow.
	XPStyle = lipgloss.NewStyle().
		Foreground(ColorCyan).
		Bold(true)

	// StreakStyle renders streak counter in hot magenta.
	StreakStyle = lipgloss.NewStyle().
			Foreground(ColorMagenta).
			Bold(true)

	// HeartStyle renders hearts in coral warmth.
	HeartStyle = lipgloss.NewStyle().
			Foreground(ColorCoral)

	// LevelStyle renders level indicator in amber glow.
	LevelStyle = lipgloss.NewStyle().
			Foreground(ColorAmberGlow).
			Bold(true)

	// ComboStyle renders combo counter with magenta fire effect.
	ComboStyle = lipgloss.NewStyle().
			Foreground(ColorMagenta).
			Bold(true).
			Blink(true)
)

// Terminal simulator styles for the embedded terminal UI.
var (
	// TerminalStyle renders the terminal container box.
	TerminalStyle = lipgloss.NewStyle().
			Border(lipgloss.Border{
			Top:         "▄",
			Bottom:      "▀",
			Left:        "▐",
			Right:       "▌",
			TopLeft:     "▗",
			TopRight:    "▖",
			BottomLeft:  "▝",
			BottomRight: "▘",
		}).
		BorderForeground(ColorCopper).
		Background(ColorBgDeep).
		Padding(0, 1)

	// PromptStyle renders the command prompt symbol.
	PromptStyle = lipgloss.NewStyle().
			Foreground(ColorLime).
			Bold(true)

	// CommandStyle renders user-typed commands.
	CommandStyle = lipgloss.NewStyle().
			Foreground(ColorAmberGlow)

	// OutputStyle renders command output text.
	OutputStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	// ErrorOutputStyle renders error output in coral italic.
	ErrorOutputStyle = lipgloss.NewStyle().
				Foreground(ColorCoral).
				Italic(true)

	// CursorStyle renders the blinking cursor.
	CursorStyle = lipgloss.NewStyle().
			Foreground(ColorAmber).
			Bold(true).
			Blink(true)
)

// Skill tree styles.
var (
	SkillLockedStyle = lipgloss.NewStyle().
				Foreground(ColorDim)

	SkillLearningStyle = lipgloss.NewStyle().
				Foreground(ColorAmberDim)

	SkillPracticingStyle = lipgloss.NewStyle().
				Foreground(ColorAmber)

	SkillMasteredStyle = lipgloss.NewStyle().
				Foreground(ColorLime).
				Bold(true)

	CategoryStyle = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Bold(true).
			MarginTop(1).
			MarginBottom(0)
)

// Badge/achievement styles.
var (
	BadgeStyle = lipgloss.NewStyle().
			Foreground(ColorAmberGlow).
			Background(ColorBgLight).
			Padding(0, 1).
			Bold(true)

	BadgeLockedStyle = lipgloss.NewStyle().
				Foreground(ColorDim).
				Background(ColorBgMid).
				Padding(0, 1)
)

// ProgressBar creates a gradient-feel progress bar.
func ProgressBar(width int, percent float64) string {
	if percent > 1.0 {
		percent = 1.0
	}
	if percent < 0 {
		percent = 0
	}

	filled := int(float64(width) * percent)
	var bar strings.Builder

	for i := 0; i < width; i++ {
		switch {
		case i < filled:
			// Gradient effect: brighter toward the leading edge.
			if i >= filled-2 && percent < 1.0 {
				bar.WriteString(ProgressBarStyle.Render(BarHigh))
			} else {
				bar.WriteString(ProgressBarStyle.Render(BarFull))
			}
		case i == filled && percent > 0 && percent < 1.0:
			// Transitional character at the edge.
			bar.WriteString(lipgloss.NewStyle().Foreground(ColorCyanDim).Render(BarMid))
		default:
			bar.WriteString(ProgressBarEmptyStyle.Render(BarEmpty))
		}
	}
	return bar.String()
}

// ProgressBarWithGlow creates a progress bar with percentage glow effect.
func ProgressBarWithGlow(width int, percent float64) string {
	bar := ProgressBar(width, percent)
	pctStr := lipgloss.NewStyle().
		Foreground(ColorCyan).
		Bold(true).
		Render(strings.Repeat(" ", 1) + formatPercent(percent))
	return bar + pctStr
}

func formatPercent(p float64) string {
	pct := int(p * 100)
	if pct > 100 {
		pct = 100
	}
	return lipgloss.NewStyle().Foreground(ColorCyan).Render(
		strings.Repeat(" ", 3-len(string(rune('0'+pct/100)))) + string(rune('0'+pct/100)) + string(rune('0'+(pct%100)/10)) + string(rune('0'+pct%10)) + "%",
	)
}

// Hearts renders heart indicators for lives.
func Hearts(current, max int) string {
	var result strings.Builder
	for i := 0; i < max; i++ {
		if i < current {
			result.WriteString(HeartStyle.Render(Heart + " "))
		} else {
			result.WriteString(lipgloss.NewStyle().Foreground(ColorDim).Render(HeartEmpty + " "))
		}
	}
	return result.String()
}

// Stars renders star rating.
func Stars(current, max int) string {
	var result strings.Builder
	for i := 0; i < max; i++ {
		if i < current {
			result.WriteString(AccentStyle.Render(Star))
		} else {
			result.WriteString(lipgloss.NewStyle().Foreground(ColorDim).Render(StarEmpty))
		}
	}
	return result.String()
}

// MasteryIndicator returns the appropriate mastery symbol.
func MasteryIndicator(strength float64, unlocked bool) string {
	if !unlocked {
		return SkillLockedStyle.Render(MasteryLocked)
	}
	switch {
	case strength >= 0.8:
		return SkillMasteredStyle.Render(MasteryFull)
	case strength >= 0.5:
		return SkillPracticingStyle.Render(MasteryHalf)
	case strength > 0:
		return SkillLearningStyle.Render(MasteryLearning)
	default:
		return SkillLockedStyle.Render(MasteryNone)
	}
}

// Divider creates a styled divider line.
func Divider(width int, style string) string {
	var char string
	var color lipgloss.Color

	switch style {
	case "heavy":
		char = "━"
		color = ColorAmber
	case "double":
		char = "═"
		color = ColorAmberDim
	case "dot":
		char = "•"
		color = ColorDim
	case "light":
		char = "┈"
		color = ColorDim
	default:
		char = "─"
		color = ColorDim
	}

	return lipgloss.NewStyle().Foreground(color).Render(strings.Repeat(char, width))
}

// Logo returns the stylized turtle logo.
func Logo() string {
	logo := `
   ╭──────────────────────────╮
   │  ` + TitleStyle.Render("🐢 T U R T L E") + `  │
   │  ` + MutedStyle.Render("terminal trainer") + `       │
   ╰──────────────────────────╯`
	return lipgloss.NewStyle().Foreground(ColorAmber).Render(logo)
}

// SmallLogo returns a compact logo.
func SmallLogo() string {
	return TitleStyle.Render("🐢 TURTLE")
}

// HeroBanner returns a dramatic ASCII art banner for the landing screen.
func HeroBanner() string {
	// Stylized ASCII turtle with shell pattern.
	turtleArt := lipgloss.NewStyle().Foreground(ColorAmberGlow).Render(`
                    __
         .,-;-;-,. /'_\
       _/_/_/_|_\_\) /
     '-<_><_><_><_>=/\
       ` + "`" + `/_/____/_/` + "`" + `\_\
        _____|_|___|_|___
       (_)   (_)(_)   (_)`)

	// Big blocky title.
	title := lipgloss.NewStyle().Foreground(ColorAmber).Bold(true).Render(`
  ▀█▀ █ █ █▀█ ▀█▀ █   █▀▀
   █  █ █ █▀▄  █  █   █▀▀
   █  ▀▄▀ █ █  █  █▄▄ █▄▄`)

	return lipgloss.JoinVertical(lipgloss.Center, turtleArt, title)
}

// Tagline returns the app tagline with styling.
func Tagline() string {
	return lipgloss.NewStyle().
		Foreground(ColorCyan).
		Italic(true).
		Render("Master your terminal, one command at a time")
}

// WhatIsTurtle returns the explanation block with default width.
func WhatIsTurtle() string {
	return WhatIsTurtleWidth(46)
}

// WhatIsTurtleWidth returns the explanation block with custom width.
func WhatIsTurtleWidth(width int) string {
	header := lipgloss.NewStyle().
		Foreground(ColorAmberDim).
		Bold(true).
		Render("WHAT IS TURTLE?")

	desc := lipgloss.NewStyle().
		Foreground(ColorText).
		Width(width).
		Render("A Duolingo-style trainer for terminal and " +
			"tmux commands. Practice with spaced " +
			"repetition, earn XP, and level up!")

	bullets := []string{
		lipgloss.NewStyle().Foreground(ColorCyan).Render("◆") + " " +
			lipgloss.NewStyle().Foreground(ColorMuted).Render("Learn") + " " +
			lipgloss.NewStyle().Foreground(ColorText).Render("terminal navigation & files"),
		lipgloss.NewStyle().Foreground(ColorMagenta).Render("◆") + " " +
			lipgloss.NewStyle().Foreground(ColorMuted).Render("Practice") + " " +
			lipgloss.NewStyle().Foreground(ColorText).Render("with spaced repetition"),
		lipgloss.NewStyle().Foreground(ColorLime).Render("◆") + " " +
			lipgloss.NewStyle().Foreground(ColorMuted).Render("Master") + " " +
			lipgloss.NewStyle().Foreground(ColorText).Render("tmux panes and windows"),
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		desc,
		"",
		strings.Join(bullets, "\n"),
	)
}

// PlayerCard renders the player stats in a card format.
func PlayerCard(level int, xp int, streak int, xpProgress float64) string {
	// Level badge.
	levelBadge := lipgloss.NewStyle().
		Foreground(ColorBgDeep).
		Background(ColorAmber).
		Bold(true).
		Padding(0, 1).
		Render(fmt.Sprintf("LVL %d", level))

	// XP with glow.
	xpText := lipgloss.NewStyle().
		Foreground(ColorCyan).
		Bold(true).
		Render(fmt.Sprintf("%d XP", xp))

	// Streak.
	var streakText string
	if streak > 0 {
		streakText = lipgloss.NewStyle().
			Foreground(ColorMagenta).
			Bold(true).
			Render(fmt.Sprintf("🔥 %d day streak", streak))
	} else {
		streakText = lipgloss.NewStyle().
			Foreground(ColorDim).
			Italic(true).
			Render("start your streak today!")
	}

	// Progress bar to next level.
	bar := ProgressBar(20, xpProgress)
	barLabel := lipgloss.NewStyle().
		Foreground(ColorDim).
		Render(fmt.Sprintf(" %.0f%% to next level", xpProgress*100))

	// Combine into card.
	statsRow := lipgloss.JoinHorizontal(
		lipgloss.Center,
		levelBadge,
		"  ",
		xpText,
		"  "+Bullet+"  ",
		streakText,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		statsRow,
		bar+barLabel,
	)
}

// LevelNames returns display names for mission levels.
func LevelNames() []string {
	return []string{
		"Shell Basics",
		"Navigation",
		"File Operations",
		"Text Processing",
		"Advanced",
	}
}

// LevelGoals returns goal descriptions for mission levels.
func LevelGoals() []string {
	return []string{
		"Master the fundamentals",
		"Move around with confidence",
		"Create, copy, move files",
		"Search and manipulate text",
		"Power user techniques",
	}
}
