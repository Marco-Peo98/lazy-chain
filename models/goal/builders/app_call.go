package builders

import (
	goal "lazychain/models/goal/iface"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Docs: https://developer.algorand.org/docs/clis/goal/app/app/
type AppCallBuilder struct {
	status string
	RunWith func(argv []string)
}

func NewAppCallBuilder() *AppCallBuilder { return &AppCallBuilder{} }
func (b *AppCallBuilder) Title() string { return "App Call (goal app call/method)" }
func (b *AppCallBuilder) Init() tea.Cmd { return nil }
func (b *AppCallBuilder) Validate() error { return nil }
func (b *AppCallBuilder) Args() []string { return []string{"app","call","--app-id","0"} }
func (b *AppCallBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil { b.status = strings.TrimSpace(stderr); return }
	b.status = strings.TrimSpace(stdout)
}
func (b *AppCallBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) { return b, nil }
func (b *AppCallBuilder) View() string { return b.Title()+"\n\n"+b.status+"\n\n(coming soon)" }
