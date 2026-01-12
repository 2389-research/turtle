// ABOUTME: Goal DSL parser and evaluator for mission completion
// ABOUTME: Supports AND/OR/NOT logic combinators and filesystem predicates

package content

import (
	"fmt"
	"strings"
)

// GoalEvaluator provides filesystem operations needed for goal evaluation.
// This interface is implemented by sandbox.Filesystem.
type GoalEvaluator interface {
	Pwd() string
	Exists(path string) bool
	IsDir(path string) bool
	ReadFile(path string) (string, error)
}

// GoalNode represents a parsed goal condition that can be evaluated.
type GoalNode interface {
	Evaluate(fs GoalEvaluator) bool
}

// AlwaysGoal always returns true (for command-match missions).
type AlwaysGoal struct{}

func (g *AlwaysGoal) Evaluate(_ GoalEvaluator) bool {
	return true
}

// PwdEqualsGoal checks if current directory matches the expected path.
type PwdEqualsGoal struct {
	Path string
}

func (g *PwdEqualsGoal) Evaluate(fs GoalEvaluator) bool {
	return fs.Pwd() == g.Path
}

// PathExistsGoal checks if a path exists.
type PathExistsGoal struct {
	Path string
}

func (g *PathExistsGoal) Evaluate(fs GoalEvaluator) bool {
	return fs.Exists(g.Path)
}

// PathNotExistsGoal checks if a path does not exist.
type PathNotExistsGoal struct {
	Path string
}

func (g *PathNotExistsGoal) Evaluate(fs GoalEvaluator) bool {
	return !fs.Exists(g.Path)
}

// IsDirGoal checks if a path is a directory.
type IsDirGoal struct {
	Path string
}

func (g *IsDirGoal) Evaluate(fs GoalEvaluator) bool {
	return fs.IsDir(g.Path)
}

// IsFileGoal checks if a path is a regular file (exists and not a directory).
type IsFileGoal struct {
	Path string
}

func (g *IsFileGoal) Evaluate(fs GoalEvaluator) bool {
	return fs.Exists(g.Path) && !fs.IsDir(g.Path)
}

// FileContainsGoal checks if a file contains a substring.
type FileContainsGoal struct {
	Path    string
	Content string
}

func (g *FileContainsGoal) Evaluate(fs GoalEvaluator) bool {
	content, err := fs.ReadFile(g.Path)
	if err != nil {
		return false
	}
	return strings.Contains(content, g.Content)
}

// AndGoal requires all child conditions to be true.
type AndGoal struct {
	Conditions []GoalNode
}

func (g *AndGoal) Evaluate(fs GoalEvaluator) bool {
	for _, c := range g.Conditions {
		if !c.Evaluate(fs) {
			return false
		}
	}
	return true
}

// OrGoal requires at least one child condition to be true.
type OrGoal struct {
	Conditions []GoalNode
}

func (g *OrGoal) Evaluate(fs GoalEvaluator) bool {
	for _, c := range g.Conditions {
		if c.Evaluate(fs) {
			return true
		}
	}
	return false
}

// NotGoal negates a condition.
type NotGoal struct {
	Condition GoalNode
}

func (g *NotGoal) Evaluate(fs GoalEvaluator) bool {
	return !g.Condition.Evaluate(fs)
}

// ParseGoal converts a YAML goal map to a GoalNode tree.
func ParseGoal(raw map[string]any) (GoalNode, error) {
	if len(raw) == 0 {
		return &AlwaysGoal{}, nil
	}

	// Handle single-key primitives and combinators
	for key, value := range raw {
		node, err := parseGoalKey(key, value)
		if err != nil {
			return nil, err
		}
		if node != nil {
			return node, nil
		}
	}

	return nil, fmt.Errorf("empty goal")
}

func parseGoalKey(key string, value any) (GoalNode, error) {
	switch key {
	case "always":
		return &AlwaysGoal{}, nil
	case "pwd_equals":
		return parseStringGoal(key, value, func(s string) GoalNode { return &PwdEqualsGoal{Path: s} })
	case "path_exists":
		return parseStringGoal(key, value, func(s string) GoalNode { return &PathExistsGoal{Path: s} })
	case "path_not_exists":
		return parseStringGoal(key, value, func(s string) GoalNode { return &PathNotExistsGoal{Path: s} })
	case "is_dir":
		return parseStringGoal(key, value, func(s string) GoalNode { return &IsDirGoal{Path: s} })
	case "is_file":
		return parseStringGoal(key, value, func(s string) GoalNode { return &IsFileGoal{Path: s} })
	case "file_contains":
		return parseFileContains(value)
	case "and":
		return parseAnd(value)
	case "or":
		return parseOr(value)
	case "not":
		return parseNot(value)
	default:
		return nil, fmt.Errorf("unknown goal operation: %s", key)
	}
}

func parseStringGoal(key string, value any, create func(string) GoalNode) (GoalNode, error) {
	path, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%s expects string, got %T", key, value)
	}
	return create(path), nil
}

func parseFileContains(value any) (GoalNode, error) {
	m, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("file_contains expects map with path and content, got %T", value)
	}

	path, ok := m["path"].(string)
	if !ok {
		return nil, fmt.Errorf("file_contains.path expects string")
	}

	content, ok := m["content"].(string)
	if !ok {
		return nil, fmt.Errorf("file_contains.content expects string")
	}

	return &FileContainsGoal{Path: path, Content: content}, nil
}

func parseAnd(value any) (GoalNode, error) {
	items, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("and expects array, got %T", value)
	}

	conditions := make([]GoalNode, 0, len(items))
	for i, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("and[%d] expects map, got %T", i, item)
		}
		node, err := ParseGoal(m)
		if err != nil {
			return nil, fmt.Errorf("and[%d]: %w", i, err)
		}
		conditions = append(conditions, node)
	}

	return &AndGoal{Conditions: conditions}, nil
}

func parseOr(value any) (GoalNode, error) {
	items, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("or expects array, got %T", value)
	}

	conditions := make([]GoalNode, 0, len(items))
	for i, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("or[%d] expects map, got %T", i, item)
		}
		node, err := ParseGoal(m)
		if err != nil {
			return nil, fmt.Errorf("or[%d]: %w", i, err)
		}
		conditions = append(conditions, node)
	}

	return &OrGoal{Conditions: conditions}, nil
}

func parseNot(value any) (GoalNode, error) {
	m, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("not expects map, got %T", value)
	}

	node, err := ParseGoal(m)
	if err != nil {
		return nil, fmt.Errorf("not: %w", err)
	}

	return &NotGoal{Condition: node}, nil
}
