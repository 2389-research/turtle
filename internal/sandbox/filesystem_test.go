// ABOUTME: Tests for the simulated filesystem
// ABOUTME: Ensures safe sandbox behavior for learning

package sandbox

import (
	"testing"
)

func TestNewFilesystem(t *testing.T) {
	fs := NewFilesystem()

	if fs.Root == nil {
		t.Error("Root should not be nil")
	}
	if fs.Pwd() != "/" {
		t.Errorf("Expected pwd /, got %s", fs.Pwd())
	}
}

func TestMkdir(t *testing.T) {
	fs := NewFilesystem()

	err := fs.Mkdir("/test/nested/dir")
	if err != nil {
		t.Errorf("Mkdir failed: %v", err)
	}

	if !fs.Exists("/test/nested/dir") {
		t.Error("/test/nested/dir should exist")
	}
	if !fs.IsDir("/test/nested/dir") {
		t.Error("/test/nested/dir should be a directory")
	}
}

func TestTouch(t *testing.T) {
	fs := NewFilesystem()

	err := fs.Touch("/test/file.txt")
	if err != nil {
		t.Errorf("Touch failed: %v", err)
	}

	if !fs.Exists("/test/file.txt") {
		t.Error("/test/file.txt should exist")
	}
	if fs.IsDir("/test/file.txt") {
		t.Error("/test/file.txt should not be a directory")
	}
}

func TestCd(t *testing.T) {
	fs := NewFilesystem()
	_ = fs.Mkdir("/home/user/projects")

	err := fs.Cd("/home/user/projects")
	if err != nil {
		t.Errorf("Cd failed: %v", err)
	}
	if fs.Pwd() != "/home/user/projects" {
		t.Errorf("Expected /home/user/projects, got %s", fs.Pwd())
	}

	// Test cd ..
	err = fs.Cd("..")
	if err != nil {
		t.Errorf("Cd .. failed: %v", err)
	}
	if fs.Pwd() != "/home/user" {
		t.Errorf("Expected /home/user, got %s", fs.Pwd())
	}

	// Test cd ~ (home)
	fs.Home = "/home/user"
	err = fs.Cd("~")
	if err != nil {
		t.Errorf("Cd ~ failed: %v", err)
	}
	if fs.Pwd() != "/home/user" {
		t.Errorf("Expected /home/user, got %s", fs.Pwd())
	}
}

func TestLs(t *testing.T) {
	fs := NewFilesystem()
	_ = fs.Mkdir("/test")
	_ = fs.Touch("/test/file1.txt")
	_ = fs.Touch("/test/file2.txt")
	_ = fs.Touch("/test/.hidden")

	// Without hidden
	files, err := fs.Ls("/test", false)
	if err != nil {
		t.Errorf("Ls failed: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}

	// With hidden
	files, err = fs.Ls("/test", true)
	if err != nil {
		t.Errorf("Ls failed: %v", err)
	}
	if len(files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(files))
	}
}

func TestWriteAndReadFile(t *testing.T) {
	fs := NewFilesystem()

	content := "Hello, World!"
	err := fs.WriteFile("/test.txt", content)
	if err != nil {
		t.Errorf("WriteFile failed: %v", err)
	}

	read, err := fs.ReadFile("/test.txt")
	if err != nil {
		t.Errorf("ReadFile failed: %v", err)
	}
	if read != content {
		t.Errorf("Expected %q, got %q", content, read)
	}
}

func TestCp(t *testing.T) {
	fs := NewFilesystem()
	_ = fs.WriteFile("/original.txt", "content")

	err := fs.Cp("/original.txt", "/copy.txt")
	if err != nil {
		t.Errorf("Cp failed: %v", err)
	}

	if !fs.Exists("/copy.txt") {
		t.Error("/copy.txt should exist")
	}

	content, _ := fs.ReadFile("/copy.txt")
	if content != "content" {
		t.Errorf("Expected 'content', got %q", content)
	}
}

func TestMv(t *testing.T) {
	fs := NewFilesystem()
	_ = fs.WriteFile("/old.txt", "content")

	err := fs.Mv("/old.txt", "/new.txt")
	if err != nil {
		t.Errorf("Mv failed: %v", err)
	}

	if fs.Exists("/old.txt") {
		t.Error("/old.txt should not exist")
	}
	if !fs.Exists("/new.txt") {
		t.Error("/new.txt should exist")
	}
}

func TestRm(t *testing.T) {
	fs := NewFilesystem()
	_ = fs.Touch("/delete.txt")

	err := fs.Rm("/delete.txt")
	if err != nil {
		t.Errorf("Rm failed: %v", err)
	}

	if fs.Exists("/delete.txt") {
		t.Error("/delete.txt should not exist")
	}
}

func TestRmDirectory(t *testing.T) {
	fs := NewFilesystem()
	_ = fs.Mkdir("/dir")

	err := fs.Rm("/dir")
	if err == nil {
		t.Error("Rm on directory should fail")
	}
}

func TestGrep(t *testing.T) {
	fs := NewFilesystem()
	_ = fs.WriteFile("/log.txt", "line 1\nERROR: something failed\nline 3\n")

	matches, err := fs.Grep("ERROR", "/log.txt")
	if err != nil {
		t.Errorf("Grep failed: %v", err)
	}
	if len(matches) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matches))
	}
}

func TestFind(t *testing.T) {
	fs := NewFilesystem()
	_ = fs.Mkdir("/a/b/c")
	_ = fs.Touch("/a/test.txt")
	_ = fs.Touch("/a/b/test.txt")
	_ = fs.Touch("/a/b/c/other.txt")

	results, err := fs.Find("/a", "test.txt")
	if err != nil {
		t.Errorf("Find failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d: %v", len(results), results)
	}
}

func TestClone(t *testing.T) {
	fs := NewFilesystem()
	_ = fs.Mkdir("/test")
	_ = fs.Touch("/test/file.txt")
	_ = fs.Cd("/test")

	clone := fs.Clone()

	// Modify original
	_ = fs.Touch("/test/new.txt")

	// Clone should not have new file
	if clone.Exists("/test/new.txt") {
		t.Error("Clone should be independent")
	}
	if clone.Pwd() != "/test" {
		t.Errorf("Clone should preserve cwd, got %s", clone.Pwd())
	}
}

func TestDefaultFilesystem(t *testing.T) {
	fs := NewDefaultFilesystem()

	// Should start in home
	if fs.Pwd() != "/home/learner" {
		t.Errorf("Expected /home/learner, got %s", fs.Pwd())
	}

	// Should have basic structure
	if !fs.Exists("/home/learner/projects") {
		t.Error("/home/learner/projects should exist")
	}
	if !fs.Exists("/tmp") {
		t.Error("/tmp should exist")
	}
}
