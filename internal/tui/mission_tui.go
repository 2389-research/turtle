// ABOUTME: Mission-based TUI for interactive terminal learning
// ABOUTME: Provides sandbox environment with goal-based challenges and immediate feedback

package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/2389-research/turtle/internal/sandbox"
	"github.com/2389-research/turtle/internal/skills"
)

// MissionScreen represents the current screen state.
type MissionScreen int

const (
	ScreenMenu MissionScreen = iota
	ScreenLevelSelect
	ScreenMission
	ScreenComplete
	ScreenStats
	ScreenFlashcards
)

// MissionTUI handles mission-based learning.
type MissionTUI struct {
	Screen   MissionScreen
	Missions map[int][]*sandbox.Mission

	// Mission state
	CurrentLevel   int
	CurrentMission int
	Runner         *sandbox.MissionRunner
	Input          string
	History        []historyEntry
	ShowHint       bool

	// Menu state
	MenuIndex  int
	LevelIndex int

	// Stats
	MissionsCompleted int
	CommandsUsed      int

	// Terminal dimensions
	Width  int
	Height int

	// Flashcard mode (shares Progress and SkillGraph).
	FlashcardModel Model
}

type historyEntry struct {
	Command string
	Output  string
	Error   string
	Success bool
}

// mainMenuItem defines a menu entry.
type mainMenuItem struct {
	label string
	desc  string
	icon  string
}

// mainMenuItems is the single source of truth for main menu options.
var mainMenuItems = []mainMenuItem{
	{"Missions", "Interactive sandbox tutorials", Lightning},
	{"Flashcards", "Quick spaced repetition drill", Star},
	{"Speed Round", "Timed flashcard challenge", Fire},
	{"Skill Tree", "View your skill progress", Bullet},
	{"Stats", "Session and overall stats", Bullet},
	{"Quit", "Exit Turtle", ArrowLeft},
}

// NewMissionTUI creates a new mission-based learning interface.
func NewMissionTUI() *MissionTUI {
	return &MissionTUI{
		Screen:         ScreenMenu,
		Missions:       sandbox.GetAllMissions(),
		MenuIndex:      0,
		Width:          80,
		Height:         24,
		FlashcardModel: NewModel(),
	}
}

// Init implements tea.Model.
func (m *MissionTUI) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *MissionTUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		// Pass window size to flashcard model too
		m.FlashcardModel.Width = msg.Width
		m.FlashcardModel.Height = msg.Height

	case tea.KeyMsg:
		switch m.Screen {
		case ScreenMenu:
			return m.updateMenu(msg)
		case ScreenLevelSelect:
			return m.updateLevelSelect(msg)
		case ScreenMission:
			return m.updateMission(msg)
		case ScreenComplete:
			return m.updateComplete(msg)
		case ScreenStats:
			return m.updateStats(msg)
		case ScreenFlashcards:
			return m.updateFlashcards(msg)
		}
	}
	return m, nil
}

func (m *MissionTUI) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	menuCount := len(mainMenuItems)

	switch msg.String() {
	case "up", "k":
		if m.MenuIndex > 0 {
			m.MenuIndex--
		} else {
			m.MenuIndex = menuCount - 1
		}
	case "down", "j":
		if m.MenuIndex < menuCount-1 {
			m.MenuIndex++
		} else {
			m.MenuIndex = 0
		}
	case "enter", " ":
		switch m.MenuIndex {
		case 0: // Missions
			m.Screen = ScreenLevelSelect
			m.LevelIndex = 0
		case 1: // Flashcards - go directly to practice
			m.Screen = ScreenFlashcards
			m.FlashcardModel.CurrentView = ViewLesson
			m.FlashcardModel.LessonModel = NewLessonModel(m.FlashcardModel.Progress, m.FlashcardModel.SkillGraph, nil)
			return m, m.FlashcardModel.LessonModel.Init()
		case 2: // Speed Round - go directly to speed round
			m.Screen = ScreenFlashcards
			m.FlashcardModel.CurrentView = ViewLesson
			m.FlashcardModel.LessonModel = NewSpeedRoundModel(m.FlashcardModel.Progress, m.FlashcardModel.SkillGraph, nil)
			return m, m.FlashcardModel.LessonModel.Init()
		case 3: // Skill Tree
			m.Screen = ScreenFlashcards
			m.FlashcardModel.CurrentView = ViewSkillTree
		case 4: // Stats
			m.Screen = ScreenStats
		case 5: // Quit
			return m, tea.Quit
		}
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m *MissionTUI) updateLevelSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	maxLevel := len(m.Missions) - 1

	switch msg.String() {
	case "up", "k":
		if m.LevelIndex > 0 {
			m.LevelIndex--
		}
	case "down", "j":
		if m.LevelIndex < maxLevel {
			m.LevelIndex++
		}
	case "enter", " ":
		m.startLevel(m.LevelIndex)
	case "q", "esc":
		m.Screen = ScreenMenu
	}
	return m, nil
}

