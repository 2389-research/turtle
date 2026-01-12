// ABOUTME: Main Bubble Tea application model for Turtle
// ABOUTME: Handles view routing, state management, and keyboard input

package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/2389-research/turtle/internal/progress"
	"github.com/2389-research/turtle/internal/skills"
)

// View represents the current screen
type View int

const (
	ViewHome View = iota
	ViewLesson
	ViewSkillTree
	ViewStats
	ViewSettings
)

// Model is the main application state
type Model struct {
	CurrentView View
	Width       int
	Height      int
	Progress    *skills.UserProgress
	SkillGraph  *skills.SkillGraph
	MenuIndex   int
	LessonModel *LessonModel
	SavePath    string
	quitting    bool
}

// NewModel creates the initial application state
func NewModel() Model {
	graph := buildDefaultSkillGraph()
	savePath := progress.GetDefaultPath()

	// Load existing progress or start fresh
	userProgress, err := progress.Load(savePath)
	if err != nil {
		// If load fails, start fresh
		userProgress = skills.NewUserProgress()
	}

	return Model{
		CurrentView: ViewHome,
		Width:       80,
		Height:      24,
		Progress:    userProgress,
		SkillGraph:  graph,
		MenuIndex:   0,
		SavePath:    savePath,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// save persists the current progress to disk
func (m *Model) save() {
	if m.SavePath != "" {
		_ = progress.Save(m.Progress, m.SavePath)
	}
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeypress(msg)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	}

	// Delegate to sub-views
	switch m.CurrentView {
	case ViewLesson:
		if m.LessonModel != nil {
			newLesson, cmd := m.LessonModel.Update(msg)
			m.LessonModel = newLesson.(*LessonModel)
			return m, cmd
		}
	}

	return m, nil
}

// handleKeypress processes keyboard input
func (m Model) handleKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys
	switch msg.String() {
	case "ctrl+c", "q":
		if m.CurrentView == ViewHome {
			m.save() // Save progress before quitting
			m.quitting = true
			return m, tea.Quit
		}
		// Return to home from other views
		m.CurrentView = ViewHome
		return m, nil
	}

	// View-specific keys
	switch m.CurrentView {
	case ViewHome:
		return m.handleHomeKeys(msg)
	case ViewLesson:
		if m.LessonModel != nil {
			newLesson, cmd := m.LessonModel.Update(msg)
			m.LessonModel = newLesson.(*LessonModel)

			// Check if lesson is done
			if m.LessonModel.Done {
				m.save() // Save progress after lesson
				m.CurrentView = ViewHome
				m.LessonModel = nil
			}
			return m, cmd
		}
	case ViewSkillTree:
		return m.handleSkillTreeKeys(msg)
	case ViewStats:
		return m.handleStatsKeys(msg)
	}

	return m, nil
}

// handleHomeKeys handles input on the home screen
func (m Model) handleHomeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	menuItems := []string{"Start Practice", "Skill Tree", "Stats", "Quit"}

	switch msg.String() {
	case "up", "k":
		m.MenuIndex--
		if m.MenuIndex < 0 {
			m.MenuIndex = len(menuItems) - 1
		}
	case "down", "j":
		m.MenuIndex++
		if m.MenuIndex >= len(menuItems) {
			m.MenuIndex = 0
		}
	case "enter", " ":
		switch m.MenuIndex {
		case 0: // Start Practice
			m.CurrentView = ViewLesson
			m.LessonModel = NewLessonModel(m.Progress, m.SkillGraph)
		case 1: // Skill Tree
			m.CurrentView = ViewSkillTree
		case 2: // Stats
			m.CurrentView = ViewStats
		case 3: // Quit
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

// handleSkillTreeKeys handles input on skill tree view
func (m Model) handleSkillTreeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "backspace":
		m.CurrentView = ViewHome
	}
	return m, nil
}

// handleStatsKeys handles input on stats view
func (m Model) handleStatsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "backspace":
		m.CurrentView = ViewHome
	}
	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	if m.quitting {
		return "\n  🐢 See you next time! Keep practicing!\n\n"
	}

	switch m.CurrentView {
	case ViewHome:
		return m.renderHome()
	case ViewLesson:
		if m.LessonModel != nil {
			return m.LessonModel.View()
		}
		return "Loading lesson..."
	case ViewSkillTree:
		return m.renderSkillTree()
	case ViewStats:
		return m.renderStats()
	default:
		return "Unknown view"
	}
}

// renderHome renders the main menu
func (m Model) renderHome() string {
	// Header with stats
	header := m.renderHeader()

	// Menu
	menuItems := []string{"Start Practice", "Skill Tree", "Stats", "Quit"}
	menu := "\n"
	for i, item := range menuItems {
		cursor := "  "
		style := MenuItemStyle
		if i == m.MenuIndex {
			cursor = "> "
			style = MenuItemSelectedStyle
		}
		menu += style.Render(cursor+item) + "\n"
	}

	// Today's goals
	goals := m.renderTodaysGoals()

	// Combine
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		BoxStyle.Render(menu),
		"",
		goals,
	)

	footer := FooterStyle.Render("↑/↓ navigate • enter select • q quit")

	return lipgloss.JoinVertical(lipgloss.Left, content, "", footer)
}

