# Turtle 🐢

**Duolingo for the terminal.** Learn tmux and shell commands through gamified, spaced microlearning.

## What is this?

Turtle is a TUI app that teaches complete beginners how to use the terminal and tmux through:

- **Gamified learning** - XP, streaks, levels, combos
- **Microlearning** - Bite-sized 2-5 minute lessons
- **Spaced repetition** - SM-2 algorithm schedules reviews based on forgetting curves
- **Mastery-based progression** - Skills unlock as you prove competence
- **Immediate feedback** - Know right away if you got it

## Installation

### Homebrew (macOS/Linux)

```bash
brew install 2389-research/tap/turtle
```

### Go Install

```bash
go install github.com/2389-research/turtle/cmd/turtle@latest
```

### Build from source

```bash
git clone https://github.com/2389-research/turtle
cd turtle
go build -o turtle ./cmd/turtle
./turtle
```

## Usage

Just run `turtle`:

```bash
./turtle
```

Navigate with arrow keys or `j`/`k`, select with Enter, quit with `q`.

## Learning Path

```
Level 0-1: Navigation
├── pwd (where am I?)
├── ls (what's here?)
├── cd (move around)
└── tab completion

Level 2: File Operations
├── mkdir, touch, rm
├── cp, mv
└── cat

Level 3+: Tmux
├── Sessions (new, detach, attach)
├── Panes (split, navigate)
└── Windows (create, switch)
```

## Progress

Your progress is automatically saved to `~/.local/share/turtle/progress.json`. Streaks, XP, and skill mastery persist between sessions.

## Tech Stack

- Go
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- SM-2 algorithm for spaced repetition

## License

MIT
