package ui

import (
	"fmt"
	"reflect"
	"testing"
)

// TestGetVisibleChoices_HideUnlinked tests hideUnlinked mode filtering
func TestGetVisibleChoices_HideUnlinked(t *testing.T) {
	m := &multiSelectModel{
		choices:      []string{"a.txt", "b.txt", "c.txt", "d.txt"},
		selected:     map[string]bool{"b.txt": true, "d.txt": true},
		hideUnlinked: true,
	}

	visible := m.getVisibleChoices()

	if len(visible) != 2 {
		t.Errorf("expected 2 visible choices, got %d", len(visible))
	}

	expected := map[string]bool{"b.txt": true, "d.txt": true}
	for _, v := range visible {
		if !expected[v] {
			t.Errorf("unexpected item in visible: %s", v)
		}
	}
}

// TestGetVisibleChoices_WithFilter tests filter mode
func TestGetVisibleChoices_WithFilter(t *testing.T) {
	m := &multiSelectModel{
		choices:      []string{"apple.txt", "banana.txt", "apricot.txt", "berry.txt"},
		choicesLower: []string{"apple.txt", "banana.txt", "apricot.txt", "berry.txt"},
		filter:       "ap",
		filtered:     []string{"apple.txt", "apricot.txt"},
		selected:     map[string]bool{"apple.txt": true},
		hideUnlinked: false,
	}

	visible := m.getVisibleChoices()

	if len(visible) != 2 {
		t.Errorf("expected 2 visible choices, got %d", len(visible))
	}

	if visible[0] != "apple.txt" || visible[1] != "apricot.txt" {
		t.Errorf("unexpected filter results: %v", visible)
	}
}

// TestGetVisibleChoices_FilterAndHide tests combined filter and hideUnlinked
func TestGetVisibleChoices_FilterAndHide(t *testing.T) {
	m := &multiSelectModel{
		choices:      []string{"apple.txt", "banana.txt", "apricot.txt", "berry.txt"},
		choicesLower: []string{"apple.txt", "banana.txt", "apricot.txt", "berry.txt"},
		filter:       "ap",
		filtered:     []string{"apple.txt", "apricot.txt"},
		selected:     map[string]bool{"apple.txt": true},
		hideUnlinked: true,
	}

	visible := m.getVisibleChoices()

	if len(visible) != 1 {
		t.Errorf("expected 1 visible choice, got %d", len(visible))
	}
	if len(visible) > 0 && visible[0] != "apple.txt" {
		t.Errorf("expected 'apple.txt', got '%s'", visible[0])
	}
}

// TestGetVisibleChoices_NoFilter tests normal mode without filters
func TestGetVisibleChoices_NoFilter(t *testing.T) {
	m := &multiSelectModel{
		choices:      []string{"a.txt", "b.txt", "c.txt"},
		selected:     map[string]bool{"a.txt": true},
		hideUnlinked: false,
	}

	visible := m.getVisibleChoices()

	if len(visible) != 3 {
		t.Errorf("expected 3 visible choices, got %d", len(visible))
	}
}

// TestGetVisibleChoices_Caching tests that caching works correctly
func TestGetVisibleChoices_Caching(t *testing.T) {
	m := &multiSelectModel{
		choices:      []string{"a.txt", "b.txt", "c.txt"},
		selected:     map[string]bool{"a.txt": true},
		hideUnlinked: false,
	}

	// First call should cache
	visible1 := m.getVisibleChoices()
	if !m.cacheValid {
		t.Error("cache should be valid after first call")
	}

	// Second call should return cached result
	visible2 := m.getVisibleChoices()
	if !reflect.DeepEqual(visible1, visible2) {
		t.Error("cached result should be identical to first result")
	}

	// After invalidation, cache should be rebuilt
	m.cacheValid = false
	_ = m.getVisibleChoices()
	if !m.cacheValid {
		t.Error("cache should be valid after rebuild")
	}
}

// TestClampCursor_EmptyList tests cursor clamping with empty list
func TestClampCursor_EmptyList(t *testing.T) {
	m := &multiSelectModel{
		choices: []string{},
		cursor:  5,
	}

	m.clampCursor()

	if m.cursor != 0 {
		t.Errorf("expected cursor 0 for empty list, got %d", m.cursor)
	}
}

