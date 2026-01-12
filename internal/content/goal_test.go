// ABOUTME: Tests for Goal DSL parser and evaluator
// ABOUTME: Covers primitives, logic combinators, and edge cases

package content

import (
	"testing"
)

// mockFS is a simple GoalEvaluator implementation for tests.
type mockFS struct {
	pwd   string
	paths map[string]bool // true = directory, false = file
	files map[string]string
}

func newMockFS() *mockFS {
	return &mockFS{
		pwd:   "/",
		paths: make(map[string]bool),
		files: make(map[string]string),
	}
}

func (m *mockFS) Pwd() string {
	return m.pwd
}

func (m *mockFS) Exists(path string) bool {
	_, pathExists := m.paths[path]
	_, fileExists := m.files[path]
	return pathExists || fileExists
}

func (m *mockFS) IsDir(path string) bool {
	isDir, exists := m.paths[path]
	return exists && isDir
}

func (m *mockFS) ReadFile(path string) (string, error) {
	content, ok := m.files[path]
	if !ok {
		return "", &mockError{msg: "file not found"}
	}
	return content, nil
}

type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}

// Helper methods for setting up test state.
func (m *mockFS) setCwd(path string) {
	m.pwd = path
}

func (m *mockFS) addDir(path string) {
	m.paths[path] = true
}

func (m *mockFS) addFile(path string) {
	m.paths[path] = false
}

func (m *mockFS) writeFile(path, content string) {
	m.paths[path] = false
	m.files[path] = content
}

func (m *mockFS) remove(path string) {
	delete(m.paths, path)
	delete(m.files, path)
}

func TestAlwaysGoal(t *testing.T) {
	goal := map[string]any{"always": true}
	node, err := ParseGoal(goal)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}

	fs := newMockFS()
	if !node.Evaluate(fs) {
		t.Error("always goal should return true")
	}
}

func TestPwdEqualsGoal(t *testing.T) {
	goal := map[string]any{"pwd_equals": "/home/learner/projects"}
	node, err := ParseGoal(goal)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}

	fs := newMockFS()
	fs.addDir("/home/learner/projects")
	fs.setCwd("/home/learner/projects")

	if !node.Evaluate(fs) {
		t.Errorf("pwd_equals should match, pwd=%s", fs.Pwd())
	}

	fs.setCwd("/")
	if node.Evaluate(fs) {
		t.Error("pwd_equals should not match after cd /")
	}
}

func TestPathExistsGoal(t *testing.T) {
	goal := map[string]any{"path_exists": "/test/file.txt"}
	node, err := ParseGoal(goal)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}

	fs := newMockFS()

	if node.Evaluate(fs) {
		t.Error("path_exists should return false for non-existent path")
	}

	fs.addDir("/test")
	fs.addFile("/test/file.txt")

	if !node.Evaluate(fs) {
		t.Error("path_exists should return true after file created")
	}
}

func TestPathNotExistsGoal(t *testing.T) {
	goal := map[string]any{"path_not_exists": "/should/not/exist"}
	node, err := ParseGoal(goal)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}

	fs := newMockFS()

	if !node.Evaluate(fs) {
		t.Error("path_not_exists should return true for non-existent path")
	}

	fs.addDir("/should/not")
	fs.addFile("/should/not/exist")

	if node.Evaluate(fs) {
		t.Error("path_not_exists should return false after path created")
	}
}

func TestIsDirGoal(t *testing.T) {
	goal := map[string]any{"is_dir": "/mydir"}
	node, err := ParseGoal(goal)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}

	fs := newMockFS()

	if node.Evaluate(fs) {
		t.Error("is_dir should return false for non-existent path")
	}

	fs.addFile("/mydir") // Create as file
	if node.Evaluate(fs) {
		t.Error("is_dir should return false for file")
	}

	fs.remove("/mydir")
	fs.addDir("/mydir")
	if !node.Evaluate(fs) {
		t.Error("is_dir should return true for directory")
	}
}

func TestIsFileGoal(t *testing.T) {
	fs := newMockFS()

	// Test non-existent path
	goal1 := map[string]any{"is_file": "/nonexistent"}
	node1, err := ParseGoal(goal1)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}
	if node1.Evaluate(fs) {
		t.Error("is_file should return false for non-existent path")
	}

	// Test directory
	goal2 := map[string]any{"is_file": "/mydir"}
	node2, err := ParseGoal(goal2)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}
	fs.addDir("/mydir")
	if node2.Evaluate(fs) {
		t.Error("is_file should return false for directory")
	}

	// Test file
	goal3 := map[string]any{"is_file": "/myfile"}
	node3, err := ParseGoal(goal3)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}
	fs.addFile("/myfile")
	if !node3.Evaluate(fs) {
		t.Error("is_file should return true for file")
	}
}

