package builders

import (
	"errors"
	"strings"

	goal "lazychain/models/goal/iface"

	tea "github.com/charmbracelet/bubbletea"
)

// Docs: https://developer.algorand.org/docs/clis/goal/asset/send/
type AssetTransferBuilder struct {
	// TODO: fields for: asset-id (-a), amount (-o? differs), sender (-f), receiver (-t), close-to, revocation-target, note, out, sign, etc.
	status string
	RunWith func(argv []string)
}

func NewAssetTransferBuilder() *AssetTransferBuilder { return &AssetTransferBuilder{} }
func (b *AssetTransferBuilder) Title() string { return "ASA Transfer (goal asset send)" }
func (b *AssetTransferBuilder) Init() tea.Cmd { return nil }
func (b *AssetTransferBuilder) Validate() error { return errors.New("not implemented yet") }
func (b *AssetTransferBuilder) Args() []string { return []string{"asset","send"} }
func (b *AssetTransferBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil { b.status = strings.TrimSpace(stderr); return }
	b.status = strings.TrimSpace(stdout)
}
func (b *AssetTransferBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) { return b, nil }
func (b *AssetTransferBuilder) View() string { return b.Title()+"\n\n"+b.status+"\n\n(coming soon)" }
