// ABOUTME: Lesson model for individual practice sessions
// ABOUTME: Handles challenge presentation, input, scoring, and feedback

package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/2389-research/turtle/internal/skills"
	"github.com/2389-research/turtle/internal/srs"
)

// ChallengeType represents different kinds of exercises
type ChallengeType int

const (
	ChallengeTypeCommand ChallengeType = iota // Type the command
	ChallengePredict                          // Predict output
	ChallengeTranslate                        // Natural language -> command
)

// Challenge represents a single practice exercise
type Challenge struct {
	Type        ChallengeType
	SkillID     string
	Prompt      string
	Expected    string   // Expected answer
	Hint        string   // Help text
	Options     []string // For multiple choice
	Explanation string   // Shown after answering
}

// LessonModel handles a practice session
type LessonModel struct {
	Progress      *skills.UserProgress
	SkillGraph    *skills.SkillGraph
	Challenges    []Challenge
	CurrentIndex  int
	Input         string
	ShowFeedback  bool
	WasCorrect    bool
	StartTime     time.Time
	XPEarned      int
	Done          bool
	Hearts        int // Lives remaining
	Combo         int // Consecutive correct
}

// NewLessonModel creates a new practice session
func NewLessonModel(progress *skills.UserProgress, graph *skills.SkillGraph) *LessonModel {
	// Generate challenges based on due skills and unlocked skills
	challenges := generateChallenges(progress, graph)

	return &LessonModel{
		Progress:     progress,
		SkillGraph:   graph,
		Challenges:   challenges,
		CurrentIndex: 0,
		Input:        "",
		ShowFeedback: false,
		StartTime:    time.Now(),
		Hearts:       3,
		Combo:        0,
	}
}

// generateChallenges creates a set of challenges for this session
func generateChallenges(progress *skills.UserProgress, graph *skills.SkillGraph) []Challenge {
	var challenges []Challenge

	// Get skills to practice (due for review + new unlocked)
	unlocked := graph.GetUnlockedSkills(progress)

	for _, skillID := range unlocked {
		// Add challenges for this skill
		skillChallenges := getChallengesForSkill(skillID)
		challenges = append(challenges, skillChallenges...)

		// Limit total challenges
		if len(challenges) >= 5 {
			break
		}
	}

	// Ensure at least some challenges
	if len(challenges) == 0 {
		challenges = getDefaultChallenges()
	}

	return challenges
}

