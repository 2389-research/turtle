// ABOUTME: Tests for the skill graph and mastery tracking system
// ABOUTME: Validates prerequisites, unlocking, and skill strength decay mechanics

package skills

import (
	"testing"
)

func TestNewSkillGraph(t *testing.T) {
	graph := NewSkillGraph()

	if graph == nil {
		t.Fatal("expected non-nil skill graph")
	}
	if len(graph.Skills) != 0 {
		t.Errorf("expected empty skills map, got %d skills", len(graph.Skills))
	}
}

func TestAddSkill(t *testing.T) {
	graph := NewSkillGraph()

	skill := &Skill{
		ID:          "pwd",
		Name:        "Print Working Directory",
		Description: "Learn where you are in the filesystem",
		Category:    CategoryNavigation,
	}

	graph.AddSkill(skill)

	if len(graph.Skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(graph.Skills))
	}
	if graph.Skills["pwd"] != skill {
		t.Error("skill not found in graph")
	}
}

func TestSkillPrerequisites(t *testing.T) {
	graph := NewSkillGraph()

	// Add prerequisite skill
	pwd := &Skill{ID: "pwd", Name: "pwd", Category: CategoryNavigation}
	graph.AddSkill(pwd)

	// Add skill with prerequisite
	cd := &Skill{
		ID:            "cd",
		Name:          "cd",
		Category:      CategoryNavigation,
		Prerequisites: []string{"pwd"},
	}
	graph.AddSkill(cd)

	// pwd has no prerequisites
	prereqs := graph.GetPrerequisites("pwd")
	if len(prereqs) != 0 {
		t.Errorf("expected pwd to have no prerequisites, got %d", len(prereqs))
	}

	// cd requires pwd
	prereqs = graph.GetPrerequisites("cd")
	if len(prereqs) != 1 || prereqs[0] != "pwd" {
		t.Errorf("expected cd to require pwd, got %v", prereqs)
	}
}

func TestSkillUnlocking(t *testing.T) {
	graph := NewSkillGraph()

	// Build a small skill tree
	graph.AddSkill(&Skill{ID: "pwd", Name: "pwd", Category: CategoryNavigation})
	graph.AddSkill(&Skill{
		ID:            "cd",
		Name:          "cd",
		Category:      CategoryNavigation,
		Prerequisites: []string{"pwd"},
	})
	graph.AddSkill(&Skill{
		ID:            "cd-relative",
		Name:          "Relative paths",
		Category:      CategoryNavigation,
		Prerequisites: []string{"cd"},
	})

	progress := NewUserProgress()

	// Initially only pwd is unlocked (no prerequisites)
	unlocked := graph.GetUnlockedSkills(progress)
	if len(unlocked) != 1 || unlocked[0] != "pwd" {
		t.Errorf("expected only pwd unlocked, got %v", unlocked)
	}

	// Master pwd (strength > unlock threshold)
	progress.SetStrength("pwd", UnlockThreshold+0.1)

	// Now cd should be unlocked
	unlocked = graph.GetUnlockedSkills(progress)
	if len(unlocked) != 2 {
		t.Errorf("expected pwd and cd unlocked, got %v", unlocked)
	}

	// cd-relative still locked (cd not mastered)
	if contains(unlocked, "cd-relative") {
		t.Error("cd-relative should still be locked")
	}
}

func TestCrackingSkills(t *testing.T) {
	progress := NewUserProgress()

	// Learn a skill to high strength
	progress.SetStrength("pwd", 0.9)

	// Simulate time passing - strength should decay
	progress.SimulateDecay("pwd", 30) // 30 days

	strength := progress.GetStrength("pwd")
	if strength >= 0.9 {
		t.Errorf("expected strength to decay, still at %.2f", strength)
	}

	// Check if skill is "cracking" (below crack threshold)
	if !progress.IsCracking("pwd") {
		t.Error("skill should be cracking after significant decay")
	}
}

