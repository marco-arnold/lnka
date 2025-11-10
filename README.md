# lnka

**A fast, interactive CLI tool for managing symlinks with a beautiful Terminal UI.**

Manage your configuration files, systemd services, or nginx sites with ease. lnka provides an intuitive multi-select interface to enable/disable files by creating or removing symlinks.

<div align="center">
  <img src="assets/logo.png" alt="lnka logo" width="200"/>
</div>


## What is lnka?

lnka (pronounced "link-a") simplifies symlink management between directories. Think of it as a universal tool for the pattern used by nginx (`sites-available` → `sites-enabled`) and systemd (`/lib/systemd/system` → `/etc/systemd/system`).

**Perfect for:**
- 🔧 Managing nginx/apache site configurations
- ⚙️  Enabling/disabling systemd services
- 📝 Organizing dotfiles and config files
- 🔗 Any workflow involving symlinks between directories

**Key Features:**
- ✨ Beautiful interactive Terminal UI powered by Bubble Tea
- ⚡ Fast async file loading
- 🔍 Built-in fuzzy search/filter
- 🎯 Pre-selects currently enabled files
- 🧹 Automatic detection and cleanup of broken symlinks
- 📏 Smart relative symlinks when possible
- ⌨️  Vim-style keyboard navigation

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

## Quick Start

### Basic Usage

```bash
lnka <source-dir> <target-dir>
```

**Example: Managing nginx sites**
```bash
lnka /etc/nginx/sites-available /etc/nginx/sites-enabled
```

**Example: Managing systemd services**
```bash
lnka /lib/systemd/system /etc/systemd/system --title "System Services"
```

**Example: Dotfiles**
```bash
lnka ~/dotfiles ~/.config
```

### Command-Line Options

```bash
# With custom title
lnka /path/to/source /path/to/target --title "My Configs"
lnka /path/to/source /path/to/target -t "My Configs"

# With debug logging
lnka /path/to/source /path/to/target --debug debug.log

# Show version
lnka --version
```

### Environment Variables

```bash
export LNKA_TITLE="My Services"
lnka /path/to/source /path/to/target
```

## How It Works

1. **Launch** - Run lnka with your source and target directories
2. **Auto-detect** - Broken symlinks? You'll be prompted to clean them
3. **Select** - Interactive UI shows all files, with currently enabled files pre-selected
4. **Navigate** - Use keyboard shortcuts to browse, filter, and select files
5. **Apply** - Press Enter to create/remove symlinks based on your selection
6. **Done** - Exit silently on success

## Keyboard Shortcuts

### Essential Shortcuts
| Key | Action |
|-----|--------|
| `Space` | Select/deselect item at cursor |
| `Enter` | Confirm selection and apply changes |
| `↑/k` `↓/j` | Navigate up/down (Vim-style) |
| `/` | Enter filter mode (fuzzy search) |
| `h` | Toggle hide mode (show only linked items) |
| `?` | Toggle help (short/full) |
| `Ctrl+C` | Abort without changes |

### Advanced Shortcuts (shown in full help with `?`)
| Key | Action |
|-----|--------|
| `g` / `G` | Jump to top/bottom |
| `PgUp/PgDn` | Page up/down |
| `Ctrl+B` / `Ctrl+F` | Page up/down (Vim-style) |
| `Ctrl+A` | Select all visible items |
| `Ctrl+D` | Deselect all items |

### Filter Mode
| Key | Action |
|-----|--------|
| `Type...` | Filter list (fuzzy search) |
| `Backspace` | Remove filter characters |
| `Enter` | Exit filter mode |
| `Esc` | Clear filter and exit filter mode |

## Configuration

### Required Arguments

```bash
lnka <source-dir> <target-dir>
```

Both directories must exist. The tool will:
- Read available files from `<source-dir>`
- Create/remove symlinks in `<target-dir>`

