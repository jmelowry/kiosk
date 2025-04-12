package menu

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// execCommand is a variable that holds the exec.Command function
// It can be replaced with a mock function during tests
var execCommand = exec.Command

// InitTmux ensures tmux is properly configured for terminal restoration
func InitTmux() error {
	// Create a custom tmux configuration that helps preserve terminal settings
	cmds := []string{
		// Set terminal overrides to fix common issues
		"tmux set-option -g set-clipboard on",
		"tmux set-option -g default-terminal \"screen-256color\"",
		"tmux set-option -g terminal-overrides \"xterm*:smcup@:rmcup@\"",
		// Enable mouse support
		"tmux set-option -g mouse on",
		// Set escape-time to reduce delay issues
		"tmux set-option -sg escape-time 10",
		// Ensure alternate screen is properly used
		"tmux set-window-option -g alternate-screen on",
	}

	for _, cmd := range cmds {
		parts := strings.Fields(cmd)
		command := execCommand(parts[0], parts[1:]...)
		if err := command.Run(); err != nil {
			// Don't fail if the command doesn't work, just continue
			fmt.Printf("Warning: %s failed: %v\n", cmd, err)
		}
	}
	return nil
}

// CreateSession creates a new tmux session with the given name.
func CreateSession(name string) error {
	// Initialize tmux with proper settings first
	if err := InitTmux(); err != nil {
		return err
	}

	// Create the session with terminal restoration options
	cmd := execCommand("tmux", "new-session", "-d",
		"-s", name,
		"-e", "TERM=screen-256color")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tmux session: %w", err)
	}
	return nil
}

// ListSessions lists all active tmux sessions.
func ListSessions() ([]string, error) {
	cmd := execCommand("tmux", "list-sessions", "-F", "#S")
	output, err := cmd.Output()
	if err != nil {
		// Handle case where no sessions exist (tmux returns error)
		if strings.Contains(err.Error(), "no server running") {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to list tmux sessions: %w", err)
	}

	// Handle empty output case
	if len(strings.TrimSpace(string(output))) == 0 {
		return []string{}, nil
	}

	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

// AttachSession attaches to an existing tmux session by name.
// This function properly attaches to an existing tmux session by using syscall.Exec
// to replace the current process with a tmux process.
func AttachSession(name string) error {
	// Check if the session exists first
	sessions, err := ListSessions()
	if err != nil {
		return fmt.Errorf("failed to list tmux sessions: %w", err)
	}

	sessionExists := false
	for _, session := range sessions {
		if session == name {
			sessionExists = true
			break
		}
	}

	if !sessionExists {
		return fmt.Errorf("tmux session '%s' does not exist", name)
	}

	// Get the path to the tmux binary
	tmuxPath, err := exec.LookPath("tmux")
	if err != nil {
		return fmt.Errorf("tmux not found in PATH: %w", err)
	}

	// Execute tmux attach-session with options to properly restore terminal
	// -d: detaches other clients
	// The tput reset is used to reset terminal after detaching
	args := []string{"tmux", "attach-session", "-d", "-t", name}

	// Register a cleanup function to reset the terminal if this process exits
	cleanupCmd := exec.Command("sh", "-c", "tput reset")
	defer func() {
		cleanupCmd.Stdout = os.Stdout
		cleanupCmd.Run()
	}()

	return syscall.Exec(tmuxPath, args, os.Environ())
}

// KillSession kills a tmux session by name.
func KillSession(name string) error {
	cmd := execCommand("tmux", "kill-session", "-t", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to kill tmux session: %w", err)
	}
	return nil
}

// LaunchBtopInSession launches btop in a tmux session named "btop"
func LaunchBtopInSession() error {
	// Check if the session already exists and kill it if it does
	sessions, err := ListSessions()
	if err == nil {
		for _, session := range sessions {
			if session == "btop" {
				KillSession("btop")
				break
			}
		}
	}

	// Create a new session for btop
	if err := CreateSession("btop"); err != nil {
		return fmt.Errorf("failed to create tmux session for btop: %w", err)
	}

	// Send the command to run btop in the session
	cmd := execCommand("tmux", "send-keys", "-t", "btop", "btop", "C-m")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start btop: %w", err)
	}

	// Attach to the session
	return AttachSession("btop")
}

// GetTmuxCheatSheet returns a string containing common tmux commands
func GetTmuxCheatSheet() string {
	cheatSheet := `
TMUX CHEAT SHEET

Session Management:
  Ctrl+b d          Detach from session
  Ctrl+b $          Rename session
  
Window Management:  
  Ctrl+b c          Create new window
  Ctrl+b ,          Rename window
  Ctrl+b n          Next window
  Ctrl+b p          Previous window
  Ctrl+b w          List windows
  Ctrl+b &          Kill window
  
Pane Management:
  Ctrl+b "          Split horizontally
  Ctrl+b %          Split vertically
  Ctrl+b arrow      Switch to pane in that direction
  Ctrl+b z          Toggle pane zoom
  Ctrl+b x          Kill pane
  
Copy Mode:
  Ctrl+b [          Enter copy mode (use vi or emacs keys)
  q                 Quit copy mode
`
	return cheatSheet
}
