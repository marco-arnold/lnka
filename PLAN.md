# lnka - Implementation Plan

## Overview
Build a Go CLI tool with Terminal UI for managing symlinks between source and target directories.

## Architecture

### Project Structure
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
├── go.mod                           # Module: github.com/marco-arnold/lnka
├── go.sum                           # Dependencies lock
├── LICENSE                          # MIT License
├── README.md                        # User documentation
├── CLAUDE.md                        # AI assistant guidance
└── PLAN.md                          # This plan document
```

## Dependencies
- **github.com/charmbracelet/bubbletea** v1.3.10 - TUI framework
- **github.com/spf13/cobra** v1.10.1 - CLI framework

## Implementation Status

### ✅ Completed Features

#### 1. Project Initialization
- ✅ Go module initialized: `github.com/marco-arnold/lnka`
- ✅ Directory structure created
- ✅ Dependencies installed

#### 2. Configuration Layer (`internal/config/config.go`)
- ✅ `Config` struct with:
  - `SourceDir` string (first positional argument)
  - `TargetDir` string (second positional argument)
  - `Title` string (from `--title`/`-t` flag or `LNKA_TITLE` env)
  - `MaxVisibleItems` int (from `--max`/`-m` flag or `LNKA_MAX` env, default: 10)
- ✅ Validation: ensure both directories are provided and exist
- ✅ No defaults for directories - always require explicit arguments

#### 3. Filesystem Operations (`internal/filesystem/symlinks.go`)
- ✅ `ListAvailableFiles(dir string) ([]string, error)` - List files in source directory
- ✅ `ListEnabledSymlinks(dir, availableDir string) (map[string]string, error)` - Map of symlinks
- ✅ `GetEnabledFiles(enabledDir, availableDir string) ([]string, error)` - List enabled files
- ✅ `CreateSymlink(availableDir, enabledDir, filename string) error` - Create symlink with relative path support
- ✅ `RemoveSymlink(enabledDir, filename string) error` - Remove symlink
- ✅ `ValidateSymlinks(availableDir, enabledDir string) ([]string, error)` - Find orphaned symlinks
- ✅ `CleanOrphanedSymlinks(enabledDir string, orphaned []string) error` - Remove broken symlinks
- ✅ `ApplyChanges(availableDir, enabledDir string, selectedFiles []string) error` - Apply selection

#### 4. Terminal UI (`internal/ui/tui.go`)
Using **pure bubbletea** (no huh):
- ✅ Multi-select interface with all files from source directory
- ✅ Pre-select files that are already enabled (symlinks exist)
- ✅ Toggle selection with space, navigate with arrow keys
- ✅ Filter mode: press `/` to filter list
- ✅ Enter confirms selection, ESC aborts
- ✅ Cursor: filled triangle `▶`
- ✅ Visual feedback: dim gray for unselected, normal for selected
- ✅ Pagination with configurable max visible items
- ✅ ANSI color codes for styling
- ✅ Confirmation dialog for orphaned symlink cleanup

#### 5. Main Entry Point (`main.go`)
- ✅ Cobra CLI framework with positional arguments
- ✅ Flags: `--title`/`-t`, `--max`/`-m`, `--version`/`-v`
- ✅ Environment variables: `LNKA_TITLE`, `LNKA_MAX`
- ✅ Load configuration
- ✅ Validate symlinks and prompt for cleanup
- ✅ Launch TUI
- ✅ Apply symlink changes based on user selection
- ✅ Exit silently on success, exit code 1 on abort
- ✅ Version information embedded via ldflags

#### 6. Release Management
- ✅ `.goreleaser.yml` configuration
  - Static binaries for Linux, macOS, Windows (amd64, arm64)
  - Archives with README and LICENSE
  - SHA256 checksums
  - Automatic changelog generation
- ✅ GitHub Actions workflow (`.github/workflows/release.yml`)
  - Triggers on `v*` tags
  - Automated release creation
- ✅ MIT License file
- ✅ `.gitignore` updated for build artifacts and goreleaser

## CLI Interface

### Usage
```bash
# Basic usage with positional arguments
lnka /path/to/source /path/to/target

# With options
lnka /path/to/source /path/to/target --title "My Files" --max 15

# Using shorthand flags
lnka /path/to/source /path/to/target -t "My Files" -m 15

# Using environment variables
export LNKA_TITLE="My Files"
export LNKA_MAX=15
lnka /path/to/source /path/to/target

# Show version
lnka --version
```

## User Flow
1. Start tool with source and target directory arguments
2. If orphaned symlinks detected → prompt user to clean them (yes/no confirmation)
3. Display multi-select TUI with all available files from source directory
4. Files already enabled (existing symlinks in target) are pre-selected
5. User can:
   - Navigate with arrow keys
   - Toggle selection with space
   - Filter list by pressing `/`
   - Confirm with Enter
   - Abort with ESC
6. Tool calculates diff:
   - Create symlinks in target for newly selected files
   - Remove symlinks from target for deselected files
   - Use relative paths when directories are close together
7. Exit silently on success
8. Exit with code 1 on abort (no error message)

## Key Features
- ✅ Multi-select TUI for file selection
- ✅ Configuration via positional arguments, CLI flags, and ENV vars
- ✅ Automatic detection and cleanup of orphaned symlinks
- ✅ Validation of existing symlinks
- ✅ Pre-selection of currently enabled files
- ✅ Filter mode for searching files
- ✅ Relative symlink creation when beneficial
- ✅ Clear error handling and user feedback
- ✅ Pagination for long file lists
- ✅ Version information embedded in binary
- ✅ Automated releases via GoReleaser and GitHub Actions

## Release Process
```bash
# Create and push a tag
git tag v0.1.0
git push origin v0.1.0

# GitHub Actions automatically:
# - Builds binaries for all platforms
# - Creates GitHub release
# - Generates changelog
# - Uploads artifacts
```

## Testing Strategy
- Unit tests for filesystem operations (use temp directories)
- Integration tests for config loading
- Manual testing for TUI interactions

## Future Enhancements
- Dry-run mode with `--dry-run` flag
- Verbose/debug logging with `--verbose`
- List-only mode (no UI, just show status)
- Homebrew tap integration
- Support for nested directories (currently only top-level files)
- Configuration file support (e.g., `.lnka.yml`)
