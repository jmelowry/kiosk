package menu

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// CreateSession creates a new tmux session with the given name.
func CreateSession(name string) error {
	cmd := exec.Command("tmux", "new-session", "-d", "-s", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tmux session: %w", err)
	}
	return nil
}

// ListSessions lists all active tmux sessions.
func ListSessions() ([]string, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#S")
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

	// Execute tmux attach-session by replacing the current process
	// This is necessary for proper terminal handling
	args := []string{"tmux", "attach-session", "-t", name}
	return syscall.Exec(tmuxPath, args, os.Environ())
}

// KillSession kills a tmux session by name.
func KillSession(name string) error {
	cmd := exec.Command("tmux", "kill-session", "-t", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to kill tmux session: %w", err)
	}
	return nil
}