// renderHeader shows XP, level, streak
func (m Model) renderHeader() string {
	turtle := TitleStyle.Render("🐢 TURTLE")

	level := fmt.Sprintf("Level %d", m.Progress.Level)
	xp := fmt.Sprintf("%d XP", m.Progress.XP)
	streak := fmt.Sprintf("🔥 %d day streak", m.Progress.CurrentStreak)

	stats := lipgloss.JoinHorizontal(
		lipgloss.Center,
		AccentStyle.Render(level),
		MutedStyle.Render(" • "),
		XPStyle.Render(xp),
		MutedStyle.Render(" • "),
		StreakStyle.Render(streak),
	)

	// XP progress to next level
	xpInLevel := m.Progress.XP % skills.XPPerLevel
	xpProgress := float64(xpInLevel) / float64(skills.XPPerLevel)
	progressBar := ProgressBar(30, xpProgress)
	nextLevel := fmt.Sprintf(" %d/%d to Level %d", xpInLevel, skills.XPPerLevel, m.Progress.Level+1)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		turtle,
		stats,
		progressBar+MutedStyle.Render(nextLevel),
	)
}

// renderTodaysGoals shows what the user should practice today
func (m Model) renderTodaysGoals() string {
	title := SubtitleStyle.Render("TODAY'S TRAINING")

	// Count skills needing review
	dueSkills := m.SkillGraph.GetDueSkills(m.Progress)
	reviewCount := len(dueSkills)

	// Find next unlockable skill
	unlocked := m.SkillGraph.GetUnlockedSkills(m.Progress)
	var newSkill string
	for _, id := range unlocked {
		if m.Progress.GetStrength(id) == 0 {
			if skill, ok := m.SkillGraph.Skills[id]; ok {
				newSkill = skill.Name
			}
			break
		}
	}

	// Build goals list
	goals := ""
	if reviewCount > 0 {
		goals += fmt.Sprintf("  🔄 Review: %d skills need practice\n", reviewCount)
	}
	if newSkill != "" {
		goals += fmt.Sprintf("  📚 New: Learn %s\n", newSkill)
	}
	if reviewCount == 0 && newSkill == "" {
		goals += "  ✨ All caught up! Great work!\n"
	}

	return HighlightBoxStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		goals,
	))
}

// renderSkillTree shows the skill progression tree
func (m Model) renderSkillTree() string {
	title := TitleStyle.Render("🌳 SKILL TREE")

	categories := []skills.Category{
		skills.CategoryNavigation,
		skills.CategoryFileOps,
		skills.CategoryTmuxBasics,
		skills.CategoryTmuxPanes,
		skills.CategoryTmuxWindows,
	}

	content := ""
	for _, cat := range categories {
		catSkills := m.SkillGraph.GetSkillsByCategory(cat)
		if len(catSkills) == 0 {
			continue
		}

		content += "\n" + SubtitleStyle.Render(string(cat)) + "\n"

		for _, skillID := range catSkills {
			skill := m.SkillGraph.Skills[skillID]
			strength := m.Progress.GetStrength(skillID)

			// Status indicator
			var status string
			if strength == 0 {
				status = MutedStyle.Render("○")
			} else if strength < skills.CrackThreshold {
				status = DangerStyle.Render("◐") // cracking
			} else if strength < skills.UnlockThreshold {
				status = AccentStyle.Render("◑")
			} else {
				status = SuccessStyle.Render("●")
			}

			// Progress bar
			bar := ProgressBar(10, strength)

			line := fmt.Sprintf("  %s %s %s %.0f%%",
				status,
				skill.Name,
				bar,
				strength*100,
			)
			content += line + "\n"
		}
	}

	footer := FooterStyle.Render("esc back")

	return lipgloss.JoinVertical(lipgloss.Left, title, content, "", footer)
}

// renderStats shows detailed user statistics
func (m Model) renderStats() string {
	title := TitleStyle.Render("📊 STATISTICS")

	stats := fmt.Sprintf(`
  %s %s
  %s %s
  %s %s
  %s %s
`,
		StatLabelStyle.Render("Level:"),
		StatValueStyle.Render(fmt.Sprintf("%d", m.Progress.Level)),
		StatLabelStyle.Render("Total XP:"),
		XPStyle.Render(fmt.Sprintf("%d", m.Progress.XP)),
		StatLabelStyle.Render("Streak:"),
		StreakStyle.Render(fmt.Sprintf("%d days", m.Progress.CurrentStreak)),
		StatLabelStyle.Render("Best Streak:"),
		StatValueStyle.Render(fmt.Sprintf("%d days", m.Progress.BestStreak)),
	)

	// Skill summary
	total := len(m.SkillGraph.Skills)
	practiced := 0
	mastered := 0
	for skillID := range m.SkillGraph.Skills {
		strength := m.Progress.GetStrength(skillID)
		if strength > 0 {
			practiced++
		}
		if strength >= skills.UnlockThreshold {
			mastered++
		}
	}

	skillStats := fmt.Sprintf(`
  %s %s
  %s %s
  %s %s
`,
		StatLabelStyle.Render("Skills:"),
		StatValueStyle.Render(fmt.Sprintf("%d total", total)),
		StatLabelStyle.Render("Practiced:"),
		StatValueStyle.Render(fmt.Sprintf("%d", practiced)),
		StatLabelStyle.Render("Mastered:"),
		SuccessStyle.Render(fmt.Sprintf("%d", mastered)),
	)

	footer := FooterStyle.Render("esc back")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		BoxStyle.Render(stats),
		BoxStyle.Render(skillStats),
		"",
		footer,
	)
}
