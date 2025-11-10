// Package ui provides terminal user interface components using Bubble Tea.
//
// The package includes:
//   - Multi-select list with filtering and hiding capabilities
//   - Confirmation dialogs for yes/no prompts
//   - Keyboard navigation and visual feedback
//
// # Key Features
//
//   - Filter mode: Press '/' to search through items
//   - Hide mode: Press 'h' to toggle between all/linked items
//   - Automatic sizing: List adapts to terminal size
//   - Smart cursor positioning: Maintains cursor position across mode switches
//   - Vim-style navigation: j/k for up/down, g/G for top/bottom
//   - Bulk operations: ctrl+a to select all, ctrl+d to deselect all
//   - Visual feedback: Bold for linked items, gray for unlinked, bold green for cursor
//
// # Multi-Select UI
//
// The multi-select interface allows users to select multiple items from a list
// with keyboard navigation. Features include:
//
//   - Space: Select/deselect item at cursor
//   - j/k or ↑/↓: Navigate items
//   - g/G: Jump to top/bottom
//   - PgUp/PgDn or ctrl+b/ctrl+f: Page up/down
//   - ctrl+a: Select all visible items
//   - ctrl+d: Deselect all items
//   - /: Enter filter mode to search
//   - h: Toggle between showing all items or only linked items
//   - Enter: Confirm selection
//   - ?: Toggle help (ctrl+c to abort in extended help)
//   - ctrl+c: Abort (shown in extended help with ?)
//
// Example usage:
//
//	sourceDir := "/path/to/source"
//	targetDir := "/path/to/target"
//	selected, err := ui.ShowMultiSelect(sourceDir, targetDir, "Select files")
//	if err != nil {
//	    // Handle error (user aborted or other error)
//	}
//	// Use selected files
//
// # Confirmation Dialog
//
// The confirmation dialog shows a simple yes/no prompt:
//
//	confirmed, err := ui.ShowConfirmation("Delete all files?")
//	if err != nil {
//	    // Handle error (user aborted)
//	}
//	if confirmed {
//	    // Perform action
//	}
//
// # Performance Considerations
//
// The UI is optimized for large lists (1000+ items) with:
//   - O(1) selection/deselection using indexed maps
//   - Per-cycle caching of visible choices
//   - Early exit optimization in hideUnlinked mode
//   - Efficient pagination with smart viewport management
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UI layout constants
const (
	// helpBarReservedLines is the number of lines reserved for the help bar
	// and optional list chrome (title when set, help bar)
	helpBarReservedLines = 4
)

// lipgloss styles for terminal UI
var (
	stylePrompt = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")) // Bold Green

	// Help bar style for confirmation dialog - inverse video spanning full width
	styleHelpBar = lipgloss.NewStyle().
			Reverse(true).
			Width(0) // Will be set dynamically based on terminal width
)

// keyMap defines all keyboard shortcuts for the multi-select UI
type keyMap struct {
	Quit        key.Binding // Abort operation (ctrl+c) - shown only in full help
	Confirm     key.Binding // Confirm selection (enter)
	Filter      key.Binding // Enter filter mode (/)
	HideToggle  key.Binding // Toggle hide unlinked items (h)
	Select      key.Binding // Select/deselect item at cursor (space)
	Up          key.Binding // Move cursor up (↑/k)
	Down        key.Binding // Move cursor down (↓/j)
	GoTop       key.Binding // Jump to top (g)
	GoBottom    key.Binding // Jump to bottom (G)
	SelectAll   key.Binding // Select all visible items (ctrl+a)
	DeselectAll key.Binding // Deselect all items (ctrl+d)
	PageDown    key.Binding // Page down (pgdn/ctrl+f)
	PageUp      key.Binding // Page up (pgup/ctrl+b)
}

// defaultKeyMap returns the default keyboard shortcuts for the multi-select UI.
// All shortcuts are defined with their keys and help text for the built-in help system.
// The keyMap is used in normal mode (not during filtering).
func defaultKeyMap() *keyMap {
	return &keyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "abort"),
		),
		Confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		HideToggle: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "toggle"),
		),
		Select: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "select"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		GoTop: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "top"),
		),
		GoBottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "bottom"),
		),
		SelectAll: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "select all"),
		),
		DeselectAll: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "deselect all"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "ctrl+f"),
			key.WithHelp("pgdn/ctrl+f", "page down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "ctrl+b"),
			key.WithHelp("pgup/ctrl+b", "page up"),
		),
	}
}

