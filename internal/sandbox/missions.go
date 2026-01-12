// ABOUTME: Mission definitions organized by pedagogical level
// ABOUTME: Each mission has a clear goal, constrained environment, and immediate feedback

package sandbox

import "strings"

// GetAllMissions returns missions organized by level.
func GetAllMissions() map[int][]*Mission {
	return map[int][]*Mission{
		0: Level0Missions(), // Orientation
		1: Level1Missions(), // Reading the filesystem
		2: Level2Missions(), // File operations
		3: Level3Missions(), // Search + inspection
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
