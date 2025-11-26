package iface

import (
	tea "github.com/charmbracelet/bubbletea"
)

// RunFunc is the type signature for the function that actually runs the command.
type RunFunc func(args []string)

// Builder is the interface for a transaction type builder.
type Builder interface {
	// Transaction metadata
	Title() string
	TxnType() string

	// Lifecycle
	Init() tea.Cmd
	Update(msg tea.Msg) (Builder, tea.Cmd)
	Validate() error
	Args() []string
	AfterRun(stdout, stderr string, runErr error)

	// State
	SetRunFunc(fn RunFunc)
	IsEditing() bool

	// Layout integration
	SetAvailableHeight(lines int) // Sets available height in terminal lines for viewport calculation

	// Rendering
	RenderFields() string
	RenderOutput() string
	View() string
}