func (m *MissionTUI) updateMission(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if m.Input != "" {
			m.executeCommand()
		}
	case "backspace":
		if len(m.Input) > 0 {
			m.Input = m.Input[:len(m.Input)-1]
		}
	case "ctrl+c":
		m.Input = ""
	case "ctrl+h", "?":
		m.ShowHint = !m.ShowHint
	case "ctrl+r":
		if m.Runner != nil {
			m.Runner.Reset()
			m.History = nil
			m.Input = ""
		}
	case "esc":
		m.Screen = ScreenLevelSelect
	default:
		// Regular character input
		switch {
		case len(msg.String()) == 1:
			m.Input += msg.String()
		case msg.String() == "space":
			m.Input += " "
		case msg.String() == "tab":
			m.Input += "\t" // Could implement tab completion here
		}
	}
	return m, nil
}

func (m *MissionTUI) updateComplete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ":
		m.nextMission()
	case "q", "esc":
		m.Screen = ScreenLevelSelect
	}
	return m, nil
}

func (m *MissionTUI) updateStats(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle ctrl+c to quit app.
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}
	// Any other key returns to menu.
	m.Screen = ScreenMenu
	return m, nil
}

func (m *MissionTUI) updateFlashcards(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Always intercept esc/q to return to main menu.
	if msg.String() == "esc" || msg.String() == "q" {
		m.Screen = ScreenMenu
		m.FlashcardModel.CurrentView = ViewHome
		m.FlashcardModel.LessonModel = nil
		return m, nil
	}

	// Delegate to flashcard model for other keys.
	updatedModel, cmd := m.FlashcardModel.Update(msg)
	if model, ok := updatedModel.(Model); ok {
		m.FlashcardModel = model
	}

	// If flashcard model navigated to its home, override back to main menu.
	if m.FlashcardModel.CurrentView == ViewHome {
		m.Screen = ScreenMenu
		m.FlashcardModel.LessonModel = nil
		return m, nil
	}

	return m, cmd
}

func (m *MissionTUI) startLevel(level int) {
	m.CurrentLevel = level
	m.CurrentMission = 0
	m.startCurrentMission()
}

func (m *MissionTUI) startCurrentMission() {
	missions := m.Missions[m.CurrentLevel]
	if m.CurrentMission >= len(missions) {
		// Level complete
		m.Screen = ScreenLevelSelect
		return
	}

	mission := missions[m.CurrentMission]
	m.Runner = sandbox.NewMissionRunner(mission)
	m.History = nil
	m.Input = ""
	m.ShowHint = false
	m.Screen = ScreenMission
}

func (m *MissionTUI) nextMission() {
	m.CurrentMission++
	m.startCurrentMission()
}

func (m *MissionTUI) executeCommand() {
	if m.Runner == nil {
		return
	}

	cmd := strings.TrimSpace(m.Input)
	result := m.Runner.Execute(cmd)
	m.CommandsUsed++

	entry := historyEntry{
		Command: cmd,
		Output:  result.Output,
		Error:   result.Error,
		Success: result.Success,
	}
	m.History = append(m.History, entry)
	m.Input = ""

	if result.Completed {
		m.MissionsCompleted++
		m.FlashcardModel.Progress.Practice(m.Runner.Mission.SkillID, 5) // Perfect score for completion
		m.Screen = ScreenComplete
	}
}

// View implements tea.Model.
func (m *MissionTUI) View() string {
	switch m.Screen {
	case ScreenMenu:
		return m.viewMenu()
	case ScreenLevelSelect:
		return m.viewLevelSelect()
	case ScreenMission:
		return m.viewMission()
	case ScreenComplete:
		return m.viewComplete()
	case ScreenStats:
		return m.viewStats()
	case ScreenFlashcards:
		return m.FlashcardModel.View()
	}
	return ""
}

func (m *MissionTUI) viewMenu() string {
	// For very narrow terminals, use compact layout.
	if m.Width < 60 {
		return m.viewMenuCompact()
	}

	// Hero banner with ASCII art (hide on medium terminals).
	var banner, tagline string
	if m.Width >= 80 {
		banner = HeroBanner()
		tagline = Tagline()
	} else {
		banner = TitleStyle.Render("🐢 TURTLE")
		tagline = Tagline()
	}

	// Player stats card.
	xpInLevel := m.FlashcardModel.Progress.XP % skills.XPPerLevel
	xpProgress := float64(xpInLevel) / float64(skills.XPPerLevel)
	playerCard := PlayerCard(m.FlashcardModel.Progress.Level, m.FlashcardModel.Progress.XP, m.FlashcardModel.Progress.CurrentStreak, xpProgress)

	// Build menu and columns.
	menu := m.buildMainMenu()
	columns := m.buildMenuColumns(menu)

	// Session stats if any progress.
	sessionStats := m.buildSessionStats()

	// Combine everything.
	content := lipgloss.JoinVertical(lipgloss.Center, banner, tagline, "", playerCard, "", columns)
	if sessionStats != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, content, "", sessionStats)
	}

	footer := FooterStyle.Render("  " + Arrow + "↑↓  " + Bullet + " enter  " + BulletEmpty + " q")
	return lipgloss.JoinVertical(lipgloss.Center, content, "", footer)
}

