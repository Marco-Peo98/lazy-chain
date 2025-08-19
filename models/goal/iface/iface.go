package iface

import tea "github.com/charmbracelet/bubbletea"

// RunFunc lets the host model execute a goal command with argv.
type RunFunc func(argv []string)

// Builder is the interface every TUI builder implements.
type Builder interface {
	Title() string
	Init() tea.Cmd
	Update(tea.Msg) (Builder, tea.Cmd)
	View() string

	Validate() error
	Args() []string
	AfterRun(stdout, stderr string, runErr error)
}
