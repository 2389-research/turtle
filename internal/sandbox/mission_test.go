// ABOUTME: Tests for the mission system
// ABOUTME: Verifies goal-based learning challenges work correctly

package sandbox

import (
	"strings"
	"testing"
)

func TestMissionRunner_Execute(t *testing.T) {
	mission := &Mission{
		ID: "test",
		Setup: func(fs *Filesystem) {
			_ = fs.Mkdir("/tmp")
			_ = fs.Cd("/tmp")
		},
		Goal: func(fs *Filesystem) bool { return fs.Pwd() == "/" },
	}

	runner := NewMissionRunner(mission)

	// Initial state - should be in /tmp after setup
	if runner.FS.Pwd() != "/tmp" {
		t.Errorf("Expected /tmp, got %s", runner.FS.Pwd())
	}

	// Execute command
	result := runner.Execute("cd /")
	if !result.Success {
		t.Errorf("Command failed: %s", result.Error)
	}
	if !result.Completed {
		t.Error("Mission should be completed")
	}
}

func TestMissionRunner_Reset(t *testing.T) {
	mission := &Mission{
		ID: "test",
		Setup: func(fs *Filesystem) {
			_ = fs.Mkdir("/test")
			_ = fs.Cd("/test")
		},
	}

	runner := NewMissionRunner(mission)

	// Make changes
	runner.Execute("touch new.txt")
	if !runner.FS.Exists("/test/new.txt") {
		t.Error("File should exist after touch")
	}

	// Reset
	runner.Reset()
	if runner.FS.Exists("/test/new.txt") {
		t.Error("File should not exist after reset")
	}
}

func TestMissionRunner_Pwd(t *testing.T) {
	mission := &Mission{
		Setup: func(fs *Filesystem) {
			// Default filesystem already starts in /home/learner
		},
	}
	runner := NewMissionRunner(mission)

	result := runner.Execute("pwd")
	if !result.Success {
		t.Errorf("pwd failed: %s", result.Error)
	}
	// Default filesystem starts at /home/learner
	if result.Output != "/home/learner" {
		t.Errorf("Expected /home/learner, got %s", result.Output)
	}
}

func TestMissionRunner_LsHidden(t *testing.T) {
	mission := &Mission{
		Setup: func(fs *Filesystem) {
			_ = fs.Mkdir("/test")
			_ = fs.Touch("/test/.hidden")
			_ = fs.Touch("/test/visible")
			_ = fs.Cd("/test")
		},
	}
	runner := NewMissionRunner(mission)

	// Without -a
	result := runner.Execute("ls")
	if strings.Contains(result.Output, ".hidden") {
		t.Error("Hidden file should not appear without -a")
	}

	// With -a
	result = runner.Execute("ls -a")
	if !strings.Contains(result.Output, ".hidden") {
		t.Error("Hidden file should appear with -a")
	}
}

func TestMissionRunner_Echo(t *testing.T) {
	mission := &Mission{Setup: func(fs *Filesystem) { _ = fs.Cd("/tmp") }}
	runner := NewMissionRunner(mission)

	result := runner.Execute("echo hello world")
	if result.Output != "hello world" {
		t.Errorf("Expected 'hello world', got %q", result.Output)
	}
}

func TestMissionRunner_EchoRedirect(t *testing.T) {
	mission := &Mission{
		Setup: func(fs *Filesystem) {
			_ = fs.Mkdir("/tmp")
			_ = fs.Cd("/tmp")
		},
	}
	runner := NewMissionRunner(mission)

	result := runner.Execute("echo test > file.txt")
	if !result.Success {
		t.Errorf("Echo redirect failed: %s", result.Error)
	}

	// File should be created in current directory (/tmp)
	content, err := runner.FS.ReadFile("/tmp/file.txt")
	if err != nil {
		t.Errorf("ReadFile failed: %v", err)
	}
	if !strings.Contains(content, "test") {
		t.Errorf("Expected 'test' in file, got %q", content)
	}
}

