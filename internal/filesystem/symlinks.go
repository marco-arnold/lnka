package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

// ListAvailableFiles lists all files (not directories) in the source directory
func ListAvailableFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read source directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		// Only include regular files, skip directories
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// ListEnabledSymlinks returns a map of symlink names to their targets
// Only includes symlinks that point to files in the source directory
func ListEnabledSymlinks(sourceDir string, targetDir string) (map[string]string, error) {
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read target directory: %w", err)
	}

	symlinks := make(map[string]string)
	for _, entry := range entries {
		// Check if it's a symlink
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.Mode()&os.ModeSymlink != 0 {
			linkPath := filepath.Join(targetDir, entry.Name())
			target, err := os.Readlink(linkPath)
			if err != nil {
				continue
			}

			// Store the symlink name and its target
			symlinks[entry.Name()] = target
		}
	}

	return symlinks, nil
}

// GetEnabledFiles returns a list of file names that are currently enabled
// (have symlinks pointing to them in the target directory)
func GetEnabledFiles(sourceDir string, targetDir string) ([]string, error) {
	symlinks, err := ListEnabledSymlinks(sourceDir, targetDir)
	if err != nil {
		return nil, err
	}

	enabled := make([]string, 0, len(symlinks))
	for name, target := range symlinks {
		// Resolve the target path (could be relative or absolute)
		var resolvedTarget string
		if filepath.IsAbs(target) {
			resolvedTarget = target
		} else {
			resolvedTarget = filepath.Join(targetDir, target)
		}

		// Check if the resolved target points to a file in sourceDir
		expectedPath := filepath.Join(sourceDir, name)

		// Compare resolved paths
		resolvedTargetAbs, err1 := filepath.Abs(resolvedTarget)
		expectedPathAbs, err2 := filepath.Abs(expectedPath)

		if err1 == nil && err2 == nil && resolvedTargetAbs == expectedPathAbs {
			enabled = append(enabled, name)
		}
	}

	return enabled, nil
}

// CreateSymlink creates a symlink in the target directory pointing to a file in the source directory
// Uses relative paths when source and target are close together
func CreateSymlink(sourceDir, targetDir, filename string) error {
	sourcePath := filepath.Join(sourceDir, filename)
	linkPath := filepath.Join(targetDir, filename)

	// Check if source file exists
	if _, err := os.Stat(sourcePath); err != nil {
		return fmt.Errorf("source file %s does not exist: %w", filename, err)
	}

	// Check if symlink already exists
	if _, err := os.Lstat(linkPath); err == nil {
		// Symlink exists, remove it first
		if err := os.Remove(linkPath); err != nil {
			return fmt.Errorf("failed to remove existing symlink %s: %w", filename, err)
		}
	}

	// Try to create a relative symlink if possible
	symlinkTarget := sourcePath
	relPath, err := filepath.Rel(targetDir, sourcePath)
	if err == nil && !filepath.IsAbs(relPath) && len(relPath) < len(sourcePath) {
		// Use relative path if it's shorter and valid
		symlinkTarget = relPath
	}

	// Create the symlink
	if err := os.Symlink(symlinkTarget, linkPath); err != nil {
		return fmt.Errorf("failed to create symlink %s: %w", filename, err)
	}

	return nil
}

// RemoveSymlink removes a symlink from the target directory
func RemoveSymlink(targetDir, filename string) error {
	linkPath := filepath.Join(targetDir, filename)

	// Check if symlink exists
	info, err := os.Lstat(linkPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Symlink doesn't exist, nothing to do
			return nil
		}
		return fmt.Errorf("failed to check symlink %s: %w", filename, err)
	}

	// Verify it's a symlink before removing
	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("%s is not a symlink, refusing to remove", filename)
	}

	// Remove the symlink
	if err := os.Remove(linkPath); err != nil {
		return fmt.Errorf("failed to remove symlink %s: %w", filename, err)
	}

	return nil
}

// ValidateSymlinks finds orphaned or broken symlinks in the target directory
// Returns a list of symlink names that are broken (point to non-existent files)
func ValidateSymlinks(sourceDir, targetDir string) ([]string, error) {
	symlinks, err := ListEnabledSymlinks(sourceDir, targetDir)
	if err != nil {
		return nil, err
	}

	var orphaned []string
	for name, target := range symlinks {
		// Resolve target path relative to target directory if it's a relative path
		targetPath := target
		if !filepath.IsAbs(target) {
			targetPath = filepath.Join(targetDir, target)
		}

		// Check if target exists
		if _, err := os.Stat(targetPath); err != nil {
			if os.IsNotExist(err) {
				orphaned = append(orphaned, name)
			}
		}
	}

	return orphaned, nil
}

// CleanOrphanedSymlinks removes broken symlinks from the target directory
func CleanOrphanedSymlinks(targetDir string, orphaned []string) error {
	for _, name := range orphaned {
		if err := RemoveSymlink(targetDir, name); err != nil {
			return fmt.Errorf("failed to clean orphaned symlink %s: %w", name, err)
		}
	}

	return nil
}

// ApplyChanges applies the user's selection by creating and removing symlinks
func ApplyChanges(sourceDir, targetDir string, selectedFiles []string) error {
	// Get currently enabled files
	currentlyEnabled, err := GetEnabledFiles(sourceDir, targetDir)
	if err != nil {
		return fmt.Errorf("failed to get currently enabled files: %w", err)
	}

	// Convert to maps for easier lookup
	selectedMap := make(map[string]bool)
	for _, name := range selectedFiles {
		selectedMap[name] = true
	}

	currentMap := make(map[string]bool)
	for _, name := range currentlyEnabled {
		currentMap[name] = true
	}

	// Remove symlinks for files that are no longer selected
	for _, name := range currentlyEnabled {
		if !selectedMap[name] {
			if err := RemoveSymlink(targetDir, name); err != nil {
				return err
			}
		}
	}

	// Create symlinks for newly selected files
	for _, name := range selectedFiles {
		if !currentMap[name] {
			if err := CreateSymlink(sourceDir, targetDir, name); err != nil {
				return err
			}
		}
	}

	return nil
}
