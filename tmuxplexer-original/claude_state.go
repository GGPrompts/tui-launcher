package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	claudeStateDir = "/tmp/claude-code-state"
	staleThreshold = 60 * time.Second // Consider state stale after 60 seconds (was 5s, caused flashing)
)

// ClaudeState represents the state of a Claude Code session
type ClaudeState struct {
	SessionID   string                 `json:"session_id"`
	Status      string                 `json:"status"` // idle, processing, tool_use, awaiting_input, working
	CurrentTool string                 `json:"current_tool"`
	WorkingDir  string                 `json:"working_dir"`
	LastUpdated string                 `json:"last_updated"`
	TmuxPane    string                 `json:"tmux_pane"`
	PID         int                    `json:"pid"`
	HookType    string                 `json:"hook_type"`
	Details     map[string]interface{} `json:"details"`
}

// detectClaudeSession checks if a tmux session is running Claude Code
func detectClaudeSession(sessionName string) bool {
	command := getPaneCommand(sessionName)
	commandLower := strings.ToLower(command)
	// Check specifically for claude or claude-code commands (case-insensitive)
	// Also check for node running claude (common pattern)
	return strings.Contains(commandLower, "claude") ||
	       (strings.Contains(commandLower, "node") && hasClaudeInCmdline(sessionName))
}

// detectCodexSession checks if a tmux session is running Codex
func detectCodexSession(sessionName string) bool {
	command := getPaneCommand(sessionName)
	commandLower := strings.ToLower(command)
	return strings.Contains(commandLower, "codex")
}

// detectGeminiSession checks if a tmux session is running Gemini
func detectGeminiSession(sessionName string) bool {
	command := getPaneCommand(sessionName)
	commandLower := strings.ToLower(command)
	return strings.Contains(commandLower, "gemini")
}

