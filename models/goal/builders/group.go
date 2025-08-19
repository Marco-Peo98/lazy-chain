package builders

import (
	"strings"

	goal "lazychain/models/goal/iface"

	tea "github.com/charmbracelet/bubbletea"
)

// Docs: https://developer.algorand.org/docs/clis/goal/clerk/group/
type GroupBuilder struct {
	status string
	RunWith func(argv []string)
}

func NewGroupBuilder() *GroupBuilder { return &GroupBuilder{} }
func (b *GroupBuilder) Title() string { return "Atomic Group (goal clerk group)" }
func (b *GroupBuilder) Init() tea.Cmd { return nil }
func (b *GroupBuilder) Validate() error { return nil }
func (b *GroupBuilder) Args() []string { return []string{"clerk","group","-i","group.json","-o","group.txn"} }
func (b *GroupBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil { b.status = strings.TrimSpace(stderr); return }
	b.status = strings.TrimSpace(stdout)
}
func (b *GroupBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) { return b, nil }
func (b *GroupBuilder) View() string { return b.Title()+"\n\n"+b.status+"\n\n(coming soon)" }
