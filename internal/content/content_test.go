// ABOUTME: Integration tests for content loading
// ABOUTME: Verifies YAML content loads correctly at startup

package content

import (
	"testing"
)

func TestLoadContent(t *testing.T) {
	err := LoadContent()
	if err != nil {
		t.Fatalf("LoadContent failed: %v", err)
	}
}

func TestGetSkillGraph(t *testing.T) {
	graph, err := GetSkillGraph()
	if err != nil {
		t.Fatalf("GetSkillGraph failed: %v", err)
	}

	// Verify we have skills loaded by checking a known category
	navSkills := graph.GetSkillsByCategory("navigation")
	if len(navSkills) == 0 {
		t.Error("expected navigation skills to be loaded, got 0")
	}

	// Verify a known skill exists by checking prerequisites
	prereqs := graph.GetPrerequisites("ls")
	if len(prereqs) == 0 {
		t.Error("expected ls to have prerequisites")
	}
	if prereqs[0] != "pwd" {
		t.Errorf("expected ls prerequisite to be pwd, got %s", prereqs[0])
	}
}

func TestGetRawChallenges(t *testing.T) {
	challenges, err := GetRawChallenges()
	if err != nil {
		t.Fatalf("GetRawChallenges failed: %v", err)
	}

	if len(challenges) == 0 {
		t.Error("expected challenges to be loaded, got 0")
	}

	// Verify a known skill has challenges
	pwdChallenges, ok := challenges["pwd"]
	if !ok {
		t.Error("expected challenges for pwd skill")
	}
	if len(pwdChallenges) == 0 {
		t.Error("expected at least one challenge for pwd skill")
	}
}

func TestGetRawMissions(t *testing.T) {
	missions, err := GetRawMissions()
	if err != nil {
		t.Fatalf("GetRawMissions failed: %v", err)
	}

	if len(missions) == 0 {
		t.Error("expected missions to be loaded, got 0")
	}

	// Verify missions have required fields
	for _, m := range missions {
		if m.ID == "" {
			t.Error("mission ID should not be empty")
		}
		if m.SkillID == "" {
			t.Error("mission SkillID should not be empty")
		}
		if m.Title == "" {
			t.Error("mission Title should not be empty")
		}
	}
}
