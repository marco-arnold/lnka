# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go project called "lnka" - a CLI tool for managing symlinks between source and target directories with an interactive Terminal UI.

**Key Features:**
- Interactive multi-select TUI using Bubble Tea
- Relative symlink creation when directories are close together
- Orphaned symlink detection and cleanup
- Filter mode for searching files
- Configuration via CLI flags and environment variables
- Built with Cobra for CLI framework
- Released via GoReleaser with GitHub Actions

## Development Commands

### Go Module Management
```bash
# Initialize Go module (if not already done)
go mod init github.com/marco-arnold/lnka

# Download dependencies
go mod download

# Tidy dependencies
go mod tidy
```

### Build and Run
```bash
# Build the project
go build ./...

# Run the main program
go run . /path/to/source /path/to/target

# Build with specific output
go build -o lnka
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -v -run TestFunctionName ./path/to/package

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run linter (requires golangci-lint installation)
golangci-lint run

# Vet code for suspicious constructs
go vet ./...
```

### Release Management
```bash
# Test GoReleaser configuration locally
goreleaser check
goreleaser release --snapshot --clean

# Create a new release (via git tag)
git tag v0.1.0
git push origin v0.1.0
# GitHub Actions will automatically build and release
```

## Repository Structure

```
lnka/
├── main.go                           # Entry point with cobra CLI & version info
├── internal/
│   ├── config/
│   │   └── config.go                # Configuration management
│   ├── filesystem/
│   │   └── symlinks.go              # Symlink operations (create, remove, validate)
│   └── ui/
│       └── tui.go                   # Terminal UI with bubbletea (multi-select, filter)
├── .github/
│   └── workflows/
│       └── release.yml              # GitHub Actions for automated releases
├── .goreleaser.yml                  # GoReleaser configuration
├── go.mod                           # Go module definition
├── LICENSE                          # MIT License
├── README.md                        # User documentation
├── CLAUDE.md                        # AI assistant guidance
└── PLAN.md                          # Implementation plan
```

## Notes

- The .gitignore is configured for Go projects with common exclusions (binaries, test artifacts, coverage profiles, vendor directories, GoReleaser artifacts)
- Go workspace files (go.work) are ignored
- Binary names `/enabler` and `/lnka` are explicitly ignored
- Test artifacts (`/test-data/`, `output.txt`) are ignored

## Configuration

**Environment Variables:**
- `LNKA_TITLE`: Optional title for the TUI
- `LNKA_MAX`: Maximum visible items before pagination (default: 10)

**CLI Flags:**
- `--title`, `-t`: Title to display in UI
- `--max`, `-m`: Maximum items to show before pagination
- `--version`, `-v`: Print version information

## Dependencies

- **github.com/charmbracelet/bubbletea**: TUI framework for interactive interface
- **github.com/spf13/cobra**: CLI framework for command-line parsing