// multiSelectModel is the Bubble Tea model for multi-select UI
// It manages the state for selecting multiple items from a list
type multiSelectModel struct {
	list           list.Model      // Bubble Tea list component (replaces: choices, cursor, filter, filtered)
	selectedMap    map[string]bool // Selected items (renamed from 'selected' for clarity)
	selectedOrder  []string        // Order of selection for result (preserved for consistent output)
	sourceDir      string          // Source directory for Commands
	targetDir      string          // Target directory for Commands
	availableFiles []string        // Unfiltered source list (for rebuilding items after mode changes)
	aborted        bool            // User pressed ctrl+c
	hideUnlinked   bool            // Hide unlinked items when true
	loading        bool            // Files are being loaded
	err            error           // Error during loading
	keys           *keyMap         // Keyboard shortcuts (now a pointer following Go conventions)
}

// Init initializes the model
// Returns command to load available and enabled files asynchronously
func (m multiSelectModel) Init() tea.Cmd {
	logDebug("Init: starting async load from sourceDir=%s, targetDir=%s", m.sourceDir, m.targetDir)
	return loadFilesCmd(m.sourceDir, m.targetDir)
}

// Update handles messages
func (m multiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Handle async file loading message
	case filesLoadedMsg:
		if msg.err != nil {
			logDebug("filesLoadedMsg: error loading files: %v", msg.err)
			m.err = msg.err
			m.aborted = true
			return m, tea.Quit
		}

		logDebug("filesLoadedMsg: loaded %d available files, %d enabled files",
			len(msg.availableFiles), len(msg.enabledFiles))

		// Store available files
		m.availableFiles = msg.availableFiles

		// Build initial selection map from enabled files
		for _, file := range msg.enabledFiles {
			m.selectedMap[file] = true
			m.selectedOrder = append(m.selectedOrder, file)
		}

		// Build item list and display
		items := m.buildItemList()
		cmd := m.list.SetItems(items)
		m.loading = false
		logDebug("filesLoadedMsg: loading complete, displaying %d items", len(items))

		return m, cmd

	case itemsRefreshedMsg:
		// Item list was rebuilt (e.g., after hideUnlinked toggle)
		cmd := m.list.SetItems(msg.items)

		// If a cursor filename was specified, try to position cursor on that item
		if msg.cursorFileName != "" {
			m.setCursorToFile(msg.cursorFileName)
		}

		return m, cmd

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-helpBarReservedLines)
		return m, nil

	case tea.KeyMsg:
		// Don't handle keys while loading
		if m.loading {
			return m, nil
		}

		// Check if list is in filter mode
		wasFiltering := m.list.FilterState() == list.Filtering
		isFiltering := wasFiltering

		// Handle quit keys
		if key.Matches(msg, m.keys.Quit) {
			logDebug("Quit: user aborted")
			m.aborted = true
			return m, tea.Quit
		}

		// Handle confirm key (Enter)
		if key.Matches(msg, m.keys.Confirm) {
			if !isFiltering {
				logDebug("Confirm: user confirmed selection with %d items", len(m.selectedMap))
				return m, tea.Quit
			}
			// If filtering, let list.Model handle it
		}

		// Handle toggle selection (Space)
		if key.Matches(msg, m.keys.Select) {
			if !isFiltering {
				// Remember current cursor position before toggling
				var currentFileName string
				if item := m.list.SelectedItem(); item != nil {
					if fi, ok := item.(fileItem); ok {
						currentFileName = fi.name
					}
				}

				modeChanged := m.handleToggleSelection()
				logDebug("Toggle: selectedCount=%d", len(m.selectedMap))

				// If mode changed (hideUnlinked was auto-disabled), rebuild entire list
				// and preserve cursor on the toggled file
				if modeChanged {
					return m, m.rebuildItemsCmdWithCursor(currentFileName)
				}

				// Otherwise just refresh current item
				cmd := m.refreshCurrentItem()
				return m, cmd
			}
		}

		// Handle select all (Ctrl+A)
		if key.Matches(msg, m.keys.SelectAll) {
			if !isFiltering {
				// Remember current cursor position before selecting all
				var currentFileName string
				if item := m.list.SelectedItem(); item != nil {
					if fi, ok := item.(fileItem); ok {
						currentFileName = fi.name
					}
				}

				// Select all visible items
				countBefore := len(m.selectedMap)
				for _, item := range m.list.VisibleItems() {
					if fi, ok := item.(fileItem); ok {
						if !m.selectedMap[fi.name] {
							m.selectedMap[fi.name] = true
							m.selectedOrder = append(m.selectedOrder, fi.name)
						}
					}
				}
				logDebug("SelectAll: selected %d new items (total: %d), preserving cursor on: %s", len(m.selectedMap)-countBefore, len(m.selectedMap), currentFileName)
				// Refresh all items while preserving cursor position
				return m, m.rebuildItemsCmdWithCursor(currentFileName)
			}
		}

		// Handle deselect all (Ctrl+D)
		if key.Matches(msg, m.keys.DeselectAll) {
			if !isFiltering {
				// Remember current cursor position before deselecting all
				var currentFileName string
				if item := m.list.SelectedItem(); item != nil {
					if fi, ok := item.(fileItem); ok {
						currentFileName = fi.name
					}
				}

				logDebug("DeselectAll: clearing all selections")
				m.selectedMap = make(map[string]bool)
				m.selectedOrder = []string{}

				// Auto-disable hideUnlinked if no items are selected
				if m.shouldDisableHideMode() {
					logDebug("DeselectAll: disabling hideUnlinked mode, preserving cursor on: %s", currentFileName)
					m.hideUnlinked = false
				}

				return m, m.rebuildItemsCmdWithCursor(currentFileName)
			}
		}

		// Handle hide toggle (H)
		if key.Matches(msg, m.keys.HideToggle) {
			if !isFiltering && len(m.selectedMap) > 0 {
				// Remember current cursor position before toggling
				var currentFileName string
				if item := m.list.SelectedItem(); item != nil {
					if fi, ok := item.(fileItem); ok {
						currentFileName = fi.name
					}
				}

				m.hideUnlinked = !m.hideUnlinked
				logDebug("HideToggle: hideUnlinked=%t, preserving cursor on: %s", m.hideUnlinked, currentFileName)
				return m, m.rebuildItemsCmdWithCursor(currentFileName)
			}
		}

		// Delegate all other keys to list.Model (navigation, filtering, etc.)
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)

		// Log filter mode changes
		nowFiltering := m.list.FilterState() == list.Filtering
		if !wasFiltering && nowFiltering {
			logDebug("Filter: entered filter mode")
		} else if wasFiltering && !nowFiltering {
			logDebug("Filter: exited filter mode")
		}

		return m, cmd
	}

	// Delegate other messages to list.Model
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// buildItemList builds the list of items from availableFiles
// Respects hideUnlinked mode
func (m *multiSelectModel) buildItemList() []list.Item {
	// Preallocate with capacity to avoid reallocation
	items := make([]list.Item, 0, len(m.availableFiles))
	for _, name := range m.availableFiles {
		// In hideUnlinked mode, only show selected files
		if m.hideUnlinked && !m.selectedMap[name] {
			continue
		}

		items = append(items, fileItem{
			name:      name,
			isEnabled: m.selectedMap[name],
		})
	}
	return items
}

