// ABOUTME: Pedagogically-organized skill graph for mission-based learning
// ABOUTME: Follows Duolingo-style progressive levels with clear learning goals

package tui

import "github.com/2389-research/turtle/internal/skills"

// buildPedagogicalSkillGraph creates the learning curriculum organized by level
// Level 0: Orientation - "I can move without fear"
// Level 1: Reading - "I can find things"
// Level 2: File ops - "I can create and organize"
// Level 3: Search - "I can find anything"
// Level 4: Tmux - "I can work across connections"
// Level 5: Muscle memory - "I'm fast and confident".
//
//nolint:funlen // Data definition function
func buildPedagogicalSkillGraph() *skills.SkillGraph {
	graph := skills.NewSkillGraph()

	// ===================
	// LEVEL 0: ORIENTATION
	// Goal: "I can move without fear"
	// ===================

	graph.AddSkill(&skills.Skill{
		ID:          "pwd",
		Name:        "pwd",
		Description: "Where am I?",
		Category:    skills.CategoryNavigation,
	})

	graph.AddSkill(&skills.Skill{
		ID:            "ls",
		Name:          "ls",
		Description:   "What's here?",
		Category:      skills.CategoryNavigation,
		Prerequisites: []string{"pwd"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "cd",
		Name:          "cd",
		Description:   "Go somewhere",
		Category:      skills.CategoryNavigation,
		Prerequisites: []string{"ls"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "cd-parent",
		Name:          "cd ..",
		Description:   "Go back up",
		Category:      skills.CategoryNavigation,
		Prerequisites: []string{"cd"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tab-completion",
		Name:          "Tab key",
		Description:   "Autocomplete paths",
		Category:      skills.CategoryNavigation,
		Prerequisites: []string{"cd"},
	})

	// ===================
	// LEVEL 1: READING THE FILESYSTEM
	// Goal: "I can find things"
	// ===================

	graph.AddSkill(&skills.Skill{
		ID:            "ls-hidden",
		Name:          "ls -a",
		Description:   "See hidden files",
		Category:      skills.CategoryNavigation,
		Prerequisites: []string{"ls"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "ls-details",
		Name:          "ls -l",
		Description:   "See file details",
		Category:      skills.CategoryNavigation,
		Prerequisites: []string{"ls"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "paths",
		Name:          "Paths",
		Description:   ". and .. and /",
		Category:      skills.CategoryNavigation,
		Prerequisites: []string{"cd-parent"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "cat",
		Name:          "cat",
		Description:   "Read file contents",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"ls"},
	})

	// ===================
	// LEVEL 2: FILE OPERATIONS
	// Goal: "I can create and organize"
	// ===================

	graph.AddSkill(&skills.Skill{
		ID:            "mkdir",
		Name:          "mkdir",
		Description:   "Create folders",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"cd"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "touch",
		Name:          "touch",
		Description:   "Create files",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"mkdir"},
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
		Description:   "Move and rename",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"cp"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "rm",
		Name:          "rm",
		Description:   "Delete files (careful!)",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"mv"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "less",
		Name:          "less",
		Description:   "Page through files",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"cat"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "redirect",
		Name:          "echo > file",
		Description:   "Write to files",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"cat"},
	})

	// ===================
	// LEVEL 3: SEARCH + INSPECTION
	// Goal: "I can find anything"
	// ===================

	graph.AddSkill(&skills.Skill{
		ID:            "grep",
		Name:          "grep",
		Description:   "Search in files",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"cat"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "find",
		Name:          "find",
		Description:   "Find files by name",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"ls-hidden"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "pipes",
		Name:          "Pipes |",
		Description:   "Chain commands",
		Category:      skills.CategoryFileOps,
		Prerequisites: []string{"grep"},
	})

	// ===================
	// LEVEL 4: TMUX
	// Goal: "I can work across connections"
	// ===================

	graph.AddSkill(&skills.Skill{
		ID:                "tmux-why",
		Name:              "Why tmux?",
		Description:       "Sessions survive disconnect",
		Category:          skills.CategoryTmuxBasics,
		RequiresCategory:  skills.CategoryNavigation,
		CategoryThreshold: 0.5,
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-new",
		Name:          "tmux new",
		Description:   "Start a session",
		Category:      skills.CategoryTmuxBasics,
		Prerequisites: []string{"tmux-why"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-detach",
		Name:          "Detach",
		Description:   "Leave session running",
		Category:      skills.CategoryTmuxBasics,
		Prerequisites: []string{"tmux-new"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-attach",
		Name:          "Attach",
		Description:   "Reconnect to session",
		Category:      skills.CategoryTmuxBasics,
		Prerequisites: []string{"tmux-detach"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-split",
		Name:          "Split panes",
		Description:   "Divide the screen",
		Category:      skills.CategoryTmuxPanes,
		Prerequisites: []string{"tmux-attach"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-pane-nav",
		Name:          "Pane navigation",
		Description:   "Move between panes",
		Category:      skills.CategoryTmuxPanes,
		Prerequisites: []string{"tmux-split"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "tmux-window",
		Name:          "Windows",
		Description:   "Multiple workspaces",
		Category:      skills.CategoryTmuxWindows,
		Prerequisites: []string{"tmux-pane-nav"},
	})

	// ===================
	// LEVEL 5: MUSCLE MEMORY
	// Goal: "I'm fast and confident"
	// ===================

	graph.AddSkill(&skills.Skill{
		ID:            "history",
		Name:          "History",
		Description:   "Arrow keys & !!",
		Category:      skills.CategoryNavigation,
		Prerequisites: []string{"cd"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "ctrl-r",
		Name:          "Ctrl-R",
		Description:   "Search history",
		Category:      skills.CategoryNavigation,
		Prerequisites: []string{"history"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "ctrl-c",
		Name:          "Ctrl-C",
		Description:   "Cancel/interrupt",
		Category:      skills.CategoryNavigation,
		Prerequisites: []string{"pwd"},
	})

	graph.AddSkill(&skills.Skill{
		ID:            "ctrl-d",
		Name:          "Ctrl-D",
		Description:   "Exit/EOF",
		Category:      skills.CategoryNavigation,
		Prerequisites: []string{"ctrl-c"},
	})

	return graph
}

// GetSkillLevel returns the pedagogical level for a skill (0-5).
func GetSkillLevel(skillID string) int {
	levelMap := map[string]int{
		// Level 0: Orientation
		"pwd": 0, "ls": 0, "cd": 0, "cd-parent": 0, "tab-completion": 0,
		// Level 1: Reading
		"ls-hidden": 1, "ls-details": 1, "paths": 1, "cat": 1,
		// Level 2: File ops
		"mkdir": 2, "touch": 2, "cp": 2, "mv": 2, "rm": 2, "less": 2, "redirect": 2,
		// Level 3: Search
		"grep": 3, "find": 3, "pipes": 3,
		// Level 4: Tmux
		"tmux-why": 4, "tmux-new": 4, "tmux-detach": 4, "tmux-attach": 4,
		"tmux-split": 4, "tmux-pane-nav": 4, "tmux-window": 4,
		// Level 5: Muscle memory
		"history": 5, "ctrl-r": 5, "ctrl-c": 5, "ctrl-d": 5,
	}

	if level, ok := levelMap[skillID]; ok {
		return level
	}
	return -1
}

// LevelNames returns human-readable names for each level.
func LevelNames() map[int]string {
	return map[int]string{
		0: "Orientation",
		1: "Reading",
		2: "File Operations",
		3: "Search & Inspect",
		4: "Tmux",
		5: "Muscle Memory",
	}
}

// LevelGoals returns the learning goal for each level.
func LevelGoals() map[int]string {
	return map[int]string{
		0: "I can move without fear",
		1: "I can find things",
		2: "I can create and organize",
		3: "I can find anything",
		4: "I can work across connections",
		5: "I'm fast and confident",
	}
}
