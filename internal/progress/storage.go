// ABOUTME: Persistence layer for user progress
// ABOUTME: Saves and loads progress data to/from JSON file

package progress

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/2389-research/turtle/internal/skills"
	"github.com/2389-research/turtle/internal/srs"
)

// SaveData is the serializable representation of user progress
type SaveData struct {
	XP            int                `json:"xp"`
	Level         int                `json:"level"`
	CurrentStreak int                `json:"current_streak"`
	BestStreak    int                `json:"best_streak"`
	LastActive    time.Time          `json:"last_active"`
	Cards         map[string]*CardData `json:"cards"`
}

// CardData is the serializable representation of an SRS card
type CardData struct {
	SkillID      string    `json:"skill_id"`
	EaseFactor   float64   `json:"ease_factor"`
	Interval     int       `json:"interval"`
	Repetitions  int       `json:"repetitions"`
	LastReviewed time.Time `json:"last_reviewed"`
}

// GetDefaultPath returns the default save file location
func GetDefaultPath() string {
	// Use XDG_DATA_HOME if set, otherwise ~/.local/share
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			// Fallback to current directory
			return "turtle_progress.json"
		}
		dataHome = filepath.Join(home, ".local", "share")
	}

	return filepath.Join(dataHome, "turtle", "progress.json")
}

// Save persists user progress to disk
func Save(progress *skills.UserProgress, path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Convert to saveable format
	data := SaveData{
		XP:            progress.XP,
		Level:         progress.Level,
		CurrentStreak: progress.CurrentStreak,
		BestStreak:    progress.BestStreak,
		LastActive:    progress.LastActive,
		Cards:         make(map[string]*CardData),
	}

	for skillID, card := range progress.Cards {
		data.Cards[skillID] = &CardData{
			SkillID:      card.SkillID,
			EaseFactor:   card.EaseFactor,
			Interval:     card.Interval,
			Repetitions:  card.Repetitions,
			LastReviewed: card.LastReviewed,
		}
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Write atomically using temp file
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, jsonData, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, path)
}

// Load reads user progress from disk
func Load(path string) (*skills.UserProgress, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return fresh progress
		return skills.NewUserProgress(), nil
	}

	// Read file
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var data SaveData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	// Convert back to UserProgress
	progress := skills.NewUserProgress()
	progress.XP = data.XP
	progress.Level = data.Level
	progress.CurrentStreak = data.CurrentStreak
	progress.BestStreak = data.BestStreak
	progress.LastActive = data.LastActive

	for skillID, cardData := range data.Cards {
		card := srs.NewCard(skillID)
		card.EaseFactor = cardData.EaseFactor
		card.Interval = cardData.Interval
		card.Repetitions = cardData.Repetitions
		card.LastReviewed = cardData.LastReviewed
		progress.Cards[skillID] = card
	}

	return progress, nil
}