// handleToggleSelection toggles selection of the current item
// Returns true if hideUnlinked mode was auto-disabled (requires full list rebuild)
func (m *multiSelectModel) handleToggleSelection() bool {
	item := m.list.SelectedItem()
	if item == nil {
		return false
	}

	fi, ok := item.(fileItem)
	if !ok {
		return false
	}

	modeChanged := false

	// Toggle selection
	if m.selectedMap[fi.name] {
		// Deselect
		delete(m.selectedMap, fi.name)
		m.removeFromOrder(fi.name)

		// Auto-disable hideUnlinked if no items are selected
		if m.shouldDisableHideMode() {
			logDebug("Toggle: auto-disabling hideUnlinked mode (last item deselected)")
			m.hideUnlinked = false
			modeChanged = true
		}
	} else {
		// Select
		m.selectedMap[fi.name] = true
		m.selectedOrder = append(m.selectedOrder, fi.name)
	}

	return modeChanged
}

// removeFromOrder removes a file from selectedOrder
func (m *multiSelectModel) removeFromOrder(file string) {
	for i, f := range m.selectedOrder {
		if f == file {
			m.selectedOrder = append(m.selectedOrder[:i], m.selectedOrder[i+1:]...)
			return
		}
	}
}

// shouldDisableHideMode checks if hideUnlinked mode should be automatically disabled
// This happens when there are no selected items left
func (m *multiSelectModel) shouldDisableHideMode() bool {
	return m.hideUnlinked && len(m.selectedMap) == 0
}

