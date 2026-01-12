// ABOUTME: Content loading with go:embed for skills, challenges, and missions
// ABOUTME: Provides adapter functions to convert YAML to existing Go types

package content

import (
	"embed"
	"fmt"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/2389-research/turtle/internal/skills"
)

//go:embed skills.yaml challenges.yaml missions.yaml
var contentFS embed.FS

var (
	loadOnce       sync.Once
	skillsData     *SkillsFile
	challengesData *ChallengesFile
	missionsData   *MissionsFile
	loadErr        error
)

// LoadContent loads all embedded YAML content files.
// Safe to call multiple times - only loads once.
func LoadContent() error {
	loadOnce.Do(func() {
		loadErr = loadAll()
	})
	return loadErr
}

func loadAll() error {
	// Load skills
	data, err := contentFS.ReadFile("skills.yaml")
	if err != nil {
		return fmt.Errorf("reading skills.yaml: %w", err)
	}
	skillsData = &SkillsFile{}
	if err := yaml.Unmarshal(data, skillsData); err != nil {
		return fmt.Errorf("parsing skills.yaml: %w", err)
	}

	// Load challenges
	data, err = contentFS.ReadFile("challenges.yaml")
	if err != nil {
		return fmt.Errorf("reading challenges.yaml: %w", err)
	}
	challengesData = &ChallengesFile{}
	if err := yaml.Unmarshal(data, challengesData); err != nil {
		return fmt.Errorf("parsing challenges.yaml: %w", err)
	}

	// Load missions
	data, err = contentFS.ReadFile("missions.yaml")
	if err != nil {
		return fmt.Errorf("reading missions.yaml: %w", err)
	}
	missionsData = &MissionsFile{}
	if err := yaml.Unmarshal(data, missionsData); err != nil {
		return fmt.Errorf("parsing missions.yaml: %w", err)
	}

	return validate()
}

func validate() error {
	// Build skill ID set for validation
	skillIDs := make(map[string]bool)
	for _, s := range skillsData.Skills {
		if skillIDs[s.ID] {
			return fmt.Errorf("duplicate skill ID: %s", s.ID)
		}
		skillIDs[s.ID] = true
	}

	// Validate skill prerequisites
	for _, s := range skillsData.Skills {
		for _, prereq := range s.Prerequisites {
			if !skillIDs[prereq] {
				return fmt.Errorf("skill %s: unknown prerequisite %s", s.ID, prereq)
			}
		}
	}

	// Validate challenge skill references
	for skillID := range challengesData.Challenges {
		if !skillIDs[skillID] {
			return fmt.Errorf("challenges reference unknown skill: %s", skillID)
		}
	}

	// Validate mission skill references, goals, and uniqueness
	missionIDs := make(map[string]bool)
	for _, m := range missionsData.Missions {
		if missionIDs[m.ID] {
			return fmt.Errorf("duplicate mission ID: %s", m.ID)
		}
		missionIDs[m.ID] = true

		if !skillIDs[m.SkillID] {
			return fmt.Errorf("mission %s: unknown skill_id %s", m.ID, m.SkillID)
		}
		if _, err := ParseGoal(m.Goal); err != nil {
			return fmt.Errorf("mission %s: invalid goal: %w", m.ID, err)
		}
	}

	return nil
}

// GetSkillGraph returns a populated SkillGraph from the YAML content.
func GetSkillGraph() (*skills.SkillGraph, error) {
	if err := LoadContent(); err != nil {
		return nil, err
	}

	graph := skills.NewSkillGraph()

	for _, s := range skillsData.Skills {
		graph.AddSkill(&skills.Skill{
			ID:                s.ID,
			Name:              s.Name,
			Description:       s.Description,
			Category:          skills.Category(s.Category),
			Prerequisites:     s.Prerequisites,
			RequiresCategory:  skills.Category(s.RequiresCategory),
			CategoryThreshold: s.CategoryThreshold,
		})
	}

	return graph, nil
}

// GetRawChallenges returns raw challenge data organized by skill ID.
// Consumers are responsible for converting to their own types.
func GetRawChallenges() (map[string][]YAMLChallenge, error) {
	if err := LoadContent(); err != nil {
		return nil, err
	}
	return challengesData.Challenges, nil
}

// GetRawMissions returns raw mission data.
// Consumers are responsible for converting to their own types.
func GetRawMissions() ([]YAMLMission, error) {
	if err := LoadContent(); err != nil {
		return nil, err
	}
	return missionsData.Missions, nil
}
