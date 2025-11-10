package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// TestPrintVersion tests the printVersion function
func TestPrintVersion(t *testing.T) {
	// Save original stdout
	oldStdout := os.Stdout

	tests := []struct {
		name    string
		version string
		commit  string
		date    string
		want    []string
	}{
		{
			name:    "default dev version",
			version: "dev",
			commit:  "none",
			date:    "unknown",
			want:    []string{"lnka dev", "commit: none", "built at: unknown"},
		},
		{
			name:    "release version",
			version: "1.0.0",
			commit:  "abc123def456",
			date:    "2024-01-15T10:30:00Z",
			want:    []string{"lnka 1.0.0", "commit: abc123def456", "built at: 2024-01-15T10:30:00Z"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Temporarily override version variables
			oldVersion := version
			oldCommit := commit
			oldDate := date

			version = tt.version
			commit = tt.commit
			date = tt.date

			// Restore after test
			defer func() {
				version = oldVersion
				commit = oldCommit
				date = oldDate
			}()

			// Capture stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Call the function
			printVersion()

			// Restore stdout and read output
			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Check if all expected strings are present
			for _, want := range tt.want {
				if !strings.Contains(output, want) {
					t.Errorf("printVersion() output missing %q, got:\n%s", want, output)
				}
			}
		})
	}
}

// TestExecute tests that the Execute function can be called
// Note: Full integration testing of cobra command is complex and would require mocking
func TestExecute_VersionFlag(t *testing.T) {
	// Save original args and stdout
	oldArgs := os.Args
	oldStdout := os.Stdout

	defer func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
	}()

	// Set version flag
	os.Args = []string{"lnka", "--version"}

	// Note: Testing Execute with --version is complex because it calls os.Exit(0)
	// This would terminate the test process. Full integration testing would require
	// refactoring Execute to be more testable (e.g., dependency injection, return
	// values instead of os.Exit).
	//
	// For now, we rely on:
	// - TestPrintVersion() for version output testing
	// - Manual integration testing for full command behavior
}

// TestVersionVariables tests that version variables are initialized
func TestVersionVariables(t *testing.T) {
	// These should have default values
	if version == "" {
		t.Error("version variable should not be empty")
	}
	if commit == "" {
		t.Error("commit variable should not be empty")
	}
	if date == "" {
		t.Error("date variable should not be empty")
	}

	// Check default values
	expectedDefaults := map[string]string{
		"version": "dev",
		"commit":  "none",
		"date":    "unknown",
	}

	// Note: In a real build, these would be overridden by ldflags
	// but in tests they should have the default values
	if version != expectedDefaults["version"] {
		t.Logf("version = %q (may be overridden by build)", version)
	}
	if commit != expectedDefaults["commit"] {
		t.Logf("commit = %q (may be overridden by build)", commit)
	}
	if date != expectedDefaults["date"] {
		t.Logf("date = %q (may be overridden by build)", date)
	}
}
