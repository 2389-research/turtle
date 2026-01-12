// ABOUTME: Skill selector screen for choosing which skills to practice
// ABOUTME: Allows users to select specific skills before starting flashcard practice

package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/2389-research/turtle/internal/skills"
)

// SkillSelectorModel manages skill selection state.
type SkillSelectorModel struct {
	graph            *skills.SkillGraph
	progress         *skills.UserProgress
	categories       []skills.Category
	skillsByCategory map[skills.Category][]string
	unlockedSkills   map[string]bool
	selected         map[string]bool
	flatSkills       []flatSkillItem // Flat list for navigation
	cursorIndex      int
	width            int
	height           int
	Done             bool
	Cancelled        bool
	errorMsg         string
}

// flatSkillItem represents a navigable item (either category header or skill).
type flatSkillItem struct {
	isCategory bool
	category   skills.Category
	skillID    string
}

// categoryDisplayNames maps category IDs to display names.
var categoryDisplayNames = map[skills.Category]string{
	skills.CategoryNavigation:  "Navigation",
	skills.CategoryFileOps:     "File Operations",
	skills.CategoryTmuxBasics:  "Tmux Basics",
	skills.CategoryTmuxPanes:   "Tmux Panes",
	skills.CategoryTmuxWindows: "Tmux Windows",
	skills.CategoryAdvanced:    "Advanced",
}

// NewSkillSelectorModel creates a new skill selector.
func NewSkillSelectorModel(progress *skills.UserProgress, graph *skills.SkillGraph) *SkillSelectorModel {
	m := &SkillSelectorModel{
		graph:            graph,
		progress:         progress,
		selected:         make(map[string]bool),
		skillsByCategory: make(map[skills.Category][]string),
		unlockedSkills:   make(map[string]bool),
		width:            80,
		height:           24,
	}

	// Define category order
	m.categories = []skills.Category{
		skills.CategoryNavigation,
		skills.CategoryFileOps,
		skills.CategoryTmuxBasics,
		skills.CategoryTmuxPanes,
		skills.CategoryTmuxWindows,
		skills.CategoryAdvanced,
	}

	// Get unlocked skills
	unlockedList := graph.GetUnlockedSkills(progress)
	for _, skillID := range unlockedList {
		m.unlockedSkills[skillID] = true
	}

	// Organize skills by category
	for _, cat := range m.categories {
		skillIDs := graph.GetSkillsByCategory(cat)
		if len(skillIDs) > 0 {
			m.skillsByCategory[cat] = skillIDs
		}
	}

	// Build flat navigation list
	m.buildFlatList()

	// Pre-select all unlocked skills
	for skillID := range m.unlockedSkills {
		m.selected[skillID] = true
	}

	// Position cursor on first skill (skip first category header)
	m.cursorIndex = 0
	m.moveToNextSkill()

	return m
}

// buildFlatList creates a flat list of navigable items.
func (m *SkillSelectorModel) buildFlatList() {
	m.flatSkills = nil

	for _, cat := range m.categories {
		skillIDs, ok := m.skillsByCategory[cat]
		if !ok || len(skillIDs) == 0 {
			continue
		}

		// Add category header
		m.flatSkills = append(m.flatSkills, flatSkillItem{
			isCategory: true,
			category:   cat,
		})

		// Add skills in this category
		for _, skillID := range skillIDs {
			m.flatSkills = append(m.flatSkills, flatSkillItem{
				isCategory: false,
				skillID:    skillID,
				category:   cat,
			})
		}
	}
}

// moveToNextSkill moves cursor to next skill (skipping category headers).
func (m *SkillSelectorModel) moveToNextSkill() {
	for m.cursorIndex < len(m.flatSkills) && m.flatSkills[m.cursorIndex].isCategory {
		m.cursorIndex++
	}
}

