package goal

import (
	"lazychain/models/goal/iface"
)

// Builder is a screen that configures a specific `goal` command.
/* type Builder interface {
	Title() string
	Init() tea.Cmd
	Update(tea.Msg) (Builder, tea.Cmd)
	View() string

	// Validate returns nil if the form is ready to run.
	Validate() error

	// Args returns the full argv for `goal` (e.g., ["clerk","send","-a","123",...])
	Args() []string

	// AfterRun is called with stdout/stderr + error after execution.
	AfterRun(stdout, stderr string, runErr error)
} */


type Builder = iface.Builder
type RunFunc = iface.RunFunc