// ABOUTME: Default skill graph configuration for Turtle
// ABOUTME: Defines the learning curriculum - skills, prerequisites, and progression

package tui

import "github.com/2389-research/turtle/internal/skills"

// buildDefaultSkillGraph creates the full learning curriculum
func buildDefaultSkillGraph() *skills.SkillGraph {
	graph := skills.NewSkillGraph()

	// ===================
	// NAVIGATION (Level 0-1)
	// ===================

	graph.AddSkill(&skills.Skill{
		ID:          "pwd",
		Name:        "pwd",
		Description: "Print working directory - know where you are",
		Category:    skills.CategoryNavigation,
	})

	graph.AddSkill(&skills.Skill{
		ID:            "ls",
		Name:          "ls",
		Description:   "List directory contents - see what's here",
		Category:      skills.CategoryNavigation,
		Prerequisites: []string{"pwd"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "cd",
		Name:          "cd",
		Description:   "Change directory - move around",
		Category:      skills.CategoryNavigation,
		Prerequisites: []string{"ls"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "cd-relative",
		Name:          "cd (relative)",
		Description:   "Navigate using . and ..",
		Category:      skills.CategoryNavigation,
		Prerequisites: []string{"cd"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tab-completion",
		Name:          "Tab completion",
		Description:   "Let the shell complete paths for you",
		Category:      skills.CategoryNavigation,
		Prerequisites: []string{"cd"},
	})

	// ===================
	// FILE OPERATIONS (Level 2)
	// ===================

	graph.AddSkill(&skills.Skill{
		ID:            "mkdir",
		Name:          "mkdir",
		Description:   "Make directories",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"cd"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "touch",
		Name:          "touch",
		Description:   "Create empty files",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"ls"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "rm",
		Name:          "rm",
		Description:   "Remove files (careful!)",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"touch"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "cp",
		Name:          "cp",
		Description:   "Copy files",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"touch"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "mv",
		Name:          "mv",
		Description:   "Move or rename files",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"cp"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "cat",
		Name:          "cat",
		Description:   "Display file contents",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"touch"},
	})

	// ===================
	// TMUX BASICS (Level 3)
	// Requires navigation mastery
	// ===================

	graph.AddSkill(&skills.Skill{
		ID:                "tmux-why",
		Name:              "Why tmux?",
		Description:       "Understanding terminal multiplexing",
		Category:          skills.CategoryTmuxBasics,
		RequiresCategory:  skills.CategoryNavigation,
		CategoryThreshold: 0.6,
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-new",
		Name:          "tmux new",
		Description:   "Create a new session",
		Category:      skills.CategoryTmuxBasics,
		Prerequisites: []string{"tmux-why"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-prefix",
		Name:          "Prefix key",
		Description:   "The Ctrl-b gateway",
		Category:      skills.CategoryTmuxBasics,
		Prerequisites: []string{"tmux-new"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-detach",
		Name:          "Detach",
		Description:   "Leave session running in background",
		Category:      skills.CategoryTmuxBasics,
		Prerequisites: []string{"tmux-prefix"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-attach",
		Name:          "Attach",
		Description:   "Reconnect to a session",
		Category:      skills.CategoryTmuxBasics,
		Prerequisites: []string{"tmux-detach"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-list",
		Name:          "List sessions",
		Description:   "See all running sessions",
		Category:      skills.CategoryTmuxBasics,
		Prerequisites: []string{"tmux-attach"},
	})

	// ===================
	// TMUX PANES (Level 4)
	// ===================

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-split-h",
		Name:          "Split horizontal",
		Description:   "Divide pane side by side",
		Category:      skills.CategoryTmuxPanes,
		Prerequisites: []string{"tmux-prefix"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-split-v",
		Name:          "Split vertical",
		Description:   "Divide pane top/bottom",
		Category:      skills.CategoryTmuxPanes,
		Prerequisites: []string{"tmux-split-h"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-pane-nav",
		Name:          "Pane navigation",
		Description:   "Move between panes",
		Category:      skills.CategoryTmuxPanes,
		Prerequisites: []string{"tmux-split-v"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-pane-resize",
		Name:          "Resize panes",
		Description:   "Adjust pane sizes",
		Category:      skills.CategoryTmuxPanes,
		Prerequisites: []string{"tmux-pane-nav"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-pane-close",
		Name:          "Close pane",
		Description:   "Close current pane",
		Category:      skills.CategoryTmuxPanes,
		Prerequisites: []string{"tmux-pane-nav"},
	})

	// ===================
	// TMUX WINDOWS (Level 5)
	// ===================

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-window-new",
		Name:          "New window",
		Description:   "Create a new window",
		Category:      skills.CategoryTmuxWindows,
		Prerequisites: []string{"tmux-pane-nav"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-window-nav",
		Name:          "Window navigation",
		Description:   "Switch between windows",
		Category:      skills.CategoryTmuxWindows,
		Prerequisites: []string{"tmux-window-new"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-window-rename",
		Name:          "Rename window",
		Description:   "Give windows meaningful names",
		Category:      skills.CategoryTmuxWindows,
		Prerequisites: []string{"tmux-window-nav"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-window-close",
		Name:          "Close window",
		Description:   "Close current window",
		Category:      skills.CategoryTmuxWindows,
		Prerequisites: []string{"tmux-window-nav"},
	})

	return graph
}
