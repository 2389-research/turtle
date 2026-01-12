// ABOUTME: Mission definitions organized by pedagogical level
// ABOUTME: Each mission has a clear goal, constrained environment, and immediate feedback

package sandbox

import (
	"log"
	"strings"

	"github.com/2389-research/turtle/internal/content"
)

// GetAllMissions returns missions organized by level.
func GetAllMissions() map[int][]*Mission {
	rawMissions, err := content.GetRawMissions()
	if err != nil {
		log.Printf("Error loading missions from YAML: %v, falling back to legacy missions", err)
		return GetAllMissionsLegacy()
	}
	return convertMissions(rawMissions)
}

// convertMissions converts raw YAML missions to Mission types.
func convertMissions(raw []content.YAMLMission) map[int][]*Mission {
	result := make(map[int][]*Mission)

	for i := range raw {
		ym := &raw[i]
		mission, err := convertMission(ym)
		if err != nil {
			log.Printf("Error converting mission %s: %v, skipping", ym.ID, err)
			continue
		}
		result[ym.Level] = append(result[ym.Level], mission)
	}

	return result
}

func convertMission(ym *content.YAMLMission) (*Mission, error) {
	goalNode, err := content.ParseGoal(ym.Goal)
	if err != nil {
		return nil, err
	}

	// Capture setup actions for closure
	setupActions := ym.Setup

	return &Mission{
		ID:          ym.ID,
		SkillID:     ym.SkillID,
		Level:       ym.Level,
		Title:       ym.Title,
		Briefing:    ym.Briefing,
		Hint:        ym.Hint,
		Explanation: ym.Explanation,
		Commands:    ym.Commands,
		Setup: func(fs *Filesystem) {
			executeSetup(fs, setupActions)
		},
		Goal: func(fs *Filesystem) bool {
			return goalNode.Evaluate(fs)
		},
	}, nil
}

func executeSetup(fs *Filesystem, actions []content.SetupAction) {
	for _, action := range actions {
		switch {
		case action.Mkdir != "":
			_ = fs.Mkdir(action.Mkdir)
		case action.Cd != "":
			_ = fs.Cd(action.Cd)
		case action.Touch != "":
			_ = fs.Touch(action.Touch)
		case action.WriteFile != nil:
			_ = fs.WriteFile(action.WriteFile.Path, action.WriteFile.Content)
		}
	}
}

// GetAllMissionsLegacy returns missions organized by level (legacy hardcoded version).
// Kept for reference during migration.
func GetAllMissionsLegacy() map[int][]*Mission {
	return map[int][]*Mission{
		0: Level0Missions(), // Orientation
		1: Level1Missions(), // Reading the filesystem
		2: Level2Missions(), // File operations
		3: Level3Missions(), // Search + inspection
		4: Level4Missions(), // Tmux basics
		5: Level5Missions(), // Muscle memory
	}
}

// GetMissionsForSkill returns missions that teach a specific skill.
func GetMissionsForSkill(skillID string) []*Mission {
	var missions []*Mission
	for _, levelMissions := range GetAllMissions() {
		for _, m := range levelMissions {
			if m.SkillID == skillID {
				missions = append(missions, m)
			}
		}
	}
	return missions
}

// ===================
// LEVEL 0: ORIENTATION
// Goal: "I can move without fear"
// ===================