// getChallengesForSkill returns challenges for a specific skill
func getChallengesForSkill(skillID string) []Challenge {
	switch skillID {
	case "pwd":
		return []Challenge{
			{
				Type:        ChallengeTypeCommand,
				SkillID:     "pwd",
				Prompt:      "Print the current working directory",
				Expected:    "pwd",
				Hint:        "Three letters: print working directory",
				Explanation: "pwd shows where you are in the filesystem",
			},
			{
				Type:        ChallengeTranslate,
				SkillID:     "pwd",
				Prompt:      "\"Where am I?\" translates to:",
				Expected:    "pwd",
				Hint:        "The command shows your current location",
				Explanation: "pwd = print working directory",
			},
		}
	case "ls":
		return []Challenge{
			{
				Type:        ChallengeTypeCommand,
				SkillID:     "ls",
				Prompt:      "List files in the current directory",
				Expected:    "ls",
				Hint:        "Two letters, rhymes with 'miss'",
				Explanation: "ls lists directory contents",
			},
			{
				Type:        ChallengeTypeCommand,
				SkillID:     "ls",
				Prompt:      "List ALL files (including hidden)",
				Expected:    "ls -a",
				Hint:        "Add the -a flag for 'all'",
				Explanation: "ls -a shows hidden files (starting with .)",
			},
			{
				Type:        ChallengeTypeCommand,
				SkillID:     "ls",
				Prompt:      "List files with details (long format)",
				Expected:    "ls -l",
				Hint:        "Add the -l flag for 'long'",
				Explanation: "ls -l shows permissions, size, dates",
			},
		}
	case "cd":
		return []Challenge{
			{
				Type:        ChallengeTypeCommand,
				SkillID:     "cd",
				Prompt:      "Change to your home directory",
				Expected:    "cd",
				Hint:        "Just the command, no arguments",
				Explanation: "cd alone takes you home",
			},
			{
				Type:        ChallengeTypeCommand,
				SkillID:     "cd",
				Prompt:      "Go to the /tmp directory",
				Expected:    "cd /tmp",
				Hint:        "cd followed by the path",
				Explanation: "cd <path> changes to that directory",
			},
			{
				Type:        ChallengeTypeCommand,
				SkillID:     "cd",
				Prompt:      "Go up one directory level",
				Expected:    "cd ..",
				Hint:        "Two dots represent parent directory",
				Explanation: ".. always means 'parent directory'",
			},
		}
	case "mkdir":
		return []Challenge{
			{
				Type:        ChallengeTypeCommand,
				SkillID:     "mkdir",
				Prompt:      "Create a directory called 'projects'",
				Expected:    "mkdir projects",
				Hint:        "make directory",
				Explanation: "mkdir creates a new directory",
			},
		}
	case "touch":
		return []Challenge{
			{
				Type:        ChallengeTypeCommand,
				SkillID:     "touch",
				Prompt:      "Create an empty file called 'notes.txt'",
				Expected:    "touch notes.txt",
				Hint:        "The command 'touches' files into existence",
				Explanation: "touch creates empty files",
			},
		}
	case "tmux-new":
		return []Challenge{
			{
				Type:        ChallengeTypeCommand,
				SkillID:     "tmux-new",
				Prompt:      "Start a new tmux session named 'work'",
				Expected:    "tmux new -s work",
				Hint:        "tmux new -s <name>",
				Explanation: "tmux new -s creates a named session",
			},
		}
	case "tmux-detach":
		return []Challenge{
			{
				Type:        ChallengeTypeCommand,
				SkillID:     "tmux-detach",
				Prompt:      "Detach from current tmux session (key combo)",
				Expected:    "Ctrl-b d",
				Hint:        "Prefix key + d",
				Explanation: "Ctrl-b is the prefix, then d for detach",
			},
		}
	case "tmux-attach":
		return []Challenge{
			{
				Type:        ChallengeTypeCommand,
				SkillID:     "tmux-attach",
				Prompt:      "Attach to a tmux session named 'work'",
				Expected:    "tmux attach -t work",
				Hint:        "tmux attach -t <name>",
				Explanation: "attach -t lets you reconnect to a named session",
			},
		}
	default:
		return nil
	}
}

// getDefaultChallenges returns starter challenges if none generated
func getDefaultChallenges() []Challenge {
	return getChallengesForSkill("pwd")
}

// Init implements tea.Model for LessonModel
func (m *LessonModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model for LessonModel
func (m *LessonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeypress(msg)
	}
	return m, nil
}

// handleKeypress processes input during lesson
func (m *LessonModel) handleKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.ShowFeedback {
		// Any key continues to next challenge
		switch msg.String() {
		case "enter", " ":
			m.ShowFeedback = false
			m.CurrentIndex++
			m.Input = ""

			if m.CurrentIndex >= len(m.Challenges) {
				m.Done = true
			}
		case "q", "esc":
			m.Done = true
		}
		return m, nil
	}

	// Regular input handling
	switch msg.String() {
	case "enter":
		m.checkAnswer()
	case "backspace":
		if len(m.Input) > 0 {
			m.Input = m.Input[:len(m.Input)-1]
		}
	case "ctrl+c":
		m.Done = true
	case "?":
		// Show hint (could implement hint reveal here)
	default:
		// Regular character input
		if len(msg.String()) == 1 {
			m.Input += msg.String()
		} else if msg.String() == "space" {
			m.Input += " "
		}
	}

	return m, nil
}