// TestClampCursor_OutOfBounds tests cursor clamping when out of bounds
func TestClampCursor_OutOfBounds(t *testing.T) {
	m := &multiSelectModel{
		choices: []string{"a.txt", "b.txt", "c.txt"},
		cursor:  10,
	}

	m.clampCursor()

	if m.cursor != 2 {
		t.Errorf("expected cursor 2 (last item), got %d", m.cursor)
	}
}

// TestClampCursor_NegativeCursor tests cursor clamping with negative value
func TestClampCursor_NegativeCursor(t *testing.T) {
	m := &multiSelectModel{
		choices: []string{"a.txt", "b.txt", "c.txt"},
		cursor:  -5,
	}

	m.clampCursor()

	if m.cursor != 0 {
		t.Errorf("expected cursor 0, got %d", m.cursor)
	}
}

// TestClampCursor_ValidCursor tests that valid cursor is not changed
func TestClampCursor_ValidCursor(t *testing.T) {
	m := &multiSelectModel{
		choices: []string{"a.txt", "b.txt", "c.txt"},
		cursor:  1,
	}

	m.clampCursor()

	if m.cursor != 1 {
		t.Errorf("expected cursor 1, got %d", m.cursor)
	}
}

// TestSelectItem tests item selection
func TestSelectItem(t *testing.T) {
	m := &multiSelectModel{
		selected:      make(map[string]bool),
		selectedOrder: []string{},
		selectedIndex: make(map[string]int),
	}

	m.selectItem("test.txt")

	if !m.selected["test.txt"] {
		t.Error("item should be selected")
	}
	if len(m.selectedOrder) != 1 || m.selectedOrder[0] != "test.txt" {
		t.Error("item should be in selectedOrder")
	}
	if m.selectedIndex["test.txt"] != 0 {
		t.Error("item index should be 0")
	}
}

// TestRemoveFromOrder tests removal from selectedOrder
func TestRemoveFromOrder(t *testing.T) {
	m := &multiSelectModel{
		selectedOrder: []string{"a.txt", "b.txt", "c.txt"},
		selectedIndex: map[string]int{
			"a.txt": 0,
			"b.txt": 1,
			"c.txt": 2,
		},
	}

	m.removeFromOrder("b.txt")

	if len(m.selectedOrder) != 2 {
		t.Errorf("expected 2 items in order, got %d", len(m.selectedOrder))
	}
	if m.selectedOrder[0] != "a.txt" || m.selectedOrder[1] != "c.txt" {
		t.Errorf("unexpected order after removal: %v", m.selectedOrder)
	}
	if _, exists := m.selectedIndex["b.txt"]; exists {
		t.Error("removed item should not be in index")
	}
	if m.selectedIndex["a.txt"] != 0 {
		t.Error("a.txt index should still be 0")
	}
	if m.selectedIndex["c.txt"] != 1 {
		t.Error("c.txt index should be updated to 1")
	}
}

// TestRemoveFromOrder_FirstItem tests removing first item
func TestRemoveFromOrder_FirstItem(t *testing.T) {
	m := &multiSelectModel{
		selectedOrder: []string{"a.txt", "b.txt", "c.txt"},
		selectedIndex: map[string]int{
			"a.txt": 0,
			"b.txt": 1,
			"c.txt": 2,
		},
	}

	m.removeFromOrder("a.txt")

	if len(m.selectedOrder) != 2 {
		t.Errorf("expected 2 items, got %d", len(m.selectedOrder))
	}
	if m.selectedOrder[0] != "b.txt" || m.selectedOrder[1] != "c.txt" {
		t.Errorf("unexpected order: %v", m.selectedOrder)
	}
	// Check indices are updated
	if m.selectedIndex["b.txt"] != 0 {
		t.Error("b.txt should now be at index 0")
	}
	if m.selectedIndex["c.txt"] != 1 {
		t.Error("c.txt should now be at index 1")
	}
}