//nolint:funlen // Data definition function
func Level0Missions() []*Mission {
	return []*Mission{
		// Mission 0.1: Where am I?
		{
			ID:       "0.1-where-am-i",
			SkillID:  "pwd",
			Level:    0,
			Title:    "Where Am I?",
			Briefing: "You've just opened a terminal. You're somewhere in the filesystem, but where? Find out your current location.",
			Hint:     "There's a command that prints the working directory...",
			Setup: func(fs *Filesystem) {
				_ = fs.Cd("/home/learner/projects")
			},
			Goal: func(fs *Filesystem) bool {
				// Mission completes when they run pwd (any state is fine)
				return true // This is triggered by checking if they ran pwd
			},
			Explanation: "pwd = Print Working Directory. It shows the full path to where you are. You're in /home/learner/projects.",
			Commands:    []string{"pwd"},
		},

		// Mission 0.2: What's here?
		{
			ID:       "0.2-whats-here",
			SkillID:  "ls",
			Level:    0,
			Title:    "What's Here?",
			Briefing: "You know where you are. Now look around. What files and folders are in this directory?",
			Hint:     "Two letters. Short for 'list'.",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/projects")
				_ = fs.Touch("/home/learner/projects/readme.txt")
				_ = fs.Mkdir("/home/learner/projects/src")
				_ = fs.Mkdir("/home/learner/projects/docs")
				_ = fs.Cd("/home/learner/projects")
			},
			Goal: func(fs *Filesystem) bool {
				return true // Completes when they run ls
			},
			Explanation: "ls = list. It shows you everything in the current directory. You can see readme.txt, src/, and docs/.",
			Commands:    []string{"ls"},
		},

		// Mission 0.3: Go home
		{
			ID:       "0.3-go-home",
			SkillID:  "cd",
			Level:    0,
			Title:    "Go Home",
			Briefing: "You're in /tmp. That's not where you belong. Navigate to your home directory.",
			Hint:     "cd without any arguments takes you home. Or try cd ~",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/tmp")
				_ = fs.Cd("/tmp")
			},
			Goal: func(fs *Filesystem) bool {
				return fs.Pwd() == "/home/learner"
			},
			Explanation: "cd = change directory. By itself, cd takes you home. ~ is shorthand for your home directory.",
			Commands:    []string{"cd", "cd ~"},
		},

		// Mission 0.4: Enter a folder
		{
			ID:       "0.4-enter-folder",
			SkillID:  "cd",
			Level:    0,
			Title:    "Enter the Projects Folder",
			Briefing: "You're home. There's a 'projects' folder here. Go inside it.",
			Hint:     "cd followed by the folder name",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/projects")
				_ = fs.Cd("/home/learner")
			},
			Goal: func(fs *Filesystem) bool {
				return fs.Pwd() == "/home/learner/projects"
			},
			Explanation: "cd foldername moves you into that folder. You can always use pwd to check where you ended up.",
			Commands:    []string{"cd projects"},
		},

		// Mission 0.5: Go back up
		{
			ID:       "0.5-go-up",
			SkillID:  "cd",
			Level:    0,
			Title:    "Back Up One Level",
			Briefing: "You're deep in /home/learner/projects/src/utils. Go up one level to the src folder.",
			Hint:     "Two dots (..) means 'parent directory'",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/projects/src/utils")
				_ = fs.Cd("/home/learner/projects/src/utils")
			},
			Goal: func(fs *Filesystem) bool {
				return fs.Pwd() == "/home/learner/projects/src"
			},
			Explanation: "cd .. takes you up one level. .. always means 'the parent directory'.",
			Commands:    []string{"cd .."},
		},

		// Mission 0.6: Navigate with confidence
		{
			ID:       "0.6-navigate-path",
			SkillID:  "cd",
			Level:    0,
			Title:    "Navigate a Path",
			Briefing: "From /home/learner, get to /home/learner/documents/work in one command.",
			Hint:     "You can type a path: cd folder/subfolder",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/documents/work")
				_ = fs.Cd("/home/learner")
			},
			Goal: func(fs *Filesystem) bool {
				return fs.Pwd() == "/home/learner/documents/work"
			},
			Explanation: "You can navigate multiple levels at once: cd folder/subfolder. This is faster than cd folder then cd subfolder.",
			Commands:    []string{"cd documents/work"},
		},
	}
}

// ===================
// LEVEL 1: READING THE FILESYSTEM
// Goal: "I can find things"
// ===================

