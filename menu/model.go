package menu

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type menuState int

const (
	mainMenu menuState = iota
	tmuxMenu
	sessionInputState
	sessionListState
)

type sessionNameInput struct {
	action    string
	value     string
	complete  bool
	cancelled bool
}

type sessionListDisplay struct {
	sessions []string
	selected int
}

type model struct {
	cursor       int
	choices      []string
	state        menuState
	sessionInput sessionNameInput
	sessionList  sessionListDisplay
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
			"📂  Manage Tmux Sessions",
			"🚪  Exit",
		},
		state: mainMenu,
		sessionInput: sessionNameInput{
			value: "",
		},
	}
}

func tmuxMenuModel() model {
	return model{
		choices: []string{
			"➕  Create New Session",
			"📜  List Sessions",
			"🔗  Attach to Session",
			"❌  Kill Session",
			"⬅️  Back to Main Menu",
		},
		state: tmuxMenu,
		sessionInput: sessionNameInput{
			value: "",
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case sessionInputState:
			return m.updateSessionInput(msg)
		case sessionListState:
			return m.updateSessionList(msg)
		default:
			return m.updateMenus(msg)
		}
	}
	return m, nil
}

func (m model) updateMenus(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		switch m.state {
		case mainMenu:
			switch m.choices[m.cursor] {
			case "📂  Manage Tmux Sessions":
				return tmuxMenuModel(), nil
			case "💾  Start Dev Session":
				err := CreateSession("dev-session")
				if err != nil {
					fmt.Println("Error starting tmux session:", err)
				} else {
					fmt.Println("Tmux session 'dev-session' started successfully.")
					if err := AttachSession("dev-session"); err != nil {
						fmt.Println("Error attaching to tmux session:", err)
					}
				}
				return m, tea.Quit
			case "🚪  Exit":
				return m, tea.Quit
			default:
				fmt.Println("Launching:", m.choices[m.cursor])
				return m, tea.Quit
			}
		case tmuxMenu:
			switch m.choices[m.cursor] {
			case "➕  Create New Session":
				m.state = sessionInputState
				m.sessionInput = sessionNameInput{
					action:   "create",
					value:    "",
					complete: false,
				}
				return m, nil
			case "📜  List Sessions":
				sessions, err := ListSessions()
				if err != nil {
					m.sessionList.sessions = []string{"Error listing sessions: " + err.Error()}
				} else if len(sessions) == 0 {
					m.sessionList.sessions = []string{"No active tmux sessions"}
				} else {
					m.sessionList.sessions = sessions
				}
				m.sessionList.selected = -1
				m.state = sessionListState
				return m, nil
			case "🔗  Attach to Session":
				sessions, err := ListSessions()
				if err != nil {
					m.sessionList.sessions = []string{"Error listing sessions: " + err.Error()}
				} else if len(sessions) == 0 {
					m.sessionList.sessions = []string{"No active tmux sessions"}
				} else {
					m.sessionList.sessions = sessions
				}
				m.sessionList.selected = 0
				m.state = sessionListState
				return m, nil
			case "❌  Kill Session":
				sessions, err := ListSessions()
				if err != nil {
					m.sessionList.sessions = []string{"Error listing sessions: " + err.Error()}
				} else if len(sessions) == 0 {
					m.sessionList.sessions = []string{"No active tmux sessions"}
				} else {
					m.sessionList.sessions = sessions
				}
				m.state = sessionListState
				m.sessionList.selected = 0
				return m, nil
			case "⬅️  Back to Main Menu":
				return initialModel(), nil
			}
		}
	}
	return m, nil
}

func (m model) updateSessionInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		if m.sessionInput.value == "" {
			return m, nil
		}
		m.sessionInput.complete = true

		switch m.sessionInput.action {
		case "create":
			err := CreateSession(m.sessionInput.value)
			if err != nil {
				fmt.Println("Error creating tmux session:", err)
			} else {
				fmt.Println("Tmux session '" + m.sessionInput.value + "' created successfully.")
				if err := AttachSession(m.sessionInput.value); err != nil {
					fmt.Println("Error attaching to tmux session:", err)
				}
			}
			return m, tea.Quit
		}
	case tea.KeyEsc:
		m.sessionInput.cancelled = true
		return tmuxMenuModel(), nil
	case tea.KeyBackspace:
		if len(m.sessionInput.value) > 0 {
			m.sessionInput.value = m.sessionInput.value[:len(m.sessionInput.value)-1]
		}
	case tea.KeyRunes:
		m.sessionInput.value += string(msg.Runes)
	}
	return m, nil
}

