// ABOUTME: Simulated filesystem for safe terminal learning
// ABOUTME: Provides a fake filesystem that learners can explore and modify without risk

package sandbox

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// FileType represents the type of filesystem entry.
type FileType int

const (
	FileTypeRegular FileType = iota
	FileTypeDirectory
	FileTypeHidden // starts with .
)

// File represents a file or directory in the sandbox.
type File struct {
	Name     string
	Type     FileType
	Content  string    // For regular files
	Children []*File   // For directories
	Parent   *File     // Parent directory (nil for root)
	ModTime  time.Time // Modification time
	Size     int       // File size in bytes
}

// Filesystem represents the complete sandbox environment.
type Filesystem struct {
	Root    *File
	Cwd     *File  // Current working directory
	CwdPath string // Current path as string
	Home    string // Home directory path
	User    string // Current username
}

// NewFilesystem creates a new empty filesystem.
func NewFilesystem() *Filesystem {
	root := &File{
		Name:     "/",
		Type:     FileTypeDirectory,
		Children: []*File{},
		ModTime:  time.Now(),
	}

	fs := &Filesystem{
		Root:    root,
		Cwd:     root,
		CwdPath: "/",
		Home:    "/home/learner",
		User:    "learner",
	}

	return fs
}

// NewDefaultFilesystem creates a filesystem with a realistic structure for learning.
func NewDefaultFilesystem() *Filesystem {
	fs := NewFilesystem()

	// Build basic structure (errors ignored - internal setup on fresh filesystem)
	_ = fs.Mkdir("/home")
	_ = fs.Mkdir("/home/learner")
	_ = fs.Mkdir("/home/learner/projects")
	_ = fs.Mkdir("/home/learner/documents")
	_ = fs.Mkdir("/home/learner/downloads")
	_ = fs.Mkdir("/tmp")
	_ = fs.Mkdir("/var")
	_ = fs.Mkdir("/var/log")
	_ = fs.Mkdir("/etc")

	// Add some starter files (errors ignored - internal setup)
	_ = fs.Touch("/home/learner/.bashrc")
	_ = fs.WriteFile("/home/learner/.bashrc", "# Bash configuration\nexport PATH=$PATH:~/bin\n")

	_ = fs.Touch("/home/learner/readme.txt")
	_ = fs.WriteFile("/home/learner/readme.txt", "Welcome to the terminal!\nThis is your home directory.\n")

	_ = fs.Touch("/etc/passwd")
	_ = fs.WriteFile("/etc/passwd", "root:x:0:0:root:/root:/bin/bash\nlearner:x:1000:1000::/home/learner:/bin/bash\n")

	// Set cwd to home
	_ = fs.Cd("/home/learner")

	return fs
}

// Mkdir creates a directory (and parents if needed).
func (fs *Filesystem) Mkdir(path string) error {
	path = fs.resolvePath(path)
	parts := splitPath(path)

	current := fs.Root
	for _, part := range parts {
		if part == "" {
			continue
		}

		child := findChild(current, part)
		switch {
		case child == nil:
			// Create new directory
			newDir := &File{
				Name:     part,
				Type:     FileTypeDirectory,
				Children: []*File{},
				Parent:   current,
				ModTime:  time.Now(),
			}
			if strings.HasPrefix(part, ".") {
				newDir.Type = FileTypeHidden
			}
			current.Children = append(current.Children, newDir)
			current = newDir
		case child.Type == FileTypeDirectory || child.Type == FileTypeHidden:
			current = child
		default:
			return fmt.Errorf("not a directory: %s", part)
		}
	}

	return nil
}

// Touch creates an empty file or updates modification time.
func (fs *Filesystem) Touch(path string) error {
	path = fs.resolvePath(path)
	dir := filepath.Dir(path)
	name := filepath.Base(path)

	// Ensure parent directory exists
	if err := fs.Mkdir(dir); err != nil {
		return err
	}

	parent, err := fs.getNode(dir)
	if err != nil {
		return err
	}

	// Check if file exists
	existing := findChild(parent, name)
	if existing != nil {
		existing.ModTime = time.Now()
		return nil
	}

	// Create new file
	fileType := FileTypeRegular
	if strings.HasPrefix(name, ".") {
		fileType = FileTypeHidden
	}

	newFile := &File{
		Name:    name,
		Type:    fileType,
		Content: "",
		Parent:  parent,
		ModTime: time.Now(),
		Size:    0,
	}
	parent.Children = append(parent.Children, newFile)

	return nil
}