func TestSkillCategories(t *testing.T) {
	graph := NewSkillGraph()

	graph.AddSkill(&Skill{ID: "pwd", Category: CategoryNavigation})
	graph.AddSkill(&Skill{ID: "ls", Category: CategoryNavigation})
	graph.AddSkill(&Skill{ID: "tmux-new", Category: CategoryTmuxBasics})
	graph.AddSkill(&Skill{ID: "tmux-prefix", Category: CategoryTmuxBasics})

	navSkills := graph.GetSkillsByCategory(CategoryNavigation)
	if len(navSkills) != 2 {
		t.Errorf("expected 2 navigation skills, got %d", len(navSkills))
	}

	tmuxSkills := graph.GetSkillsByCategory(CategoryTmuxBasics)
	if len(tmuxSkills) != 2 {
		t.Errorf("expected 2 tmux skills, got %d", len(tmuxSkills))
	}
}

func TestCategoryUnlocking(t *testing.T) {
	graph := NewSkillGraph()

	// Navigation skills
	graph.AddSkill(&Skill{ID: "pwd", Category: CategoryNavigation})
	graph.AddSkill(&Skill{ID: "ls", Category: CategoryNavigation})
	graph.AddSkill(&Skill{ID: "cd", Category: CategoryNavigation, Prerequisites: []string{"pwd"}})

	// Tmux category requires navigation mastery
	graph.AddSkill(&Skill{
		ID:                "tmux-why",
		Category:          CategoryTmuxBasics,
		RequiresCategory:  CategoryNavigation,
		CategoryThreshold: 0.7, // Need 70% avg strength in navigation
	})

	progress := NewUserProgress()

	// Tmux not unlocked yet
	unlocked := graph.GetUnlockedSkills(progress)
	if contains(unlocked, "tmux-why") {
		t.Error("tmux-why should not be unlocked without navigation mastery")
	}

	// Master navigation skills
	progress.SetStrength("pwd", 0.8)
	progress.SetStrength("ls", 0.8)
	progress.SetStrength("cd", 0.7)

	// Now tmux should unlock
	unlocked = graph.GetUnlockedSkills(progress)
	if !contains(unlocked, "tmux-why") {
		t.Errorf("tmux-why should be unlocked now, got %v", unlocked)
	}
}

func TestXPAndLevels(t *testing.T) {
	progress := NewUserProgress()

	if progress.Level != 1 {
		t.Errorf("expected starting level 1, got %d", progress.Level)
	}
	if progress.XP != 0 {
		t.Errorf("expected starting XP 0, got %d", progress.XP)
	}

	// Award XP
	progress.AddXP(50)
	if progress.XP != 50 {
		t.Errorf("expected 50 XP, got %d", progress.XP)
	}

	// Level up (assuming 100 XP per level)
	progress.AddXP(60)
	if progress.Level != 2 {
		t.Errorf("expected level 2 after 110 XP, got level %d", progress.Level)
	}
}

func TestStreaks(t *testing.T) {
	progress := NewUserProgress()

	if progress.CurrentStreak != 0 {
		t.Errorf("expected starting streak 0, got %d", progress.CurrentStreak)
	}

	// Record activity today
	progress.RecordActivity()
	if progress.CurrentStreak != 1 {
		t.Errorf("expected streak 1 after activity, got %d", progress.CurrentStreak)
	}

	// Simulate next day activity
	progress.SimulateNextDay()
	progress.RecordActivity()
	if progress.CurrentStreak != 2 {
		t.Errorf("expected streak 2, got %d", progress.CurrentStreak)
	}

	// Simulate missing a day
	progress.SimulateNextDay()
	progress.SimulateNextDay() // Skip a day
	progress.RecordActivity()
	if progress.CurrentStreak != 1 {
		t.Errorf("expected streak reset to 1, got %d", progress.CurrentStreak)
	}
}

func TestDueSkills(t *testing.T) {
	graph := NewSkillGraph()
	graph.AddSkill(&Skill{ID: "pwd", Category: CategoryNavigation})
	graph.AddSkill(&Skill{ID: "ls", Category: CategoryNavigation})

	progress := NewUserProgress()

	// Practice skills
	progress.Practice("pwd", 5) // Perfect
	progress.Practice("ls", 5)

	// Nothing due right after practice
	due := graph.GetDueSkills(progress)
	if len(due) != 0 {
		t.Errorf("expected no skills due right after practice, got %v", due)
	}

	// Simulate time passing
	progress.SimulateDecay("pwd", 2) // 2 days
	progress.SimulateDecay("ls", 10) // 10 days

	due = graph.GetDueSkills(progress)
	// ls should be due (longer since review)
	if !contains(due, "ls") {
		t.Errorf("expected ls to be due, got %v", due)
	}
}

// Helper.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
