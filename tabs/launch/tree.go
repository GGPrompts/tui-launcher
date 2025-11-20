package launch

// buildTreeFromConfig converts config into separate global and project tree structures
// Returns (globalItems, projectItems) for the two-pane layout
func buildTreeFromConfig(config Config) ([]launchItem, []launchItem) {
	var globalItems []launchItem
	var projectItems []launchItem

	// Projects go to right pane
	if len(config.Projects) > 0 {
		for _, proj := range config.Projects {
			item := launchItem{
				Name:     proj.Name,
				Path:     "projects/" + proj.Name,
				ItemType: typeCategory,
				Icon:     proj.Icon,
				Cwd:      expandPath(proj.Path), // Store project directory for CD
				Children: []launchItem{},
			}

			// Add commands
			for _, cmd := range proj.Commands {
				cmdItem := launchItem{
					Name:         cmd.Name,
					Path:         "projects/" + proj.Name + "/" + cmd.Name,
					ItemType:     typeCommand,
					Icon:         cmd.Icon,
					Command:      cmd.Command,
					Cwd:          expandPath(cmd.Cwd),
					SpawnStr:     cmd.Spawn,
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

			projectItems = append(projectItems, item)
		}
	}

	// Tools go to left pane (global)
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
					Name:         cmd.Name,
					Path:         "tools/" + cat.Category + "/" + cmd.Name,
					ItemType:     typeCommand,
					Icon:         cmd.Icon,
					Command:      cmd.Command,
					Cwd:          expandPath(cmd.Cwd),
					SpawnStr:     cmd.Spawn,
					DefaultSpawn: parseSpawnMode(cmd.Spawn),
				}
				item.Children = append(item.Children, cmdItem)
			}

			globalItems = append(globalItems, item)
		}
	}

	// Scripts go to left pane (global)
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
					Name:         cmd.Name,
					Path:         "scripts/" + cat.Category + "/" + cmd.Name,
					ItemType:     typeCommand,
					Icon:         cmd.Icon,
					Command:      cmd.Command,
					Cwd:          expandPath(cmd.Cwd),
					SpawnStr:     cmd.Spawn,
					DefaultSpawn: parseSpawnMode(cmd.Spawn),
				}
				item.Children = append(item.Children, cmdItem)
			}

			globalItems = append(globalItems, item)
		}
	}

	return globalItems, projectItems
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