// WriteFile writes content to a file (creates if doesn't exist).
func (fs *Filesystem) WriteFile(path, content string) error {
	if err := fs.Touch(path); err != nil {
		return err
	}

	node, err := fs.getNode(fs.resolvePath(path))
	if err != nil {
		return err
	}

	node.Content = content
	node.Size = len(content)
	node.ModTime = time.Now()

	return nil
}

// ReadFile reads content from a file.
func (fs *Filesystem) ReadFile(path string) (string, error) {
	node, err := fs.getNode(fs.resolvePath(path))
	if err != nil {
		return "", err
	}

	if node.Type == FileTypeDirectory {
		return "", fmt.Errorf("is a directory: %s", path)
	}

	return node.Content, nil
}

// Cd changes the current working directory.
func (fs *Filesystem) Cd(path string) error {
	if path == "" || path == "~" {
		path = fs.Home
	} else if strings.HasPrefix(path, "~/") {
		path = fs.Home + path[1:]
	}

	resolved := fs.resolvePath(path)
	node, err := fs.getNode(resolved)
	if err != nil {
		return err
	}

	if node.Type != FileTypeDirectory && node.Type != FileTypeHidden {
		return fmt.Errorf("not a directory: %s", path)
	}

	fs.Cwd = node
	fs.CwdPath = resolved

	return nil
}

// Pwd returns the current working directory path.
func (fs *Filesystem) Pwd() string {
	return fs.CwdPath
}

// Ls lists directory contents.
func (fs *Filesystem) Ls(path string, showHidden bool) ([]string, error) {
	if path == "" {
		path = fs.CwdPath
	}

	node, err := fs.getNode(fs.resolvePath(path))
	if err != nil {
		return nil, err
	}

	if node.Type != FileTypeDirectory && node.Type != FileTypeHidden {
		return []string{node.Name}, nil
	}

	var names []string
	for _, child := range node.Children {
		isHidden := strings.HasPrefix(child.Name, ".")
		if !isHidden || showHidden {
			name := child.Name
			if child.Type == FileTypeDirectory {
				name += "/"
			}
			names = append(names, name)
		}
	}

	sort.Strings(names)
	return names, nil
}

// Rm removes a file.
func (fs *Filesystem) Rm(path string) error {
	resolved := fs.resolvePath(path)
	node, err := fs.getNode(resolved)
	if err != nil {
		return err
	}

	if node.Type == FileTypeDirectory {
		return fmt.Errorf("is a directory (use rm -r): %s", path)
	}

	if node.Parent == nil {
		return fmt.Errorf("cannot remove root")
	}

	// Remove from parent's children
	parent := node.Parent
	for i, child := range parent.Children {
		if child == node {
			parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
			break
		}
	}

	return nil
}

// Cp copies a file.
func (fs *Filesystem) Cp(src, dst string) error {
	srcNode, err := fs.getNode(fs.resolvePath(src))
	if err != nil {
		return err
	}

	if srcNode.Type == FileTypeDirectory {
		return fmt.Errorf("is a directory (use cp -r): %s", src)
	}

	return fs.WriteFile(dst, srcNode.Content)
}

// Mv moves/renames a file or directory.
func (fs *Filesystem) Mv(src, dst string) error {
	srcResolved := fs.resolvePath(src)
	srcNode, err := fs.getNode(srcResolved)
	if err != nil {
		return err
	}

	// Check if dst is a directory
	dstResolved := fs.resolvePath(dst)
	dstNode, _ := fs.getNode(dstResolved)
	if dstNode != nil && dstNode.Type == FileTypeDirectory {
		// Move into directory
		dstResolved = filepath.Join(dstResolved, srcNode.Name)
	}

	// Remove from old parent
	if srcNode.Parent != nil {
		parent := srcNode.Parent
		for i, child := range parent.Children {
			if child == srcNode {
				parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
				break
			}
		}
	}

	// Add to new parent
	dstDir := filepath.Dir(dstResolved)
	dstName := filepath.Base(dstResolved)

	if err := fs.Mkdir(dstDir); err != nil {
		return err
	}

	newParent, err := fs.getNode(dstDir)
	if err != nil {
		return err
	}

	srcNode.Name = dstName
	srcNode.Parent = newParent
	newParent.Children = append(newParent.Children, srcNode)

	return nil
}

