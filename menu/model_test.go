package menu

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestCursorMovement(t *testing.T) {
	m := initialModel()

	msg := tea.KeyMsg{Type: tea.KeyDown}
	updated, _ := m.Update(msg)
	m = updated.(model)

	if m.cursor != 1 {
		t.Errorf("expected cursor to be at 1, got %d", m.cursor)
	}

	msg = tea.KeyMsg{Type: tea.KeyUp}
	updated, _ = m.Update(msg)
	m = updated.(model)

	if m.cursor != 0 {
		t.Errorf("expected cursor to be at 0, got %d", m.cursor)
	}
}

func TestQuitShortcut(t *testing.T) {
	m := initialModel()
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Error("expected tea.Quit command on 'q' press, got nil")
	}
}

func TestViewIncludesHeader(t *testing.T) {
	m := initialModel()
	view := m.View()

	if !strings.Contains(view, "KIOSK") {
		t.Error("expected header to contain 'KIOSK'")
	}
}

func TestTmuxMenuNavigation(t *testing.T) {
	m := initialModel()

	// Select the "Manage Tmux Sessions" option
	m.cursor = 1 // Position of "📂  Manage Tmux Sessions" (now at index 1)
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updated, _ := m.Update(msg)
	m = updated.(model)

	if m.state != tmuxMenu {
		t.Errorf("expected state to be tmuxMenu, got %v", m.state)
	}

	// Test going back to main menu
	m.cursor = 4 // Position of "⬅️  Back to Main Menu"
	updated, _ = m.Update(msg)
	m = updated.(model)

	if m.state != mainMenu {
		t.Errorf("expected state to be mainMenu, got %v", m.state)
	}
}

func TestSessionInputView(t *testing.T) {
	m := initialModel()
	m.state = sessionInputState
	m.sessionInput = sessionNameInput{
		action: "create",
		value:  "test-session",
	}

	view := m.viewSessionInput()

	if !strings.Contains(view, "Create New Tmux Session") {
		t.Error("expected session input view to contain the title")
	}

	if !strings.Contains(view, "test-session") {
		t.Error("expected session input view to contain the input value")
	}
}

func TestSessionListView(t *testing.T) {
	m := initialModel()
	m.state = sessionListState
	m.choices = tmuxMenuModel().choices
	m.cursor = 1 // Position of "📜  List Sessions"
	m.sessionList = sessionListDisplay{
		sessions: []string{"session1", "session2"},
		selected: 0,
	}

	view := m.viewSessionList()

	if !strings.Contains(view, "Tmux Sessions") {
		t.Error("expected session list view to contain the title")
	}

	if !strings.Contains(view, "session1") || !strings.Contains(view, "session2") {
		t.Error("expected session list view to contain the session names")
	}
}
