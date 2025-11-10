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

# Run with debug logging
go run . /path/to/source /path/to/target --debug debug.log

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

**Using Makefile (recommended):**
```bash
# Run all checks (fmt, vet, test) - use before every commit
make check

# Format code with goimports and go fmt
make fmt

# Run go vet
make vet

# Run tests
make test

# Generate coverage report (creates coverage.html)
make coverage

# Build the project
make build

# Clean build artifacts
make clean

# Show all available targets
make help
```

**Manual commands:**
```bash
# Format code with goimports (auto-fixes unused imports)
goimports -w .

# Format code with go fmt
go fmt ./...

# Run linter (requires golangci-lint installation)
golangci-lint run

# Vet code for suspicious constructs
go vet ./...
```

**Note:** `goimports` is preferred over `go fmt` as it automatically manages imports (adds missing, removes unused). Install with:
```bash
go install golang.org/x/tools/cmd/goimports@latest
```

### Debugging
```bash
# Enable debug logging during development
go run . /path/to/source /path/to/target --debug debug.log

# View debug logs in real-time
tail -f debug.log

# Or watch logs while using the TUI
tail -f debug.log &
go run . /path/to/source /path/to/target --debug debug.log
```

Debug logging captures key events:
- Async file loading (available files and enabled files)
- User interactions (toggle, select all, deselect all)
- Mode changes (filter mode, hide mode)
- Selection state changes

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
├── Makefile                          # Build automation (check, fmt, test, build, etc.)
├── internal/
│   ├── config/
│   │   └── config.go                # Configuration management
│   ├── filesystem/
│   │   └── symlinks.go              # Symlink operations (create, remove, validate)
│   └── ui/
│       ├── tui.go                   # Terminal UI with bubbletea (multi-select, filter)
│       ├── types.go                 # Message types and list item implementation
│       ├── commands.go              # Async command functions
│       └── debug.go                 # Debug logging utility
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

**CLI Flags:**
- `--title`, `-t`: Title to display in UI
- `--version`, `-v`: Print version information
- `--debug`, `-d`: Enable debug logging to specified file (e.g., `--debug debug.log`)

## Dependencies

- **github.com/charmbracelet/bubbletea**: TUI framework for interactive interface
- **github.com/spf13/cobra**: CLI framework for command-line parsing

## Important Rules

- Run `make check` after each change
