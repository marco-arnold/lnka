package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// TestValidate tests the Validate method
func TestValidate(t *testing.T) {
	// Create temporary directories for testing
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	targetDir := filepath.Join(tempDir, "target")

	// Create the directories
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	tests := []struct {
		name      string
		config    Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid config",
			config: Config{
				SourceDir: sourceDir,
				TargetDir: targetDir,
				Title:     "Test Title",
			},
			wantError: false,
		},
		{
			name: "missing source directory",
			config: Config{
				SourceDir: "",
				TargetDir: targetDir,
			},
			wantError: true,
			errorMsg:  "source directory not specified",
		},
		{
			name: "missing target directory",
			config: Config{
				SourceDir: sourceDir,
				TargetDir: "",
			},
			wantError: true,
			errorMsg:  "target directory not specified",
		},
		{
			name: "non-existent source directory",
			config: Config{
				SourceDir: filepath.Join(tempDir, "nonexistent"),
				TargetDir: targetDir,
			},
			wantError: true,
			errorMsg:  "does not exist",
		},
		{
			name: "non-existent target directory",
			config: Config{
				SourceDir: sourceDir,
				TargetDir: filepath.Join(tempDir, "nonexistent"),
			},
			wantError: true,
			errorMsg:  "does not exist",
		},
		{
			name: "source is a file not directory",
			config: Config{
				SourceDir: filepath.Join(sourceDir, "file.txt"),
				TargetDir: targetDir,
			},
			wantError: true,
			errorMsg:  "is not a directory",
		},
	}

	// Create a file in source dir for file test
	testFile := filepath.Join(sourceDir, "file.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantError {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Validate() error = %v, want error containing %q", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

// TestLoad tests the Load function with cobra command
func TestLoad(t *testing.T) {
	// Create temporary directories for testing
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	targetDir := filepath.Join(tempDir, "target")

	// Create the directories
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	tests := []struct {
		name      string
		args      []string
		title     string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid config with title",
			args:      []string{sourceDir, targetDir},
			title:     "Test Title",
			wantError: false,
		},
		{
			name:      "valid config without title",
			args:      []string{sourceDir, targetDir},
			title:     "",
			wantError: false,
		},
		{
			name:      "missing target directory",
			args:      []string{sourceDir},
			title:     "",
			wantError: true,
			errorMsg:  "target directory not specified",
		},
		{
			name:      "missing both directories",
			args:      []string{},
			title:     "",
			wantError: true,
			errorMsg:  "source directory not specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new cobra command for each test
			cmd := &cobra.Command{
				Use: "test",
			}
			cmd.Flags().StringP("title", "t", "", "Title")
			if tt.title != "" {
				_ = cmd.Flags().Set("title", tt.title)
			}

			cfg, err := Load(cmd, tt.args)
			if tt.wantError {
				if err == nil {
					t.Errorf("Load() expected error but got none")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Load() error = %v, want error containing %q", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Load() unexpected error = %v", err)
				}
				if cfg == nil {
					t.Errorf("Load() returned nil config")
					return
				}
				if cfg.SourceDir != sourceDir {
					t.Errorf("Load() SourceDir = %v, want %v", cfg.SourceDir, sourceDir)
				}
				if cfg.TargetDir != targetDir {
					t.Errorf("Load() TargetDir = %v, want %v", cfg.TargetDir, targetDir)
				}
				if cfg.Title != tt.title {
					t.Errorf("Load() Title = %v, want %v", cfg.Title, tt.title)
				}
			}
		})
	}
}

// TestCheckDirExists tests the checkDirExists function indirectly through Validate
func TestCheckDirExists(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name      string
		setupFunc func() string
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid directory",
			setupFunc: func() string {
				dir := filepath.Join(tempDir, "valid")
				_ = os.MkdirAll(dir, 0755)
				return dir
			},
			wantError: false,
		},
		{
			name: "non-existent directory",
			setupFunc: func() string {
				return filepath.Join(tempDir, "nonexistent")
			},
			wantError: true,
			errorMsg:  "does not exist",
		},
		{
			name: "file instead of directory",
			setupFunc: func() string {
				file := filepath.Join(tempDir, "file.txt")
				_ = os.WriteFile(file, []byte("test"), 0644)
				return file
			},
			wantError: true,
			errorMsg:  "is not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupFunc()
			err := checkDirExists(path)
			if tt.wantError {
				if err == nil {
					t.Errorf("checkDirExists() expected error but got none")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("checkDirExists() error = %v, want error containing %q", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("checkDirExists() unexpected error = %v", err)
				}
			}
		})
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			len(s) > len(substr)+1 && containsInner(s, substr)))
}

func containsInner(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
