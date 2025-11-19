package main

import (
	"fmt"
	"os"
	"strings"
)

// buildTreeFromConfig converts config into a flat tree structure
func buildTreeFromConfig(config Config) []launchItem {
	var items []launchItem

	// Add projects
	if len(config.Projects) > 0 {
		for _, proj := range config.Projects {
			item := launchItem{
				Name:     proj.Name,
				Path:     "projects/" + proj.Name,
				ItemType: typeCategory,
				Icon:     proj.Icon,
				Children: []launchItem{},
			}

			// Add commands
			for _, cmd := range proj.Commands {
				cmdItem := launchItem{
					Name:     cmd.Name,
					Path:     "projects/" + proj.Name + "/" + cmd.Name,
					ItemType: typeCommand,
					Icon:     cmd.Icon,
					Command:  cmd.Command,
					Cwd:      expandPath(cmd.Cwd),
					SpawnStr: cmd.Spawn,
					DefaultSpawn: parseSpawnMode(cmd.Spawn),
				}
				item.Children = append(item.Children, cmdItem)
			}

			// Add profiles
			for _, prof := range proj.Profiles {
				profItem := launchItem{
					Name:      prof.Name,
					Path:      "projects/" + proj.Name + "/" + prof.Name,
					ItemType:  typeProfile,
					Icon:      prof.Icon,
					IsProfile: true,
					LayoutStr: prof.Layout,
					Layout:    parseLayoutMode(prof.Layout),
					Panes:     prof.Panes,
				}
				item.Children = append(item.Children, profItem)
			}

			items = append(items, item)
		}
	}

	// Add tools
	if len(config.Tools) > 0 {
		for _, cat := range config.Tools {
			item := launchItem{
				Name:     cat.Category,
				Path:     "tools/" + cat.Category,
				ItemType: typeCategory,
				Icon:     cat.Icon,
				Children: []launchItem{},
			}

			for _, cmd := range cat.Items {
				cmdItem := launchItem{
					Name:     cmd.Name,
					Path:     "tools/" + cat.Category + "/" + cmd.Name,
					ItemType: typeCommand,
					Icon:     cmd.Icon,
					Command:  cmd.Command,
					Cwd:      expandPath(cmd.Cwd),
					SpawnStr: cmd.Spawn,
					DefaultSpawn: parseSpawnMode(cmd.Spawn),
				}
				item.Children = append(item.Children, cmdItem)
			}

			items = append(items, item)
		}
	}

	// Add scripts
	if len(config.Scripts) > 0 {
		for _, cat := range config.Scripts {
			item := launchItem{
				Name:     cat.Category,
				Path:     "scripts/" + cat.Category,
				ItemType: typeCategory,
				Icon:     cat.Icon,
				Children: []launchItem{},
			}

			for _, cmd := range cat.Items {
				cmdItem := launchItem{
					Name:     cmd.Name,
					Path:     "scripts/" + cat.Category + "/" + cmd.Name,
					ItemType: typeCommand,
					Icon:     cmd.Icon,
					Command:  cmd.Command,
					Cwd:      expandPath(cmd.Cwd),
					SpawnStr: cmd.Spawn,
					DefaultSpawn: parseSpawnMode(cmd.Spawn),
				}
				item.Children = append(item.Children, cmdItem)
			}

			items = append(items, item)
		}
	}

	return items
}

// flattenTree converts hierarchical items into a flat list for display
func flattenTree(items []launchItem, expandedItems map[string]bool) []launchTreeItem {
	var result []launchTreeItem
	for i, item := range items {
		isLast := i == len(items)-1
		flattenTreeRecursive(item, 0, isLast, []bool{}, expandedItems, &result)
	}
	return result
}

// flattenTreeRecursive recursively flattens the tree
func flattenTreeRecursive(item launchItem, depth int, isLast bool, parentLasts []bool, expandedItems map[string]bool, result *[]launchTreeItem) {
	treeItem := launchTreeItem{
		item:        item,
		depth:       depth,
		isLast:      isLast,
		parentLasts: parentLasts,
	}
	*result = append(*result, treeItem)

	// If expanded and has children, add children
	if expandedItems[item.Path] && len(item.Children) > 0 {
		newParentLasts := append([]bool{}, parentLasts...)
		newParentLasts = append(newParentLasts, isLast)

		for i, child := range item.Children {
			childIsLast := i == len(item.Children)-1
			flattenTreeRecursive(child, depth+1, childIsLast, newParentLasts, expandedItems, result)
		}
	}
}

// parseSpawnMode converts spawn string to spawnMode
func parseSpawnMode(spawn string) spawnMode {
	switch spawn {
	case "xterm-window":
		return spawnXtermWindow
	case "tmux-window":
		return spawnTmuxWindow
	case "tmux-split-h":
		return spawnTmuxSplitH
	case "tmux-split-v":
		return spawnTmuxSplitV
	case "tmux-layout":
		return spawnTmuxLayout
	case "current-pane":
		return spawnCurrentPane
	default:
		return spawnTmuxWindow
	}
}

// parseLayoutMode converts layout string to tmuxLayout
func parseLayoutMode(layout string) tmuxLayout {
	switch layout {
	case "main-vertical":
		return layoutMainVertical
	case "main-horizontal":
		return layoutMainHorizontal
	case "tiled":
		return layoutTiled
	case "even-horizontal":
		return layoutEvenHorizontal
	case "even-vertical":
		return layoutEvenVertical
	default:
		return layoutTiled
	}
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return strings.Replace(path, "~", homeDir, 1)
	}
	return path
}

// renderTreeItem renders a single tree item with proper indentation
func renderTreeItem(ti launchTreeItem, cursor int, index int, selected bool, expanded bool) string {
	var sb strings.Builder

	// Cursor indicator
	if index == cursor {
		sb.WriteString("> ")
	} else {
		sb.WriteString("  ")
	}

	// Tree lines
	for i, isLast := range ti.parentLasts {
		if i == len(ti.parentLasts)-1 {
			continue
		}
		if isLast {
			sb.WriteString("  ")
		} else {
			sb.WriteString("│ ")
		}
	}

	// Branch character
	if ti.depth > 0 {
		if ti.isLast {
			sb.WriteString("└─")
		} else {
			sb.WriteString("├─")
		}
	}

	// Selection checkbox
	if selected {
		sb.WriteString(emojiSelected + " ")
	} else if ti.item.ItemType == typeCommand || ti.item.ItemType == typeProfile {
		sb.WriteString(emojiUnselected + " ")
	}

	// Expansion indicator for categories
	if ti.item.ItemType == typeCategory {
		if expanded {
			sb.WriteString(emojiExpanded + " ")
		} else {
			sb.WriteString(emojiCollapsed + " ")
		}
	}

	// Icon
	if ti.item.Icon != "" {
		sb.WriteString(ti.item.Icon + " ")
	}

	// Name
	sb.WriteString(ti.item.Name)

	// Type indicator (only show for profiles)
	switch ti.item.ItemType {
	case typeProfile:
		sb.WriteString(fmt.Sprintf(" [%s]", ti.item.LayoutStr))
	}

	return sb.String()
}
