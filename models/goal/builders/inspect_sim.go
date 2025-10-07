package builders

import (
	"strings"

	goal "lazychain/models/goal/iface"

	tea "github.com/charmbracelet/bubbletea"
)

// Docs:
// inspect:  https://developer.algorand.org/docs/clis/goal/clerk/inspect/
// simulate: https://developer.algorand.org/docs/clis/goal/clerk/simulate/
type InspectSimBuilder struct {
	status string
	RunWith func(argv []string)
}

func NewInspectSimBuilder() *InspectSimBuilder { return &InspectSimBuilder{} }
func (b *InspectSimBuilder) Title() string { return "Inspect / Simulate" }
func (b *InspectSimBuilder) Init() tea.Cmd { return nil }
func (b *InspectSimBuilder) Validate() error { return nil }
func (b *InspectSimBuilder) Args() []string { return []string{"clerk","inspect","out.stxn"} }
func (b *InspectSimBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil { b.status = strings.TrimSpace(stderr); return }
	b.status = strings.TrimSpace(stdout)
}
func (b *InspectSimBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) { return b, nil }
func (b *InspectSimBuilder) View() string { return b.Title()+"\n\n"+b.status+"\n\n(coming soon)" }
