# lnka - Implementation Plan

## Overview
Build a Go CLI tool with Terminal UI for managing symlinks between source and target directories.

## Architecture

### Project Structure
```
lnka/
├── main.go                           # Entry point with cobra CLI & version info
├── Makefile                          # Build automation (check, fmt, vet, test, coverage, build)
├── internal/
│   ├── config/
│   │   └── config.go                # Configuration management
│   ├── filesystem/
│   │   └── symlinks.go              # Symlink operations (create, remove, validate)
│   └── ui/
│       ├── tui.go                   # Terminal UI main logic (models, update, view)
│       ├── types.go                 # Message types and list item implementation
│       ├── commands.go              # Async command functions (file loading)
│       └── debug.go                 # Debug logging utility
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
- **github.com/charmbracelet/bubbles** v0.21.0 - Reusable Bubble Tea components
- **github.com/charmbracelet/lipgloss** v1.1.0 - Terminal styling
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
  - `DebugLog` string (from `--debug`/`-d` flag for debug logging to file)
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

#### 4. Terminal UI (`internal/ui/`)
Refactored architecture using **Bubble Tea with bubbles/list component**:

**Core Components:**
- ✅ `tui.go` - Main UI logic:
  - `multiSelectModel` with bubbles/list.Model integration
  - Keyboard navigation with custom keyMap
  - Async file loading via Init()
  - Smart selection management with selectedMap and selectedOrder
  - Auto-disable hideUnlinked mode when no items selected
  - Help system integration (short/full help with ?)
- ✅ `types.go` - Data structures:
  - `filesLoadedMsg` for async file loading
  - `itemsRefreshedMsg` for item list rebuilds
  - `fileItem` implementing list.Item interface
  - `fileItemDelegate` for custom rendering
- ✅ `commands.go` - Async commands:
  - `loadFilesCmd()` for parallel file loading
- ✅ `debug.go` - Debug logging:
  - `logDebug()` function for development tracing

**Features:**
- ✅ Multi-select interface with all files from source directory
- ✅ Pre-select files that are already enabled (symlinks exist)
- ✅ Toggle selection with space, navigate with arrow keys
- ✅ Vim-style navigation: j/k for up/down, g/G for top/bottom
- ✅ Page navigation: PgUp/PgDn or ctrl+b/ctrl+f
- ✅ Bulk operations: ctrl+a (select all), ctrl+d (deselect all)
- ✅ Filter mode: press `/` to filter list (case-insensitive)
- ✅ Hide mode: press `h` to toggle between all/linked items
- ✅ Enter confirms selection, ctrl+c aborts
- ✅ Visual feedback:
  - Bold green cursor marker (`>`) for selected item at cursor
  - Green cursor marker (not bold) for unselected item at cursor
  - Bold text for linked items
  - Gray text for unlinked items
- ✅ Built-in help system (press `?` to toggle short/full help)
- ✅ Async file loading with loading state
- ✅ Smart item list rebuilding after mode changes
- ✅ Confirmation dialog for orphaned symlink cleanup
- ✅ Comprehensive package documentation with examples
- ✅ Inline documentation for all functions and methods

#### 5. Main Entry Point (`main.go`)
- ✅ Cobra CLI framework with positional arguments
- ✅ Flags: `--title`/`-t`, `--debug`/`-d`, `--version`/`-v`
- ✅ Environment variables: `LNKA_TITLE`
- ✅ Debug logging to file via `--debug` flag (e.g., `--debug debug.log`)
- ✅ Load configuration
- ✅ Validate symlinks and prompt for cleanup
- ✅ Launch TUI with async file loading
- ✅ Apply symlink changes based on user selection
- ✅ Exit silently on success, exit code 1 on abort
- ✅ Version information embedded via ldflags

#### 6. Build Management (`Makefile`)
- ✅ Automated build tasks:
  - `make check` - Run all checks (fmt, vet, test) - use before commits
  - `make fmt` - Format code with goimports and go fmt
  - `make vet` - Run go vet for suspicious constructs
  - `make test` - Run all tests
  - `make coverage` - Generate HTML coverage report
  - `make build` - Build the binary
  - `make clean` - Clean build artifacts
  - `make help` - Show all available targets
- ✅ Uses goimports for automatic import management
- ✅ Integrated into development workflow

#### 7. Release Management
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
lnka /path/to/source /path/to/target --title "My Files"

# Using shorthand flags
lnka /path/to/source /path/to/target -t "My Files"

# With debug logging
lnka /path/to/source /path/to/target --debug debug.log

# Using environment variables
export LNKA_TITLE="My Files"
lnka /path/to/source /path/to/target

# Show version
lnka --version
```

