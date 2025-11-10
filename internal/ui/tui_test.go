package ui

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/charmbracelet/bubbles/list"
)

// TestRemoveFromOrder tests removing items from the selection order
func TestRemoveFromOrder(t *testing.T) {
	m := &multiSelectModel{
		selectedOrder: []string{"a.txt", "b.txt", "c.txt"},
	}

	m.removeFromOrder("b.txt")

	expected := []string{"a.txt", "c.txt"}
	if !reflect.DeepEqual(m.selectedOrder, expected) {
		t.Errorf("expected %v, got %v", expected, m.selectedOrder)
	}
}

func TestRemoveFromOrder_FirstItem(t *testing.T) {
	m := &multiSelectModel{
		selectedOrder: []string{"a.txt", "b.txt", "c.txt"},
	}

	m.removeFromOrder("a.txt")

	expected := []string{"b.txt", "c.txt"}
	if !reflect.DeepEqual(m.selectedOrder, expected) {
		t.Errorf("expected %v, got %v", expected, m.selectedOrder)
	}
}

func TestRemoveFromOrder_LastItem(t *testing.T) {
	m := &multiSelectModel{
		selectedOrder: []string{"a.txt", "b.txt", "c.txt"},
	}

	m.removeFromOrder("c.txt")

	expected := []string{"a.txt", "b.txt"}
	if !reflect.DeepEqual(m.selectedOrder, expected) {
		t.Errorf("expected %v, got %v", expected, m.selectedOrder)
	}
}

func TestRemoveFromOrder_NonExistent(t *testing.T) {
	m := &multiSelectModel{
		selectedOrder: []string{"a.txt", "b.txt", "c.txt"},
	}

	m.removeFromOrder("nonexistent.txt")

	// Should remain unchanged
	expected := []string{"a.txt", "b.txt", "c.txt"}
	if !reflect.DeepEqual(m.selectedOrder, expected) {
		t.Errorf("expected %v, got %v", expected, m.selectedOrder)
	}
}

// TestBuildItemList tests the item list building logic
func TestBuildItemList(t *testing.T) {
	m := &multiSelectModel{
		availableFiles: []string{"a.txt", "b.txt", "c.txt"},
		selectedMap:    map[string]bool{"b.txt": true},
		hideUnlinked:   false,
	}

	items := m.buildItemList()

	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}

	// Check that b.txt is marked as enabled
	for _, item := range items {
		fi, ok := item.(fileItem)
		if !ok {
			t.Fatal("item is not fileItem")
		}
		if fi.name == "b.txt" && !fi.isEnabled {
			t.Error("b.txt should be enabled")
		}
		if fi.name != "b.txt" && fi.isEnabled {
			t.Errorf("%s should not be enabled", fi.name)
		}
	}
}

func TestBuildItemList_HideUnlinked(t *testing.T) {
	m := &multiSelectModel{
		availableFiles: []string{"a.txt", "b.txt", "c.txt", "d.txt"},
		selectedMap:    map[string]bool{"b.txt": true, "d.txt": true},
		hideUnlinked:   true,
	}

	items := m.buildItemList()

	if len(items) != 2 {
		t.Fatalf("expected 2 items in hideUnlinked mode, got %d", len(items))
	}

	// Verify only selected items are present
	itemNames := make(map[string]bool)
	for _, item := range items {
		fi, ok := item.(fileItem)
		if !ok {
			t.Fatal("item is not fileItem")
		}
		itemNames[fi.name] = true
	}

	if !itemNames["b.txt"] || !itemNames["d.txt"] {
		t.Error("hideUnlinked should only show b.txt and d.txt")
	}
}

func TestBuildItemList_EmptySelection(t *testing.T) {
	m := &multiSelectModel{
		availableFiles: []string{"a.txt", "b.txt"},
		selectedMap:    make(map[string]bool),
		hideUnlinked:   false,
	}

	items := m.buildItemList()

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	// All should be disabled
	for _, item := range items {
		fi, ok := item.(fileItem)
		if !ok {
			t.Fatal("item is not fileItem")
		}
		if fi.isEnabled {
			t.Errorf("%s should not be enabled", fi.name)
		}
	}
}

