package iface

import tea "github.com/charmbracelet/bubbletea"

// RunFunc lets the host model execute a goal command with argv.
type RunFunc func(argv []string)

// Builder is the interface every TUI builder implements.
type Builder interface {
	// Identification
	Title() string
	TxnType() string

	// Lifecycle
	Init() tea.Cmd
	Update(tea.Msg) (Builder, tea.Cmd)

	// Rendering
	RenderFields() string // Fields panel content
	RenderOutput() string // Output panel content

	// Legacy kept for compatibility
	View() string

	// Validation and Execution
	Validate() error
	Args() []string
	AfterRun(stdout, stderr string, runErr error)
}