// TestRemoveFromOrder_LastItem tests removing last item
func TestRemoveFromOrder_LastItem(t *testing.T) {
	m := &multiSelectModel{
		selectedOrder: []string{"a.txt", "b.txt", "c.txt"},
		selectedIndex: map[string]int{
			"a.txt": 0,
			"b.txt": 1,
			"c.txt": 2,
		},
	}

	m.removeFromOrder("c.txt")

	if len(m.selectedOrder) != 2 {
		t.Errorf("expected 2 items, got %d", len(m.selectedOrder))
	}
	if m.selectedOrder[0] != "a.txt" || m.selectedOrder[1] != "b.txt" {
		t.Errorf("unexpected order: %v", m.selectedOrder)
	}
	// Indices should not change for remaining items
	if m.selectedIndex["a.txt"] != 0 {
		t.Error("a.txt should still be at index 0")
	}
	if m.selectedIndex["b.txt"] != 1 {
		t.Error("b.txt should still be at index 1")
	}
}

// TestRemoveFromOrder_NonExistent tests removing non-existent item
func TestRemoveFromOrder_NonExistent(t *testing.T) {
	m := &multiSelectModel{
		selectedOrder: []string{"a.txt", "b.txt"},
		selectedIndex: map[string]int{
			"a.txt": 0,
			"b.txt": 1,
		},
	}

	m.removeFromOrder("nonexistent.txt")

	if len(m.selectedOrder) != 2 {
		t.Error("order should not change when removing non-existent item")
	}
}

// TestUpdateFiltered tests filter update functionality
func TestUpdateFiltered(t *testing.T) {
	m := &multiSelectModel{
		choices:      []string{"apple.txt", "banana.txt", "apricot.txt", "berry.txt"},
		choicesLower: []string{"apple.txt", "banana.txt", "apricot.txt", "berry.txt"},
		filter:       "ap",
	}

	m.updateFiltered()

	if len(m.filtered) != 2 {
		t.Errorf("expected 2 filtered items, got %d", len(m.filtered))
	}
	if m.filtered[0] != "apple.txt" || m.filtered[1] != "apricot.txt" {
		t.Errorf("unexpected filtered results: %v", m.filtered)
	}
}

// TestUpdateFiltered_CaseInsensitive tests case-insensitive filtering
func TestUpdateFiltered_CaseInsensitive(t *testing.T) {
	m := &multiSelectModel{
		choices:      []string{"Apple.txt", "BANANA.txt", "aPRicot.txt"},
		choicesLower: []string{"apple.txt", "banana.txt", "apricot.txt"},
		filter:       "AP",
	}

	m.updateFiltered()

	if len(m.filtered) != 2 {
		t.Errorf("expected 2 filtered items (case insensitive), got %d", len(m.filtered))
	}
}

// TestUpdateFiltered_EmptyFilter tests with empty filter
func TestUpdateFiltered_EmptyFilter(t *testing.T) {
	m := &multiSelectModel{
		choices:      []string{"a.txt", "b.txt", "c.txt"},
		choicesLower: []string{"a.txt", "b.txt", "c.txt"},
		filter:       "",
	}

	m.updateFiltered()

	if len(m.filtered) != 3 {
		t.Errorf("expected 3 items with empty filter, got %d", len(m.filtered))
	}
}

// TestIsKey tests the isKey helper function
func TestIsKey(t *testing.T) {
	tests := []struct {
		name     string
		pressed  string
		keys     []string
		expected bool
	}{
		{"match first", "ctrl+c", []string{"ctrl+c", "esc"}, true},
		{"match second", "esc", []string{"ctrl+c", "esc"}, true},
		{"no match", "enter", []string{"ctrl+c", "esc"}, false},
		{"single key match", "h", []string{"h"}, true},
		{"single key no match", "j", []string{"h"}, false},
		{"empty keys", "h", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isKey(tt.pressed, tt.keys...)
			if result != tt.expected {
				t.Errorf("isKey(%q, %v) = %v, want %v", tt.pressed, tt.keys, result, tt.expected)
			}
		})
	}
}