//nolint:funlen // Data definition function
func Level1Missions() []*Mission {
	return []*Mission{
		// Mission 1.1: Find the hidden file
		{
			ID:       "1.1-hidden-file",
			SkillID:  "ls-hidden",
			Level:    1,
			Title:    "Find the Hidden File",
			Briefing: "There's a hidden configuration file in this directory. Normal ls won't show it. Find it.",
			Hint:     "Hidden files start with a dot. ls has a flag to show ALL files...",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/project")
				_ = fs.Touch("/home/learner/project/readme.md")
				_ = fs.Touch("/home/learner/project/.env")
				_ = fs.WriteFile("/home/learner/project/.env", "SECRET_KEY=abc123\n")
				_ = fs.Cd("/home/learner/project")
			},
			Goal: func(fs *Filesystem) bool {
				return true // Completes when they see .env
			},
			Explanation: "ls -a shows ALL files, including hidden ones (files starting with .). The .env file contains secrets!",
			Commands:    []string{"ls -a"},
		},

		// Mission 1.2: Understand paths
		{
			ID:       "1.2-absolute-path",
			SkillID:  "paths",
			Level:    1,
			Title:    "Jump to Absolute Path",
			Briefing: "No matter where you are, get to /var/log using an absolute path.",
			Hint:     "Absolute paths start with /",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/var/log")
				_ = fs.Touch("/var/log/system.log")
				_ = fs.Cd("/home/learner/documents")
			},
			Goal: func(fs *Filesystem) bool {
				return fs.Pwd() == "/var/log"
			},
			Explanation: "Absolute paths start with / and work from anywhere. Relative paths depend on where you are.",
			Commands:    []string{"cd /var/log"},
		},

		// Mission 1.3: Read a file
		{
			ID:       "1.3-read-file",
			SkillID:  "cat",
			Level:    1,
			Title:    "Read the Secret",
			Briefing: "There's a file called secret.txt in this directory. What's inside?",
			Hint:     "cat displays file contents. Think of a cat knocking things off tables - it dumps everything out.",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/mission")
				_ = fs.WriteFile("/home/learner/mission/secret.txt", "The password is: turtlepower\n")
				_ = fs.Cd("/home/learner/mission")
			},
			Goal: func(fs *Filesystem) bool {
				return true
			},
			Explanation: "cat filename shows the contents of a file. It's called cat because it conCATenates files.",
			Commands:    []string{"cat secret.txt"},
		},

		// Mission 1.4: What's the current directory?
		{
			ID:       "1.4-dot-directory",
			SkillID:  "paths",
			Level:    1,
			Title:    "The Mysterious Dot",
			Briefing: "Run: ls . — What does the single dot mean?",
			Hint:     "A single dot represents something special...",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/test")
				_ = fs.Touch("/home/learner/test/file1.txt")
				_ = fs.Touch("/home/learner/test/file2.txt")
				_ = fs.Cd("/home/learner/test")
			},
			Goal: func(fs *Filesystem) bool {
				return true
			},
			Explanation: "A single dot (.) means 'the current directory'. ls . lists the current folder. This becomes useful later.",
			Commands:    []string{"ls ."},
		},

		// Mission 1.5: List with details
		{
			ID:       "1.5-long-listing",
			SkillID:  "ls-details",
			Level:    1,
			Title:    "Get the Details",
			Briefing: "List the files, but this time show the full details - sizes, dates, permissions.",
			Hint:     "There's a 'long' format flag for ls...",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/project")
				_ = fs.WriteFile("/home/learner/project/big.txt", strings.Repeat("data", 1000))
				_ = fs.Touch("/home/learner/project/small.txt")
				_ = fs.Cd("/home/learner/project")
			},
			Goal: func(fs *Filesystem) bool {
				return true
			},
			Explanation: "ls -l gives the 'long' listing with permissions, owner, size, and modification date.",
			Commands:    []string{"ls -l"},
		},
	}
}

// ===================
// LEVEL 2: FILE OPERATIONS
// Goal: "I can create and organize"
// ===================