func (m *MissionTUI) viewMenuCompact() string {
	// Minimal layout for very narrow terminals.
	title := TitleStyle.Render("🐢 TURTLE")

	menu := m.buildMainMenu()
	menuBox := HighlightBoxStyle.Width(m.Width - 4).Render(menu)

	footer := MutedStyle.Render("↑↓ nav • enter • q quit")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		title,
		"",
		menuBox,
		"",
		footer,
	)
}

func (m *MissionTUI) buildMainMenu() string {
	menu := ""
	for i, item := range mainMenuItems {
		if i == m.MenuIndex {
			line := MenuItemSelectedStyle.Render(item.icon+" "+item.label) +
				"  " + MutedStyle.Render(item.desc)
			menu += line + "\n"
		} else {
			menu += MenuItemStyle.Render("  "+item.label) + "\n"
		}
	}
	return menu
}

func (m *MissionTUI) buildMenuColumns(menu string) string {
	// For narrow terminals, stack vertically instead of side-by-side.
	if m.Width < 100 {
		menuBox := HighlightBoxStyle.Width(m.Width - 4).Render(menu)
		return menuBox
	}

	// Wide terminals: two-column layout.
	menuWidth := (m.Width - 10) * 2 / 5 // 40% for menu
	infoWidth := (m.Width - 10) * 3 / 5 // 60% for info
	whatIs := WhatIsTurtleWidth(infoWidth - 6)
	menuBox := HighlightBoxStyle.Width(menuWidth).Render(menu)
	infoBox := BoxStyle.Width(infoWidth).Render(whatIs)
	return lipgloss.JoinHorizontal(lipgloss.Top, menuBox, "  ", infoBox)
}

func (m *MissionTUI) buildSessionStats() string {
	if m.MissionsCompleted == 0 && m.CommandsUsed == 0 {
		return ""
	}
	return GlowBoxStyle.Render(
		SubtitleStyle.Render("THIS SESSION") + "\n" +
			MutedStyle.Render(fmt.Sprintf("  %s Missions completed: ", MasteryFull)) +
			StatValueStyle.Render(fmt.Sprintf("%d", m.MissionsCompleted)) + "\n" +
			MutedStyle.Render(fmt.Sprintf("  %s Commands executed: ", Lightning)) +
			StatValueStyle.Render(fmt.Sprintf("%d", m.CommandsUsed)),
	)
}

func (m *MissionTUI) viewLevelSelect() string {
	// Header with small logo.
	header := lipgloss.JoinHorizontal(lipgloss.Center,
		TitleStyle.Render("🐢"),
		MutedStyle.Render(" TURTLE "),
		AccentStyle.Render(Arrow+" Missions"),
	)

	title := SubtitleStyle.Render("SELECT LEVEL")

	levelNames := LevelNames()
	levelGoals := LevelGoals()

	var levels string
	for i := 0; i < len(m.Missions); i++ {
		name := levelNames[i]
		goal := levelGoals[i]
		missionCount := len(m.Missions[i])

		if i == m.LevelIndex {
			levels += MenuItemSelectedStyle.Render(fmt.Sprintf("%s Level %d: %s", Arrow, i, name)) + "\n"
			levels += MutedStyle.Render(fmt.Sprintf("     \"%s\" (%d missions)", goal, missionCount)) + "\n"
		} else {
			levels += MenuItemStyle.Render(fmt.Sprintf("  Level %d: %s", i, name)) + "\n"
		}
	}

	content := HighlightBoxStyle.Width(50).Render(levels)
	footer := FooterStyle.Render("  " + Arrow + "↑↓ navigate  " + Bullet + " enter start  " + BulletEmpty + " esc back")

	return lipgloss.JoinVertical(lipgloss.Left, "", header, "", title, "", content, "", footer)
}

