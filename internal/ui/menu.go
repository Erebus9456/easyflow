package ui

// MainMenuItem maps an executable item inside our dashboard choice menu
type MainMenuItem struct {
	Title       string
	Description string
}

// GetMainMenuOptions yields the structured collection of core operational tracking steps
func GetMainMenuOptions() []MainMenuItem {
	return []MainMenuItem{
		{
			Title:       "🐛 Manage Issues Menu",
			Description: "Access CRUD: List open issues, create trackers, or close existing items",
		},
		{
			Title:       "🌿 Manage Branches Menu",
			Description: "Access CRUD: Checkout existing, create custom, or drop local branches",
		},
		{
			Title:       "💾 Manage Commits Menu",
			Description: "Access CRUD: View mini git log history, stage everything, or undo check-ins",
		},
		{
			Title:       "🚀 Start Pipeline Work Loop",
			Description: "Automated loop: Select Issue ➔ Branch ➔ Commit ➔ Push ➔ PR ➔ Merge",
		},
		{
			Title:       "Stage & Commit Local Modifications",
			Description: "Runs 'git add .' and prompts for a message checkpoint locally",
		},
		{
			Title:       "Sync Tracked Upstream Modifications",
			Description: "Pushes the current active local tracking context branch directly to origin",
		},
		{
			Title:       "Reset Context State Engine",
			Description: "Clears current operational memory flags back to factory defaults safely",
		},
	}
}

// GetSubMenuOptions yields localized sub-choices for fine-grained management control
func GetSubMenuOptions(category string) []MainMenuItem {
	switch category {
	case "issues":
		return []MainMenuItem{
			{Title: "List Repository Issues", Description: "Browse open tracked issues on your remote GitHub repository"},
			{Title: "Create Tracker Issue", Description: "Open a brand new ticket via standard title text inputs"},
			{Title: "Close Issue by Number", Description: "Settle and flag an item as closed using its identification number"},
		}
	case "branches":
		return []MainMenuItem{
			{Title: "Select / Checkout Existing Branch", Description: "Interactively browse and jump onto an existing workspace branch"},
			{Title: "Create Custom Local Branch", Description: "Generate and auto-checkout a brand new localized working branch"},
			{Title: "Delete Local Working Branch", Description: "Safely drop a chosen workspace branch using underlying 'git branch -d'"},
		}
	case "commits":
		return []MainMenuItem{
			{Title: "View Recent Commit Log", Description: "Inspect the last 5 local commit messages checked into this branch"},
			{Title: "Stage & Commit Modifications", Description: "Instantly index all modified tracked files and record a commit message"},
			{Title: "Undo Last Local Commit", Description: "Execute a soft reset back one revision to unstage files into your workspace"},
		}
	default:
		return []MainMenuItem{}
	}
}
