package shared

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// shared/tmux.go - Unified Tmux Operations Layer
// Merges tui-launcher/spawn.go + tmuxplexer/tmux.go
// All tmux operations for both Launch tab (spawning) and Sessions/Templates tabs (management)

// ===== DETECTION =====

// InsideTmux checks if we're currently inside a tmux session
func InsideTmux() bool {
	return os.Getenv("TMUX") != ""
}

// ===== SPAWN OPERATIONS (from spawn.go) =====

// SpawnSingle spawns a single command
func SpawnSingle(item LaunchItem, mode SpawnMode) tea.Cmd {
	return func() tea.Msg {
		var err error

		switch mode {
		case SpawnTmuxSplitH:
			err = tmuxSplitHorizontal(item)
		case SpawnTmuxSplitV:
			err = tmuxSplitVertical(item)
		case SpawnTmuxWindow:
			err = tmuxNewWindow(item)
		case SpawnXtermWindow:
			err = xtermWindow(item)
		case SpawnCurrentPane:
			err = tmuxCurrentPane(item)
		default:
			// Default: use current pane (foreground)
			err = tmuxCurrentPane(item)
		}

		return SpawnCompleteMsg{Err: err}
	}
}

// SpawnMultiple spawns multiple commands with a layout
// Uses the tmuxplexer strategy: create all panes, then apply layout
func SpawnMultiple(items []LaunchItem, layout TmuxLayout) tea.Cmd {
	return func() tea.Msg {
		if len(items) == 0 {
			return SpawnCompleteMsg{Err: fmt.Errorf("no items to spawn")}
		}

		// Get common working directory (use first item's)
		baseDir := items[0].Cwd
		if baseDir == "" {
			baseDir = os.Getenv("HOME")
		}

		var err error
		if InsideTmux() {
			// Inside tmux: spawn in current session
			err = spawnInCurrentSession(items, layout, baseDir)
		} else {
			// Outside tmux: create new session
			err = spawnNewSession(items, layout, baseDir)
		}

		return SpawnCompleteMsg{Err: err}
	}
}

// spawnInCurrentSession spawns items in the current tmux session
func spawnInCurrentSession(items []LaunchItem, layout TmuxLayout, baseDir string) error {
	// Get current session name
	cmd := exec.Command("tmux", "display-message", "-p", "#{session_name}")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get session name: %w", err)
	}
	sessionName := strings.TrimSpace(string(output))

	// Get current window index
	cmd = exec.Command("tmux", "display-message", "-p", "#{window_index}")
	output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get window index: %w", err)
	}
	windowIndex := strings.TrimSpace(string(output))

	target := sessionName + ":" + windowIndex

	// Strategy: Create all panes first, then apply layout

	// First pane: send command to current pane
	if items[0].Command != "" {
		cwd := items[0].Cwd
		if cwd == "" {
			cwd = baseDir
		}

		// Change directory and run command
		if err := TmuxSendKeys(target+".0", fmt.Sprintf("cd '%s' && %s", cwd, items[0].Command)); err != nil {
			return fmt.Errorf("failed to send keys to pane 0: %w", err)
		}
	}

	// Create remaining panes
	for i := 1; i < len(items); i++ {
		cwd := items[i].Cwd
		if cwd == "" {
			cwd = baseDir
		}

		// Create pane with working directory
		cmd := exec.Command("tmux", "split-window", "-t", target, "-c", cwd)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create pane %d: %w", i, err)
		}

		// Small delay for stability (tmuxplexer uses 10ms)
		time.Sleep(10 * time.Millisecond)
	}

	// Apply selected layout
	layoutStr := layout.String()
	cmd = exec.Command("tmux", "select-layout", "-t", target, layoutStr)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to apply layout %s: %w", layoutStr, err)
	}

	// Send commands to panes (after layout is set)
	for i := 1; i < len(items); i++ {
		if items[i].Command != "" {
			paneTarget := fmt.Sprintf("%s.%d", target, i)
			if err := TmuxSendKeys(paneTarget, items[i].Command); err != nil {
				return fmt.Errorf("failed to send keys to pane %d: %w", i, err)
			}
		}
	}

	return nil
}