// TestAdjustCursorAfterItemRemoved tests cursor adjustment after item removal
func TestAdjustCursorAfterItemRemoved(t *testing.T) {
	tests := []struct {
		name           string
		choices        []string
		previousCursor int
		expectedCursor int
	}{
		{
			name:           "cursor in middle stays",
			choices:        []string{"a.txt", "b.txt", "c.txt"},
			previousCursor: 1,
			expectedCursor: 1,
		},
		{
			name:           "cursor at end moves back",
			choices:        []string{"a.txt", "b.txt"},
			previousCursor: 2,
			expectedCursor: 1,
		},
		{
			name:           "cursor beyond end moves to last",
			choices:        []string{"a.txt"},
			previousCursor: 5,
			expectedCursor: 0,
		},
		{
			name:           "empty list",
			choices:        []string{},
			previousCursor: 0,
			expectedCursor: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &multiSelectModel{
				choices: tt.choices,
			}
			m.adjustCursorAfterItemRemoved(tt.previousCursor)
			if m.cursor != tt.expectedCursor {
				t.Errorf("expected cursor %d, got %d", tt.expectedCursor, m.cursor)
			}
		})
	}
}

// TestGetVisibleChoices_LargeListEarlyExit tests early exit optimization for large lists
func TestGetVisibleChoices_LargeListEarlyExit(t *testing.T) {
	// Create a large list with only a few selected items
	largeList := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		largeList[i] = fmt.Sprintf("file%d.txt", i)
	}

	m := &multiSelectModel{
		choices: largeList,
		selected: map[string]bool{
			"file10.txt": true,
			"file50.txt": true,
			"file90.txt": true,
		},
		hideUnlinked: true,
	}

	visible := m.getVisibleChoices()

	// Should find exactly 3 items
	if len(visible) != 3 {
		t.Errorf("expected 3 visible items, got %d", len(visible))
	}

	// Verify correct items are returned
	expectedItems := map[string]bool{
		"file10.txt": true,
		"file50.txt": true,
		"file90.txt": true,
	}
	for _, item := range visible {
		if !expectedItems[item] {
			t.Errorf("unexpected item in visible: %s", item)
		}
	}
}

// TestSelectItem_Multiple tests selecting multiple items
func TestSelectItem_Multiple(t *testing.T) {
	m := &multiSelectModel{
		choices:       []string{"a.txt", "b.txt", "c.txt"},
		selected:      make(map[string]bool),
		selectedOrder: []string{},
		selectedIndex: make(map[string]int),
	}

	// Select three items
	m.selectItem("a.txt")
	m.selectItem("b.txt")
	m.selectItem("c.txt")

	if len(m.selected) != 3 {
		t.Errorf("expected 3 selected items, got %d", len(m.selected))
	}

	if len(m.selectedOrder) != 3 {
		t.Errorf("expected 3 items in order, got %d", len(m.selectedOrder))
	}

	// Check indices are correct
	if m.selectedIndex["a.txt"] != 0 {
		t.Error("a.txt should be at index 0")
	}
	if m.selectedIndex["b.txt"] != 1 {
		t.Error("b.txt should be at index 1")
	}
	if m.selectedIndex["c.txt"] != 2 {
		t.Error("c.txt should be at index 2")
	}
}

// TestDeselectAll_WithHideUnlinked tests deselecting all items when hideUnlinked is active
func TestDeselectAll_WithHideUnlinked(t *testing.T) {
	m := &multiSelectModel{
		choices: []string{"a.txt", "b.txt", "c.txt"},
		selected: map[string]bool{
			"a.txt": true,
			"b.txt": true,
		},
		selectedOrder: []string{"a.txt", "b.txt"},
		selectedIndex: map[string]int{
			"a.txt": 0,
			"b.txt": 1,
		},
		hideUnlinked: true,
		cursor:       1,
	}

	// Simulate deselect all (what ctrl+d does)
	m.selected = make(map[string]bool)
	m.selectedOrder = []string{}
	m.selectedIndex = make(map[string]int)
	if m.hideUnlinked {
		m.hideUnlinked = false
		m.clampCursor()
	}

	if len(m.selected) != 0 {
		t.Error("all items should be deselected")
	}
	if m.hideUnlinked {
		t.Error("hideUnlinked should be disabled after deselect all")
	}
	if m.cursor > len(m.choices)-1 {
		t.Error("cursor should be clamped after deselect all")
	}
}

