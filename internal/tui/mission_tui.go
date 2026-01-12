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
	Screen     MissionScreen
	Progress   *skills.UserProgress
	SkillGraph *skills.SkillGraph
	Missions   map[int][]*sandbox.Mission

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

	// Flashcard mode
	FlashcardModel Model
}

type historyEntry struct {
	Command string
	Output  string
	Error   string
	Success bool
}

// NewMissionTUI creates a new mission-based learning interface.
func NewMissionTUI() *MissionTUI {
	return &MissionTUI{
		Screen:         ScreenMenu,
		Progress:       skills.NewUserProgress(),
		SkillGraph:     buildPedagogicalSkillGraph(),
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
	switch msg.String() {
	case "up", "k":
		if m.MenuIndex > 0 {
			m.MenuIndex--
		}
	case "down", "j":
		if m.MenuIndex < 3 {
			m.MenuIndex++
		}
	case "enter", " ":
		switch m.MenuIndex {
		case 0: // Missions
			m.Screen = ScreenLevelSelect
			m.LevelIndex = 0
		case 1: // Flashcards
			m.Screen = ScreenFlashcards
		case 2: // Stats
			m.Screen = ScreenStats
		case 3: // Quit
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

func (m *MissionTUI) updateStats(_ tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Any key returns to menu
	m.Screen = ScreenMenu
	return m, nil
}

func (m *MissionTUI) updateFlashcards(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Check for escape to return to main menu
	if msg.String() == "esc" && m.FlashcardModel.CurrentView == ViewHome {
		m.Screen = ScreenMenu
		return m, nil
	}

	// Delegate to flashcard model
	updatedModel, cmd := m.FlashcardModel.Update(msg)
	m.FlashcardModel = updatedModel.(Model)
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
		m.Progress.Practice(m.Runner.Mission.SkillID, 5) // Perfect score for completion
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
	title := TitleStyle.Render("🐢 TURTLE")
	subtitle := MutedStyle.Render("Learn the terminal, one mission at a time")

	menuItems := []string{"Missions", "Flashcards", "Stats", "Quit"}
	var menu string
	for i, item := range menuItems {
		cursor := "  "
		style := MutedStyle
		if i == m.MenuIndex {
			cursor = "> "
			style = AccentStyle
		}
		menu += style.Render(cursor+item) + "\n"
	}

	footer := MutedStyle.Render("↑/↓ navigate • enter select • q quit")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		title,
		subtitle,
		"",
		menu,
		"",
		footer,
	)
}

func (m *MissionTUI) viewLevelSelect() string {
	title := TitleStyle.Render("Select Level")

	levelNames := LevelNames()
	levelGoals := LevelGoals()

	var levels string
	for i := 0; i < len(m.Missions); i++ {
		cursor := "  "
		style := MutedStyle
		if i == m.LevelIndex {
			cursor = "> "
			style = AccentStyle
		}

		name := levelNames[i]
		goal := levelGoals[i]
		missionCount := len(m.Missions[i])

		levelLine := fmt.Sprintf("%sLevel %d: %s", cursor, i, name)
		levels += style.Render(levelLine) + "\n"
		if i == m.LevelIndex {
			levels += MutedStyle.Render(fmt.Sprintf("    \"%s\" (%d missions)", goal, missionCount)) + "\n"
		}
	}

	footer := MutedStyle.Render("↑/↓ navigate • enter start • esc back")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		"",
		levels,
		"",
		footer,
	)
}

func (m *MissionTUI) viewMission() string {
	if m.Runner == nil {
		return "Loading..."
	}

	mission := m.Runner.Mission

	// Header
	levelName := LevelNames()[m.CurrentLevel]
	header := MutedStyle.Render(fmt.Sprintf("Level %d: %s • Mission %d/%d",
		m.CurrentLevel, levelName, m.CurrentMission+1, len(m.Missions[m.CurrentLevel])))

	// Mission title and briefing
	title := TitleStyle.Render(mission.Title)
	briefing := BoxStyle.Render(mission.Briefing)

	// Hint (if shown)
	var hint string
	if m.ShowHint {
		hint = AccentStyle.Render("💡 " + mission.Hint)
	} else {
		hint = MutedStyle.Render("Press ? for hint")
	}

	// Terminal output (history)
	var terminalContent string
	for _, entry := range m.History {
		terminalContent += PromptStyle.Render("$ ") + CommandStyle.Render(entry.Command) + "\n"
		if entry.Error != "" {
			terminalContent += DangerStyle.Render(entry.Error) + "\n"
		} else if entry.Output != "" {
			terminalContent += entry.Output + "\n"
		}
	}

	// Current location indicator
	location := MutedStyle.Render("📍 " + m.Runner.GetCurrentLocation())

	// Input line
	cursor := "▋"
	inputLine := TerminalStyle.Render(
		PromptStyle.Render("$ ") + CommandStyle.Render(m.Input+cursor),
	)

	// Build terminal view
	terminalView := TerminalStyle.Render(terminalContent)
	if terminalContent == "" {
		terminalView = TerminalStyle.Render(MutedStyle.Render("Type a command and press Enter"))
	}

	// Footer
	footer := MutedStyle.Render("enter execute • ctrl+r reset • esc exit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		title,
		briefing,
		"",
		hint,
		"",
		location,
		terminalView,
		inputLine,
		"",
		footer,
	)
}

func (m *MissionTUI) viewComplete() string {
	if m.Runner == nil {
		return ""
	}

	mission := m.Runner.Mission

	title := SuccessStyle.Render("✓ Mission Complete!")
	missionTitle := TitleStyle.Render(mission.Title)

	explanation := BoxStyle.Render(mission.Explanation)

	// Commands that could solve it
	var commandsUsed string
	if len(mission.Commands) > 0 {
		commandsUsed = MutedStyle.Render("Example solutions: ") +
			CommandStyle.Render(strings.Join(mission.Commands, ", "))
	}

	// Stats for this mission
	stats := fmt.Sprintf("Commands used: %d", m.Runner.Attempts)

	footer := MutedStyle.Render("Press enter for next mission • esc to menu")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		missionTitle,
		"",
		explanation,
		"",
		commandsUsed,
		MutedStyle.Render(stats),
		"",
		footer,
	)
}

func (m *MissionTUI) viewStats() string {
	title := TitleStyle.Render("Stats")

	stats := fmt.Sprintf(`
  Missions Completed: %d
  Commands Used: %d
  Level: %d
  XP: %d
  Streak: %d days
`,
		m.MissionsCompleted,
		m.CommandsUsed,
		m.Progress.Level,
		m.Progress.XP,
		m.Progress.CurrentStreak,
	)

	footer := MutedStyle.Render("Press any key to return")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		stats,
		"",
		footer,
	)
}
