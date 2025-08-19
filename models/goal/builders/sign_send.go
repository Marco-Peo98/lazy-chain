package builders

import (
	"strings"

	goal "lazychain/models/goal/iface"

	tea "github.com/charmbracelet/bubbletea"
)

// Docs:
// sign:   https://developer.algorand.org/docs/clis/goal/clerk/sign/
// rawsend:https://developer.algorand.org/docs/clis/goal/clerk/rawsend/
type SignSendBuilder struct {
	status  string
	RunWith func(argv []string)
}

func NewSignSendBuilder() *SignSendBuilder { return &SignSendBuilder{} }
func (b *SignSendBuilder) Title() string { return "Sign / Send (goal clerk sign/rawsend)" }
func (b *SignSendBuilder) Init() tea.Cmd { return nil }
func (b *SignSendBuilder) Validate() error { return nil }
func (b *SignSendBuilder) Args() []string { return []string{"clerk","sign","-i","in.txn","-o","out.stxn"} }
func (b *SignSendBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil { b.status = strings.TrimSpace(stderr); return }
	b.status = strings.TrimSpace(stdout)
}
func (b *SignSendBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) { return b, nil }
func (b *SignSendBuilder) View() string { return b.Title()+"\n\n"+b.status+"\n\n(coming soon)" }
