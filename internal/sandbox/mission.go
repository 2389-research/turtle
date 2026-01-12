// ABOUTME: Mission-based learning system for terminal commands
// ABOUTME: Defines goal-oriented challenges with immediate feedback

package sandbox

import "strings"

// Mission represents a goal-based learning challenge.
type Mission struct {
	ID          string
	SkillID     string                 // Which skill this teaches
	Level       int                    // 0-5 difficulty
	Title       string                 // Short mission name
	Briefing    string                 // What the learner needs to do
	Hint        string                 // Help if stuck
	Setup       func(*Filesystem)      // Prepares the filesystem
	Goal        func(*Filesystem) bool // Returns true if mission complete
	Explanation string                 // Shown after success
	Commands    []string               // Commands that could solve this (for reference)
}

// MissionResult represents the outcome of a command.
type MissionResult struct {
	Output    string // What to show the learner
	Success   bool   // Command executed successfully
	Error     string // Error message if failed
	Completed bool   // Mission goal achieved
}

// MissionRunner executes missions in the sandbox.
type MissionRunner struct {
	FS        *Filesystem
	Mission   *Mission
	InitialFS *Filesystem // For reset
	Attempts  int
	Completed bool
	History   []string // Commands entered
}

// NewMissionRunner creates a runner for a mission.
func NewMissionRunner(m *Mission) *MissionRunner {
	// Start with a default filesystem that has common directories
	fs := NewDefaultFilesystem()

	// Run mission-specific setup
	if m.Setup != nil {
		m.Setup(fs)
	}

	return &MissionRunner{
		FS:        fs,
		Mission:   m,
		InitialFS: fs.Clone(),
		Attempts:  0,
		Completed: false,
		History:   []string{},
	}
}

// Reset restores the filesystem to initial state.
func (r *MissionRunner) Reset() {
	r.FS = r.InitialFS.Clone()
	r.Attempts = 0
	r.History = []string{}
}

// Execute runs a command and returns the result.
func (r *MissionRunner) Execute(input string) MissionResult {
	r.Attempts++
	r.History = append(r.History, input)

	input = strings.TrimSpace(input)
	if input == "" {
		return MissionResult{Output: "", Success: true}
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return MissionResult{Output: "", Success: true}
	}

	cmd := parts[0]
	args := parts[1:]

	result := r.executeCommand(cmd, args)

	// Check if mission is complete
	if r.Mission.Goal != nil && r.Mission.Goal(r.FS) {
		result.Completed = true
		r.Completed = true
	}

	return result
}