// Update handles input.
func (m *SkillSelectorModel) Update(msg tea.Msg) (*SkillSelectorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.errorMsg = ""

		switch msg.String() {
		case "up", "k":
			m.cursorIndex--
			if m.cursorIndex < 0 {
				m.cursorIndex = len(m.flatSkills) - 1
			}
			// Skip category headers when navigating up
			for m.cursorIndex > 0 && m.flatSkills[m.cursorIndex].isCategory {
				m.cursorIndex--
			}

		case "down", "j":
			m.cursorIndex++
			if m.cursorIndex >= len(m.flatSkills) {
				m.cursorIndex = 0
				m.moveToNextSkill()
			} else {
				// Skip category headers when navigating down
				for m.cursorIndex < len(m.flatSkills) && m.flatSkills[m.cursorIndex].isCategory {
					m.cursorIndex++
				}
				if m.cursorIndex >= len(m.flatSkills) {
					m.cursorIndex = 0
					m.moveToNextSkill()
				}
			}

		case " ": // Space to toggle
			m.toggleCurrentSkill()

		case "a": // Select all unlocked
			for skillID := range m.unlockedSkills {
				m.selected[skillID] = true
			}

		case "n": // Select none
			m.selected = make(map[string]bool)

		case "enter", "s": // Start practice
			selectedCount := len(m.GetSelectedSkills())
			if selectedCount == 0 {
				m.errorMsg = "Please select at least one skill"
			} else {
				m.Done = true
			}

		case "esc", "backspace":
			m.Cancelled = true
			m.Done = true
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// toggleCurrentSkill toggles selection of the skill under cursor.
func (m *SkillSelectorModel) toggleCurrentSkill() {
	if m.cursorIndex < 0 || m.cursorIndex >= len(m.flatSkills) {
		return
	}

	item := m.flatSkills[m.cursorIndex]
	if item.isCategory {
		return
	}

	// Only allow toggling unlocked skills
	if !m.unlockedSkills[item.skillID] {
		m.errorMsg = "This skill is locked"
		return
	}

	if m.selected[item.skillID] {
		delete(m.selected, item.skillID)
	} else {
		m.selected[item.skillID] = true
	}
}

// GetSelectedSkills returns the list of selected skill IDs.
func (m *SkillSelectorModel) GetSelectedSkills() []string {
	result := make([]string, 0, len(m.selected))
	for skillID := range m.selected {
		result = append(result, skillID)
	}
	return result
}

// View renders the skill selector.
func (m *SkillSelectorModel) View() string {
	var b strings.Builder

	// Title
	title := SubtitleStyle.Render("SELECT SKILLS TO PRACTICE")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Skill list
	for i, item := range m.flatSkills {
		if item.isCategory {
			// Category header
			name := categoryDisplayNames[item.category]
			b.WriteString("\n")
			b.WriteString(AccentStyle.Render(name))
			b.WriteString("\n")
		} else {
			// Skill item
			isSelected := m.selected[item.skillID]
			isUnlocked := m.unlockedSkills[item.skillID]
			isCursor := i == m.cursorIndex

			line := m.renderSkillLine(item.skillID, isSelected, isUnlocked, isCursor)
			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	// Selection count
	selectedCount := len(m.GetSelectedSkills())
	b.WriteString("\n")
	countStyle := MutedStyle
	if selectedCount == 0 {
		countStyle = DangerStyle
	}
	b.WriteString(countStyle.Render(fmt.Sprintf("%d skill(s) selected", selectedCount)))
	b.WriteString("\n")

	// Error message
	if m.errorMsg != "" {
		b.WriteString(DangerStyle.Render(m.errorMsg))
		b.WriteString("\n")
	}

	// Help
	help := MutedStyle.Render("[Space] Toggle  [A] All  [N] None  [Enter] Start  [Esc] Cancel")
	b.WriteString("\n")
	b.WriteString(help)

	return BoxStyle.Render(b.String())
}

// renderSkillLine renders a single skill line.
func (m *SkillSelectorModel) renderSkillLine(skillID string, isSelected, isUnlocked, isCursor bool) string {
	// Cursor indicator
	cursor := "  "
	if isCursor {
		cursor = "> "
	}

	// Skill status indicator (mastery level)
	var status string
	if !isUnlocked {
		status = "○" // Locked
	} else {
		strength := m.progress.GetStrength(skillID)
		switch {
		case strength >= 0.8:
			status = "●" // Mastered
		case strength >= 0.5:
			status = "◑" // Practicing
		default:
			status = "◐" // Learning
		}
	}

	// Checkbox
	checkbox := "[ ]"
	if isSelected {
		checkbox = "[✓]"
	}

	// Get skill name (fallback to ID)
	name := skillID

	// Style based on state
	var lineStyle lipgloss.Style
	switch {
	case !isUnlocked:
		lineStyle = MutedStyle
		checkbox = "[-]" // Show as unavailable
	case isCursor:
		lineStyle = MenuItemSelectedStyle
	default:
		lineStyle = MenuItemStyle
	}

	return lineStyle.Render(fmt.Sprintf("%s%s %s %s", cursor, status, checkbox, name))
}
