package menu

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	cursor  int
	choices []string
}

var (
	colorAccent   = lipgloss.Color("#00d9ff")
	colorText     = lipgloss.Color("#bbbbbb")
	colorCursor   = lipgloss.Color("#00ffd9")
	colorSelect   = lipgloss.Color("#c2f0fc")
	headerStyle   = lipgloss.NewStyle().Foreground(colorAccent).Bold(true).Padding(1, 4)
	cursorStyle   = lipgloss.NewStyle().Foreground(colorCursor).Bold(true)
	selectedStyle = lipgloss.NewStyle().Foreground(colorSelect).Bold(true)
	footerStyle   = lipgloss.NewStyle().Foreground(colorText).Italic(true).PaddingTop(1)
)

func Start() error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func initialModel() model {
	return model{
		choices: []string{
			"💾  Start Dev Session",
			"📡  Attach to Logs",
			"🗃️  Launch Notes Panel",
			"🧼  Clean Workspace",
			"🚪  Exit",
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			// TODO: handle selection
			fmt.Println("Launching:", m.choices[m.cursor])
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := headerStyle.Render("📟 KIOSK // Terminal Session Portal") + "\n\n"

	for i, choice := range m.choices {
		cursor := "   " // Default spacing for unselected items
		if m.cursor == i {
			cursor = cursorStyle.Render("▶ ") // Add a space after the cursor for alignment
		} else {
			cursor = "   " // Ensure unselected items have the same spacing
		}
		line := fmt.Sprintf("%s%s", cursor, choice) // Combine cursor and choice
		if m.cursor == i {
			line = selectedStyle.Render(line) // Apply selected style
		}
		s += line + "\n"
	}

	s += "\n" + footerStyle.Render("↑ ↓ to navigate  •  ⏎ to launch  •  q to quit")

	return s
}