func (m model) updateSessionList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "escape":
		return tmuxMenuModel(), nil
	case "up", "k":
		if m.sessionList.selected > 0 {
			m.sessionList.selected--
		}
	case "down", "j":
		if m.sessionList.selected < len(m.sessionList.sessions)-1 {
			m.sessionList.selected++
		}
	case "enter":
		if m.sessionList.selected < 0 || len(m.sessionList.sessions) == 0 ||
			m.sessionList.sessions[0] == "No active tmux sessions" ||
			strings.HasPrefix(m.sessionList.sessions[0], "Error") {
			return tmuxMenuModel(), nil
		}

		selectedSession := m.sessionList.sessions[m.sessionList.selected]

		switch m.choices[m.cursor] {
		case "🔗  Attach to Session":
			fmt.Println("Attaching to session:", selectedSession)
			if err := AttachSession(selectedSession); err != nil {
				fmt.Println("Error attaching to tmux session:", err)
			}
			return m, tea.Quit
		case "❌  Kill Session":
			if err := KillSession(selectedSession); err != nil {
				fmt.Println("Error killing tmux session:", err)
			} else {
				fmt.Println("Tmux session '" + selectedSession + "' killed successfully.")
			}
			return tmuxMenuModel(), nil
		case "📜  List Sessions":
			return tmuxMenuModel(), nil
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case sessionInputState:
		return m.viewSessionInput()
	case sessionListState:
		return m.viewSessionList()
	default:
		return m.viewMenu()
	}
}

func (m model) viewMenu() string {
	s := headerStyle.Render("📟 KIOSK // Terminal Session Portal") + "\n\n"

	for i, choice := range m.choices {
		cursor := "   "
		if m.cursor == i {
			cursor = cursorStyle.Render("▶ ")
		} else {
			cursor = "   "
		}
		line := fmt.Sprintf("%s%s", cursor, choice)
		if m.cursor == i {
			line = selectedStyle.Render(line)
		}
		s += line + "\n"
	}

	s += "\n" + footerStyle.Render("↑ ↓ to navigate  •  ⏎ to launch  •  q to quit")

	return s
}

func (m model) viewSessionInput() string {
	var title string
	var prompt string
	switch m.sessionInput.action {
	case "create":
		title = "Create New Tmux Session"
		prompt = "Enter a name for the new session"
	}

	s := headerStyle.Render("📟 KIOSK // "+title) + "\n\n"

	s += prompt + ":\n\n"
	inputStyle := lipgloss.NewStyle().
		Foreground(colorText).
		Background(lipgloss.Color("#333333")).
		Padding(0, 1).
		Width(40)

	s += inputStyle.Render(m.sessionInput.value+"█") + "\n\n"

	s += footerStyle.Render("Enter to confirm • Esc to cancel")

	return s
}

func (m model) viewSessionList() string {
	var title string
	var action string
	switch m.choices[m.cursor] {
	case "📜  List Sessions":
		title = "Tmux Sessions"
		action = ""
	case "🔗  Attach to Session":
		title = "Select Session to Attach"
		action = "attach to"
	case "❌  Kill Session":
		title = "Select Session to Kill"
		action = "kill"
	}

	s := headerStyle.Render("📟 KIOSK // "+title) + "\n\n"

	if len(m.sessionList.sessions) == 0 {
		s += "No active tmux sessions\n"
	} else if strings.HasPrefix(m.sessionList.sessions[0], "Error") ||
		m.sessionList.sessions[0] == "No active tmux sessions" {
		s += m.sessionList.sessions[0] + "\n"
	} else {
		for i, session := range m.sessionList.sessions {
			cursor := "   "
			if m.sessionList.selected == i {
				cursor = cursorStyle.Render("▶ ")
			} else {
				cursor = "   "
			}

			line := fmt.Sprintf("%s%s", cursor, session)
			if m.sessionList.selected == i {
				line = selectedStyle.Render(line)
			}
			s += line + "\n"
		}

		if action != "" {
			s += "\n" + "Press Enter to " + action + " selected session"
		}
	}

	s += "\n\n" + footerStyle.Render("↑ ↓ to navigate  •  ⏎ to select  •  Esc to go back")

	return s
}