// checkAnswer evaluates the user's response
func (m *LessonModel) checkAnswer() {
	if m.CurrentIndex >= len(m.Challenges) {
		return
	}

	challenge := m.Challenges[m.CurrentIndex]
	elapsed := time.Since(m.StartTime).Milliseconds()

	// Normalize both for comparison
	userAnswer := strings.TrimSpace(strings.ToLower(m.Input))
	expected := strings.TrimSpace(strings.ToLower(challenge.Expected))

	m.WasCorrect = userAnswer == expected

	// Calculate grade for SRS
	grade := srs.CalculateGrade(m.WasCorrect, elapsed)

	// Update skill progress
	m.Progress.Practice(challenge.SkillID, grade)

	// Update game state
	if m.WasCorrect {
		m.Combo++
		baseXP := 10
		comboBonus := m.Combo * 2
		m.XPEarned += baseXP + comboBonus
		m.Progress.AddXP(baseXP + comboBonus)
	} else {
		m.Combo = 0
		m.Hearts--
		if m.Hearts <= 0 {
			m.Done = true
		}
	}

	// Record activity for streaks
	m.Progress.RecordActivity()

	m.ShowFeedback = true
	m.StartTime = time.Now() // Reset for next challenge
}

// View renders the lesson screen
func (m *LessonModel) View() string {
	if m.Done {
		return m.renderComplete()
	}

	if m.CurrentIndex >= len(m.Challenges) {
		return m.renderComplete()
	}

	challenge := m.Challenges[m.CurrentIndex]

	// Header with hearts and combo
	hearts := strings.Repeat("❤️ ", m.Hearts) + strings.Repeat("🖤 ", 3-m.Hearts)
	combo := ""
	if m.Combo > 1 {
		combo = AccentStyle.Render(fmt.Sprintf("🔥 %dx combo!", m.Combo))
	}

	header := lipgloss.JoinHorizontal(
		lipgloss.Center,
		hearts,
		"  ",
		combo,
	)

	// Progress indicator
	progress := fmt.Sprintf("Question %d/%d", m.CurrentIndex+1, len(m.Challenges))
	progressBar := ProgressBar(20, float64(m.CurrentIndex)/float64(len(m.Challenges)))

	// Challenge prompt
	prompt := BoxStyle.Render(
		SubtitleStyle.Render(challenge.Prompt),
	)

	// Input area
	cursor := "▋"
	if m.ShowFeedback {
		cursor = ""
	}
	inputDisplay := TerminalStyle.Render(
		PromptStyle.Render("$ ") + CommandStyle.Render(m.Input+cursor),
	)

	// Feedback (if showing)
	var feedback string
	if m.ShowFeedback {
		if m.WasCorrect {
			feedback = SuccessStyle.Render("✓ Correct!") + "\n"
			if m.Combo > 1 {
				feedback += AccentStyle.Render(fmt.Sprintf("+%d XP (combo bonus!)", 10+m.Combo*2)) + "\n"
			} else {
				feedback += MutedStyle.Render("+10 XP") + "\n"
			}
		} else {
			feedback = DangerStyle.Render("✗ Not quite") + "\n"
			feedback += MutedStyle.Render("Expected: ") + CommandStyle.Render(challenge.Expected) + "\n"
		}
		feedback += "\n" + MutedStyle.Render(challenge.Explanation)
		feedback += "\n\n" + MutedStyle.Render("Press enter to continue...")
	}

	// Hint
	hint := MutedStyle.Render("💡 " + challenge.Hint)

	// Combine everything
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		MutedStyle.Render(progress),
		progressBar,
		"",
		prompt,
		"",
		inputDisplay,
		"",
		hint,
		"",
		feedback,
	)
}

// renderComplete shows the end-of-lesson summary
func (m *LessonModel) renderComplete() string {
	var title string
	if m.Hearts <= 0 {
		title = DangerStyle.Render("💔 Out of hearts!")
	} else {
		title = SuccessStyle.Render("🎉 Lesson Complete!")
	}

	stats := fmt.Sprintf(`
  Questions: %d
  XP Earned: %s
  Best Combo: %dx
`,
		len(m.Challenges),
		XPStyle.Render(fmt.Sprintf("+%d", m.XPEarned)),
		m.Combo,
	)

	motivation := ""
	if m.XPEarned > 50 {
		motivation = AccentStyle.Render("🐢 Turtle-ific! You're on fire!")
	} else if m.XPEarned > 20 {
		motivation = AccentStyle.Render("🐢 Slow and steady wins the race!")
	} else {
		motivation = MutedStyle.Render("🐢 Every step counts. Keep going!")
	}

	footer := MutedStyle.Render("\nPress any key to continue...")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		BoxStyle.Render(stats),
		"",
		motivation,
		footer,
	)
}