func TestMissionRunner_Cat(t *testing.T) {
	mission := &Mission{
		Setup: func(fs *Filesystem) {
			_ = fs.WriteFile("/test.txt", "hello")
			_ = fs.Cd("/")
		},
	}
	runner := NewMissionRunner(mission)

	result := runner.Execute("cat test.txt")
	if result.Output != "hello" {
		t.Errorf("Expected 'hello', got %q", result.Output)
	}
}

func TestMissionRunner_Grep(t *testing.T) {
	mission := &Mission{
		Setup: func(fs *Filesystem) {
			_ = fs.WriteFile("/log.txt", "line1\nERROR here\nline3")
			_ = fs.Cd("/")
		},
	}
	runner := NewMissionRunner(mission)

	result := runner.Execute("grep ERROR log.txt")
	if !strings.Contains(result.Output, "ERROR") {
		t.Errorf("Expected ERROR in output, got %q", result.Output)
	}
}

func TestMissionRunner_CommandNotFound(t *testing.T) {
	mission := &Mission{Setup: func(fs *Filesystem) {}}
	runner := NewMissionRunner(mission)

	result := runner.Execute("notacommand")
	if result.Error == "" {
		t.Error("Unknown command should return error")
	}
	if !strings.Contains(result.Error, "not found") {
		t.Errorf("Expected 'not found' in error, got %q", result.Error)
	}
}

func TestLevel0Missions(t *testing.T) {
	missions := Level0Missions()

	if len(missions) == 0 {
		t.Error("Level 0 should have missions")
	}

	// Test each mission can be created and run
	for _, m := range missions {
		runner := NewMissionRunner(m)
		if runner.FS == nil {
			t.Errorf("Mission %s: filesystem is nil", m.ID)
		}
		if m.Briefing == "" {
			t.Errorf("Mission %s: missing briefing", m.ID)
		}
	}
}

func TestMission_WhereAmI(t *testing.T) {
	missions := Level0Missions()
	mission := missions[0] // "Where Am I?"

	runner := NewMissionRunner(mission)

	// Verify setup worked - should be in projects
	if runner.FS.Pwd() != "/home/learner/projects" {
		t.Errorf("Setup should put us in /home/learner/projects, got %s", runner.FS.Pwd())
	}

	result := runner.Execute("pwd")
	if !result.Success {
		t.Errorf("pwd should succeed: %s", result.Error)
	}
	if result.Output != "/home/learner/projects" {
		t.Errorf("Expected /home/learner/projects, got %s", result.Output)
	}
}

func TestMission_GoHome(t *testing.T) {
	// Find "Go Home" mission
	missions := Level0Missions()
	var mission *Mission
	for _, m := range missions {
		if m.ID == "0.3-go-home" {
			mission = m
			break
		}
	}

	if mission == nil {
		t.Fatal("Go Home mission not found")
	}

	runner := NewMissionRunner(mission)

	// Should start in /tmp
	if runner.FS.Pwd() != "/tmp" {
		t.Errorf("Expected /tmp, got %s", runner.FS.Pwd())
	}

	// Complete the mission
	result := runner.Execute("cd")
	if !result.Completed {
		t.Error("Mission should be completed after cd")
	}
}

func TestMission_CreateDir(t *testing.T) {
	missions := Level2Missions()
	var mission *Mission
	for _, m := range missions {
		if m.ID == "2.1-create-dir" {
			mission = m
			break
		}
	}

	if mission == nil {
		t.Fatal("Create dir mission not found")
	}

	runner := NewMissionRunner(mission)

	// Should not be complete initially
	if mission.Goal(runner.FS) {
		t.Error("Should not be complete before mkdir")
	}

	// Complete the mission
	runner.Execute("mkdir workspace")
	if !mission.Goal(runner.FS) {
		t.Error("Should be complete after mkdir workspace")
	}
}