### Optional Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--title` | `-t` | Title displayed in UI | (empty) |
| `--debug` | `-d` | Enable debug logging to file | (disabled) |
| `--version` | `-v` | Show version information | - |

### Environment Variables

| Variable | Description |
|----------|-------------|
| `LNKA_TITLE` | Default title for UI |

## Real-World Examples

### nginx Site Management

Enable/disable nginx virtual hosts:

```bash
# Interactive selection
sudo lnka /etc/nginx/sites-available /etc/nginx/sites-enabled -t "nginx Sites"

# Reload nginx after changes
sudo nginx -s reload
```

### systemd Service Management

Enable/disable systemd services:

```bash
# System services
sudo lnka /lib/systemd/system /etc/systemd/system -t "Services"

# User services
lnka ~/.local/share/systemd/user ~/.config/systemd/user -t "User Services"

# Reload systemd after changes
sudo systemctl daemon-reload
```

### Dotfiles Management

Organize your dotfiles:

```bash
# Create directory structure
mkdir -p ~/dotfiles/{bash,vim,git}
mv ~/.bashrc ~/dotfiles/bash/
mv ~/.vimrc ~/dotfiles/vim/
mv ~/.gitconfig ~/dotfiles/git/

# Manage symlinks
lnka ~/dotfiles/bash ~/
lnka ~/dotfiles/vim ~/
lnka ~/dotfiles/git ~/
```

### Apache Site Management

Similar to nginx:

```bash
sudo lnka /etc/apache2/sites-available /etc/apache2/sites-enabled -t "Apache Sites"
sudo systemctl reload apache2
```

## Troubleshooting

### Permission Denied

If you get permission errors when creating symlinks in system directories:

```bash
# Run with sudo
sudo lnka /etc/nginx/sites-available /etc/nginx/sites-enabled
```

### Broken Symlinks

lnka automatically detects broken symlinks and offers to clean them:

```bash
$ lnka source target
Found 2 orphaned symlinks: old-site.conf, deprecated.conf
Clean orphaned symlinks? (Y/n): y
✓ Cleaned 2 orphaned symlinks
```

### Debug Mode

Enable debug logging to troubleshoot issues:

```bash
lnka source target --debug debug.log
tail -f debug.log  # View logs in real-time
```

## Technical Details

### Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Modern TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Cobra](https://github.com/spf13/cobra) - CLI framework

### Features Under the Hood

- **Async Loading**: Files load asynchronously for instant startup
- **Smart Symlinks**: Creates relative paths when beneficial, absolute when necessary
- **Safe Operations**: Refuses to delete regular files, only removes symlinks
- **Idempotent**: Safe to run multiple times, won't duplicate or break existing setups
- **Cross-Platform**: Works on Linux, macOS, and Windows (with symlink support)

## Contributing

Contributions are welcome! Here's how you can help:

- 🐛 Report bugs via [GitHub Issues](https://github.com/marco-arnold/lnka/issues)
- 💡 Suggest features or improvements
- 📖 Improve documentation
- 🧪 Write tests
- 🔧 Submit pull requests

For development setup and guidelines, see [CLAUDE.md](CLAUDE.md).

## FAQ

**Q: Can I use this with version-controlled config files?**
A: Absolutely! lnka works great with git-managed dotfiles or configs.

**Q: What happens if I abort (Ctrl+C)?**
A: No changes are made. Your symlinks remain exactly as they were.

**Q: Can I manage subdirectories?**
A: Currently, lnka only manages files in the top level of the source directory.

**Q: Does it work on Windows?**
A: Yes, but you need Windows 10+ with Developer Mode enabled for symlink support.

**Q: How do I update to the latest version?**
A: If installed via Homebrew: `brew upgrade lnka`. Via Go: `go install github.com/marco-arnold/lnka@latest`

## License

MIT License - see [LICENSE](LICENSE) for details

---

**Made with ❤️  using [Bubble Tea](https://github.com/charmbracelet/bubbletea)**
