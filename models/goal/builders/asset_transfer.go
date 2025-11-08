package builders

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	goal "lazychain/models/goal/iface"
)

type AssetTransferBuilder struct {
	assetID  string
	sender   string
	receiver string
	amount   string

	feeAuto         bool
	validRoundsAuto bool
	note            string

	status  string
	runFunc goal.RunFunc
}

func NewAssetTransferBuilder() *AssetTransferBuilder {
	return &AssetTransferBuilder{
		feeAuto:         true,
		validRoundsAuto: true,
	}
}

func (b *AssetTransferBuilder) SetRunFunc(fn goal.RunFunc) { b.runFunc = fn }
func (b *AssetTransferBuilder) Title() string              { return "Asset Transfer (axfer)" }
func (b *AssetTransferBuilder) TxnType() string            { return "axfer" }
func (b *AssetTransferBuilder) Init() tea.Cmd              { return nil }

func (b *AssetTransferBuilder) Validate() error {
	if strings.TrimSpace(b.assetID) == "" {
		return errors.New("asset ID required")
	}
	if strings.TrimSpace(b.sender) == "" {
		return errors.New("sender required")
	}
	if strings.TrimSpace(b.receiver) == "" {
		return errors.New("receiver required")
	}
	if strings.TrimSpace(b.amount) == "" {
		return errors.New("amount required")
	}
	return nil
}

func (b *AssetTransferBuilder) Args() []string {
	return []string{"asset", "send", "--assetid", b.assetID, "-f", b.sender, "-t", b.receiver, "-a", b.amount}
}

func (b *AssetTransferBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil {
		b.status = fmt.Sprintf("Error: %v\n%s", runErr, stderr)
		return
	}
	b.status = stdout
}

func (b *AssetTransferBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "enter" {
		if err := b.Validate(); err != nil {
			b.status = "Validation Error: " + err.Error()
		} else if b.runFunc != nil {
			b.runFunc(b.Args())
		}
	}
	return b, nil
}

func (b *AssetTransferBuilder) RenderFields() string {
	var l []string
	l = append(l, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#cba6f7")).Render("Asset Transfer"), "")
	l = append(l, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#89b4fa")).Render("Required:"), "")
	l = append(l, "  Asset ID:", "    "+getPlaceholder(b.assetID, "[asset id]"), "")
	l = append(l, "  Sender:", "    "+getPlaceholder(b.sender, "[address]"), "")
	l = append(l, "  Receiver:", "    "+getPlaceholder(b.receiver, "[address]"), "")
	l = append(l, "  Amount:", "    "+getPlaceholder(b.amount, "[amount]"), "")
	l = append(l, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f9e2af")).Render("Optional:"), "")
	l = append(l, "  Fee auto: "+getBoolStr(b.feeAuto), "")
	l = append(l, "  Valid rounds auto: "+getBoolStr(b.validRoundsAuto), "")
	l = append(l, "  Note:", "    "+getPlaceholder(b.note, "[optional]"), "")
	return strings.Join(l, "\n")
}

func (b *AssetTransferBuilder) RenderOutput() string { return renderOutput(b.status) }

// View returns the complete view (delegates to RenderFields for compatibility)
func (b *AssetTransferBuilder) View() string {
	return b.RenderFields()
}
