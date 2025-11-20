package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// spawn.go - Tmux/Xterm spawn logic
// Based on tmuxplexer's proven implementation

// insideTmux checks if we're currently inside a tmux session
func insideTmux() bool {
	return os.Getenv("TMUX") != ""
}

// spawnSingle spawns a single command
func spawnSingle(item launchItem, mode spawnMode) tea.Cmd {
	return func() tea.Msg {
		var err error

		switch mode {
		case spawnTmuxSplitH:
			err = tmuxSplitHorizontal(item)
		case spawnTmuxSplitV:
			err = tmuxSplitVertical(item)
		case spawnTmuxWindow:
			err = tmuxNewWindow(item)
		case spawnXtermWindow:
			err = xtermWindow(item)
		case spawnCurrentPane:
			err = tmuxCurrentPane(item)
		default:
			// Default: use current pane (foreground)
			// Non-tmux mode is handled at higher level in model.go
			err = tmuxCurrentPane(item)
		}

		return spawnCompleteMsg{err: err}
	}
}

// spawnMultiple spawns multiple commands with a layout
// Uses the tmuxplexer strategy: create all panes, then apply layout
func spawnMultiple(items []launchItem, layout tmuxLayout) tea.Cmd {
	return func() tea.Msg {
		if len(items) == 0 {
			return spawnCompleteMsg{err: fmt.Errorf("no items to spawn")}
		}

		// Get common working directory (use first item's)
		baseDir := items[0].Cwd
		if baseDir == "" {
			baseDir = os.Getenv("HOME")
		}

		var err error
		if insideTmux() {
			// Inside tmux: spawn in current session
			err = spawnInCurrentSession(items, layout, baseDir)
		} else {
			// Outside tmux: create new session
			err = spawnNewSession(items, layout, baseDir)
		}

		return spawnCompleteMsg{err: err}
	}
}

// spawnInCurrentSession spawns items in the current tmux session
func spawnInCurrentSession(items []launchItem, layout tmuxLayout, baseDir string) error {
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
	// (This is the tmuxplexer pattern - don't track pane indices!)

	// First pane: send command to current pane
	if items[0].Command != "" {
		cwd := items[0].Cwd
		if cwd == "" {
			cwd = baseDir
		}

		// Change directory and run command
		if err := tmuxSendKeys(target+".0", fmt.Sprintf("cd '%s' && %s", cwd, items[0].Command)); err != nil {
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
			if err := tmuxSendKeys(paneTarget, items[i].Command); err != nil {
				return fmt.Errorf("failed to send keys to pane %d: %w", i, err)
			}
		}
	}

	return nil
}

// spawnNewSession creates a new tmux session with multiple panes
func spawnNewSession(items []launchItem, layout tmuxLayout, baseDir string) error {
	// Generate unique session name
	sessionName := generateSessionName(items[0].Name)

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
		if err := tmuxSendKeys(sessionName+":0.0", items[0].Command); err != nil {
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
			if err := tmuxSendKeys(paneTarget, items[i].Command); err != nil {
				return err
			}
		}
	}

	// Attach to session
	cmd = exec.Command("tmux", "attach", "-t", sessionName)
	return cmd.Run()
}

// tmuxSplitHorizontal splits the current pane horizontally
func tmuxSplitHorizontal(item launchItem) error {
	cwd := item.Cwd
	if cwd == "" {
		cwd = os.Getenv("HOME")
	}

	// Use shell to properly execute the command
	cmd := exec.Command("tmux", "split-window", "-h", "-c", cwd, "sh", "-c", item.Command)
	return cmd.Run()
}

// tmuxSplitVertical splits the current pane vertically
func tmuxSplitVertical(item launchItem) error {
	cwd := item.Cwd
	if cwd == "" {
		cwd = os.Getenv("HOME")
	}

	// Use shell to properly execute the command
	cmd := exec.Command("tmux", "split-window", "-v", "-c", cwd, "sh", "-c", item.Command)
	return cmd.Run()
}

// tmuxNewWindow creates a new tmux window
func tmuxNewWindow(item launchItem) error {
	cwd := item.Cwd
	if cwd == "" {
		cwd = os.Getenv("HOME")
	}

	// Use shell to properly execute the command
	cmd := exec.Command("tmux", "new-window", "-c", cwd, "-n", item.Name, "sh", "-c", item.Command)
	return cmd.Run()
}

// tmuxCurrentPane runs command in current pane
func tmuxCurrentPane(item launchItem) error {
	cwd := item.Cwd
	if cwd == "" {
		cwd = os.Getenv("HOME")
	}

	// Change directory and run command
	commandStr := fmt.Sprintf("cd '%s' && %s", cwd, item.Command)
	return tmuxSendKeys("", commandStr)
}

// xtermWindow spawns a new xterm window
func xtermWindow(item launchItem) error {
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

// tmuxSendKeys sends keys to a tmux pane (with Enter)
func tmuxSendKeys(target, keys string) error {
	args := []string{"send-keys"}
	if target != "" {
		args = append(args, "-t", target)
	}
	args = append(args, keys, "C-m") // C-m = Enter

	cmd := exec.Command("tmux", args...)
	return cmd.Run()
}

// generateSessionName creates a unique session name
func generateSessionName(baseName string) string {
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
