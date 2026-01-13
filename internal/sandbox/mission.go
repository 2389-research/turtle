// ABOUTME: Mission-based learning system for terminal commands
// ABOUTME: Defines goal-oriented challenges with immediate feedback

package sandbox

import (
	"fmt"
	"strings"

	"github.com/2389-research/turtle/internal/content"
)

// goalContext adapts Filesystem to the GoalEvaluator interface, adding command tracking.
type goalContext struct {
	fs          *Filesystem
	lastCommand string
}

func (g *goalContext) Pwd() string {
	return g.fs.Pwd()
}

func (g *goalContext) Exists(path string) bool {
	return g.fs.Exists(path)
}

func (g *goalContext) IsDir(path string) bool {
	return g.fs.IsDir(path)
}

func (g *goalContext) ReadFile(path string) (string, error) {
	return g.fs.ReadFile(path)
}

func (g *goalContext) LastCommand() string {
	return g.lastCommand
}

// Mission represents a goal-based learning challenge.
type Mission struct {
	ID          string
	SkillID     string                           // Which skill this teaches
	Level       int                              // 0-5 difficulty
	Title       string                           // Short mission name
	Briefing    string                           // What the learner needs to do
	Hint        string                           // Help if stuck
	Setup       func(*Filesystem)                // Prepares the filesystem
	Goal        func(content.GoalEvaluator) bool // Returns true if mission complete
	Explanation string                           // Shown after success
	Commands    []string                         // Commands that could solve this (for reference)
}

// MissionResult represents the outcome of a command.
type MissionResult struct {
	Output    string // What to show the learner
	Success   bool   // Command executed successfully
	Error     string // Error message if failed
	Completed bool   // Mission goal achieved
}

// TmuxState tracks simulated tmux session state.
type TmuxState struct {
	InSession     bool
	SessionName   string
	Panes         int
	Windows       int
	CurrentPane   int
	CurrentWindow int
	Detached      bool
}