## User Flow
1. Start tool with source and target directory arguments (optional: debug logging)
2. Files are loaded asynchronously from source and target directories
3. If orphaned symlinks detected → prompt user to clean them (yes/no confirmation)
4. Display multi-select TUI with all available files from source directory
5. Files already enabled (existing symlinks in target) are pre-selected
6. User can:
   - Navigate with arrow keys or j/k (Vim-style)
   - Jump to top/bottom with g/G
   - Page up/down with PgUp/PgDn or ctrl+b/ctrl+f
   - Toggle selection with space
   - Select all visible items with ctrl+a
   - Deselect all items with ctrl+d
   - Filter list by pressing `/` (case-insensitive search built into bubbles/list)
   - Hide unlinked items by pressing `h` (only shows selected items)
   - Show help by pressing `?` (toggle short/full help)
   - Confirm with Enter
   - Abort with ctrl+c
7. Tool calculates diff:
   - Create symlinks in target for newly selected files
   - Remove symlinks from target for deselected files
   - Use relative paths when directories are close together
8. Exit silently on success
9. Exit with code 1 on abort (no error message)

## Key Features
- ✅ Multi-select TUI for file selection using bubbles/list component
- ✅ Configuration via positional arguments, CLI flags, and ENV vars
- ✅ Asynchronous file loading for responsive startup
- ✅ Debug logging to file for development and troubleshooting
- ✅ Automatic detection and cleanup of orphaned symlinks
- ✅ Validation of existing symlinks
- ✅ Pre-selection of currently enabled files
- ✅ Built-in filter mode with fuzzy search (bubbles/list)
- ✅ Hide mode for showing only linked items
- ✅ Integrated help system (short/full help with `?`)
- ✅ Vim-style keyboard shortcuts (j/k, g/G)
- ✅ Bulk selection operations (ctrl+a, ctrl+d)
- ✅ Relative symlink creation when beneficial
- ✅ Clear error handling and user feedback
- ✅ Custom rendering with lipgloss styling
- ✅ Version information embedded in binary
- ✅ Automated releases via GoReleaser and GitHub Actions
- ✅ Makefile for development workflow automation

## TUI Architecture Refactoring (Bubble Tea Migration)

The Terminal UI has been refactored to use the official Bubble Tea component library:

### Architecture Changes
- ✅ **Migrated to bubbles/list**: Replaced custom list implementation with `github.com/charmbracelet/bubbles/list`
  - Built-in filtering with fuzzy search
  - Built-in pagination and scrolling
  - Built-in help system integration
  - Standard keyboard navigation (handled by component)
- ✅ **Modular code organization**: Split into multiple files for better maintainability
  - `tui.go` - Main UI logic, models, Update/View functions
  - `types.go` - Message types and list item implementation
  - `commands.go` - Async command functions
  - `debug.go` - Debug logging utilities
- ✅ **Custom rendering**: Implemented `fileItemDelegate` for styled item rendering
  - Bold green cursor marker for selected items
  - Gray styling for unlinked items
  - Bold styling for linked items
  - Lipgloss integration for consistent styling

### Performance & Behavior
- ✅ **Async file loading**: Files loaded in Init() via tea.Cmd
  - Non-blocking startup
  - Loading state displayed to user
  - Error handling for file loading failures
- ✅ **O(1) selection/deselection**: Using `selectedMap map[string]bool`
  - Instant lookup for item selection state
  - Separate `selectedOrder` slice maintains selection order
- ✅ **Smart item list rebuilding**: Efficient updates after mode changes
  - `buildItemList()` respects hideUnlinked mode
  - `rebuildItemsCmd()` triggers list refresh via message
  - `refreshCurrentItem()` updates single item efficiently

### User Experience Enhancements
- ✅ **Built-in help system**: Press `?` to toggle short/full help
  - Custom keybindings integrated into help display
  - Short help shows essential shortcuts
  - Full help includes all available commands (including ctrl+c abort)
- ✅ **Auto-disable hideUnlinked**: Automatically disables when no items selected
  - Prevents empty list confusion
  - Checked after deselection operations
- ✅ **Smart cursor management**: Cursor position maintained by bubbles/list
- ✅ **Filter integration**: Built-in filter mode with `/` key
  - Fuzzy search through bubbles/list
  - Visual feedback during filtering
  - Enter to exit filter mode

### Documentation
- ✅ **Comprehensive package documentation**: Complete overview with usage examples
- ✅ **Function documentation**: Detailed docs for all public functions
- ✅ **Inline documentation**: Comments for models, methods, and key logic
- ✅ **Architecture notes**: Explains async loading, message types, and component integration

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
