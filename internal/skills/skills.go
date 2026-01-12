// ABOUTME: Skill graph and mastery tracking system for progressive learning
// ABOUTME: Manages skill prerequisites, unlocking, categories, and decay mechanics

package skills

import (
	"time"

	"github.com/2389-research/turtle/internal/srs"
)

// Skill categories for organizing the learning path
type Category string

const (
	CategoryNavigation  Category = "navigation"
	CategoryFileOps     Category = "file-operations"
	CategoryTmuxBasics  Category = "tmux-basics"
	CategoryTmuxPanes   Category = "tmux-panes"
	CategoryTmuxWindows Category = "tmux-windows"
	CategoryAdvanced    Category = "advanced"
)

// Thresholds for game mechanics
const (
	UnlockThreshold = 0.6 // 60% strength to unlock dependent skills
	CrackThreshold  = 0.4 // Below 40% = skill is "cracking"
	XPPerLevel      = 100 // XP required per level
)

// Skill represents a learnable concept in the skill tree
type Skill struct {
	ID                string
	Name              string
	Description       string
	Category          Category
	Prerequisites     []string // Skill IDs that must be mastered first
	RequiresCategory  Category // Optional: requires avg mastery in this category
	CategoryThreshold float64  // Required avg strength in RequiresCategory
}

// SkillGraph manages the entire skill tree and relationships
type SkillGraph struct {
	Skills map[string]*Skill
}

// NewSkillGraph creates an empty skill graph
func NewSkillGraph() *SkillGraph {
	return &SkillGraph{
		Skills: make(map[string]*Skill),
	}
}

// AddSkill registers a skill in the graph
func (g *SkillGraph) AddSkill(skill *Skill) {
	g.Skills[skill.ID] = skill
}

// GetPrerequisites returns the prerequisite skill IDs for a skill
func (g *SkillGraph) GetPrerequisites(skillID string) []string {
	skill, ok := g.Skills[skillID]
	if !ok {
		return nil
	}
	return skill.Prerequisites
}

// GetSkillsByCategory returns all skills in a category
func (g *SkillGraph) GetSkillsByCategory(cat Category) []string {
	var result []string
	for id, skill := range g.Skills {
		if skill.Category == cat {
			result = append(result, id)
		}
	}
	return result
}

// GetUnlockedSkills returns all skills the user can currently practice
func (g *SkillGraph) GetUnlockedSkills(progress *UserProgress) []string {
	var unlocked []string

	for id, skill := range g.Skills {
		if g.isUnlocked(skill, progress) {
			unlocked = append(unlocked, id)
		}
	}

	return unlocked
}

// isUnlocked checks if a specific skill is unlocked for the user
func (g *SkillGraph) isUnlocked(skill *Skill, progress *UserProgress) bool {
	// Check individual prerequisites
	for _, prereqID := range skill.Prerequisites {
		if progress.GetStrength(prereqID) < UnlockThreshold {
			return false
		}
	}

	// Check category requirement
	if skill.RequiresCategory != "" {
		avgStrength := g.getCategoryAvgStrength(skill.RequiresCategory, progress)
		if avgStrength < skill.CategoryThreshold {
			return false
		}
	}

	return true
}

// getCategoryAvgStrength calculates average strength across a category
func (g *SkillGraph) getCategoryAvgStrength(cat Category, progress *UserProgress) float64 {
	skills := g.GetSkillsByCategory(cat)
	if len(skills) == 0 {
		return 0
	}

	var total float64
	for _, skillID := range skills {
		total += progress.GetStrength(skillID)
	}
	return total / float64(len(skills))
}

// GetDueSkills returns skills that need review based on SRS scheduling
func (g *SkillGraph) GetDueSkills(progress *UserProgress) []string {
	var due []string

	for skillID := range g.Skills {
		card := progress.GetCard(skillID)
		if card != nil && card.IsDue() {
			due = append(due, skillID)
		}
	}

	return due
}