func (m *MissionTUI) viewMission() string {
	if m.Runner == nil {
		return "Loading..."
	}

	mission := m.Runner.Mission
	levelName := LevelNames()[m.CurrentLevel]

	// Header breadcrumb.
	header := lipgloss.JoinHorizontal(lipgloss.Center,
		TitleStyle.Render("🐢"),
		MutedStyle.Render(" "),
		AccentStyle.Render(levelName),
		MutedStyle.Render(fmt.Sprintf(" %s Mission %d/%d",
			Arrow, m.CurrentMission+1, len(m.Missions[m.CurrentLevel]))),
	)

	// Mission title and briefing.
	title := TitleStyle.Render(mission.Title)
	briefing := GlowBoxStyle.Width(60).Render(mission.Briefing)

	// Hint.
	var hint string
	if m.ShowHint {
		hint = AccentStyle.Render("💡 " + mission.Hint)
	} else {
		hint = MutedStyle.Render("Press ? for hint")
	}

	// Terminal output and input.
	terminalView := m.renderTerminalView()
	location := MutedStyle.Render("📍 " + m.Runner.GetCurrentLocation())
	inputLine := TerminalStyle.Render(PromptStyle.Render("$ ") + CommandStyle.Render(m.Input+"▋"))

	footer := FooterStyle.Render("  enter execute  " + Bullet + " ? hint  " + Bullet + " ctrl+r reset  " + Bullet + " esc exit")

	return lipgloss.JoinVertical(lipgloss.Left,
		header, "", title, briefing, "", hint, "", location, terminalView, inputLine, "", footer)
}

func (m *MissionTUI) renderTerminalView() string {
	if len(m.History) == 0 {
		return TerminalStyle.Render(MutedStyle.Render("Type a command and press Enter"))
	}

	var content string
	for _, entry := range m.History {
		content += PromptStyle.Render("$ ") + CommandStyle.Render(entry.Command) + "\n"
		if entry.Error != "" {
			content += DangerStyle.Render(entry.Error) + "\n"
		} else if entry.Output != "" {
			content += entry.Output + "\n"
		}
	}
	return TerminalStyle.Render(content)
}

func (m *MissionTUI) viewComplete() string {
	if m.Runner == nil {
		return ""
	}

	mission := m.Runner.Mission

	// Success header.
	header := lipgloss.JoinHorizontal(lipgloss.Center,
		TitleStyle.Render("🐢"),
		MutedStyle.Render(" TURTLE "),
	)

	title := SuccessStyle.Render(Star + " Mission Complete! " + Star)
	missionTitle := TitleStyle.Render(mission.Title)

	explanation := GlowBoxStyle.Width(60).Render(mission.Explanation)

	// Example solutions.
	var solutions string
	if len(mission.Commands) > 0 {
		solutions = MutedStyle.Render("Example solutions: ") +
			CommandStyle.Render(strings.Join(mission.Commands, ", "))
	}

	// Stats.
	stats := MutedStyle.Render(fmt.Sprintf("Commands used: %d", m.Runner.Attempts))

	footer := FooterStyle.Render("  enter next mission  " + Bullet + " esc menu")

	return lipgloss.JoinVertical(lipgloss.Left,
		"", header, "", title, missionTitle, "", explanation, "", solutions, stats, "", footer)
}

func (m *MissionTUI) viewStats() string {
	// Header.
	header := lipgloss.JoinHorizontal(lipgloss.Center,
		TitleStyle.Render("🐢"),
		MutedStyle.Render(" TURTLE "),
		AccentStyle.Render(Arrow+" Stats"),
	)

	title := SubtitleStyle.Render("YOUR STATISTICS")

	// Session stats.
	sessionStats := lipgloss.JoinVertical(lipgloss.Left,
		StatLabelStyle.Render("Missions:")+StatValueStyle.Render(fmt.Sprintf(" %d", m.MissionsCompleted)),
		StatLabelStyle.Render("Commands:")+StatValueStyle.Render(fmt.Sprintf(" %d", m.CommandsUsed)),
	)

	// Overall progress.
	progressStats := lipgloss.JoinVertical(lipgloss.Left,
		StatLabelStyle.Render("Level:")+LevelStyle.Render(fmt.Sprintf(" %d", m.FlashcardModel.Progress.Level)),
		StatLabelStyle.Render("XP:")+XPStyle.Render(fmt.Sprintf(" %d", m.FlashcardModel.Progress.XP)),
		StatLabelStyle.Render("Streak:")+StreakStyle.Render(fmt.Sprintf(" %d days", m.FlashcardModel.Progress.CurrentStreak)),
	)

	sessionBox := BoxStyle.Width(30).Render(SubtitleStyle.Render("This Session") + "\n\n" + sessionStats)
	progressBox := GlowBoxStyle.Width(30).Render(SubtitleStyle.Render("All Time") + "\n\n" + progressStats)

	content := lipgloss.JoinHorizontal(lipgloss.Top, sessionBox, "  ", progressBox)

	footer := FooterStyle.Render("  press any key to return")

	return lipgloss.JoinVertical(lipgloss.Left, "", header, "", title, "", content, "", footer)
}