// hasClaudeInCmdline checks if the full command line contains "claude"
func hasClaudeInCmdline(sessionName string) bool {
	cmd := exec.Command("tmux", "display-message", "-p", "-t", sessionName+":0.0", "#{pane_current_command} #{pane_start_command}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(output)), "claude")
}

// getPaneCommand gets the command running in the first pane of a session
func getPaneCommand(sessionName string) string {
	cmd := exec.Command("tmux", "display-message", "-p", "-t", sessionName+":0.0", "#{pane_current_command}")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// getClaudeStateForSession retrieves Claude state for a tmux session
func getClaudeStateForSession(sessionName string, paneID string) (*ClaudeState, error) {
	// Try to find state file by tmux pane ID
	if paneID != "" && paneID != "none" {
		state, err := findStateByPane(paneID)
		if err == nil {
			return state, nil
		}
	}

	// Fallback: find state file by working directory
	// Get working directory from tmux pane
	cmd := exec.Command("tmux", "display-message", "-p", "-t", sessionName+":0.0", "#{pane_current_path}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	workingDir := strings.TrimSpace(string(output))
	return findStateByWorkingDir(workingDir)
}

// findStateByPane finds a state file by tmux pane ID
func findStateByPane(paneID string) (*ClaudeState, error) {
	files, err := filepath.Glob(filepath.Join(claudeStateDir, "*.json"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		state, err := readStateFile(file)
		if err != nil {
			continue
		}

		if state.TmuxPane == paneID {
			// Return state even if stale - let the display layer handle staleness
			return state, nil
		}
	}

	return nil, fmt.Errorf("no state found for pane %s", paneID)
}

// findStateByWorkingDir finds a state file by working directory
func findStateByWorkingDir(workingDir string) (*ClaudeState, error) {
	files, err := filepath.Glob(filepath.Join(claudeStateDir, "*.json"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		state, err := readStateFile(file)
		if err != nil {
			continue
		}

		// Only match states that are actually running in tmux
		// Skip states with tmux_pane="none" (Claude running outside tmux)
		if state.TmuxPane == "" || state.TmuxPane == "none" {
			continue
		}

		if state.WorkingDir == workingDir {
			// Return state even if stale - let the display layer handle staleness
			return state, nil
		}
	}

	return nil, fmt.Errorf("no state found for working dir %s", workingDir)
}

// readStateFile reads and parses a Claude state file
func readStateFile(path string) (*ClaudeState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var state ClaudeState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// isStateFresh checks if the state was updated recently
func isStateFresh(state *ClaudeState) bool {
	updated, err := time.Parse(time.RFC3339, state.LastUpdated)
	if err != nil {
		return false
	}

	age := time.Since(updated)
	return age < staleThreshold
}

// formatClaudeStatus returns a human-readable status string with icon
func formatClaudeStatus(state *ClaudeState) string {
	if state == nil {
		return "Unknown"
	}

	// Check if state is stale
	if !isStateFresh(state) {
		return "⚪ Stale (no updates)"
	}

	switch state.Status {
	case "idle":
		return "🟢 Idle"
	case "processing":
		return "🟡 Processing"
	case "tool_use":
		if state.CurrentTool != "" {
			// Try to get detailed info from args
			detail := extractToolDetail(state)
			if detail != "" {
				return fmt.Sprintf("🔧 %s: %s", state.CurrentTool, detail)
			}
			return fmt.Sprintf("🔧 Using %s", state.CurrentTool)
		}
		return "🔧 Using Tool"
	case "awaiting_input":
		return "⏸️  Awaiting Input"
	case "working":
		// Show what tool just finished if available
		if state.CurrentTool != "" {
			detail := extractToolDetail(state)
			if detail != "" {
				return fmt.Sprintf("⚙️  Processing %s: %s", state.CurrentTool, detail)
			}
			return fmt.Sprintf("⚙️  Processing %s", state.CurrentTool)
		}
		return "⚙️  Working"
	default:
		return fmt.Sprintf("❓ %s", state.Status)
	}
}

// extractToolDetail extracts relevant detail from tool args (file path, command, etc.)
func extractToolDetail(state *ClaudeState) string {
	if state.Details == nil {
		return ""
	}

	// Get args from details
	args, ok := state.Details["args"].(map[string]interface{})
	if !ok {
		return ""
	}

	// Extract based on tool type
	switch state.CurrentTool {
	case "Read", "Edit", "Write":
		// Show file path
		if filePath, ok := args["file_path"].(string); ok {
			// Show just the filename for brevity
			parts := strings.Split(filePath, "/")
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
			return filePath
		}

	case "Bash":
		// Show command (truncated)
		if command, ok := args["command"].(string); ok {
			// Truncate long commands
			maxLen := 40
			if len(command) > maxLen {
				return command[:maxLen] + "..."
			}
			return command
		}

	case "Grep", "Glob":
		// Show search pattern
		if pattern, ok := args["pattern"].(string); ok {
			maxLen := 30
			if len(pattern) > maxLen {
				return pattern[:maxLen] + "..."
			}
			return pattern
		}

	case "Task":
		// Show task description
		if description, ok := args["description"].(string); ok {
			return description
		}
	}

	return ""
}

// getClaudeStatusIcon returns just the icon for compact display
func getClaudeStatusIcon(state *ClaudeState) string {
	if state == nil {
		return "○"
	}

	// Check if state is stale
	if !isStateFresh(state) {
		return "⚪"
	}

	switch state.Status {
	case "idle":
		return "🟢"
	case "processing":
		return "🟡"
	case "tool_use":
		return "🔧"
	case "awaiting_input":
		return "⏸️"
	case "working":
		return "⚙️"
	default:
		return "❓"
	}
}

// formatTimeAgo returns a human-readable time duration
func formatTimeAgo(timestamp string) string {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return "unknown"
	}

	duration := time.Since(t)
	if duration < time.Second {
		return "just now"
	} else if duration < time.Minute {
		return fmt.Sprintf("%d seconds ago", int(duration.Seconds()))
	} else if duration < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	} else {
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	}
}