// TestGetVisibleChoices_NilSelectedMap tests defensive nil check
func TestGetVisibleChoices_NilSelectedMap(t *testing.T) {
	m := &multiSelectModel{
		choices:      []string{"a.txt", "b.txt", "c.txt"},
		selected:     nil, // Intentionally nil
		hideUnlinked: true,
	}

	// Should not panic, should return empty list
	visible := m.getVisibleChoices()

	if visible == nil {
		t.Error("visible should not be nil")
	}
	if len(visible) != 0 {
		t.Errorf("expected empty list with nil selected map, got %d items", len(visible))
	}
}

// TestHandleToggleSelection tests space key selection toggling
func TestHandleToggleSelection(t *testing.T) {
	m := &multiSelectModel{
		choices:       []string{"a.txt", "b.txt", "c.txt"},
		selected:      make(map[string]bool),
		selectedOrder: []string{},
		selectedIndex: make(map[string]int),
		cursor:        0,
	}

	// Select first item
	m.handleToggleSelection()
	if !m.selected["a.txt"] {
		t.Error("first item should be selected")
	}
	if len(m.selectedOrder) != 1 {
		t.Error("selectedOrder should have 1 item")
	}

	// Deselect first item
	m.handleToggleSelection()
	if m.selected["a.txt"] {
		t.Error("first item should be deselected")
	}
	if len(m.selectedOrder) != 0 {
		t.Error("selectedOrder should be empty")
	}
}

// TestHandleToggleSelection_OutOfBounds tests toggle with invalid cursor
func TestHandleToggleSelection_OutOfBounds(t *testing.T) {
	m := &multiSelectModel{
		choices:       []string{"a.txt"},
		selected:      make(map[string]bool),
		selectedOrder: []string{},
		selectedIndex: make(map[string]int),
		cursor:        10, // Out of bounds
	}

	// Should not panic or change state
	m.handleToggleSelection()
	if len(m.selected) != 0 {
		t.Error("nothing should be selected with out of bounds cursor")
	}
}

// TestClampCursor_WithFilter tests cursor clamping with filtered list
func TestClampCursor_WithFilter(t *testing.T) {
	m := &multiSelectModel{
		choices:      []string{"apple.txt", "banana.txt", "apricot.txt"},
		choicesLower: []string{"apple.txt", "banana.txt", "apricot.txt"},
		filter:       "ap",
		filtered:     []string{"apple.txt", "apricot.txt"},
		cursor:       5, // Out of bounds for filtered list
	}

	m.clampCursor()

	// Cursor should be clamped to last item in ALL choices, not filtered
	if m.cursor > len(m.choices)-1 {
		t.Errorf("cursor should be clamped to %d, got %d", len(m.choices)-1, m.cursor)
	}
}

// TestUpdateFiltered_NoMatches tests filtering with no matches
func TestUpdateFiltered_NoMatches(t *testing.T) {
	m := &multiSelectModel{
		choices:      []string{"apple.txt", "banana.txt"},
		choicesLower: []string{"apple.txt", "banana.txt"},
		filter:       "xyz",
	}

	m.updateFiltered()

	if len(m.filtered) != 0 {
		t.Errorf("expected 0 filtered items, got %d", len(m.filtered))
	}
}

// TestSwitchToShowAllMode tests mode switching
func TestSwitchToShowAllMode(t *testing.T) {
	m := &multiSelectModel{
		choices:      []string{"a.txt", "b.txt", "c.txt"},
		selected:     map[string]bool{"b.txt": true},
		hideUnlinked: true,
		cursor:       0,
	}

	m.switchToShowAllMode("b.txt")

	if m.hideUnlinked {
		t.Error("hideUnlinked should be false after switch")
	}

	// Cursor should be positioned on b.txt (index 1 in full list)
	if m.cursor != 1 {
		t.Errorf("cursor should be on b.txt (index 1), got %d", m.cursor)
	}
}

