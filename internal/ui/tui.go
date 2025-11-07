package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ANSI color codes for terminal styling
const (
	colorReset  = "\033[0m"
	colorPrompt = "\033[1;32m" // Bold Green
	colorDim    = "\033[2;90m" // Dim Gray
	colorHelp   = "\033[90m"   // Gray
	colorCursor = "\033[7m"    // Reverse video
	colorNormal = "\033[0m"    // Normal/White
)

// multiSelectModel is the Bubble Tea model for multi-select UI
// It manages the state for selecting multiple items from a list
type multiSelectModel struct {
	choices         []string        // Available choices
	choicesLower    []string        // Lowercase versions for efficient filtering
	selected        map[string]bool // Selected items
	selectedOrder   []string        // Order of selection for result
	cursor          int             // Cursor position
	filter          string          // Filter text
	filtering       bool            // Filter mode active
	filtered        []string        // Filtered choices
	aborted         bool            // User pressed ESC
	title           string          // Optional title
	maxVisibleItems int             // Maximum items to show before pagination
}

// Init initializes the model
func (m multiSelectModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m multiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.aborted = true
			return m, tea.Quit
		case "esc":
			m.aborted = true
			return m, tea.Quit
		case "enter":
			if !m.filtering {
				return m, tea.Quit
			}
			// Exit filter mode on enter, keep the filter
			m.filtering = false
			m.clampCursor()
			return m, nil
		case "/":
			if !m.filtering {
				m.filtering = true
				// Don't clear the existing filter, just allow editing
				return m, nil
			}
		case "up":
			if !m.filtering && m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if !m.filtering {
				choices := m.getVisibleChoices()
				if m.cursor < len(choices)-1 {
					m.cursor++
				}
			}
		case " ":
			if !m.filtering {
				choices := m.getVisibleChoices()
				if m.cursor >= 0 && m.cursor < len(choices) {
					choice := choices[m.cursor]
					// Toggle selection
					if m.selected[choice] {
						delete(m.selected, choice)
						// Remove from order
						for i, s := range m.selectedOrder {
							if s == choice {
								m.selectedOrder = append(m.selectedOrder[:i], m.selectedOrder[i+1:]...)
								break
							}
						}
					} else {
						m.selected[choice] = true
						m.selectedOrder = append(m.selectedOrder, choice)
					}
				}
			}
		case "backspace":
			if m.filtering && len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.updateFiltered()
				m.clampCursor()
			}
		default:
			// Add character to filter
			if m.filtering && len(msg.String()) == 1 {
				m.filter += msg.String()
				m.updateFiltered()
				m.clampCursor()
			}
		}
	}
	return m, nil
}

// View renders the UI
func (m multiSelectModel) View() string {
	if m.aborted {
		return ""
	}

	var b strings.Builder

	// Title (optional)
	if m.title != "" {
		b.WriteString(m.title)
		b.WriteString("\n\n")
	}

	// Filter prompt
	if m.filtering {
		b.WriteString(colorPrompt)
		b.WriteString("$ ")
		b.WriteString(colorReset)
		b.WriteString(m.filter)
		b.WriteString(colorCursor)
		b.WriteString(" ")
		b.WriteString(colorReset)
		b.WriteString("\n\n")
	}

	// Choices
	choices := m.getVisibleChoices()
	visibleStart := 0
	visibleEnd := len(choices)

	// Pagination
	if len(choices) > m.maxVisibleItems {
		if m.cursor >= m.maxVisibleItems {
			visibleStart = m.cursor - m.maxVisibleItems + 1
		}
		visibleEnd = visibleStart + m.maxVisibleItems
		if visibleEnd > len(choices) {
			visibleEnd = len(choices)
			visibleStart = max(0, visibleEnd-m.maxVisibleItems)
		}
	}

	for i := visibleStart; i < visibleEnd; i++ {
		choice := choices[i]
		cursor := " "
		if i == m.cursor && !m.filtering {
			cursor = "▶"
		}

		// Text color: dim gray if not selected, normal if selected
		textStyle := colorDim
		if m.selected[choice] {
			textStyle = colorNormal
		}

		b.WriteString(cursor)
		b.WriteString(" ")
		b.WriteString(textStyle)
		b.WriteString(choice)
		b.WriteString(colorReset)
		b.WriteString("\n")
	}

	// Help text
	b.WriteString("\n")
	b.WriteString(colorHelp)
	if !m.filtering {
		b.WriteString("space: toggle | /: filter | enter: confirm | esc: abort")
	} else {
		b.WriteString("type to filter | enter: exit filter | esc: abort")
	}
	b.WriteString(colorReset)

	return b.String()
}

// getVisibleChoices returns filtered or all choices
func (m *multiSelectModel) getVisibleChoices() []string {
	if m.filter != "" {
		return m.filtered
	}
	return m.choices
}

// updateFiltered updates the filtered list
func (m *multiSelectModel) updateFiltered() {
	m.filtered = []string{}
	filterLower := strings.ToLower(m.filter)
	for i, choice := range m.choices {
		if strings.Contains(m.choicesLower[i], filterLower) {
			m.filtered = append(m.filtered, choice)
		}
	}
}

// clampCursor ensures cursor is within valid range
func (m *multiSelectModel) clampCursor() {
	choices := m.getVisibleChoices()
	if len(choices) == 0 {
		m.cursor = 0
		return
	}
	if m.cursor >= len(choices) {
		m.cursor = len(choices) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// ShowMultiSelect displays a multi-select UI for choosing files to enable
func ShowMultiSelect(availableFiles []string, currentlyEnabled []string, title string, maxVisibleItems int) ([]string, error) {
	if len(availableFiles) == 0 {
		return nil, fmt.Errorf("no files available to enable")
	}

	// Create initial selection map and order
	selected := make(map[string]bool)
	selectedOrder := []string{}
	for _, file := range currentlyEnabled {
		selected[file] = true
		selectedOrder = append(selectedOrder, file)
	}

	// Pre-compute lowercase versions for efficient filtering
	choicesLower := make([]string, len(availableFiles))
	for i, choice := range availableFiles {
		choicesLower[i] = strings.ToLower(choice)
	}

	// Create model
	m := multiSelectModel{
		choices:         availableFiles,
		choicesLower:    choicesLower,
		selected:        selected,
		selectedOrder:   selectedOrder,
		cursor:          0,
		title:           title,
		maxVisibleItems: maxVisibleItems,
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

	// Return selected items in order
	return model.selectedOrder, nil
}

// confirmModel is the Bubble Tea model for confirmation dialog
// It manages the state for a yes/no confirmation prompt
type confirmModel struct {
	message  string
	selected bool // true = yes, false = no
	aborted  bool
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
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

func (m confirmModel) View() string {
	if m.aborted {
		return ""
	}

	var b strings.Builder
	b.WriteString(m.message)
	b.WriteString("\n\n")

	yesStyle := "[ Yes ]"
	noStyle := "[ No ]"

	if m.selected {
		yesStyle = colorPrompt + "[ Yes ]" + colorReset
	} else {
		noStyle = colorPrompt + "[ No ]" + colorReset
	}

	b.WriteString(yesStyle)
	b.WriteString("  ")
	b.WriteString(noStyle)
	b.WriteString("\n\n")
	b.WriteString(colorHelp)
	b.WriteString("arrows: move | enter/y/n: select | esc: abort")
	b.WriteString(colorReset)

	return b.String()
}

// ShowConfirmation displays a confirmation dialog
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