func TestFileContainsGoal(t *testing.T) {
	goal := map[string]any{
		"file_contains": map[string]any{
			"path":    "/test.txt",
			"content": "Hello World",
		},
	}
	node, err := ParseGoal(goal)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}

	fs := newMockFS()

	if node.Evaluate(fs) {
		t.Error("file_contains should return false for non-existent file")
	}

	fs.writeFile("/test.txt", "Goodbye")
	if node.Evaluate(fs) {
		t.Error("file_contains should return false when content not present")
	}

	fs.writeFile("/test.txt", "Hello World!")
	if !node.Evaluate(fs) {
		t.Error("file_contains should return true when content present")
	}
}

func TestAndGoal(t *testing.T) {
	goal := map[string]any{
		"and": []any{
			map[string]any{"path_exists": "/a"},
			map[string]any{"path_exists": "/b"},
		},
	}
	node, err := ParseGoal(goal)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}

	fs := newMockFS()

	if node.Evaluate(fs) {
		t.Error("and should return false when no conditions met")
	}

	fs.addDir("/a")
	if node.Evaluate(fs) {
		t.Error("and should return false when only one condition met")
	}

	fs.addDir("/b")
	if !node.Evaluate(fs) {
		t.Error("and should return true when all conditions met")
	}
}

func TestOrGoal(t *testing.T) {
	goal := map[string]any{
		"or": []any{
			map[string]any{"path_exists": "/a"},
			map[string]any{"path_exists": "/b"},
		},
	}
	node, err := ParseGoal(goal)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}

	fs := newMockFS()

	if node.Evaluate(fs) {
		t.Error("or should return false when no conditions met")
	}

	fs.addDir("/a")
	if !node.Evaluate(fs) {
		t.Error("or should return true when one condition met")
	}

	fs.addDir("/b")
	if !node.Evaluate(fs) {
		t.Error("or should return true when all conditions met")
	}
}

func TestNotGoal(t *testing.T) {
	goal := map[string]any{
		"not": map[string]any{"path_exists": "/deleted"},
	}
	node, err := ParseGoal(goal)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}

	fs := newMockFS()

	if !node.Evaluate(fs) {
		t.Error("not should return true when inner condition false")
	}

	fs.addDir("/deleted")
	if node.Evaluate(fs) {
		t.Error("not should return false when inner condition true")
	}
}

func TestNestedCombinators(t *testing.T) {
	// Test: file moved from downloads to documents
	// and:
	//   - path_exists: /documents/report.pdf
	//   - not:
	//       path_exists: /downloads/report.pdf
	goal := map[string]any{
		"and": []any{
			map[string]any{"path_exists": "/documents/report.pdf"},
			map[string]any{
				"not": map[string]any{"path_exists": "/downloads/report.pdf"},
			},
		},
	}
	node, err := ParseGoal(goal)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}

	fs := newMockFS()
	fs.addDir("/downloads")
	fs.addDir("/documents")
	fs.addFile("/downloads/report.pdf")

	if node.Evaluate(fs) {
		t.Error("should fail: file not moved yet")
	}

	// Move the file (remove from downloads, add to documents)
	fs.remove("/downloads/report.pdf")
	fs.addFile("/documents/report.pdf")

	if !node.Evaluate(fs) {
		t.Error("should pass: file moved successfully")
	}
}

func TestParseGoalErrors(t *testing.T) {
	tests := []struct {
		name string
		goal map[string]any
	}{
		{
			name: "unknown operation",
			goal: map[string]any{"invalid_op": "value"},
		},
		{
			name: "pwd_equals wrong type",
			goal: map[string]any{"pwd_equals": 123},
		},
		{
			name: "and not array",
			goal: map[string]any{"and": "not an array"},
		},
		{
			name: "file_contains missing path",
			goal: map[string]any{
				"file_contains": map[string]any{"content": "test"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseGoal(tt.goal)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestEmptyGoal(t *testing.T) {
	goal := map[string]any{}
	node, err := ParseGoal(goal)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}

	fs := newMockFS()
	if !node.Evaluate(fs) {
		t.Error("empty goal should return true (like always)")
	}
}

func TestEmptyAndGoal(t *testing.T) {
	goal := map[string]any{
		"and": []any{},
	}
	node, err := ParseGoal(goal)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}

	fs := newMockFS()
	// Empty AND should return true (vacuously true)
	if !node.Evaluate(fs) {
		t.Error("empty and should return true (vacuously true)")
	}
}

func TestEmptyOrGoal(t *testing.T) {
	goal := map[string]any{
		"or": []any{},
	}
	node, err := ParseGoal(goal)
	if err != nil {
		t.Fatalf("ParseGoal failed: %v", err)
	}

	fs := newMockFS()
	// Empty OR should return false (vacuously false)
	if node.Evaluate(fs) {
		t.Error("empty or should return false (vacuously false)")
	}
}
