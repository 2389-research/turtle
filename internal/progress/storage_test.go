// ABOUTME: Tests for progress persistence functionality
// ABOUTME: Validates save/load of user progress to disk

package progress

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/2389-research/turtle/internal/skills"
)

func TestSaveAndLoad(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "turtle-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	savePath := filepath.Join(tmpDir, "progress.json")

	// Create progress with some data
	original := skills.NewUserProgress()
	original.AddXP(150)
	original.Practice("pwd", 5)
	original.Practice("ls", 4)
	original.RecordActivity()

	// Save it
	err = Save(original, savePath)
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	// Load it back
	loaded, err := Load(savePath)
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	// Verify data
	if loaded.XP != original.XP {
		t.Errorf("XP mismatch: got %d, want %d", loaded.XP, original.XP)
	}
	if loaded.Level != original.Level {
		t.Errorf("Level mismatch: got %d, want %d", loaded.Level, original.Level)
	}
	if loaded.CurrentStreak != original.CurrentStreak {
		t.Errorf("Streak mismatch: got %d, want %d", loaded.CurrentStreak, original.CurrentStreak)
	}

	// Check card data
	pwdCard := loaded.GetCard("pwd")
	if pwdCard == nil {
		t.Fatal("pwd card not found after load")
	}
	if pwdCard.Repetitions == 0 {
		t.Error("pwd card should have repetitions")
	}
}

func TestLoadNonExistent(t *testing.T) {
	// Loading non-existent file should return fresh progress
	progress, err := Load("/nonexistent/path/progress.json")
	if err != nil {
		t.Fatalf("load should not error on missing file: %v", err)
	}
	if progress.Level != 1 {
		t.Errorf("expected fresh progress with level 1, got %d", progress.Level)
	}
	if progress.XP != 0 {
		t.Errorf("expected fresh progress with 0 XP, got %d", progress.XP)
	}
}

func TestGetDefaultPath(t *testing.T) {
	path := GetDefaultPath()

	if path == "" {
		t.Error("default path should not be empty")
	}

	// Should contain turtle
	if filepath.Base(filepath.Dir(path)) != "turtle" && filepath.Base(path) != "progress.json" {
		t.Errorf("unexpected path structure: %s", path)
	}
}

func TestEnsureDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "turtle-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	nestedPath := filepath.Join(tmpDir, "a", "b", "c", "progress.json")

	// Directory doesn't exist yet
	dir := filepath.Dir(nestedPath)
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatal("directory should not exist yet")
	}

	// Create progress and save
	progress := skills.NewUserProgress()
	err = Save(progress, nestedPath)
	if err != nil {
		t.Fatalf("save should create directories: %v", err)
	}

	// Now directory should exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("directory should exist after save")
	}
}