// TestHandleToggleSelection tests the selection toggle logic
func TestHandleToggleSelection(t *testing.T) {
	delegate := list.NewDefaultDelegate()
	l := list.New([]list.Item{
		fileItem{name: "a.txt", isEnabled: false},
		fileItem{name: "b.txt", isEnabled: true},
	}, delegate, 80, 10)

	m := &multiSelectModel{
		list:          l,
		selectedMap:   map[string]bool{"b.txt": true},
		selectedOrder: []string{"b.txt"},
	}

	// Select index 0 (a.txt)
	m.list.Select(0)
	m.handleToggleSelection()

	if !m.selectedMap["a.txt"] {
		t.Error("a.txt should be selected")
	}
	if len(m.selectedOrder) != 2 {
		t.Errorf("expected 2 selected items, got %d", len(m.selectedOrder))
	}

	// Deselect index 1 (b.txt)
	m.list.Select(1)
	m.handleToggleSelection()

	if m.selectedMap["b.txt"] {
		t.Error("b.txt should be deselected")
	}
	if len(m.selectedOrder) != 1 {
		t.Errorf("expected 1 selected item, got %d", len(m.selectedOrder))
	}
}

func TestHandleToggleSelection_EmptyList(t *testing.T) {
	delegate := fileItemDelegate{}
	l := list.New([]list.Item{}, delegate, 80, 10)

	m := &multiSelectModel{
		list:          l,
		selectedMap:   make(map[string]bool),
		selectedOrder: []string{},
	}

	// Should not panic
	m.handleToggleSelection()

	if len(m.selectedMap) != 0 {
		t.Error("selectedMap should remain empty")
	}
}

func TestHandleToggleSelection_LastItemInHideMode(t *testing.T) {
	delegate := list.NewDefaultDelegate()
	l := list.New([]list.Item{
		fileItem{name: "a.txt", isEnabled: true},
	}, delegate, 80, 10)

	m := &multiSelectModel{
		list:           l,
		availableFiles: []string{"a.txt", "b.txt"},
		selectedMap:    map[string]bool{"a.txt": true},
		selectedOrder:  []string{"a.txt"},
		hideUnlinked:   true,
	}

	// Deselect the last item
	m.list.Select(0)
	m.handleToggleSelection()

	// Should auto-disable hideUnlinked mode
	if m.hideUnlinked {
		t.Error("hideUnlinked should be disabled when last item is deselected")
	}
}

// TestFileItem tests the list.Item interface implementation
func TestFileItem_FilterValue(t *testing.T) {
	item := fileItem{name: "test.txt", isEnabled: true}
	if item.FilterValue() != "test.txt" {
		t.Errorf("expected 'test.txt', got '%s'", item.FilterValue())
	}
}

// TestUpdateMessageHandling tests the Update() message handling
func TestFilesLoadedMsg_Success(t *testing.T) {
	delegate := fileItemDelegate{}
	l := list.New([]list.Item{}, delegate, 80, 10)

	m := multiSelectModel{
		list:          l,
		selectedMap:   make(map[string]bool),
		selectedOrder: []string{},
		loading:       true,
	}

	msg := filesLoadedMsg{
		availableFiles: []string{"a.txt", "b.txt", "c.txt"},
		enabledFiles:   []string{"b.txt"},
		err:            nil,
	}

	result, _ := m.Update(msg)
	resultModel := result.(multiSelectModel)

	// Check available files
	if len(resultModel.availableFiles) != 3 {
		t.Errorf("expected 3 available files, got %d", len(resultModel.availableFiles))
	}

	// Check enabled file is selected
	if !resultModel.selectedMap["b.txt"] {
		t.Error("b.txt should be selected")
	}

	if len(resultModel.selectedOrder) != 1 {
		t.Errorf("expected 1 item in selectedOrder, got %d", len(resultModel.selectedOrder))
	}

	// Loading should be complete
	if resultModel.loading {
		t.Error("loading should be false after filesLoadedMsg")
	}

	// List should have items
	if len(resultModel.list.Items()) == 0 {
		t.Error("list should have items after loading")
	}
}

func TestFilesLoadedMsg_Error(t *testing.T) {
	m := multiSelectModel{
		selectedMap:   make(map[string]bool),
		selectedOrder: []string{},
		loading:       true,
	}

	msg := filesLoadedMsg{
		availableFiles: nil,
		enabledFiles:   nil,
		err:            fmt.Errorf("test error"),
	}

	result, _ := m.Update(msg)
	resultModel := result.(multiSelectModel)

	if resultModel.err == nil {
		t.Error("expected error to be set")
	}

	if !resultModel.aborted {
		t.Error("expected aborted to be true on error")
	}
}

