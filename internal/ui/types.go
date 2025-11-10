package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Lipgloss styles for file item rendering
// These are defined at package level to avoid repeated allocations during rendering
var (
	// Cursor styles (item under cursor with ">")
	styleCursorEnabled  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")) // Bold green for cursor on linked
	styleCursorDisabled = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))            // Green (not bold) for cursor on unlinked

	// Normal item styles (not under cursor)
	styleEnabled  = lipgloss.NewStyle().Bold(true)                        // Bold for linked items
	styleDisabled = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray for unlinked
)

// Message types for async operations

// filesLoadedMsg is sent when both available and enabled files have been loaded
type filesLoadedMsg struct {
	availableFiles []string
	enabledFiles   []string
	err            error
}

// itemsRefreshedMsg is sent when the item list needs to be rebuilt
// (e.g., after toggling hideUnlinked mode or changing selections)
type itemsRefreshedMsg struct {
	items          []list.Item
	cursorFileName string // Optional: filename to position cursor on after rebuild
}

// fileItem represents a single file in the list
// It implements the list.Item interface for use with bubbles/list
type fileItem struct {
	name      string
	isEnabled bool // Whether this file is currently selected/linked
}

// FilterValue implements list.Item interface
// Returns the string to be used for filtering
func (i fileItem) FilterValue() string {
	return i.name
}

// fileItemDelegate is a custom delegate for rendering file items
type fileItemDelegate struct{}

// Height returns the height of each list item (1 line)
func (d fileItemDelegate) Height() int { return 1 }

// Spacing returns spacing between items
func (d fileItemDelegate) Spacing() int { return 0 }

// Update handles delegate-specific updates
func (d fileItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// Render draws a single item in the list
// Uses pre-defined package-level styles to avoid repeated allocations
func (d fileItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	fi, ok := listItem.(fileItem)
	if !ok {
		return
	}

	// Render based on cursor position
	if index == m.Index() {
		// Current cursor position with ">"
		if fi.isEnabled {
			// Linked item at cursor: bold green
			fmt.Fprint(w, styleCursorEnabled.Render("> "+fi.name))
		} else {
			// Unlinked item at cursor: green (not bold)
			fmt.Fprint(w, styleCursorDisabled.Render("> "+fi.name))
		}
	} else {
		// Normal item: styled based on selection status
		if fi.isEnabled {
			// Linked items are bold
			fmt.Fprint(w, styleEnabled.Render("  "+fi.name))
		} else {
			// Unlinked items are gray
			fmt.Fprint(w, styleDisabled.Render("  "+fi.name))
		}
	}
}