//nolint:funlen // Data definition function
func Level2Missions() []*Mission {
	return []*Mission{
		// Mission 2.1: Create a directory
		{
			ID:       "2.1-create-dir",
			SkillID:  "mkdir",
			Level:    2,
			Title:    "Build a Workspace",
			Briefing: "Create a new directory called 'workspace' in your home folder.",
			Hint:     "mkdir = make directory",
			Setup: func(fs *Filesystem) {
				_ = fs.Cd("/home/learner")
			},
			Goal: func(fs *Filesystem) bool {
				return fs.Exists("/home/learner/workspace") && fs.IsDir("/home/learner/workspace")
			},
			Explanation: "mkdir creates a new directory. Now you have a place to organize your work!",
			Commands:    []string{"mkdir workspace"},
		},

		// Mission 2.2: Create a file
		{
			ID:       "2.2-create-file",
			SkillID:  "touch",
			Level:    2,
			Title:    "Create a Note",
			Briefing: "Create an empty file called 'notes.txt' in the current directory.",
			Hint:     "touch creates empty files (or updates timestamps of existing ones)",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/work")
				_ = fs.Cd("/home/learner/work")
			},
			Goal: func(fs *Filesystem) bool {
				return fs.Exists("/home/learner/work/notes.txt")
			},
			Explanation: "touch creates an empty file. It's called touch because it 'touches' the file, updating its timestamp.",
			Commands:    []string{"touch notes.txt"},
		},

		// Mission 2.3: Copy a file
		{
			ID:       "2.3-copy-file",
			SkillID:  "cp",
			Level:    2,
			Title:    "Backup the Config",
			Briefing: "There's a config.json here. Make a backup copy called config.backup.json",
			Hint:     "cp source destination",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/app")
				_ = fs.WriteFile("/home/learner/app/config.json", "{\n  \"debug\": false\n}\n")
				_ = fs.Cd("/home/learner/app")
			},
			Goal: func(fs *Filesystem) bool {
				return fs.Exists("/home/learner/app/config.backup.json")
			},
			Explanation: "cp copies files. Always make backups before editing important config files!",
			Commands:    []string{"cp config.json config.backup.json"},
		},

		// Mission 2.4: Move/rename a file
		{
			ID:       "2.4-move-file",
			SkillID:  "mv",
			Level:    2,
			Title:    "Organize the Download",
			Briefing: "There's a report.pdf in downloads/. Move it to documents/.",
			Hint:     "mv moves files. It's also how you rename things.",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/downloads")
				_ = fs.Mkdir("/home/learner/documents")
				_ = fs.Touch("/home/learner/downloads/report.pdf")
				_ = fs.Cd("/home/learner")
			},
			Goal: func(fs *Filesystem) bool {
				return fs.Exists("/home/learner/documents/report.pdf") &&
					!fs.Exists("/home/learner/downloads/report.pdf")
			},
			Explanation: "mv moves files. Unlike cp, the original is gone. mv is also how you rename files!",
			Commands:    []string{"mv downloads/report.pdf documents/"},
		},

		// Mission 2.5: Rename a file
		{
			ID:       "2.5-rename-file",
			SkillID:  "mv",
			Level:    2,
			Title:    "Fix the Typo",
			Briefing: "Someone created a file called 'importnat.txt'. Rename it to 'important.txt'.",
			Hint:     "mv also renames when source and destination are in the same directory",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/docs")
				_ = fs.Touch("/home/learner/docs/importnat.txt")
				_ = fs.Cd("/home/learner/docs")
			},
			Goal: func(fs *Filesystem) bool {
				return fs.Exists("/home/learner/docs/important.txt") &&
					!fs.Exists("/home/learner/docs/importnat.txt")
			},
			Explanation: "mv oldname newname renames a file. There's no separate rename command in Unix!",
			Commands:    []string{"mv importnat.txt important.txt"},
		},

		// Mission 2.6: Delete a file
		{
			ID:       "2.6-delete-file",
			SkillID:  "rm",
			Level:    2,
			Title:    "Clean Up",
			Briefing: "Delete the file called 'temp.txt'. Be careful - there's no undo!",
			Hint:     "rm = remove. It's permanent.",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/project")
				_ = fs.Touch("/home/learner/project/temp.txt")
				_ = fs.Touch("/home/learner/project/keep.txt")
				_ = fs.Cd("/home/learner/project")
			},
			Goal: func(fs *Filesystem) bool {
				return !fs.Exists("/home/learner/project/temp.txt") &&
					fs.Exists("/home/learner/project/keep.txt")
			},
			Explanation: "rm permanently deletes files. There's no trash can! Always double-check before using rm.",
			Commands:    []string{"rm temp.txt"},
		},

		// Mission 2.7: Write to a file
		{
			ID:       "2.7-echo-to-file",
			SkillID:  "redirect",
			Level:    2,
			Title:    "Leave a Note",
			Briefing: "Create a file called message.txt containing 'Hello World'",
			Hint:     "echo prints text. > redirects output to a file.",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/notes")
				_ = fs.Cd("/home/learner/notes")
			},
			Goal: func(fs *Filesystem) bool {
				content, err := fs.ReadFile("/home/learner/notes/message.txt")
				if err != nil {
					return false
				}
				return strings.Contains(content, "Hello World")
			},
			Explanation: "echo text > file creates a file with that content. > means 'send output to file'.",
			Commands:    []string{"echo Hello World > message.txt"},
		},

		// Mission 2.8: Full workflow
		{
			ID:       "2.8-organize-project",
			SkillID:  "file-ops",
			Level:    2,
			Title:    "Organize a Project",
			Briefing: "Create a 'src' folder, then create a file called 'main.py' inside it.",
			Hint:     "First mkdir, then touch (or you can do touch src/main.py if src exists)",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/myproject")
				_ = fs.Cd("/home/learner/myproject")
			},
			Goal: func(fs *Filesystem) bool {
				return fs.Exists("/home/learner/myproject/src/main.py")
			},
			Explanation: "Real projects need organization. Create folders with mkdir, files with touch. Build habits!",
			Commands:    []string{"mkdir src", "touch src/main.py"},
		},
	}
}

// ===================
// LEVEL 3: SEARCH + INSPECTION
// Goal: "I can find anything"
// ===================