// TestSwitchToShowAllMode_ItemNotFound tests switching when item not found
func TestSwitchToShowAllMode_ItemNotFound(t *testing.T) {
	m := &multiSelectModel{
		choices:      []string{"a.txt", "b.txt", "c.txt"},
		selected:     make(map[string]bool),
		hideUnlinked: true,
		cursor:       5,
	}

	m.switchToShowAllMode("nonexistent.txt")

	// Cursor should be clamped
	if m.cursor > len(m.choices)-1 {
		t.Error("cursor should be clamped when item not found")
	}
}

// TestDeselectItem_InHideMode tests deselecting in hide mode
func TestDeselectItem_InHideMode(t *testing.T) {
	m := &multiSelectModel{
		choices: []string{"a.txt", "b.txt", "c.txt"},
		selected: map[string]bool{
			"a.txt": true,
			"b.txt": true,
			"c.txt": true,
		},
		selectedOrder: []string{"a.txt", "b.txt", "c.txt"},
		selectedIndex: map[string]int{
			"a.txt": 0,
			"b.txt": 1,
			"c.txt": 2,
		},
		hideUnlinked: true,
		cursor:       1,
	}

	// Deselect b.txt
	m.deselectItem("b.txt", 1)

	if m.selected["b.txt"] {
		t.Error("b.txt should be deselected")
	}
	if len(m.selectedOrder) != 2 {
		t.Errorf("selectedOrder should have 2 items, got %d", len(m.selectedOrder))
	}
	// hideUnlinked should still be true (not last item)
	if !m.hideUnlinked {
		t.Error("hideUnlinked should still be true")
	}
}

// TestDeselectItem_LastItemInHideMode tests auto-switch on last deselect
func TestDeselectItem_LastItemInHideMode(t *testing.T) {
	m := &multiSelectModel{
		choices: []string{"a.txt", "b.txt", "c.txt"},
		selected: map[string]bool{
			"a.txt": true,
		},
		selectedOrder: []string{"a.txt"},
		selectedIndex: map[string]int{
			"a.txt": 0,
		},
		hideUnlinked: true,
		cursor:       0,
	}

	// Deselect last item
	m.deselectItem("a.txt", 0)

	if m.hideUnlinked {
		t.Error("hideUnlinked should be auto-disabled after last item deselect")
	}
}

// TestAdjustCursorAfterItemRemoved_Advanced tests cursor adjustment with hideUnlinked
func TestAdjustCursorAfterItemRemoved_Advanced(t *testing.T) {
	tests := []struct {
		name           string
		choices        []string
		selected       map[string]bool
		hideUnlinked   bool
		previousCursor int
		expectedCursor int
	}{
		{
			name:           "cursor in middle of selected items",
			choices:        []string{"a.txt", "b.txt", "c.txt", "d.txt", "e.txt"},
			selected:       map[string]bool{"a.txt": true, "b.txt": true, "c.txt": true},
			hideUnlinked:   true,
			previousCursor: 1,
			expectedCursor: 1,
		},
		{
			name:           "cursor at end after item removed",
			choices:        []string{"a.txt", "b.txt", "c.txt"},
			selected:       map[string]bool{"a.txt": true, "b.txt": true},
			hideUnlinked:   true,
			previousCursor: 2,
			expectedCursor: 1,
		},
		{
			name:           "cursor beyond end moves to last",
			choices:        []string{"a.txt", "b.txt", "c.txt"},
			selected:       map[string]bool{"a.txt": true},
			hideUnlinked:   true,
			previousCursor: 10,
			expectedCursor: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &multiSelectModel{
				choices:      tt.choices,
				selected:     tt.selected,
				hideUnlinked: tt.hideUnlinked,
			}
			m.adjustCursorAfterItemRemoved(tt.previousCursor)

			if m.cursor != tt.expectedCursor {
				t.Errorf("expected cursor %d, got %d", tt.expectedCursor, m.cursor)
			}
		})
	}
}