// TestInit verifies that Init returns proper commands
func TestInit(t *testing.T) {
	m := multiSelectModel{
		sourceDir: "/test/source",
		targetDir: "/test/target",
	}

	cmd := m.Init()

	if cmd == nil {
		t.Fatal("Init should return a command")
	}

	// The command should be a batch (we can't easily test the contents without executing it)
	// At least verify it doesn't panic when called
	// Note: In real code, the command would trigger filesystem operations
}

// TestView tests the View rendering
func TestView_Aborted(t *testing.T) {
	m := multiSelectModel{
		aborted: true,
	}

	view := m.View()
	if view != "" {
		t.Error("aborted view should be empty")
	}
}

func TestView_Loading(t *testing.T) {
	m := multiSelectModel{
		loading: true,
	}

	view := m.View()
	if view != "Loading files...\n" {
		t.Errorf("unexpected loading view: %s", view)
	}
}

func TestView_Error(t *testing.T) {
	m := multiSelectModel{
		err: fmt.Errorf("test error"),
	}

	view := m.View()
	if !contains(view, "Error") {
		t.Errorf("error view should contain 'Error': %s", view)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestSetCursorToFile tests positioning cursor on a specific filename
func TestSetCursorToFile(t *testing.T) {
	items := []list.Item{
		fileItem{name: "a.txt", isEnabled: false},
		fileItem{name: "b.txt", isEnabled: true},
		fileItem{name: "c.txt", isEnabled: false},
		fileItem{name: "d.txt", isEnabled: true},
	}

	l := list.New(items, fileItemDelegate{}, 80, 20)
	m := &multiSelectModel{
		list: l,
	}

	// Test finding existing file
	m.setCursorToFile("c.txt")
	selected := m.list.SelectedItem()
	if fi, ok := selected.(fileItem); ok {
		if fi.name != "c.txt" {
			t.Errorf("Expected cursor on c.txt, got %s", fi.name)
		}
	} else {
		t.Error("Expected fileItem, got different type")
	}

	// Test cursor position
	if m.list.Index() != 2 {
		t.Errorf("Expected cursor at index 2, got %d", m.list.Index())
	}
}

// TestSetCursorToFile_NotFound tests behavior when file not in list
func TestSetCursorToFile_NotFound(t *testing.T) {
	items := []list.Item{
		fileItem{name: "a.txt", isEnabled: false},
		fileItem{name: "b.txt", isEnabled: true},
	}

	l := list.New(items, fileItemDelegate{}, 80, 20)
	l.Select(1) // Position on b.txt

	m := &multiSelectModel{
		list: l,
	}

	originalIndex := m.list.Index()

	// Try to find non-existent file
	m.setCursorToFile("nonexistent.txt")

	// Cursor should stay at original position
	if m.list.Index() != originalIndex {
		t.Errorf("Expected cursor to stay at index %d, got %d", originalIndex, m.list.Index())
	}
}

// TestSetCursorToFile_EmptyString tests behavior with empty filename
func TestSetCursorToFile_EmptyString(t *testing.T) {
	items := []list.Item{
		fileItem{name: "a.txt", isEnabled: false},
		fileItem{name: "b.txt", isEnabled: true},
	}

	l := list.New(items, fileItemDelegate{}, 80, 20)
	l.Select(1)

	m := &multiSelectModel{
		list: l,
	}

	originalIndex := m.list.Index()

	// Try with empty string
	m.setCursorToFile("")

	// Cursor should stay at original position
	if m.list.Index() != originalIndex {
		t.Errorf("Expected cursor to stay at index %d, got %d", originalIndex, m.list.Index())
	}
}

// TestRebuildItemsCmdWithCursor tests that cursor filename is preserved in message
func TestRebuildItemsCmdWithCursor(t *testing.T) {
	m := &multiSelectModel{
		availableFiles: []string{"a.txt", "b.txt", "c.txt"},
		selectedMap:    map[string]bool{"b.txt": true},
	}

	cmd := m.rebuildItemsCmdWithCursor("b.txt")
	msg := cmd()

	refreshMsg, ok := msg.(itemsRefreshedMsg)
	if !ok {
		t.Fatal("Expected itemsRefreshedMsg")
	}

	if refreshMsg.cursorFileName != "b.txt" {
		t.Errorf("Expected cursorFileName to be b.txt, got %s", refreshMsg.cursorFileName)
	}

	if len(refreshMsg.items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(refreshMsg.items))
	}
}