//nolint:funlen // Data definition function
func Level3Missions() []*Mission {
	return []*Mission{
		// Mission 3.1: Find a string in a file
		{
			ID:       "3.1-grep-basic",
			SkillID:  "grep",
			Level:    3,
			Title:    "Find the Error",
			Briefing: "The log file has an error somewhere. Find the line containing 'ERROR'.",
			Hint:     "grep searches for patterns inside files",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/var/log")
				logContent := `2024-01-01 INFO: Server started
2024-01-01 INFO: Connection established
2024-01-01 ERROR: Database connection failed
2024-01-01 INFO: Retrying...
2024-01-01 INFO: Connection restored`
				_ = fs.WriteFile("/var/log/app.log", logContent)
				_ = fs.Cd("/var/log")
			},
			Goal: func(fs *Filesystem) bool {
				return true
			},
			Explanation: "grep pattern file shows all lines containing the pattern. Essential for debugging logs!",
			Commands:    []string{"grep ERROR app.log"},
		},

		// Mission 3.2: Find a file by name
		{
			ID:       "3.2-find-file",
			SkillID:  "find",
			Level:    3,
			Title:    "Hunt for the Config",
			Briefing: "There's a file called 'config.yaml' somewhere in /home/learner. Find it.",
			Hint:     "find searches for files by name. Use -name pattern",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/projects/secret/deeply/nested")
				_ = fs.WriteFile("/home/learner/projects/secret/deeply/nested/config.yaml", "key: value")
				_ = fs.Cd("/home/learner")
			},
			Goal: func(fs *Filesystem) bool {
				return true
			},
			Explanation: "find startpath -name pattern searches for files. Unlike grep, find looks at filenames, not contents.",
			Commands:    []string{"find . -name config.yaml", "find . -name \"config.yaml\""},
		},

		// Mission 3.3: Use a pipe
		{
			ID:       "3.3-pipe-intro",
			SkillID:  "pipes",
			Level:    3,
			Title:    "Chain Commands",
			Briefing: "List the files, but only show ones containing 'log' in the name. Use ls and grep together.",
			Hint:     "The | symbol sends output from one command to another",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/logs")
				_ = fs.Touch("/home/learner/logs/app.log")
				_ = fs.Touch("/home/learner/logs/error.log")
				_ = fs.Touch("/home/learner/logs/readme.txt")
				_ = fs.Touch("/home/learner/logs/debug.log")
				_ = fs.Cd("/home/learner/logs")
			},
			Goal: func(fs *Filesystem) bool {
				return true
			},
			Explanation: "The pipe | sends output from one command as input to another. ls | grep log shows only entries containing 'log'.",
			Commands:    []string{"ls | grep log"},
		},

		// Mission 3.4: Find what mentions a term
		{
			ID:       "3.4-grep-recursive",
			SkillID:  "grep",
			Level:    3,
			Title:    "Find All Debug Flags",
			Briefing: "Search ALL files in this directory for the word 'DEBUG'. Which files mention it?",
			Hint:     "grep -r searches recursively through directories",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/project")
				_ = fs.WriteFile("/home/learner/project/app.py", "DEBUG = True\nprint('hello')")
				_ = fs.WriteFile("/home/learner/project/config.py", "# No debug here")
				_ = fs.WriteFile("/home/learner/project/test.py", "DEBUG = False")
				_ = fs.Cd("/home/learner/project")
			},
			Goal: func(fs *Filesystem) bool {
				return true
			},
			Explanation: "grep -r searches all files in a directory tree. Great for finding where something is defined!",
			Commands:    []string{"grep -r DEBUG .", "grep DEBUG *"},
		},

		// Mission 3.5: Combine powers
		{
			ID:       "3.5-find-and-read",
			SkillID:  "find",
			Level:    3,
			Title:    "Detective Work",
			Briefing: "Find the hidden .secret file somewhere under /home/learner, then read its contents.",
			Hint:     "First find it with find, then read it with cat",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/stuff/more/things")
				_ = fs.WriteFile("/home/learner/stuff/more/.secret", "The treasure is buried under the old oak tree.\n")
				_ = fs.Cd("/home/learner")
			},
			Goal: func(fs *Filesystem) bool {
				return true // Complete when they read the file
			},
			Explanation: "Combining commands is powerful. find locates, cat reads. You can even pipe find's output to other commands!",
			Commands:    []string{"find . -name .secret", "cat stuff/more/.secret"},
		},
	}
}

// ===================
// LEVEL 4: TMUX BASICS
// Goal: "I can manage terminal sessions"
// ===================