// refreshCurrentItem refreshes the currently selected item to update its description
func (m *multiSelectModel) refreshCurrentItem() tea.Cmd {
	// Get current index
	index := m.list.Index()
	if index < 0 || index >= len(m.list.Items()) {
		return nil
	}

	// Get current item
	item := m.list.Items()[index]
	fi, ok := item.(fileItem)
	if !ok {
		return nil
	}

	// Update item with new enabled state
	updatedItem := fileItem{
		name:      fi.name,
		isEnabled: m.selectedMap[fi.name],
	}

	// Replace item in list
	return m.list.SetItem(index, updatedItem)
}

// rebuildItemsCmdWithCursor returns a command that rebuilds the item list
// and preserves cursor position on the specified filename
// Pass empty string to skip cursor positioning
func (m *multiSelectModel) rebuildItemsCmdWithCursor(fileName string) tea.Cmd {
	return func() tea.Msg {
		items := m.buildItemList()
		return itemsRefreshedMsg{
			items:          items,
			cursorFileName: fileName,
		}
	}
}

// setCursorToFile positions the cursor on the item with the specified filename
// If the file is not found in the current list, cursor stays at current position
func (m *multiSelectModel) setCursorToFile(fileName string) {
	if fileName == "" {
		return
	}

	items := m.list.Items()
	for i, item := range items {
		if fi, ok := item.(fileItem); ok {
			if fi.name == fileName {
				m.list.Select(i)
				logDebug("setCursorToFile: positioned cursor on %s at index %d", fileName, i)
				return
			}
		}
	}

	logDebug("setCursorToFile: file %s not found in list, cursor unchanged", fileName)
}

// View renders the UI
func (m multiSelectModel) View() string {
	// Handle aborted state
	if m.aborted {
		return ""
	}

	// Show loading state
	if m.loading {
		return "Loading files...\n"
	}

	// Show error state
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	// Delegate everything to list.Model (includes built-in help bar)
	return m.list.View()
}

// ShowFileSelect displays an interactive multi-select list in the terminal.
//
// The function loads files from the source directory and checks which ones are
// currently enabled (linked) in the target directory. It presents an interactive
// list where users can select/deselect files with keyboard navigation, filtering,
// and bulk operations.
//
// Visual feedback:
//   - Bold text: Linked/selected items
//   - Gray text: Unlinked items
//   - Bold green with ">": Current cursor position
//
// UI elements (conditional):
//   - Title: Shown only if title parameter is not empty
//   - Status bar: Shown only when title is set
//   - Help bar: Always visible (press ? to toggle short/full help)
//
// Parameters:
//   - sourceDir: Path to the source directory containing available files
//   - targetDir: Path to the target directory with symlinks
//   - title: Optional title to display above the list (empty = no title/status bar)
//
// Returns:
//   - []string: Ordered list of selected items (in selection order)
//   - error: Returns an error if user aborts (ctrl+c), loading fails, or no files found
//
// Keyboard shortcuts (short help):
//   - Space: Select/deselect item at cursor
//   - ↑/k, ↓/j: Move cursor up/down
//   - h: Toggle hide unlinked items (only when items are selected)
//   - /: Enter filter mode
//   - Enter: Confirm selection and exit
//   - ?: Show full help
//
// Additional shortcuts (full help with ?):
//   - g/G: Jump to top/bottom of list
//   - PgUp/PgDn, ctrl+b/ctrl+f: Page up/down
//   - ctrl+a: Select all visible items
//   - ctrl+d: Deselect all items
//   - ctrl+c: Abort without saving
//
// Example:
//
//	sourceDir := "/path/to/source/configs"
//	targetDir := "/path/to/target/configs"
//	selected, err := ShowFileSelect(sourceDir, targetDir, "Select files to link")
//	if err != nil {
//	    if strings.Contains(err.Error(), "user aborted") {
//	        fmt.Println("Operation cancelled")
//	        return
//	    }
//	    log.Fatal(err)
//	}
//	fmt.Printf("Selected: %v\n", selected)
func ShowFileSelect(sourceDir, targetDir, title string) ([]string, error) {
	// Create empty list (items loaded asynchronously in Init())
	// Use our custom delegate for simple rendering
	delegate := fileItemDelegate{}

	l := list.New([]list.Item{}, delegate, 0, 0) // width=0, height=0 (set via WindowSizeMsg)

	// Show status bar only if title is set
	if title != "" {
		l.Title = title
		l.SetShowTitle(true)
	} else {
		l.SetShowTitle(false)
	}

	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetFilteringEnabled(true)

	// Create model with our custom keys
	keys := defaultKeyMap()

	// Add our custom keybindings to the list's help
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Select, keys.HideToggle, keys.Filter, keys.Confirm}
	}

	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.Select, keys.SelectAll, keys.DeselectAll,
			keys.HideToggle, keys.Filter, keys.Confirm, keys.Quit,
		}
	}

	m := multiSelectModel{
		list:          l,
		sourceDir:     sourceDir,
		targetDir:     targetDir,
		selectedMap:   make(map[string]bool),
		selectedOrder: []string{},
		loading:       true,
		keys:          keys,
	}

	// Run the program
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("program error: %w", err)
	}

	// Type assert with check
	model, ok := finalModel.(multiSelectModel)
	if !ok {
		return nil, fmt.Errorf("unexpected model type")
	}

	// Check if aborted
	if model.aborted {
		return nil, fmt.Errorf("user aborted")
	}

	// Check for errors during loading
	if model.err != nil {
		return nil, model.err
	}

	// Check if no files were found
	if len(model.availableFiles) == 0 {
		return nil, fmt.Errorf("no files available to enable")
	}

	// Return selected items in order
	return model.selectedOrder, nil
}