// UserProgress tracks a user's learning state
type UserProgress struct {
	XP            int
	Level         int
	CurrentStreak int
	BestStreak    int
	LastActive    time.Time
	Cards         map[string]*srs.Card // SRS cards for each skill

	// testNow is used for testing time-dependent behavior (nil in production)
	testNow *time.Time
}

// NewUserProgress creates a fresh user progress tracker
func NewUserProgress() *UserProgress {
	return &UserProgress{
		XP:            0,
		Level:         1,
		CurrentStreak: 0,
		BestStreak:    0,
		LastActive:    time.Time{},
		Cards:         make(map[string]*srs.Card),
		testNow:       nil,
	}
}

// now returns the current time (or simulated time for testing)
func (p *UserProgress) now() time.Time {
	if p.testNow != nil {
		return *p.testNow
	}
	return time.Now()
}

// GetStrength returns the current strength for a skill
func (p *UserProgress) GetStrength(skillID string) float64 {
	card := p.Cards[skillID]
	if card == nil {
		return 0
	}
	return card.Strength()
}

// SetStrength manually sets a skill's strength (for testing/simulation)
func (p *UserProgress) SetStrength(skillID string, strength float64) {
	card := p.getOrCreateCard(skillID)
	// Approximate the strength by setting repetitions
	// This is a simplified approach for direct strength setting
	if strength > 0 {
		card.Repetitions = int(strength * 5) // Rough mapping
		card.LastReviewed = time.Now()
	}
}

// GetCard returns the SRS card for a skill, or nil if not practiced
func (p *UserProgress) GetCard(skillID string) *srs.Card {
	return p.Cards[skillID]
}

// getOrCreateCard ensures a card exists for a skill
func (p *UserProgress) getOrCreateCard(skillID string) *srs.Card {
	if p.Cards[skillID] == nil {
		p.Cards[skillID] = srs.NewCard(skillID)
	}
	return p.Cards[skillID]
}

// Practice records a practice session for a skill
func (p *UserProgress) Practice(skillID string, grade int) {
	card := p.getOrCreateCard(skillID)
	card.Review(grade)
}

// IsCracking returns true if a skill's strength has dropped below the crack threshold
func (p *UserProgress) IsCracking(skillID string) bool {
	strength := p.GetStrength(skillID)
	card := p.Cards[skillID]

	// Must have been practiced at least once to be "cracking"
	if card == nil || card.Repetitions == 0 {
		return false
	}

	return strength < CrackThreshold
}

// SimulateDecay simulates the passage of time for testing decay mechanics
func (p *UserProgress) SimulateDecay(skillID string, days int) {
	card := p.Cards[skillID]
	if card != nil {
		card.LastReviewed = card.LastReviewed.AddDate(0, 0, -days)
	}
}

// AddXP awards experience points and handles level ups
func (p *UserProgress) AddXP(amount int) {
	p.XP += amount

	// Calculate new level
	newLevel := (p.XP / XPPerLevel) + 1
	if newLevel > p.Level {
		p.Level = newLevel
	}
}

// RecordActivity records that the user practiced today (for streaks)
func (p *UserProgress) RecordActivity() {
	now := p.now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lastActiveDay := time.Date(p.LastActive.Year(), p.LastActive.Month(), p.LastActive.Day(), 0, 0, 0, 0, p.LastActive.Location())

	if p.LastActive.IsZero() {
		// First activity ever
		p.CurrentStreak = 1
	} else {
		daysSince := int(today.Sub(lastActiveDay).Hours() / 24)

		if daysSince == 0 {
			// Already active today, no change
			return
		} else if daysSince == 1 {
			// Consecutive day
			p.CurrentStreak++
		} else {
			// Streak broken
			p.CurrentStreak = 1
		}
	}

	if p.CurrentStreak > p.BestStreak {
		p.BestStreak = p.CurrentStreak
	}

	p.LastActive = now
}

// SimulateNextDay advances the virtual clock by one day (for testing streaks)
func (p *UserProgress) SimulateNextDay() {
	current := p.now()
	next := current.AddDate(0, 0, 1)
	p.testNow = &next
}