//nolint:funlen // Data definition function
func Level4Missions() []*Mission {
	return []*Mission{
		// Mission 4.1: Start tmux
		{
			ID:       "4.1-start-tmux",
			SkillID:  "tmux-new",
			Level:    4,
			Title:    "Enter the Multiplexer",
			Briefing: "Start a new tmux session. Just type 'tmux' to begin.",
			Hint:     "Simply run: tmux",
			Setup:    func(_ *Filesystem) {},
			Goal: func(_ *Filesystem) bool {
				return true // Checked via TmuxState
			},
			Explanation: "tmux is a terminal multiplexer - it lets you have multiple terminal sessions in one window, and they persist even if you disconnect!",
			Commands:    []string{"tmux"},
		},

		// Mission 4.2: Create named session
		{
			ID:       "4.2-named-session",
			SkillID:  "tmux-new",
			Level:    4,
			Title:    "Name Your Session",
			Briefing: "Create a new tmux session named 'work'. Named sessions are easier to manage.",
			Hint:     "tmux new -s sessionname",
			Setup:    func(_ *Filesystem) {},
			Goal: func(_ *Filesystem) bool {
				return true
			},
			Explanation: "tmux new -s name creates a named session. When you have multiple sessions, names help you remember which is which.",
			Commands:    []string{"tmux new -s work", "tmux new-session -s work"},
		},

		// Mission 4.3: Detach from session
		{
			ID:       "4.3-detach",
			SkillID:  "tmux-detach",
			Level:    4,
			Title:    "Step Away",
			Briefing: "You're in a tmux session. Detach from it without closing it. The session will keep running!",
			Hint:     "tmux detach (or the keyboard shortcut Ctrl-b d)",
			Setup:    func(_ *Filesystem) {},
			Goal: func(_ *Filesystem) bool {
				return true
			},
			Explanation: "tmux detach leaves the session running in the background. Your processes continue even when you're not attached!",
			Commands:    []string{"tmux detach", "tmux d"},
		},

		// Mission 4.4: List sessions
		{
			ID:       "4.4-list-sessions",
			SkillID:  "tmux-list",
			Level:    4,
			Title:    "What's Running?",
			Briefing: "You've detached. Now check what tmux sessions are running.",
			Hint:     "tmux ls (short for list-sessions)",
			Setup:    func(_ *Filesystem) {},
			Goal: func(_ *Filesystem) bool {
				return true
			},
			Explanation: "tmux ls shows all your tmux sessions. You can see which ones are attached and how many windows each has.",
			Commands:    []string{"tmux ls", "tmux list-sessions"},
		},

		// Mission 4.5: Reattach to session
		{
			ID:       "4.5-attach",
			SkillID:  "tmux-attach",
			Level:    4,
			Title:    "Return to Work",
			Briefing: "Reattach to your detached tmux session.",
			Hint:     "tmux attach (or tmux a for short)",
			Setup:    func(_ *Filesystem) {},
			Goal: func(_ *Filesystem) bool {
				return true
			},
			Explanation: "tmux attach reconnects you to a session. Add -t sessionname to attach to a specific one.",
			Commands:    []string{"tmux attach", "tmux a", "tmux attach-session"},
		},

		// Mission 4.6: Split horizontally
		{
			ID:       "4.6-split-horizontal",
			SkillID:  "tmux-split-h",
			Level:    4,
			Title:    "Split the Screen",
			Briefing: "Split your tmux pane horizontally (left and right).",
			Hint:     "tmux split-window (or Ctrl-b %)",
			Setup:    func(_ *Filesystem) {},
			Goal: func(_ *Filesystem) bool {
				return true
			},
			Explanation: "tmux split-window creates a new pane. By default it splits horizontally. Now you can see two terminals at once!",
			Commands:    []string{"tmux split-window", "tmux split"},
		},

		// Mission 4.7: Split vertically
		{
			ID:       "4.7-split-vertical",
			SkillID:  "tmux-split-v",
			Level:    4,
			Title:    "Stack the Panes",
			Briefing: "Split your tmux pane vertically (top and bottom).",
			Hint:     "Add -v flag for vertical split (or Ctrl-b \")",
			Setup:    func(_ *Filesystem) {},
			Goal: func(_ *Filesystem) bool {
				return true
			},
			Explanation: "tmux split-window -v splits vertically. -v means vertical division (panes stacked top/bottom).",
			Commands:    []string{"tmux split-window -v", "tmux split -v"},
		},

		// Mission 4.8: Navigate panes
		{
			ID:       "4.8-select-pane",
			SkillID:  "tmux-pane-nav",
			Level:    4,
			Title:    "Jump Between Panes",
			Briefing: "You have multiple panes. Move to the pane on your right.",
			Hint:     "tmux select-pane -R (or Ctrl-b arrow key)",
			Setup:    func(_ *Filesystem) {},
			Goal: func(_ *Filesystem) bool {
				return true
			},
			Explanation: "tmux select-pane -L/R/U/D moves between panes. L=left, R=right, U=up, D=down.",
			Commands:    []string{"tmux select-pane -R"},
		},

		// Mission 4.9: Create new window
		{
			ID:       "4.9-new-window",
			SkillID:  "tmux-window-new",
			Level:    4,
			Title:    "Open a New Window",
			Briefing: "Create a new tmux window. Windows are like tabs in a browser.",
			Hint:     "tmux new-window (or Ctrl-b c)",
			Setup:    func(_ *Filesystem) {},
			Goal: func(_ *Filesystem) bool {
				return true
			},
			Explanation: "Windows are full-screen views. Use them to organize different tasks. Panes split one window; windows give you fresh space.",
			Commands:    []string{"tmux new-window"},
		},

		// Mission 4.10: Switch windows
		{
			ID:       "4.10-select-window",
			SkillID:  "tmux-window-nav",
			Level:    4,
			Title:    "Switch Windows",
			Briefing: "Switch to the next tmux window.",
			Hint:     "tmux select-window -n (or Ctrl-b n)",
			Setup:    func(_ *Filesystem) {},
			Goal: func(_ *Filesystem) bool {
				return true
			},
			Explanation: "select-window -n goes to next window, -p to previous. Or use Ctrl-b followed by the window number.",
			Commands:    []string{"tmux select-window -n"},
		},

		// Mission 4.11: Kill session
		{
			ID:       "4.11-kill-session",
			SkillID:  "tmux-kill",
			Level:    4,
			Title:    "Clean Up",
			Briefing: "You're done with this tmux session. Kill it entirely.",
			Hint:     "tmux kill-session",
			Setup:    func(_ *Filesystem) {},
			Goal: func(_ *Filesystem) bool {
				return true
			},
			Explanation: "tmux kill-session destroys the session and all its windows/panes. Use -t name to kill a specific session.",
			Commands:    []string{"tmux kill-session"},
		},
	}
}

