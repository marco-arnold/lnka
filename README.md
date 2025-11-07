# lnka

A Go CLI tool with Terminal UI for managing symlinks between source and target directories.

<div align="center">
  <img src="assets/logo.png" alt="lnka logo" width="200"/>
</div>

## Overview

lnka simplifies the process of managing symlinks between directories. Similar to how nginx and systemd handle configurations, this tool allows you to:

- View all available files from the source directory in a multi-select Terminal UI
- Enable files by creating symlinks in the target directory pointing to the source directory
- Disable files by removing symlinks from the target directory
- Automatically detect and clean orphaned (broken) symlinks
- Create relative symlinks when directories are close together

## Installation

### Via Homebrew (macOS/Linux)

```bash
# Install directly
brew install marco-arnold/lnka/lnka

# Or tap first, then install
brew tap marco-arnold/lnka
brew install lnka
```

### Via Go Install

```bash
go install github.com/marco-arnold/lnka@latest
```

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/marco-arnold/lnka/releases).

### Build from Source

```bash
git clone https://github.com/marco-arnold/lnka.git
cd lnka
go build -o lnka
```

## Usage

### Basic Usage

```bash
lnka /path/to/source /path/to/target
```

The first argument is the source directory containing files you want to manage, the second is the target directory where symlinks will be created.

### With Options

```bash
lnka /path/to/source /path/to/target --title "My Services" --max 15
```

Or using shorthand flags:

```bash
lnka /path/to/source /path/to/target -t "My Services" -m 15
```

### Using Environment Variables

```bash
export LNKA_TITLE="My Services"
export LNKA_MAX=15
lnka /path/to/source /path/to/target
```

## Features

- **Interactive Terminal UI**: Multi-select interface using charmbracelet/bubbletea
- **Pre-selection**: Currently enabled files (existing symlinks) are automatically selected
- **Orphaned Symlink Detection**: Automatically finds and optionally removes broken symlinks
- **Relative Symlinks**: Creates relative symlinks when source and target are close together
- **Filter Mode**: Press `/` to filter the list of available files
- **Flexible Configuration**: Configure via positional arguments, CLI flags and environment variables

## User Flow

1. Run the tool with source and target directories as arguments
2. If orphaned symlinks are detected, you'll be prompted to clean them
3. The multi-select UI displays all available files from the source directory
4. Files that already have symlinks in the target directory are pre-selected
5. Use arrow keys to navigate, space to toggle, `/` to filter, enter to confirm
6. The tool creates/removes symlinks in the target directory based on your selection
7. Exit silently on success or with code 1 on ESC/abort

## Configuration

### Positional Arguments

- First argument: Path to source directory (required)
- Second argument: Path to target directory (required)

### CLI Flags

- `--title`, `-t`: Title to display in UI (default: empty)
- `--max`, `-m`: Maximum items to show before pagination (default: 10)

### Environment Variables

- `LNKA_TITLE`: Title to display in UI
- `LNKA_MAX`: Maximum items to show before pagination

**Note**: No default directories are provided. You must always specify both directories explicitly.

## Project Structure

```
lnka/
├── main.go                           # Entry point
├── internal/
│   ├── config/
│   │   └── config.go                # Configuration management
│   ├── filesystem/
│   │   └── symlinks.go              # Symlink operations
│   └── ui/
│       └── tui.go                   # Terminal UI
└── PLAN.md                          # Implementation plan
```

## Dependencies

- [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [github.com/spf13/cobra](https://github.com/spf13/cobra) - CLI framework

## Example

```bash
# Setup test directories
mkdir -p source target

# Create some test files in the source directory
touch source/nginx.conf
touch source/redis.conf
touch source/postgres.conf

# Run lnka
lnka source target
```

The interactive UI will appear, allowing you to select which files to enable. Selected files will have symlinks created in the `target` directory pointing to the corresponding files in the `source` directory.

## License

MIT