// executeCommand handles individual commands.
//
//nolint:gocognit,gocyclo,funlen // Command dispatcher requires many branches
func (r *MissionRunner) executeCommand(cmd string, args []string) MissionResult {
	switch cmd {
	case "pwd":
		return MissionResult{
			Output:  r.FS.Pwd(),
			Success: true,
		}

	case "ls":
		path := ""
		showHidden := false
		for _, arg := range args {
			if arg == "-a" || arg == "-la" || arg == "-al" {
				showHidden = true
			} else if !strings.HasPrefix(arg, "-") {
				path = arg
			}
		}
		files, err := r.FS.Ls(path, showHidden)
		if err != nil {
			return MissionResult{Error: err.Error()}
		}
		return MissionResult{
			Output:  strings.Join(files, "  "),
			Success: true,
		}

	case "cd":
		path := ""
		if len(args) > 0 {
			path = args[0]
		}
		err := r.FS.Cd(path)
		if err != nil {
			return MissionResult{Error: err.Error()}
		}
		return MissionResult{Success: true}

	case "mkdir":
		if len(args) == 0 {
			return MissionResult{Error: "mkdir: missing operand"}
		}
		for _, path := range args {
			if !strings.HasPrefix(path, "-") {
				if err := r.FS.Mkdir(path); err != nil {
					return MissionResult{Error: err.Error()}
				}
			}
		}
		return MissionResult{Success: true}

	case "touch":
		if len(args) == 0 {
			return MissionResult{Error: "touch: missing file operand"}
		}
		for _, path := range args {
			if err := r.FS.Touch(path); err != nil {
				return MissionResult{Error: err.Error()}
			}
		}
		return MissionResult{Success: true}

	case "cat":
		if len(args) == 0 {
			return MissionResult{Error: "cat: missing file operand"}
		}
		var contents []string
		for _, path := range args {
			content, err := r.FS.Cat(path)
			if err != nil {
				return MissionResult{Error: err.Error()}
			}
			contents = append(contents, content)
		}
		return MissionResult{
			Output:  strings.Join(contents, "\n"),
			Success: true,
		}

	case "cp":
		if len(args) < 2 {
			return MissionResult{Error: "cp: missing destination file operand"}
		}
		src := args[len(args)-2]
		dst := args[len(args)-1]
		if err := r.FS.Cp(src, dst); err != nil {
			return MissionResult{Error: err.Error()}
		}
		return MissionResult{Success: true}

	case "mv":
		if len(args) < 2 {
			return MissionResult{Error: "mv: missing destination file operand"}
		}
		src := args[len(args)-2]
		dst := args[len(args)-1]
		if err := r.FS.Mv(src, dst); err != nil {
			return MissionResult{Error: err.Error()}
		}
		return MissionResult{Success: true}

	case "rm":
		if len(args) == 0 {
			return MissionResult{Error: "rm: missing operand"}
		}
		for _, path := range args {
			if !strings.HasPrefix(path, "-") {
				if err := r.FS.Rm(path); err != nil {
					return MissionResult{Error: err.Error()}
				}
			}
		}
		return MissionResult{Success: true}

	case "grep":
		if len(args) < 2 {
			return MissionResult{Error: "grep: missing pattern or file"}
		}
		pattern := args[0]
		file := args[1]
		matches, err := r.FS.Grep(pattern, file)
		if err != nil {
			return MissionResult{Error: err.Error()}
		}
		return MissionResult{
			Output:  strings.Join(matches, "\n"),
			Success: true,
		}

	case "find":
		if len(args) < 3 {
			return MissionResult{Error: "find: missing arguments"}
		}
		path := args[0]
		pattern := ""
		for i, arg := range args {
			if arg == "-name" && i+1 < len(args) {
				pattern = strings.Trim(args[i+1], "\"'")
			}
		}
		if pattern == "" {
			return MissionResult{Error: "find: missing -name pattern"}
		}
		results, err := r.FS.Find(path, pattern)
		if err != nil {
			return MissionResult{Error: err.Error()}
		}
		return MissionResult{
			Output:  strings.Join(results, "\n"),
			Success: true,
		}

	case "echo":
		// Handle redirects
		outputFile := ""
		appendMode := false
		var echoArgs []string

		for i := 0; i < len(args); i++ {
			switch {
			case args[i] == ">" && i+1 < len(args):
				outputFile = args[i+1]
				i++
			case args[i] == ">>" && i+1 < len(args):
				outputFile = args[i+1]
				appendMode = true
				i++
			case strings.HasPrefix(args[i], ">"):
				outputFile = strings.TrimPrefix(args[i], ">")
			default:
				echoArgs = append(echoArgs, args[i])
			}
		}

		output := strings.Join(echoArgs, " ")

		if outputFile != "" {
			if appendMode {
				existing, _ := r.FS.ReadFile(outputFile)
				output = existing + output + "\n"
			} else {
				output += "\n"
			}
			if err := r.FS.WriteFile(outputFile, output); err != nil {
				return MissionResult{Error: err.Error()}
			}
			return MissionResult{Success: true}
		}

		return MissionResult{
			Output:  output,
			Success: true,
		}

	case "clear":
		return MissionResult{
			Output:  "\033[2J\033[H", // ANSI clear screen
			Success: true,
		}

	case "help":
		return MissionResult{
			Output:  "Available: pwd, ls, cd, mkdir, touch, cat, cp, mv, rm, echo, grep, find, clear",
			Success: true,
		}

	default:
		return MissionResult{
			Error: cmd + ": command not found",
		}
	}
}

// GetCurrentLocation returns a user-friendly description of where they are.
func (r *MissionRunner) GetCurrentLocation() string {
	path := r.FS.Pwd()
	if path == r.FS.Home {
		return "~ (your home directory)"
	}
	if strings.HasPrefix(path, r.FS.Home) {
		return "~" + strings.TrimPrefix(path, r.FS.Home)
	}
	return path
}