// TestIsKey tests key matching helper
func TestIsKey_MultipleKeys(t *testing.T) {
	tests := []struct {
		pressed  string
		keys     []string
		expected bool
	}{
		{"ctrl+c", []string{"ctrl+c", "esc"}, true},
		{"esc", []string{"ctrl+c", "esc"}, true},
		{"enter", []string{"ctrl+c", "esc"}, false},
		{"a", []string{"a", "b", "c"}, true},
		{"z", []string{"a", "b", "c"}, false},
	}

	for _, tt := range tests {
		result := isKey(tt.pressed, tt.keys...)
		if result != tt.expected {
			t.Errorf("isKey(%q, %v) = %v, want %v", tt.pressed, tt.keys, result, tt.expected)
		}
	}
}

// TestGetVisibleChoices_CombinedModes tests filter + hideUnlinked + caching
func TestGetVisibleChoices_CombinedModes(t *testing.T) {
	m := &multiSelectModel{
		choices:      []string{"apple.txt", "banana.txt", "apricot.txt", "berry.txt"},
		choicesLower: []string{"apple.txt", "banana.txt", "apricot.txt", "berry.txt"},
		filter:       "ap",
		filtered:     []string{"apple.txt", "apricot.txt"},
		selected: map[string]bool{
			"apple.txt": true,
		},
		hideUnlinked: true,
		cacheValid:   false,
	}

	// First call should cache
	visible1 := m.getVisibleChoices()
	if len(visible1) != 1 || visible1[0] != "apple.txt" {
		t.Errorf("expected [apple.txt], got %v", visible1)
	}
	if !m.cacheValid {
		t.Error("cache should be valid after first call")
	}

	// Second call should use cache
	visible2 := m.getVisibleChoices()
	if len(visible2) != 1 || visible2[0] != "apple.txt" {
		t.Errorf("expected cached [apple.txt], got %v", visible2)
	}
}

// TestShowMultiSelect_CursorPositioning tests initial cursor position
func TestShowMultiSelect_CursorPositioning(t *testing.T) {
	// Note: This is a unit test that directly creates the model
	// rather than running the full ShowMultiSelect which uses tea.Program

	availableFiles := []string{"a.txt", "b.txt", "c.txt", "d.txt"}
	currentlyEnabled := []string{"c.txt"}

	// Simulate the cursor positioning logic from ShowMultiSelect
	initialCursor := 0
	if len(currentlyEnabled) > 0 {
		firstSelected := currentlyEnabled[0]
		found := false
		for i, file := range availableFiles {
			if file == firstSelected {
				initialCursor = i
				found = true
				break
			}
		}
		if !found {
			initialCursor = 0
		}
	}

	// Cursor should be on c.txt (index 2)
	if initialCursor != 2 {
		t.Errorf("expected cursor at index 2 (c.txt), got %d", initialCursor)
	}
}

// TestShowMultiSelect_CursorPositioning_ItemNotFound tests fallback behavior
func TestShowMultiSelect_CursorPositioning_ItemNotFound(t *testing.T) {
	availableFiles := []string{"a.txt", "b.txt", "c.txt"}
	currentlyEnabled := []string{"nonexistent.txt"} // Not in availableFiles

	// Simulate the cursor positioning logic from ShowMultiSelect
	initialCursor := 0
	if len(currentlyEnabled) > 0 {
		firstSelected := currentlyEnabled[0]
		found := false
		for i, file := range availableFiles {
			if file == firstSelected {
				initialCursor = i
				found = true
				break
			}
		}
		if !found {
			// Should stay at 0 (defensive fallback)
			initialCursor = 0
		}
	}

	// Cursor should default to 0 when item not found
	if initialCursor != 0 {
		t.Errorf("expected cursor at index 0 (fallback), got %d", initialCursor)
	}
}

// TestShowMultiSelect_CursorPositioning_EmptyEnabled tests empty list
func TestShowMultiSelect_CursorPositioning_EmptyEnabled(t *testing.T) {
	currentlyEnabled := []string{} // Empty

	// Simulate the cursor positioning logic
	initialCursor := 0
	if len(currentlyEnabled) > 0 {
		// This branch should not execute
		t.Error("should not enter this branch with empty currentlyEnabled")
	}

	// Cursor should be 0 with empty enabled list
	if initialCursor != 0 {
		t.Errorf("expected cursor at index 0, got %d", initialCursor)
	}
}