// Exists checks if a path exists.
func (fs *Filesystem) Exists(path string) bool {
	_, err := fs.getNode(fs.resolvePath(path))
	return err == nil
}

// IsDir checks if path is a directory.
func (fs *Filesystem) IsDir(path string) bool {
	node, err := fs.getNode(fs.resolvePath(path))
	if err != nil {
		return false
	}
	return node.Type == FileTypeDirectory || node.Type == FileTypeHidden
}

// Cat returns file contents (for display).
func (fs *Filesystem) Cat(path string) (string, error) {
	return fs.ReadFile(path)
}

// Grep searches for pattern in file.
func (fs *Filesystem) Grep(pattern, path string) ([]string, error) {
	content, err := fs.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var matches []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, pattern) {
			matches = append(matches, line)
		}
	}

	return matches, nil
}

// Find searches for files matching a pattern.
func (fs *Filesystem) Find(startPath, pattern string) ([]string, error) {
	start, err := fs.getNode(fs.resolvePath(startPath))
	if err != nil {
		return nil, err
	}

	var results []string
	fs.findRecursive(start, fs.resolvePath(startPath), pattern, &results)
	return results, nil
}

func (fs *Filesystem) findRecursive(node *File, currentPath, pattern string, results *[]string) {
	matched, _ := filepath.Match(pattern, node.Name)
	if matched {
		*results = append(*results, currentPath)
	}

	if node.Type == FileTypeDirectory || node.Type == FileTypeHidden {
		for _, child := range node.Children {
			childPath := filepath.Join(currentPath, child.Name)
			fs.findRecursive(child, childPath, pattern, results)
		}
	}
}

// Clone creates a deep copy of the filesystem (for reset).
func (fs *Filesystem) Clone() *Filesystem {
	newFs := &Filesystem{
		Home: fs.Home,
		User: fs.User,
	}

	newFs.Root = cloneFile(fs.Root, nil)
	newFs.Cwd = newFs.Root
	newFs.CwdPath = "/"

	// Navigate to same cwd (error ignored - path exists in cloned fs)
	_ = newFs.Cd(fs.CwdPath)

	return newFs
}

func cloneFile(f *File, parent *File) *File {
	if f == nil {
		return nil
	}

	clone := &File{
		Name:    f.Name,
		Type:    f.Type,
		Content: f.Content,
		Parent:  parent,
		ModTime: f.ModTime,
		Size:    f.Size,
	}

	for _, child := range f.Children {
		clone.Children = append(clone.Children, cloneFile(child, clone))
	}

	return clone
}

// resolvePath converts relative paths to absolute.
func (fs *Filesystem) resolvePath(path string) string {
	if path == "" {
		return fs.CwdPath
	}

	// Handle home directory
	if path == "~" {
		return fs.Home
	}
	if strings.HasPrefix(path, "~/") {
		path = fs.Home + path[1:]
	}

	// Handle absolute paths
	if strings.HasPrefix(path, "/") {
		return filepath.Clean(path)
	}

	// Handle relative paths
	return filepath.Clean(filepath.Join(fs.CwdPath, path))
}

// getNode finds a node by path.
func (fs *Filesystem) getNode(path string) (*File, error) {
	if path == "/" {
		return fs.Root, nil
	}

	parts := splitPath(path)
	current := fs.Root

	for _, part := range parts {
		if part == "" {
			continue
		}

		child := findChild(current, part)
		if child == nil {
			return nil, fmt.Errorf("no such file or directory: %s", path)
		}
		current = child
	}

	return current, nil
}

func splitPath(path string) []string {
	path = filepath.Clean(path)
	if path == "/" {
		return []string{}
	}
	return strings.Split(strings.TrimPrefix(path, "/"), "/")
}

func findChild(dir *File, name string) *File {
	for _, child := range dir.Children {
		if child.Name == name {
			return child
		}
	}
	return nil
}
