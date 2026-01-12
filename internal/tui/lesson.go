// ABOUTME: Lesson model for individual practice sessions
// ABOUTME: Handles challenge presentation, input, scoring, and feedback

package tui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/2389-research/turtle/internal/skills"
	"github.com/2389-research/turtle/internal/srs"
)

// SpeedRoundDuration is the time limit for speed rounds.
const SpeedRoundDuration = 30 // seconds

// tickMsg is sent every second during speed rounds.
type tickMsg time.Time

// tickCmd returns a command that sends a tick every second.
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// ChallengeType represents different kinds of exercises.
type ChallengeType int

const (
	ChallengeTypeCommand    ChallengeType = iota // Type the command
	ChallengeMultipleChoice                      // Pick from options A/B/C/D
	ChallengeFixError                            // Fix a broken command
	ChallengePredictOutput                       // What does this output?
	ChallengeTranslate                           // Natural language -> command
	ChallengeSpeedRound                          // Rapid fire timed questions
)

// Challenge represents a single practice exercise.
type Challenge struct {
	Type          ChallengeType
	SkillID       string
	Prompt        string
	Expected      string   // Expected answer (or correct option index for MC)
	Hint          string   // Help text
	Options       []string // For multiple choice
	Explanation   string   // Shown after answering
	BrokenCommand string   // For fix-the-error: the broken command to fix
	CommandOutput string   // For predict-output: what the command actually outputs
}

// LessonModel handles a practice session.
type LessonModel struct {
	Progress       *skills.UserProgress
	SkillGraph     *skills.SkillGraph
	Challenges     []Challenge
	CurrentIndex   int
	Input          string
	ShowFeedback   bool
	WasCorrect     bool
	StartTime      time.Time
	XPEarned       int
	Done           bool
	Hearts         int // Lives remaining
	Combo          int // Consecutive correct
	SelectedOption int // For multiple choice (0-3)
	IsSpeedRound   bool
	SpeedTimeLeft  int // Seconds remaining in speed round
}

