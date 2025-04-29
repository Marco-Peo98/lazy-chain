package models

type SessionState int

const (
	MainView SessionState = iota
	ProjectView
	SettingsView
	ApplicationsView
	CmdGoalsView
	ExploreView
)