// ===================
// LEVEL 5: MUSCLE MEMORY
// Goal: "I can work efficiently"
// ===================

//nolint:funlen // Data definition function
func Level5Missions() []*Mission {
	return []*Mission{
		// Mission 5.1: Project setup workflow
		{
			ID:       "5.1-project-setup",
			SkillID:  "workflow",
			Level:    5,
			Title:    "Set Up a Project",
			Briefing: "Create a project structure: mkdir myapp, then inside it create src/, tests/, and docs/ folders, plus a README.md file.",
			Hint:     "Use mkdir for folders, touch for files. You can chain commands.",
			Setup: func(fs *Filesystem) {
				_ = fs.Cd("/home/learner")
			},
			Goal: func(fs *Filesystem) bool {
				return fs.Exists("/home/learner/myapp/src") &&
					fs.Exists("/home/learner/myapp/tests") &&
					fs.Exists("/home/learner/myapp/docs") &&
					fs.Exists("/home/learner/myapp/README.md")
			},
			Explanation: "Good project structure is a habit. src/ for code, tests/ for tests, docs/ for documentation. README.md explains your project.",
			Commands:    []string{"mkdir myapp", "mkdir myapp/src myapp/tests myapp/docs", "touch myapp/README.md"},
		},

		// Mission 5.2: Find and edit
		{
			ID:       "5.2-find-and-edit",
			SkillID:  "workflow",
			Level:    5,
			Title:    "Hunt and Fix",
			Briefing: "Find all .py files under /home/learner/project, then create a backup of main.py as main.py.bak",
			Hint:     "find for searching, cp for backup",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/project/src")
				_ = fs.Mkdir("/home/learner/project/tests")
				_ = fs.WriteFile("/home/learner/project/src/main.py", "print('hello')")
				_ = fs.WriteFile("/home/learner/project/src/utils.py", "def helper(): pass")
				_ = fs.WriteFile("/home/learner/project/tests/test_main.py", "def test(): pass")
				_ = fs.Cd("/home/learner/project")
			},
			Goal: func(fs *Filesystem) bool {
				return fs.Exists("/home/learner/project/src/main.py.bak")
			},
			Explanation: "Real workflow: find files, understand structure, make backups before editing. Professionals always backup first!",
			Commands:    []string{"find . -name \"*.py\"", "cp src/main.py src/main.py.bak"},
		},

		// Mission 5.3: Tmux dev environment
		{
			ID:       "5.3-tmux-dev-setup",
			SkillID:  "tmux-workflow",
			Level:    5,
			Title:    "Dev Environment",
			Briefing: "Create a tmux session named 'dev', then split it into two panes.",
			Hint:     "tmux new -s dev, then tmux split-window",
			Setup:    func(_ *Filesystem) {},
			Goal: func(_ *Filesystem) bool {
				return true
			},
			Explanation: "Pro setup: one pane for editing, one for running commands. Name your sessions so you can find them later!",
			Commands:    []string{"tmux new -s dev", "tmux split-window"},
		},

		// Mission 5.4: Log investigation
		{
			ID:       "5.4-log-investigation",
			SkillID:  "workflow",
			Level:    5,
			Title:    "Debug Detective",
			Briefing: "Find all ERROR lines in /var/log/app.log, and also check what WARNING messages exist.",
			Hint:     "Use grep twice with different patterns",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/var/log")
				logContent := `2024-01-01 10:00:00 INFO: Server started
2024-01-01 10:00:01 INFO: Loading configuration
2024-01-01 10:00:02 WARNING: Config file not found, using defaults
2024-01-01 10:00:03 INFO: Database connecting
2024-01-01 10:00:05 ERROR: Database connection timeout
2024-01-01 10:00:06 INFO: Retrying connection
2024-01-01 10:00:08 WARNING: Slow connection detected
2024-01-01 10:00:10 INFO: Connected successfully
2024-01-01 10:05:00 ERROR: Request failed: 500
2024-01-01 10:05:01 INFO: Error logged to monitoring`
				_ = fs.WriteFile("/var/log/app.log", logContent)
				_ = fs.Cd("/var/log")
			},
			Goal: func(_ *Filesystem) bool {
				return true
			},
			Explanation: "Real debugging: check errors first, then warnings. grep is your best friend for log analysis!",
			Commands:    []string{"grep ERROR app.log", "grep WARNING app.log"},
		},

		// Mission 5.5: Cleanup workflow
		{
			ID:       "5.5-cleanup",
			SkillID:  "workflow",
			Level:    5,
			Title:    "Spring Cleaning",
			Briefing: "In /home/learner/downloads, delete all .tmp files but keep everything else.",
			Hint:     "Use find to locate .tmp files, then rm each one",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/downloads")
				_ = fs.Touch("/home/learner/downloads/report.pdf")
				_ = fs.Touch("/home/learner/downloads/cache.tmp")
				_ = fs.Touch("/home/learner/downloads/session.tmp")
				_ = fs.Touch("/home/learner/downloads/photo.jpg")
				_ = fs.Touch("/home/learner/downloads/backup.tmp")
				_ = fs.Cd("/home/learner/downloads")
			},
			Goal: func(fs *Filesystem) bool {
				return !fs.Exists("/home/learner/downloads/cache.tmp") &&
					!fs.Exists("/home/learner/downloads/session.tmp") &&
					!fs.Exists("/home/learner/downloads/backup.tmp") &&
					fs.Exists("/home/learner/downloads/report.pdf") &&
					fs.Exists("/home/learner/downloads/photo.jpg")
			},
			Explanation: "Targeted cleanup: find what to delete, verify, then remove. Never use rm * blindly!",
			Commands:    []string{"find . -name \"*.tmp\"", "rm cache.tmp session.tmp backup.tmp"},
		},

		// Mission 5.6: Multi-window workflow
		{
			ID:       "5.6-multi-window",
			SkillID:  "tmux-workflow",
			Level:    5,
			Title:    "Window Manager",
			Briefing: "Create a tmux session with two windows: one for 'code', one for 'tests'.",
			Hint:     "tmux new-session, then tmux new-window twice",
			Setup:    func(_ *Filesystem) {},
			Goal: func(_ *Filesystem) bool {
				return true
			},
			Explanation: "Multiple windows let you organize by task. Window 1 for coding, window 2 for tests, window 3 for servers...",
			Commands:    []string{"tmux new -s project", "tmux new-window", "tmux new-window"},
		},

		// Mission 5.7: Organize a mess
		{
			ID:       "5.7-organize-mess",
			SkillID:  "workflow",
			Level:    5,
			Title:    "Untangle the Mess",
			Briefing: "Move all .log files to logs/, all .py files to src/, and all .md files to docs/. Create folders if needed.",
			Hint:     "mkdir first, then mv files to appropriate folders",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/chaos")
				_ = fs.Touch("/home/learner/chaos/app.log")
				_ = fs.Touch("/home/learner/chaos/error.log")
				_ = fs.Touch("/home/learner/chaos/main.py")
				_ = fs.Touch("/home/learner/chaos/utils.py")
				_ = fs.Touch("/home/learner/chaos/README.md")
				_ = fs.Touch("/home/learner/chaos/CHANGELOG.md")
				_ = fs.Cd("/home/learner/chaos")
			},
			Goal: func(fs *Filesystem) bool {
				return fs.Exists("/home/learner/chaos/logs/app.log") &&
					fs.Exists("/home/learner/chaos/logs/error.log") &&
					fs.Exists("/home/learner/chaos/src/main.py") &&
					fs.Exists("/home/learner/chaos/src/utils.py") &&
					fs.Exists("/home/learner/chaos/docs/README.md") &&
					fs.Exists("/home/learner/chaos/docs/CHANGELOG.md")
			},
			Explanation: "Organization is a skill. Group related files, use clear folder names. Your future self will thank you!",
			Commands:    []string{"mkdir logs src docs", "mv app.log error.log logs/", "mv main.py utils.py src/", "mv README.md CHANGELOG.md docs/"},
		},

		// Mission 5.8: Quick environment check
		{
			ID:       "5.8-env-check",
			SkillID:  "workflow",
			Level:    5,
			Title:    "Environment Recon",
			Briefing: "Figure out: Where are you? What's here? Are there any hidden files? Do this in three commands.",
			Hint:     "pwd, ls, ls -a",
			Setup: func(fs *Filesystem) {
				_ = fs.Mkdir("/home/learner/secret-project")
				_ = fs.Touch("/home/learner/secret-project/visible.txt")
				_ = fs.Touch("/home/learner/secret-project/.hidden-config")
				_ = fs.Touch("/home/learner/secret-project/.env")
				_ = fs.Cd("/home/learner/secret-project")
			},
			Goal: func(_ *Filesystem) bool {
				return true
			},
			Explanation: "First thing in any directory: pwd, ls, ls -a. Know where you are and what's there. Hidden files often contain secrets!",
			Commands:    []string{"pwd", "ls", "ls -a"},
		},
	}
}
