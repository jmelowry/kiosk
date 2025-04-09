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