// spawnNewSession creates a new tmux session with multiple panes
func spawnNewSession(items []LaunchItem, layout TmuxLayout, baseDir string) error {
	// Generate unique session name
	sessionName := GenerateSessionName(items[0].Name)

	// Get working directory for first pane
	firstDir := items[0].Cwd
	if firstDir == "" {
		firstDir = baseDir
	}

	// Create new session (detached)
	cmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-c", firstDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	// Send command to first pane
	if items[0].Command != "" {
		if err := TmuxSendKeys(sessionName+":0.0", items[0].Command); err != nil {
			return err
		}
	}

	// Create remaining panes (tmuxplexer pattern)
	target := sessionName + ":0"
	for i := 1; i < len(items); i++ {
		cwd := items[i].Cwd
		if cwd == "" {
			cwd = baseDir
		}

		cmd := exec.Command("tmux", "split-window", "-t", target, "-c", cwd)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create pane %d: %w", i, err)
		}

		time.Sleep(10 * time.Millisecond)
	}

	// Apply layout
	layoutStr := layout.String()
	cmd = exec.Command("tmux", "select-layout", "-t", target, layoutStr)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to apply layout: %w", err)
	}

	// Send commands to remaining panes
	for i := 1; i < len(items); i++ {
		if items[i].Command != "" {
			paneTarget := fmt.Sprintf("%s:0.%d", sessionName, i)
			if err := TmuxSendKeys(paneTarget, items[i].Command); err != nil {
				return err
			}
		}
	}

	// Attach to session
	cmd = exec.Command("tmux", "attach", "-t", sessionName)
	return cmd.Run()
}

// tmuxSplitHorizontal splits the current pane horizontally
func tmuxSplitHorizontal(item LaunchItem) error {
	cwd := item.Cwd
	if cwd == "" {
		cwd = os.Getenv("HOME")
	}

	// Use shell to properly execute the command
	cmd := exec.Command("tmux", "split-window", "-h", "-c", cwd, "sh", "-c", item.Command)
	return cmd.Run()
}

// tmuxSplitVertical splits the current pane vertically
func tmuxSplitVertical(item LaunchItem) error {
	cwd := item.Cwd
	if cwd == "" {
		cwd = os.Getenv("HOME")
	}

	// Use shell to properly execute the command
	cmd := exec.Command("tmux", "split-window", "-v", "-c", cwd, "sh", "-c", item.Command)
	return cmd.Run()
}

// tmuxNewWindow creates a new tmux window
func tmuxNewWindow(item LaunchItem) error {
	cwd := item.Cwd
	if cwd == "" {
		cwd = os.Getenv("HOME")
	}

	// Use shell to properly execute the command
	cmd := exec.Command("tmux", "new-window", "-c", cwd, "-n", item.Name, "sh", "-c", item.Command)
	return cmd.Run()
}

// tmuxCurrentPane runs command in current pane
func tmuxCurrentPane(item LaunchItem) error {
	cwd := item.Cwd
	if cwd == "" {
		cwd = os.Getenv("HOME")
	}

	// Change directory and run command
	commandStr := fmt.Sprintf("cd '%s' && %s", cwd, item.Command)
	return TmuxSendKeys("", commandStr)
}

// xtermWindow spawns a new xterm window
func xtermWindow(item LaunchItem) error {
	cwd := item.Cwd
	if cwd == "" {
		cwd = os.Getenv("HOME")
	}

	// Use shell -c to run cd + command
	shellCmd := fmt.Sprintf("cd '%s' && %s", cwd, item.Command)
	cmd := exec.Command("xterm", "-e", "sh", "-c", shellCmd)

	// Start in background
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to spawn xterm: %w", err)
	}

	return nil
}

// TmuxSendKeys sends keys to a tmux pane (with Enter)
func TmuxSendKeys(target, keys string) error {
	args := []string{"send-keys"}
	if target != "" {
		args = append(args, "-t", target)
	}
	args = append(args, keys, "C-m") // C-m = Enter

	cmd := exec.Command("tmux", args...)
	return cmd.Run()
}

// GenerateSessionName creates a unique session name
func GenerateSessionName(baseName string) string {
	// Clean the name
	cleaned := strings.ToLower(baseName)
	cleaned = strings.ReplaceAll(cleaned, " ", "-")

	// Remove special characters (keep alphanumeric and hyphens)
	cleaned = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, cleaned)

	// Add timestamp to ensure uniqueness
	timestamp := time.Now().Format("150405") // HHMMSS
	return fmt.Sprintf("%s-%s", cleaned, timestamp)
}

// ===== SESSION MANAGEMENT (from tmuxplexer/tmux.go) =====

