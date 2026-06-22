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
			Title:       "Start Pipeline Work Loop",
			Description: "Step-by-step pipeline: Select Issue ➔ Branch ➔ Commit ➔ Push ➔ PR ➔ Merge",
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