// confirmModel is the Bubble Tea model for confirmation dialog
// It manages the state for a yes/no confirmation prompt
type confirmModel struct {
	message  string
	selected bool // true = yes, false = no
	aborted  bool
	width    int // Terminal width
}

// Init initializes the confirmation dialog model.
// No commands are needed for initialization.
func (m confirmModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the confirmation dialog.
// Supported keys:
//   - ctrl+c: Abort dialog
//   - enter: Confirm current selection
//   - left/right: Navigate between Yes/No
//   - y/Y: Quick select Yes and confirm
//   - n/N: Quick select No and confirm
func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.aborted = true
			return m, tea.Quit
		case "enter":
			return m, tea.Quit
		case "left":
			m.selected = true
		case "right":
			m.selected = false
		case "y", "Y":
			m.selected = true
			return m, tea.Quit
		case "n", "N":
			m.selected = false
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the confirmation dialog UI.
// Shows the message, Yes/No buttons with highlighting, and a help bar at the bottom.
// Returns empty string if dialog was aborted.
func (m confirmModel) View() string {
	if m.aborted {
		return ""
	}

	var b strings.Builder
	b.WriteString(m.message)
	b.WriteString("\n\n")

	var yesText, noText string
	if m.selected {
		yesText = stylePrompt.Render("[ Yes ]")
		noText = "[ No ]"
	} else {
		yesText = "[ Yes ]"
		noText = stylePrompt.Render("[ No ]")
	}

	b.WriteString(yesText)
	b.WriteString("  ")
	b.WriteString(noText)
	b.WriteString("\n\n")

	// Help text as inverse bar spanning full width
	helpText := "arrows: move | enter/y/n: select | ctrl+c: abort"
	helpBar := styleHelpBar.Width(m.width).Render(" " + helpText)
	b.WriteString(helpBar)

	return b.String()
}

// ShowConfirmation displays a yes/no confirmation dialog in the terminal.
//
// The function shows a message with two options (Yes/No) and returns the
// user's choice. The cursor starts on "Yes" by default.
//
// Parameters:
//   - message: The question or message to display to the user
//
// Returns:
//   - bool: true if user confirmed (pressed enter on "Yes"), false if declined
//   - error: Returns an error if user aborts (ctrl+c) or if there's a program error
//
// Keyboard shortcuts:
//   - ←/→: Move between Yes/No
//   - y/n: Quick select Yes/No and confirm
//   - Enter: Confirm current selection
//   - ctrl+c: Abort (returns error with "user aborted")
//
// Example:
//
//	confirmed, err := ShowConfirmation("Delete all files?")
//	if err != nil {
//	    if strings.Contains(err.Error(), "user aborted") {
//	        fmt.Println("Cancelled")
//	        return
//	    }
//	    log.Fatal(err)
//	}
//	if confirmed {
//	    // User selected "Yes"
//	    deleteFiles()
//	} else {
//	    // User selected "No"
//	    fmt.Println("Keeping files")
//	}
func ShowConfirmation(message string) (bool, error) {
	m := confirmModel{
		message:  message,
		selected: true, // Default to Yes
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return false, fmt.Errorf("program error: %w", err)
	}

	// Type assert with check
	model, ok := finalModel.(confirmModel)
	if !ok {
		return false, fmt.Errorf("unexpected model type")
	}

	if model.aborted {
		return false, fmt.Errorf("user aborted")
	}

	return model.selected, nil
}