// MissionRunner executes missions in the sandbox.
type MissionRunner struct {
	FS        *Filesystem
	Mission   *Mission
	InitialFS *Filesystem // For reset
	Attempts  int
	Completed bool
	History   []string  // Commands entered
	Tmux      TmuxState // Tmux simulation state
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

	// Check if mission is complete using goalContext for command tracking
	if r.Mission.Goal != nil {
		ctx := &goalContext{fs: r.FS, lastCommand: input}
		if r.Mission.Goal(ctx) {
			result.Completed = true
			r.Completed = true
		}
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
			Output:  "Available: pwd, ls, cd, mkdir, touch, cat, cp, mv, rm, echo, grep, find, clear, tmux",
			Success: true,
		}

	case "tmux":
		return r.executeTmux(args)

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

// executeTmux handles tmux commands.
//
//nolint:funlen,gocognit,gocyclo // Tmux command dispatcher requires many branches
func (r *MissionRunner) executeTmux(args []string) MissionResult {
	if len(args) == 0 {
		// Just 'tmux' starts a new session
		if r.Tmux.InSession {
			return MissionResult{Error: "sessions should be nested with care, unset $TMUX to force"}
		}
		r.Tmux.InSession = true
		r.Tmux.SessionName = "0"
		r.Tmux.Panes = 1
		r.Tmux.Windows = 1
		r.Tmux.CurrentPane = 0
		r.Tmux.CurrentWindow = 0
		r.Tmux.Detached = false
		return MissionResult{
			Output:  "[new session 0 created]",
			Success: true,
		}
	}

	subCmd := args[0]
	subArgs := args[1:]

	switch subCmd {
	case "new", "new-session":
		if r.Tmux.InSession && !r.Tmux.Detached {
			return MissionResult{Error: "sessions should be nested with care, unset $TMUX to force"}
		}
		sessionName := "0"
		for i, arg := range subArgs {
			if arg == "-s" && i+1 < len(subArgs) {
				sessionName = subArgs[i+1]
			}
		}
		r.Tmux.InSession = true
		r.Tmux.SessionName = sessionName
		r.Tmux.Panes = 1
		r.Tmux.Windows = 1
		r.Tmux.CurrentPane = 0
		r.Tmux.CurrentWindow = 0
		r.Tmux.Detached = false
		return MissionResult{
			Output:  fmt.Sprintf("[new session %s created]", sessionName),
			Success: true,
		}

	case "attach", "attach-session", "a":
		if r.Tmux.InSession && !r.Tmux.Detached {
			return MissionResult{Error: "sessions should be nested with care"}
		}
		if !r.Tmux.Detached && r.Tmux.SessionName == "" {
			return MissionResult{Error: "no sessions"}
		}
		r.Tmux.Detached = false
		r.Tmux.InSession = true
		return MissionResult{
			Output:  fmt.Sprintf("[attached to session %s]", r.Tmux.SessionName),
			Success: true,
		}

	case "detach", "detach-client", "d":
		if !r.Tmux.InSession || r.Tmux.Detached {
			return MissionResult{Error: "no current client"}
		}
		r.Tmux.Detached = true
		return MissionResult{
			Output:  fmt.Sprintf("[detached (from session %s)]", r.Tmux.SessionName),
			Success: true,
		}

	case "ls", "list-sessions":
		if r.Tmux.SessionName == "" {
			return MissionResult{
				Output:  "no server running",
				Success: true,
			}
		}
		attached := ""
		if r.Tmux.InSession && !r.Tmux.Detached {
			attached = " (attached)"
		}
		return MissionResult{
			Output:  fmt.Sprintf("%s: %d windows%s", r.Tmux.SessionName, r.Tmux.Windows, attached),
			Success: true,
		}

	case "split-window", "split":
		if !r.Tmux.InSession || r.Tmux.Detached {
			return MissionResult{Error: "no current session"}
		}
		r.Tmux.Panes++
		direction := "horizontally"
		for _, arg := range subArgs {
			if arg == "-v" {
				direction = "vertically"
			}
		}
		return MissionResult{
			Output:  fmt.Sprintf("[split %s, now %d panes]", direction, r.Tmux.Panes),
			Success: true,
		}

	case "select-pane":
		if !r.Tmux.InSession || r.Tmux.Detached {
			return MissionResult{Error: "no current session"}
		}
		// Handle -L, -R, -U, -D for direction
		direction := "next"
		for _, arg := range subArgs {
			switch arg {
			case "-L":
				direction = "left"
			case "-R":
				direction = "right"
			case "-U":
				direction = "up"
			case "-D":
				direction = "down"
			}
		}
		r.Tmux.CurrentPane = (r.Tmux.CurrentPane + 1) % r.Tmux.Panes
		return MissionResult{
			Output:  fmt.Sprintf("[moved %s to pane %d]", direction, r.Tmux.CurrentPane),
			Success: true,
		}

	case "new-window":
		if !r.Tmux.InSession || r.Tmux.Detached {
			return MissionResult{Error: "no current session"}
		}
		r.Tmux.Windows++
		r.Tmux.CurrentWindow = r.Tmux.Windows - 1
		return MissionResult{
			Output:  fmt.Sprintf("[new window %d created]", r.Tmux.CurrentWindow),
			Success: true,
		}

	case "select-window":
		if !r.Tmux.InSession || r.Tmux.Detached {
			return MissionResult{Error: "no current session"}
		}
		// Handle -n (next) and -p (previous)
		for _, arg := range subArgs {
			switch arg {
			case "-n":
				r.Tmux.CurrentWindow = (r.Tmux.CurrentWindow + 1) % r.Tmux.Windows
			case "-p":
				r.Tmux.CurrentWindow = (r.Tmux.CurrentWindow - 1 + r.Tmux.Windows) % r.Tmux.Windows
			}
		}
		return MissionResult{
			Output:  fmt.Sprintf("[switched to window %d]", r.Tmux.CurrentWindow),
			Success: true,
		}

	case "kill-session":
		if r.Tmux.SessionName == "" {
			return MissionResult{Error: "no sessions"}
		}
		name := r.Tmux.SessionName
		r.Tmux = TmuxState{}
		return MissionResult{
			Output:  fmt.Sprintf("[killed session %s]", name),
			Success: true,
		}

	default:
		return MissionResult{
			Error: fmt.Sprintf("tmux: unknown command: %s", subCmd),
		}
	}
}

// InTmuxSession returns true if currently in a tmux session.
func (r *MissionRunner) InTmuxSession() bool {
	return r.Tmux.InSession && !r.Tmux.Detached
}

// GetTmuxStatus returns a description of the current tmux state.
func (r *MissionRunner) GetTmuxStatus() string {
	if !r.Tmux.InSession && r.Tmux.SessionName == "" {
		return "not in tmux"
	}
	if r.Tmux.Detached {
		return fmt.Sprintf("detached from session '%s'", r.Tmux.SessionName)
	}
	return fmt.Sprintf("session '%s' [%d windows, %d panes]",
		r.Tmux.SessionName, r.Tmux.Windows, r.Tmux.Panes)
}