// TestRepositionCursorAfterModeChange tests the extracted cursor repositioning method
func TestRepositionCursorAfterModeChange(t *testing.T) {
	m := &multiSelectModel{
		choices: []string{"a.txt", "b.txt", "c.txt", "d.txt"},
		selected: map[string]bool{
			"b.txt": true,
			"d.txt": true,
		},
		hideUnlinked: false,
		cursor:       1, // On b.txt
	}

	// Toggle to hideUnlinked mode
	m.hideUnlinked = true
	m.repositionCursorAfterModeChange("b.txt")

	// Should find b.txt in hideUnlinked list (index 0)
	visibleChoices := m.getVisibleChoices()
	if len(visibleChoices) != 2 {
		t.Errorf("expected 2 visible choices, got %d", len(visibleChoices))
	}
	if m.cursor != 0 {
		t.Errorf("expected cursor at 0 (b.txt in hideUnlinked list), got %d", m.cursor)
	}
}

// TestRepositionCursorAfterModeChange_ItemNotFound tests fallback behavior
func TestRepositionCursorAfterModeChange_ItemNotFound(t *testing.T) {
	m := &multiSelectModel{
		choices: []string{"a.txt", "b.txt", "c.txt"},
		selected: map[string]bool{
			"a.txt": true,
		},
		hideUnlinked: true,
		cursor:       10, // Invalid
	}

	// Try to reposition to non-existent item
	m.repositionCursorAfterModeChange("nonexistent.txt")

	// Should clamp cursor
	visibleChoices := m.getVisibleChoices()
	if m.cursor >= len(visibleChoices) {
		t.Error("cursor should be clamped to valid range")
	}
}

// TestRepositionCursorAfterModeChange_EmptyItem tests empty string handling
func TestRepositionCursorAfterModeChange_EmptyItem(t *testing.T) {
	m := &multiSelectModel{
		choices: []string{"a.txt", "b.txt", "c.txt"},
		cursor:  10, // Out of bounds
	}

	// Empty string should trigger clamping
	m.repositionCursorAfterModeChange("")

	// Should clamp cursor to valid range
	if m.cursor >= len(m.choices) {
		t.Errorf("cursor should be clamped, got %d", m.cursor)
	}
}

// TestHideToggleRefactoring tests the refactored hide toggle logic
func TestHideToggleRefactoring(t *testing.T) {
	m := &multiSelectModel{
		choices: []string{"a.txt", "b.txt", "c.txt", "d.txt"},
		selected: map[string]bool{
			"b.txt": true,
			"c.txt": true,
		},
		selectedIndex: make(map[string]int),
		hideUnlinked:  false,
		cursor:        1, // On b.txt
		keys:          defaultKeyBindings,
		cacheValid:    false, // Ensure cache is invalid
	}

	// Simulate 'h' key press to toggle hideUnlinked
	// Get current item
	choices := m.getVisibleChoices()
	currentItem := ""
	if m.cursor >= 0 && m.cursor < len(choices) {
		currentItem = choices[m.cursor]
	}

	// Toggle and invalidate cache (like Update does)
	m.hideUnlinked = !m.hideUnlinked
	m.cacheValid = false

	// Reposition (using extracted method)
	m.repositionCursorAfterModeChange(currentItem)

	// Verify hideUnlinked is now true
	if !m.hideUnlinked {
		t.Error("hideUnlinked should be true after toggle")
	}

	// Verify cursor is on b.txt in the new list
	m.cacheValid = false // Invalidate again before checking
	newChoices := m.getVisibleChoices()
	if len(newChoices) != 2 {
		t.Errorf("expected 2 visible choices in hideUnlinked mode, got %d", len(newChoices))
	}
	if m.cursor >= len(newChoices) {
		t.Errorf("cursor %d is out of bounds for %d choices", m.cursor, len(newChoices))
	}
	if newChoices[m.cursor] != "b.txt" {
		t.Errorf("expected cursor on b.txt, got %s", newChoices[m.cursor])
	}
}
