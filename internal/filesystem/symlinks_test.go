package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCreateSymlink_SiblingDirectories tests that symlinks are created correctly
// when source and target directories are siblings (in the same parent directory)
func TestCreateSymlink_SiblingDirectories(t *testing.T) {
	// Create a temporary directory structure:
	// temp/
	//   ├── services-available/
	//   │   └── test-file.yml
	//   └── services-enabled/
	//       └── test-file.yml -> ../services-available/test-file.yml

	tempDir, err := os.MkdirTemp("", "lnka-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "services-available")
	targetDir := filepath.Join(tempDir, "services-enabled")

	// Create directories
	if err := os.Mkdir(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	if err := os.Mkdir(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	// Create a test file in source directory
	testFile := "test-file.yml"
	sourceFile := filepath.Join(sourceDir, testFile)
	if err := os.WriteFile(sourceFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create the symlink
	if err := CreateSymlink(sourceDir, targetDir, testFile); err != nil {
		t.Fatalf("CreateSymlink failed: %v", err)
	}

	// Verify the symlink was created
	linkPath := filepath.Join(targetDir, testFile)
	linkTarget, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}

	// The symlink should be relative and start with ../
	expectedTarget := filepath.Join("..", "services-available", testFile)
	if linkTarget != expectedTarget {
		t.Errorf("Symlink target incorrect:\n  got:  %q\n  want: %q", linkTarget, expectedTarget)
	}

	// Verify the symlink actually works (can resolve to the source file)
	resolvedPath := filepath.Join(targetDir, linkTarget)
	resolvedAbs, err := filepath.Abs(resolvedPath)
	if err != nil {
		t.Fatalf("Failed to resolve symlink path: %v", err)
	}

	sourceAbs, err := filepath.Abs(sourceFile)
	if err != nil {
		t.Fatalf("Failed to get absolute source path: %v", err)
	}

	if resolvedAbs != sourceAbs {
		t.Errorf("Symlink doesn't resolve to source file:\n  resolved: %s\n  source:   %s", resolvedAbs, sourceAbs)
	}

	// Verify we can actually read the file through the symlink
	content, err := os.ReadFile(linkPath)
	if err != nil {
		t.Errorf("Failed to read through symlink: %v", err)
	}
	if string(content) != "test content" {
		t.Errorf("Content mismatch: got %q, want %q", string(content), "test content")
	}
}

// TestCreateSymlink_NestedDirectories tests symlink creation with nested directories
func TestCreateSymlink_NestedDirectories(t *testing.T) {
	// Create a more complex directory structure:
	// temp/
	//   ├── config/
	//   │   └── available/
	//   │       └── test.conf
	//   └── active/
	//       └── test.conf -> ../config/available/test.conf

	tempDir, err := os.MkdirTemp("", "lnka-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "config", "available")
	targetDir := filepath.Join(tempDir, "active")

	// Create directories
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	if err := os.Mkdir(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	// Create a test file
	testFile := "test.conf"
	sourceFile := filepath.Join(sourceDir, testFile)
	if err := os.WriteFile(sourceFile, []byte("config data"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create the symlink
	if err := CreateSymlink(sourceDir, targetDir, testFile); err != nil {
		t.Fatalf("CreateSymlink failed: %v", err)
	}

	// Verify the symlink
	linkPath := filepath.Join(targetDir, testFile)
	linkTarget, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}

	// Should be a relative path
	expectedTarget := filepath.Join("..", "config", "available", testFile)
	if linkTarget != expectedTarget {
		t.Errorf("Symlink target incorrect:\n  got:  %q\n  want: %q", linkTarget, expectedTarget)
	}

	// Verify the symlink resolves correctly
	content, err := os.ReadFile(linkPath)
	if err != nil {
		t.Errorf("Failed to read through symlink: %v", err)
	}
	if string(content) != "config data" {
		t.Errorf("Content mismatch: got %q, want %q", string(content), "config data")
	}
}

// TestListAvailableFiles tests listing files in a directory
func TestListAvailableFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create some test files and directories
	files := []string{"file1.txt", "file2.yml", "config.json"}
	for _, f := range files {
		path := filepath.Join(tempDir, f)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", f, err)
		}
	}

	// Create a directory (should be ignored)
	if err := os.Mkdir(filepath.Join(tempDir, "subdir"), 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// List available files
	result, err := ListAvailableFiles(tempDir)
	if err != nil {
		t.Fatalf("ListAvailableFiles failed: %v", err)
	}

	// Verify count (should only include files, not the directory)
	if len(result) != 3 {
		t.Errorf("Expected 3 files, got %d: %v", len(result), result)
	}

	// Verify all files are present
	fileMap := make(map[string]bool)
	for _, f := range result {
		fileMap[f] = true
	}

	for _, expected := range files {
		if !fileMap[expected] {
			t.Errorf("Expected file %s not found in result", expected)
		}
	}

	// Verify directory is not included
	if fileMap["subdir"] {
		t.Error("Directory 'subdir' should not be included in file list")
	}
}

// TestListAvailableFiles_NonExistentDir tests error handling
func TestListAvailableFiles_NonExistentDir(t *testing.T) {
	_, err := ListAvailableFiles("/nonexistent/directory")
	if err == nil {
		t.Error("Expected error for non-existent directory, got nil")
	}
}

// TestListEnabledSymlinks tests listing symlinks in target directory
func TestListEnabledSymlinks(t *testing.T) {
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	targetDir := filepath.Join(tempDir, "target")

	// Create directories
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	// Create source files
	files := []string{"file1.txt", "file2.yml"}
	for _, f := range files {
		path := filepath.Join(sourceDir, f)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create source file %s: %v", f, err)
		}
	}

	// Create symlinks
	for _, f := range files {
		sourcePath := filepath.Join(sourceDir, f)
		linkPath := filepath.Join(targetDir, f)
		relPath, _ := filepath.Rel(targetDir, sourcePath)
		if err := os.Symlink(relPath, linkPath); err != nil {
			t.Fatalf("Failed to create symlink for %s: %v", f, err)
		}
	}

	// Create a regular file (should be ignored)
	regularFile := filepath.Join(targetDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("not a symlink"), 0644); err != nil {
		t.Fatalf("Failed to create regular file: %v", err)
	}

	// List enabled symlinks
	result, err := ListEnabledSymlinks(sourceDir, targetDir)
	if err != nil {
		t.Fatalf("ListEnabledSymlinks failed: %v", err)
	}

	// Verify count (only symlinks, not regular file)
	if len(result) != 2 {
		t.Errorf("Expected 2 symlinks, got %d: %v", len(result), result)
	}

	// Verify symlinks are present
	for _, f := range files {
		if _, exists := result[f]; !exists {
			t.Errorf("Expected symlink %s not found", f)
		}
	}

	// Verify regular file is not included
	if _, exists := result["regular.txt"]; exists {
		t.Error("Regular file should not be included in symlink list")
	}
}

// TestGetEnabledFiles tests getting list of enabled files
func TestGetEnabledFiles(t *testing.T) {
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	targetDir := filepath.Join(tempDir, "target")

	// Create directories
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	// Create source files
	files := []string{"enabled1.txt", "enabled2.yml"}
	for _, f := range files {
		path := filepath.Join(sourceDir, f)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create source file %s: %v", f, err)
		}
	}

	// Create symlinks pointing to source files
	for _, f := range files {
		if err := CreateSymlink(sourceDir, targetDir, f); err != nil {
			t.Fatalf("Failed to create symlink for %s: %v", f, err)
		}
	}

	// Create a symlink pointing to a different location (should be ignored)
	otherFile := filepath.Join(tempDir, "other.txt")
	if err := os.WriteFile(otherFile, []byte("other"), 0644); err != nil {
		t.Fatalf("Failed to create other file: %v", err)
	}
	otherLink := filepath.Join(targetDir, "other.txt")
	if err := os.Symlink(otherFile, otherLink); err != nil {
		t.Fatalf("Failed to create other symlink: %v", err)
	}

	// Get enabled files
	result, err := GetEnabledFiles(sourceDir, targetDir)
	if err != nil {
		t.Fatalf("GetEnabledFiles failed: %v", err)
	}

	// Verify count (only files pointing to sourceDir)
	if len(result) != 2 {
		t.Errorf("Expected 2 enabled files, got %d: %v", len(result), result)
	}

	// Verify correct files are enabled
	enabledMap := make(map[string]bool)
	for _, f := range result {
		enabledMap[f] = true
	}

	for _, expected := range files {
		if !enabledMap[expected] {
			t.Errorf("Expected enabled file %s not found", expected)
		}
	}

	// Verify other.txt is not included (points elsewhere)
	if enabledMap["other.txt"] {
		t.Error("Symlink pointing outside sourceDir should not be included")
	}
}

// TestRemoveSymlink tests removing a symlink
func TestRemoveSymlink(t *testing.T) {
	tempDir := t.TempDir()
	targetDir := filepath.Join(tempDir, "target")

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	// Create a symlink
	linkPath := filepath.Join(targetDir, "testlink.txt")
	if err := os.Symlink("/tmp/source.txt", linkPath); err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	// Remove the symlink
	if err := RemoveSymlink(targetDir, "testlink.txt"); err != nil {
		t.Fatalf("RemoveSymlink failed: %v", err)
	}

	// Verify symlink is gone
	if _, err := os.Lstat(linkPath); !os.IsNotExist(err) {
		t.Error("Symlink should have been removed")
	}
}

// TestRemoveSymlink_NonExistent tests removing non-existent symlink
func TestRemoveSymlink_NonExistent(t *testing.T) {
	tempDir := t.TempDir()

	// Removing non-existent symlink should succeed (idempotent)
	err := RemoveSymlink(tempDir, "nonexistent.txt")
	if err != nil {
		t.Errorf("RemoveSymlink should be idempotent for non-existent files, got error: %v", err)
	}
}

// TestRemoveSymlink_RegularFile tests refusing to remove regular files
func TestRemoveSymlink_RegularFile(t *testing.T) {
	tempDir := t.TempDir()

	// Create a regular file
	regularFile := filepath.Join(tempDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create regular file: %v", err)
	}

	// Try to remove it (should fail)
	err := RemoveSymlink(tempDir, "regular.txt")
	if err == nil {
		t.Error("RemoveSymlink should refuse to remove regular files")
	}

	// Verify file still exists
	if _, err := os.Stat(regularFile); err != nil {
		t.Error("Regular file should not have been removed")
	}
}

// TestValidateSymlinks tests finding broken symlinks
func TestValidateSymlinks(t *testing.T) {
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	targetDir := filepath.Join(tempDir, "target")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	// Create a valid symlink
	validFile := filepath.Join(sourceDir, "valid.txt")
	if err := os.WriteFile(validFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create valid file: %v", err)
	}
	validLink := filepath.Join(targetDir, "valid.txt")
	relPath, _ := filepath.Rel(targetDir, validFile)
	if err := os.Symlink(relPath, validLink); err != nil {
		t.Fatalf("Failed to create valid symlink: %v", err)
	}

	// Create a broken symlink
	brokenLink := filepath.Join(targetDir, "broken.txt")
	if err := os.Symlink("../source/nonexistent.txt", brokenLink); err != nil {
		t.Fatalf("Failed to create broken symlink: %v", err)
	}

	// Validate symlinks
	orphaned, err := ValidateSymlinks(sourceDir, targetDir)
	if err != nil {
		t.Fatalf("ValidateSymlinks failed: %v", err)
	}

	// Should find one broken symlink
	if len(orphaned) != 1 {
		t.Errorf("Expected 1 orphaned symlink, got %d: %v", len(orphaned), orphaned)
	}

	if len(orphaned) > 0 && orphaned[0] != "broken.txt" {
		t.Errorf("Expected broken.txt to be orphaned, got %s", orphaned[0])
	}
}

// TestCleanOrphanedSymlinks tests removing broken symlinks
func TestCleanOrphanedSymlinks(t *testing.T) {
	tempDir := t.TempDir()
	targetDir := filepath.Join(tempDir, "target")

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	// Create broken symlinks
	orphaned := []string{"orphan1.txt", "orphan2.txt"}
	for _, name := range orphaned {
		linkPath := filepath.Join(targetDir, name)
		if err := os.Symlink("/nonexistent/"+name, linkPath); err != nil {
			t.Fatalf("Failed to create orphaned symlink %s: %v", name, err)
		}
	}

	// Clean orphaned symlinks
	if err := CleanOrphanedSymlinks(targetDir, orphaned); err != nil {
		t.Fatalf("CleanOrphanedSymlinks failed: %v", err)
	}

	// Verify symlinks are removed
	for _, name := range orphaned {
		linkPath := filepath.Join(targetDir, name)
		if _, err := os.Lstat(linkPath); !os.IsNotExist(err) {
			t.Errorf("Orphaned symlink %s should have been removed", name)
		}
	}
}

// TestApplyChanges tests creating and removing symlinks based on selection
func TestApplyChanges(t *testing.T) {
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	targetDir := filepath.Join(tempDir, "target")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	// Create source files
	allFiles := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, f := range allFiles {
		path := filepath.Join(sourceDir, f)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create source file %s: %v", f, err)
		}
	}

	// Initially enable file1 and file2
	for _, f := range []string{"file1.txt", "file2.txt"} {
		if err := CreateSymlink(sourceDir, targetDir, f); err != nil {
			t.Fatalf("Failed to create initial symlink for %s: %v", f, err)
		}
	}

	// Apply changes: keep file1, remove file2, add file3
	selectedFiles := []string{"file1.txt", "file3.txt"}
	if err := ApplyChanges(sourceDir, targetDir, selectedFiles); err != nil {
		t.Fatalf("ApplyChanges failed: %v", err)
	}

	// Verify file1 still exists
	link1 := filepath.Join(targetDir, "file1.txt")
	if _, err := os.Lstat(link1); err != nil {
		t.Error("file1.txt symlink should still exist")
	}

	// Verify file2 was removed
	link2 := filepath.Join(targetDir, "file2.txt")
	if _, err := os.Lstat(link2); !os.IsNotExist(err) {
		t.Error("file2.txt symlink should have been removed")
	}

	// Verify file3 was created
	link3 := filepath.Join(targetDir, "file3.txt")
	if _, err := os.Lstat(link3); err != nil {
		t.Error("file3.txt symlink should have been created")
	}

	// Verify file3 points to correct location
	target, err := os.Readlink(link3)
	if err != nil {
		t.Fatalf("Failed to read file3 symlink: %v", err)
	}

	// Should be a relative path
	if filepath.IsAbs(target) {
		t.Errorf("Expected relative symlink, got absolute: %s", target)
	}
}
