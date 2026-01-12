// ABOUTME: YAML struct definitions for content loading
// ABOUTME: Defines schemas for skills, challenges, and missions

package content

// SkillsFile represents the top-level skills.yaml structure.
type SkillsFile struct {
	Version int         `yaml:"version"`
	Skills  []YAMLSkill `yaml:"skills"`
}

// YAMLSkill represents a skill definition in YAML.
type YAMLSkill struct {
	ID                string   `yaml:"id"`
	Name              string   `yaml:"name"`
	Description       string   `yaml:"description"`
	Category          string   `yaml:"category"`
	Prerequisites     []string `yaml:"prerequisites,omitempty"`
	RequiresCategory  string   `yaml:"requires_category,omitempty"`
	CategoryThreshold float64  `yaml:"category_threshold,omitempty"`
}

// ChallengesFile represents the top-level challenges.yaml structure.
type ChallengesFile struct {
	Version    int                        `yaml:"version"`
	Challenges map[string][]YAMLChallenge `yaml:"challenges"`
}

// YAMLChallenge represents a flashcard challenge in YAML.
type YAMLChallenge struct {
	Type        string   `yaml:"type"`
	Prompt      string   `yaml:"prompt"`
	Expected    string   `yaml:"expected,omitempty"`
	Hint        string   `yaml:"hint,omitempty"`
	Explanation string   `yaml:"explanation,omitempty"`
	Options     []string `yaml:"options,omitempty"`
	Correct     int      `yaml:"correct,omitempty"`
	Broken      string   `yaml:"broken,omitempty"`
	Command     string   `yaml:"command,omitempty"`
}

// MissionsFile represents the top-level missions.yaml structure.
type MissionsFile struct {
	Version  int           `yaml:"version"`
	Missions []YAMLMission `yaml:"missions"`
}

// YAMLMission represents a mission definition in YAML.
type YAMLMission struct {
	ID          string         `yaml:"id"`
	SkillID     string         `yaml:"skill_id"`
	Level       int            `yaml:"level"`
	Title       string         `yaml:"title"`
	Briefing    string         `yaml:"briefing"`
	Hint        string         `yaml:"hint,omitempty"`
	Explanation string         `yaml:"explanation,omitempty"`
	Commands    []string       `yaml:"commands,omitempty"`
	Setup       []SetupAction  `yaml:"setup,omitempty"`
	Goal        map[string]any `yaml:"goal"`
}

// SetupAction represents a single setup operation.
// Only one field should be set per action.
type SetupAction struct {
	Mkdir     string           `yaml:"mkdir,omitempty"`
	Cd        string           `yaml:"cd,omitempty"`
	Touch     string           `yaml:"touch,omitempty"`
	WriteFile *WriteFileAction `yaml:"write_file,omitempty"`
}

// WriteFileAction represents a write_file setup operation.
type WriteFileAction struct {
	Path    string `yaml:"path"`
	Content string `yaml:"content"`
}