// NewLessonModel creates a new practice session.
// If selectedSkills is nil or empty, all unlocked skills are used.
func NewLessonModel(progress *skills.UserProgress, graph *skills.SkillGraph, selectedSkills []string) *LessonModel {
	// Generate challenges based on selected skills (or all unlocked if none selected)
	challenges := generateChallenges(progress, graph, selectedSkills)

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

// NewSpeedRoundModel creates a timed speed round session.
// If selectedSkills is nil or empty, all unlocked skills are used.
func NewSpeedRoundModel(progress *skills.UserProgress, graph *skills.SkillGraph, selectedSkills []string) *LessonModel {
	// Generate many quick challenges for speed round
	challenges := generateSpeedChallenges(progress, graph, selectedSkills)

	return &LessonModel{
		Progress:      progress,
		SkillGraph:    graph,
		Challenges:    challenges,
		CurrentIndex:  0,
		Input:         "",
		ShowFeedback:  false,
		StartTime:     time.Now(),
		Hearts:        1, // One mistake ends speed round
		Combo:         0,
		IsSpeedRound:  true,
		SpeedTimeLeft: SpeedRoundDuration,
	}
}

// generateSpeedChallenges creates rapid-fire challenges for speed rounds.
// If selectedSkills is nil or empty, all unlocked skills are used.
func generateSpeedChallenges(progress *skills.UserProgress, graph *skills.SkillGraph, selectedSkills []string) []Challenge {
	// Determine which skills to use
	var skillsToUse []string
	if len(selectedSkills) > 0 {
		skillsToUse = selectedSkills
	} else {
		skillsToUse = graph.GetUnlockedSkills(progress)
	}

	// Collect all quick challenges from selected skills
	var allChallenges []Challenge
	for _, skillID := range skillsToUse {
		skillChallenges := getChallengesForSkill(skillID)
		// Only include type-command and multiple choice (quick to answer)
		for _, c := range skillChallenges {
			if c.Type == ChallengeTypeCommand || c.Type == ChallengeMultipleChoice {
				allChallenges = append(allChallenges, c)
			}
		}
	}

	// Shuffle and take up to 20 challenges
	rand.Shuffle(len(allChallenges), func(i, j int) {
		allChallenges[i], allChallenges[j] = allChallenges[j], allChallenges[i]
	})

	if len(allChallenges) > 20 {
		allChallenges = allChallenges[:20]
	}

	// If not enough, repeat some
	for len(allChallenges) < 10 {
		allChallenges = append(allChallenges, getChallengesForSkill("pwd")...)
	}

	return allChallenges
}

// generateChallenges creates a set of challenges for this session.
// If selectedSkills is nil or empty, all unlocked skills are used.
func generateChallenges(progress *skills.UserProgress, graph *skills.SkillGraph, selectedSkills []string) []Challenge {
	var challenges []Challenge

	// Determine which skills to use
	var skillsToUse []string
	if len(selectedSkills) > 0 {
		skillsToUse = selectedSkills
	} else {
		skillsToUse = graph.GetUnlockedSkills(progress)
	}

	for _, skillID := range skillsToUse {
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

// getChallengesForSkill returns challenges for a specific skill from the challenge database.
func getChallengesForSkill(skillID string) []Challenge {
	allChallenges := getAllChallenges()
	if challenges, ok := allChallenges[skillID]; ok {
		return challenges
	}
	return nil
}

// getDefaultChallenges returns starter challenges if none generated.
func getDefaultChallenges() []Challenge {
	return getChallengesForSkill("pwd")
}

// Init implements tea.Model for LessonModel.
func (m *LessonModel) Init() tea.Cmd {
	if m.IsSpeedRound {
		return tickCmd()
	}
	return nil
}

// Update implements tea.Model for LessonModel.
func (m *LessonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.IsSpeedRound && !m.Done && !m.ShowFeedback {
			m.SpeedTimeLeft--
			if m.SpeedTimeLeft <= 0 {
				m.Done = true
				return m, nil
			}
			return m, tickCmd()
		}
		return m, nil
	case tea.KeyMsg:
		newModel, cmd := m.handleKeypress(msg)
		// Continue ticking during speed round
		if m.IsSpeedRound && !m.Done {
			return newModel, tea.Batch(cmd, tickCmd())
		}
		return newModel, cmd
	}
	return m, nil
}

// handleKeypress processes input during lesson.
func (m *LessonModel) handleKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.ShowFeedback {
		// Any key continues to next challenge
		switch msg.String() {
		case "enter", " ":
			m.ShowFeedback = false
			m.CurrentIndex++
			m.Input = ""
			m.SelectedOption = 0

			if m.CurrentIndex >= len(m.Challenges) {
				m.Done = true
			} else {
				m.StartTime = time.Now() // Reset timer for next question
			}
		case "q", "esc":
			m.Done = true
		}
		return m, nil
	}

	// Get current challenge type
	var challengeType ChallengeType
	if m.CurrentIndex < len(m.Challenges) {
		challengeType = m.Challenges[m.CurrentIndex].Type
	}

	// Multiple choice handling
	if challengeType == ChallengeMultipleChoice || challengeType == ChallengePredictOutput {
		return m.handleMultipleChoiceInput(msg)
	}

	// Regular input handling (type command, fix error, translate)
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

// handleMultipleChoiceInput processes input for multiple choice questions.
func (m *LessonModel) handleMultipleChoiceInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	challenge := m.Challenges[m.CurrentIndex]
	numOptions := len(challenge.Options)

	switch msg.String() {
	case "1", "a", "A":
		if numOptions >= 1 {
			m.SelectedOption = 0
			m.checkAnswer()
		}
	case "2", "b", "B":
		if numOptions >= 2 {
			m.SelectedOption = 1
			m.checkAnswer()
		}
	case "3", "c", "C":
		if numOptions >= 3 {
			m.SelectedOption = 2
			m.checkAnswer()
		}
	case "4", "d", "D":
		if numOptions >= 4 {
			m.SelectedOption = 3
			m.checkAnswer()
		}
	case "up", "k":
		m.SelectedOption--
		if m.SelectedOption < 0 {
			m.SelectedOption = numOptions - 1
		}
	case "down", "j":
		m.SelectedOption++
		if m.SelectedOption >= numOptions {
			m.SelectedOption = 0
		}
	case "enter", " ":
		m.checkAnswer()
	case "ctrl+c", "q":
		m.Done = true
	}

	return m, nil
}

// checkAnswer evaluates the user's response.
func (m *LessonModel) checkAnswer() {
	if m.CurrentIndex >= len(m.Challenges) {
		return
	}

	challenge := m.Challenges[m.CurrentIndex]
	elapsed := time.Since(m.StartTime).Milliseconds()

	// Check answer based on challenge type
	switch challenge.Type {
	case ChallengeMultipleChoice, ChallengePredictOutput:
		// For multiple choice, Expected contains the correct option index as string ("0", "1", etc.)
		expectedIdx := 0
		_, _ = fmt.Sscanf(challenge.Expected, "%d", &expectedIdx)
		m.WasCorrect = m.SelectedOption == expectedIdx

	case ChallengeFixError:
		// For fix-the-error, user types the corrected command
		userAnswer := strings.TrimSpace(strings.ToLower(m.Input))
		expected := strings.TrimSpace(strings.ToLower(challenge.Expected))
		m.WasCorrect = userAnswer == expected

	default:
		// Type command, translate - direct text comparison
		userAnswer := strings.TrimSpace(strings.ToLower(m.Input))
		expected := strings.TrimSpace(strings.ToLower(challenge.Expected))
		m.WasCorrect = userAnswer == expected
	}

	// Calculate grade for SRS
	grade := srs.CalculateGrade(m.WasCorrect, elapsed)

	// Update skill progress
	m.Progress.Practice(challenge.SkillID, grade)

	// Calculate XP with bonuses
	baseXP := 10
	if m.IsSpeedRound {
		baseXP = 5 // Less XP per question but more questions
	}

	// Update game state
	if m.WasCorrect {
		m.Combo++
		comboBonus := m.Combo * 2
		if m.Combo > 5 {
			comboBonus = 10 // Cap combo bonus
		}
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
}

// View renders the lesson screen.
//
//nolint:funlen // View function has many UI components
func (m *LessonModel) View() string {
	if m.Done {
		return m.renderComplete()
	}

	if m.CurrentIndex >= len(m.Challenges) {
		return m.renderComplete()
	}

	challenge := m.Challenges[m.CurrentIndex]

	// Header with hearts/timer and combo
	var header string
	if m.IsSpeedRound {
		// Speed round shows timer instead of hearts
		timerStyle := AccentStyle
		if m.SpeedTimeLeft <= 10 {
			timerStyle = DangerStyle
		}
		timer := timerStyle.Render(fmt.Sprintf("%s %ds", Lightning, m.SpeedTimeLeft))
		combo := ""
		if m.Combo > 1 {
			combo = ComboStyle.Render(fmt.Sprintf("%s %dx", Fire, m.Combo))
		}
		header = lipgloss.JoinHorizontal(
			lipgloss.Center,
			TitleStyle.Render(Lightning+" SPEED ROUND"),
			"  "+Bullet+"  ",
			timer,
			"  ",
			combo,
		)
	} else {
		hearts := Hearts(m.Hearts, 3)
		combo := ""
		if m.Combo > 1 {
			combo = ComboStyle.Render(fmt.Sprintf("%s %dx combo!", Fire, m.Combo))
		}
		header = lipgloss.JoinHorizontal(
			lipgloss.Center,
			hearts,
			"  ",
			combo,
		)
	}

	// Progress indicator
	progress := fmt.Sprintf("Question %d/%d", m.CurrentIndex+1, len(m.Challenges))
	progressBar := ProgressBar(20, float64(m.CurrentIndex)/float64(len(m.Challenges)))

	// Render based on challenge type
	var inputArea string
	var prompt string

	switch challenge.Type {
	case ChallengeMultipleChoice, ChallengePredictOutput:
		prompt, inputArea = m.renderMultipleChoice(challenge)
	case ChallengeFixError:
		prompt, inputArea = m.renderFixError(challenge)
	default:
		prompt, inputArea = m.renderTypeCommand(challenge)
	}

	// Feedback (if showing)
	feedback := m.renderFeedback(challenge)

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
		inputArea,
		"",
		hint,
		"",
		feedback,
	)
}

// renderTypeCommand renders a type-the-command challenge.
func (m *LessonModel) renderTypeCommand(challenge Challenge) (string, string) {
	prompt := BoxStyle.Render(
		SubtitleStyle.Render(challenge.Prompt),
	)

	cursor := "▋"
	if m.ShowFeedback {
		cursor = ""
	}
	inputArea := TerminalStyle.Render(
		PromptStyle.Render("$ ") + CommandStyle.Render(m.Input+cursor),
	)

	return prompt, inputArea
}

// renderMultipleChoice renders a multiple choice challenge.
func (m *LessonModel) renderMultipleChoice(challenge Challenge) (string, string) {
	var promptText string
	if challenge.Type == ChallengePredictOutput {
		promptText = challenge.Prompt + "\n\n" + TerminalStyle.Render(
			PromptStyle.Render("$ ")+CommandStyle.Render(challenge.BrokenCommand),
		)
	} else {
		promptText = challenge.Prompt
	}
	prompt := BoxStyle.Render(SubtitleStyle.Render(promptText))

	// Render options
	optionLabels := []string{"A", "B", "C", "D"}
	var options string
	for i, opt := range challenge.Options {
		prefix := "  "
		style := MutedStyle
		if i == m.SelectedOption && !m.ShowFeedback {
			prefix = "> "
			style = AccentStyle
		}
		if m.ShowFeedback {
			expectedIdx := 0
			_, _ = fmt.Sscanf(challenge.Expected, "%d", &expectedIdx)
			switch i {
			case expectedIdx:
				style = SuccessStyle
				prefix = "✓ "
			case m.SelectedOption:
				style = DangerStyle
				prefix = "✗ "
			}
		}
		options += style.Render(fmt.Sprintf("%s%s) %s", prefix, optionLabels[i], opt)) + "\n"
	}

	inputArea := BoxStyle.Render(options)
	return prompt, inputArea
}

// renderFixError renders a fix-the-error challenge.
func (m *LessonModel) renderFixError(challenge Challenge) (string, string) {
	promptText := challenge.Prompt + "\n\n" + DangerStyle.Render("Broken: ") +
		TerminalStyle.Render(CommandStyle.Render(challenge.BrokenCommand))
	prompt := BoxStyle.Render(SubtitleStyle.Render(promptText))

	cursor := "▋"
	if m.ShowFeedback {
		cursor = ""
	}
	inputArea := TerminalStyle.Render(
		PromptStyle.Render("$ ") + CommandStyle.Render(m.Input+cursor),
	)

	return prompt, inputArea
}

// renderFeedback renders the feedback after answering.
//
//nolint:nestif // XP calculation has conditional logic
func (m *LessonModel) renderFeedback(challenge Challenge) string {
	if !m.ShowFeedback {
		return ""
	}

	var feedback string
	if m.WasCorrect {
		feedback = SuccessStyle.Render("✓ Correct!") + "\n"
		xpGained := 10
		if m.IsSpeedRound {
			xpGained = 5
		}
		comboBonus := m.Combo * 2
		if m.Combo > 5 {
			comboBonus = 10
		}
		if m.Combo > 1 {
			feedback += AccentStyle.Render(fmt.Sprintf("+%d XP (combo bonus!)", xpGained+comboBonus)) + "\n"
		} else {
			feedback += MutedStyle.Render(fmt.Sprintf("+%d XP", xpGained)) + "\n"
		}
	} else {
		feedback = DangerStyle.Render("✗ Not quite") + "\n"

		// Show expected answer based on type
		switch challenge.Type {
		case ChallengeMultipleChoice, ChallengePredictOutput:
			expectedIdx := 0
			_, _ = fmt.Sscanf(challenge.Expected, "%d", &expectedIdx)
			if expectedIdx < len(challenge.Options) {
				feedback += MutedStyle.Render("Answer: ") + CommandStyle.Render(challenge.Options[expectedIdx]) + "\n"
			}
		default:
			feedback += MutedStyle.Render("Expected: ") + CommandStyle.Render(challenge.Expected) + "\n"
		}
	}
	feedback += "\n" + MutedStyle.Render(challenge.Explanation)
	feedback += "\n\n" + MutedStyle.Render("Press enter to continue...")

	return feedback
}

// renderComplete shows the end-of-lesson summary.
//
//nolint:funlen,nestif // Complex summary with many stats
func (m *LessonModel) renderComplete() string {
	var title string
	var motivation string

	if m.IsSpeedRound {
		// Speed round results
		switch {
		case m.SpeedTimeLeft <= 0:
			title = DangerStyle.Render("⏱️  Time's up!")
		case m.Hearts <= 0:
			title = DangerStyle.Render("💔 Wrong answer!")
		default:
			title = SuccessStyle.Render("⚡ Speed Round Complete!")
		}

		// Calculate questions per minute
		questionsAnswered := m.CurrentIndex
		timeUsed := SpeedRoundDuration - m.SpeedTimeLeft
		var qpm float64
		if timeUsed > 0 {
			qpm = float64(questionsAnswered) / float64(timeUsed) * 60
		}

		stats := fmt.Sprintf(`
  Questions: %d answered
  Time Used: %ds
  Speed: %.1f questions/min
  XP Earned: %s
  Best Combo: %dx
`,
			questionsAnswered,
			timeUsed,
			qpm,
			XPStyle.Render(fmt.Sprintf("+%d", m.XPEarned)),
			m.Combo,
		)

		switch {
		case qpm > 10:
			motivation = AccentStyle.Render("⚡ LIGHTNING FAST! You're a terminal ninja!")
		case qpm > 5:
			motivation = AccentStyle.Render("🐢💨 Turbo turtle mode activated!")
		default:
			motivation = MutedStyle.Render("🐢 Speed comes with practice. Keep at it!")
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

	// Regular lesson results
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

	switch {
	case m.XPEarned > 50:
		motivation = AccentStyle.Render("🐢 Turtle-ific! You're on fire!")
	case m.XPEarned > 20:
		motivation = AccentStyle.Render("🐢 Slow and steady wins the race!")
	default:
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