// ListSessions returns all tmux sessions
func ListSessions() ([]TmuxSession, error) {
	// Ensure tmux server is running
	if err := startTmuxServer(); err != nil {
		return nil, fmt.Errorf("failed to start tmux server: %w", err)
	}

	// Format: session_name|session_windows|session_attached|session_created
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}|#{session_windows}|#{session_attached}|#{session_created}")
	output, err := cmd.Output()
	if err != nil {
		// No sessions exist yet - return empty list
		if strings.Contains(err.Error(), "no server running") || strings.Contains(err.Error(), "no sessions") {
			return []TmuxSession{}, nil
		}
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	sessions := make([]TmuxSession, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 4 {
			continue
		}

		// Parse attached count - session is attached if count >= 1
		attachedCount := 0
		fmt.Sscanf(parts[2], "%d", &attachedCount)

		session := TmuxSession{
			Name:       parts[0],
			Attached:   attachedCount >= 1,
			Created:    formatTime(parts[3]),
			LastActive: "active", // Will be updated
		}

		// Parse window count
		fmt.Sscanf(parts[1], "%d", &session.Windows)

		// Get working directory for the first pane
		workingDirCmd := exec.Command("tmux", "display-message", "-p", "-t", session.Name+":0.0", "#{pane_current_path}")
		workingDirOutput, err := workingDirCmd.Output()
		if err == nil {
			session.WorkingDir = strings.TrimSpace(string(workingDirOutput))
			// Get git branch if in a git repo
			session.GitBranch = getGitBranch(session.WorkingDir)
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// ListWindows returns windows for a session
func ListWindows(sessionName string) ([]TmuxWindow, error) {
	// Format: window_index|window_name|window_panes|window_active
	cmd := exec.Command("tmux", "list-windows", "-t", sessionName, "-F", "#{window_index}|#{window_name}|#{window_panes}|#{window_active}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list windows for session %s: %w", sessionName, err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	windows := make([]TmuxWindow, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 4 {
			continue
		}

		window := TmuxWindow{
			Name:   parts[1],
			Active: parts[3] == "1",
		}

		fmt.Sscanf(parts[0], "%d", &window.Index)
		fmt.Sscanf(parts[2], "%d", &window.Panes)

		windows = append(windows, window)
	}

	return windows, nil
}

// ListPanes returns panes for a window
func ListPanes(sessionName string, windowIndex int) ([]TmuxPane, error) {
	target := fmt.Sprintf("%s:%d", sessionName, windowIndex)
	// Format: pane_id|pane_index|pane_width|pane_height|pane_active|pane_current_command
	cmd := exec.Command("tmux", "list-panes", "-t", target, "-F", "#{pane_id}|#{pane_index}|#{pane_width}|#{pane_height}|#{pane_active}|#{pane_current_command}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list panes for %s: %w", target, err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	panes := make([]TmuxPane, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 6 {
			continue
		}

		pane := TmuxPane{
			ID:      parts[0],
			Active:  parts[4] == "1",
			Command: parts[5],
		}

		fmt.Sscanf(parts[1], "%d", &pane.Index)
		fmt.Sscanf(parts[2], "%d", &pane.Width)
		fmt.Sscanf(parts[3], "%d", &pane.Height)

		panes = append(panes, pane)
	}

	return panes, nil
}

// CapturePane returns full scrollback content of a pane
func CapturePane(paneID string) (string, error) {
	cmd := exec.Command("tmux", "capture-pane", "-p", "-S", "-", "-t", paneID)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to capture pane %s: %w", paneID, err)
	}
	return string(output), nil
}

// AttachToSession attaches to a session
func AttachToSession(sessionName string) tea.Cmd {
	return func() tea.Msg {
		// TODO: Implement attach logic
		return nil
	}
}

// KillSession kills a tmux session
func KillSession(sessionName string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("tmux", "kill-session", "-t", sessionName)
		err := cmd.Run()
		return SessionKilledMsg{SessionName: sessionName, Err: err}
	}
}

// RenameSession renames a session
func RenameSession(oldName, newName string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("tmux", "rename-session", "-t", oldName, newName)
		err := cmd.Run()
		return SessionRenamedMsg{OldName: oldName, NewName: newName, Err: err}
	}
}

// ===== HELPER FUNCTIONS =====

// checkTmuxRunning checks if tmux server is running
func checkTmuxRunning() bool {
	cmd := exec.Command("tmux", "info")
	err := cmd.Run()
	return err == nil
}

// startTmuxServer starts the tmux server if not running
func startTmuxServer() error {
	if checkTmuxRunning() {
		return nil
	}
	cmd := exec.Command("tmux", "start-server")
	return cmd.Run()
}

// getGitBranch returns the current git branch for a directory
func getGitBranch(dir string) string {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// formatTime formats a unix timestamp to a readable string
func formatTime(unixStr string) string {
	var timestamp int64
	fmt.Sscanf(unixStr, "%d", &timestamp)

	if timestamp == 0 {
		return "unknown"
	}

	t := time.Unix(timestamp, 0)
	now := time.Now()

	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		mins := int(diff.Minutes())
		return fmt.Sprintf("%dm ago", mins)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%dh ago", hours)
	} else {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}
